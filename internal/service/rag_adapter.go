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

	// Use a more reasonable default context size for agent queries
	// This should provide better results for complex agent queries
	contextSize := 10 // Increased from 3 to 10 for better context

	// If reranker is enabled, use a higher initial retrieval count
	if a.ragSystem.RerankerEnabled {
		// Use reranker's TopK if configured, otherwise use a reasonable default
		if a.ragSystem.RerankerTopK > 0 {
			contextSize = a.ragSystem.RerankerTopK
		} else {
			contextSize = 15 // Higher default for reranked results
		}
	}

	// Call the existing RagService with improved parameters
	return a.ragService.Query(a.ragSystem, query, contextSize)
}
