package tests

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/repository"
	"github.com/dontizi/rlama/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

		// Create mock Ollama client
		mockOllama := new(MockOllamaClient)

		// Configure the mock expectations
		embeddingModel := "bge-large"
		completionModel := "llama3.2"

		// Setup mock expectations
		mockOllama.On("CheckOllamaAndModel", mock.Anything).Return(nil)
		mockOllama.On("CheckLLMAndModel", mock.Anything).Return(nil)

		// Mock embedding generation
		mockEmbedding := []float32{0.1, 0.2, 0.3, 0.4, 0.5}
		mockOllama.On("GenerateEmbedding", mock.Anything, mock.Anything).Return(mockEmbedding, nil)

		// Mock completion generation
		mockResponse := "The test documents contain text files and a markdown file with headings."
		mockOllama.On("GenerateCompletion", mock.Anything, mock.Anything).Return(mockResponse, nil)

		// Create embedding service mock
		mockEmbeddingService := new(MockEmbeddingService)
		mockEmbeddingService.On("GenerateChunkEmbeddings", mock.Anything, mock.Anything).Return(nil)
		mockEmbeddingService.On("GenerateQueryEmbedding", mock.Anything, mock.Anything).Return(mockEmbedding, nil)

		// Create our test service - aucun besoin de BGE reranker ici car notre mock est complet
		testRagService := &TestRagService{
			mockOllama:    mockOllama,
			mockEmbedding: mockEmbeddingService,
			ragRepository: repository.NewRagRepository(),
		}

		// Create RAG with options
		options := service.DocumentLoaderOptions{
			ChunkSize:      500,
			ChunkOverlap:   50,
			EnableReranker: true,
			RerankerModel:  completionModel,
		}

		// Create a new RAG - utilise la méthode simplifiée de notre mock
		err := testRagService.CreateRagWithOptions(embeddingModel, "test-rag", tempDir, options)
		assert.NoError(t, err)

		// Test listing chunks
		filter := service.ChunkFilter{
			DocumentSubstring: "test",
			ShowContent:       true,
		}

		chunks, err := testRagService.GetRagChunks("test-rag", filter)
		assert.NoError(t, err)
		assert.NotEmpty(t, chunks)

		// Load the RAG from repository
		rag, err := testRagService.LoadRag("test-rag")
		assert.NoError(t, err)

		// Update model name to use llama3.2 for completion
		oldModel := rag.ModelName
		rag.ModelName = completionModel
		err = testRagService.UpdateRag(rag)
		assert.NoError(t, err)

		// Query the RAG
		result, err := testRagService.Query(rag, "What is in the test documents?", 5)
		assert.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.Equal(t, mockResponse, result) // Vérifier que nous avons la réponse attendue

		// Restore model
		rag.ModelName = oldModel
		err = testRagService.UpdateRag(rag)
		assert.NoError(t, err)

		// Print result
		t.Logf("RAG Query Result: %s", result)

		// Verify that all mock expectations were met
		mockOllama.AssertExpectations(t)
		mockEmbeddingService.AssertExpectations(t)
	})
}

func TestRagServiceWithCustomClient(t *testing.T) {
	// Create a test profile for custom OpenAI-compatible endpoint
	profileRepo := repository.NewProfileRepository()
	testProfileName := "test-rag-service-custom-client"
	
	// Cleanup beforehand
	_ = profileRepo.Delete(testProfileName)
	
	// Create a test profile
	profile := &domain.APIProfile{
		Name:      testProfileName,
		Provider:  "openai-api",
		APIKey:    "test-key",
		BaseURL:   "http://localhost:8080/v1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	err := profileRepo.Save(profile)
	assert.NoError(t, err)
	
	defer func() {
		_ = profileRepo.Delete(testProfileName)
	}()
	
	t.Run("NewRagServiceWithClient", func(t *testing.T) {
		// Create a client based on the profile name
		llmClient, err := client.GetLLMClientFromProfile(testProfileName)
		assert.NoError(t, err)
		assert.NotNil(t, llmClient)
		
		// Create RAG service with custom client (providing a default Ollama client)
		ollamaClient := client.NewDefaultOllamaClient()
		ragService := service.NewRagServiceWithClient(llmClient, ollamaClient)
		assert.NotNil(t, ragService)
		
		// Verify the service was created successfully
		assert.NotNil(t, ragService)
	})
	
	t.Run("SetPreferredEmbeddingModel", func(t *testing.T) {
		// Create a client based on the profile name
		llmClient, err := client.GetLLMClientFromProfile(testProfileName)
		assert.NoError(t, err)
		
		// Create RAG service with custom client
		ollamaClient := client.NewDefaultOllamaClient()
		ragService := service.NewRagServiceWithClient(llmClient, ollamaClient)
		
		// Set preferred embedding model
		customEmbeddingModel := "custom-embedding-model"
		ragService.SetPreferredEmbeddingModel(customEmbeddingModel)
		
		// The test passes if no error occurs (SetPreferredEmbeddingModel is void)
		assert.NotNil(t, ragService)
	})
	
	t.Run("NewRagServiceWithClient_MultipleProviders", func(t *testing.T) {
		// Test with traditional OpenAI profile
		openaiProfile := &domain.APIProfile{
			Name:      "test-traditional-openai",
			Provider:  "openai",
			APIKey:    "sk-test-key",
			BaseURL:   "",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		
		err := profileRepo.Save(openaiProfile)
		assert.NoError(t, err)
		defer func() { _ = profileRepo.Delete("test-traditional-openai") }()
		
		// Create client for traditional OpenAI
		openaiClient, err := client.GetLLMClientFromProfile("test-traditional-openai")
		assert.NoError(t, err)
		
		// Create RAG service with OpenAI client
		ollamaClientA := client.NewDefaultOllamaClient()
		ragServiceOpenAI := service.NewRagServiceWithClient(openaiClient, ollamaClientA)
		assert.NotNil(t, ragServiceOpenAI)
		
		// Create client for custom endpoint
		customClient, err := client.GetLLMClientFromProfile(testProfileName)
		assert.NoError(t, err)
		
		// Create RAG service with custom client
		ollamaClientB := client.NewDefaultOllamaClient()
		ragServiceCustom := service.NewRagServiceWithClient(customClient, ollamaClientB)
		assert.NotNil(t, ragServiceCustom)
		
		// Both services should be created successfully but be different instances
		assert.NotEqual(t, ragServiceOpenAI, ragServiceCustom)
	})
}
