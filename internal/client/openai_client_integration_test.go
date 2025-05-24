package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/repository"
)

func TestOpenAIClient_Integration(t *testing.T) {
	// Test with mock server to avoid real API calls
	tests := []struct {
		name           string
		setupServer    func() *httptest.Server
		setupClient    func(baseURL string) *OpenAIClient
		testCompletion bool
		testEmbedding  bool
		expectError    bool
	}{
		{
			name: "OpenAI API compatible completion",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/chat/completions" {
						response := OpenAICompletionResponse{
							ID:      "test-completion",
							Object:  "chat.completion",
							Created: time.Now().Unix(),
							Model:   "gpt-3.5-turbo",
							Choices: []struct {
								Index        int           `json:"index"`
								Message      OpenAIMessage `json:"message"`
								FinishReason string        `json:"finish_reason"`
							}{
								{
									Index:        0,
									Message:      OpenAIMessage{Role: "assistant", Content: "Test response"},
									FinishReason: "stop",
								},
							},
							Usage: struct {
								PromptTokens     int `json:"prompt_tokens"`
								CompletionTokens int `json:"completion_tokens"`
								TotalTokens      int `json:"total_tokens"`
							}{
								PromptTokens:     10,
								CompletionTokens: 5,
								TotalTokens:      15,
							},
						}
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(response)
					}
				}))
			},
			setupClient: func(baseURL string) *OpenAIClient {
				return NewGenericOpenAIClient(baseURL, "test-api-key")
			},
			testCompletion: true,
			expectError:    false,
		},
		{
			name: "OpenAI API compatible embedding",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if r.URL.Path == "/embeddings" {
						response := OpenAIEmbeddingResponse{
							Object: "list",
							Data: []struct {
								Object    string    `json:"object"`
								Embedding []float32 `json:"embedding"`
								Index     int       `json:"index"`
							}{
								{
									Object:    "embedding",
									Embedding: []float32{0.1, 0.2, 0.3, 0.4, 0.5},
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
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(response)
					}
				}))
			},
			setupClient: func(baseURL string) *OpenAIClient {
				return NewGenericOpenAIClient(baseURL, "test-api-key")
			},
			testEmbedding: true,
			expectError:   false,
		},
		{
			name: "Local OpenAI-compatible server (no API key)",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// Accept requests with no authorization or empty bearer token for local servers
					auth := r.Header.Get("Authorization")
					if auth != "" && auth != "Bearer " && auth != "Bearer" {
						w.WriteHeader(http.StatusUnauthorized)
						return
					}
					
					if r.URL.Path == "/chat/completions" {
						response := OpenAICompletionResponse{
							ID:      "local-completion",
							Object:  "chat.completion",
							Created: time.Now().Unix(),
							Model:   "local-model",
							Choices: []struct {
								Index        int           `json:"index"`
								Message      OpenAIMessage `json:"message"`
								FinishReason string        `json:"finish_reason"`
							}{
								{
									Index:        0,
									Message:      OpenAIMessage{Role: "assistant", Content: "Local response"},
									FinishReason: "stop",
								},
							},
						}
						w.Header().Set("Content-Type", "application/json")
						json.NewEncoder(w).Encode(response)
					}
				}))
			},
			setupClient: func(baseURL string) *OpenAIClient {
				return NewGenericOpenAIClient(baseURL, "") // No API key for local server
			},
			testCompletion: true,
			expectError:    false,
		},
		{
			name: "Error handling - 401 Unauthorized",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(`{"error": {"message": "Invalid API key"}}`))
				}))
			},
			setupClient: func(baseURL string) *OpenAIClient {
				return NewGenericOpenAIClient(baseURL, "invalid-key")
			},
			testCompletion: true,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupServer()
			defer server.Close()

			client := tt.setupClient(server.URL)

			if tt.testCompletion {
				response, err := client.GenerateCompletion("test-model", "Test prompt")
				if tt.expectError {
					if err == nil {
						t.Errorf("Expected error but got none")
					}
				} else {
					if err != nil {
						t.Errorf("Unexpected error: %v", err)
					}
					if response == "" {
						t.Errorf("Expected non-empty response")
					}
				}
			}

			if tt.testEmbedding {
				embedding, err := client.GenerateEmbedding("test-model", "Test text")
				if tt.expectError {
					if err == nil {
						t.Errorf("Expected error but got none")
					}
				} else {
					if err != nil {
						t.Errorf("Unexpected error: %v", err)
					}
					if len(embedding) == 0 {
						t.Errorf("Expected non-empty embedding")
					}
				}
			}
		})
	}
}

func TestOpenAIClientWithProfile_Integration(t *testing.T) {
	// Create a temporary profile for testing
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Unsetenv("HOME")

	// Create test profile
	profileRepo := repository.NewProfileRepository()
	testProfile := &domain.APIProfile{
		Name:       "test-openai",
		Provider:   "openai",
		BaseURL:    "https://api.openai.com/v1",
		APIKey:     "test-api-key",
		CreatedAt:  time.Now(),
		LastUsedAt: time.Now(),
	}
	
	err := profileRepo.Save(testProfile)
	if err != nil {
		t.Fatalf("Failed to save test profile: %v", err)
	}

	tests := []struct {
		name        string
		profileName string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid OpenAI profile",
			profileName: "test-openai",
			expectError: false,
		},
		{
			name:        "Empty profile name (fallback to env)",
			profileName: "",
			expectError: false,
		},
		{
			name:        "Non-existent profile",
			profileName: "non-existent",
			expectError: true,
			errorMsg:    "error loading profile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.profileName == "" {
				// Set environment variable for fallback test
				os.Setenv("OPENAI_API_KEY", "env-api-key")
				defer os.Unsetenv("OPENAI_API_KEY")
			}

			client, err := NewOpenAIClientWithProfile(tt.profileName)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if tt.errorMsg != "" && err != nil {
					if !containsString(err.Error(), tt.errorMsg) {
						t.Errorf("Expected error message to contain '%s', got: %v", tt.errorMsg, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if client == nil {
					t.Errorf("Expected non-nil client")
				}
				if client != nil {
					if client.BaseURL == "" {
						t.Errorf("Expected non-empty BaseURL")
					}
					if tt.profileName == "" && client.APIKey != "env-api-key" {
						t.Errorf("Expected API key from environment, got: %s", client.APIKey)
					}
				}
			}
		})
	}
}

func TestOpenAIClient_CheckOpenAIAndModel(t *testing.T) {
	tests := []struct {
		name        string
		client      *OpenAIClient
		modelName   string
		expectError bool
		errorMsg    string
	}{
		{
			name: "Official OpenAI with API key",
			client: &OpenAIClient{
				BaseURL: "https://api.openai.com/v1",
				APIKey:  "test-key",
			},
			modelName:   "gpt-3.5-turbo",
			expectError: false,
		},
		{
			name: "Official OpenAI without API key",
			client: &OpenAIClient{
				BaseURL: "https://api.openai.com/v1",
				APIKey:  "",
			},
			modelName:   "gpt-3.5-turbo",
			expectError: true,
			errorMsg:    "OPENAI_API_KEY environment variable not set",
		},
		{
			name: "Local OpenAI-compatible without API key",
			client: &OpenAIClient{
				BaseURL: "http://localhost:8080/v1",
				APIKey:  "",
			},
			modelName:   "local-model",
			expectError: false,
		},
		{
			name: "LM Studio compatible",
			client: &OpenAIClient{
				BaseURL: "http://localhost:1234/v1",
				APIKey:  "",
			},
			modelName:   "local-model",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.client.CheckOpenAIAndModel(tt.modelName)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				if tt.errorMsg != "" && err != nil {
					if !containsString(err.Error(), tt.errorMsg) {
						t.Errorf("Expected error message to contain '%s', got: %v", tt.errorMsg, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		(len(s) > len(substr) && 
			(s[:len(substr)] == substr || 
				s[len(s)-len(substr):] == substr || 
				containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}