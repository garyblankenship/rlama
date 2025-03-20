package tests

import (
	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/repository"
	"github.com/dontizi/rlama/internal/service"
	"github.com/stretchr/testify/mock"
)

// MockOllamaClient simule un client Ollama
type MockOllamaClient struct {
	mock.Mock
}

func (m *MockOllamaClient) GenerateCompletion(model, prompt string) (string, error) {
	args := m.Called(model, prompt)
	return args.String(0), args.Error(1)
}

func (m *MockOllamaClient) GenerateEmbedding(model, text string) ([]float32, error) {
	args := m.Called(model, text)
	return args.Get(0).([]float32), args.Error(1)
}

func (m *MockOllamaClient) CheckOllamaAndModel(modelName string) error {
	args := m.Called(modelName)
	return args.Error(0)
}

func (m *MockOllamaClient) CheckLLMAndModel(modelName string) error {
	return m.CheckOllamaAndModel(modelName)
}

// MockEmbeddingService simule le service d'embedding
type MockEmbeddingService struct {
	mock.Mock
}

func (m *MockEmbeddingService) GenerateChunkEmbeddings(chunks []*domain.DocumentChunk, modelName string) error {
	args := m.Called(chunks, modelName)
	return args.Error(0)
}

func (m *MockEmbeddingService) GenerateQueryEmbedding(query string, modelName string) ([]float32, error) {
	args := m.Called(query, modelName)
	return args.Get(0).([]float32), args.Error(1)
}

func (m *MockEmbeddingService) GetOllamaClient() *client.OllamaClient {
	return nil
}

// TestRagService est une implémentation de test du service RAG
type TestRagService struct {
	mockOllama    *MockOllamaClient
	mockEmbedding *MockEmbeddingService
	ragRepository *repository.RagRepository
}

func (s *TestRagService) CreateRagWithOptions(modelName, ragName, folderPath string, options service.DocumentLoaderOptions) error {
	if err := s.mockOllama.CheckOllamaAndModel(modelName); err != nil {
		return err
	}

	rag := domain.NewRagSystem(ragName, modelName)

	// Simuler l'ajout d'un document et la génération d'embeddings
	chunk := &domain.DocumentChunk{
		ID:      "test-chunk",
		Content: "Test content",
	}

	// Appeler GenerateEmbedding une seule fois
	embedding, err := s.mockOllama.GenerateEmbedding(modelName, chunk.Content)
	if err != nil {
		return err
	}
	chunk.Embedding = embedding

	// Appeler GenerateChunkEmbeddings une seule fois
	err = s.mockEmbedding.GenerateChunkEmbeddings([]*domain.DocumentChunk{chunk}, modelName)
	if err != nil {
		return err
	}

	// Appeler GenerateQueryEmbedding une seule fois
	_, err = s.mockEmbedding.GenerateQueryEmbedding("test query", modelName)
	if err != nil {
		return err
	}

	rag.AddChunk(chunk)
	return s.ragRepository.Save(rag)
}

func (s *TestRagService) LoadRag(ragName string) (*domain.RagSystem, error) {
	return s.ragRepository.Load(ragName)
}

func (s *TestRagService) DeleteRag(ragName string) error {
	return s.ragRepository.Delete(ragName)
}

func (s *TestRagService) ListRags() ([]string, error) {
	return []string{"test-rag"}, nil
}
