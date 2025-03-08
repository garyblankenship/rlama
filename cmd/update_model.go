package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/repository"
)

var updateModelCmd = &cobra.Command{
	Use:   "update-model [rag-name] [new-model]",
	Short: "Update the Ollama model used by a RAG system",
	Long: `Change the Ollama model used by an existing RAG system.
Example: rlama update-model my-docs llama3.2
	
Note: This does not regenerate embeddings. For optimal results, you may want to
recreate the RAG with the new model instead.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]
		newModel := args[1]

		// Get Ollama client from root command
		ollamaClient := GetOllamaClient()
		
		// Check if Ollama is installed and the new model is available
		if err := ollamaClient.CheckOllamaAndModel(newModel); err != nil {
			return err
		}

		// Load the RAG
		repo := repository.NewRagRepository()
		rag, err := repo.Load(ragName)
		if err != nil {
			return err
		}

		// Update the model
		oldModel := rag.ModelName
		rag.ModelName = newModel

		// Save the RAG
		if err := repo.Save(rag); err != nil {
			return fmt.Errorf("error saving the RAG: %w", err)
		}

		fmt.Printf("Successfully updated RAG '%s' model from '%s' to '%s'.\n", 
			ragName, oldModel, newModel)
		fmt.Println("Note: Embeddings have not been regenerated. For optimal results, consider recreating the RAG.")
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateModelCmd)
}