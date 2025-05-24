package agent

import (
	"context"
)

// RagService represents the interface for interacting with the RAG system
type RagService interface {
	// Query performs a RAG query and returns the result
	Query(ctx context.Context, query string) (string, error)
}
