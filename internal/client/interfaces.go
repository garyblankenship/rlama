package client

import (
	"context"
)

// LLMProvider represents the provider type
type LLMProvider string

const (
	ProviderOllama      LLMProvider = "ollama"
	ProviderOpenAI      LLMProvider = "openai"
	ProviderOpenAIAPI   LLMProvider = "openai-api"  // For OpenAI-compatible APIs
	ProviderLMStudio    LLMProvider = "lmstudio"
	ProviderVLLM        LLMProvider = "vllm"
	ProviderTGI         LLMProvider = "tgi"
)

// UnifiedLLMClient provides a unified interface for all LLM providers
type UnifiedLLMClient interface {
	// Core LLM functionality
	LLMClient
	
	// Additional unified functionality
	EmbeddingClient
	CompletionClient
	
	// Provider information
	GetProvider() LLMProvider
	GetModelName() string
	IsAvailable(ctx context.Context) error
	
	// Configuration
	SetModel(model string) error
	GetSupportedModels(ctx context.Context) ([]string, error)
}

// UnifiedClientConfig holds configuration for creating unified clients
type UnifiedClientConfig struct {
	Provider     LLMProvider
	Model        string
	BaseURL      string
	APIKey       string
	ProfileName  string
	Options      map[string]interface{}
}

// BaseUnifiedClient provides common functionality for unified clients
type BaseUnifiedClient struct {
	provider  LLMProvider
	modelName string
	client    LLMClient
}

// GetProvider returns the provider type
func (b *BaseUnifiedClient) GetProvider() LLMProvider {
	return b.provider
}

// GetModelName returns the current model name
func (b *BaseUnifiedClient) GetModelName() string {
	return b.modelName
}

// SetModel updates the model name
func (b *BaseUnifiedClient) SetModel(model string) error {
	b.modelName = model
	return nil
}

// GenerateCompletion delegates to the underlying client
func (b *BaseUnifiedClient) GenerateCompletion(model, prompt string) (string, error) {
	return b.client.GenerateCompletion(model, prompt)
}

// GenerateEmbedding delegates to the underlying client
func (b *BaseUnifiedClient) GenerateEmbedding(model, text string) ([]float32, error) {
	return b.client.GenerateEmbedding(model, text)
}

// EmbeddingClient provides embedding-specific functionality
type EmbeddingClient interface {
	// Generate embeddings for text
	GenerateEmbedding(model, text string) ([]float32, error)
	
	// Batch embedding generation
	GenerateEmbeddings(model string, texts []string) ([][]float32, error)
	
	// Get embedding dimension for a model
	GetEmbeddingDimension(model string) (int, error)
}

// CompletionClient provides completion-specific functionality
type CompletionClient interface {
	// Generate completion with options
	GenerateCompletionWithOptions(request CompletionRequest) (*CompletionResponse, error)
	
	// Stream completions
	StreamCompletion(request CompletionRequest) (<-chan CompletionChunk, error)
}

// CompletionRequest represents a completion request with all options
type CompletionRequest struct {
	Model            string
	Prompt           string
	SystemPrompt     string
	Messages         []Message
	Temperature      float64
	MaxTokens        int
	TopP             float64
	FrequencyPenalty float64
	PresencePenalty  float64
	Stop             []string
	Stream           bool
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompletionResponse represents a completion response
type CompletionResponse struct {
	ID      string
	Model   string
	Content string
	Usage   TokenUsage
	Choices []CompletionChoice
}

// CompletionChoice represents a single completion choice
type CompletionChoice struct {
	Index        int
	Message      Message
	FinishReason string
}

// TokenUsage represents token usage information
type TokenUsage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// CompletionChunk represents a streaming completion chunk
type CompletionChunk struct {
	Content string
	Error   error
	Done    bool
}

// RerankerClient provides reranking functionality
type RerankerClient interface {
	// Rerank documents based on query relevance
	Rerank(ctx context.Context, query string, documents []string, topK int) ([]RerankResult, error)
	
	// Get reranker model information
	GetRerankerModel() string
	
	// Health check
	Health(ctx context.Context) error
}

// ClientFactory creates clients based on configuration
type ClientFactory interface {
	// Create a unified LLM client
	CreateClient(config UnifiedClientConfig) (UnifiedLLMClient, error)
	
	// Create specialized clients
	CreateEmbeddingClient(config UnifiedClientConfig) (EmbeddingClient, error)
	CreateCompletionClient(config UnifiedClientConfig) (CompletionClient, error)
	CreateRerankerClient(config UnifiedClientConfig) (RerankerClient, error)
	
	// List available providers
	GetAvailableProviders() []LLMProvider
}