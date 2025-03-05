package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var forceUninstall bool

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Désinstalle RLAMA et tous ses fichiers",
	Long:  `Désinstalle complètement RLAMA en supprimant l'exécutable et tous les fichiers de données associés.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. Vérifier si l'utilisateur a confirmé la suppression
		if !forceUninstall {
			fmt.Print("Cette action va supprimer RLAMA et toutes vos données. Êtes-vous sûr ? (o/n): ")
			var response string
			fmt.Scanln(&response)
			
			response = strings.ToLower(strings.TrimSpace(response))
			if response != "o" && response != "oui" {
				fmt.Println("Désinstallation annulée.")
				return nil
			}
		}

		// 2. Supprimer le répertoire de données
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("impossible de déterminer le répertoire utilisateur: %w", err)
		}
		
		dataDir := filepath.Join(homeDir, ".rlama")
		fmt.Printf("Suppression du répertoire de données: %s\n", dataDir)
		
		if _, err := os.Stat(dataDir); err == nil {
			err = os.RemoveAll(dataDir)
			if err != nil {
				return fmt.Errorf("impossible de supprimer le répertoire de données: %w", err)
			}
			fmt.Println("✓ Répertoire de données supprimé")
		} else {
			fmt.Println("Le répertoire de données n'existe pas ou a déjà été supprimé")
		}

		// 3. Supprimer l'exécutable
		executablePath := "/usr/local/bin/rlama"
		fmt.Printf("Suppression de l'exécutable: %s\n", executablePath)
		
		if _, err := os.Stat(executablePath); err == nil {
			// Sous macOS et Linux, nous avons probablement besoin de sudo
			var err error
			if os.Geteuid() == 0 {
				// Si nous sommes déjà root
				err = os.Remove(executablePath)
			} else {
				fmt.Println("Vous devrez peut-être entrer votre mot de passe pour supprimer l'exécutable")
				err = execCommand("sudo", "rm", executablePath)
			}
			
			if err != nil {
				return fmt.Errorf("impossible de supprimer l'exécutable: %w", err)
			}
			fmt.Println("✓ Exécutable supprimé")
		} else {
			fmt.Println("L'exécutable n'existe pas ou a déjà été supprimé")
		}

		fmt.Println("\nRLAMA a été désinstallé avec succès.")
		return nil
	},
}

// execCommand exécute une commande système
func execCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
	uninstallCmd.Flags().BoolVarP(&forceUninstall, "force", "f", false, "Désinstaller sans demander de confirmation")
} 