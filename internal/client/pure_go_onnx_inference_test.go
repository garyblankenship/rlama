package client

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPureGoONNXInference_Integration(t *testing.T) {
	// Skip test if model directory doesn't exist
	modelDir := filepath.Join("..", "..", "models", "bge-reranker-large-onnx")
	modelPath := filepath.Join(modelDir, "model.onnx")
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Model directory not found, skipping pure Go ONNX inference test")
	}

	t.Run("CreateONNXInference", func(t *testing.T) {
		inference, err := NewPureGoONNXInference(modelDir)
		require.NoError(t, err, "Should create ONNX inference client")
		defer inference.Close()

		assert.True(t, inference.IsInitialized(), "Should be initialized")
	})

	t.Run("RunInferenceWithDummyData", func(t *testing.T) {
		inference, err := NewPureGoONNXInference(modelDir)
		require.NoError(t, err)
		defer inference.Close()

		// Create dummy tokenized input (similar to what the tokenizer would produce)
		inputIDs := []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
		attentionMask := []int64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
		
		// Pad to 512 tokens (standard BGE input length)
		for len(inputIDs) < 512 {
			inputIDs = append(inputIDs, 1)    // Pad token
			attentionMask = append(attentionMask, 0) // No attention for padding
		}

		request := ONNXInferenceRequest{
			InputIDs:      [][]int64{inputIDs},
			AttentionMask: [][]int64{attentionMask},
		}

		response, err := inference.RunInference(request)
		require.NoError(t, err, "Should run inference successfully")

		assert.Len(t, response.Scores, 1, "Should have 1 score")
		score := response.Scores[0]
		assert.True(t, score >= 0.0 && score <= 1.0, "Score should be between 0 and 1, got: %f", score)

		t.Logf("Inference score: %f", score)
	})

	t.Run("RunBatchInference", func(t *testing.T) {
		inference, err := NewPureGoONNXInference(modelDir)
		require.NoError(t, err)
		defer inference.Close()

		// Create multiple dummy tokenized inputs
		batchSize := 3
		inputIDs := make([][]int64, batchSize)
		attentionMasks := make([][]int64, batchSize)

		for i := 0; i < batchSize; i++ {
			ids := make([]int64, 512)
			mask := make([]int64, 512)
			
			// Fill first 10 tokens with different patterns
			for j := 0; j < 10; j++ {
				ids[j] = int64(j + i*10) // Different patterns for each batch item
				mask[j] = 1
			}
			
			// Rest are padding
			for j := 10; j < 512; j++ {
				ids[j] = 1  // Pad token
				mask[j] = 0 // No attention
			}
			
			inputIDs[i] = ids
			attentionMasks[i] = mask
		}

		request := ONNXInferenceRequest{
			InputIDs:      inputIDs,
			AttentionMask: attentionMasks,
		}

		response, err := inference.RunInference(request)
		require.NoError(t, err, "Should run batch inference successfully")

		assert.Len(t, response.Scores, batchSize, "Should have %d scores", batchSize)
		
		for i, score := range response.Scores {
			assert.True(t, score >= 0.0 && score <= 1.0, "Score %d should be between 0 and 1, got: %f", i, score)
			t.Logf("Batch item %d score: %f", i, score)
		}
	})

	t.Run("EmptyInput", func(t *testing.T) {
		inference, err := NewPureGoONNXInference(modelDir)
		require.NoError(t, err)
		defer inference.Close()

		request := ONNXInferenceRequest{
			InputIDs:      [][]int64{},
			AttentionMask: [][]int64{},
		}

		response, err := inference.RunInference(request)
		require.NoError(t, err, "Should handle empty input")
		assert.Len(t, response.Scores, 0, "Should have 0 scores for empty input")
	})

	t.Run("CleanupResources", func(t *testing.T) {
		inference, err := NewPureGoONNXInference(modelDir)
		require.NoError(t, err)

		assert.True(t, inference.IsInitialized(), "Should be initialized")

		err = inference.Close()
		assert.NoError(t, err, "Should close without error")

		assert.False(t, inference.IsInitialized(), "Should not be initialized after close")
	})
}

func TestPureGoONNXInference_EdgeCases(t *testing.T) {
	modelDir := filepath.Join("..", "..", "models", "bge-reranker-large-onnx")
	if _, err := os.Stat(filepath.Join(modelDir, "model.onnx")); os.IsNotExist(err) {
		t.Skip("Model directory not found, skipping edge case tests")
	}

	t.Run("InvalidModelPath", func(t *testing.T) {
		_, err := NewPureGoONNXInference("/nonexistent/path")
		assert.Error(t, err, "Should fail with invalid model path")
	})

	t.Run("DoubleClose", func(t *testing.T) {
		inference, err := NewPureGoONNXInference(modelDir)
		require.NoError(t, err)

		err = inference.Close()
		assert.NoError(t, err, "First close should succeed")

		err = inference.Close()
		assert.NoError(t, err, "Second close should also succeed (no-op)")
	})

	t.Run("UseAfterClose", func(t *testing.T) {
		inference, err := NewPureGoONNXInference(modelDir)
		require.NoError(t, err)

		err = inference.Close()
		require.NoError(t, err)

		request := ONNXInferenceRequest{
			InputIDs:      [][]int64{{1, 2, 3}},
			AttentionMask: [][]int64{{1, 1, 1}},
		}

		_, err = inference.RunInference(request)
		assert.Error(t, err, "Should fail when used after close")
	})
}