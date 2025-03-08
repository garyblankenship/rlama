package service

import (
	"fmt"
	"strings"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/repository"
)

// RagService manages operations related to RAG systems
type RagService struct {
	documentLoader   *DocumentLoader
	embeddingService *EmbeddingService
	ragRepository    *repository.RagRepository
	ollamaClient     *client.OllamaClient
}

// NewRagService creates a new instance of RagService
func NewRagService(ollamaClient *client.OllamaClient) *RagService {
	if ollamaClient == nil {
		ollamaClient = client.NewDefaultOllamaClient()
	}
	
	return &RagService{
		documentLoader:   NewDocumentLoader(),
		embeddingService: NewEmbeddingService(ollamaClient),
		ragRepository:    repository.NewRagRepository(),
		ollamaClient:     ollamaClient,
	}
}

// CreateRag creates a new RAG system
func (rs *RagService) CreateRag(modelName, ragName, folderPath string) error {
	// Check if Ollama is available
	if err := rs.ollamaClient.CheckOllamaAndModel(modelName); err != nil {
		return err
	}

	// Check if the RAG already exists
	if rs.ragRepository.Exists(ragName) {
		return fmt.Errorf("a RAG with name '%s' already exists", ragName)
	}

	// Load documents
	docs, err := rs.documentLoader.LoadDocumentsFromFolder(folderPath)
	if err != nil {
		return fmt.Errorf("error loading documents: %w", err)
	}

	if len(docs) == 0 {
		return fmt.Errorf("no valid documents found in folder %s", folderPath)
	}

	fmt.Printf("Successfully loaded %d documents. Generating embeddings...\n", len(docs))
	
	// Create the RAG system
	rag := domain.NewRagSystem(ragName, modelName)

	// Generate embeddings for all documents
	err = rs.embeddingService.GenerateEmbeddings(docs, modelName)
	if err != nil {
		return fmt.Errorf("error generating embeddings: %w", err)
	}

	// Add documents to the RAG
	for _, doc := range docs {
		rag.AddDocument(doc)
	}

	// Save the RAG
	err = rs.ragRepository.Save(rag)
	if err != nil {
		return fmt.Errorf("error saving the RAG: %w", err)
	}

	fmt.Printf("RAG created with %d indexed documents.\n", len(docs))
	return nil
}

// LoadRag loads a RAG system
func (rs *RagService) LoadRag(ragName string) (*domain.RagSystem, error) {
	rag, err := rs.ragRepository.Load(ragName)
	if err != nil {
		return nil, fmt.Errorf("error loading RAG '%s': %w", ragName, err)
	}

	return rag, nil
}

// Query performs a query on a RAG system
func (rs *RagService) Query(rag *domain.RagSystem, query string) (string, error) {
	// Check if Ollama is available
	if err := rs.ollamaClient.CheckOllamaAndModel(rag.ModelName); err != nil {
		return "", err
	}

	// Generate embedding for the query
	queryEmbedding, err := rs.embeddingService.GenerateQueryEmbedding(query, rag.ModelName)
	if err != nil {
		return "", fmt.Errorf("error generating embedding for query: %w", err)
	}

	// Search for the most relevant documents
	results := rag.VectorStore.Search(queryEmbedding, 3) // Top 3 documents

	// Build the context
	var context strings.Builder
	context.WriteString("Relevant information:\n\n")

	for _, result := range results {
		doc := rag.GetDocumentByID(result.ID)
		if doc != nil {
			// Limit content size to avoid prompts that are too long
			content := doc.Content
			if len(content) > 1000 {
				content = content[:1000] + "..."
			}
			context.WriteString(fmt.Sprintf("--- Document: %s ---\n%s\n\n", doc.Name, content))
		}
	}

	// Build the prompt
	prompt := fmt.Sprintf(`You are a helpful AI assistant. Use the information below to answer the question.

%s

Question: %s

Answer concisely based only on the information provided above:`, context.String(), query)

	// Generate the response
	response, err := rs.ollamaClient.GenerateCompletion(rag.ModelName, prompt)
	if err != nil {
		return "", fmt.Errorf("error generating response: %w", err)
	}

	return response, nil
}

// UpdateRag updates an existing RAG system
func (rs *RagService) UpdateRag(rag *domain.RagSystem) error {
	err := rs.ragRepository.Save(rag)
	if err != nil {
		return fmt.Errorf("error updating the RAG: %w", err)
	}
	return nil
} 