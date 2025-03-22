package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var installDependenciesCmd = &cobra.Command{
	Use:   "install-dependencies",
	Short: "Installe les dépendances nécessaires pour RLAMA",
	Long:  `Installe les dépendances système et Python pour le fonctionnement optimal de RLAMA, y compris le reranker BGE.`,
	Run: func(cmd *cobra.Command, args []string) {
		installDependencies()
	},
}

func init() {
	rootCmd.AddCommand(installDependenciesCmd)
}

func installDependencies() {
	fmt.Println("Installation des dépendances...")

	// Trouver le chemin du script d'installation
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("⚠️ Erreur lors de la détermination du chemin de l'exécutable: %v\n", err)

		// Utiliser une solution alternative
		installDependenciesFallback()
		return
	}

	// Le répertoire scripts est présumé être dans le même répertoire que l'exécutable
	// ou dans le répertoire parent pour les environnements de développement
	scriptDir := filepath.Join(filepath.Dir(execPath), "scripts")
	scriptPath := filepath.Join(scriptDir, "install_deps.sh")

	// Vérifier si le script existe
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		// Essayer dans le répertoire parent (pour le développement)
		scriptDir = filepath.Join(filepath.Dir(execPath), "..", "scripts")
		scriptPath = filepath.Join(scriptDir, "install_deps.sh")

		if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
			fmt.Println("⚠️ Script d'installation non trouvé. Utilisation de la méthode alternative.")
			installDependenciesFallback()
			return
		}
	}

	// Exécuter le script
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// Sur Windows, utiliser bash.exe (WSL ou Git Bash)
		cmd = exec.Command("bash.exe", scriptPath)
	} else {
		// Sur Unix-like, exécuter directement le script
		cmd = exec.Command("bash", scriptPath)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠️ Erreur lors de l'exécution du script d'installation: %v\n", err)
		fmt.Println("Tentative d'installation directe des dépendances Python...")
		installPythonDependencies()
	}
}

// Cette fonction est utilisée si le script install_deps.sh n'est pas trouvé
func installDependenciesFallback() {
	fmt.Println("Installation des dépendances Python nécessaires...")
	installPythonDependencies()
}

func installPythonDependencies() {
	// Déterminer la commande Python à utiliser
	pythonCmd := "python3"
	if _, err := exec.LookPath(pythonCmd); err != nil {
		pythonCmd = "python"
		if _, err := exec.LookPath(pythonCmd); err != nil {
			fmt.Println("⚠️ Python n'est pas installé ou n'est pas dans le PATH.")
			fmt.Println("Veuillez installer Python 3.x puis exécuter: pip install -U FlagEmbedding torch transformers")
			return
		}
	}

	// Installer les dépendances Python
	fmt.Println("Installation de FlagEmbedding et des dépendances associées...")
	installCmd := exec.Command(pythonCmd, "-m", "pip", "install", "--user", "-U", "FlagEmbedding", "torch", "transformers")
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr

	if err := installCmd.Run(); err != nil {
		fmt.Printf("⚠️ Erreur lors de l'installation des dépendances Python: %v\n", err)
		fmt.Println("Veuillez installer manuellement: pip install -U FlagEmbedding torch transformers")
	} else {
		fmt.Println("✅ Dépendances Python installées avec succès!")
	}
}
