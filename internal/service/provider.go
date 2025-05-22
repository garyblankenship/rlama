package service

import (
	"fmt"
	"sync"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/repository"
)

// ServiceProvider manages the creation and lifecycle of all services
// It implements dependency injection pattern for centralized service management
type ServiceProvider struct {
	config *ServiceConfig
	
	// Cached clients (lazy initialization)
	ollamaClient     *client.OllamaClient
	llmClient        client.LLMClient
	clientMutex      sync.RWMutex
	
	// Cached services (lazy initialization)  
	documentService  DocumentService
	queryService     QueryService
	watchService     WatchService
	ragService       RagService
	serviceMutex     sync.RWMutex
	
	// Repositories
	ragRepository    *repository.RagRepository
	
	// Service factories (for testing/mocking)
	documentServiceFactory func(client.LLMClient) DocumentService
	queryServiceFactory    func(client.LLMClient, *client.OllamaClient, DocumentService) QueryService
	watchServiceFactory    func(DocumentService, RagService) WatchService
}

// NewServiceProvider creates a new service provider with the given configuration
func NewServiceProvider(config *ServiceConfig) (*ServiceProvider, error) {
	if config == nil {
		config = NewServiceConfig()
	}
	
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid service configuration: %w", err)
	}
	
	return &ServiceProvider{
		config:        config,
		ragRepository: repository.NewRagRepository(),
		
		// Default factories
		documentServiceFactory: NewDocumentService,
		queryServiceFactory:    NewQueryService,
		watchServiceFactory:    NewWatchService,
	}, nil
}

// GetConfig returns a copy of the current configuration
func (sp *ServiceProvider) GetConfig() *ServiceConfig {
	return sp.config.Clone()
}

// UpdateConfig updates the service provider configuration and clears cached services
func (sp *ServiceProvider) UpdateConfig(config *ServiceConfig) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid service configuration: %w", err)
	}
	
	sp.clientMutex.Lock()
	sp.serviceMutex.Lock()
	defer sp.clientMutex.Unlock()
	defer sp.serviceMutex.Unlock()
	
	sp.config = config.Clone()
	
	// Clear cached clients and services to force recreation with new config
	sp.ollamaClient = nil
	sp.llmClient = nil
	sp.documentService = nil
	sp.queryService = nil
	sp.watchService = nil
	sp.ragService = nil
	
	return nil
}

// GetOllamaClient returns the Ollama client (cached after first creation)
func (sp *ServiceProvider) GetOllamaClient() *client.OllamaClient {
	sp.clientMutex.RLock()
	if sp.ollamaClient != nil {
		defer sp.clientMutex.RUnlock()
		return sp.ollamaClient
	}
	sp.clientMutex.RUnlock()
	
	sp.clientMutex.Lock()
	defer sp.clientMutex.Unlock()
	
	// Double-check after acquiring write lock
	if sp.ollamaClient != nil {
		return sp.ollamaClient
	}
	
	sp.ollamaClient = client.NewOllamaClient(sp.config.OllamaHost, sp.config.OllamaPort)
	return sp.ollamaClient
}

// GetLLMClient returns the appropriate LLM client based on configuration and model
func (sp *ServiceProvider) GetLLMClient(modelName string) (client.LLMClient, error) {
	// For profile-based clients, create fresh instances
	if sp.config.APIProfileName != "" {
		ollamaClient := sp.GetOllamaClient()
		return client.GetLLMClientWithProfile(modelName, sp.config.APIProfileName, ollamaClient)
	}
	
	// For direct OpenAI models
	if client.IsOpenAIModel(modelName) {
		return client.NewOpenAIClient(), nil
	}
	
	// Default to Ollama client
	return sp.GetOllamaClient(), nil
}

// GetDocumentService returns the document service (cached after first creation)
func (sp *ServiceProvider) GetDocumentService() DocumentService {
	sp.serviceMutex.RLock()
	if sp.documentService != nil {
		defer sp.serviceMutex.RUnlock()
		return sp.documentService
	}
	sp.serviceMutex.RUnlock()
	
	sp.serviceMutex.Lock()
	defer sp.serviceMutex.Unlock()
	
	// Double-check after acquiring write lock
	if sp.documentService != nil {
		return sp.documentService
	}
	
	// Create LLM client for embeddings
	llmClient, err := sp.GetLLMClient("")
	if err != nil {
		// Fallback to Ollama client
		llmClient = sp.GetOllamaClient()
	}
	
	sp.documentService = sp.documentServiceFactory(llmClient)
	return sp.documentService
}

// GetQueryService returns the query service (cached after first creation)
func (sp *ServiceProvider) GetQueryService() QueryService {
	sp.serviceMutex.RLock()
	if sp.queryService != nil {
		defer sp.serviceMutex.RUnlock()
		return sp.queryService
	}
	sp.serviceMutex.RUnlock()
	
	sp.serviceMutex.Lock()
	defer sp.serviceMutex.Unlock()
	
	// Double-check after acquiring write lock
	if sp.queryService != nil {
		return sp.queryService
	}
	
	// Create dependencies
	llmClient, err := sp.GetLLMClient("")
	if err != nil {
		// Fallback to Ollama client
		llmClient = sp.GetOllamaClient()
	}
	
	ollamaClient := sp.GetOllamaClient()
	documentService := sp.GetDocumentService()
	
	sp.queryService = sp.queryServiceFactory(llmClient, ollamaClient, documentService)
	return sp.queryService
}

// GetWatchService returns the watch service (cached after first creation)
func (sp *ServiceProvider) GetWatchService() WatchService {
	sp.serviceMutex.RLock()
	if sp.watchService != nil {
		defer sp.serviceMutex.RUnlock()
		return sp.watchService
	}
	sp.serviceMutex.RUnlock()
	
	sp.serviceMutex.Lock()
	defer sp.serviceMutex.Unlock()
	
	// Double-check after acquiring write lock
	if sp.watchService != nil {
		return sp.watchService
	}
	
	// Create dependencies
	documentService := sp.GetDocumentService()
	ragService := sp.GetRagService()
	
	sp.watchService = sp.watchServiceFactory(documentService, ragService)
	return sp.watchService
}

// GetEmbeddingService returns the embedding service
func (sp *ServiceProvider) GetEmbeddingService() *EmbeddingService {
	llmClient, err := sp.GetLLMClient("")
	if err != nil {
		// Fallback to Ollama client
		llmClient = sp.GetOllamaClient()
	}
	return NewEmbeddingService(llmClient)
}

// GetRagService returns the composite RAG service (cached after first creation)
func (sp *ServiceProvider) GetRagService() RagService {
	sp.serviceMutex.RLock()
	if sp.ragService != nil {
		defer sp.serviceMutex.RUnlock()
		return sp.ragService
	}
	sp.serviceMutex.RUnlock()
	
	sp.serviceMutex.Lock()
	defer sp.serviceMutex.Unlock()
	
	// Double-check after acquiring write lock
	if sp.ragService != nil {
		return sp.ragService
	}
	
	// Create dependencies
	llmClient, err := sp.GetLLMClient("")
	if err != nil {
		// Fallback to Ollama client
		llmClient = sp.GetOllamaClient()
	}
	
	ollamaClient := sp.GetOllamaClient()
	
	sp.ragService = NewCompositeRagService(llmClient, ollamaClient)
	return sp.ragService
}

// CreateRagServiceForModel creates a RAG service configured for a specific model
func (sp *ServiceProvider) CreateRagServiceForModel(modelName string) (RagService, error) {
	llmClient, err := sp.GetLLMClient(modelName)
	if err != nil {
		return nil, fmt.Errorf("failed to create LLM client for model %s: %w", modelName, err)
	}
	
	ollamaClient := sp.GetOllamaClient()
	
	// Use configuration-aware service if ONNX reranker is enabled
	if sp.config.UseONNXReranker {
		return NewCompositeRagServiceWithConfig(llmClient, ollamaClient, sp.config), nil
	}
	
	return NewCompositeRagService(llmClient, ollamaClient), nil
}

// SetDocumentServiceFactory allows injecting a custom document service factory (for testing)
func (sp *ServiceProvider) SetDocumentServiceFactory(factory func(client.LLMClient) DocumentService) {
	sp.serviceMutex.Lock()
	defer sp.serviceMutex.Unlock()
	
	sp.documentServiceFactory = factory
	sp.documentService = nil // Clear cached service
}

// SetQueryServiceFactory allows injecting a custom query service factory (for testing)
func (sp *ServiceProvider) SetQueryServiceFactory(factory func(client.LLMClient, *client.OllamaClient, DocumentService) QueryService) {
	sp.serviceMutex.Lock()
	defer sp.serviceMutex.Unlock()
	
	sp.queryServiceFactory = factory
	sp.queryService = nil // Clear cached service
}

// SetWatchServiceFactory allows injecting a custom watch service factory (for testing)
func (sp *ServiceProvider) SetWatchServiceFactory(factory func(DocumentService, RagService) WatchService) {
	sp.serviceMutex.Lock()
	defer sp.serviceMutex.Unlock()
	
	sp.watchServiceFactory = factory
	sp.watchService = nil // Clear cached service
}

// Reset clears all cached services and clients (useful for testing)
func (sp *ServiceProvider) Reset() {
	sp.clientMutex.Lock()
	sp.serviceMutex.Lock()
	defer sp.clientMutex.Unlock()
	defer sp.serviceMutex.Unlock()
	
	sp.ollamaClient = nil
	sp.llmClient = nil
	sp.documentService = nil
	sp.queryService = nil
	sp.watchService = nil
	sp.ragService = nil
}