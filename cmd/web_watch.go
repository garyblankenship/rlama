package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/service"
)

var (
	webWatchMaxDepth     int
	webWatchConcurrency  int
	webWatchExcludePaths []string
	webWatchChunkSize    int
	webWatchChunkOverlap int
)

var webWatchCmd = &cobra.Command{
	Use:   "web-watch [rag-name] [website-url] [interval]",
	Short: "Set up website watching for a RAG system",
	Long: `Configure a RAG system to automatically watch a website for new content and add it to the RAG.
The interval is specified in minutes. Use 0 to only check when the RAG is used.

Example: rlama web-watch my-docs https://docs.example.com 60
This will check the website every 60 minutes for new content.

Use rlama web-watch-off [rag-name] to disable watching.`,
	Args: cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]
		websiteURL := args[1]
		
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
		
		// Set up web watch options
		webWatchOptions := domain.WebWatchOptions{
			MaxDepth:     webWatchMaxDepth,
			Concurrency:  webWatchConcurrency,
			ExcludePaths: webWatchExcludePaths,
			ChunkSize:    webWatchChunkSize,
			ChunkOverlap: webWatchChunkOverlap,
		}
		
		// Set up website watching
		err := ragService.SetupWebWatching(ragName, websiteURL, interval, webWatchOptions)
		if err != nil {
			return err
		}
		
		// Provide feedback based on the interval
		if interval == 0 {
			fmt.Printf("Website watching set up for RAG '%s'. Website '%s' will be checked each time the RAG is used.\n", 
				ragName, websiteURL)
		} else {
			fmt.Printf("Website watching set up for RAG '%s'. Website '%s' will be checked every %d minutes.\n", 
				ragName, websiteURL, interval)
		}
		
		// Perform an initial check
		docsAdded, err := ragService.CheckWatchedWebsite(ragName)
		if err != nil {
			return fmt.Errorf("error during initial website check: %w", err)
		}
		
		if docsAdded > 0 {
			fmt.Printf("Added %d new pages from '%s'.\n", docsAdded, websiteURL)
		} else {
			fmt.Printf("No new content found at '%s'.\n", websiteURL)
		}
		
		return nil
	},
}

var webWatchOffCmd = &cobra.Command{
	Use:   "web-watch-off [rag-name]",
	Short: "Disable website watching for a RAG system",
	Long:  `Disable automatic website watching for a RAG system.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]
		
		// Get Ollama client from root command
		ollamaClient := GetOllamaClient()
		
		// Create RAG service
		ragService := service.NewRagService(ollamaClient)
		
		// Disable website watching
		err := ragService.DisableWebWatching(ragName)
		if err != nil {
			return err
		}
		
		fmt.Printf("Website watching disabled for RAG '%s'.\n", ragName)
		return nil
	},
}

var checkWebWatchedCmd = &cobra.Command{
	Use:   "check-web-watched [rag-name]",
	Short: "Check a RAG's watched website for new content",
	Long:  `Manually check a RAG's watched website for new content and add it to the RAG.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]
		
		// Get Ollama client from root command
		ollamaClient := GetOllamaClient()
		
		// Create RAG service
		ragService := service.NewRagService(ollamaClient)
		
		// Check the watched website
		pagesAdded, err := ragService.CheckWatchedWebsite(ragName)
		if err != nil {
			return err
		}
		
		if pagesAdded > 0 {
			fmt.Printf("Added %d new pages to RAG '%s'.\n", pagesAdded, ragName)
		} else {
			fmt.Printf("No new content found for RAG '%s'.\n", ragName)
		}
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(webWatchCmd)
	rootCmd.AddCommand(webWatchOffCmd)
	rootCmd.AddCommand(checkWebWatchedCmd)
	
	// Add web watching specific flags
	webWatchCmd.Flags().IntVar(&webWatchMaxDepth, "max-depth", 2, "Maximum crawl depth (default: 2)")
	webWatchCmd.Flags().IntVar(&webWatchConcurrency, "concurrency", 5, "Number of concurrent crawlers (default: 5)")
	webWatchCmd.Flags().StringSliceVar(&webWatchExcludePaths, "exclude-path", nil, "Paths to exclude from crawling (comma-separated)")
	webWatchCmd.Flags().IntVar(&webWatchChunkSize, "chunk-size", 1000, "Character count per chunk (default: 1000)")
	webWatchCmd.Flags().IntVar(&webWatchChunkOverlap, "chunk-overlap", 200, "Overlap between chunks in characters (default: 200)")
} 