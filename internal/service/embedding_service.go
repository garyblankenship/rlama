package service

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/domain"
)

// EmbeddingService manages the generation of embeddings for documents
type EmbeddingService struct {
	llmClient  client.LLMClient
	maxWorkers int // Number of parallel workers for embedding generation
}

// NewEmbeddingService creates a new instance of EmbeddingService
func NewEmbeddingService(llmClient client.LLMClient) *EmbeddingService {
	if llmClient == nil {
		llmClient = client.NewDefaultOllamaClient()
	}
	return &EmbeddingService{
		llmClient:  llmClient,
		maxWorkers: 3, // Default to 3 workers
	}
}

// SetMaxWorkers sets the maximum number of parallel workers for embedding generation
func (es *EmbeddingService) SetMaxWorkers(workers int) {
	if workers < 1 {
		workers = 1
	} else if workers > 5 {
		workers = 5
	}
	es.maxWorkers = workers
}

// GenerateEmbeddings generates embeddings for a list of documents
func (es *EmbeddingService) GenerateEmbeddings(docs []*domain.Document, modelName string) error {
	// Try to use snowflake-arctic-embed2 for embeddings first
	embeddingModel := "snowflake-arctic-embed2"
	
	// Process all documents
	for _, doc := range docs {
		// Generate embedding with snowflake-arctic-embed2 first
		embedding, err := es.llmClient.GenerateEmbedding(embeddingModel, doc.Content)
		
		// If snowflake-arctic-embed2 fails, try to pull it automatically (Ollama only)
		if err != nil {
			fmt.Printf("⚠️ Could not use %s for embeddings: %v\n", embeddingModel, err)
			
			// Attempt to pull the embedding model automatically (only for Ollama clients)
			var pullErr error
			if _, isOllama := es.llmClient.(*client.OllamaClient); isOllama {
				fmt.Printf("Attempting to pull %s automatically...\n", embeddingModel)
				pullErr = es.pullEmbeddingModel(embeddingModel)
			} else {
				pullErr = fmt.Errorf("model pulling not supported for this client type")
			}
			
			if pullErr == nil {
				// Try again with the pulled model
				embedding, err = es.llmClient.GenerateEmbedding(embeddingModel, doc.Content)
			}
			
			// If pulling failed or embedding still fails, fallback to the specified model
			if pullErr != nil || err != nil {
				fmt.Printf("Falling back to %s for embeddings.\n", modelName)
				embedding, err = es.llmClient.GenerateEmbedding(modelName, doc.Content)
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
	embedding, err := es.llmClient.GenerateEmbedding(embeddingModel, query)
	
	// If snowflake-arctic-embed2 fails, try to pull it (but only if not already tried)
	if err != nil {
		// We don't need to show the warning again if already shown in GenerateEmbeddings
		// Attempt to pull the model (this is a no-op if we already tried, and only for Ollama)
		var pullErr error
		if _, isOllama := es.llmClient.(*client.OllamaClient); isOllama {
			pullErr = es.pullEmbeddingModel(embeddingModel)
		} else {
			pullErr = fmt.Errorf("model pulling not supported for this client type")
		}
		
		if pullErr == nil {
			// Try again with the pulled model
			embedding, err = es.llmClient.GenerateEmbedding(embeddingModel, query)
		}
		
		// If pulling failed or embedding still fails, fallback to the specified model
		if pullErr != nil || err != nil {
			embedding, err = es.llmClient.GenerateEmbedding(modelName, query)
			if err != nil {
				return nil, fmt.Errorf("error generating embedding for query: %w", err)
			}
		}
	}

	return embedding, nil
}

// GenerateChunkEmbeddings generates embeddings for document chunks in parallel
func (es *EmbeddingService) GenerateChunkEmbeddings(chunks []*domain.DocumentChunk, modelName string) error {
	// Try to use snowflake-arctic-embed2 for embeddings first
	embeddingModel := "snowflake-arctic-embed2"
	
	// Create a wait group to synchronize goroutines
	var wg sync.WaitGroup
	
	// Create a channel to limit concurrency
	semaphore := make(chan struct{}, es.maxWorkers)
	
	// Create a channel for errors
	errorChan := make(chan error, len(chunks))
	
	// Create a mutex for printing progress
	var progressMutex sync.Mutex
	var completedChunks int
	
	// Check if we need to pull the model (attempt only once)
	modelChecked := false
	var modelCheckMutex sync.Mutex
	
	// Process chunks in parallel
	for i, chunk := range chunks {
		// Add to wait group before starting goroutine
		wg.Add(1)
		
		// Start a goroutine to process this chunk
		go func(index int, ch *domain.DocumentChunk) {
			defer wg.Done()
			
			// Acquire semaphore slot (this limits concurrency)
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			
			// Generate embedding
			embedding, err := es.llmClient.GenerateEmbedding(embeddingModel, ch.Content)
			
			// If the model fails and we haven't checked it yet
			if err != nil {
				modelCheckMutex.Lock()
				shouldCheck := !modelChecked
				modelChecked = true
				modelCheckMutex.Unlock()
				
				if shouldCheck {
					// Only print the warning and attempt to pull once
					fmt.Printf("\n⚠️ Could not use %s for embeddings: %v\n", embeddingModel, err)
					
					var pullErr error
					if _, isOllama := es.llmClient.(*client.OllamaClient); isOllama {
						fmt.Printf("Attempting to pull %s automatically...\n", embeddingModel)
						pullErr = es.pullEmbeddingModel(embeddingModel)
					} else {
						pullErr = fmt.Errorf("model pulling not supported for this client type")
					}
					
					if pullErr == nil {
						// Try again with the pulled model
						embedding, err = es.llmClient.GenerateEmbedding(embeddingModel, ch.Content)
					}
					
					if pullErr != nil || err != nil {
						fmt.Printf("Falling back to %s for embeddings.\n", modelName)
					}
				}
				
				// Use the specified model instead if the embedding model failed
				if err != nil {
					embedding, err = es.llmClient.GenerateEmbedding(modelName, ch.Content)
					if err != nil {
						errorChan <- fmt.Errorf("error generating embedding for chunk %s: %w", ch.ID, err)
						return
					}
				}
			}
			
			// Update the chunk with the embedding
			ch.Embedding = embedding
			
			// Update progress
			progressMutex.Lock()
			completedChunks++
			fmt.Printf("Generating embeddings: %d/%d chunks processed (%d%%)   \r", 
				completedChunks, len(chunks), (completedChunks * 100 / len(chunks)))
			progressMutex.Unlock()
			
		}(i, chunk)
	}
	
	// Wait for all goroutines to complete
	wg.Wait()
	close(errorChan)
	
	// Check if any errors occurred
	for err := range errorChan {
		return err // Return the first error encountered
	}
	
	fmt.Println() // Add a newline after progress indicator
	fmt.Printf("Successfully generated embeddings for %d chunks using %d parallel workers\n", 
		len(chunks), es.maxWorkers)
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