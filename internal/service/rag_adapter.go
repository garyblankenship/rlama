package service

import (
	"context"
	"fmt"

	"github.com/dontizi/rlama/internal/domain"
)

// RagServiceAdapter adapts the existing RagService to the agent's RagService interface
type RagServiceAdapter struct {
	ragService RagService
	ragSystem  *domain.RagSystem
}

// NewRagServiceAdapter creates a new adapter for RagService
func NewRagServiceAdapter(ragService RagService, ragSystem *domain.RagSystem) *RagServiceAdapter {
	return &RagServiceAdapter{
		ragService: ragService,
		ragSystem:  ragSystem,
	}
}

// Query implements the agent.RagService interface
func (a *RagServiceAdapter) Query(ctx context.Context, query string) (string, error) {
	// Check if RAG system is available
	if a.ragSystem == nil {
		return "", fmt.Errorf("no RAG system available - please specify a RAG system when running the agent")
	}

	// Call the existing RagService with default parameters
	return a.ragService.Query(a.ragSystem, query, 3) // Using 3 as a default k value
}
