package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/service"
)

var forceRemoveDoc bool

var removeDocCmd = &cobra.Command{
	Use:   "remove-doc [rag-name] [doc-id]",
	Short: "Remove a document from a RAG system",
	Long: `Remove a specific document from a RAG system by its ID.
Example: rlama remove-doc my-docs document.pdf
	
The document ID is typically the filename. You can see document IDs by using the
"rlama list-docs [rag-name]" command.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]
		docID := args[1]

		// Get Ollama client from root command
		ollamaClient := GetOllamaClient()

		// Create necessary services
		ragService := service.NewRagService(ollamaClient)

		// Load the RAG
		rag, err := ragService.LoadRag(ragName)
		if err != nil {
			return err
		}

		// Find the document
		doc := rag.GetDocumentByID(docID)
		if doc == nil {
			return fmt.Errorf("document with ID '%s' not found in RAG '%s'", docID, ragName)
		}

		// Ask for confirmation unless --force is specified
		if !forceRemoveDoc {
			fmt.Printf("Are you sure you want to remove document '%s' from RAG '%s'? (y/n): ", 
				doc.Name, ragName)
			var response string
			fmt.Scanln(&response)
			
			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" && response != "yes" {
				fmt.Println("Document removal cancelled.")
				return nil
			}
		}

		// Remove the document
		removed := rag.RemoveDocument(docID)
		if !removed {
			return fmt.Errorf("failed to remove document '%s'", docID)
		}

		// Save the RAG
		err = ragService.UpdateRag(rag)
		if err != nil {
			return fmt.Errorf("error saving the RAG: %w", err)
		}

		fmt.Printf("Successfully removed document '%s' from RAG '%s'.\n", doc.Name, ragName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeDocCmd)
	removeDocCmd.Flags().BoolVarP(&forceRemoveDoc, "force", "f", false, "Remove without asking for confirmation")
}