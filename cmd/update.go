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

// Fonction modifiée pour gérer le cas spécifique de Windows
func doUpdate(version string, force bool) error {
	// Si aucune version n'est fournie, obtenir la dernière version
	var latestVersion string
	var err error
	if version == "" {
		latestVersion, err = getLatestVersion()
		if err != nil {
			return fmt.Errorf("error checking for updates: %w", err)
		}
		version = latestVersion
	}

	// Vérifier si une mise à jour est nécessaire
	currentVersion := Version
	if currentVersion == version && !force {
		fmt.Printf("You are already using the latest version of RLAMA (%s)\n", currentVersion)
		return nil
	}

	// Demander confirmation, sauf si --force est utilisé
	if !force {
		fmt.Printf("A new version of RLAMA is available (%s). Do you want to install it? (y/n): ", version)
		var response string
		_, err := fmt.Scanln(&response)
		if err != nil || (strings.ToLower(response) != "y" && strings.ToLower(response) != "yes") {
			fmt.Println("Update cancelled.")
			return nil
		}
	}

	fmt.Printf("Installing RLAMA %s...\n", version)

	// Obtenir le chemin de l'exécutable actuel
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error getting executable path: %w", err)
	}

	// Créer un répertoire pour les fichiers de mise à jour si nécessaire
	updateDir := filepath.Dir(execPath)
	if err := os.MkdirAll(updateDir, 0755); err != nil {
		return fmt.Errorf("error creating update directory: %w", err)
	}

	// Télécharger la nouvelle version
	binaryURL := fmt.Sprintf("https://github.com/dontizi/rlama/releases/download/v%s/rlama_%s_%s", version, runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		binaryURL += ".exe"
	}

	// Chemin pour le nouveau binaire
	newBinaryPath := execPath + ".new"

	// Télécharger le nouveau binaire
	if err := downloadFile(binaryURL, newBinaryPath); err != nil {
		// Nettoyer en cas d'erreur
		os.Remove(newBinaryPath)
		return fmt.Errorf("error downloading update: %w", err)
	}

	// Rendre le nouveau binaire exécutable
	if err := os.Chmod(newBinaryPath, 0755); err != nil {
		os.Remove(newBinaryPath)
		return fmt.Errorf("error setting permissions: %w", err)
	}

	// Sous Windows, nous devons utiliser une approche différente
	if runtime.GOOS == "windows" {
		return windowsReplaceBinary(execPath, newBinaryPath)
	}

	// Sur les autres plateformes, nous pouvons remplacer directement
	if err := os.Rename(newBinaryPath, execPath); err != nil {
		os.Remove(newBinaryPath)
		return fmt.Errorf("error replacing binary: %w", err)
	}

	fmt.Printf("Successfully updated to RLAMA %s!\n", version)
	return nil
}

// Nouvelle fonction pour gérer la mise à jour sous Windows
func windowsReplaceBinary(originalPath, newPath string) error {
	// Créer un script batch pour le remplacement différé
	batchContent := `@echo off
:wait
timeout /t 1 >nul
tasklist /fi "imagename eq rlama.exe" | find "rlama.exe" >nul
if %errorlevel% equ 0 goto wait
move /y "%s" "%s"
echo Update successful!
start "" "%s"
exit
`
	batchScript := fmt.Sprintf(batchContent, newPath, originalPath, originalPath)
	
	// Créer un fichier temporaire pour le script batch
	tempBatchFile := filepath.Join(os.TempDir(), "rlama_update.bat")
	if err := os.WriteFile(tempBatchFile, []byte(batchScript), 0644); err != nil {
		return fmt.Errorf("error creating update script: %w", err)
	}
	
	// Exécuter le script en arrière-plan
	cmd := exec.Command("cmd", "/c", "start", "/min", tempBatchFile)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting update process: %w", err)
	}
	
	fmt.Println("Update will complete after you exit the RLAMA application.")
	fmt.Println("Please close this window and run 'rlama --version' to verify the update.")
	
	return nil
} 