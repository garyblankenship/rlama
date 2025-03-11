package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/service"
)

var (
	showChunkContent bool
	documentFilter   string
)

var listChunksCmd = &cobra.Command{
	Use:   "list-chunks [rag-name]",
	Short: "Inspect document chunks in a RAG system",
	Long: `List and filter document chunks with various options.
	
Examples:
  # Basic chunk listing
  rlama list-chunks my-docs
  
  # Show chunk contents
  rlama list-chunks my-docs --show-content
  
  # Filter chunks from API documentation
  rlama list-chunks my-docs --document=api
  
  # Combine filters
  rlama list-chunks my-docs --document=readme --show-content`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]

		ollamaClient := GetOllamaClient()
		ragService := service.NewRagService(ollamaClient)

		// Get chunks with filters
		chunks, err := ragService.GetRagChunks(ragName, service.ChunkFilter{
			DocumentSubstring: documentFilter,
			ShowContent:       showChunkContent,
		})
		if err != nil {
			return err
		}

		// Display results
		fmt.Printf("Found %d chunks in RAG '%s'\n", len(chunks), ragName)
		for _, chunk := range chunks {
			fmt.Printf("\nChunk ID: %s\n", chunk.ID)
			fmt.Printf("Document: %s\n", chunk.DocumentID)
			fmt.Printf("Position: %d/%d\n", chunk.ChunkNumber+1, chunk.TotalChunks)
			
			if showChunkContent {
				fmt.Printf("Content:\n%s\n", strings.TrimSpace(chunk.Content))
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listChunksCmd)
	
	listChunksCmd.Flags().BoolVar(&showChunkContent, "show-content", false, "Display full chunk content")
	listChunksCmd.Flags().StringVar(&documentFilter, "document", "", "Filter by document name substring")
} 