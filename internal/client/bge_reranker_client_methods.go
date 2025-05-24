package client

import (
	"context"
)

// Rerank implements the RerankerClient interface for BGERerankerClient
// Note: BGERerankerClient is used with a Python microservice, not directly
func (c *BGERerankerClient) Rerank(ctx context.Context, query string, documents []string, topK int) ([]RerankResult, error) {
	// BGERerankerClient uses ComputeScores method with pairs
	pairs := make([][]string, len(documents))
	for i, doc := range documents {
		pairs[i] = []string{query, doc}
	}
	
	scores, err := c.ComputeScores(pairs, false)
	if err != nil {
		return nil, err
	}
	
	// Convert scores to rerank results
	results := make([]RerankResult, len(documents))
	for i, doc := range documents {
		results[i] = RerankResult{
			Index:          i,
			Document:       doc,
			Score:          scores[i],
			RelevanceScore: scores[i],
		}
	}
	
	// Sort and return top K
	// Note: Sorting would be done by the service layer
	return results, nil
}

// GetRerankerModel returns the reranker model name for BGERerankerClient
func (c *BGERerankerClient) GetRerankerModel() string {
	return c.GetModelName()
}

// Health implements the RerankerClient interface with context for BGERerankerClient
func (c *BGERerankerClient) Health(ctx context.Context) error {
	// Check if dependencies are available
	return c.CheckDependencies()
}

// GetRerankerModel returns the reranker model name for BGEONNXRerankerClient
func (c *BGEONNXRerankerClient) GetRerankerModel() string {
	return "BAAI/bge-reranker-large"
}

// Health implements the RerankerClient interface for BGEONNXRerankerClient
func (c *BGEONNXRerankerClient) Health(ctx context.Context) error {
	// ONNX client is initialized and ready
	return nil
}

// Rerank implements the RerankerClient interface for BGEONNXRerankerClient
func (c *BGEONNXRerankerClient) Rerank(ctx context.Context, query string, documents []string, topK int) ([]RerankResult, error) {
	// Create query-document pairs
	pairs := make([][]string, len(documents))
	for i, doc := range documents {
		pairs[i] = []string{query, doc}
	}
	
	// Compute scores using the ONNX server
	scores, err := c.ComputeScores(pairs, false)
	if err != nil {
		return nil, err
	}
	
	// Convert scores to rerank results
	results := make([]RerankResult, len(documents))
	for i, doc := range documents {
		results[i] = RerankResult{
			Index:          i,
			Document:       doc,
			Score:          scores[i],
			RelevanceScore: scores[i],
		}
	}
	
	// Sort by score descending
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Score > results[i].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
	
	// Return top K
	if topK > 0 && topK < len(results) {
		return results[:topK], nil
	}
	
	return results, nil
}

// GetRerankerModel returns the reranker model name for PureGoBGEClient
func (c *PureGoBGEClient) GetRerankerModel() string {
	return "BAAI/bge-reranker-large"
}