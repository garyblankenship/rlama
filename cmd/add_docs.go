package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/service"
)

var addDocsCmd = &cobra.Command{
	Use:   "add-docs [rag-name] [folder-path]",
	Short: "Add documents to an existing RAG system",
	Long: `Add documents from a folder to an existing RAG system.
Example: rlama add-docs my-docs ./new-documents
	
This will load documents from the specified folder, generate embeddings,
and add them to the existing RAG system.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]
		folderPath := args[1]

		// Get Ollama client from root command
		ollamaClient := GetOllamaClient()
		
		// Create necessary services
		ragService := service.NewRagService(ollamaClient)
		documentLoader := service.NewDocumentLoader()
		embeddingService := service.NewEmbeddingService(ollamaClient)

		// Load the RAG
		rag, err := ragService.LoadRag(ragName)
		if err != nil {
			return err
		}

		// Check if Ollama is available with the model
		if err := ollamaClient.CheckOllamaAndModel(rag.ModelName); err != nil {
			return err
		}

		// Load documents from the folder
		fmt.Printf("Loading documents from '%s'...\n", folderPath)
		docs, err := documentLoader.LoadDocumentsFromFolder(folderPath)
		if err != nil {
			return fmt.Errorf("error loading documents: %w", err)
		}

		if len(docs) == 0 {
			return fmt.Errorf("no valid documents found in folder %s", folderPath)
		}

		fmt.Printf("Successfully loaded %d documents. Generating embeddings...\n", len(docs))

		// Generate embeddings for all documents
		err = embeddingService.GenerateEmbeddings(docs, rag.ModelName)
		if err != nil {
			return fmt.Errorf("error generating embeddings: %w", err)
		}

		// Track how many new documents were added
		docsAdded := 0
		existingDocs := make(map[string]bool)
		
		// Create a map of existing document IDs for quick lookup
		for _, doc := range rag.Documents {
			existingDocs[doc.ID] = true
		}

		// Add documents to the RAG if they don't already exist
		for _, doc := range docs {
			if !existingDocs[doc.ID] {
				rag.AddDocument(doc)
				docsAdded++
			} else {
				fmt.Printf("Skipping duplicate document: %s\n", doc.Name)
			}
		}

		// Save the RAG
		err = ragService.UpdateRag(rag)
		if err != nil {
			return fmt.Errorf("error saving the RAG: %w", err)
		}

		fmt.Printf("Successfully added %d new documents to RAG '%s'.\n", docsAdded, ragName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addDocsCmd)
}