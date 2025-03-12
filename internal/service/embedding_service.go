package service

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/domain"
)

// EmbeddingService manages the generation of embeddings for documents
type EmbeddingService struct {
	ollamaClient *client.OllamaClient
}

// NewEmbeddingService creates a new instance of EmbeddingService
func NewEmbeddingService(ollamaClient *client.OllamaClient) *EmbeddingService {
	if ollamaClient == nil {
		ollamaClient = client.NewDefaultOllamaClient()
	}
	return &EmbeddingService{
		ollamaClient: ollamaClient,
	}
}

// GenerateEmbeddings generates embeddings for a list of documents
func (es *EmbeddingService) GenerateEmbeddings(docs []*domain.Document, modelName string) error {
	// Try to use snowflake-arctic-embed2 for embeddings first
	embeddingModel := "snowflake-arctic-embed2"
	
	// Process all documents
	for _, doc := range docs {
		// Generate embedding with snowflake-arctic-embed2 first
		embedding, err := es.ollamaClient.GenerateEmbedding(embeddingModel, doc.Content)
		
		// If snowflake-arctic-embed2 fails, try to pull it automatically
		if err != nil {
			fmt.Printf("⚠️ Could not use %s for embeddings: %v\n", embeddingModel, err)
			
			// Attempt to pull the embedding model automatically
			fmt.Printf("Attempting to pull %s automatically...\n", embeddingModel)
			pullErr := es.pullEmbeddingModel(embeddingModel)
			
			if pullErr == nil {
				// Try again with the pulled model
				embedding, err = es.ollamaClient.GenerateEmbedding(embeddingModel, doc.Content)
			}
			
			// If pulling failed or embedding still fails, fallback to the specified model
			if pullErr != nil || err != nil {
				fmt.Printf("Falling back to %s for embeddings.\n", modelName)
				embedding, err = es.ollamaClient.GenerateEmbedding(modelName, doc.Content)
				if err != nil {
					return fmt.Errorf("error generating embedding for %s: %w", doc.Path, err)
				}
			}
		}

		doc.Embedding = embedding
	}

	return nil
}

// GenerateQueryEmbedding generates an embedding for a query
func (es *EmbeddingService) GenerateQueryEmbedding(query string, modelName string) ([]float32, error) {
	// Try to use snowflake-arctic-embed2 for embeddings first
	embeddingModel := "snowflake-arctic-embed2"
	
	// Generate embedding with snowflake-arctic-embed2
	embedding, err := es.ollamaClient.GenerateEmbedding(embeddingModel, query)
	
	// If snowflake-arctic-embed2 fails, try to pull it (but only if not already tried)
	if err != nil {
		// We don't need to show the warning again if already shown in GenerateEmbeddings
		// Attempt to pull the model (this is a no-op if we already tried)
		pullErr := es.pullEmbeddingModel(embeddingModel)
		
		if pullErr == nil {
			// Try again with the pulled model
			embedding, err = es.ollamaClient.GenerateEmbedding(embeddingModel, query)
		}
		
		// If pulling failed or embedding still fails, fallback to the specified model
		if pullErr != nil || err != nil {
			embedding, err = es.ollamaClient.GenerateEmbedding(modelName, query)
			if err != nil {
				return nil, fmt.Errorf("error generating embedding for query: %w", err)
			}
		}
	}

	return embedding, nil
}

// GenerateChunkEmbeddings generates embeddings for document chunks
func (es *EmbeddingService) GenerateChunkEmbeddings(chunks []*domain.DocumentChunk, modelName string) error {
	// Try to use snowflake-arctic-embed2 for embeddings first
	embeddingModel := "snowflake-arctic-embed2"
	
	// Process all chunks
	for i, chunk := range chunks {
		fmt.Printf("Generating embedding for chunk %d/%d\r", i+1, len(chunks))
		
		// Generate embedding with snowflake-arctic-embed2 first
		embedding, err := es.ollamaClient.GenerateEmbedding(embeddingModel, chunk.Content)
		
		// If snowflake-arctic-embed2 fails and this is the first chunk, try to pull it
		if err != nil {
			if i == 0 {
				fmt.Printf("\n⚠️ Could not use %s for embeddings: %v\n", embeddingModel, err)
				
				// Attempt to pull the embedding model automatically
				fmt.Printf("Attempting to pull %s automatically...\n", embeddingModel)
				pullErr := es.pullEmbeddingModel(embeddingModel)
				
				if pullErr == nil {
					// Try again with the pulled model
					embedding, err = es.ollamaClient.GenerateEmbedding(embeddingModel, chunk.Content)
				}
				
				// If pulling failed or embedding still fails, fallback to the specified model
				if pullErr != nil || err != nil {
					fmt.Printf("Falling back to %s for embeddings.\n", modelName)
				}
			}
			
			// Use the specified model instead if the embedding model failed
			if err != nil {
				embedding, err = es.ollamaClient.GenerateEmbedding(modelName, chunk.Content)
				if err != nil {
					return fmt.Errorf("error generating embedding for chunk %s: %w", chunk.ID, err)
				}
			}
		}

		chunk.Embedding = embedding
	}
	fmt.Println() // Add a newline after progress indicator
	return nil
}

// Track if we've already tried to pull the model to avoid multiple attempts
var attemptedModelPull = make(map[string]bool)

// pullEmbeddingModel attempts to pull the embedding model via Ollama
func (es *EmbeddingService) pullEmbeddingModel(modelName string) error {
	// Check if we've already tried to pull this model
	if attemptedModelPull[modelName] {
		return fmt.Errorf("already attempted to pull model")
	}
	
	// Mark that we've attempted to pull this model
	attemptedModelPull[modelName] = true
	
	// Check if Ollama CLI is available
	cmd := exec.Command("ollama", "list")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ollama command not available: %w", err)
	}
	
	fmt.Printf("Pulling %s model (this may take a while)...\n", modelName)
	
	// Run the ollama pull command
	cmd = exec.Command("ollama", "pull", modelName)
	cmd.Stdout = os.Stdout // Display output to the user
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
} 