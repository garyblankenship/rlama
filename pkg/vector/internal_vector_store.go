package vector

import (
	"encoding/gob"
	"fmt"
	"math"
	"os"
	"sort"
	"sync"
)

// InternalVectorStore implements VectorStoreInterface using a thread-safe brute-force search.
// This provides an internal, memory-based vector store with optimized linear search.
type InternalVectorStore struct {
	vectors map[string][]float32
	dims    int
	mutex   sync.RWMutex
}

// Ensure InternalVectorStore implements VectorStoreInterface
var _ VectorStoreInterface = (*InternalVectorStore)(nil)

// NewInternalVectorStore creates a new internal vector store.
func NewInternalVectorStore(dimensions int) *InternalVectorStore {
	return &InternalVectorStore{
		vectors: make(map[string][]float32),
		dims:    dimensions,
	}
}

// Add inserts or updates a vector in the store.
func (s *InternalVectorStore) Add(id string, vector []float32) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(vector) != s.dims {
		fmt.Printf("Error: vector dimension %d does not match store dimension %d for ID %s\n", len(vector), s.dims, id)
		return
	}

	// Store the vector
	s.vectors[id] = make([]float32, len(vector))
	copy(s.vectors[id], vector)
}

// computeCosineSimilarityOptimized calculates cosine similarity between two vectors
func computeCosineSimilarityOptimized(a, b []float32) float64 {
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

// Search finds the k-nearest neighbors for a query vector.
func (s *InternalVectorStore) Search(query []float32, limit int) []SearchResult {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if len(s.vectors) == 0 {
		return []SearchResult{}
	}

	results := make([]SearchResult, 0, len(s.vectors))

	// Compute similarity for all vectors
	for id, vector := range s.vectors {
		similarity := computeCosineSimilarityOptimized(query, vector)
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

// Remove removes a vector from the store.
func (s *InternalVectorStore) Remove(id string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.vectors, id)
}

// Save persists the vector store to disk.
func (s *InternalVectorStore) Save(path string) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	saveData := struct {
		Vectors map[string][]float32
		Dims    int
	}{
		Vectors: s.vectors,
		Dims:    s.dims,
	}

	err = encoder.Encode(saveData)
	if err != nil {
		return fmt.Errorf("failed to encode vectors: %w", err)
	}

	return nil
}

// Load reconstructs the vector store from disk.
func (s *InternalVectorStore) Load(path string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// File doesn't exist, initialize empty store
		s.vectors = make(map[string][]float32)
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	var saveData struct {
		Vectors map[string][]float32
		Dims    int
	}

	err = decoder.Decode(&saveData)
	if err != nil {
		return fmt.Errorf("failed to decode vectors: %w", err)
	}

	s.vectors = saveData.Vectors
	if saveData.Dims != 0 && saveData.Dims != s.dims {
		return fmt.Errorf("loaded dimension %d does not match expected dimension %d", saveData.Dims, s.dims)
	}

	return nil
}