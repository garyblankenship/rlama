package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"time"
)

// PureGoBGEClient implements BGE reranking with pure Go tokenization
// and optional fallback to Python inference
type PureGoBGEClient struct {
	tokenizer      *PureGoTokenizer
	pythonEndpoint string
	httpClient     *http.Client
	maxLength      int
	
	// Configuration
	usePureGo    bool
	fallbackURL  string
}

// NewPureGoBGEClient creates a new pure Go BGE reranker client
func NewPureGoBGEClient(modelPath string, usePureGo bool, fallbackURL string) (*PureGoBGEClient, error) {
	tokenizerPath := filepath.Join(modelPath, "tokenizer.json")
	
	tokenizer, err := NewPureGoTokenizer(tokenizerPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create tokenizer: %w", err)
	}
	
	client := &PureGoBGEClient{
		tokenizer:      tokenizer,
		pythonEndpoint: fallbackURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		maxLength:   512,
		usePureGo:   usePureGo,
		fallbackURL: fallbackURL,
	}
	
	return client, nil
}

// RerankRequest represents a rerank request
type RerankRequest struct {
	Query     string   `json:"query"`
	Documents []string `json:"documents"`
	TopK      int      `json:"top_k,omitempty"`
}

// RerankResponse represents a rerank response
type RerankResponse struct {
	Results []RerankResult `json:"results"`
}

// RerankResult represents a single rerank result
type RerankResult struct {
	Index       int     `json:"index"`
	Document    string  `json:"document"`
	Score       float64 `json:"score"`
	RelevanceScore float64 `json:"relevance_score"`
}

// BGEInferenceRequest represents the tokenized input for BGE inference
type BGEInferenceRequest struct {
	InputIDs      [][]int64 `json:"input_ids"`
	AttentionMask [][]int64 `json:"attention_mask"`
	TokenTypeIDs  [][]int64 `json:"token_type_ids"`
}

// BGEInferenceResponse represents the inference response
type BGEInferenceResponse struct {
	Scores []float64 `json:"scores"`
	Error  string    `json:"error,omitempty"`
}

// Rerank reranks documents using BGE with pure Go tokenization
func (c *PureGoBGEClient) Rerank(ctx context.Context, query string, documents []string, topK int) ([]RerankResult, error) {
	if c.usePureGo {
		return c.rerankPureGo(ctx, query, documents, topK)
	}
	return c.rerankWithFallback(ctx, query, documents, topK)
}

// rerankPureGo implements pure Go reranking (tokenization + eventual ONNX inference)
func (c *PureGoBGEClient) rerankPureGo(ctx context.Context, query string, documents []string, topK int) ([]RerankResult, error) {
	// Step 1: Tokenize all query-document pairs using pure Go
	var inputIDs [][]int64
	var attentionMasks [][]int64
	var tokenTypeIDs [][]int64
	
	for _, doc := range documents {
		tokenIDs, attMask, tokenTypes := c.tokenizer.EncodeQueryPassagePair(query, doc, c.maxLength)
		inputIDs = append(inputIDs, tokenIDs)
		attentionMasks = append(attentionMasks, attMask)
		tokenTypeIDs = append(tokenTypeIDs, tokenTypes)
	}
	
	// Step 2: For now, fall back to Python inference service
	// TODO: Replace with pure Go ONNX inference when available
	request := BGEInferenceRequest{
		InputIDs:      inputIDs,
		AttentionMask: attentionMasks,
		TokenTypeIDs:  tokenTypeIDs,
	}
	
	scores, err := c.callInferenceService(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("inference failed: %w", err)
	}
	
	// Step 3: Create results with scores
	results := make([]RerankResult, len(documents))
	for i, doc := range documents {
		score := 0.0
		if i < len(scores) {
			score = scores[i]
		}
		
		results[i] = RerankResult{
			Index:          i,
			Document:       doc,
			Score:          score,
			RelevanceScore: score,
		}
	}
	
	// Step 4: Sort by score and return top K
	return c.selectTopK(results, topK), nil
}

// rerankWithFallback uses the original Python microservice approach
func (c *PureGoBGEClient) rerankWithFallback(ctx context.Context, query string, documents []string, topK int) ([]RerankResult, error) {
	request := RerankRequest{
		Query:     query,
		Documents: documents,
		TopK:      topK,
	}
	
	return c.callRerankService(ctx, request)
}

// callInferenceService calls the Python inference service with tokenized data
func (c *PureGoBGEClient) callInferenceService(ctx context.Context, request BGEInferenceRequest) ([]float64, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	// Call inference endpoint (hypothetical)
	inferenceURL := c.fallbackURL + "/inference"
	req, err := http.NewRequestWithContext(ctx, "POST", inferenceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call inference service: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("inference service returned status %d", resp.StatusCode)
	}
	
	var response BGEInferenceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	if response.Error != "" {
		return nil, fmt.Errorf("inference error: %s", response.Error)
	}
	
	return response.Scores, nil
}

// callRerankService calls the original rerank service
func (c *PureGoBGEClient) callRerankService(ctx context.Context, request RerankRequest) ([]RerankResult, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", c.fallbackURL+"/rerank", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call rerank service: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rerank service returned status %d", resp.StatusCode)
	}
	
	var response RerankResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	return response.Results, nil
}

// selectTopK sorts results by score and returns top K
func (c *PureGoBGEClient) selectTopK(results []RerankResult, topK int) []RerankResult {
	// Sort by score (descending)
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].Score > results[i].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
	
	// Return top K
	if topK > 0 && topK < len(results) {
		return results[:topK]
	}
	
	return results
}

// SetMaxLength sets the maximum token length for tokenization
func (c *PureGoBGEClient) SetMaxLength(maxLength int) {
	c.maxLength = maxLength
}

// GetTokenizer returns the underlying tokenizer for direct access
func (c *PureGoBGEClient) GetTokenizer() *PureGoTokenizer {
	return c.tokenizer
}

// Health checks if the client is ready
func (c *PureGoBGEClient) Health(ctx context.Context) error {
	if c.tokenizer == nil {
		return fmt.Errorf("tokenizer not initialized")
	}
	
	// Test tokenization
	_, _, _ = c.tokenizer.Encode("test", 10)
	
	return nil
}