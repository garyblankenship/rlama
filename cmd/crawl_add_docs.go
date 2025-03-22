package cmd

import (
	"fmt"

	"github.com/dontizi/rlama/internal/crawler"
	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/service"
	"github.com/spf13/cobra"
)

var (
	addCrawlMaxDepth     int
	addCrawlConcurrency  int
	addCrawlExcludePaths []string
	addCrawlUseSitemap   bool
)

var crawlAddDocsCmd = &cobra.Command{
	Use:   "crawl-add-docs [rag-name] [website-url]",
	Short: "Add website content to an existing RAG system",
	Long: `Add content from a website to an existing RAG system.
Example: rlama crawl-add-docs my-docs https://blog.example.com
	
This will crawl the website, extract content, generate embeddings,
and add them to the existing RAG system.

Control the crawling behavior with these flags:
  --max-depth=3         Maximum depth of pages to crawl
  --concurrency=10      Number of concurrent crawlers
  --exclude-path=/tag   Skip specific path patterns (comma-separated)
  --use-sitemap         Use sitemap.xml if available for comprehensive coverage`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]
		websiteURL := args[1]

		// Get Ollama client from root command
		ollamaClient := GetOllamaClient()

		// Create necessary services
		ragService := service.NewRagService(ollamaClient)

		// Load existing RAG to get model name
		_, err := ragService.LoadRag(ragName)
		if err != nil {
			return err
		}

		// Create new crawler
		webCrawler, err := crawler.NewWebCrawler(websiteURL, addCrawlMaxDepth, addCrawlConcurrency, addCrawlExcludePaths)
		if err != nil {
			return fmt.Errorf("error initializing web crawler: %w", err)
		}

		// Définir l'option de sitemap
		webCrawler.SetUseSitemap(addCrawlUseSitemap)

		fmt.Printf("Crawling website '%s' to add content to RAG '%s'...\n", websiteURL, ragName)

		// Start crawling
		documents, err := webCrawler.CrawlWebsite()
		if err != nil {
			return fmt.Errorf("error crawling website: %w", err)
		}

		if len(documents) == 0 {
			return fmt.Errorf("no content found when crawling %s", websiteURL)
		}

		fmt.Printf("Retrieved %d pages from website. Processing content...\n", len(documents))

		// Convertir []domain.Document en []*domain.Document
		var docPointers []*domain.Document
		for i := range documents {
			docPointers = append(docPointers, &documents[i])
		}

		// Create temporary directory to store crawled content
		tempDir := createTempDirForDocuments(docPointers)
		if tempDir != "" {
			defer cleanupTempDir(tempDir)
		}

		// Set up loader options
		loaderOptions := service.DocumentLoaderOptions{
			ChunkSize:    chunkSize,
			ChunkOverlap: chunkOverlap,
		}

		// Pass the options to the service
		err = ragService.AddDocsWithOptions(ragName, tempDir, loaderOptions)
		if err != nil {
			return err
		}

		fmt.Printf("Website content from '%s' added to RAG '%s' successfully.\n", websiteURL, ragName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(crawlAddDocsCmd)

	// Add crawling specific flags
	crawlAddDocsCmd.Flags().IntVar(&addCrawlMaxDepth, "max-depth", 2, "Maximum crawl depth")
	crawlAddDocsCmd.Flags().IntVar(&addCrawlConcurrency, "concurrency", 5, "Number of concurrent crawlers")
	crawlAddDocsCmd.Flags().StringSliceVar(&addCrawlExcludePaths, "exclude-path", nil, "Paths to exclude from crawling (comma-separated)")
	crawlAddDocsCmd.Flags().BoolVar(&addCrawlUseSitemap, "use-sitemap", true, "Use sitemap.xml if available for comprehensive coverage")

	// Add chunking flags
	crawlAddDocsCmd.Flags().IntVar(&chunkSize, "chunk-size", 1000, "Character count per chunk (default: 1000)")
	crawlAddDocsCmd.Flags().IntVar(&chunkOverlap, "chunk-overlap", 200, "Overlap between chunks in characters (default: 200)")
}
