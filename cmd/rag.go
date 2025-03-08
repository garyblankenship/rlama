package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/service"
)

var ragCmd = &cobra.Command{
	Use:   "rag [model] [rag-name] [folder-path]",
	Short: "Create a new RAG system",
	Long: `Create a new RAG system by indexing all documents in the specified folder.
Example: rlama rag llama3.2 rag1 ./documents

The folder will be created if it doesn't exist yet.
Supported formats include: .txt, .md, .html, .json, .csv, and various source code files.`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		modelName := args[0]
		ragName := args[1]
		folderPath := args[2]

		// Get Ollama client with configured host and port
		ollamaClient := GetOllamaClient()
		if err := ollamaClient.CheckOllamaAndModel(modelName); err != nil {
			return err
		}

		// Display a message to indicate that the process has started
		fmt.Printf("Creating RAG '%s' with model '%s' from folder '%s'...\n", 
			ragName, modelName, folderPath)

		ragService := service.NewRagService(ollamaClient)
		err := ragService.CreateRag(modelName, ragName, folderPath)
		if err != nil {
			// Improve error messages related to Ollama
			if strings.Contains(err.Error(), "connection refused") {
				return fmt.Errorf("⚠️ Unable to connect to Ollama.\n"+
					"Make sure Ollama is installed and running.\n")
			}
			return err
		}

		fmt.Printf("RAG '%s' created successfully.\n", ragName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(ragCmd)
}