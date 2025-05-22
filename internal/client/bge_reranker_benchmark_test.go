package client

import (
	"os"
	"path/filepath"
	"testing"
)

func BenchmarkBGEReranker(b *testing.B) {
	// Test data
	testPairs := [][]string{
		{"machine learning", "Machine learning is a subset of artificial intelligence that involves training algorithms to make predictions"},
		{"machine learning", "I like to eat pizza on weekends with my family"},
		{"natural language processing", "NLP is a field of AI focused on enabling computers to understand human language"},
		{"natural language processing", "The weather forecast shows rain tomorrow afternoon"},
		{"computer vision", "Computer vision enables machines to interpret and understand visual information from images"},
		{"computer vision", "My favorite hobby is gardening in the spring season"},
	}

	modelDir := filepath.Join("..", "..", "models", "bge-reranker-large-onnx")
	if _, err := os.Stat(modelDir); os.IsNotExist(err) {
		b.Skip("Model directory not found, skipping benchmark")
	}

	b.Run("Original_Python", func(b *testing.B) {
		client := NewBGERerankerClientWithOptions("BAAI/bge-reranker-v2-m3", true)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := client.ComputeScores(testPairs, true)
			if err != nil {
				b.Fatalf("Failed to compute scores: %v", err)
			}
		}
	})

	b.Run("ONNX_Microservice", func(b *testing.B) {
		client, err := NewBGEONNXRerankerClient(modelDir)
		if err != nil {
			b.Fatalf("Failed to create ONNX client: %v", err)
		}
		defer client.Cleanup()
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := client.ComputeScores(testPairs, true)
			if err != nil {
				b.Fatalf("Failed to compute scores: %v", err)
			}
		}
	})
}

func BenchmarkBGEReranker_SinglePair(b *testing.B) {
	testPair := [][]string{
		{"What is machine learning?", "Machine learning is a subset of artificial intelligence"},
	}

	modelDir := filepath.Join("..", "..", "models", "bge-reranker-large-onnx")
	if _, err := os.Stat(modelDir); os.IsNotExist(err) {
		b.Skip("Model directory not found, skipping benchmark")
	}

	b.Run("Original_Python", func(b *testing.B) {
		client := NewBGERerankerClientWithOptions("BAAI/bge-reranker-v2-m3", true)
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := client.ComputeScores(testPair, true)
			if err != nil {
				b.Fatalf("Failed to compute scores: %v", err)
			}
		}
	})

	b.Run("ONNX_Microservice", func(b *testing.B) {
		client, err := NewBGEONNXRerankerClient(modelDir)
		if err != nil {
			b.Fatalf("Failed to create ONNX client: %v", err)
		}
		defer client.Cleanup()
		
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := client.ComputeScores(testPair, true)
			if err != nil {
				b.Fatalf("Failed to compute scores: %v", err)
			}
		}
	})
}