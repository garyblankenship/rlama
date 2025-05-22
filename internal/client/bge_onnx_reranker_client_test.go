package client

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBGEONNXRerankerClient(t *testing.T) {
	// Skip test if model directory doesn't exist
	modelDir := filepath.Join("..", "..", "models", "bge-reranker-large-onnx")
	if _, err := os.Stat(modelDir); os.IsNotExist(err) {
		t.Skip("Model directory not found, skipping ONNX reranker test")
	}

	t.Run("NewBGEONNXRerankerClient", func(t *testing.T) {
		client, err := NewBGEONNXRerankerClient(modelDir)
		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, "bge-reranker-large-onnx", client.GetModelName())
		
		// Cleanup
		defer client.Cleanup()
	})

	t.Run("ComputeScores_SinglePair", func(t *testing.T) {
		client, err := NewBGEONNXRerankerClient(modelDir)
		require.NoError(t, err)
		defer client.Cleanup()

		pairs := [][]string{
			{"What is a cat?", "A cat is a small domesticated carnivorous mammal."},
		}

		scores, err := client.ComputeScores(pairs, true)
		require.NoError(t, err)
		require.Len(t, scores, 1)
		
		// Score should be between 0 and 1 when normalized
		assert.Greater(t, scores[0], 0.0)
		assert.Less(t, scores[0], 1.0)
	})

	t.Run("ComputeScores_MultiplePairs", func(t *testing.T) {
		client, err := NewBGEONNXRerankerClient(modelDir)
		require.NoError(t, err)
		defer client.Cleanup()

		pairs := [][]string{
			{"What is a cat?", "A cat is a small domesticated carnivorous mammal."},
			{"What is a cat?", "The weather is nice today."},
			{"How to cook pasta?", "Boil water, add pasta, cook for 8-10 minutes."},
		}

		scores, err := client.ComputeScores(pairs, true)
		require.NoError(t, err)
		require.Len(t, scores, 3)
		
		// First pair should have higher score than second pair (more relevant)
		assert.Greater(t, scores[0], scores[1])
		
		// All scores should be normalized between 0 and 1
		for i, score := range scores {
			assert.Greater(t, score, 0.0, "Score %d should be > 0", i)
			assert.Less(t, score, 1.0, "Score %d should be < 1", i)
		}
	})

	t.Run("ComputeScores_WithoutNormalization", func(t *testing.T) {
		client, err := NewBGEONNXRerankerClient(modelDir)
		require.NoError(t, err)
		defer client.Cleanup()

		pairs := [][]string{
			{"What is a cat?", "A cat is a small domesticated carnivorous mammal."},
		}

		scores, err := client.ComputeScores(pairs, false)
		require.NoError(t, err)
		require.Len(t, scores, 1)
		
		// Without normalization, scores can be any real number (logits)
		assert.NotZero(t, scores[0])
	})

	t.Run("ComputeScores_InvalidPair", func(t *testing.T) {
		client, err := NewBGEONNXRerankerClient(modelDir)
		require.NoError(t, err)
		defer client.Cleanup()

		pairs := [][]string{
			{"single element"},
		}

		_, err = client.ComputeScores(pairs, true)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exactly 2 elements")
	})
}

func TestBGEONNXRerankerClient_Performance(t *testing.T) {
	// Skip test if model directory doesn't exist
	modelDir := filepath.Join("..", "..", "models", "bge-reranker-large-onnx")
	if _, err := os.Stat(modelDir); os.IsNotExist(err) {
		t.Skip("Model directory not found, skipping performance test")
	}

	client, err := NewBGEONNXRerankerClient(modelDir)
	require.NoError(t, err)
	defer client.Cleanup()

	// Test performance with multiple pairs
	pairs := [][]string{
		{"machine learning", "Machine learning is a subset of artificial intelligence"},
		{"machine learning", "I like to eat pizza on weekends"},
		{"cats and dogs", "Cats are independent pets while dogs are loyal companions"},
		{"cats and dogs", "Weather forecast shows rain tomorrow"},
	}

	start := time.Now()
	scores, err := client.ComputeScores(pairs, true)
	duration := time.Since(start)
	
	require.NoError(t, err)
	require.Len(t, scores, 4)
	
	// Should be faster than the original Python subprocess approach
	assert.Less(t, duration.Seconds(), 5.0, "Should complete within 5 seconds")
	
	// Relevant pairs should score higher than irrelevant ones
	assert.Greater(t, scores[0], scores[1], "ML pair should score higher than pizza pair")
	assert.Greater(t, scores[2], scores[3], "Pets pair should score higher than weather pair")
	
	t.Logf("Processed %d pairs in %v (avg: %v per pair)", 
		len(pairs), duration, duration/time.Duration(len(pairs)))
}