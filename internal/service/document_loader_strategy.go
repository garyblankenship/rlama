package service

import (
	"github.com/dontizi/rlama/internal/domain"
)

// DocumentLoaderStrategy defines the interface for different document loading strategies
type DocumentLoaderStrategy interface {
	LoadDocuments(folderPath string, options DocumentLoaderOptions) ([]*domain.Document, error)
	GetName() string
	GetSupportedFileTypes() []string
	IsAvailable() bool
}

// LegacyDocumentLoaderStrategy wraps the existing DocumentLoader
type LegacyDocumentLoaderStrategy struct {
	loader *DocumentLoader
}

// NewLegacyDocumentLoaderStrategy creates a new legacy loader strategy
func NewLegacyDocumentLoaderStrategy() *LegacyDocumentLoaderStrategy {
	return &LegacyDocumentLoaderStrategy{
		loader: NewDocumentLoader(),
	}
}

// LoadDocuments implements DocumentLoaderStrategy for the legacy loader
func (l *LegacyDocumentLoaderStrategy) LoadDocuments(folderPath string, options DocumentLoaderOptions) ([]*domain.Document, error) {
	return l.loader.LoadDocumentsFromFolderWithOptions(folderPath, options)
}

// GetName returns the strategy name
func (l *LegacyDocumentLoaderStrategy) GetName() string {
	return "legacy"
}

// GetSupportedFileTypes returns supported file extensions
func (l *LegacyDocumentLoaderStrategy) GetSupportedFileTypes() []string {
	return []string{
		".txt", ".md", ".html", ".htm", ".json", ".csv", ".log", ".xml", ".yaml", ".yml",
		".go", ".py", ".js", ".java", ".c", ".cpp", ".cxx", ".f", ".F", ".F90", ".h",
		".rb", ".php", ".rs", ".swift", ".kt", ".el", ".svelte", ".ts", ".tsx",
		".pdf", ".docx", ".doc", ".rtf", ".odt", ".pptx", ".ppt", ".xlsx", ".xls",
		".epub", ".org",
	}
}

// IsAvailable checks if the legacy loader is available (always true)
func (l *LegacyDocumentLoaderStrategy) IsAvailable() bool {
	return true
}