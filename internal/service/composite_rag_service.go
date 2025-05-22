package service

import (
	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/domain"
)

// CompositeRagService implements RagService by composing focused services
// This replaces the monolithic RagServiceImpl with a cleaner architecture
type CompositeRagService struct {
	documentService DocumentService
	queryService    QueryService
	watchService    WatchService
	ollamaClient    *client.OllamaClient
}

// NewCompositeRagService creates a new composite RAG service
func NewCompositeRagService(llmClient client.LLMClient, ollamaClient *client.OllamaClient) RagService {
	// Create focused services
	documentService := NewDocumentService(llmClient)
	queryService := NewQueryService(llmClient, ollamaClient, documentService)
	
	// Create the composite service first, then create watch service with it
	compositeService := &CompositeRagService{
		documentService: documentService,
		queryService:    queryService,
		ollamaClient:    ollamaClient,
	}
	
	// Now create watch service with the composite service
	compositeService.watchService = NewWatchService(documentService, compositeService)

	return compositeService
}

// NewCompositeRagServiceWithConfig creates a new composite RAG service with configuration options
func NewCompositeRagServiceWithConfig(llmClient client.LLMClient, ollamaClient *client.OllamaClient, config *ServiceConfig) RagService {
	// Create focused services with configuration
	documentService := NewDocumentService(llmClient)
	queryService := NewQueryServiceWithConfig(llmClient, ollamaClient, documentService, config)
	
	// Create the composite service first, then create watch service with it
	compositeService := &CompositeRagService{
		documentService: documentService,
		queryService:    queryService,
		ollamaClient:    ollamaClient,
	}
	
	// Now create watch service with the composite service
	compositeService.watchService = NewWatchService(documentService, compositeService)

	return compositeService
}

// CreateRagWithOptions implements RagService.CreateRagWithOptions
func (crs *CompositeRagService) CreateRagWithOptions(modelName, ragName, folderPath string, options DocumentLoaderOptions) error {
	return crs.documentService.CreateRAG(modelName, ragName, folderPath, options)
}

// GetRagChunks implements RagService.GetRagChunks
func (crs *CompositeRagService) GetRagChunks(ragName string, filter ChunkFilter) ([]*domain.DocumentChunk, error) {
	return crs.documentService.GetChunks(ragName, filter)
}

// LoadRag implements RagService.LoadRag
func (crs *CompositeRagService) LoadRag(ragName string) (*domain.RagSystem, error) {
	return crs.documentService.LoadRAG(ragName)
}

// Query implements RagService.Query
func (crs *CompositeRagService) Query(rag *domain.RagSystem, query string, contextSize int) (string, error) {
	return crs.queryService.Query(rag, query, contextSize)
}

// AddDocsWithOptions implements RagService.AddDocsWithOptions
func (crs *CompositeRagService) AddDocsWithOptions(ragName string, folderPath string, options DocumentLoaderOptions) error {
	return crs.documentService.AddDocuments(ragName, folderPath, options)
}

// UpdateModel implements RagService.UpdateModel
func (crs *CompositeRagService) UpdateModel(ragName string, newModel string) error {
	// Load the RAG
	rag, err := crs.documentService.LoadRAG(ragName)
	if err != nil {
		return err
	}

	// Update the model
	rag.ModelName = newModel

	// Save the updated RAG
	return crs.documentService.UpdateRAG(rag)
}

// UpdateRag implements RagService.UpdateRag
func (crs *CompositeRagService) UpdateRag(rag *domain.RagSystem) error {
	return crs.documentService.UpdateRAG(rag)
}

// UpdateRerankerModel implements RagService.UpdateRerankerModel
func (crs *CompositeRagService) UpdateRerankerModel(ragName string, model string) error {
	return crs.queryService.UpdateRerankerModel(ragName, model)
}

// ListAllRags implements RagService.ListAllRags
func (crs *CompositeRagService) ListAllRags() ([]string, error) {
	return crs.documentService.ListRAGs()
}

// GetOllamaClient implements RagService.GetOllamaClient
func (crs *CompositeRagService) GetOllamaClient() *client.OllamaClient {
	return crs.ollamaClient
}

// SetPreferredEmbeddingModel implements RagService.SetPreferredEmbeddingModel
func (crs *CompositeRagService) SetPreferredEmbeddingModel(model string) {
	crs.queryService.SetPreferredEmbeddingModel(model)
}

// Directory watching methods - delegate to WatchService
func (crs *CompositeRagService) SetupDirectoryWatching(ragName string, dirPath string, watchInterval int, options DocumentLoaderOptions) error {
	return crs.watchService.SetupDirectoryWatching(ragName, dirPath, watchInterval, options)
}

func (crs *CompositeRagService) DisableDirectoryWatching(ragName string) error {
	return crs.watchService.DisableDirectoryWatching(ragName)
}

func (crs *CompositeRagService) CheckWatchedDirectory(ragName string) (int, error) {
	return crs.watchService.CheckWatchedDirectory(ragName)
}

// Web watching methods - delegate to WatchService
func (crs *CompositeRagService) SetupWebWatching(ragName string, websiteURL string, watchInterval int, options domain.WebWatchOptions) error {
	return crs.watchService.SetupWebWatching(ragName, websiteURL, watchInterval, options)
}

func (crs *CompositeRagService) DisableWebWatching(ragName string) error {
	return crs.watchService.DisableWebWatching(ragName)
}

func (crs *CompositeRagService) CheckWatchedWebsite(ragName string) (int, error) {
	return crs.watchService.CheckWatchedWebsite(ragName)
}