package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/service"
	"github.com/dontizi/rlama/internal/util"
)

var listDocsCmd = &cobra.Command{
	Use:   "list-docs [rag-name]",
	Short: "List all documents in a RAG system",
	Long: `Display a list of all documents in a specified RAG system.
Example: rlama list-docs my-docs`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]

		// Get Ollama client from root command
		ollamaClient := GetOllamaClient()

		// Create necessary services
		ragService := service.NewRagService(ollamaClient)

		// Load the RAG
		rag, err := ragService.LoadRag(ragName)
		if err != nil {
			return err
		}

		if len(rag.Documents) == 0 {
			fmt.Printf("No documents found in RAG '%s'.\n", ragName)
			return nil
		}

		fmt.Printf("Documents in RAG '%s' (%d found):\n\n", ragName, len(rag.Documents))
		
		// Use tabwriter for aligned display
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tPATH\tSIZE\tCONTENT TYPE")
		
		for _, doc := range rag.Documents {
			sizeStr := util.FormatSize(doc.Size)
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", doc.ID, doc.Name, sizeStr, doc.ContentType)
		}
		w.Flush()
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listDocsCmd)
}

