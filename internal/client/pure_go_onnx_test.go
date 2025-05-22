package client

import (
	"os"
	"path/filepath"
	"testing"

	ort "github.com/yalue/onnxruntime_go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPureGoONNXRuntime(t *testing.T) {
	// Skip test if model directory doesn't exist
	modelDir := filepath.Join("..", "..", "models", "bge-reranker-large-onnx")
	modelPath := filepath.Join(modelDir, "model.onnx")
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Skip("Model directory not found, skipping pure Go ONNX test")
	}

	// Set the library path to our downloaded ONNX runtime (shared across all tests)
	libPath := filepath.Join("..", "..", "lib", "onnxruntime-osx-arm64-1.15.0", "lib", "libonnxruntime.dylib")
	ort.SetSharedLibraryPath(libPath)

	t.Run("InitializeONNXEnvironment", func(t *testing.T) {
		err := ort.InitializeEnvironment()
		require.NoError(t, err, "Should be able to initialize ONNX environment")
		
		defer func() {
			err := ort.DestroyEnvironment()
			assert.NoError(t, err, "Should be able to destroy ONNX environment")
		}()
	})

	t.Run("LoadModelAndInspectInputs", func(t *testing.T) {
		err := ort.InitializeEnvironment()
		require.NoError(t, err)
		defer ort.DestroyEnvironment()

		// Use DynamicAdvancedSession for mixed input/output types
		session, err := ort.NewDynamicAdvancedSession(
			modelPath,
			[]string{"input_ids", "attention_mask"},
			[]string{"logits"},
			nil, // Use default session options
		)
		require.NoError(t, err, "Should create ONNX session")
		defer session.Destroy()

		t.Logf("Successfully created ONNX session with model: %s", modelPath)
	})

	t.Run("BasicInferenceWithDummyData", func(t *testing.T) {
		err := ort.InitializeEnvironment()
		require.NoError(t, err)
		defer ort.DestroyEnvironment()

		session, err := ort.NewDynamicAdvancedSession(
			modelPath,
			[]string{"input_ids", "attention_mask"},
			[]string{"logits"},
			nil,
		)
		require.NoError(t, err)
		defer session.Destroy()

		batchSize := int64(1)
		seqLength := int64(512)
		
		inputShape := ort.NewShape(batchSize, seqLength)
		
		inputIds, err := ort.NewEmptyTensor[int64](inputShape)
		require.NoError(t, err)
		defer inputIds.Destroy()
		
		attentionMask, err := ort.NewEmptyTensor[int64](inputShape)
		require.NoError(t, err)
		defer attentionMask.Destroy()
		
		outputShape := ort.NewShape(batchSize, 1)
		output, err := ort.NewEmptyTensor[float32](outputShape)
		require.NoError(t, err)
		defer output.Destroy()

		// Fill input tensors with dummy data
		inputIdsData := inputIds.GetData()
		attentionMaskData := attentionMask.GetData()
		
		// Simple dummy tokenization:
		// token_ids: [0, 1, 2, 3, ..., 10] + padding with 1 (pad token)
		// attention_mask: [1, 1, 1, ...] for real tokens, [0, 0, ...] for padding
		for i := 0; i < len(inputIdsData); i++ {
			if i < 10 {
				inputIdsData[i] = int64(i)  // Some dummy token IDs
				attentionMaskData[i] = 1    // Attention for real tokens
			} else {
				inputIdsData[i] = 1         // Pad token ID
				attentionMaskData[i] = 0    // No attention for padding
			}
		}

		// Prepare input/output arrays for dynamic session
		inputs := []ort.ArbitraryTensor{inputIds, attentionMask}
		outputs := []ort.ArbitraryTensor{output}

		// Run inference
		err = session.Run(inputs, outputs)
		require.NoError(t, err, "Should run inference successfully")

		// Check output
		outputData := output.GetData()
		require.Len(t, outputData, 1, "Should have 1 output value")
		
		logits := outputData[0]
		t.Logf("Raw logits: %f", logits)
		
		// Convert to probability using sigmoid
		score := 1.0 / (1.0 + float64(-logits))
		t.Logf("Sigmoid score: %f", score)
		
		// Score should be a reasonable probability (0-1)
		assert.True(t, score >= 0.0 && score <= 1.0, "Score should be between 0 and 1")
	})
}

func TestONNXRuntimeCapabilities(t *testing.T) {
	// Set the library path to our downloaded ONNX runtime
	libPath := filepath.Join("..", "..", "lib", "onnxruntime-osx-arm64-1.15.0", "lib", "libonnxruntime.dylib")
	ort.SetSharedLibraryPath(libPath)

	t.Run("TensorOperations", func(t *testing.T) {
		err := ort.InitializeEnvironment()
		require.NoError(t, err)
		defer ort.DestroyEnvironment()

		// Test tensor creation and manipulation
		shape := ort.NewShape(2, 3)
		tensor, err := ort.NewEmptyTensor[int64](shape)
		require.NoError(t, err, "Should create tensor")

		// Test data access
		data := tensor.GetData()
		assert.Len(t, data, 6, "Should have 6 elements (2x3)")

		// Test data modification
		for i := range data {
			data[i] = int64(i + 1)
		}

		// Verify data was set
		assert.Equal(t, int64(1), data[0])
		assert.Equal(t, int64(6), data[5])

		tensor.Destroy()
	})

	t.Run("MultipleDataTypes", func(t *testing.T) {
		err := ort.InitializeEnvironment()
		require.NoError(t, err)
		defer ort.DestroyEnvironment()

		shape := ort.NewShape(2, 2)

		// Test int64 tensors (for input_ids, attention_mask)
		intTensor, err := ort.NewEmptyTensor[int64](shape)
		require.NoError(t, err)
		intTensor.Destroy()

		// Test float32 tensors (for outputs)
		floatTensor, err := ort.NewEmptyTensor[float32](shape)
		require.NoError(t, err)
		floatTensor.Destroy()
	})
}