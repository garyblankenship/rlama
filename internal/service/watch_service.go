package service

import (
	"fmt"

	"github.com/dontizi/rlama/internal/domain"
)

// WatchService handles file and web monitoring for RAG systems
type WatchService interface {
	// SetupDirectoryWatching enables monitoring of a directory for changes
	SetupDirectoryWatching(ragName string, dirPath string, watchInterval int, options DocumentLoaderOptions) error
	
	// DisableDirectoryWatching disables directory monitoring for a RAG system
	DisableDirectoryWatching(ragName string) error
	
	// CheckWatchedDirectory checks for changes in the watched directory and returns count of new documents
	CheckWatchedDirectory(ragName string) (int, error)
	
	// SetupWebWatching enables monitoring of a website for changes
	SetupWebWatching(ragName string, websiteURL string, watchInterval int, options domain.WebWatchOptions) error
	
	// DisableWebWatching disables web monitoring for a RAG system
	DisableWebWatching(ragName string) error
	
	// CheckWatchedWebsite checks for changes on the watched website and returns count of new documents
	CheckWatchedWebsite(ragName string) (int, error)
}

// WatchServiceImpl implements the WatchService interface
type WatchServiceImpl struct {
	documentService DocumentService
	ragService      RagService
	fileWatcher     *FileWatcher
	webWatcher      *WebWatcher
}

// NewWatchService creates a new WatchService instance
func NewWatchService(documentService DocumentService, ragService RagService) WatchService {
	return &WatchServiceImpl{
		documentService: documentService,
		ragService:      ragService,
		fileWatcher:     NewFileWatcher(ragService),
		webWatcher:      NewWebWatcher(ragService),
	}
}

// SetupDirectoryWatching implements WatchService.SetupDirectoryWatching
func (ws *WatchServiceImpl) SetupDirectoryWatching(ragName string, dirPath string, watchInterval int, options DocumentLoaderOptions) error {
	// Load the RAG system
	rag, err := ws.documentService.LoadRAG(ragName)
	if err != nil {
		return fmt.Errorf("failed to load RAG '%s': %w", ragName, err)
	}

	// Configure directory watching
	rag.WatchEnabled = true
	rag.WatchedDir = dirPath
	rag.WatchInterval = watchInterval
	rag.WatchOptions = domain.DocumentWatchOptions{
		ExcludeDirs:      options.ExcludeDirs,
		ExcludeExts:      options.ExcludeExts,
		ProcessExts:      options.ProcessExts,
		ChunkSize:        options.ChunkSize,
		ChunkOverlap:     options.ChunkOverlap,
		ChunkingStrategy: options.ChunkingStrategy,
	}

	// Save the updated RAG
	if err := ws.documentService.UpdateRAG(rag); err != nil {
		return fmt.Errorf("failed to update RAG '%s': %w", ragName, err)
	}

	fmt.Printf("Directory watching enabled for RAG '%s' on directory '%s'\n", ragName, dirPath)
	return nil
}

// DisableDirectoryWatching implements WatchService.DisableDirectoryWatching
func (ws *WatchServiceImpl) DisableDirectoryWatching(ragName string) error {
	// Load the RAG system
	rag, err := ws.documentService.LoadRAG(ragName)
	if err != nil {
		return fmt.Errorf("failed to load RAG '%s': %w", ragName, err)
	}

	// Disable directory watching
	rag.WatchEnabled = false
	rag.WatchedDir = ""
	rag.WatchInterval = 0

	// Save the updated RAG
	if err := ws.documentService.UpdateRAG(rag); err != nil {
		return fmt.Errorf("failed to update RAG '%s': %w", ragName, err)
	}

	fmt.Printf("Directory watching disabled for RAG '%s'\n", ragName)
	return nil
}

// CheckWatchedDirectory implements WatchService.CheckWatchedDirectory
func (ws *WatchServiceImpl) CheckWatchedDirectory(ragName string) (int, error) {
	// Load the RAG system
	rag, err := ws.documentService.LoadRAG(ragName)
	if err != nil {
		return 0, fmt.Errorf("failed to load RAG '%s': %w", ragName, err)
	}

	// Check if directory watching is enabled
	if !rag.WatchEnabled {
		return 0, fmt.Errorf("directory watching is not enabled for RAG '%s'", ragName)
	}

	// Use the file watcher to check for changes
	return ws.fileWatcher.CheckAndUpdateRag(rag)
}

// SetupWebWatching implements WatchService.SetupWebWatching
func (ws *WatchServiceImpl) SetupWebWatching(ragName string, websiteURL string, watchInterval int, options domain.WebWatchOptions) error {
	// Load the RAG system
	rag, err := ws.documentService.LoadRAG(ragName)
	if err != nil {
		return fmt.Errorf("failed to load RAG '%s': %w", ragName, err)
	}

	// Configure web watching
	rag.WebWatchEnabled = true
	rag.WatchedURL = websiteURL
	rag.WebWatchInterval = watchInterval
	rag.WebWatchOptions = options

	// Save the updated RAG
	if err := ws.documentService.UpdateRAG(rag); err != nil {
		return fmt.Errorf("failed to update RAG '%s': %w", ragName, err)
	}

	fmt.Printf("Web watching enabled for RAG '%s' on URL '%s'\n", ragName, websiteURL)
	return nil
}

// DisableWebWatching implements WatchService.DisableWebWatching
func (ws *WatchServiceImpl) DisableWebWatching(ragName string) error {
	// Load the RAG system
	rag, err := ws.documentService.LoadRAG(ragName)
	if err != nil {
		return fmt.Errorf("failed to load RAG '%s': %w", ragName, err)
	}

	// Disable web watching
	rag.WebWatchEnabled = false
	rag.WatchedURL = ""
	rag.WebWatchInterval = 0

	// Save the updated RAG
	if err := ws.documentService.UpdateRAG(rag); err != nil {
		return fmt.Errorf("failed to update RAG '%s': %w", ragName, err)
	}

	fmt.Printf("Web watching disabled for RAG '%s'\n", ragName)
	return nil
}

// CheckWatchedWebsite implements WatchService.CheckWatchedWebsite
func (ws *WatchServiceImpl) CheckWatchedWebsite(ragName string) (int, error) {
	// Load the RAG system
	rag, err := ws.documentService.LoadRAG(ragName)
	if err != nil {
		return 0, fmt.Errorf("failed to load RAG '%s': %w", ragName, err)
	}

	// Check if web watching is enabled
	if !rag.WebWatchEnabled {
		return 0, fmt.Errorf("web watching is not enabled for RAG '%s'", ragName)
	}

	// Use the web watcher to check for changes
	return ws.webWatcher.CheckAndUpdateRag(rag)
}