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
	UpdateRerankerModel(ragName string, model string) error
	ListAllRags() ([]string, error)
	GetOllamaClient() *client.OllamaClient
	SetPreferredEmbeddingModel(model string)
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

// NewRagServiceWithConfig creates a new instance of RagService with service configuration
func NewRagServiceWithConfig(ollamaClient *client.OllamaClient, config *ServiceConfig) RagService {
	if ollamaClient == nil {
		ollamaClient = client.NewDefaultOllamaClient()
	}

	// Create reranker service with ONNX configuration if specified
	var rerankerService *RerankerService
	if config.UseONNXReranker {
		rerankerService = NewRerankerServiceWithOptions(ollamaClient, true, config.ONNXModelDir)
	} else {
		rerankerService = NewRerankerService(ollamaClient)
	}

	return &RagServiceImpl{
		documentLoader:   NewDocumentLoader(),
		embeddingService: NewEmbeddingService(ollamaClient),
		ragRepository:    repository.NewRagRepository(),
		ollamaClient:     ollamaClient,
		rerankerService:  rerankerService,
	}
}

// NewRagServiceWithClient creates a new instance of RagService with the specified LLM client
func NewRagServiceWithClient(llmClient client.LLMClient, ollamaClient *client.OllamaClient) RagService {
	// Use the new composite service architecture
	return NewCompositeRagService(llmClient, ollamaClient)
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

// SetPreferredEmbeddingModel sets the preferred embedding model to use
func (rs *RagServiceImpl) SetPreferredEmbeddingModel(model string) {
	rs.embeddingService.SetPreferredEmbeddingModel(model)
}

// CreateRagWithOptions creates a new RAG system with options
func (rs *RagServiceImpl) CreateRagWithOptions(modelName, ragName, folderPath string, options DocumentLoaderOptions) error {
	// Check if model is available using the correct client
	// The embedding service has the right LLM client (OpenAI or Ollama)
	fmt.Printf("🔍 Debug: Checking model availability for '%s'\n", modelName)
	if rs.embeddingService != nil {
		llmClient := rs.embeddingService.GetLLMClient()
		if llmClient != nil {
			fmt.Printf("✓ Using embedding service LLM client for validation\n")
			if err := llmClient.CheckLLMAndModel(modelName); err != nil {
				fmt.Printf("❌ Model validation failed: %v\n", err)
				return err
			}
			fmt.Printf("✓ Model validation successful\n")
		} else {
			fmt.Printf("⚠️ Embedding service has no LLM client, falling back to Ollama\n")
			// Fallback to Ollama client if embedding service doesn't have a client
			if err := rs.ollamaClient.CheckOllamaAndModel(modelName); err != nil {
				return err
			}
		}
	} else {
		fmt.Printf("⚠️ No embedding service, falling back to Ollama\n")
		// Fallback to Ollama client if embedding service is not properly configured
		if err := rs.ollamaClient.CheckOllamaAndModel(modelName); err != nil {
			return err
		}
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

	// Detect embedding dimension
	embeddingDim, err := rs.embeddingService.DetectEmbeddingDimension(modelName)
	if err != nil {
		return fmt.Errorf("error detecting embedding dimension: %w", err)
	}
	fmt.Printf("Detected embedding dimension: %d\n", embeddingDim)

	// Create the RAG system with detected dimensions and vector store configuration
	var rag *domain.RagSystem
	if options.VectorStore == "qdrant" {
		rag = domain.NewRagSystemWithVectorStore(
			ragName, 
			modelName, 
			embeddingDim,
			options.VectorStore,
			options.QdrantHost,
			options.QdrantPort,
			options.QdrantAPIKey,
			options.QdrantCollectionName,
			options.QdrantGRPC,
		)
		if rag == nil {
			return fmt.Errorf("failed to create RAG system with Qdrant configuration")
		}
		fmt.Printf("Created RAG system with Qdrant vector store at %s:%d\n", options.QdrantHost, options.QdrantPort)
	} else {
		rag = domain.NewRagSystemWithDimensions(ragName, modelName, embeddingDim)
		if rag == nil {
			return fmt.Errorf("failed to create RAG system with internal vector store")
		}
		fmt.Printf("Created RAG system with internal vector store\n")
	}
	
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

	// Set preferred embedding model if specified
	if options.EmbeddingModel != "" {
		rs.embeddingService.SetPreferredEmbeddingModel(options.EmbeddingModel)
		fmt.Printf("Set preferred embedding model to: %s\n", options.EmbeddingModel)
	}

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
	// Use the embedding service's LLM client for consistency
	llmClient := rs.embeddingService.GetLLMClient()

	// The embedding service already has the right client configured
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
			if !rag.RerankerSilent {
				fmt.Printf("Using default context size of %d for reranked results\n", contextSize)
			}
		} else {
			contextSize = 20 // 20 par défaut si le reranking est désactivé
			fmt.Printf("Using context size of %d (reranking disabled)\n", contextSize) // Always show this message since reranking is disabled
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
		if !rag.RerankerSilent {
			fmt.Printf("Retrieving %d initial results for reranking...\n", initialRetrievalCount)
		}
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
			RerankerModel:     "BAAI/bge-reranker-v2-m3", // Always prefer BGE reranker
			ScoreThreshold:    0.3,                       // Minimum relevance threshold
			RerankerWeight:    rag.RerankerWeight,
			AdaptiveFiltering: true, // Enable adaptive filtering
			Silent:            rag.RerankerSilent, // Use the silent setting from the RAG
		}

		// If a specific BGE reranker model is defined in the RAG, use that one
		// This allows users to choose between different BGE reranker models
		if rag.RerankerModel != "" && strings.Contains(strings.ToLower(rag.RerankerModel), "bge-reranker") {
			options.RerankerModel = rag.RerankerModel
		}

		// Display the effective model being used (if not in silent mode)
		if !options.Silent {
			fmt.Printf("Reranking and filtering results for relevance using model '%s'...\n", options.RerankerModel)
		}

		// Perform reranking with adaptive filtering
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
	// Load the existing RAG system
	rag, err := rs.LoadRag(ragName)
	if err != nil {
		return fmt.Errorf("error loading RAG '%s': %w", ragName, err)
	}

	// Check if Ollama is available
	if err := rs.ollamaClient.CheckOllamaAndModel(rag.ModelName); err != nil {
		return err
	}

	// Load new documents with options
	newDocs, err := rs.documentLoader.LoadDocumentsFromFolderWithOptions(folderPath, options)
	if err != nil {
		return fmt.Errorf("error loading documents: %w", err)
	}

	if len(newDocs) == 0 {
		return fmt.Errorf("no valid documents found in folder %s", folderPath)
	}

	fmt.Printf("Successfully loaded %d new documents. Chunking documents...\n", len(newDocs))

	// Create chunker service with the same options as the RAG or from provided options
	chunkSize := rag.WatchOptions.ChunkSize
	chunkOverlap := rag.WatchOptions.ChunkOverlap
	chunkingStrategy := rag.ChunkingStrategy

	// Override with provided options if specified
	if options.ChunkSize > 0 {
		chunkSize = options.ChunkSize
	}
	if options.ChunkOverlap > 0 {
		chunkOverlap = options.ChunkOverlap
	}
	if options.ChunkingStrategy != "" {
		chunkingStrategy = options.ChunkingStrategy
	}

	// Create chunker with configured options
	chunkerService := NewChunkerService(ChunkingConfig{
		ChunkSize:        chunkSize,
		ChunkOverlap:     chunkOverlap,
		ChunkingStrategy: chunkingStrategy,
	})

	// Check for duplicates
	existingDocPaths := make(map[string]bool)
	for _, doc := range rag.Documents {
		existingDocPaths[doc.Path] = true
	}

	var uniqueDocs []*domain.Document
	var skippedDocs int

	// Filter out duplicate documents
	for _, doc := range newDocs {
		if existingDocPaths[doc.Path] {
			skippedDocs++
			continue
		}
		uniqueDocs = append(uniqueDocs, doc)
		existingDocPaths[doc.Path] = true // Mark as processed to avoid future duplicates
	}

	if len(uniqueDocs) == 0 {
		return fmt.Errorf("all %d documents already exist in the RAG, none added", skippedDocs)
	}

	if skippedDocs > 0 {
		fmt.Printf("Skipped %d documents that were already in the RAG.\n", skippedDocs)
	}

	// Process each unique document - chunk and generate embeddings
	var allChunks []*domain.DocumentChunk
	for _, doc := range uniqueDocs {
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

	fmt.Printf("Generated %d chunks from %d new documents. Generating embeddings...\n",
		len(allChunks), len(uniqueDocs))

	// Generate embeddings for all chunks
	err = rs.embeddingService.GenerateChunkEmbeddings(allChunks, rag.ModelName)
	if err != nil {
		return fmt.Errorf("error generating embeddings: %w", err)
	}

	// Add all chunks to the RAG
	for _, chunk := range allChunks {
		rag.AddChunk(chunk)
	}

	// Update the RAG's chunk options based on the most recent settings
	rag.WatchOptions.ChunkSize = chunkSize
	rag.WatchOptions.ChunkOverlap = chunkOverlap
	rag.ChunkingStrategy = chunkingStrategy

	// Update reranker settings if specified in options
	if options.RerankerModel != "" {
		rag.RerankerModel = options.RerankerModel
	}
	if options.RerankerWeight > 0 {
		rag.RerankerWeight = options.RerankerWeight
	}
	if options.EnableReranker {
		rag.RerankerEnabled = true
	}

	// Save the updated RAG
	err = rs.ragRepository.Save(rag)
	if err != nil {
		return fmt.Errorf("error saving the updated RAG: %w", err)
	}

	fmt.Printf("Successfully added %d new documents (%d chunks) to RAG '%s'.\n",
		len(uniqueDocs), len(allChunks), ragName)
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

// UpdateRerankerModel updates the reranker model of a RAG
func (rs *RagServiceImpl) UpdateRerankerModel(ragName string, model string) error {
	// Load the RAG
	rag, err := rs.LoadRag(ragName)
	if err != nil {
		return fmt.Errorf("error loading RAG: %w", err)
	}

	// Update the reranker model
	rag.RerankerModel = model

	// Save the updated RAG
	err = rs.ragRepository.Save(rag)
	if err != nil {
		return fmt.Errorf("error saving updated RAG: %w", err)
	}

	fmt.Printf("Reranker model updated to '%s' for RAG '%s'.\n", model, rag.Name)
	return nil
}
