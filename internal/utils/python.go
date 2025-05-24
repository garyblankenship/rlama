package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// PythonExecutor handles Python command execution with proper command detection
type PythonExecutor struct {
	pythonCmd string
	venvPath  string
}

// NewPythonExecutor creates a new Python executor with automatic command detection
func NewPythonExecutor() *PythonExecutor {
	pythonCmd := detectPythonCommand()
	venvPath := getVirtualEnvPath()

	return &PythonExecutor{
		pythonCmd: pythonCmd,
		venvPath:  venvPath,
	}
}

// detectPythonCommand finds the appropriate Python command to use
func detectPythonCommand() string {
	// Check for python3 first (preferred on modern systems)
	if _, err := exec.LookPath("python3"); err == nil {
		return "python3"
	}

	// Fall back to python
	if _, err := exec.LookPath("python"); err == nil {
		return "python"
	}

	return "python3" // Default fallback
}

// getVirtualEnvPath returns the path to the rlama virtual environment
func getVirtualEnvPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(homeDir, ".rlama", "venv")
}

// GetPythonCommand returns the Python command, preferring the virtual environment if available
func (pe *PythonExecutor) GetPythonCommand() string {
	venvPython := filepath.Join(pe.venvPath, "bin", "python")
	if _, err := os.Stat(venvPython); err == nil {
		return venvPython
	}

	// On Windows, try Scripts directory
	venvPython = filepath.Join(pe.venvPath, "Scripts", "python.exe")
	if _, err := os.Stat(venvPython); err == nil {
		return venvPython
	}

	return pe.pythonCmd
}

// CreateVirtualEnvironment creates a virtual environment for rlama
func (pe *PythonExecutor) CreateVirtualEnvironment() error {
	// Create the rlama directory if it doesn't exist
	rlamaDir := filepath.Dir(pe.venvPath)
	if err := os.MkdirAll(rlamaDir, 0755); err != nil {
		return fmt.Errorf("failed to create rlama directory: %w", err)
	}

	// Check if virtual environment already exists
	if _, err := os.Stat(pe.venvPath); err == nil {
		fmt.Println("Virtual environment already exists")
		return nil
	}

	fmt.Printf("Creating virtual environment at %s...\n", pe.venvPath)

	// Create virtual environment
	cmd := exec.Command(pe.pythonCmd, "-m", "venv", pe.venvPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set environment variables to force English locale
	cmd.Env = append(os.Environ(),
		"LC_ALL=C",
		"LANG=en_US.UTF-8",
		"LANGUAGE=en_US:en",
		"PYTHONIOENCODING=utf-8",
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create virtual environment: %w", err)
	}

	fmt.Println("✅ Virtual environment created successfully!")
	return nil
}

// InstallPackage installs a Python package in the virtual environment
func (pe *PythonExecutor) InstallPackage(packageName string) error {
	pythonCmd := pe.GetPythonCommand()

	fmt.Printf("Installing %s...\n", packageName)
	cmd := exec.Command(pythonCmd, "-m", "pip", "install", "-U", packageName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set environment variables to force English locale
	cmd.Env = append(os.Environ(),
		"LC_ALL=C",
		"LANG=en_US.UTF-8",
		"LANGUAGE=en_US:en",
		"PYTHONIOENCODING=utf-8",
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install %s: %w", packageName, err)
	}

	fmt.Printf("✅ %s installed successfully!\n", packageName)
	return nil
}

// CheckPackageInstalled checks if a package is installed in the virtual environment
func (pe *PythonExecutor) CheckPackageInstalled(packageName string) bool {
	pythonCmd := pe.GetPythonCommand()

	checkScript := fmt.Sprintf(`
import importlib.util
import sys

spec = importlib.util.find_spec("%s")
if spec is None:
    sys.exit(1)
else:
    sys.exit(0)
`, packageName)

	cmd := exec.Command(pythonCmd, "-c", checkScript)

	// Set environment variables to force English locale
	cmd.Env = append(os.Environ(),
		"LC_ALL=C",
		"LANG=en_US.UTF-8",
		"LANGUAGE=en_US:en",
		"PYTHONIOENCODING=utf-8",
	)

	return cmd.Run() == nil
}

// ExecuteScript executes a Python script using the proper Python command
func (pe *PythonExecutor) ExecuteScript(script string, stdin ...string) ([]byte, error) {
	pythonCmd := pe.GetPythonCommand()

	cmd := exec.Command(pythonCmd, "-c", script)
	if len(stdin) > 0 {
		cmd.Stdin = strings.NewReader(stdin[0])
	}

	// Set environment variables to force English locale and avoid French warnings
	cmd.Env = append(os.Environ(),
		"LC_ALL=C",
		"LANG=en_US.UTF-8",
		"LANGUAGE=en_US:en",
		"PYTHONIOENCODING=utf-8",
	)

	return cmd.CombinedOutput()
}

// GetVirtualEnvPath returns the virtual environment path
func (pe *PythonExecutor) GetVirtualEnvPath() string {
	return pe.venvPath
}
