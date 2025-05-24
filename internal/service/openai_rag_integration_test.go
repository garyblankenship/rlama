package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/pkg/vector"
)

// Mock types for testing
type mockRAGQuery struct {
	Query            string
	Model            string
	EmbeddingModel   string
	MaxResults       int
	MinSimilarity    float64
}

type mockRAGResult struct {
	Answer  string
	Sources []string
	Model   string
}

type mockRAGService struct {
	vectorStore      *vector.InternalVectorStore
	embeddingService *EmbeddingService
	llmClient        *client.OpenAIClient
	documents        map[string]*domain.Document
}

func (m *mockRAGService) Query(query *mockRAGQuery) (*mockRAGResult, error) {
	// Generate query embedding
	queryEmbedding, err := m.embeddingService.GenerateQueryEmbedding(query.Query, query.EmbeddingModel)
	if err != nil {
		return nil, err
	}

	// Search for similar documents
	results := m.vectorStore.Search(queryEmbedding, query.MaxResults)
	
	// Build context from retrieved documents
	var context string
	var sources []string
	for _, result := range results {
		if result.Score >= query.MinSimilarity {
			if doc, exists := m.documents[result.ID]; exists {
				context += doc.Content + "\n\n"
				sources = append(sources, doc.Name)
			}
		}
	}

	// Generate response using LLM
	prompt := fmt.Sprintf("Context: %s\n\nQuestion: %s\n\nAnswer:", context, query.Query)
	answer, err := m.llmClient.GenerateCompletion(query.Model, prompt)
	if err != nil {
		return nil, err
	}

	return &mockRAGResult{
		Answer:  answer,
		Sources: sources,
		Model:   query.Model,
	}, nil
}

func TestOpenAIRAGService_Integration(t *testing.T) {
	// Create mock OpenAI-compatible server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		if r.URL.Path == "/embeddings" {
			response := client.OpenAIEmbeddingResponse{
				Object: "list",
				Data: []struct {
					Object    string    `json:"object"`
					Embedding []float32 `json:"embedding"`
					Index     int       `json:"index"`
				}{
					{
						Object:    "embedding",
						Embedding: generateTestEmbedding(1536), // Standard OpenAI embedding size
						Index:     0,
					},
				},
				Model: "text-embedding-ada-002",
				Usage: struct {
					PromptTokens int `json:"prompt_tokens"`
					TotalTokens  int `json:"total_tokens"`
				}{
					PromptTokens: 5,
					TotalTokens:  5,
				},
			}
			json.NewEncoder(w).Encode(response)
		} else if r.URL.Path == "/chat/completions" {
			response := client.OpenAICompletionResponse{
				ID:      "test-completion",
				Object:  "chat.completion",
				Created: time.Now().Unix(),
				Model:   "gpt-3.5-turbo",
				Choices: []struct {
					Index        int                  `json:"index"`
					Message      client.OpenAIMessage `json:"message"`
					FinishReason string               `json:"finish_reason"`
				}{
					{
						Index: 0,
						Message: client.OpenAIMessage{
							Role:    "assistant",
							Content: "Based on the provided context, this is a test response from the OpenAI-compatible RAG system.",
						},
						FinishReason: "stop",
					},
				},
				Usage: struct {
					PromptTokens     int `json:"prompt_tokens"`
					CompletionTokens int `json:"completion_tokens"`
					TotalTokens      int `json:"total_tokens"`
				}{
					PromptTokens:     20,
					CompletionTokens: 15,
					TotalTokens:      35,
				},
			}
			json.NewEncoder(w).Encode(response)
		}
	}))
	defer server.Close()

	// Setup RAG components
	vectorStore := vector.NewInternalVectorStore(1536) // OpenAI embedding dimension
	
	// Create OpenAI client
	openaiClient := client.NewGenericOpenAIClient(server.URL, "test-api-key")
	
	// Create embedding service
	embeddingService := NewEmbeddingService(openaiClient)
	
	// Create a simple mock RAG service
	ragService := &mockRAGService{
		vectorStore:      vectorStore,
		embeddingService: embeddingService,
		llmClient:        openaiClient,
		documents:        make(map[string]*domain.Document),
	}

	tests := []struct {
		name           string
		setupDocuments func() error
		query          string
		expectError    bool
		validateResult func(result *mockRAGResult) error
	}{
		{
			name: "OpenAI RAG with document retrieval",
			setupDocuments: func() error {
				// Add test documents
				docs := []*domain.Document{
					{
						ID:        "doc1",
						Name:      "test1.txt",
						Content:   "This is a test document about artificial intelligence and machine learning.",
						CreatedAt: time.Now(),
					},
					{
						ID:        "doc2", 
						Name:      "test2.txt",
						Content:   "This document discusses natural language processing and deep learning techniques.",
						CreatedAt: time.Now(),
					},
				}
				
				for _, doc := range docs {
					// Add document to mock service
					ragService.documents[doc.ID] = doc
					
					// Generate and store embedding
					embedding, err := embeddingService.GenerateQueryEmbedding(doc.Content, "text-embedding-ada-002")
					if err != nil {
						return err
					}
					
					vectorStore.Add(doc.ID, embedding)
				}
				return nil
			},
			query:       "What is machine learning?",
			expectError: false,
			validateResult: func(result *mockRAGResult) error {
				if result.Answer == "" {
					return ErrEmptyResponse
				}
				if len(result.Sources) == 0 {
					return ErrNoSources
				}
				if result.Model == "" {
					return ErrMissingModel
				}
				return nil
			},
		},
		{
			name: "OpenAI RAG with no matching documents",
			setupDocuments: func() error {
				// Add unrelated document
				doc := &domain.Document{
					ID:        "doc3",
					Name:      "unrelated.txt", 
					Content:   "This document is about cooking recipes and has nothing to do with technology.",
					CreatedAt: time.Now(),
				}
				
				ragService.documents[doc.ID] = doc
				
				embedding, err := embeddingService.GenerateQueryEmbedding(doc.Content, "text-embedding-ada-002")
				if err != nil {
					return err
				}
				
				vectorStore.Add(doc.ID, embedding)
				return nil
			},
			query:       "What is quantum computing?",
			expectError: false,
			validateResult: func(result *mockRAGResult) error {
				if result.Answer == "" {
					return ErrEmptyResponse
				}
				// Should still return a response even with no good matches
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up vector store and documents
			vectorStore = vector.NewInternalVectorStore(1536)
			ragService.vectorStore = vectorStore
			ragService.documents = make(map[string]*domain.Document)
			
			// Setup documents for this test
			if err := tt.setupDocuments(); err != nil {
				t.Fatalf("Failed to setup documents: %v", err)
			}
			
			// Execute RAG query
			result, err := ragService.Query(&mockRAGQuery{
				Query:            tt.query,
				Model:            "gpt-3.5-turbo",
				EmbeddingModel:   "text-embedding-ada-002",
				MaxResults:       5,
				MinSimilarity:    0.1,
			})
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != nil && tt.validateResult != nil {
					if err := tt.validateResult(result); err != nil {
						t.Errorf("Result validation failed: %v", err)
					}
				}
			}
		})
	}
}

// Helper function to generate test embeddings
func generateTestEmbedding(size int) []float32 {
	embedding := make([]float32, size)
	for i := range embedding {
		embedding[i] = float32(i) / float32(size) // Simple pattern
	}
	return embedding
}

// Custom error types for validation
var (
	ErrEmptyResponse = fmt.Errorf("empty response")
	ErrNoSources     = fmt.Errorf("no sources found")
	ErrMissingModel  = fmt.Errorf("missing model information")
)