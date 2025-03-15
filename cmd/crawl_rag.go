package cmd

import (
	"fmt"
	"strings"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/crawler"
	"github.com/dontizi/rlama/internal/service"
	"github.com/dontizi/rlama/internal/domain"
)

var (
	crawlMaxDepth     int
	crawlConcurrency  int
	crawlExcludePaths []string
)

var crawlRagCmd = &cobra.Command{
	Use:   "crawl-rag [model] [rag-name] [website-url]",
	Short: "Create a new RAG system from a website",
	Long: `Create a new RAG system by crawling a website and indexing its content.
Example: rlama crawl-rag llama3 mysite-rag https://example.com

The crawler will start at the provided URL and follow links to other pages 
on the same domain up to the specified depth.

You can exclude certain paths and control other crawling parameters:
  rlama crawl-rag llama3 my-docs https://docs.example.com --max-depth=2
  rlama crawl-rag llama3 blog-rag https://blog.example.com --exclude-path=/archive,/tags`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		modelName := args[0]
		ragName := args[1]
		websiteURL := args[2]

		// Get Ollama client with configured host and port
		ollamaClient := GetOllamaClient()
		if err := ollamaClient.CheckOllamaAndModel(modelName); err != nil {
			return err
		}

		// Create new crawler
		webCrawler, err := crawler.NewWebCrawler(websiteURL, crawlMaxDepth, crawlConcurrency, crawlExcludePaths)
		if err != nil {
			return fmt.Errorf("error initializing web crawler: %w", err)
		}

		// Display a message to indicate that the process has started
		fmt.Printf("Creating RAG '%s' with model '%s' by crawling website '%s'...\n", 
			ragName, modelName, websiteURL)
		fmt.Printf("Max crawl depth: %d, Concurrency: %d\n", crawlMaxDepth, crawlConcurrency)
		
		// Start crawling
		documents, err := webCrawler.CrawlWebsite()
		if err != nil {
			return fmt.Errorf("error crawling website: %w", err)
		}

		if len(documents) == 0 {
			return fmt.Errorf("no content found when crawling %s", websiteURL)
		}

		fmt.Printf("Retrieved %d pages from website. Processing content...\n", len(documents))

		// Create RAG service
		ragService := service.NewRagService(ollamaClient)

		// Set chunking options
		loaderOptions := service.DocumentLoaderOptions{
			ChunkSize:    chunkSize,
			ChunkOverlap: chunkOverlap,
		}

		// Create temporary directory to store crawled content
		tempDir := createTempDirForDocuments(documents)
		if tempDir != "" {
			defer cleanupTempDir(tempDir)
		}

		// Create RAG system
		err = ragService.CreateRagWithOptions(modelName, ragName, tempDir, loaderOptions)
		if err != nil {
			if strings.Contains(err.Error(), "connection refused") {
				return fmt.Errorf("⚠️ Unable to connect to Ollama.\n"+
					"Make sure Ollama is installed and running.\n")
			}
			return err
		}

		fmt.Printf("RAG '%s' created successfully with content from %s.\n", ragName, websiteURL)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(crawlRagCmd)
	
	// Add crawling specific flags
	crawlRagCmd.Flags().IntVar(&crawlMaxDepth, "max-depth", 2, "Maximum crawl depth (default: 2)")
	crawlRagCmd.Flags().IntVar(&crawlConcurrency, "concurrency", 5, "Number of concurrent crawlers (default: 5)")
	crawlRagCmd.Flags().StringSliceVar(&crawlExcludePaths, "exclude-path", nil, "Paths to exclude from crawling (comma-separated)")
	
	// Add chunking flags
	crawlRagCmd.Flags().IntVar(&chunkSize, "chunk-size", 1000, "Character count per chunk (default: 1000)")
	crawlRagCmd.Flags().IntVar(&chunkOverlap, "chunk-overlap", 200, "Overlap between chunks in characters (default: 200)")
}

// Helper function to create a temporary directory and save crawled documents as files
func createTempDirForDocuments(documents []*domain.Document) string {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "rlama-crawl-*")
	if err != nil {
		fmt.Printf("Error creating temporary directory: %v\n", err)
		return ""
	}
	
	fmt.Printf("Created temporary directory for documents: %s\n", tempDir)
	
	// Save each document as a file in the temporary directory
	for i, doc := range documents {
		// Default to index-based filename
		filename := fmt.Sprintf("page_%d.md", i+1)
		
		// Try to use Path if available (more likely to exist than URL)
		if doc.Path != "" {
			// Create a safe filename from the Path
			safePath := strings.Map(func(r rune) rune {
				if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
					return r
				}
				return '-'
			}, doc.Path)
			
			// Trim leading/trailing dashes
			safePath = strings.Trim(safePath, "-")
			if safePath != "" {
				filename = fmt.Sprintf("%s.md", safePath)
			}
		}
		
		// Full path to the file
		filePath := filepath.Join(tempDir, filename)
		
		// Write the document content to the file
		err := os.WriteFile(filePath, []byte(doc.Content), 0644)
		if err != nil {
			fmt.Printf("Error writing document to file %s: %v\n", filePath, err)
			continue
		}
	}
	
	return tempDir
}

func cleanupTempDir(path string) {
	if path != "" {
		if err := os.RemoveAll(path); err != nil {
			fmt.Printf("Warning: Failed to clean up temporary directory %s: %v\n", path, err)
		}
	}
}