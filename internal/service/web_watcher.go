package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dontizi/rlama/internal/crawler"
	"github.com/dontizi/rlama/internal/domain"
)

// WebWatcher is responsible for watching websites for content changes
type WebWatcher struct {
	ragService RagService
}

// NewWebWatcher creates a new web watcher service
func NewWebWatcher(ragService RagService) *WebWatcher {
	return &WebWatcher{
		ragService: ragService,
	}
}

// CheckAndUpdateRag checks for new content on the watched website and updates the RAG
func (ww *WebWatcher) CheckAndUpdateRag(rag *domain.RagSystem) (int, error) {
	if !rag.WebWatchEnabled || rag.WatchedURL == "" {
		return 0, nil // Watching not enabled
	}

	// Create a webcrawler to fetch the site content
	webCrawler, err := crawler.NewWebCrawler(
		rag.WatchedURL,
		rag.WebWatchOptions.MaxDepth,
		rag.WebWatchOptions.Concurrency,
		rag.WebWatchOptions.ExcludePaths,
	)
	if err != nil {
		return 0, fmt.Errorf("error initializing web crawler: %w", err)
	}

	// Start crawling
	documents, err := webCrawler.CrawlWebsite()
	if err != nil {
		return 0, fmt.Errorf("error crawling website: %w", err)
	}

	if len(documents) == 0 {
		// Update last watched time even if no new documents
		rag.LastWebWatchAt = time.Now()
		err = ww.ragService.UpdateRag(rag)
		return 0, err
	}

	// Get existing document URLs to avoid re-processing
	existingURLs := make(map[string]bool)
	for _, doc := range rag.Documents {
		if doc.URL != "" {
			existingURLs[doc.URL] = true
		}
	}

	// Filter out existing documents
	var newDocs []*domain.Document
	for _, doc := range documents {
		if !existingURLs[doc.URL] {
			newDocs = append(newDocs, doc)
		}
	}

	if len(newDocs) == 0 {
		// Update last watched time even if no new documents
		rag.LastWebWatchAt = time.Now()
		err = ww.ragService.UpdateRag(rag)
		return 0, err
	}

	// Create temporary directory to store crawled content
	tempDir := createTempDirForDocuments(newDocs)
	if tempDir == "" {
		return 0, fmt.Errorf("failed to create temporary directory for documents")
	}
	defer cleanupTempDir(tempDir)

	// Set chunking options
	loaderOptions := DocumentLoaderOptions{
		ChunkSize:    rag.WebWatchOptions.ChunkSize,
		ChunkOverlap: rag.WebWatchOptions.ChunkOverlap,
	}

	// Create chunker service
	chunkerService := NewChunkerService(ChunkingConfig{
		ChunkSize:    loaderOptions.ChunkSize,
		ChunkOverlap: loaderOptions.ChunkOverlap,
	})

	// Process each new document - chunk and prepare for embeddings
	var allChunks []*domain.DocumentChunk
	for _, doc := range newDocs {
		// Chunk the document
		chunks := chunkerService.ChunkDocument(doc)
		
		// Update total chunks in metadata
		for i, chunk := range chunks {
			chunk.ChunkNumber = i
			chunk.TotalChunks = len(chunks)
		}
		
		allChunks = append(allChunks, chunks...)
	}

	// Generate embeddings for all chunks
	embeddingService := NewEmbeddingService(ww.ragService.GetOllamaClient())
	err = embeddingService.GenerateChunkEmbeddings(allChunks, rag.ModelName)
	if err != nil {
		return 0, fmt.Errorf("error generating embeddings for new documents: %w", err)
	}

	// Add documents and chunks to the RAG
	for _, doc := range newDocs {
		rag.AddDocument(doc)
	}
	
	for _, chunk := range allChunks {
		rag.AddChunk(chunk)
	}

	// Update last watched time
	rag.LastWebWatchAt = time.Now()
	
	// Save the updated RAG
	err = ww.ragService.UpdateRag(rag)
	if err != nil {
		return 0, fmt.Errorf("error saving updated RAG: %w", err)
	}

	return len(newDocs), nil
}

// Ajouter ce au fichier file_watcher.go ou implémenter ici si c'est un nouveau fichier
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

// StartWebWatcherDaemon starts a background daemon to watch websites
func (ww *WebWatcher) StartWebWatcherDaemon(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			<-ticker.C
			ww.checkAllRags()
		}
	}()
}

// checkAllRags checks all RAGs with web watching enabled
func (ww *WebWatcher) checkAllRags() {
	// Get all RAGs
	rags, err := ww.ragService.ListAllRags()
	if err != nil {
		fmt.Printf("Error listing RAGs for web watching: %v\n", err)
		return
	}

	now := time.Now()
	
	for _, ragName := range rags {
		rag, err := ww.ragService.LoadRag(ragName)
		if err != nil {
			fmt.Printf("Error loading RAG %s: %v\n", ragName, err)
			continue
		}

		// Check if web watching is enabled and if interval has passed
		if rag.WebWatchEnabled && rag.WebWatchInterval > 0 {
			intervalDuration := time.Duration(rag.WebWatchInterval) * time.Minute
			if now.Sub(rag.LastWebWatchAt) >= intervalDuration {
				docsAdded, err := ww.CheckAndUpdateRag(rag)
				if err != nil {
					fmt.Printf("Error checking for updates in RAG %s website: %v\n", ragName, err)
				} else if docsAdded > 0 {
					fmt.Printf("Added %d new pages to RAG %s from watched website\n", docsAdded, ragName)
				}
			}
		}
	}
} 