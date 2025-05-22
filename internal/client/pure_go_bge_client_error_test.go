package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPureGoBGEClient_CreationErrors(t *testing.T) {
	tests := []struct {
		name        string
		modelPath   string
		usePureGo   bool
		fallbackURL string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "NonexistentModelPath",
			modelPath:   "/nonexistent/model/path",
			usePureGo:   true,
			fallbackURL: "http://localhost:8000",
			expectError: true,
			errorMsg:    "failed to create tokenizer",
		},
		{
			name:        "EmptyModelPath",
			modelPath:   "",
			usePureGo:   true,
			fallbackURL: "http://localhost:8000",
			expectError: true,
			errorMsg:    "failed to create tokenizer",
		},
		{
			name:        "ValidConfigBothModes",
			modelPath:   createTempModelDir(t),
			usePureGo:   true,
			fallbackURL: "http://localhost:8000",
			expectError: false,
		},
		{
			name:        "ValidConfigFallbackMode",
			modelPath:   createTempModelDir(t),
			usePureGo:   false,
			fallbackURL: "http://localhost:8000",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.modelPath != "" && tt.modelPath != "/nonexistent/model/path" {
				defer os.RemoveAll(tt.modelPath)
			}

			client, err := NewPureGoBGEClient(tt.modelPath, tt.usePureGo, tt.fallbackURL)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain '%s', got: %s", tt.errorMsg, err.Error())
				}
				if client != nil {
					t.Errorf("Expected nil client on error, got non-nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if client == nil {
					t.Errorf("Expected non-nil client on success")
				}
			}
		})
	}
}

func TestPureGoBGEClient_HealthCheckErrors(t *testing.T) {
	modelPath := createTempModelDir(t)
	defer os.RemoveAll(modelPath)

	tests := []struct {
		name        string
		setupClient func() *PureGoBGEClient
		expectError bool
	}{
		{
			name: "ValidClient",
			setupClient: func() *PureGoBGEClient {
				client, _ := NewPureGoBGEClient(modelPath, true, "http://localhost:8000")
				return client
			},
			expectError: false,
		},
		{
			name: "NilTokenizer",
			setupClient: func() *PureGoBGEClient {
				return &PureGoBGEClient{
					tokenizer: nil,
				}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := tt.setupClient()
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := client.Health(ctx)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestPureGoBGEClient_NetworkErrors(t *testing.T) {
	modelPath := createTempModelDir(t)
	defer os.RemoveAll(modelPath)

	tests := []struct {
		name        string
		serverSetup func() *httptest.Server
		expectError bool
		errorMsg    string
	}{
		{
			name: "ServerReturns500",
			serverSetup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
			},
			expectError: true,
			errorMsg:    "returned status 500",
		},
		{
			name: "ServerReturns404",
			serverSetup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}))
			},
			expectError: true,
			errorMsg:    "returned status 404",
		},
		{
			name: "ServerReturnsInvalidJSON",
			serverSetup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte("invalid json"))
				}))
			},
			expectError: true,
			errorMsg:    "failed to decode response",
		},
		{
			name: "ServerReturnsErrorMessage",
			serverSetup: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.Write([]byte(`{"scores": [], "error": "Model loading failed"}`))
				}))
			},
			expectError: true,
			errorMsg:    "inference error: Model loading failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.serverSetup()
			defer server.Close()

			client, err := NewPureGoBGEClient(modelPath, true, server.URL)
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			_, err = client.Rerank(ctx, "test query", []string{"doc1", "doc2"}, 2)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain '%s', got: %s", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestPureGoBGEClient_TimeoutHandling(t *testing.T) {
	modelPath := createTempModelDir(t)
	defer os.RemoveAll(modelPath)

	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Longer than client timeout
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"scores": [0.5]}`))
	}))
	defer server.Close()

	client, err := NewPureGoBGEClient(modelPath, true, server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Set a short timeout
	client.httpClient.Timeout = 500 * time.Millisecond

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err = client.Rerank(ctx, "test query", []string{"doc1"}, 1)

	if err == nil {
		t.Errorf("Expected timeout error but got none")
	} else if !strings.Contains(err.Error(), "timeout") && !strings.Contains(err.Error(), "deadline") {
		t.Logf("Got error (may be timeout-related): %v", err)
	}
}

func TestPureGoBGEClient_InvalidInputHandling(t *testing.T) {
	modelPath := createTempModelDir(t)
	defer os.RemoveAll(modelPath)

	client, err := NewPureGoBGEClient(modelPath, true, "http://localhost:8000")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []struct {
		name      string
		query     string
		documents []string
		topK      int
	}{
		{
			name:      "EmptyQuery",
			query:     "",
			documents: []string{"doc1", "doc2"},
			topK:      2,
		},
		{
			name:      "EmptyDocuments",
			query:     "test query",
			documents: []string{},
			topK:      2,
		},
		{
			name:      "NilDocuments",
			query:     "test query",
			documents: nil,
			topK:      2,
		},
		{
			name:      "NegativeTopK",
			query:     "test query",
			documents: []string{"doc1", "doc2"},
			topK:      -1,
		},
		{
			name:      "ZeroTopK",
			query:     "test query",
			documents: []string{"doc1", "doc2"},
			topK:      0,
		},
		{
			name:      "TopKLargerThanDocs",
			query:     "test query",
			documents: []string{"doc1"},
			topK:      10,
		},
		{
			name:      "VeryLongQuery",
			query:     strings.Repeat("very long query ", 1000),
			documents: []string{"doc1"},
			topK:      1,
		},
		{
			name:      "VeryLongDocument",
			query:     "test query",
			documents: []string{strings.Repeat("very long document ", 1000)},
			topK:      1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// This should not panic or return invalid results
			results, err := client.Rerank(ctx, tt.query, tt.documents, tt.topK)

			// We expect either an error or valid results, but no panics
			if err != nil {
				t.Logf("Got expected error for %s: %v", tt.name, err)
			} else {
				// If no error, validate results
				if tt.documents != nil && len(tt.documents) > 0 {
					expectedLen := len(tt.documents)
					if tt.topK > 0 && tt.topK < len(tt.documents) {
						expectedLen = tt.topK
					}
					
					if len(results) > expectedLen {
						t.Errorf("Got more results than expected: %d > %d", len(results), expectedLen)
					}
				}
				
				t.Logf("Got valid results for %s: %d results", tt.name, len(results))
			}
		})
	}
}

func TestPureGoBGEClient_ConfigurationEdgeCases(t *testing.T) {
	modelPath := createTempModelDir(t)
	defer os.RemoveAll(modelPath)

	client, err := NewPureGoBGEClient(modelPath, true, "http://localhost:8000")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tests := []struct {
		name      string
		maxLength int
		expectErr bool
	}{
		{"NegativeMaxLength", -1, false}, // Should handle gracefully
		{"ZeroMaxLength", 0, false},      // Should handle gracefully
		{"VeryLargeMaxLength", 100000, false}, // Should handle gracefully
		{"NormalMaxLength", 512, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("SetMaxLength panicked with value %d: %v", tt.maxLength, r)
				}
			}()

			client.SetMaxLength(tt.maxLength)

			// Test that tokenization still works
			tokenizer := client.GetTokenizer()
			if tokenizer != nil {
				_, _, _ = tokenizer.Encode("test", 128)
			}

			t.Logf("SetMaxLength(%d) completed successfully", tt.maxLength)
		})
	}
}

// Helper function to create a temporary model directory with tokenizer.json
func createTempModelDir(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "test_model_*")
	if err != nil {
		t.Fatal(err)
	}

	// Create tokenizer.json
	tokenizerPath := filepath.Join(tmpDir, "tokenizer.json")
	minimalConfigPath := createMinimalTokenizerConfig(t)
	defer os.Remove(minimalConfigPath)

	// Copy the minimal config to the model directory
	data, err := os.ReadFile(minimalConfigPath)
	if err != nil {
		os.RemoveAll(tmpDir)
		t.Fatal(err)
	}

	if err := os.WriteFile(tokenizerPath, data, 0644); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatal(err)
	}

	return tmpDir
}