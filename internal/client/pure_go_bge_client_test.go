package client

import (
	"context"
	"path/filepath"
	"testing"
	"time"
)

func TestPureGoBGEClient_Creation(t *testing.T) {
	modelPath := filepath.Join("..", "..", "models", "bge-reranker-large-onnx")
	
	client, err := NewPureGoBGEClient(modelPath, true, "http://localhost:8000")
	if err != nil {
		t.Fatalf("Failed to create pure Go BGE client: %v", err)
	}
	
	if client == nil {
		t.Fatal("Client is nil")
	}
	
	if client.tokenizer == nil {
		t.Fatal("Tokenizer is nil")
	}
	
	t.Logf("Pure Go BGE client created successfully")
}

func TestPureGoBGEClient_Health(t *testing.T) {
	modelPath := filepath.Join("..", "..", "models", "bge-reranker-large-onnx")
	
	client, err := NewPureGoBGEClient(modelPath, true, "http://localhost:8000")
	if err != nil {
		t.Fatalf("Failed to create pure Go BGE client: %v", err)
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	err = client.Health(ctx)
	if err != nil {
		t.Fatalf("Health check failed: %v", err)
	}
	
	t.Logf("Health check passed")
}

func TestPureGoBGEClient_TokenizationOnly(t *testing.T) {
	modelPath := filepath.Join("..", "..", "models", "bge-reranker-large-onnx")
	
	client, err := NewPureGoBGEClient(modelPath, true, "http://localhost:8000")
	if err != nil {
		t.Fatalf("Failed to create pure Go BGE client: %v", err)
	}
	
	// Test tokenization functionality
	query := "What is artificial intelligence?"
	documents := []string{
		"Artificial intelligence is a branch of computer science.",
		"Machine learning is a subset of AI.",
		"Deep learning uses neural networks.",
		"Natural language processing handles text.",
	}
	
	// Test tokenization for each query-document pair
	tokenizer := client.GetTokenizer()
	
	for i, doc := range documents {
		tokenIDs, attentionMask, tokenTypeIDs := tokenizer.EncodeQueryPassagePair(query, doc, 512)
		
		if len(tokenIDs) != 512 {
			t.Errorf("Document %d: Expected 512 token IDs, got %d", i, len(tokenIDs))
		}
		
		if len(attentionMask) != 512 {
			t.Errorf("Document %d: Expected 512 attention mask values, got %d", i, len(attentionMask))
		}
		
		if len(tokenTypeIDs) != 512 {
			t.Errorf("Document %d: Expected 512 token type IDs, got %d", i, len(tokenTypeIDs))
		}
		
		// Count actual tokens
		actualTokens := 0
		for _, mask := range attentionMask {
			if mask == 1 {
				actualTokens++
			}
		}
		
		t.Logf("Document %d: %d actual tokens", i, actualTokens)
		t.Logf("Document %d: First 10 token IDs: %v", i, tokenIDs[:10])
	}
}

func TestPureGoBGEClient_Configuration(t *testing.T) {
	modelPath := filepath.Join("..", "..", "models", "bge-reranker-large-onnx")
	
	client, err := NewPureGoBGEClient(modelPath, true, "http://localhost:8000")
	if err != nil {
		t.Fatalf("Failed to create pure Go BGE client: %v", err)
	}
	
	// Test max length configuration
	originalMaxLength := client.maxLength
	newMaxLength := 256
	
	client.SetMaxLength(newMaxLength)
	
	if client.maxLength != newMaxLength {
		t.Errorf("Expected max length %d, got %d", newMaxLength, client.maxLength)
	}
	
	// Test tokenization with new max length
	tokenizer := client.GetTokenizer()
	tokenIDs, _, _ := tokenizer.Encode("This is a test sentence", newMaxLength)
	
	if len(tokenIDs) != newMaxLength {
		t.Errorf("Expected %d tokens with new max length, got %d", newMaxLength, len(tokenIDs))
	}
	
	t.Logf("Original max length: %d, New max length: %d", originalMaxLength, newMaxLength)
}

func TestPureGoBGEClient_EdgeCases(t *testing.T) {
	modelPath := filepath.Join("..", "..", "models", "bge-reranker-large-onnx")
	
	client, err := NewPureGoBGEClient(modelPath, true, "http://localhost:8000")
	if err != nil {
		t.Fatalf("Failed to create pure Go BGE client: %v", err)
	}
	
	tokenizer := client.GetTokenizer()
	
	// Test edge cases
	testCases := []struct {
		name    string
		query   string
		passage string
	}{
		{"Empty query", "", "This is a passage"},
		{"Empty passage", "This is a query", ""},
		{"Both empty", "", ""},
		{"Very long query", "This is a very long query that repeats itself. " + 
			"This is a very long query that repeats itself. " +
			"This is a very long query that repeats itself.", "Short passage"},
		{"Special characters", "Query with !@#$%^&*()", "Passage with émoji 🤖"},
		{"Numbers and symbols", "Query 123 + 456 = ?", "Answer: 579"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenIDs, attentionMask, tokenTypeIDs := tokenizer.EncodeQueryPassagePair(tc.query, tc.passage, 128)
			
			if len(tokenIDs) != 128 {
				t.Errorf("Expected 128 token IDs, got %d", len(tokenIDs))
			}
			
			if len(attentionMask) != 128 {
				t.Errorf("Expected 128 attention mask values, got %d", len(attentionMask))
			}
			
			if len(tokenTypeIDs) != 128 {
				t.Errorf("Expected 128 token type IDs, got %d", len(tokenTypeIDs))
			}
			
			actualTokens := 0
			for _, mask := range attentionMask {
				if mask == 1 {
					actualTokens++
				}
			}
			
			t.Logf("Case '%s': %d actual tokens", tc.name, actualTokens)
		})
	}
}

// Benchmark tokenization performance
func BenchmarkPureGoBGEClient_Tokenization(b *testing.B) {
	modelPath := filepath.Join("..", "..", "models", "bge-reranker-large-onnx")
	
	client, err := NewPureGoBGEClient(modelPath, true, "http://localhost:8000")
	if err != nil {
		b.Fatalf("Failed to create pure Go BGE client: %v", err)
	}
	
	tokenizer := client.GetTokenizer()
	query := "What is the best way to implement machine learning algorithms?"
	passage := "Machine learning algorithms can be implemented using various frameworks and libraries. Popular choices include TensorFlow, PyTorch, and scikit-learn for Python, or MLlib for Scala and Java applications."
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_, _, _ = tokenizer.EncodeQueryPassagePair(query, passage, 512)
	}
}

func BenchmarkPureGoBGEClient_BatchTokenization(b *testing.B) {
	modelPath := filepath.Join("..", "..", "models", "bge-reranker-large-onnx")
	
	client, err := NewPureGoBGEClient(modelPath, true, "http://localhost:8000")
	if err != nil {
		b.Fatalf("Failed to create pure Go BGE client: %v", err)
	}
	
	tokenizer := client.GetTokenizer()
	query := "What is machine learning?"
	
	documents := []string{
		"Machine learning is a method of data analysis.",
		"It automates analytical model building.",
		"ML is a branch of artificial intelligence.",
		"It uses algorithms that iteratively learn from data.",
		"ML allows computers to find hidden insights.",
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		for _, doc := range documents {
			_, _, _ = tokenizer.EncodeQueryPassagePair(query, doc, 512)
		}
	}
}