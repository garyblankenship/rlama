package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tmc/langchaingo/documentloaders"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/textsplitter"

	"github.com/dontizi/rlama/internal/domain"
)

// LangChainDocumentLoaderStrategy implements document loading using LangChainGo with available APIs
type LangChainDocumentLoaderStrategy struct {
	textSplitter textsplitter.TextSplitter
	ctx          context.Context
	timeout      time.Duration
}

// NewLangChainDocumentLoaderStrategy creates a new LangChain-based loader strategy
func NewLangChainDocumentLoaderStrategy() *LangChainDocumentLoaderStrategy {
	return &LangChainDocumentLoaderStrategy{
		textSplitter: textsplitter.NewRecursiveCharacter(
			textsplitter.WithChunkSize(1000),
			textsplitter.WithChunkOverlap(200),
		),
		ctx:     context.Background(),
		timeout: 5 * time.Minute,
	}
}

// LoadDocuments implements DocumentLoaderStrategy using LangChain with manual directory walking
func (lcl *LangChainDocumentLoaderStrategy) LoadDocuments(folderPath string, options DocumentLoaderOptions) ([]*domain.Document, error) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(lcl.ctx, lcl.timeout)
	defer cancel()

	// Find all files in the directory using our own walking
	files, err := lcl.findFiles(folderPath, options)
	if err != nil {
		return nil, fmt.Errorf("failed to find files in directory: %w", err)
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no supported files found in %s", folderPath)
	}

	// Load documents from files using appropriate LangChain loaders
	var allDocs []schema.Document
	var errors []string

	for _, file := range files {
		docs, err := lcl.loadFile(ctx, file)
		if err != nil {
			errors = append(errors, fmt.Sprintf("failed to load %s: %v", file, err))
			continue
		}
		allDocs = append(allDocs, docs...)
	}

	if len(allDocs) == 0 {
		return nil, fmt.Errorf("no documents loaded successfully. Errors: %s", strings.Join(errors, "; "))
	}

	// Convert to RLAMA format
	rlamaDocs, err := lcl.convertToRLAMADocuments(allDocs, folderPath)
	if err != nil {
		return nil, fmt.Errorf("document conversion failed: %w", err)
	}

	if len(errors) > 0 {
		fmt.Printf("⚠️ Some files failed to load: %s\n", strings.Join(errors, "; "))
	}

	fmt.Printf("✅ LangChain loaded %d documents from %s\n", len(rlamaDocs), folderPath)
	return rlamaDocs, nil
}

// findFiles walks the directory and finds files to process
func (lcl *LangChainDocumentLoaderStrategy) findFiles(folderPath string, options DocumentLoaderOptions) ([]string, error) {
	var files []string

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			// Check if this directory should be excluded
			for _, excludeDir := range options.ExcludeDirs {
				if strings.Contains(path, excludeDir) {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Skip hidden files
		if strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))

		// Apply extension filters
		if lcl.shouldExcludeByExt(path, options.ExcludeExts, options.ProcessExts) {
			return nil
		}

		// Check if we support this file type
		if lcl.isSupportedFileType(ext) {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}

// loadFile loads a single file using the appropriate LangChain loader
func (lcl *LangChainDocumentLoaderStrategy) loadFile(ctx context.Context, filePath string) ([]schema.Document, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".txt", ".md", ".log", ".go", ".py", ".js", ".java", ".c", ".cpp", ".h", ".rb", ".php", ".rs", ".swift", ".kt", ".ts", ".tsx", ".jsx":
		return lcl.loadTextFile(ctx, filePath)
	case ".pdf":
		return lcl.loadPDFFile(ctx, filePath)
	case ".html", ".htm":
		return lcl.loadHTMLFile(ctx, filePath)
	case ".csv":
		return lcl.loadCSVFile(ctx, filePath)
	default:
		// Fallback to text loading
		return lcl.loadTextFile(ctx, filePath)
	}
}

// loadTextFile loads a text file using LangChain Text loader
func (lcl *LangChainDocumentLoaderStrategy) loadTextFile(ctx context.Context, filePath string) ([]schema.Document, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	loader := documentloaders.NewText(file)
	docs, err := loader.Load(ctx)
	if err != nil {
		return nil, err
	}

	// Add source metadata
	for i := range docs {
		if docs[i].Metadata == nil {
			docs[i].Metadata = make(map[string]any)
		}
		docs[i].Metadata["source"] = filePath
	}

	return docs, nil
}

// loadPDFFile loads a PDF file using LangChain PDF loader
func (lcl *LangChainDocumentLoaderStrategy) loadPDFFile(ctx context.Context, filePath string) ([]schema.Document, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get file size
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	loader := documentloaders.NewPDF(file, fileInfo.Size())
	docs, err := loader.Load(ctx)
	if err != nil {
		return nil, err
	}

	// Add source metadata
	for i := range docs {
		if docs[i].Metadata == nil {
			docs[i].Metadata = make(map[string]any)
		}
		docs[i].Metadata["source"] = filePath
	}

	return docs, nil
}

// loadHTMLFile loads an HTML file using LangChain HTML loader
func (lcl *LangChainDocumentLoaderStrategy) loadHTMLFile(ctx context.Context, filePath string) ([]schema.Document, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	loader := documentloaders.NewHTML(file)
	docs, err := loader.Load(ctx)
	if err != nil {
		return nil, err
	}

	// Add source metadata
	for i := range docs {
		if docs[i].Metadata == nil {
			docs[i].Metadata = make(map[string]any)
		}
		docs[i].Metadata["source"] = filePath
	}

	return docs, nil
}

// loadCSVFile loads a CSV file using LangChain CSV loader
func (lcl *LangChainDocumentLoaderStrategy) loadCSVFile(ctx context.Context, filePath string) ([]schema.Document, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	loader := documentloaders.NewCSV(file)
	docs, err := loader.Load(ctx)
	if err != nil {
		return nil, err
	}

	// Add source metadata
	for i := range docs {
		if docs[i].Metadata == nil {
			docs[i].Metadata = make(map[string]any)
		}
		docs[i].Metadata["source"] = filePath
	}

	return docs, nil
}

// Helper methods (reused from original implementation)

// GetName returns the strategy name
func (lcl *LangChainDocumentLoaderStrategy) GetName() string {
	return "langchain"
}

// GetSupportedFileTypes returns file types supported by LangChain
func (lcl *LangChainDocumentLoaderStrategy) GetSupportedFileTypes() []string {
	return []string{
		".txt", ".md", ".html", ".htm", ".json", ".csv", ".log", ".xml", ".yaml", ".yml",
		".go", ".py", ".js", ".java", ".c", ".cpp", ".h", ".rb", ".php", ".rs", ".swift",
		".kt", ".ts", ".tsx", ".jsx", ".vue", ".svelte",
		".pdf", ".docx", ".doc", ".rtf", ".odt",
		".epub", ".org",
	}
}

// IsAvailable checks if LangChain loader dependencies are available
func (lcl *LangChainDocumentLoaderStrategy) IsAvailable() bool {
	defer func() {
		if r := recover(); r != nil {
			// If LangChain panics, consider it unavailable
		}
	}()

	// Test if we can create a basic text splitter
	_ = textsplitter.NewRecursiveCharacter()
	return true
}

// isSupportedFileType checks if a file extension is supported
func (lcl *LangChainDocumentLoaderStrategy) isSupportedFileType(ext string) bool {
	supportedTypes := lcl.GetSupportedFileTypes()
	for _, supportedExt := range supportedTypes {
		if ext == supportedExt {
			return true
		}
	}
	return false
}

// shouldExcludeByExt checks if a file should be excluded based on extension filters
func (lcl *LangChainDocumentLoaderStrategy) shouldExcludeByExt(filePath string, excludeExts, processExts []string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))

	// Check exclude list first
	for _, excludeExt := range excludeExts {
		if !strings.HasPrefix(excludeExt, ".") {
			excludeExt = "." + excludeExt
		}
		if ext == excludeExt {
			return true
		}
	}

	// If processExts is specified, only include those extensions
	if len(processExts) > 0 {
		for _, processExt := range processExts {
			if !strings.HasPrefix(processExt, ".") {
				processExt = "." + processExt
			}
			if ext == processExt {
				return false
			}
		}
		return true // Exclude if not in processExts list
	}

	return false
}

// convertToRLAMADocuments converts LangChain documents to RLAMA domain objects
func (lcl *LangChainDocumentLoaderStrategy) convertToRLAMADocuments(lcDocs []schema.Document, basePath string) ([]*domain.Document, error) {
	var rlamaDocs []*domain.Document
	var errors []string

	for i, lcDoc := range lcDocs {
		rlamaDoc, err := lcl.convertSingleDocument(lcDoc, basePath, i)
		if err != nil {
			errors = append(errors, fmt.Sprintf("document %d: %v", i, err))
			continue
		}
		rlamaDocs = append(rlamaDocs, rlamaDoc)
	}

	if len(rlamaDocs) == 0 && len(errors) > 0 {
		return nil, fmt.Errorf("failed to convert any documents: %s", strings.Join(errors, "; "))
	}

	if len(errors) > 0 {
		fmt.Printf("⚠️ Some documents failed conversion: %s\n", strings.Join(errors, "; "))
	}

	return rlamaDocs, nil
}

// convertSingleDocument converts a single LangChain document to RLAMA format
func (lcl *LangChainDocumentLoaderStrategy) convertSingleDocument(lcDoc schema.Document, basePath string, index int) (*domain.Document, error) {
	// Extract source path
	sourcePath, ok := lcDoc.Metadata["source"].(string)
	if !ok || sourcePath == "" {
		return nil, fmt.Errorf("missing or invalid source path")
	}

	// Create relative path for better identification
	relPath, err := filepath.Rel(basePath, sourcePath)
	if err != nil {
		relPath = filepath.Base(sourcePath)
	}

	// Validate content
	content := strings.TrimSpace(lcDoc.PageContent)
	if len(content) == 0 {
		return nil, fmt.Errorf("empty content")
	}

	// Generate document ID
	docID := lcl.generateDocumentID(sourcePath, index)

	// Create RLAMA document
	rlamaDoc := &domain.Document{
		ID:          docID,
		Path:        sourcePath,
		Name:        relPath,
		Content:     content,
		CreatedAt:   time.Now(),
		ContentType: lcl.guessContentType(sourcePath),
		Size:        int64(len(content)),
		Metadata:    lcl.convertMetadata(lcDoc.Metadata),
	}

	return rlamaDoc, nil
}

// generateDocumentID creates a unique document ID
func (lcl *LangChainDocumentLoaderStrategy) generateDocumentID(sourcePath string, index int) string {
	baseName := strings.ReplaceAll(filepath.Base(sourcePath), ".", "_")
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("lc_%s_%d_%d", baseName, index, timestamp)
}

// guessContentType determines content type from file extension
func (lcl *LangChainDocumentLoaderStrategy) guessContentType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	
	contentTypes := map[string]string{
		".txt":   "text/plain",
		".md":    "text/markdown",
		".html":  "text/html",
		".htm":   "text/html",
		".json":  "application/json",
		".xml":   "application/xml",
		".yaml":  "application/yaml",
		".yml":   "application/yaml",
		".csv":   "text/csv",
		".pdf":   "application/pdf",
		".go":    "text/x-go",
		".py":    "text/x-python",
		".js":    "application/javascript",
		".ts":    "application/typescript",
		".java":  "text/x-java-source",
		".c":     "text/x-c",
		".cpp":   "text/x-c++",
		".h":     "text/x-c-header",
		".rb":    "text/x-ruby",
		".php":   "application/x-httpd-php",
		".rs":    "text/x-rust",
		".swift": "text/x-swift",
		".kt":    "text/x-kotlin",
	}

	if contentType, exists := contentTypes[ext]; exists {
		return contentType
	}

	return "application/octet-stream"
}

// convertMetadata converts LangChain metadata to RLAMA metadata string
func (lcl *LangChainDocumentLoaderStrategy) convertMetadata(metadata map[string]any) string {
	var parts []string
	
	for key, value := range metadata {
		if key == "source" {
			continue // Skip source as it's handled separately
		}
		parts = append(parts, fmt.Sprintf("%s: %v", key, value))
	}
	
	return strings.Join(parts, ", ")
}