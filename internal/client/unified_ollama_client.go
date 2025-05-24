package client

import (
	"context"
	"fmt"
)

// UnifiedOllamaClient provides a unified interface for Ollama
type UnifiedOllamaClient struct {
	BaseUnifiedClient
	ollamaClient *OllamaClient
}

// Ensure UnifiedOllamaClient implements UnifiedLLMClient
var _ UnifiedLLMClient = (*UnifiedOllamaClient)(nil)

// CheckLLMAndModel implements the LLMClient interface
func (u *UnifiedOllamaClient) CheckLLMAndModel(model string) error {
	return u.ollamaClient.CheckOllamaAndModel(model)
}

// NewUnifiedOllamaClient creates a new unified Ollama client
func NewUnifiedOllamaClient(config UnifiedClientConfig) (*UnifiedOllamaClient, error) {
	// Create base Ollama client
	ollamaClient := NewDefaultOllamaClient()
	
	// Configure with custom URL if provided
	// Note: This would need to be implemented in OllamaClient
	
	client := &UnifiedOllamaClient{
		BaseUnifiedClient: BaseUnifiedClient{
			provider:  ProviderOllama,
			modelName: config.Model,
			client:    ollamaClient,
		},
		ollamaClient: ollamaClient,
	}
	
	return client, nil
}

// IsAvailable checks if Ollama is available
func (u *UnifiedOllamaClient) IsAvailable(ctx context.Context) error {
	// Check if Ollama is running with a test model
	return u.ollamaClient.CheckOllamaAndModel("llama3.2")
}

// GetSupportedModels returns available Ollama models
func (u *UnifiedOllamaClient) GetSupportedModels(ctx context.Context) ([]string, error) {
	// For now, return a list of common models
	// TODO: Implement actual model listing when available in OllamaClient
	return []string{
		"llama3.2",
		"llama3.1",
		"llama2",
		"mistral",
		"mixtral",
		"deepseek-coder",
		"nomic-embed-text",
		"snowflake-arctic-embed2",
	}, nil
}

// GenerateEmbeddings generates embeddings for multiple texts
func (u *UnifiedOllamaClient) GenerateEmbeddings(model string, texts []string) ([][]float32, error) {
	embeddings := make([][]float32, 0, len(texts))
	
	for _, text := range texts {
		embedding, err := u.ollamaClient.GenerateEmbedding(model, text)
		if err != nil {
			return nil, fmt.Errorf("failed to generate embedding: %w", err)
		}
		embeddings = append(embeddings, embedding)
	}
	
	return embeddings, nil
}

// GetEmbeddingDimension returns the embedding dimension for a model
func (u *UnifiedOllamaClient) GetEmbeddingDimension(model string) (int, error) {
	// Generate a test embedding to determine dimension
	testEmbedding, err := u.ollamaClient.GenerateEmbedding(model, "test")
	if err != nil {
		return 0, fmt.Errorf("failed to determine embedding dimension: %w", err)
	}
	
	return len(testEmbedding), nil
}

// GenerateCompletionWithOptions generates a completion with advanced options
func (u *UnifiedOllamaClient) GenerateCompletionWithOptions(request CompletionRequest) (*CompletionResponse, error) {
	// For now, use the simple completion method
	// TODO: Implement advanced options support for Ollama
	completion, err := u.ollamaClient.GenerateCompletion(request.Model, request.Prompt)
	if err != nil {
		return nil, err
	}
	
	return &CompletionResponse{
		Model:   request.Model,
		Content: completion,
		Choices: []CompletionChoice{
			{
				Index: 0,
				Message: Message{
					Role:    "assistant",
					Content: completion,
				},
				FinishReason: "stop",
			},
		},
	}, nil
}

// StreamCompletion streams completions (not yet implemented)
func (u *UnifiedOllamaClient) StreamCompletion(request CompletionRequest) (<-chan CompletionChunk, error) {
	// TODO: Implement streaming for Ollama
	return nil, fmt.Errorf("streaming not yet implemented for Ollama")
}