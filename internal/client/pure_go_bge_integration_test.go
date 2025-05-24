package client

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPureGoBGEClient_WithONNXInference(t *testing.T) {
	// Skip test if model directory doesn't exist
	modelDir := filepath.Join("..", "..", "models", "bge-reranker-large-onnx")
	if _, err := os.Stat(filepath.Join(modelDir, "model.onnx")); os.IsNotExist(err) {
		t.Skip("Model directory not found, skipping pure Go BGE client test")
	}

	t.Run("CreatePureGoBGEClient", func(t *testing.T) {
		client, err := NewPureGoBGEClient(modelDir, true, "")
		require.NoError(t, err, "Should create pure Go BGE client")
		defer client.Close()

		err = client.Health(context.Background())
		assert.NoError(t, err, "Client should be healthy")
	})

	t.Run("RerankWithPureGoONNX", func(t *testing.T) {
		client, err := NewPureGoBGEClient(modelDir, true, "")
		require.NoError(t, err)
		defer client.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		query := "What is machine learning?"
		documents := []string{
			"Machine learning is a subset of artificial intelligence that enables computers to learn without explicit programming.",
			"Cooking involves preparing food by combining ingredients with heat and various techniques.",
			"Deep learning uses neural networks with multiple layers to process complex patterns in data.",
			"Basketball is a sport played between two teams of five players on a court.",
		}

		results, err := client.Rerank(ctx, query, documents, 2)
		require.NoError(t, err, "Should rerank successfully")

		assert.Len(t, results, 2, "Should return top 2 results")
		
		// First result should be most relevant (machine learning related)
		assert.Contains(t, results[0].Document, "Machine learning", 
			"First result should be most relevant to machine learning")
		
		// Scores should be between 0 and 1
		for i, result := range results {
			assert.True(t, result.Score >= 0.0 && result.Score <= 1.0, 
				"Score %d should be between 0 and 1, got: %f", i, result.Score)
			assert.Equal(t, result.Score, result.RelevanceScore, 
				"Score and RelevanceScore should be equal")
			t.Logf("Result %d: Score=%.6f, Document=%s", 
				result.Index, result.Score, result.Document[:50]+"...")
		}
		
		// Results should be sorted by score (descending)
		if len(results) > 1 {
			assert.True(t, results[0].Score >= results[1].Score, 
				"Results should be sorted by score (descending)")
		}
	})

	t.Run("RerankWithVariousQueries", func(t *testing.T) {
		client, err := NewPureGoBGEClient(modelDir, true, "")
		require.NoError(t, err)
		defer client.Close()

		documents := []string{
			"Python is a programming language used for web development, data science, and automation.",
			"The Python snake is a non-venomous constrictor found in Africa and Asia.",
			"JavaScript is a scripting language primarily used for web development.",
			"Java is an object-oriented programming language used for enterprise applications.",
		}

		testCases := []struct {
			query             string
			expectedTopResult string
		}{
			{
				query:             "programming languages for web development",
				expectedTopResult: "Python is a programming language",
			},
			{
				query:             "animals and wildlife",
				expectedTopResult: "Python snake",
			},
			{
				query:             "enterprise software development",
				expectedTopResult: "Java is an object-oriented",
			},
		}

		for _, tc := range testCases {
			t.Run("Query: "+tc.query, func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				results, err := client.Rerank(ctx, tc.query, documents, 1)
				require.NoError(t, err, "Should rerank successfully")
				require.Len(t, results, 1, "Should return 1 result")

				assert.Contains(t, results[0].Document, tc.expectedTopResult,
					"Top result should be relevant to query")
				
				t.Logf("Query: %s", tc.query)
				t.Logf("Top result (score=%.6f): %s", 
					results[0].Score, results[0].Document)
			})
		}
	})

	t.Run("EmptyDocuments", func(t *testing.T) {
		client, err := NewPureGoBGEClient(modelDir, true, "")
		require.NoError(t, err)
		defer client.Close()

		ctx := context.Background()
		results, err := client.Rerank(ctx, "test query", []string{}, 5)
		require.NoError(t, err, "Should handle empty documents")
		assert.Len(t, results, 0, "Should return empty results")
	})

	t.Run("SingleDocument", func(t *testing.T) {
		client, err := NewPureGoBGEClient(modelDir, true, "")
		require.NoError(t, err)
		defer client.Close()

		ctx := context.Background()
		documents := []string{"This is a single test document."}
		
		results, err := client.Rerank(ctx, "test", documents, 1)
		require.NoError(t, err, "Should handle single document")
		assert.Len(t, results, 1, "Should return 1 result")
		assert.Equal(t, documents[0], results[0].Document)
	})

	t.Run("TokenizerPerformance", func(t *testing.T) {
		client, err := NewPureGoBGEClient(modelDir, true, "")
		require.NoError(t, err)
		defer client.Close()

		// Test tokenizer performance with rapid successive calls
		query := "Performance test query for tokenizer speed measurement"
		documents := []string{
			"Document 1: Performance testing is crucial for ensuring application scalability.",
			"Document 2: Speed optimization requires careful analysis of bottlenecks.",
			"Document 3: Tokenization performance directly impacts inference latency.",
		}

		start := time.Now()
		numRuns := 10

		for i := 0; i < numRuns; i++ {
			ctx := context.Background()
			_, err := client.Rerank(ctx, query, documents, 3)
			require.NoError(t, err, "Rerank should succeed on run %d", i)
		}

		elapsed := time.Since(start)
		avgTime := elapsed / time.Duration(numRuns)
		
		t.Logf("Performed %d rerank operations in %v (avg: %v per operation)", 
			numRuns, elapsed, avgTime)
		
		// Should be reasonably fast (less than 1 second per operation on average)
		assert.True(t, avgTime < time.Second, 
			"Average rerank time should be less than 1 second, got: %v", avgTime)
	})
}

func TestPureGoBGEClient_Fallback(t *testing.T) {
	modelDir := filepath.Join("..", "..", "models", "bge-reranker-large-onnx")
	if _, err := os.Stat(filepath.Join(modelDir, "model.onnx")); os.IsNotExist(err) {
		t.Skip("Model directory not found, skipping fallback test")
	}

	t.Run("FallbackMode", func(t *testing.T) {
		// Create client without pure Go ONNX (fallback mode)
		client, err := NewPureGoBGEClient(modelDir, false, "http://localhost:8001")
		require.NoError(t, err, "Should create BGE client in fallback mode")
		defer client.Close()

		// Health check should pass even without ONNX inference
		err = client.Health(context.Background())
		assert.NoError(t, err, "Health check should pass in fallback mode")

		// Note: Actual reranking would require the Python service to be running
		// So we only test client creation and health check here
	})

	t.Run("PureGoModeWithoutONNX", func(t *testing.T) {
		// Test creating pure Go client with invalid model path
		// This should fail during ONNX initialization
		invalidModelDir := "/nonexistent/path"
		_, err := NewPureGoBGEClient(invalidModelDir, true, "")
		assert.Error(t, err, "Should fail with invalid model path")
	})
}