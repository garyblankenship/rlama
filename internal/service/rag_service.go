package service

import (
	"fmt"
	"strings"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/repository"
)

// RagService interface defines the contract for RAG operations
type RagService interface {
	CreateRagWithOptions(modelName, ragName, folderPath string, options DocumentLoaderOptions) error
	GetRagChunks(ragName string, filter ChunkFilter) ([]*domain.DocumentChunk, error)
	LoadRag(ragName string) (*domain.RagSystem, error)
	Query(rag *domain.RagSystem, query string, contextSize int) (string, error)
	AddDocsWithOptions(ragName string, folderPath string, options DocumentLoaderOptions) error
	UpdateModel(ragName string, newModel string) error
	UpdateRag(rag *domain.RagSystem) error
	ListAllRags() ([]string, error)
	GetOllamaClient() *client.OllamaClient
	// Directory watching methods
	SetupDirectoryWatching(ragName string, dirPath string, watchInterval int, options DocumentLoaderOptions) error
	DisableDirectoryWatching(ragName string) error
	CheckWatchedDirectory(ragName string) (int, error)
	// Web watching methods
	SetupWebWatching(ragName string, websiteURL string, watchInterval int, options domain.WebWatchOptions) error
	DisableWebWatching(ragName string) error
	CheckWatchedWebsite(ragName string) (int, error)
}

// ChunkFilter defines filtering criteria for retrieving chunks
type ChunkFilter struct {
	DocumentSubstring string
	ShowContent       bool
}

// RagServiceImpl implements the RagService interface
type RagServiceImpl struct {
	documentLoader   *DocumentLoader
	embeddingService *EmbeddingService
	ragRepository    *repository.RagRepository
	ollamaClient     *client.OllamaClient
	rerankerService  *RerankerService
}

// NewRagService creates a new instance of RagService
func NewRagService(ollamaClient *client.OllamaClient) RagService {
	if ollamaClient == nil {
		ollamaClient = client.NewDefaultOllamaClient()
	}

	return &RagServiceImpl{
		documentLoader:   NewDocumentLoader(),
		embeddingService: NewEmbeddingService(ollamaClient),
		ragRepository:    repository.NewRagRepository(),
		ollamaClient:     ollamaClient,
		rerankerService:  NewRerankerService(ollamaClient),
	}
}

// NewRagServiceWithEmbedding creates a new RagService with a specific embedding service
func NewRagServiceWithEmbedding(ollamaClient *client.OllamaClient, embeddingService *EmbeddingService) RagService {
	if ollamaClient == nil {
		ollamaClient = client.NewDefaultOllamaClient()
	}

	return &RagServiceImpl{
		documentLoader:   NewDocumentLoader(),
		embeddingService: embeddingService,
		ragRepository:    repository.NewRagRepository(),
		ollamaClient:     ollamaClient,
		rerankerService:  NewRerankerService(ollamaClient),
	}
}

// GetOllamaClient returns the Ollama client
func (rs *RagServiceImpl) GetOllamaClient() *client.OllamaClient {
	return rs.ollamaClient
}

// CreateRagWithOptions creates a new RAG system with options
func (rs *RagServiceImpl) CreateRagWithOptions(modelName, ragName, folderPath string, options DocumentLoaderOptions) error {
	// Check if Ollama is available
	if err := rs.ollamaClient.CheckOllamaAndModel(modelName); err != nil {
		return err
	}

	// Check if the RAG already exists
	if rs.ragRepository.Exists(ragName) {
		return fmt.Errorf("a RAG with name '%s' already exists", ragName)
	}

	// Load documents with options
	docs, err := rs.documentLoader.LoadDocumentsFromFolderWithOptions(folderPath, options)
	if err != nil {
		return fmt.Errorf("error loading documents: %w", err)
	}

	if len(docs) == 0 {
		return fmt.Errorf("no valid documents found in folder %s", folderPath)
	}

	fmt.Printf("Successfully loaded %d documents. Chunking documents...\n", len(docs))

	// Create the RAG system
	rag := domain.NewRagSystem(ragName, modelName)
	rag.ChunkingStrategy = options.ChunkingStrategy
	rag.APIProfileName = options.APIProfileName

	// Configure reranking options - enable by default
	rag.RerankerEnabled = true // Always enable reranking by default
	fmt.Println("Reranking enabled for better retrieval accuracy")

	// Only disable if explicitly set to false in options
	if !options.EnableReranker && options.RerankerModel == "" {
		// Check if EnableReranker field was explicitly set
		// This prevents the zero-value (false) from disabling reranking when the field isn't set
		rag.RerankerEnabled = false
		fmt.Println("Reranking disabled by user configuration")
	}

	// Set reranker model if specified, otherwise use the same model
	if options.RerankerModel != "" {
		rag.RerankerModel = options.RerankerModel
	} else {
		rag.RerankerModel = modelName
	}

	// Set reranker weight
	if options.RerankerWeight > 0 {
		rag.RerankerWeight = options.RerankerWeight
	} else {
		rag.RerankerWeight = 0.7 // Default to 70% reranker, 30% vector
	}

	// Set default TopK if not already set
	if rag.RerankerTopK <= 0 {
		rag.RerankerTopK = 5 // Default to 5 results
	}

	// Set chunking options in WatchOptions too
	rag.WatchOptions.ChunkSize = options.ChunkSize
	rag.WatchOptions.ChunkOverlap = options.ChunkOverlap
	rag.WatchOptions.ChunkingStrategy = options.ChunkingStrategy

	// Create chunker service
	chunkerService := NewChunkerService(ChunkingConfig{
		ChunkSize:        options.ChunkSize,
		ChunkOverlap:     options.ChunkOverlap,
		ChunkingStrategy: options.ChunkingStrategy,
	})

	// Process each document - chunk and generate embeddings
	var allChunks []*domain.DocumentChunk
	for _, doc := range docs {
		// Add the document to the RAG
		rag.AddDocument(doc)

		// Chunk the document
		chunks := chunkerService.ChunkDocument(doc)

		// Update total chunks in metadata
		for i, chunk := range chunks {
			chunk.ChunkNumber = i
			chunk.TotalChunks = len(chunks)
		}

		allChunks = append(allChunks, chunks...)
	}

	fmt.Printf("Generated %d chunks from %d documents. Generating embeddings...\n",
		len(allChunks), len(docs))

	// Generate embeddings for all chunks
	err = rs.embeddingService.GenerateChunkEmbeddings(allChunks, modelName)
	if err != nil {
		return fmt.Errorf("error generating embeddings: %w", err)
	}

	// Add all chunks to the RAG
	for _, chunk := range allChunks {
		rag.AddChunk(chunk)
	}

	// Save the RAG
	err = rs.ragRepository.Save(rag)
	if err != nil {
		return fmt.Errorf("error saving the RAG: %w", err)
	}

	fmt.Printf("RAG created with %d indexed documents (%d chunks).\n", len(docs), len(allChunks))
	return nil
}

// GetRagChunks gets chunks from a RAG with filtering
func (rs *RagServiceImpl) GetRagChunks(ragName string, filter ChunkFilter) ([]*domain.DocumentChunk, error) {
	// Load the RAG
	rag, err := rs.LoadRag(ragName)
	if err != nil {
		return nil, fmt.Errorf("error loading RAG: %w", err)
	}

	var filteredChunks []*domain.DocumentChunk

	// Apply filters
	for _, chunk := range rag.Chunks {
		// Apply document name filter if provided
		if filter.DocumentSubstring != "" {
			docID := chunk.DocumentID
			doc := rag.GetDocumentByID(docID)
			if doc != nil && !strings.Contains(doc.Name, filter.DocumentSubstring) {
				continue
			}
		}

		filteredChunks = append(filteredChunks, chunk)
	}

	return filteredChunks, nil
}

// LoadRag loads a RAG system
func (rs *RagServiceImpl) LoadRag(ragName string) (*domain.RagSystem, error) {
	return rs.ragRepository.Load(ragName)
}

// Query performs a query on a RAG system
func (rs *RagServiceImpl) Query(rag *domain.RagSystem, query string, contextSize int) (string, error) {
	// Check if Ollama is available
	var llmClient client.LLMClient

	// Determine which client to use based on the model
	if client.IsOpenAIModel(rag.ModelName) {
		// For OpenAI, use the specified profile or default
		openAIClient, err := client.NewOpenAIClientWithProfile(rag.APIProfileName)
		if err != nil {
			return "", err
		}
		llmClient = openAIClient
	} else {
		llmClient = rs.ollamaClient
	}

	if err := llmClient.CheckLLMAndModel(rag.ModelName); err != nil {
		return "", err
	}

	// Generate embedding for the query
	queryEmbedding, err := rs.embeddingService.GenerateQueryEmbedding(query, rag.ModelName)
	if err != nil {
		return "", fmt.Errorf("error generating embedding for query: %w", err)
	}

	// Use the provided context size or default value based on settings
	rerankerOpts := DefaultRerankerOptions()

	// Si contextSize est 0 (auto), utiliser:
	// - RerankerTopK du RAG si défini
	// - Sinon le TopK par défaut (5)
	// - 20 si le reranking est désactivé
	if contextSize <= 0 {
		if rag.RerankerEnabled {
			if rag.RerankerTopK > 0 {
				contextSize = rag.RerankerTopK
			} else {
				contextSize = rerankerOpts.TopK // 5 par défaut
			}
			fmt.Printf("Using default context size of %d for reranked results\n", contextSize)
		} else {
			contextSize = 20 // 20 par défaut si le reranking est désactivé
			fmt.Printf("Using context size of %d (reranking disabled)\n", contextSize)
		}
	}

	// First-stage retrieval: Get initial results using vector search
	// Get more results than needed for reranking
	initialRetrievalCount := contextSize
	if rag.RerankerEnabled {
		// If reranking is enabled, retrieve more documents initially (20 or 2*contextSize, whichever is larger)
		initialRetrievalCount = rerankerOpts.InitialK
		if initialRetrievalCount < contextSize {
			initialRetrievalCount = contextSize * 2 // Ensure we get enough documents for reranking
		}
		fmt.Printf("Retrieving %d initial results for reranking...\n", initialRetrievalCount)
	}

	// Search for the most relevant chunks
	results := rag.HybridStore.Search(queryEmbedding, initialRetrievalCount)

	// Second-stage retrieval: Re-rank if enabled
	var rankedResults []RankedResult
	var includedDocs = make(map[string]bool)

	if rag.RerankerEnabled {
		// Set reranker options for adaptive content-based filtering
		options := RerankerOptions{
			// Don't limit by fixed TopK but use minimum threshold
			TopK:              100, // Set to a high value to avoid arbitrary limit
			InitialK:          initialRetrievalCount,
			RerankerModel:     rag.RerankerModel,
			ScoreThreshold:    0.3, // Minimum relevance threshold
			RerankerWeight:    rag.RerankerWeight,
			AdaptiveFiltering: true, // Enable adaptive filtering
		}

		// If no reranker model specified, use the same as the main model
		if options.RerankerModel == "" {
			options.RerankerModel = rag.ModelName
		}

		// Perform reranking with adaptive filtering
		fmt.Printf("Reranking and filtering results for relevance using model '%s'...\n", options.RerankerModel)
		rerankedResults, err := rs.rerankerService.Rerank(query, rag, results, options)
		if err != nil {
			return "", fmt.Errorf("error during reranking: %w", err)
		}

		rankedResults = rerankedResults

		// Track documents included after adaptive filtering
		for _, result := range rankedResults {
			includedDocs[result.Chunk.DocumentID] = true
		}

		// Show information about filtered results
		fmt.Printf("Selected %d relevant chunks from %d initial results\n",
			len(rankedResults), len(results))
	}

	// Build the context
	var context strings.Builder
	context.WriteString("Relevant information:\n\n")

	// Use the reranked results if available, otherwise use the initial results
	if rag.RerankerEnabled && len(rankedResults) > 0 {
		for _, result := range rankedResults {
			chunk := result.Chunk
			// Add chunk content with its metadata
			context.WriteString(fmt.Sprintf("--- %s (Score: %.4f) ---\n%s\n\n",
				chunk.GetMetadataString(), result.FinalScore, chunk.Content))
		}
	} else {
		// Use original vector search results if reranking is disabled or failed
		for _, result := range results {
			chunk := rag.GetChunkByID(result.ID)
			if chunk != nil {
				// Add chunk content with its metadata
				context.WriteString(fmt.Sprintf("--- %s ---\n%s\n\n",
					chunk.GetMetadataString(), chunk.Content))

				includedDocs[chunk.DocumentID] = true
			}
		}
	}

	// Build the prompt with better formatting and instructions for citing sources
	systemMessage := "You are a helpful assistant that provides accurate information based on the documents you've been given. Answer the question based on the context provided. If you don't know the answer based on the context, say that you don't know rather than making up an answer. Important: Always respond in the same language as the user's query."
	prompt := fmt.Sprintf(`System: %s

Context:
%s

Question: %s

Answer based on the provided information. If the information doesn't contain the answer, say so clearly.
Include references to the source documents in your answer using the format (Source: document name).`,
		systemMessage, context.String(), query)

	// Show search results to the user
	fmt.Println()
	fmt.Printf("Found %d relevant sections across %d documents\n",
		len(includedDocs), len(includedDocs))

	// Generate the response with the appropriate client
	response, err := llmClient.GenerateCompletion(rag.ModelName, prompt)
	if err != nil {
		return "", fmt.Errorf("error generating response: %w", err)
	}

	return response, nil
}

// AddDocsWithOptions adds documents to a RAG with options
func (rs *RagServiceImpl) AddDocsWithOptions(ragName string, folderPath string, options DocumentLoaderOptions) error {
	return nil
}

// UpdateModel updates the model of a RAG
func (rs *RagServiceImpl) UpdateModel(ragName string, newModel string) error {
	return nil
}

// UpdateRag updates a RAG system
func (rs *RagServiceImpl) UpdateRag(rag *domain.RagSystem) error {
	// Save the updated RAG
	err := rs.ragRepository.Save(rag)
	if err != nil {
		return fmt.Errorf("error saving updated RAG: %w", err)
	}

	fmt.Printf("RAG '%s' updated successfully.\n", rag.Name)
	return nil
}

// ListAllRags lists all available RAGs
func (rs *RagServiceImpl) ListAllRags() ([]string, error) {
	return nil, nil
}

// SetupDirectoryWatching sets up directory watching for a RAG
func (rs *RagServiceImpl) SetupDirectoryWatching(ragName string, dirPath string, watchInterval int, options DocumentLoaderOptions) error {
	return nil
}

// DisableDirectoryWatching disables directory watching for a RAG
func (rs *RagServiceImpl) DisableDirectoryWatching(ragName string) error {
	return nil
}

// CheckWatchedDirectory checks a watched directory for changes
func (rs *RagServiceImpl) CheckWatchedDirectory(ragName string) (int, error) {
	return 0, nil
}

// SetupWebWatching sets up web watching for a RAG
func (rs *RagServiceImpl) SetupWebWatching(ragName string, websiteURL string, watchInterval int, options domain.WebWatchOptions) error {
	return nil
}

// DisableWebWatching disables web watching for a RAG
func (rs *RagServiceImpl) DisableWebWatching(ragName string) error {
	return nil
}

// CheckWatchedWebsite checks a watched website for changes
func (rs *RagServiceImpl) CheckWatchedWebsite(ragName string) (int, error) {
	return 0, nil
}
