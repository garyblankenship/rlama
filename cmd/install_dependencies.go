package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/dontizi/rlama/internal/utils"
	"github.com/spf13/cobra"
)

var installDependenciesCmd = &cobra.Command{
	Use:   "install-dependencies",
	Short: "Install necessary dependencies for RLAMA",
	Long:  `Install system and Python dependencies for optimal RLAMA performance, including the BGE reranker.`,
	Run: func(cmd *cobra.Command, args []string) {
		installDependencies()
	},
}

func init() {
	rootCmd.AddCommand(installDependenciesCmd)
}

func installDependencies() {
	fmt.Println("Installing dependencies...")

	// Find the installation script path
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("⚠️ Error determining executable path: %v\n", err)

		// Use an alternative solution
		installDependenciesFallback()
		return
	}

	// The scripts directory is presumed to be in the same directory as the executable
	// or in the parent directory for development environments
	scriptDir := filepath.Join(filepath.Dir(execPath), "scripts")
	scriptPath := filepath.Join(scriptDir, "install_deps.sh")

	// Check if the script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		// Try in the parent directory (for development)
		scriptDir = filepath.Join(filepath.Dir(execPath), "..", "scripts")
		scriptPath = filepath.Join(scriptDir, "install_deps.sh")

		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			fmt.Println("⚠️ Installation script not found. Using alternative method.")
			installDependenciesFallback()
			return
		}
	}

	// Execute the script
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// On Windows, use bash.exe (WSL or Git Bash)
		cmd = exec.Command("bash.exe", scriptPath)
	} else {
		// On Unix-like, execute the script directly
		cmd = exec.Command("bash", scriptPath)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠️ Error executing installation script: %v\n", err)
		fmt.Println("Attempting direct installation of Python dependencies...")
		installPythonDependencies()
	}
}

// This function is used if the install_deps.sh script is not found
func installDependenciesFallback() {
	fmt.Println("Installing necessary Python dependencies...")
	installPythonDependencies()
}

func installPythonDependencies() {
	pythonExecutor := utils.NewPythonExecutor()

	// Create virtual environment
	fmt.Println("Setting up Python virtual environment...")
	if err := pythonExecutor.CreateVirtualEnvironment(); err != nil {
		fmt.Printf("⚠️ Error creating virtual environment: %v\n", err)
		fmt.Println("Falling back to system-wide installation...")
		installSystemWide()
		return
	}

	// Install required packages in virtual environment
	packages := []string{"FlagEmbedding==1.2.10", "torch", "transformers", "pdfminer.six", "docx2txt", "xlsx2csv"}

	for _, pkg := range packages {
		if !pythonExecutor.CheckPackageInstalled(pkg) {
			if err := pythonExecutor.InstallPackage(pkg); err != nil {
				fmt.Printf("⚠️ Error installing %s: %v\n", pkg, err)
				continue
			}
		} else {
			fmt.Printf("✅ %s is already installed\n", pkg)
		}
	}

	fmt.Println("✅ All Python dependencies installed successfully in virtual environment!")
	fmt.Printf("Virtual environment location: %s\n", pythonExecutor.GetVirtualEnvPath())
}

// installSystemWide tries to install packages system-wide as a fallback
func installSystemWide() {
	pythonCmd := "python3"
	if _, err := exec.LookPath(pythonCmd); err != nil {
		pythonCmd = "python"
		if _, err := exec.LookPath(pythonCmd); err != nil {
			fmt.Println("⚠️ Python is not installed or is not in the PATH.")
			fmt.Println("Please install Python 3.x then run: rlama install-dependencies")
			return
		}
	}

	fmt.Println("Attempting system-wide installation...")
	packages := []string{"FlagEmbedding==1.2.10", "torch", "transformers", "pdfminer.six", "docx2txt", "xlsx2csv"}

	// Try --user flag first
	installCmd := exec.Command(pythonCmd, "-m", "pip", "install", "--user", "-U")
	installCmd.Args = append(installCmd.Args, packages...)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	if err := installCmd.Run(); err != nil {
		fmt.Printf("⚠️ Error with --user installation: %v\n", err)
		fmt.Println()
		fmt.Println("=== MANUAL INSTALLATION REQUIRED ===")
		fmt.Println("Due to PEP 668 (externally-managed-environment), you need to install dependencies manually.")
		fmt.Println("Please run the following commands:")
		fmt.Println()
		fmt.Printf("# Create virtual environment\n")
		fmt.Printf("python3 -m venv ~/.rlama/venv\n")
		fmt.Printf("\n# Activate virtual environment\n")
		fmt.Printf("source ~/.rlama/venv/bin/activate\n")
		fmt.Printf("\n# Install dependencies\n")
		fmt.Printf("pip install -U %s\n", "FlagEmbedding==1.2.10 torch transformers pdfminer.six docx2txt xlsx2csv")
		fmt.Printf("\n# Deactivate (optional)\n")
		fmt.Printf("deactivate\n")
	} else {
		fmt.Println("✅ Python dependencies installed successfully!")
	}
}
