package agent

import (
	"context"
)

// LLMClient represents a client for interacting with a language model
type LLMClient interface {
	// GenerateCompletion generates a completion for the given prompt
	GenerateCompletion(ctx context.Context, prompt string) (string, error)
	// GenerateEmbedding generates an embedding for the given text
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
	// GenerateStructuredCompletion generates a structured JSON completion
	GenerateStructuredCompletion(ctx context.Context, prompt string, schema map[string]interface{}) (string, error)
}
