package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"runtime"

	"github.com/spf13/cobra"
)

var forceUninstall bool

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall RLAMA and all its files",
	Long:  `Completely uninstall RLAMA by removing the executable and all associated data files.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 1. Check if the user confirmed the deletion
		if !forceUninstall {
			fmt.Print("This action will remove RLAMA and all your data. Are you sure? (y/n): ")
			var response string
			fmt.Scanln(&response)
			
			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" && response != "yes" {
				fmt.Println("Uninstallation cancelled.")
				return nil
			}
		}

		// 2. Delete the data directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("unable to determine user directory: %w", err)
		}
		
		dataDir := filepath.Join(homeDir, ".rlama")
		fmt.Printf("Removing data directory: %s\n", dataDir)
		
		if _, err := os.Stat(dataDir); err == nil {
			err = os.RemoveAll(dataDir)
			if err != nil {
				return fmt.Errorf("unable to remove data directory: %w", err)
			}
			fmt.Println("✓ Data directory removed")
		} else {
			fmt.Println("Data directory doesn't exist or has already been removed")
		}

		// 3. Remove the executable
		var executablePath string
		if runtime.GOOS == "windows" {
			// Try different locations where it might be installed
			localAppData := os.Getenv("LOCALAPPDATA")
			possiblePaths := []string{
				filepath.Join(localAppData, "RLAMA", "rlama.exe"),
				filepath.Join(os.Getenv("ProgramFiles"), "RLAMA", "rlama.exe"),
				filepath.Join(os.Getenv("ProgramFiles(x86)"), "RLAMA", "rlama.exe"),
				filepath.Join(homeDir, "AppData", "Local", "RLAMA", "rlama.exe"),
			}
			
			for _, path := range possiblePaths {
				if _, err := os.Stat(path); err == nil {
					executablePath = path
					break
				}
			}
		} else {
			executablePath = "/usr/local/bin/rlama"
		}
		
		fmt.Printf("Removing executable: %s\n", executablePath)
		
		if _, err := os.Stat(executablePath); err == nil {
			// On macOS and Linux, we probably need sudo
			var err error
			if os.Geteuid() == 0 {
				// If we're already root
				err = os.Remove(executablePath)
			} else {
				fmt.Println("You may need to enter your password to remove the executable")
				err = execCommand("sudo", "rm", executablePath)
			}
			
			if err != nil {
				return fmt.Errorf("unable to remove executable: %w", err)
			}
			fmt.Println("✓ Executable removed")
		} else {
			fmt.Println("Executable doesn't exist or has already been removed")
		}

		fmt.Println("\nRLAMA has been successfully uninstalled.")
		return nil
	},
}

// execCommand executes a system command
func execCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
	uninstallCmd.Flags().BoolVarP(&forceUninstall, "force", "f", false, "Uninstall without asking for confirmation")
} 