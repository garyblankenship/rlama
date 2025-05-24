package client

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// DefaultClientFactory is the default implementation of ClientFactory
type DefaultClientFactory struct {
	// Configuration options
	defaultProvider LLMProvider
	providers       map[LLMProvider]bool
}

// NewDefaultClientFactory creates a new default client factory
func NewDefaultClientFactory() *DefaultClientFactory {
	return &DefaultClientFactory{
		defaultProvider: ProviderOllama,
		providers: map[LLMProvider]bool{
			ProviderOllama:    true,
			ProviderOpenAI:    true,
			ProviderOpenAIAPI: true,
			ProviderLMStudio:  true,
			ProviderVLLM:      true,
			ProviderTGI:       true,
		},
	}
}

// CreateClient creates a unified LLM client based on configuration
func (f *DefaultClientFactory) CreateClient(config UnifiedClientConfig) (UnifiedLLMClient, error) {
	// Auto-detect provider if not specified
	if config.Provider == "" {
		config.Provider = f.detectProvider(config)
	}
	
	// Validate provider
	if !f.providers[config.Provider] {
		return nil, fmt.Errorf("unsupported provider: %s", config.Provider)
	}
	
	// Create client based on provider
	switch config.Provider {
	case ProviderOllama:
		return NewUnifiedOllamaClient(config)
		
	case ProviderOpenAI, ProviderOpenAIAPI, ProviderLMStudio, ProviderVLLM, ProviderTGI:
		return NewUnifiedOpenAIClient(config)
		
	default:
		return nil, fmt.Errorf("provider not implemented: %s", config.Provider)
	}
}

// CreateEmbeddingClient creates a specialized embedding client
func (f *DefaultClientFactory) CreateEmbeddingClient(config UnifiedClientConfig) (EmbeddingClient, error) {
	client, err := f.CreateClient(config)
	if err != nil {
		return nil, err
	}
	
	// Check if the client supports embedding
	if embeddingClient, ok := client.(EmbeddingClient); ok {
		return embeddingClient, nil
	}
	
	return nil, fmt.Errorf("provider %s does not support embeddings", config.Provider)
}

// CreateCompletionClient creates a specialized completion client
func (f *DefaultClientFactory) CreateCompletionClient(config UnifiedClientConfig) (CompletionClient, error) {
	client, err := f.CreateClient(config)
	if err != nil {
		return nil, err
	}
	
	// Check if the client supports advanced completions
	if completionClient, ok := client.(CompletionClient); ok {
		return completionClient, nil
	}
	
	return nil, fmt.Errorf("provider %s does not support advanced completions", config.Provider)
}

// CreateRerankerClient creates a specialized reranker client
func (f *DefaultClientFactory) CreateRerankerClient(config UnifiedClientConfig) (RerankerClient, error) {
	// Reranker clients are specialized - only BGE reranker for now
	switch config.Provider {
	case "bge", "bge-reranker":
		if config.Options["use_onnx"] == true {
			modelDir := ""
			if dir, ok := config.Options["model_dir"].(string); ok {
				modelDir = dir
			}
			return NewBGEONNXRerankerClient(modelDir)
		}
		// Default to HTTP client
		baseURL := config.BaseURL
		if baseURL == "" {
			baseURL = "http://localhost:8001"
		}
		return NewBGERerankerClient(baseURL), nil
		
	case "pure-go-bge":
		modelPath := ""
		if path, ok := config.Options["model_path"].(string); ok {
			modelPath = path
		}
		usePureGo := true
		if val, ok := config.Options["use_pure_go"].(bool); ok {
			usePureGo = val
		}
		fallbackURL := config.BaseURL
		if fallbackURL == "" {
			fallbackURL = "http://localhost:8001"
		}
		return NewPureGoBGEClient(modelPath, usePureGo, fallbackURL)
		
	default:
		return nil, fmt.Errorf("unsupported reranker provider: %s", config.Provider)
	}
}

// GetAvailableProviders returns the list of available providers
func (f *DefaultClientFactory) GetAvailableProviders() []LLMProvider {
	providers := make([]LLMProvider, 0, len(f.providers))
	for provider, enabled := range f.providers {
		if enabled {
			providers = append(providers, provider)
		}
	}
	return providers
}

// detectProvider attempts to detect the provider from configuration
func (f *DefaultClientFactory) detectProvider(config UnifiedClientConfig) LLMProvider {
	// Check model name patterns
	modelLower := strings.ToLower(config.Model)
	if strings.HasPrefix(modelLower, "gpt-") || 
	   strings.Contains(modelLower, "text-embedding-") {
		return ProviderOpenAI
	}
	
	// Check base URL patterns
	if config.BaseURL != "" {
		return DetectProviderFromURL(config.BaseURL)
	}
	
	// Default to Ollama
	return f.defaultProvider
}

// CreateClientFromModel creates a client based on model name
func CreateClientFromModel(modelName string) (UnifiedLLMClient, error) {
	factory := NewDefaultClientFactory()
	
	config := UnifiedClientConfig{
		Model: modelName,
	}
	
	// Auto-detect provider from model name
	return factory.CreateClient(config)
}

// CreateClientFromProfile creates a client based on a profile name
func CreateClientFromProfile(profileName string) (UnifiedLLMClient, error) {
	factory := NewDefaultClientFactory()
	
	config := UnifiedClientConfig{
		ProfileName: profileName,
	}
	
	return factory.CreateClient(config)
}

// TestClientConnectivity tests if a client can connect to its service
func TestClientConnectivity(client UnifiedLLMClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	return client.IsAvailable(ctx)
}