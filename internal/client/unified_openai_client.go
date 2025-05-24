package client

import (
	"context"
	"fmt"
	"strings"
)

// UnifiedOpenAIClient provides a unified interface for OpenAI and compatible APIs
type UnifiedOpenAIClient struct {
	BaseUnifiedClient
	openaiClient *OpenAIClient
}

// Ensure UnifiedOpenAIClient implements UnifiedLLMClient
var _ UnifiedLLMClient = (*UnifiedOpenAIClient)(nil)

// CheckLLMAndModel implements the LLMClient interface
func (u *UnifiedOpenAIClient) CheckLLMAndModel(model string) error {
	return u.openaiClient.CheckOpenAIAndModel(model)
}

// NewUnifiedOpenAIClient creates a new unified OpenAI client
func NewUnifiedOpenAIClient(config UnifiedClientConfig) (*UnifiedOpenAIClient, error) {
	var openaiClient *OpenAIClient
	var err error
	
	// Determine provider type
	provider := config.Provider
	if provider == "" {
		provider = ProviderOpenAI
	}
	
	// Create OpenAI client based on configuration
	if config.ProfileName != "" {
		openaiClient, err = NewOpenAIClientWithProfile(config.ProfileName)
		if err != nil {
			return nil, fmt.Errorf("failed to create OpenAI client with profile: %w", err)
		}
	} else if config.BaseURL != "" && config.BaseURL != "https://api.openai.com/v1" {
		// Custom OpenAI-compatible endpoint
		openaiClient = NewGenericOpenAIClient(config.BaseURL, config.APIKey)
		provider = ProviderOpenAIAPI
	} else {
		// Standard OpenAI client
		openaiClient = NewOpenAIClient()
	}
	
	client := &UnifiedOpenAIClient{
		BaseUnifiedClient: BaseUnifiedClient{
			provider:  provider,
			modelName: config.Model,
			client:    openaiClient,
		},
		openaiClient: openaiClient,
	}
	
	return client, nil
}

// IsAvailable checks if OpenAI API is available
func (u *UnifiedOpenAIClient) IsAvailable(ctx context.Context) error {
	// For OpenAI-compatible APIs, we might not need to check specific models
	if u.provider == ProviderOpenAIAPI {
		// Simple connectivity check could be implemented here
		return nil
	}
	
	// For official OpenAI, check with a known model
	testModel := "gpt-3.5-turbo"
	if u.modelName != "" {
		testModel = u.modelName
	}
	
	return u.openaiClient.CheckOpenAIAndModel(testModel)
}

// GetSupportedModels returns available OpenAI models
func (u *UnifiedOpenAIClient) GetSupportedModels(ctx context.Context) ([]string, error) {
	// Return common OpenAI models
	// In a real implementation, this could query the models endpoint
	models := []string{
		"gpt-4-turbo-preview",
		"gpt-4-turbo",
		"gpt-4",
		"gpt-3.5-turbo",
		"gpt-3.5-turbo-16k",
		"text-embedding-ada-002",
		"text-embedding-3-small",
		"text-embedding-3-large",
	}
	
	return models, nil
}

// GenerateEmbeddings generates embeddings for multiple texts
func (u *UnifiedOpenAIClient) GenerateEmbeddings(model string, texts []string) ([][]float32, error) {
	embeddings := make([][]float32, 0, len(texts))
	
	// OpenAI supports batch embedding in a single request
	// For now, we'll do individual requests
	for _, text := range texts {
		embedding, err := u.openaiClient.GenerateEmbedding(model, text)
		if err != nil {
			return nil, fmt.Errorf("failed to generate embedding: %w", err)
		}
		embeddings = append(embeddings, embedding)
	}
	
	return embeddings, nil
}

// GetEmbeddingDimension returns the embedding dimension for a model
func (u *UnifiedOpenAIClient) GetEmbeddingDimension(model string) (int, error) {
	// Known OpenAI embedding dimensions
	dimensions := map[string]int{
		"text-embedding-ada-002":   1536,
		"text-embedding-3-small":   1536,
		"text-embedding-3-large":   3072,
	}
	
	if dim, exists := dimensions[model]; exists {
		return dim, nil
	}
	
	// For unknown models, generate a test embedding
	testEmbedding, err := u.openaiClient.GenerateEmbedding(model, "test")
	if err != nil {
		return 0, fmt.Errorf("failed to determine embedding dimension: %w", err)
	}
	
	return len(testEmbedding), nil
}

// GenerateCompletionWithOptions generates a completion with advanced options
func (u *UnifiedOpenAIClient) GenerateCompletionWithOptions(request CompletionRequest) (*CompletionResponse, error) {
	// Convert to OpenAI format
	openaiReq := OpenAICompletionRequest{
		Model:       request.Model,
		Temperature: request.Temperature,
		MaxTokens:   request.MaxTokens,
	}
	
	// Build messages
	if request.SystemPrompt != "" {
		openaiReq.Messages = append(openaiReq.Messages, OpenAIMessage{
			Role:    "system",
			Content: request.SystemPrompt,
		})
	}
	
	if len(request.Messages) > 0 {
		for _, msg := range request.Messages {
			openaiReq.Messages = append(openaiReq.Messages, OpenAIMessage{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}
	} else if request.Prompt != "" {
		openaiReq.Messages = append(openaiReq.Messages, OpenAIMessage{
			Role:    "user",
			Content: request.Prompt,
		})
	}
	
	// For now, use the simple completion method
	completion, err := u.openaiClient.GenerateCompletion(request.Model, request.Prompt)
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

// StreamCompletion streams completions
func (u *UnifiedOpenAIClient) StreamCompletion(request CompletionRequest) (<-chan CompletionChunk, error) {
	// TODO: Implement streaming for OpenAI
	return nil, fmt.Errorf("streaming not yet implemented for OpenAI")
}

// DetectProviderFromURL detects the provider type from a base URL
func DetectProviderFromURL(baseURL string) LLMProvider {
	url := strings.ToLower(baseURL)
	
	switch {
	case strings.Contains(url, "api.openai.com"):
		return ProviderOpenAI
	case strings.Contains(url, "localhost:1234"):
		return ProviderLMStudio
	case strings.Contains(url, "vllm"):
		return ProviderVLLM
	case strings.Contains(url, "tgi"):
		return ProviderTGI
	default:
		return ProviderOpenAIAPI
	}
}