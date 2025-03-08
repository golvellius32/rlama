package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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
	Short: "Check and install RLAMA updates",
	Long: `Check if a new version of RLAMA is available and install it if so.
Example: rlama update

By default, the command asks for confirmation before installing the update.
Use the --force flag to update without confirmation.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Checking for RLAMA updates...")
		
		// Check the latest available version
		latestRelease, hasUpdates, err := checkForUpdates()
		if err != nil {
			return fmt.Errorf("error checking for updates: %w", err)
		}
		
		if !hasUpdates {
			fmt.Printf("You are already using the latest version of RLAMA (%s).\n", Version)
			return nil
		}
		
		latestVersion := strings.TrimPrefix(latestRelease.TagName, "v")
		
		// Ask for confirmation unless --force is specified
		if !forceUpdate {
			fmt.Printf("A new version of RLAMA is available (%s). Do you want to install it? (y/n): ", latestVersion)
			var response string
			fmt.Scanln(&response)
			
			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" && response != "yes" {
				fmt.Println("Update cancelled.")
				return nil
			}
		}
		
		fmt.Printf("Installing RLAMA %s...\n", latestVersion)
		
		// Determine which binary to download based on OS and architecture
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
			return fmt.Errorf("no binary found for your system (%s_%s)", osName, archName)
		}
		
		// Download the binary
		execPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("unable to determine executable location: %w", err)
		}
		
		// Create a temporary file for the download
		tempFile := execPath + ".new"
		out, err := os.Create(tempFile)
		if err != nil {
			return fmt.Errorf("error creating temporary file: %w", err)
		}
		defer out.Close()
		
		// Download the binary
		resp, err := http.Get(assetURL)
		if err != nil {
			return fmt.Errorf("download error: %w", err)
		}
		defer resp.Body.Close()
		
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return fmt.Errorf("error writing file: %w", err)
		}
		
		// Make the binary executable
		err = os.Chmod(tempFile, 0755)
		if err != nil {
			return fmt.Errorf("error setting permissions: %w", err)
		}
		
		// Replace the old binary with the new one
		backupPath := execPath + ".bak"
		os.Rename(execPath, backupPath) // Backup the old binary
		err = os.Rename(tempFile, execPath)
		if err != nil {
			// In case of error, restore the old binary
			os.Rename(backupPath, execPath)
			return fmt.Errorf("error replacing binary: %w", err)
		}
		
		fmt.Printf("RLAMA has been updated to version %s.\n", latestVersion)
		return nil
	},
}

// checkForUpdates checks if updates are available by querying the GitHub API
func checkForUpdates() (*GitHubRelease, bool, error) {
	// Query the GitHub API to get the latest release
	resp, err := http.Get("https://api.github.com/repos/dontizi/rlama/releases/latest")
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()
	
	// Parse the JSON response
	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, false, err
	}
	
	// Check if the version is newer
	latestVersion := strings.TrimPrefix(release.TagName, "v")
	hasUpdates := latestVersion != Version
	
	return &release, hasUpdates, nil
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVarP(&forceUpdate, "force", "f", false, "Update without asking for confirmation")
} 