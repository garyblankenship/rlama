package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestPureGoBGEClient_EndToEndReranking(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create mock inference server
	server := createMockInferenceServer(t)
	defer server.Close()

	modelPath := createTempModelDir(t)
	defer os.RemoveAll(modelPath)

	client, err := NewPureGoBGEClient(modelPath, true, server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := "What is machine learning?"
	documents := []string{
		"Machine learning is a subset of artificial intelligence that focuses on algorithms.",
		"Deep learning uses neural networks with multiple layers.",
		"Natural language processing deals with human language.",
		"Computer vision enables machines to interpret visual information.",
		"Reinforcement learning learns through interaction with environment.",
	}

	results, err := client.Rerank(ctx, query, documents, 3)
	if err != nil {
		t.Fatalf("Reranking failed: %v", err)
	}

	// Validate results
	if len(results) != 3 {
		t.Errorf("Expected 3 results, got %d", len(results))
	}

	for i, result := range results {
		if result.Index < 0 || result.Index >= len(documents) {
			t.Errorf("Result %d has invalid index: %d", i, result.Index)
		}

		if result.Document != documents[result.Index] {
			t.Errorf("Result %d document mismatch", i)
		}

		if result.Score < 0 || result.Score > 1 {
			t.Errorf("Result %d has invalid score: %f", i, result.Score)
		}

		// Results should be sorted by score (descending)
		if i > 0 && result.Score > results[i-1].Score {
			t.Errorf("Results not sorted by score: result %d score %f > result %d score %f", 
				i, result.Score, i-1, results[i-1].Score)
		}

		t.Logf("Result %d: Index=%d, Score=%f, Doc=%s", 
			i, result.Index, result.Score, result.Document[:50]+"...")
	}
}

func TestPureGoBGEClient_FallbackMode(t *testing.T) {
	// Create mock rerank server (simulates original Python service)
	server := createMockRerankServer(t)
	defer server.Close()

	modelPath := createTempModelDir(t)
	defer os.RemoveAll(modelPath)

	// Create client in fallback mode
	client, err := NewPureGoBGEClient(modelPath, false, server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := "Test query"
	documents := []string{"Doc 1", "Doc 2", "Doc 3"}

	results, err := client.Rerank(ctx, query, documents, 2)
	if err != nil {
		t.Fatalf("Fallback reranking failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	t.Logf("Fallback mode test passed with %d results", len(results))
}

func TestPureGoBGEClient_BatchProcessing(t *testing.T) {
	server := createMockInferenceServer(t)
	defer server.Close()

	modelPath := createTempModelDir(t)
	defer os.RemoveAll(modelPath)

	client, err := NewPureGoBGEClient(modelPath, true, server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with varying batch sizes
	batchSizes := []int{1, 5, 10, 20, 50}
	query := "Test query for batch processing"

	for _, batchSize := range batchSizes {
		t.Run(fmt.Sprintf("BatchSize%d", batchSize), func(t *testing.T) {
			// Generate documents
			documents := make([]string, batchSize)
			for i := range documents {
				documents[i] = fmt.Sprintf("Document %d content for testing batch processing", i)
			}

			startTime := time.Now()
			results, err := client.Rerank(ctx, query, documents, -1) // Get all results
			duration := time.Since(startTime)

			if err != nil {
				t.Fatalf("Batch processing failed for size %d: %v", batchSize, err)
			}

			if len(results) != batchSize {
				t.Errorf("Expected %d results, got %d", batchSize, len(results))
			}

			t.Logf("Batch size %d: %d results in %v", batchSize, len(results), duration)
		})
	}
}

func TestPureGoBGEClient_DifferentMaxLengths(t *testing.T) {
	server := createMockInferenceServer(t)
	defer server.Close()

	modelPath := createTempModelDir(t)
	defer os.RemoveAll(modelPath)

	client, err := NewPureGoBGEClient(modelPath, true, server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	maxLengths := []int{64, 128, 256, 512, 1024}
	query := "Test query with variable length"
	documents := []string{
		"Short document",
		"Medium length document with more content to test tokenization",
		"Very long document with extensive content that will test the tokenization limits and see how the system handles longer sequences of text that might exceed certain length boundaries",
	}

	for _, maxLen := range maxLengths {
		t.Run(fmt.Sprintf("MaxLength%d", maxLen), func(t *testing.T) {
			client.SetMaxLength(maxLen)

			results, err := client.Rerank(ctx, query, documents, 3)
			if err != nil {
				t.Fatalf("Reranking failed with max length %d: %v", maxLen, err)
			}

			if len(results) != len(documents) {
				t.Errorf("Expected %d results, got %d", len(documents), len(results))
			}

			t.Logf("Max length %d: Successfully processed %d documents", maxLen, len(results))
		})
	}
}

func TestPureGoBGEClient_TokenizationConsistency(t *testing.T) {
	modelPath := createTempModelDir(t)
	defer os.RemoveAll(modelPath)

	client, err := NewPureGoBGEClient(modelPath, true, "http://localhost:8000")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tokenizer := client.GetTokenizer()
	testCases := []struct {
		name    string
		query   string
		passage string
		maxLen  int
	}{
		{
			name:    "ShortTexts",
			query:   "AI",
			passage: "Artificial Intelligence",
			maxLen:  128,
		},
		{
			name:    "MediumTexts", 
			query:   "What is machine learning?",
			passage: "Machine learning is a subset of AI.",
			maxLen:  256,
		},
		{
			name:    "LongTexts",
			query:   "Explain the differences between supervised and unsupervised learning",
			passage: "Supervised learning uses labeled data while unsupervised learning finds patterns in unlabeled data.",
			maxLen:  512,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test consistency across multiple runs
			var firstTokenIDs []int64
			var firstAttention []int64
			var firstTokenTypes []int64

			for run := 0; run < 5; run++ {
				tokenIDs, attention, tokenTypes := tokenizer.EncodeQueryPassagePair(tc.query, tc.passage, tc.maxLen)

				if run == 0 {
					firstTokenIDs = tokenIDs
					firstAttention = attention  
					firstTokenTypes = tokenTypes
				} else {
					// Verify consistency
					if !slicesEqual(tokenIDs, firstTokenIDs) {
						t.Errorf("Run %d: Token IDs differ from first run", run)
					}
					if !slicesEqual(attention, firstAttention) {
						t.Errorf("Run %d: Attention mask differs from first run", run)
					}
					if !slicesEqual(tokenTypes, firstTokenTypes) {
						t.Errorf("Run %d: Token types differ from first run", run)
					}
				}
			}

			t.Logf("Tokenization consistency verified for %s", tc.name)
		})
	}
}

// Mock servers for testing

func createMockInferenceServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/inference" {
			var req BGEInferenceRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Generate mock scores based on input size
			scores := make([]float64, len(req.InputIDs))
			for i := range scores {
				// Mock scoring: higher scores for shorter sequences (simulating relevance)
				actualTokens := 0
				for _, mask := range req.AttentionMask[i] {
					if mask == 1 {
						actualTokens++
					}
				}
				// Normalize score based on token count (shorter = more relevant)
				scores[i] = 1.0 - float64(actualTokens)/512.0
				if scores[i] < 0.1 {
					scores[i] = 0.1
				}
			}

			response := BGEInferenceResponse{Scores: scores}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func createMockRerankServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/rerank" {
			var req RerankRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// Generate mock results
			results := make([]RerankResult, len(req.Documents))
			for i, doc := range req.Documents {
				results[i] = RerankResult{
					Index:          i,
					Document:       doc,
					Score:          0.9 - float64(i)*0.1, // Decreasing scores
					RelevanceScore: 0.9 - float64(i)*0.1,
				}
			}

			// Apply top-k if specified
			if req.TopK > 0 && req.TopK < len(results) {
				results = results[:req.TopK]
			}

			response := RerankResponse{Results: results}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

// Helper functions

func slicesEqual(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}