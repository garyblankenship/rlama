package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dontizi/rlama/internal/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestONNXRerankerIntegration(t *testing.T) {
	// Skip test if model directory doesn't exist
	modelDir := filepath.Join("..", "..", "models", "bge-reranker-large-onnx")
	if _, err := os.Stat(modelDir); os.IsNotExist(err) {
		t.Skip("Model directory not found, skipping ONNX integration test")
	}

	t.Run("ServiceConfig_ONNXConfiguration", func(t *testing.T) {
		config := NewServiceConfig()
		config.UseONNXReranker = true
		config.ONNXModelDir = modelDir

		// Validate configuration
		err := config.Validate()
		require.NoError(t, err)

		// Check that DocumentLoaderOptions includes ONNX settings
		options := config.ToDocumentLoaderOptions()
		assert.True(t, options.UseONNXReranker)
		assert.Equal(t, modelDir, options.ONNXModelDir)
	})

	t.Run("QueryService_WithONNXConfig", func(t *testing.T) {
		config := &ServiceConfig{
			UseONNXReranker: true,
			ONNXModelDir:    modelDir,
		}

		ollamaClient := client.NewDefaultOllamaClient()
		llmClient := client.NewOllamaClient("localhost", "11434")
		documentService := NewDocumentService(llmClient)

		queryService := NewQueryServiceWithConfig(llmClient, ollamaClient, documentService, config)
		require.NotNil(t, queryService)

		// Verify that the query service was created successfully
		impl, ok := queryService.(*QueryServiceImpl)
		require.True(t, ok)
		assert.NotNil(t, impl.rerankerService)
		
		// Check if it's using ONNX (this is indirect since we can't easily check the internal state)
		assert.True(t, impl.rerankerService.IsUsingONNX())
	})

	t.Run("CompositeRagService_WithONNXConfig", func(t *testing.T) {
		config := &ServiceConfig{
			UseONNXReranker: true,
			ONNXModelDir:    modelDir,
		}

		ollamaClient := client.NewDefaultOllamaClient()
		llmClient := client.NewOllamaClient("localhost", "11434")

		ragService := NewCompositeRagServiceWithConfig(llmClient, ollamaClient, config)
		require.NotNil(t, ragService)

		// Verify service creation was successful
		impl, ok := ragService.(*CompositeRagService)
		require.True(t, ok)
		assert.NotNil(t, impl.queryService)
	})

	t.Run("ServiceProvider_WithONNXConfig", func(t *testing.T) {
		config := NewServiceConfig()
		config.UseONNXReranker = true
		config.ONNXModelDir = modelDir

		provider, err := NewServiceProvider(config)
		require.NoError(t, err)

		// Test creating a RAG service for a model
		ragService, err := provider.CreateRagServiceForModel("llama3.2")
		require.NoError(t, err)
		assert.NotNil(t, ragService)
	})
}

func TestRerankerServiceInterface(t *testing.T) {
	// Skip test if model directory doesn't exist
	modelDir := filepath.Join("..", "..", "models", "bge-reranker-large-onnx")
	if _, err := os.Stat(modelDir); os.IsNotExist(err) {
		t.Skip("Model directory not found, skipping reranker interface test")
	}

	t.Run("BGERerankerClient_ImplementsInterface", func(t *testing.T) {
		// Test that both implementations satisfy the RerankerClient interface
		pythonClient := client.NewBGERerankerClient("BAAI/bge-reranker-v2-m3")
		var _ RerankerClient = pythonClient

		onnxClient, err := client.NewBGEONNXRerankerClient(modelDir)
		require.NoError(t, err)
		defer onnxClient.Cleanup()
		var _ RerankerClient = onnxClient
		var _ CleanupableRerankerClient = onnxClient
	})

	t.Run("RerankerService_Cleanup", func(t *testing.T) {
		ollamaClient := client.NewDefaultOllamaClient()
		
		// Test with ONNX reranker (which needs cleanup)
		rerankerService := NewRerankerServiceWithOptions(ollamaClient, true, modelDir)
		require.NotNil(t, rerankerService)
		
		// Cleanup should not error
		err := rerankerService.Cleanup()
		assert.NoError(t, err)
		
		// Test with Python reranker (which doesn't need cleanup)
		rerankerService2 := NewRerankerServiceWithOptions(ollamaClient, false, modelDir)
		require.NotNil(t, rerankerService2)
		
		// Cleanup should not error even when not needed
		err = rerankerService2.Cleanup()
		assert.NoError(t, err)
	})
}