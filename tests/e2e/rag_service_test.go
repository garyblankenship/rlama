package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/repository"
	"github.com/dontizi/rlama/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestRagServiceOperations(t *testing.T) {
	// Create a temporary test directory
	tempDir, err := os.MkdirTemp("", "rag-service-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	files := map[string]string{
		"test.txt": "This is a test document.",
		"test.md":  "# Test\nThis is a markdown file.",
	}

	for path, content := range files {
		fullPath := filepath.Join(tempDir, path)
		err := os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", path, err)
		}
	}

	t.Run("CreateAndQueryRag", func(t *testing.T) {
		// Clean up any existing test RAG
		repo := repository.NewRagRepository()
		repo.Delete("test-rag")

		// Create Ollama client
		ollamaClient := client.NewDefaultOllamaClient()

		// Use real embedding model - "bge-large" from Ollama
		embeddingModel := "bge-large"
		completionModel := "llama3.2"

		// Create embedding service with real Ollama client
		embeddingService := service.NewEmbeddingService(ollamaClient)

		// Create RAG service with real embedding
		ragService := service.NewRagServiceWithEmbedding(ollamaClient, embeddingService)

		// Create RAG with options
		options := service.DocumentLoaderOptions{
			ChunkSize:      500,
			ChunkOverlap:   50,
			EnableReranker: true,
			RerankerModel:  completionModel,
		}

		err := ragService.CreateRagWithOptions(embeddingModel, "test-rag", tempDir, options)
		assert.NoError(t, err)

		// Test listing chunks
		filter := service.ChunkFilter{
			DocumentSubstring: "test",
			ShowContent:       true,
		}

		chunks, err := ragService.GetRagChunks("test-rag", filter)
		assert.NoError(t, err)
		assert.NotEmpty(t, chunks)

		// Load the RAG from repository
		rag, err := ragService.LoadRag("test-rag")
		assert.NoError(t, err)

		// Update model name to use llama3.2 for completion
		oldModel := rag.ModelName
		rag.ModelName = completionModel
		err = ragService.UpdateRag(rag)
		assert.NoError(t, err)

		// Query the RAG (this will use real Ollama)
		result, err := ragService.Query(rag, "What is in the test documents?", 5)
		assert.NoError(t, err)
		assert.NotEmpty(t, result)

		// Restore model
		rag.ModelName = oldModel
		err = ragService.UpdateRag(rag)
		assert.NoError(t, err)

		// Print result
		t.Logf("RAG Query Result: %s", result)
	})
}
