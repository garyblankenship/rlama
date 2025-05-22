package service

import (
	"strings"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/repository"
)

// DocumentService handles document management operations for RAG systems
type DocumentService interface {
	// CreateRAG creates a new RAG system with documents from the specified folder
	CreateRAG(modelName, ragName, folderPath string, options DocumentLoaderOptions) error
	
	// AddDocuments adds documents from a folder to an existing RAG system
	AddDocuments(ragName string, folderPath string, options DocumentLoaderOptions) error
	
	// GetChunks retrieves document chunks from a RAG system with optional filtering
	GetChunks(ragName string, filter ChunkFilter) ([]*domain.DocumentChunk, error)
	
	// LoadRAG loads a RAG system from the repository
	LoadRAG(ragName string) (*domain.RagSystem, error)
	
	// UpdateRAG saves changes to a RAG system
	UpdateRAG(rag *domain.RagSystem) error
	
	// ListRAGs returns all available RAG system names
	ListRAGs() ([]string, error)
}

// DocumentServiceImpl implements the DocumentService interface
type DocumentServiceImpl struct {
	documentLoader   *EnhancedDocumentLoader
	embeddingService *EmbeddingService
	ragRepository    *repository.RagRepository
}

// NewDocumentService creates a new DocumentService instance with enhanced document loading
func NewDocumentService(llmClient client.LLMClient) DocumentService {
	return &DocumentServiceImpl{
		documentLoader:   NewEnhancedDocumentLoader(),
		embeddingService: NewEmbeddingService(llmClient),
		ragRepository:    repository.NewRagRepository(),
	}
}

// CreateRAG implements DocumentService.CreateRAG
func (ds *DocumentServiceImpl) CreateRAG(modelName, ragName, folderPath string, options DocumentLoaderOptions) error {
	// Load documents from folder
	documents, err := ds.documentLoader.LoadDocumentsFromFolderWithOptions(folderPath, options)
	if err != nil {
		return err
	}

	// Create new RAG system
	rag := &domain.RagSystem{
		Name:                   ragName,
		ModelName:              modelName,
		Documents:              documents,
		ChunkingStrategy:       options.ChunkingStrategy,
		VectorStoreType:        options.VectorStore,
		QdrantHost:             options.QdrantHost,
		QdrantPort:             options.QdrantPort,
		QdrantAPIKey:           options.QdrantAPIKey,
		QdrantCollectionName:   options.QdrantCollectionName,
		QdrantGRPC:             options.QdrantGRPC,
		RerankerEnabled:        options.EnableReranker,
		RerankerModel:          options.RerankerModel,
		RerankerWeight:         options.RerankerWeight,
		APIProfileName:         options.APIProfileName,
	}

	// Generate embeddings for all documents
	if err := ds.generateEmbeddings(rag, modelName); err != nil {
		return err
	}

	// Save the RAG system
	return ds.ragRepository.Save(rag)
}

// AddDocuments implements DocumentService.AddDocuments
func (ds *DocumentServiceImpl) AddDocuments(ragName string, folderPath string, options DocumentLoaderOptions) error {
	// Load existing RAG
	rag, err := ds.LoadRAG(ragName)
	if err != nil {
		return err
	}

	// Load new documents
	newDocuments, err := ds.documentLoader.LoadDocumentsFromFolderWithOptions(folderPath, options)
	if err != nil {
		return err
	}

	// Add documents to RAG
	for _, doc := range newDocuments {
		rag.AddDocument(doc)
	}

	// Generate embeddings for new documents
	if err := ds.generateEmbeddings(rag, rag.ModelName); err != nil {
		return err
	}

	// Save updated RAG
	return ds.ragRepository.Save(rag)
}

// GetChunks implements DocumentService.GetChunks
func (ds *DocumentServiceImpl) GetChunks(ragName string, filter ChunkFilter) ([]*domain.DocumentChunk, error) {
	rag, err := ds.LoadRAG(ragName)
	if err != nil {
		return nil, err
	}

	var filteredChunks []*domain.DocumentChunk
	for _, chunk := range rag.Chunks {
		if ds.matchesFilter(chunk, filter) {
			filteredChunks = append(filteredChunks, chunk)
		}
	}

	return filteredChunks, nil
}

// LoadRAG implements DocumentService.LoadRAG
func (ds *DocumentServiceImpl) LoadRAG(ragName string) (*domain.RagSystem, error) {
	return ds.ragRepository.Load(ragName)
}

// UpdateRAG implements DocumentService.UpdateRAG
func (ds *DocumentServiceImpl) UpdateRAG(rag *domain.RagSystem) error {
	return ds.ragRepository.Save(rag)
}

// ListRAGs implements DocumentService.ListRAGs
func (ds *DocumentServiceImpl) ListRAGs() ([]string, error) {
	return ds.ragRepository.ListAll()
}

// generateEmbeddings generates embeddings for all chunks in the RAG system
func (ds *DocumentServiceImpl) generateEmbeddings(rag *domain.RagSystem, modelName string) error {
	// Create chunker service with default options since RAG doesn't store these directly
	chunkerService := NewChunkerService(ChunkingConfig{
		ChunkSize:    1000, // Use sensible defaults
		ChunkOverlap: 200,
	})

	// Generate chunks for all documents
	var allChunks []*domain.DocumentChunk
	for _, doc := range rag.Documents {
		chunks := chunkerService.ChunkDocument(doc)
		for i, chunk := range chunks {
			chunk.ChunkNumber = i
			chunk.TotalChunks = len(chunks)
		}
		allChunks = append(allChunks, chunks...)
	}

	// Generate embeddings for all chunks
	if err := ds.embeddingService.GenerateChunkEmbeddings(allChunks, modelName); err != nil {
		return err
	}

	// Add chunks to RAG
	for _, chunk := range allChunks {
		rag.AddChunk(chunk)
	}

	return nil
}

// matchesFilter checks if a chunk matches the given filter criteria
func (ds *DocumentServiceImpl) matchesFilter(chunk *domain.DocumentChunk, filter ChunkFilter) bool {
	if filter.DocumentSubstring != "" && !strings.Contains(chunk.DocumentID, filter.DocumentSubstring) {
		return false
	}
	return true
}