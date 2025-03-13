package service

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dontizi/rlama/internal/domain"
)

// FileWatcher is responsible for watching directories for file changes
type FileWatcher struct {
	ragService RagService
}

// NewFileWatcher creates a new file watcher service
func NewFileWatcher(ragService RagService) *FileWatcher {
	return &FileWatcher{
		ragService: ragService,
	}
}

// CheckAndUpdateRag checks for new files in the watched directory and updates the RAG
func (fw *FileWatcher) CheckAndUpdateRag(rag *domain.RagSystem) (int, error) {
	if !rag.WatchEnabled || rag.WatchedDir == "" {
		return 0, nil // Watching not enabled
	}

	// Check if the directory exists
	dirInfo, err := os.Stat(rag.WatchedDir)
	if os.IsNotExist(err) {
		return 0, fmt.Errorf("watched directory '%s' does not exist", rag.WatchedDir)
	} else if err != nil {
		return 0, fmt.Errorf("error accessing watched directory: %w", err)
	}

	if !dirInfo.IsDir() {
		return 0, fmt.Errorf("'%s' is not a directory", rag.WatchedDir)
	}

	// Get the last modified time of the directory
	lastModified := getLastModifiedTime(rag.WatchedDir)
	
	// If the directory hasn't been modified since last check, no need to proceed
	if !lastModified.After(rag.LastWatchedAt) && !rag.LastWatchedAt.IsZero() {
		return 0, nil
	}

	// Convert watch options to document loader options
	loaderOptions := DocumentLoaderOptions{
		ExcludeDirs:  rag.WatchOptions.ExcludeDirs,
		ExcludeExts:  rag.WatchOptions.ExcludeExts,
		ProcessExts:  rag.WatchOptions.ProcessExts,
		ChunkSize:    rag.WatchOptions.ChunkSize,
		ChunkOverlap: rag.WatchOptions.ChunkOverlap,
	}

	// Get existing document paths to avoid re-processing
	existingPaths := make(map[string]bool)
	for _, doc := range rag.Documents {
		existingPaths[doc.Path] = true
	}

	// Create a document loader
	docLoader := NewDocumentLoader()
	
	// Load all documents from the directory
	allDocs, err := docLoader.LoadDocumentsFromFolderWithOptions(rag.WatchedDir, loaderOptions)
	if err != nil {
		return 0, fmt.Errorf("error loading documents from watched directory: %w", err)
	}

	// Filter out existing documents
	var newDocs []*domain.Document
	for _, doc := range allDocs {
		if !existingPaths[doc.Path] {
			newDocs = append(newDocs, doc)
		}
	}

	if len(newDocs) == 0 {
		// Update last watched time even if no new documents
		rag.LastWatchedAt = time.Now()
		err = fw.ragService.UpdateRag(rag)
		return 0, err
	}

	// Create chunker service with options from the RAG
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
	embeddingService := NewEmbeddingService(fw.ragService.GetOllamaClient())
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
	rag.LastWatchedAt = time.Now()
	
	// Save the updated RAG
	err = fw.ragService.UpdateRag(rag)
	if err != nil {
		return 0, fmt.Errorf("error saving updated RAG: %w", err)
	}

	return len(newDocs), nil
}

// getLastModifiedTime gets the latest modification time in a directory
func getLastModifiedTime(dirPath string) time.Time {
	var lastModTime time.Time

	// Walk through the directory and find the latest modification time
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		
		if info.ModTime().After(lastModTime) {
			lastModTime = info.ModTime()
		}
		
		return nil
	})

	return lastModTime
}

// StartWatcherDaemon starts a background daemon to watch directories
// for all RAGs that have watching enabled with intervals
func (fw *FileWatcher) StartWatcherDaemon(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			<-ticker.C
			fw.checkAllRags()
		}
	}()
}

// checkAllRags checks all RAGs with watching enabled
func (fw *FileWatcher) checkAllRags() {
	// Get all RAGs
	rags, err := fw.ragService.ListAllRags()
	if err != nil {
		fmt.Printf("Error listing RAGs for file watching: %v\n", err)
		return
	}

	now := time.Now()
	
	for _, ragName := range rags {
		rag, err := fw.ragService.LoadRag(ragName)
		if err != nil {
			fmt.Printf("Error loading RAG %s: %v\n", ragName, err)
			continue
		}

		// Check if watching is enabled and if interval has passed
		if rag.WatchEnabled && rag.WatchInterval > 0 {
			intervalDuration := time.Duration(rag.WatchInterval) * time.Minute
			if now.Sub(rag.LastWatchedAt) >= intervalDuration {
				docsAdded, err := fw.CheckAndUpdateRag(rag)
				if err != nil {
					fmt.Printf("Error checking for updates in RAG %s: %v\n", ragName, err)
				} else if docsAdded > 0 {
					fmt.Printf("Added %d new documents to RAG %s from watched directory\n", docsAdded, ragName)
				}
			}
		}
	}
} 