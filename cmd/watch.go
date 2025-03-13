package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/service"
)

var (
	watchExcludeDirs  []string
	watchExcludeExts  []string
	watchProcessExts  []string
	watchChunkSize    int
	watchChunkOverlap int
)

var watchCmd = &cobra.Command{
	Use:   "watch [rag-name] [directory-path] [interval]",
	Short: "Set up directory watching for a RAG system",
	Long: `Configure a RAG system to automatically watch a directory for new files and add them to the RAG.
The interval is specified in minutes. Use 0 to only check when the RAG is used.

Example: rlama watch my-docs ./documents 60
This will check the ./documents directory every 60 minutes for new files.

Use rlama watch-off [rag-name] to disable watching.`,
	Args: cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]
		dirPath := args[1]
		
		// Default interval is 0 (only check when RAG is used)
		interval := 0
		
		// If an interval is provided, parse it
		if len(args) > 2 {
			var err error
			interval, err = strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("invalid interval: %s", args[2])
			}
			
			if interval < 0 {
				return fmt.Errorf("interval must be non-negative")
			}
		}
		
		// Get Ollama client from root command
		ollamaClient := GetOllamaClient()
		
		// Create RAG service
		ragService := service.NewRagService(ollamaClient)
		
		// Set up loader options based on flags
		loaderOptions := service.DocumentLoaderOptions{
			ExcludeDirs:  watchExcludeDirs,
			ExcludeExts:  watchExcludeExts,
			ProcessExts:  watchProcessExts,
			ChunkSize:    watchChunkSize,
			ChunkOverlap: watchChunkOverlap,
		}
		
		// Set up directory watching
		err := ragService.SetupDirectoryWatching(ragName, dirPath, interval, loaderOptions)
		if err != nil {
			return err
		}
		
		// Provide feedback based on the interval
		if interval == 0 {
			fmt.Printf("Directory watching set up for RAG '%s'. Directory '%s' will be checked each time the RAG is used.\n", 
				ragName, dirPath)
		} else {
			fmt.Printf("Directory watching set up for RAG '%s'. Directory '%s' will be checked every %d minutes.\n", 
				ragName, dirPath, interval)
		}
		
		// Perform an initial check
		docsAdded, err := ragService.CheckWatchedDirectory(ragName)
		if err != nil {
			return fmt.Errorf("error during initial directory check: %w", err)
		}
		
		if docsAdded > 0 {
			fmt.Printf("Added %d new documents from '%s'.\n", docsAdded, dirPath)
		} else {
			fmt.Printf("No new documents found in '%s'.\n", dirPath)
		}
		
		return nil
	},
}

var watchOffCmd = &cobra.Command{
	Use:   "watch-off [rag-name]",
	Short: "Disable directory watching for a RAG system",
	Long:  `Disable automatic directory watching for a RAG system.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]
		
		// Get Ollama client from root command
		ollamaClient := GetOllamaClient()
		
		// Create RAG service
		ragService := service.NewRagService(ollamaClient)
		
		// Disable directory watching
		err := ragService.DisableDirectoryWatching(ragName)
		if err != nil {
			return err
		}
		
		fmt.Printf("Directory watching disabled for RAG '%s'.\n", ragName)
		return nil
	},
}

var checkWatchedCmd = &cobra.Command{
	Use:   "check-watched [rag-name]",
	Short: "Check a RAG's watched directory for new files",
	Long:  `Manually check a RAG's watched directory for new files and add them to the RAG.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]
		
		// Get Ollama client from root command
		ollamaClient := GetOllamaClient()
		
		// Create RAG service
		ragService := service.NewRagService(ollamaClient)
		
		// Check the watched directory
		docsAdded, err := ragService.CheckWatchedDirectory(ragName)
		if err != nil {
			return err
		}
		
		if docsAdded > 0 {
			fmt.Printf("Added %d new documents to RAG '%s'.\n", docsAdded, ragName)
		} else {
			fmt.Printf("No new documents found for RAG '%s'.\n", ragName)
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(watchCmd)
	rootCmd.AddCommand(watchOffCmd)
	rootCmd.AddCommand(checkWatchedCmd)
	
	// Add exclusion and processing flags
	watchCmd.Flags().StringSliceVar(&watchExcludeDirs, "exclude-dir", nil, "Directories to exclude (comma-separated)")
	watchCmd.Flags().StringSliceVar(&watchExcludeExts, "exclude-ext", nil, "File extensions to exclude (comma-separated)")
	watchCmd.Flags().StringSliceVar(&watchProcessExts, "process-ext", nil, "Only process these file extensions (comma-separated)")
	watchCmd.Flags().IntVar(&watchChunkSize, "chunk-size", 1000, "Character count per chunk (default: 1000)")
	watchCmd.Flags().IntVar(&watchChunkOverlap, "chunk-overlap", 200, "Overlap between chunks in characters (default: 200)")
} 