package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"os/exec"
	"bytes"
	
	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/service"
)

// Structure pour parser la sortie JSON d'Ollama list
type OllamaModel struct {
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	ModifiedAt string `json:"modified_at"`
	Digest     string `json:"digest"`
}

var (
	// Variables pour le wizard local
	localWizardModel       string
	localWizardName        string
	localWizardPath        string
	localWizardChunkSize   int
	localWizardChunkOverlap int
	localWizardExcludeDirs []string
	localWizardExcludeExts []string
	localWizardProcessExts []string
)

// Renommé pour éviter le conflit avec snowflake_wizard.go
var localWizardCmd = &cobra.Command{
	Use:   "wizard",
	Short: "Interactive wizard to create a local RAG",
	Long: `Start an interactive wizard that guides you through creating a RAG system.
This makes it easy to set up a new RAG without remembering all command options.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("\n🧙 Welcome to the RLAMA Local RAG Wizard! 🧙\n")
		
		reader := bufio.NewReader(os.Stdin)
		
		// Étape 1: Nom du RAG
		fmt.Print("Enter a name for your RAG: ")
		ragName, _ := reader.ReadString('\n')
		ragName = strings.TrimSpace(ragName)
		if ragName == "" {
			return fmt.Errorf("RAG name cannot be empty")
		}
		
		// Déclarer modelName au niveau de la fonction pour qu'il soit disponible partout
		var modelName string
		
		// Étape 2: Sélection du modèle
		fmt.Println("\nStep 2: Select a model")
		
		// Récupérer la liste des modèles disponibles via la commande ollama list
		fmt.Println("Retrieving available Ollama models...")
		
		// Tester d'abord avec ollama list sans --json pour plus de compatibilité
		// et capturer stderr pour le débogage
		var stdout, stderr bytes.Buffer
		ollamaCmd := exec.Command("ollama", "list")
		ollamaCmd.Stdout = &stdout
		ollamaCmd.Stderr = &stderr
		
		// Configuration pour l'exécution de la commande
		ollamaHost := os.Getenv("OLLAMA_HOST")
		if cmd.Flag("host").Changed {
			ollamaHost = cmd.Flag("host").Value.String()
		}
		
		if ollamaHost != "" {
			// Définir la variable d'environnement OLLAMA_HOST pour la commande
			ollamaCmd.Env = append(os.Environ(), fmt.Sprintf("OLLAMA_HOST=%s", ollamaHost))
		}
		
		// Exécuter la commande
		err := ollamaCmd.Run()
		if err != nil {
			fmt.Println("❌ Failed to list Ollama models.")
			if stderr.Len() > 0 {
				fmt.Printf("Error details: %s\n", stderr.String())
			}
			fmt.Println("Make sure Ollama is installed and running.")
			fmt.Println("Continuing without model list. You'll need to enter a model name manually.")
		}
		
		// Analyser la sortie de ollama list (format texte)
		modelsOutput := stdout.String()
		var modelNames []string
		
		if modelsOutput != "" {
			// Format typique:
			// NAME             ID            SIZE    MODIFIED
			// llama3           xxx...xxx     4.7 GB  X days ago
			
			// Ignorer la première ligne (en-têtes)
			lines := strings.Split(modelsOutput, "\n")
			for i, line := range lines {
				if i == 0 || strings.TrimSpace(line) == "" {
					continue
				}
				
				fields := strings.Fields(line)
				if len(fields) >= 1 {
					modelNames = append(modelNames, fields[0])
				}
			}
			
			// Afficher les modèles dans notre format
			if len(modelNames) > 0 {
				fmt.Println("\nAvailable models:")
				for i, name := range modelNames {
					fmt.Printf("  %d. %s\n", i+1, name)
				}
				
				// Permettre à l'utilisateur de choisir un modèle
				fmt.Print("\nChoose a model (number) or enter model name: ")
				modelChoice, _ := reader.ReadString('\n')
				modelChoice = strings.TrimSpace(modelChoice)
				
				// Vérifier si l'utilisateur a entré un numéro
				var modelNumber int
				modelName = "" // Initialiser ici aussi
				
				if _, err := fmt.Sscanf(modelChoice, "%d", &modelNumber); err == nil {
					// L'utilisateur a entré un numéro
					if modelNumber > 0 && modelNumber <= len(modelNames) {
						modelName = modelNames[modelNumber-1]
					} else {
						fmt.Println("Invalid selection. Please enter a valid model name manually.")
					}
				} else {
					// L'utilisateur a entré un nom directement
					modelName = modelChoice
				}
			}
		}
		
		// Si aucun modèle n'a été sélectionné, demander manuellement
		if modelName == "" {
			fmt.Print("Enter model name [llama3]: ")
			inputName, _ := reader.ReadString('\n')
			inputName = strings.TrimSpace(inputName)
			if inputName == "" {
				modelName = "llama3"
			} else {
				modelName = inputName
			}
		}
		
		// Étape 3: Chemin des documents
		fmt.Println("\nStep 3: Specify documents path")
		fmt.Print("Enter path to document folder: ")
		folderPath, _ := reader.ReadString('\n')
		folderPath = strings.TrimSpace(folderPath)
		if folderPath == "" {
			return fmt.Errorf("folder path cannot be empty")
		}
		
		// Étape 4: Options de chunking
		fmt.Println("\nStep 4: Chunking options")
		
		fmt.Print("Chunk size [1000]: ")
		chunkSizeStr, _ := reader.ReadString('\n')
		chunkSizeStr = strings.TrimSpace(chunkSizeStr)
		chunkSize := 1000
		if chunkSizeStr != "" {
			fmt.Sscanf(chunkSizeStr, "%d", &chunkSize)
		}
		
		fmt.Print("Chunk overlap [200]: ")
		overlapStr, _ := reader.ReadString('\n')
		overlapStr = strings.TrimSpace(overlapStr)
		overlap := 200
		if overlapStr != "" {
			fmt.Sscanf(overlapStr, "%d", &overlap)
		}
		
		// Étape 5: Filtrer les fichiers (optionnel)
		fmt.Println("\nStep 5: File filtering (optional)")
		
		fmt.Print("Exclude directories (comma-separated): ")
		excludeDirsStr, _ := reader.ReadString('\n')
		excludeDirsStr = strings.TrimSpace(excludeDirsStr)
		var excludeDirs []string
		if excludeDirsStr != "" {
			excludeDirs = strings.Split(excludeDirsStr, ",")
			for i := range excludeDirs {
				excludeDirs[i] = strings.TrimSpace(excludeDirs[i])
			}
		}
		
		fmt.Print("Exclude extensions (comma-separated): ")
		excludeExtsStr, _ := reader.ReadString('\n')
		excludeExtsStr = strings.TrimSpace(excludeExtsStr)
		var excludeExts []string
		if excludeExtsStr != "" {
			excludeExts = strings.Split(excludeExtsStr, ",")
			for i := range excludeExts {
				excludeExts[i] = strings.TrimSpace(excludeExts[i])
			}
		}
		
		fmt.Print("Process only these extensions (comma-separated): ")
		processExtsStr, _ := reader.ReadString('\n')
		processExtsStr = strings.TrimSpace(processExtsStr)
		var processExts []string
		if processExtsStr != "" {
			processExts = strings.Split(processExtsStr, ",")
			for i := range processExts {
				processExts[i] = strings.TrimSpace(processExts[i])
			}
		}
		
		// Étape 6: Confirmation et création
		fmt.Println("\nStep 6: Review and create")
		fmt.Println("RAG configuration:")
		fmt.Printf("- Name: %s\n", ragName)
		fmt.Printf("- Model: %s\n", modelName)
		fmt.Printf("- Documents: %s\n", folderPath)
		fmt.Printf("- Chunk size: %d\n", chunkSize)
		fmt.Printf("- Chunk overlap: %d\n", overlap)
		
		if len(excludeDirs) > 0 {
			fmt.Printf("- Exclude directories: %s\n", strings.Join(excludeDirs, ", "))
		}
		if len(excludeExts) > 0 {
			fmt.Printf("- Exclude extensions: %s\n", strings.Join(excludeExts, ", "))
		}
		if len(processExts) > 0 {
			fmt.Printf("- Process only: %s\n", strings.Join(processExts, ", "))
		}
		
		fmt.Print("\nCreate RAG with these settings? (y/n): ")
		confirm, _ := reader.ReadString('\n')
		confirm = strings.ToLower(strings.TrimSpace(confirm))
		
		if confirm != "y" && confirm != "yes" {
			fmt.Println("RAG creation cancelled.")
			return nil
		}
		
		// Créer le RAG
		fmt.Println("\nCreating RAG...")
		
		// Obtenir le client Ollama configuré
		ollamaClient := GetOllamaClient()
		
		// Vérifier que le modèle est disponible avant de continuer
		// Cette étape est importante pour éviter les erreurs plus tard
		fmt.Printf("Checking if model '%s' is available...\n", modelName)
		err = ollamaClient.CheckOllamaAndModel(modelName)
		if err != nil {
			return fmt.Errorf("model '%s' is not available: %w", modelName, err)
		}
		
		// Utiliser RagService pour créer le RAG
		ragService := service.NewRagService(ollamaClient)
		
		// Préparation des options
		loaderOptions := service.DocumentLoaderOptions{
			ExcludeDirs:  excludeDirs,
			ExcludeExts:  excludeExts,
			ProcessExts:  processExts,
			ChunkSize:    chunkSize,
			ChunkOverlap: overlap,
		}
		
		err = ragService.CreateRagWithOptions(modelName, ragName, folderPath, loaderOptions)
		if err != nil {
			return fmt.Errorf("error creating RAG: %w", err)
		}
		
		fmt.Println("\n🎉 RAG created successfully! 🎉")
		fmt.Printf("\nYou can now use your RAG with: rlama run %s\n", ragName)
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(localWizardCmd)
} 