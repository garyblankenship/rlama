package vector

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
)

// VectorStoreInterface defines the common interface for vector stores
type VectorStoreInterface interface {
	Add(id string, vector []float32)
	Search(query []float32, limit int) []SearchResult
	Remove(id string)
	Save(path string) error
	Load(path string) error
}

// VectorItem represents an item in the vector storage
type VectorItem struct {
	ID      string    `json:"id"`
	Vector  []float32 `json:"vector"`
}

// SearchResult represents a search result
type SearchResult struct {
	ID       string  `json:"id"`
	Score    float64 `json:"score"`
}

// Store is a simple vector storage with cosine similarity search
type Store struct {
	Items []VectorItem `json:"items"`
}

// Ensure Store implements VectorStoreInterface
var _ VectorStoreInterface = (*Store)(nil)

// NewStore creates a new vector storage
func NewStore() *Store {
	return &Store{
		Items: []VectorItem{},
	}
}

// Add adds a vector to the storage
func (s *Store) Add(id string, vector []float32) {
	// Check if the ID already exists
	for i, item := range s.Items {
		if item.ID == id {
			// Replace the existing vector
			s.Items[i].Vector = vector
			return
		}
	}
	
	// Add a new vector
	s.Items = append(s.Items, VectorItem{
		ID:     id,
		Vector: vector,
	})
}

// Search searches for the most similar vectors
func (s *Store) Search(query []float32, limit int) []SearchResult {
	var results []SearchResult
	
	// Calculate cosine similarity for each vector
	for _, item := range s.Items {
		score := cosineSimilarity(query, item.Vector)
		results = append(results, SearchResult{
			ID:    item.ID,
			Score: score,
		})
	}
	
	// Sort by descending score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	
	// Limit the number of results
	if limit > 0 && limit < len(results) {
		results = results[:limit]
	}
	
	return results
}

// cosineSimilarity calculates the cosine similarity between two vectors
func cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}
	
	var dotProduct float64
	var normA, normB float64
	
	for i := range a {
		dotProduct += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}
	
	if normA == 0 || normB == 0 {
		return 0
	}
	
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// Save saves the vector storage to a file
func (s *Store) Save(path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to serialize vector storage: %w", err)
	}
	
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("unable to save vector storage: %w", err)
	}
	
	return nil
}

// Load loads the vector storage from a file
func (s *Store) Load(path string) error {
	// Check if the file exists
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, use empty storage
			s.Items = []VectorItem{}
			return nil
		}
		return fmt.Errorf("unable to access vector storage file: %w", err)
	}
	
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("unable to read vector storage file: %w", err)
	}
	
	err = json.Unmarshal(data, s)
	if err != nil {
		return fmt.Errorf("unable to deserialize vector storage: %w", err)
	}
	
	return nil
}

// Remove removes a vector from the storage by its ID
func (s *Store) Remove(id string) {
	for i, item := range s.Items {
		if item.ID == id {
			s.Items = append(s.Items[:i], s.Items[i+1:]...)
			return
		}
	}
} 