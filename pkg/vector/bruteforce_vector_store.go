package vector

import (
	"encoding/gob"
	"fmt"
	"math"
	"os"
	"sort"
)

// Define key type for vector IDs
type vectorID string

// BruteForceVectorStore implements a vector store using brute-force linear search
// This provides a simple, straightforward implementation without any indexing optimizations
type BruteForceVectorStore struct {
	items map[string][]float32
	dims  int
}

// Ensure BruteForceVectorStore implements VectorStoreInterface
var _ VectorStoreInterface = (*BruteForceVectorStore)(nil)

// NewBruteForceVectorStore creates a new vector store
func NewBruteForceVectorStore(dimensions int) *BruteForceVectorStore {
	return &BruteForceVectorStore{
		items: make(map[string][]float32),
		dims:  dimensions,
	}
}

// Add adds a vector to the store
func (s *BruteForceVectorStore) Add(id string, vector []float32) {
	// Store vector in items map
	s.items[id] = vector
}

// Remove removes a vector from the store
func (s *BruteForceVectorStore) Remove(id string) {
	// Remove from items map
	delete(s.items, id)
}

// computeCosineSimilarity calculates cosine similarity between two vectors
func computeCosineSimilarity(a, b []float32) float64 {
	// Check for empty vectors to prevent index out of range errors
	if len(a) == 0 || len(b) == 0 {
		return 0.0
	}
	
	// Check for length mismatch
	if len(a) != len(b) {
		// Log the error but return a default value instead of panicking
		fmt.Printf("Warning: Vector length mismatch (%d vs %d), cannot compute similarity\n", len(a), len(b))
		return 0.0
	}

	var dotProduct float64
	var normA float64
	var normB float64

	for i := 0; i < len(a); i++ {
		dotProduct += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}

	// Handle the case where one of the norms is zero
	if normA == 0 || normB == 0 {
		return 0.0
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// Search searches for similar vectors using cosine similarity
func (s *BruteForceVectorStore) Search(query []float32, limit int) []SearchResult {
	results := make([]SearchResult, 0, len(s.items))
	
	// Compute similarity for all vectors
	for id, vector := range s.items {
		similarity := computeCosineSimilarity(query, vector)
		results = append(results, SearchResult{
			ID:    id,
			Score: similarity,
		})
	}
	
	// Sort by similarity score in descending order
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	
	// Limit results
	if len(results) > limit {
		results = results[:limit]
	}
	
	return results
}

// Save saves the vector store to disk
func (s *BruteForceVectorStore) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()
	
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(s.items)
	if err != nil {
		return fmt.Errorf("failed to encode vectors: %w", err)
	}
	
	return nil
}

// Load loads the vector store from disk
func (s *BruteForceVectorStore) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, use empty storage
			return nil
		}
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()
	
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&s.items)
	if err != nil {
		return fmt.Errorf("failed to decode vectors: %w", err)
	}
	
	return nil
} 