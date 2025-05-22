package service

import (
    "fmt"
    "path/filepath"
    "strings"
    "time"

    "github.com/dontizi/rlama/internal/domain"
)

// SimpleDocumentProcessor implements DocumentProcessor using the existing DocumentLoader
// This provides a compatible interface without external langchain dependencies
type SimpleDocumentProcessor struct {
    documentLoader *DocumentLoader
    chunkService   *ChunkerService
}

// NewSimpleDocumentProcessor creates a new document processor using existing services
func NewSimpleDocumentProcessor() *SimpleDocumentProcessor {
    return &SimpleDocumentProcessor{
        documentLoader: NewDocumentLoader(),
        chunkService:   NewChunkerService(DefaultChunkingConfig()),
    }
}

// LoadDocuments loads documents from a folder path with the given options
func (sdp *SimpleDocumentProcessor) LoadDocuments(folderPath string, options DocumentLoaderOptions) ([]*domain.Document, error) {
    if folderPath == "" {
        return nil, fmt.Errorf("folder path cannot be empty")
    }

    // Validate folder path exists
    if !filepath.IsAbs(folderPath) {
        absPath, err := filepath.Abs(folderPath)
        if err != nil {
            return nil, fmt.Errorf("failed to resolve absolute path for %s: %w", folderPath, err)
        }
        folderPath = absPath
    }

    // Use the existing document loader
    documents, err := sdp.documentLoader.LoadDocumentsFromFolderWithOptions(folderPath, options)
    if err != nil {
        return nil, fmt.Errorf("failed to load documents from %s: %w", folderPath, err)
    }

    if len(documents) == 0 {
        return nil, fmt.Errorf("no documents found in folder %s", folderPath)
    }

    // Validate and clean documents
    var validDocuments []*domain.Document
    for _, doc := range documents {
        if err := sdp.validateDocument(doc); err != nil {
            fmt.Printf("Warning: Skipping invalid document %s: %v\n", doc.Name, err)
            continue
        }
        validDocuments = append(validDocuments, doc)
    }

    if len(validDocuments) == 0 {
        return nil, fmt.Errorf("no valid documents found after validation")
    }

    fmt.Printf("Successfully loaded %d valid documents from %s\n", len(validDocuments), folderPath)
    return validDocuments, nil
}

// ProcessContent processes raw content into chunks using the specified configuration
func (sdp *SimpleDocumentProcessor) ProcessContent(content string, config ChunkingConfig) ([]*domain.DocumentChunk, error) {
    if strings.TrimSpace(content) == "" {
        return nil, fmt.Errorf("content cannot be empty")
    }

    // Validate chunking configuration
    if err := sdp.validateChunkingConfig(config); err != nil {
        return nil, fmt.Errorf("invalid chunking configuration: %w", err)
    }

    // Create a temporary document for chunking
    tempDoc := &domain.Document{
        ID:      fmt.Sprintf("temp_%d", time.Now().UnixNano()),
        Content: content,
        Name:    "temporary_content",
        Path:    "",
    }

    // Use the chunker service to create chunks (update chunker config first)
    sdp.chunkService = NewChunkerService(config)
    chunks := sdp.chunkService.ChunkDocument(tempDoc)

    if len(chunks) == 0 {
        return nil, fmt.Errorf("no chunks generated from content")
    }

    return chunks, nil
}

// ProcessDocuments processes multiple documents into chunks
func (sdp *SimpleDocumentProcessor) ProcessDocuments(documents []*domain.Document, config ChunkingConfig) (map[string][]*domain.DocumentChunk, error) {
    if len(documents) == 0 {
        return nil, fmt.Errorf("no documents to process")
    }

    // Validate chunking configuration
    if err := sdp.validateChunkingConfig(config); err != nil {
        return nil, fmt.Errorf("invalid chunking configuration: %w", err)
    }

    result := make(map[string][]*domain.DocumentChunk)
    var errors []string

    for _, doc := range documents {
        if err := sdp.validateDocument(doc); err != nil {
            errors = append(errors, fmt.Sprintf("invalid document %s: %v", doc.Name, err))
            continue
        }

        // Update chunker config and chunk document
        chunkerService := NewChunkerService(config)
        chunks := chunkerService.ChunkDocument(doc)

        if len(chunks) > 0 {
            result[doc.ID] = chunks
        } else {
            errors = append(errors, fmt.Sprintf("no chunks generated for document %s", doc.Name))
        }
    }

    if len(result) == 0 {
        return nil, fmt.Errorf("no documents were successfully processed. Errors: %s", strings.Join(errors, "; "))
    }

    if len(errors) > 0 {
        fmt.Printf("Warning: Some documents failed to process: %s\n", strings.Join(errors, "; "))
    }

    fmt.Printf("Successfully processed %d documents into chunks\n", len(result))
    return result, nil
}

// EstimateChunkCount estimates how many chunks will be generated from content
func (sdp *SimpleDocumentProcessor) EstimateChunkCount(content string, config ChunkingConfig) int {
    if strings.TrimSpace(content) == "" {
        return 0
    }

    contentLength := len([]rune(content))
    if config.ChunkSize <= 0 {
        return 1
    }

    stepSize := config.ChunkSize - config.ChunkOverlap
    if stepSize <= 0 {
        stepSize = config.ChunkSize
    }

    return (contentLength + stepSize - 1) / stepSize
}

// validateDocument validates a document for processing
func (sdp *SimpleDocumentProcessor) validateDocument(doc *domain.Document) error {
    if doc == nil {
        return fmt.Errorf("document is nil")
    }

    if doc.ID == "" {
        return fmt.Errorf("document ID is empty")
    }

    if strings.TrimSpace(doc.Content) == "" {
        return fmt.Errorf("document content is empty")
    }

    if doc.Name == "" {
        return fmt.Errorf("document name is empty")
    }

    // Check for minimum content length
    if len(strings.TrimSpace(doc.Content)) < 10 {
        return fmt.Errorf("document content too short (minimum 10 characters)")
    }

    return nil
}

// validateChunkingConfig validates chunking configuration
func (sdp *SimpleDocumentProcessor) validateChunkingConfig(config ChunkingConfig) error {
    if config.ChunkSize <= 0 {
        return fmt.Errorf("chunk size must be positive, got %d", config.ChunkSize)
    }

    if config.ChunkSize < 50 {
        return fmt.Errorf("chunk size too small (minimum 50 characters), got %d", config.ChunkSize)
    }

    if config.ChunkOverlap < 0 {
        return fmt.Errorf("chunk overlap cannot be negative, got %d", config.ChunkOverlap)
    }

    if config.ChunkOverlap >= config.ChunkSize {
        return fmt.Errorf("chunk overlap (%d) must be less than chunk size (%d)", config.ChunkOverlap, config.ChunkSize)
    }

    validStrategies := map[string]bool{
        "fixed":        true,
        "semantic":     true,
        "hybrid":       true,
        "hierarchical": true,
    }

    if config.ChunkingStrategy != "" && !validStrategies[config.ChunkingStrategy] {
        return fmt.Errorf("invalid chunking strategy: %s. Valid options: fixed, semantic, hybrid, hierarchical", config.ChunkingStrategy)
    }

    return nil
}

// GetSupportedFileTypes returns the file types supported by this processor
func (sdp *SimpleDocumentProcessor) GetSupportedFileTypes() []string {
    return []string{
        ".txt", ".md", ".html", ".htm", ".json", ".csv", ".log", ".xml", ".yaml", ".yml",
        ".go", ".py", ".js", ".java", ".c", ".cpp", ".cxx", ".f", ".F", ".F90", ".h", 
        ".rb", ".php", ".rs", ".swift", ".kt", ".el", ".svelte", ".ts", ".tsx",
        ".pdf", ".docx", ".doc", ".rtf", ".odt", ".pptx", ".ppt", ".xlsx", ".xls", 
        ".epub", ".org",
    }
}

// GetProcessingStatistics returns statistics about the last processing operation
func (sdp *SimpleDocumentProcessor) GetProcessingStatistics() map[string]interface{} {
    return map[string]interface{}{
        "processor_type":     "simple",
        "supported_formats":  len(sdp.GetSupportedFileTypes()),
        "last_updated":      time.Now(),
        "features": []string{
            "text_extraction",
            "pdf_processing", 
            "document_chunking",
            "content_validation",
            "parallel_processing",
        },
    }
}