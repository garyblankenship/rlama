package service

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/repository"
)

// Remove duplicate interface declaration and keep this one
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
	// Add new methods for watching
	SetupDirectoryWatching(ragName string, dirPath string, watchInterval int, options DocumentLoaderOptions) error
	DisableDirectoryWatching(ragName string) error
	CheckWatchedDirectory(ragName string) (int, error)
	// Add any other required methods here
	// Web watching methods (new)
	SetupWebWatching(ragName string, websiteURL string, watchInterval int, options domain.WebWatchOptions) error
	DisableWebWatching(ragName string) error
	CheckWatchedWebsite(ragName string) (int, error)
}

// Update the struct implementation to match the interface
type RagServiceImpl struct {
	documentLoader   *DocumentLoader
	embeddingService *EmbeddingService
	ragRepository    *repository.RagRepository
	ollamaClient     *client.OllamaClient
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
	}
}

// CreateRagWithOptions creates a new RAG system with the specified options
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

// Modify the existing CreateRag to use CreateRagWithOptions
func (rs *RagServiceImpl) CreateRag(modelName, ragName, folderPath string) error {
	return rs.CreateRagWithOptions(modelName, ragName, folderPath, DocumentLoaderOptions{})
}

// LoadRag loads a RAG system
func (rs *RagServiceImpl) LoadRag(ragName string) (*domain.RagSystem, error) {
	rag, err := rs.ragRepository.Load(ragName)
	if err != nil {
		return nil, fmt.Errorf("error loading RAG '%s': %w", ragName, err)
	}

	return rag, nil
}

// Query performs a query on a RAG system
func (rs *RagServiceImpl) Query(rag *domain.RagSystem, query string, contextSize int) (string, error) {
	// Check if Ollama is available
	var llmClient client.LLMClient

	// Déterminer quel client utiliser en fonction du modèle
	if client.IsOpenAIModel(rag.ModelName) {
		// Pour OpenAI, utiliser le profil spécifié ou celui par défaut
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

	// Generate embedding for the query (toujours avec Ollama)
	queryEmbedding, err := rs.embeddingService.GenerateQueryEmbedding(query, rag.ModelName)
	if err != nil {
		return "", fmt.Errorf("error generating embedding for query: %w", err)
	}

	// Use the provided context size or default to 20
	if contextSize <= 0 {
		contextSize = 20
	}

	// Search for the most relevant chunks
	results := rag.HybridStore.Search(queryEmbedding, contextSize)

	// Build the context
	var context strings.Builder
	context.WriteString("Relevant information:\n\n")

	// Track which documents we've included for reference
	includedDocs := make(map[string]bool)

	for _, result := range results {
		chunk := rag.GetChunkByID(result.ID)
		if chunk != nil {
			// Add chunk content with its metadata
			context.WriteString(fmt.Sprintf("--- %s ---\n%s\n\n",
				chunk.GetMetadataString(), chunk.Content))

			includedDocs[chunk.DocumentID] = true
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
		len(results), len(includedDocs))

	// Generate the response with le client approprié
	response, err := llmClient.GenerateCompletion(rag.ModelName, prompt)
	if err != nil {
		return "", fmt.Errorf("error generating response: %w", err)
	}

	return response, nil
}

// UpdateRag updates an existing RAG system
func (rs *RagServiceImpl) UpdateRag(rag *domain.RagSystem) error {
	err := rs.ragRepository.Save(rag)
	if err != nil {
		return fmt.Errorf("error updating the RAG: %w", err)
	}
	return nil
}

// AddDocsWithOptions adds documents to an existing RAG with the specified options
func (rs *RagServiceImpl) AddDocsWithOptions(ragName string, folderPath string, options DocumentLoaderOptions) error {
	// Load existing RAG
	rag, err := rs.LoadRag(ragName)
	if err != nil {
		return err
	}

	// Load documents with options
	docs, err := rs.documentLoader.LoadDocumentsFromFolderWithOptions(folderPath, options)
	if err != nil {
		return fmt.Errorf("error loading documents: %w", err)
	}

	if len(docs) == 0 {
		return fmt.Errorf("no valid documents found in folder %s", folderPath)
	}

	var allChunks []*domain.DocumentChunk

	// Create chunker service with config from options
	chunkerConfig := ChunkingConfig{
		ChunkSize:        options.ChunkSize,
		ChunkOverlap:     options.ChunkOverlap,
		ChunkingStrategy: options.ChunkingStrategy,
		IncludeMetadata:  true,
	}
	chunkerService := NewChunkerService(chunkerConfig)

	// Process documents
	for _, doc := range docs {
		chunks := chunkerService.ChunkDocument(doc)
		for i, chunk := range chunks {
			chunk.ChunkNumber = i
			chunk.TotalChunks = len(chunks)
		}
		allChunks = append(allChunks, chunks...)
	}

	fmt.Printf("Generated %d chunks using '%s' strategy.\n", len(allChunks), options.ChunkingStrategy)

	// Generate embeddings
	err = rs.embeddingService.GenerateChunkEmbeddings(allChunks, rag.ModelName)
	if err != nil {
		return err
	}

	// Add new chunks
	chunksAdded := 0
	existingChunks := make(map[string]bool)
	for _, chunk := range rag.Chunks {
		existingChunks[chunk.ID] = true
	}
	for _, chunk := range allChunks {
		if !existingChunks[chunk.ID] {
			rag.AddChunk(chunk)
			chunksAdded++
		}
	}

	// Save updated RAG
	return rs.UpdateRag(rag)
}

// Add chunk filter struct
type ChunkFilter struct {
	DocumentSubstring string
	ShowContent       bool
}

func (rs *RagServiceImpl) GetRagChunks(ragName string, filter ChunkFilter) ([]*domain.DocumentChunk, error) {
	rag, err := rs.ragRepository.Load(ragName)
	if err != nil {
		return nil, fmt.Errorf("error loading RAG: %w", err)
	}

	var filtered []*domain.DocumentChunk
	for _, chunk := range rag.Chunks {
		// Apply document filter
		if filter.DocumentSubstring != "" &&
			!strings.Contains(strings.ToLower(chunk.DocumentID), strings.ToLower(filter.DocumentSubstring)) {
			continue
		}

		// Clone chunk to avoid modifying original
		c := *chunk

		// Clear content if not requested
		if !filter.ShowContent {
			c.Content = ""
		}

		filtered = append(filtered, &c)
	}

	return filtered, nil
}

// UpdateModel updates the model of an existing RAG system
func (rs *RagServiceImpl) UpdateModel(ragName string, newModel string) error {
	rag, err := rs.LoadRag(ragName)
	if err != nil {
		return fmt.Errorf("error loading RAG: %w", err)
	}

	rag.ModelName = newModel
	return rs.UpdateRag(rag)
}

// Add method to list all RAGs
func (rs *RagServiceImpl) ListAllRags() ([]string, error) {
	return rs.ragRepository.ListAll()
}

// GetOllamaClient returns the Ollama client
func (rs *RagServiceImpl) GetOllamaClient() *client.OllamaClient {
	return rs.ollamaClient
}

// SetupDirectoryWatching configures a RAG to watch a directory for changes
func (rs *RagServiceImpl) SetupDirectoryWatching(ragName string, dirPath string, watchInterval int, options DocumentLoaderOptions) error {
	// Load the RAG
	rag, err := rs.LoadRag(ragName)
	if err != nil {
		return fmt.Errorf("error loading RAG: %w", err)
	}

	// Check if the directory exists
	dirInfo, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("directory '%s' does not exist", dirPath)
	} else if err != nil {
		return fmt.Errorf("error accessing directory: %w", err)
	}

	if !dirInfo.IsDir() {
		return fmt.Errorf("'%s' is not a directory", dirPath)
	}

	// Set up watching configuration
	rag.WatchedDir = dirPath
	rag.WatchInterval = watchInterval
	rag.WatchEnabled = true
	rag.LastWatchedAt = time.Time{} // Zero time to force first check

	// Save watch options
	rag.WatchOptions = domain.DocumentWatchOptions{
		ExcludeDirs:  options.ExcludeDirs,
		ExcludeExts:  options.ExcludeExts,
		ProcessExts:  options.ProcessExts,
		ChunkSize:    options.ChunkSize,
		ChunkOverlap: options.ChunkOverlap,
	}

	// Update the RAG
	return rs.UpdateRag(rag)
}

// DisableDirectoryWatching disables directory watching for a RAG
func (rs *RagServiceImpl) DisableDirectoryWatching(ragName string) error {
	// Load the RAG
	rag, err := rs.LoadRag(ragName)
	if err != nil {
		return fmt.Errorf("error loading RAG: %w", err)
	}

	// Disable watching
	rag.WatchEnabled = false

	// Update the RAG
	return rs.UpdateRag(rag)
}

// CheckWatchedDirectory manually checks a RAG's watched directory
func (rs *RagServiceImpl) CheckWatchedDirectory(ragName string) (int, error) {
	// Load the RAG
	rag, err := rs.LoadRag(ragName)
	if err != nil {
		return 0, fmt.Errorf("error loading RAG: %w", err)
	}

	// Check if watching is enabled
	if !rag.WatchEnabled || rag.WatchedDir == "" {
		return 0, fmt.Errorf("directory watching is not enabled for RAG '%s'", ragName)
	}

	// Create a file watcher and check for updates
	fileWatcher := NewFileWatcher(rs)
	return fileWatcher.CheckAndUpdateRag(rag)
}

// SetupWebWatching configures a RAG to watch a website for changes
func (rs *RagServiceImpl) SetupWebWatching(ragName string, websiteURL string, watchInterval int, options domain.WebWatchOptions) error {
	// Load the RAG
	rag, err := rs.LoadRag(ragName)
	if err != nil {
		return fmt.Errorf("error loading RAG: %w", err)
	}

	// Validate URL
	_, err = url.Parse(websiteURL)
	if err != nil {
		return fmt.Errorf("invalid website URL: %w", err)
	}

	// Set up watching configuration
	rag.WatchedURL = websiteURL
	rag.WebWatchInterval = watchInterval
	rag.WebWatchEnabled = true
	rag.LastWebWatchAt = time.Time{} // Zero time to force first check

	// Save watch options
	rag.WebWatchOptions = options

	// Update the RAG
	return rs.UpdateRag(rag)
}

// DisableWebWatching disables website watching for a RAG
func (rs *RagServiceImpl) DisableWebWatching(ragName string) error {
	// Load the RAG
	rag, err := rs.LoadRag(ragName)
	if err != nil {
		return fmt.Errorf("error loading RAG: %w", err)
	}

	// Disable watching
	rag.WebWatchEnabled = false

	// Update the RAG
	return rs.UpdateRag(rag)
}

// CheckWatchedWebsite manually checks a RAG's watched website for new content
func (rs *RagServiceImpl) CheckWatchedWebsite(ragName string) (int, error) {
	// Load the RAG
	rag, err := rs.LoadRag(ragName)
	if err != nil {
		return 0, fmt.Errorf("error loading RAG: %w", err)
	}

	// Check if watching is enabled
	if !rag.WebWatchEnabled || rag.WatchedURL == "" {
		return 0, fmt.Errorf("website watching is not enabled for RAG '%s'", ragName)
	}

	// Create a web watcher and check for updates
	webWatcher := NewWebWatcher(rs)
	return webWatcher.CheckAndUpdateRag(rag)
}

// Helper function to truncate string for display
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
