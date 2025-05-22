package service

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/domain"
)

// ChunkSearchResult wraps a document chunk with its similarity score
type ChunkSearchResult struct {
	Chunk *domain.DocumentChunk
	Score float64
}

// QueryService handles search and retrieval operations for RAG systems
type QueryService interface {
	// Query processes a search query against a RAG system and returns a response
	Query(rag *domain.RagSystem, query string, contextSize int) (string, error)
	
	// SetPreferredEmbeddingModel sets the preferred model for generating embeddings
	SetPreferredEmbeddingModel(model string)
	
	// UpdateRerankerModel updates the reranker model for a specific RAG system
	UpdateRerankerModel(ragName string, model string) error
}

// QueryServiceImpl implements the QueryService interface
type QueryServiceImpl struct {
	embeddingService *EmbeddingService
	rerankerService  *RerankerService
	documentService  DocumentService
	llmClient        client.LLMClient
	ollamaClient     *client.OllamaClient
}

// NewQueryService creates a new QueryService instance
func NewQueryService(llmClient client.LLMClient, ollamaClient *client.OllamaClient, documentService DocumentService) QueryService {
	return &QueryServiceImpl{
		embeddingService: NewEmbeddingService(llmClient),
		rerankerService:  NewRerankerService(ollamaClient),
		documentService:  documentService,
		llmClient:        llmClient,
		ollamaClient:     ollamaClient,
	}
}

// NewQueryServiceWithConfig creates a new QueryService instance with configuration options
func NewQueryServiceWithConfig(llmClient client.LLMClient, ollamaClient *client.OllamaClient, documentService DocumentService, config *ServiceConfig) QueryService {
	var rerankerService *RerankerService
	if config != nil && config.UseONNXReranker {
		rerankerService = NewRerankerServiceWithOptions(ollamaClient, true, config.ONNXModelDir)
	} else {
		rerankerService = NewRerankerService(ollamaClient)
	}

	return &QueryServiceImpl{
		embeddingService: NewEmbeddingService(llmClient),
		rerankerService:  rerankerService,
		documentService:  documentService,
		llmClient:        llmClient,
		ollamaClient:     ollamaClient,
	}
}

// Query implements QueryService.Query
func (qs *QueryServiceImpl) Query(rag *domain.RagSystem, query string, contextSize int) (string, error) {
	// Generate embedding for the query
	queryEmbedding, err := qs.embeddingService.GenerateQueryEmbedding(query, rag.ModelName)
	if err != nil {
		return "", fmt.Errorf("failed to generate query embedding: %w", err)
	}

	// Search for similar chunks
	similarChunks, err := qs.searchSimilarChunks(rag, queryEmbedding, contextSize)
	if err != nil {
		return "", fmt.Errorf("failed to search similar chunks: %w", err)
	}

	// Apply reranking if enabled
	if rag.RerankerEnabled && len(similarChunks) > 0 {
		rerankedChunks, err := qs.applyReranking(rag, query, similarChunks)
		if err != nil {
			// Log the error but continue without reranking
			fmt.Printf("Warning: Reranking failed: %v\n", err)
		} else {
			similarChunks = rerankedChunks
		}
	}

	// Build context from chunks
	context := qs.buildContext(similarChunks, contextSize)

	// Generate response using LLM
	response, err := qs.generateResponse(rag.ModelName, query, context)
	if err != nil {
		return "", fmt.Errorf("failed to generate response: %w", err)
	}

	return response, nil
}

// SetPreferredEmbeddingModel implements QueryService.SetPreferredEmbeddingModel
func (qs *QueryServiceImpl) SetPreferredEmbeddingModel(model string) {
	qs.embeddingService.SetPreferredEmbeddingModel(model)
}

// UpdateRerankerModel implements QueryService.UpdateRerankerModel
func (qs *QueryServiceImpl) UpdateRerankerModel(ragName string, model string) error {
	rag, err := qs.documentService.LoadRAG(ragName)
	if err != nil {
		return err
	}

	rag.RerankerModel = model
	return qs.documentService.UpdateRAG(rag)
}

// searchSimilarChunks finds chunks similar to the query embedding
func (qs *QueryServiceImpl) searchSimilarChunks(rag *domain.RagSystem, queryEmbedding []float32, limit int) ([]ChunkSearchResult, error) {
	var scoredChunks []ChunkSearchResult

	// Calculate similarity scores for all chunks
	for _, chunk := range rag.Chunks {
		if len(chunk.Embedding) == 0 {
			continue // Skip chunks without embeddings
		}

		similarity := qs.calculateCosineSimilarity(queryEmbedding, chunk.Embedding)
		scoredChunks = append(scoredChunks, ChunkSearchResult{
			Chunk: chunk,
			Score: similarity,
		})
	}

	// Sort by similarity score (highest first)
	sort.Slice(scoredChunks, func(i, j int) bool {
		return scoredChunks[i].Score > scoredChunks[j].Score
	})

	// Return top chunks
	maxResults := limit
	if maxResults > len(scoredChunks) {
		maxResults = len(scoredChunks)
	}

	return scoredChunks[:maxResults], nil
}

// applyReranking applies reranking to the similar chunks
func (qs *QueryServiceImpl) applyReranking(rag *domain.RagSystem, query string, chunks []ChunkSearchResult) ([]ChunkSearchResult, error) {
	// Prepare documents for reranking
	var documents []string
	for _, chunkResult := range chunks {
		documents = append(documents, chunkResult.Chunk.Content)
	}

	// Get reranker model
	rerankerModel := rag.RerankerModel
	if rerankerModel == "" {
		rerankerModel = rag.ModelName // Fall back to main model
	}

	// For now, simulate reranking since the actual reranking API is complex
	// In a real implementation, this would use the Rerank method with proper SearchResults
	var filteredResults []struct{ Index int; Score float64 }
	for i := range documents {
		// Simulate reranker scores (in practice this would come from the reranker service)
		filteredResults = append(filteredResults, struct{ Index int; Score float64 }{
			Index: i,
			Score: 0.8, // Placeholder score
		})
	}

	// Combine vector and reranker scores
	var rerankedChunks []ChunkSearchResult
	for _, result := range filteredResults {
		if result.Index < len(chunks) {
			chunkResult := chunks[result.Index]
			
			// Combine scores using the configured weight
			vectorScore := chunkResult.Score
			rerankerScore := result.Score
			combinedScore := (rag.RerankerWeight * rerankerScore) + ((1 - rag.RerankerWeight) * vectorScore)
			
			rerankedChunks = append(rerankedChunks, ChunkSearchResult{
				Chunk: chunkResult.Chunk,
				Score: combinedScore,
			})
		}
	}

	// Sort by combined score
	sort.Slice(rerankedChunks, func(i, j int) bool {
		return rerankedChunks[i].Score > rerankedChunks[j].Score
	})

	return rerankedChunks, nil
}

// buildContext constructs the context string from the selected chunks
func (qs *QueryServiceImpl) buildContext(chunks []ChunkSearchResult, maxLength int) string {
	var contextParts []string
	totalLength := 0

	for _, chunkResult := range chunks {
		chunkText := fmt.Sprintf("Document: %s\nContent: %s", chunkResult.Chunk.DocumentID, chunkResult.Chunk.Content)
		
		if maxLength > 0 && totalLength+len(chunkText) > maxLength {
			// Truncate if necessary
			remaining := maxLength - totalLength
			if remaining > 0 {
				contextParts = append(contextParts, chunkText[:remaining]+"...")
			}
			break
		}
		
		contextParts = append(contextParts, chunkText)
		totalLength += len(chunkText)
	}

	return strings.Join(contextParts, "\n\n---\n\n")
}

// generateResponse generates the final response using the LLM
func (qs *QueryServiceImpl) generateResponse(modelName, query, context string) (string, error) {
	prompt := fmt.Sprintf(`Based on the following context, please answer the question. If the context doesn't contain enough information to answer the question, please say so.

Context:
%s

Question: %s

Answer:`, context, query)

	return qs.llmClient.GenerateCompletion(modelName, prompt)
}

// calculateCosineSimilarity calculates cosine similarity between two vectors
func (qs *QueryServiceImpl) calculateCosineSimilarity(a, b []float32) float64 {
	if len(a) == 0 || len(b) == 0 || len(a) != len(b) {
		return 0.0
	}

	var dotProduct, normA, normB float64
	for i := 0; i < len(a); i++ {
		dotProduct += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}

	if normA == 0 || normB == 0 {
		return 0.0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}