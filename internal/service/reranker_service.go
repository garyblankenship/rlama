package service

import (
	"fmt"
	"sort"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/pkg/vector"
)

// RerankerOptions defines configuration options for the reranking process
type RerankerOptions struct {
	// TopK is the number of documents to return after reranking
	TopK int

	// InitialK is the number of documents to retrieve from the initial search
	// before reranking (should be >= TopK)
	InitialK int

	// RerankerModel is the model to use for reranking
	// If empty, will default to the same model used for embedding
	RerankerModel string

	// ScoreThreshold is the minimum relevance score (0-1) for a document to be included
	// Documents with scores below this threshold will be filtered out
	ScoreThreshold float64

	// RerankerWeight defines the weight of the reranker score vs vector similarity
	// 0.0 = use only vector similarity, 1.0 = use only reranker scores
	RerankerWeight float64
}

// DefaultRerankerOptions returns the default options for reranking
func DefaultRerankerOptions() RerankerOptions {
	return RerankerOptions{
		TopK:           5,
		InitialK:       20,
		RerankerModel:  "",
		ScoreThreshold: 0.0,
		RerankerWeight: 0.7, // 70% reranker score, 30% vector similarity
	}
}

// RerankerService handles document reranking using cross-encoder models
type RerankerService struct {
	ollamaClient *client.OllamaClient
}

// NewRerankerService creates a new instance of RerankerService
func NewRerankerService(ollamaClient *client.OllamaClient) *RerankerService {
	if ollamaClient == nil {
		ollamaClient = client.NewDefaultOllamaClient()
	}

	return &RerankerService{
		ollamaClient: ollamaClient,
	}
}

// RankedResult represents a document with its relevance score after reranking
type RankedResult struct {
	Chunk         *domain.DocumentChunk
	VectorScore   float64
	RerankerScore float64
	FinalScore    float64
}

// Rerank takes initial retrieval results and reruns them through a cross-encoder for more accurate ranking
func (rs *RerankerService) Rerank(
	query string,
	rag *domain.RagSystem,
	initialResults []vector.SearchResult,
	options RerankerOptions,
) ([]RankedResult, error) {
	// Create an empty result if no documents were found
	if len(initialResults) == 0 {
		return []RankedResult{}, nil
	}

	// Use model from options or fall back to the RAG model
	modelName := options.RerankerModel
	if modelName == "" {
		modelName = rag.ModelName
	}

	fmt.Printf("Using reranker model: %s\n", modelName)

	// Prepare the cross-encoder prompt template
	promptTemplate := `You are a document relevance scoring system. Rate how relevant a document is to a query on a scale from 0 to 1, where 0 is completely irrelevant and 1 is highly relevant.

Query: %s

Document:
%s

Relevance score (output only a single number between 0 and 1):
`

	// Get chunks and score them
	var rankedResults []RankedResult
	for _, result := range initialResults {
		chunk := rag.GetChunkByID(result.ID)
		if chunk == nil {
			continue
		}

		// Prepare prompt for this chunk
		prompt := fmt.Sprintf(promptTemplate, query, chunk.Content)

		// Get reranking score from the model
		response, err := rs.ollamaClient.GenerateCompletion(modelName, prompt)
		if err != nil {
			return nil, fmt.Errorf("error generating reranking score: %w", err)
		}

		// Parse the response as a float (score)
		var score float64
		_, err = fmt.Sscanf(response, "%f", &score)
		if err != nil {
			// If parsing fails, use a default score (based on vector similarity)
			score = result.Score
		}

		// Ensure score is in range [0,1]
		if score < 0 {
			score = 0
		} else if score > 1 {
			score = 1
		}

		// Calculate final score as weighted combination of vector and reranker scores
		finalScore := (options.RerankerWeight * score) + ((1 - options.RerankerWeight) * result.Score)

		// Add to results if above threshold
		if finalScore >= options.ScoreThreshold {
			rankedResults = append(rankedResults, RankedResult{
				Chunk:         chunk,
				VectorScore:   result.Score,
				RerankerScore: score,
				FinalScore:    finalScore,
			})
		}
	}

	// Sort by final score (descending)
	sort.Slice(rankedResults, func(i, j int) bool {
		return rankedResults[i].FinalScore > rankedResults[j].FinalScore
	})

	// Apply top-k limit
	if options.TopK > 0 && len(rankedResults) > options.TopK {
		fmt.Printf("Limiting reranked results from %d to top %d\n", len(rankedResults), options.TopK)
		rankedResults = rankedResults[:options.TopK]
	}

	return rankedResults, nil
}
