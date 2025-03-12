package service

import (
	"fmt"
	"strings"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/repository"
)

// Remove duplicate interface declaration and keep this one
type RagService interface {
	CreateRagWithOptions(modelName, ragName, folderPath string, options DocumentLoaderOptions) error
	GetRagChunks(ragName string, filter ChunkFilter) ([]*domain.DocumentChunk, error)
	LoadRag(ragName string) (*domain.RagSystem, error)
	Query(rag *domain.RagSystem, query string, contextSize int) (string, error)
	AddDocsWithOptions(ragName string, folderPath string, options DocumentLoaderOptions) error
	UpdateModel(ragName string, newModel string) error
	UpdateRag(rag *domain.RagSystem) error
	// Add any other required methods here
}

// Update the struct implementation to match the interface
type RagServiceImpl struct {
	documentLoader   *DocumentLoader
	embeddingService *EmbeddingService
	ragRepository    *repository.RagRepository
	ollamaClient     *client.OllamaClient
}

// NewRagService creates a new instance of RagService
func NewRagService(ollamaClient *client.OllamaClient) RagService {
	if ollamaClient == nil {
		ollamaClient = client.NewDefaultOllamaClient()
	}
	
	return &RagServiceImpl{
		documentLoader:   NewDocumentLoader(),
		embeddingService: NewEmbeddingService(ollamaClient),
		ragRepository:    repository.NewRagRepository(),
		ollamaClient:     ollamaClient,
	}
}

// CreateRagWithOptions creates a new RAG system with the specified options
func (rs *RagServiceImpl) CreateRagWithOptions(modelName, ragName, folderPath string, options DocumentLoaderOptions) error {
	// Check if Ollama is available
	if err := rs.ollamaClient.CheckOllamaAndModel(modelName); err != nil {
		return err
	}

	// Check if the RAG already exists
	if rs.ragRepository.Exists(ragName) {
		return fmt.Errorf("a RAG with name '%s' already exists", ragName)
	}

	// Load documents with options
	docs, err := rs.documentLoader.LoadDocumentsFromFolderWithOptions(folderPath, options)
	if err != nil {
		return fmt.Errorf("error loading documents: %w", err)
	}

	if len(docs) == 0 {
		return fmt.Errorf("no valid documents found in folder %s", folderPath)
	}

	fmt.Printf("Successfully loaded %d documents. Chunking documents...\n", len(docs))
	
	// Create the RAG system
	rag := domain.NewRagSystem(ragName, modelName)
	
	// Create chunker service
	chunkerService := NewChunkerService(ChunkingConfig{
		ChunkSize:    options.ChunkSize,
		ChunkOverlap: options.ChunkOverlap,
	})
	
	// Process each document - chunk and generate embeddings
	var allChunks []*domain.DocumentChunk
	for _, doc := range docs {
		// Add the document to the RAG
		rag.AddDocument(doc)
		
		// Chunk the document
		chunks := chunkerService.ChunkDocument(doc)
		
		// Update total chunks in metadata
		for i, chunk := range chunks {
			chunk.ChunkNumber = i
			chunk.TotalChunks = len(chunks)
		}
		
		allChunks = append(allChunks, chunks...)
	}
	
	fmt.Printf("Generated %d chunks from %d documents. Generating embeddings...\n", 
		len(allChunks), len(docs))
	
	// Generate embeddings for all chunks
	err = rs.embeddingService.GenerateChunkEmbeddings(allChunks, modelName)
	if err != nil {
		return fmt.Errorf("error generating embeddings: %w", err)
	}
	
	// Add all chunks to the RAG
	for _, chunk := range allChunks {
		rag.AddChunk(chunk)
	}

	// Save the RAG
	err = rs.ragRepository.Save(rag)
	if err != nil {
		return fmt.Errorf("error saving the RAG: %w", err)
	}

	fmt.Printf("RAG created with %d indexed documents (%d chunks).\n", len(docs), len(allChunks))
	return nil
}

// Modify the existing CreateRag to use CreateRagWithOptions
func (rs *RagServiceImpl) CreateRag(modelName, ragName, folderPath string) error {
	return rs.CreateRagWithOptions(modelName, ragName, folderPath, DocumentLoaderOptions{})
}

// LoadRag loads a RAG system
func (rs *RagServiceImpl) LoadRag(ragName string) (*domain.RagSystem, error) {
	rag, err := rs.ragRepository.Load(ragName)
	if err != nil {
		return nil, fmt.Errorf("error loading RAG '%s': %w", ragName, err)
	}

	return rag, nil
}

// Query performs a query on a RAG system
func (rs *RagServiceImpl) Query(rag *domain.RagSystem, query string, contextSize int) (string, error) {
	// Check if Ollama is available
	if err := rs.ollamaClient.CheckOllamaAndModel(rag.ModelName); err != nil {
		return "", err
	}

	// Generate embedding for the query
	queryEmbedding, err := rs.embeddingService.GenerateQueryEmbedding(query, rag.ModelName)
	if err != nil {
		return "", fmt.Errorf("error generating embedding for query: %w", err)
	}

	// Use the provided context size or default to 20
	if contextSize <= 0 {
		contextSize = 20
	}
	
	// Search for the most relevant chunks
	results := rag.HybridStore.Search(queryEmbedding, contextSize)
	
	// Build the context
	var context strings.Builder
	context.WriteString("Relevant information:\n\n")
	
	// Track which documents we've included for reference
	includedDocs := make(map[string]bool)
	
	for _, result := range results {
		chunk := rag.GetChunkByID(result.ID)
		if chunk != nil {
			// Add chunk content with its metadata
			context.WriteString(fmt.Sprintf("--- %s ---\n%s\n\n", 
				chunk.GetMetadataString(), chunk.Content))
				
			includedDocs[chunk.DocumentID] = true
		}
	}
	
	// Build the prompt with better formatting and instructions for citing sources
	prompt := fmt.Sprintf(`You are a helpful AI assistant. Use the information below to answer the question.

%s

Question: %s

Answer based on the provided information. If the information doesn't contain the answer, say so clearly.
Include references to the source documents in your answer using the format (Source: document name).`, 
	context.String(), query)
	
	// Show search results to the user
	fmt.Println("\nSearching documents...\n")
	fmt.Printf("Found %d relevant sections across %d documents\n", 
		len(results), len(includedDocs))
	
	// Generate the response
	response, err := rs.ollamaClient.GenerateCompletion(rag.ModelName, prompt)
	if err != nil {
		return "", fmt.Errorf("error generating response: %w", err)
	}
	
	return response, nil
}

// UpdateRag updates an existing RAG system
func (rs *RagServiceImpl) UpdateRag(rag *domain.RagSystem) error {
	err := rs.ragRepository.Save(rag)
	if err != nil {
		return fmt.Errorf("error updating the RAG: %w", err)
	}
	return nil
}

// AddDocsWithOptions adds documents to an existing RAG system with options
func (rs *RagServiceImpl) AddDocsWithOptions(ragName string, folderPath string, options DocumentLoaderOptions) error {
	// Load existing RAG
	rag, err := rs.LoadRag(ragName)
	if err != nil {
		return err
	}

	// Load documents with options
	docs, err := rs.documentLoader.LoadDocumentsFromFolderWithOptions(folderPath, options)
	if err != nil {
		return fmt.Errorf("error loading documents: %w", err)
	}

	if len(docs) == 0 {
		return fmt.Errorf("no valid documents found in folder %s", folderPath)
	}

	// Create chunker service with default config
	chunkerService := NewChunkerService(DefaultChunkingConfig())
	
	// Process documents
	var allChunks []*domain.DocumentChunk
	for _, doc := range docs {
		chunks := chunkerService.ChunkDocument(doc)
		for _, chunk := range chunks {
			chunk.UpdateTotalChunks(len(chunks))
		}
		allChunks = append(allChunks, chunks...)
	}

	// Generate embeddings
	err = rs.embeddingService.GenerateChunkEmbeddings(allChunks, rag.ModelName)
	if err != nil {
		return err
	}

	// Add new chunks
	chunksAdded := 0
	existingChunks := make(map[string]bool)
	for _, chunk := range rag.Chunks {
		existingChunks[chunk.ID] = true
	}
	for _, chunk := range allChunks {
		if !existingChunks[chunk.ID] {
			rag.AddChunk(chunk)
			chunksAdded++
		}
	}

	// Save updated RAG
	return rs.UpdateRag(rag)
}

// Add chunk filter struct
type ChunkFilter struct {
	DocumentSubstring string
	ShowContent       bool
}

func (rs *RagServiceImpl) GetRagChunks(ragName string, filter ChunkFilter) ([]*domain.DocumentChunk, error) {
	rag, err := rs.ragRepository.Load(ragName)
	if err != nil {
		return nil, fmt.Errorf("error loading RAG: %w", err)
	}

	var filtered []*domain.DocumentChunk
	for _, chunk := range rag.Chunks {
		// Apply document filter
		if filter.DocumentSubstring != "" && 
		   !strings.Contains(strings.ToLower(chunk.DocumentID), strings.ToLower(filter.DocumentSubstring)) {
			continue
		}

		// Clone chunk to avoid modifying original
		c := *chunk
		
		// Clear content if not requested
		if !filter.ShowContent {
			c.Content = ""
		}

		filtered = append(filtered, &c)
	}

	return filtered, nil
}

// UpdateModel updates the model of an existing RAG system
func (rs *RagServiceImpl) UpdateModel(ragName string, newModel string) error {
	rag, err := rs.LoadRag(ragName)
	if err != nil {
		return fmt.Errorf("error loading RAG: %w", err)
	}

	rag.ModelName = newModel
	return rs.UpdateRag(rag)
} 