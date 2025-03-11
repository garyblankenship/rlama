package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/service"
)

var (
	addDocsExcludeDirs  []string
	addDocsExcludeExts  []string
	addDocsProcessExts  []string
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

		// Set up loader options based on flags
		loaderOptions := service.DocumentLoaderOptions{
			ExcludeDirs: addDocsExcludeDirs,
			ExcludeExts: addDocsExcludeExts,
			ProcessExts: addDocsProcessExts,
		}

		// Pass the options to the service
		err := ragService.AddDocsWithOptions(ragName, folderPath, loaderOptions)
		if err != nil {
			return err
		}

		fmt.Printf("Documents from '%s' added to RAG '%s' successfully.\n", folderPath, ragName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addDocsCmd)

	// Add exclusion and processing flags
	addDocsCmd.Flags().StringSliceVar(&addDocsExcludeDirs, "excludedir", nil, "Directories to exclude (comma-separated)")
	addDocsCmd.Flags().StringSliceVar(&addDocsExcludeExts, "excludeext", nil, "File extensions to exclude (comma-separated)")
	addDocsCmd.Flags().StringSliceVar(&addDocsProcessExts, "processext", nil, "Only process these file extensions (comma-separated)")
}