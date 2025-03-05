package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var forceUpdate bool

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Vérifie et installe les mises à jour de RLAMA",
	Long: `Vérifie si une nouvelle version de RLAMA est disponible et l'installe si c'est le cas.
Exemple: rlama update

Par défaut, la commande demande une confirmation avant d'installer la mise à jour.
Utilisez le flag --force pour mettre à jour sans confirmation.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Vérification des mises à jour de RLAMA...")
		
		// Vérifier la dernière version disponible
		latestRelease, hasUpdates, err := checkForUpdates()
		if err != nil {
			return fmt.Errorf("erreur lors de la vérification des mises à jour: %w", err)
		}
		
		if !hasUpdates {
			fmt.Printf("Vous utilisez déjà la dernière version de RLAMA (%s).\n", Version)
			return nil
		}
		
		latestVersion := strings.TrimPrefix(latestRelease.TagName, "v")
		
		// Demander confirmation sauf si --force est spécifié
		if !forceUpdate {
			fmt.Printf("Une nouvelle version de RLAMA est disponible (%s). Voulez-vous l'installer? (o/n): ", latestVersion)
			var response string
			fmt.Scanln(&response)
			
			response = strings.ToLower(strings.TrimSpace(response))
			if response != "o" && response != "oui" {
				fmt.Println("Mise à jour annulée.")
				return nil
			}
		}
		
		fmt.Printf("Installation de RLAMA %s...\n", latestVersion)
		
		// Déterminer le binaire à télécharger en fonction du système d'exploitation et de l'architecture
		var assetURL string
		osName := runtime.GOOS
		archName := runtime.GOARCH
		assetPattern := fmt.Sprintf("rlama_%s_%s", osName, archName)
		
		for _, asset := range latestRelease.Assets {
			if strings.Contains(asset.Name, assetPattern) {
				assetURL = asset.BrowserDownloadURL
				break
			}
		}
		
		if assetURL == "" {
			return fmt.Errorf("aucun binaire trouvé pour votre système (%s_%s)", osName, archName)
		}
		
		// Télécharger le binaire
		execPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("impossible de déterminer l'emplacement de l'exécutable: %w", err)
		}
		
		// Créer un fichier temporaire pour le téléchargement
		tempFile := execPath + ".new"
		out, err := os.Create(tempFile)
		if err != nil {
			return fmt.Errorf("erreur lors de la création du fichier temporaire: %w", err)
		}
		defer out.Close()
		
		// Télécharger le binaire
		resp, err := http.Get(assetURL)
		if err != nil {
			return fmt.Errorf("erreur lors du téléchargement: %w", err)
		}
		defer resp.Body.Close()
		
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return fmt.Errorf("erreur lors de l'écriture du fichier: %w", err)
		}
		
		// Rendre le binaire exécutable
		err = os.Chmod(tempFile, 0755)
		if err != nil {
			return fmt.Errorf("erreur lors de la définition des permissions: %w", err)
		}
		
		// Remplacer l'ancien binaire par le nouveau
		backupPath := execPath + ".bak"
		os.Rename(execPath, backupPath) // Sauvegarde de l'ancien binaire
		err = os.Rename(tempFile, execPath)
		if err != nil {
			// En cas d'erreur, restaurer l'ancien binaire
			os.Rename(backupPath, execPath)
			return fmt.Errorf("erreur lors du remplacement du binaire: %w", err)
		}
		
		fmt.Printf("RLAMA a été mis à jour vers la version %s.\n", latestVersion)
		return nil
	},
}

// checkForUpdates vérifie si des mises à jour sont disponibles en interrogeant l'API GitHub
func checkForUpdates() (*GitHubRelease, bool, error) {
	// Interroger l'API GitHub pour obtenir la dernière release
	resp, err := http.Get("https://api.github.com/repos/dontizi/rlama/releases/latest")
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()
	
	// Analyser la réponse JSON
	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, false, err
	}
	
	// Vérifier si la version est plus récente
	latestVersion := strings.TrimPrefix(release.TagName, "v")
	hasUpdates := latestVersion != Version
	
	return &release, hasUpdates, nil
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVarP(&forceUpdate, "force", "f", false, "Mettre à jour sans demander de confirmation")
} 