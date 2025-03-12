package vector

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/blevesearch/bleve/v2"
)

// DocumentData represents the data structure for Bleve indexing
type DocumentData struct {
	ID       string `json:"id"`
	Content  string `json:"content"`
	Metadata string `json:"metadata"`
}

// EnhancedHybridStore combines HNSW vector search and BM25 text search
type EnhancedHybridStore struct {
	VectorStore VectorStoreInterface `json:"-"`
	TextIndex   bleve.Index          `json:"-"`
	WeightBM25  float64              `json:"weight_bm25"`
	// Maps for quick access to content and metadata
	contentCache  map[string]string `json:"-"`
	metadataCache map[string]string `json:"-"`
}

// Ensure EnhancedHybridStore implements VectorStoreInterface
var _ VectorStoreInterface = (*EnhancedHybridStore)(nil)

// NewEnhancedHybridStore creates a new enhanced hybrid store
func NewEnhancedHybridStore(indexPath string, dimensions int) (*EnhancedHybridStore, error) {
	// Create index directory if needed
	if indexPath != "" && indexPath != ":memory:" {
		err := os.MkdirAll(filepath.Dir(indexPath), 0755)
		if err != nil {
			return nil, fmt.Errorf("unable to create index directory: %w", err)
		}
	}

	// Create or open Bleve index
	var textIndex bleve.Index
	var err error

	if indexPath == "" || indexPath == ":memory:" {
		// In-memory index
		indexMapping := bleve.NewIndexMapping()
		textIndex, err = bleve.NewMemOnly(indexMapping)
	} else {
		// Check if index already exists
		_, err := os.Stat(indexPath)
		if os.IsNotExist(err) {
			// Create new index
			indexMapping := bleve.NewIndexMapping()
			textIndex, err = bleve.New(indexPath, indexMapping)
		} else {
			// Open existing index
			textIndex, err = bleve.Open(indexPath)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("error creating/opening Bleve index: %w", err)
	}

	return &EnhancedHybridStore{
		VectorStore:   NewHNSWStore(dimensions),
		TextIndex:     textIndex,
		WeightBM25:    0.3, // 30% BM25, 70% vector by default
		contentCache:  make(map[string]string),
		metadataCache: make(map[string]string),
	}, nil
}

// AddDocument adds a document to both the vector and text indexes
func (hs *EnhancedHybridStore) AddDocument(id string, content string, metadata string, vector []float32) error {
	// Add to vector store
	hs.VectorStore.Add(id, vector)

	// Add to cache
	hs.contentCache[id] = content
	hs.metadataCache[id] = metadata

	// Add to text index
	doc := DocumentData{
		ID:       id,
		Content:  content,
		Metadata: metadata,
	}
	
	err := hs.TextIndex.Index(id, doc)
	if err != nil {
		return fmt.Errorf("error indexing text: %w", err)
	}
	
	return nil
}

// Add implements the VectorStoreInterface
func (hs *EnhancedHybridStore) Add(id string, vector []float32) {
	hs.VectorStore.Add(id, vector)
}

// Remove removes a document from both indexes
func (hs *EnhancedHybridStore) Remove(id string) {
	// Remove from vector store
	hs.VectorStore.Remove(id)
	
	// Remove from caches
	delete(hs.contentCache, id)
	delete(hs.metadataCache, id)
	
	// Remove from text index (ignore errors for interface compatibility)
	hs.TextIndex.Delete(id)
}

// GetContent returns a document's content
func (hs *EnhancedHybridStore) GetContent(id string) string {
	return hs.contentCache[id]
}

// GetMetadata returns a document's metadata
func (hs *EnhancedHybridStore) GetMetadata(id string) string {
	return hs.metadataCache[id]
}

// HybridSearchResult représente un résultat de recherche hybride
type HybridSearchResult struct {
	ID             string  `json:"id"`
	VectorScore    float64 `json:"vector_score"`
	TextScore      float64 `json:"text_score"`
	CombinedScore  float64 `json:"combined_score"`
}

// HybridSearch performs a combined vector and text search
func (hs *EnhancedHybridStore) HybridSearch(queryVector []float32, queryText string, limit int) ([]HybridSearchResult, error) {
	// Execute vector search
	vectorResults := hs.VectorStore.Search(queryVector, limit*2) // Get more results for fusion
	
	// Execute BM25 text search
	textQuery := bleve.NewQueryStringQuery(queryText)
	textSearch := bleve.NewSearchRequest(textQuery)
	textSearch.Size = limit * 2
	textSearchResults, err := hs.TextIndex.Search(textSearch)
	if err != nil {
		return nil, fmt.Errorf("error during text search: %w", err)
	}
	
	// Store normalized scores in maps
	vectorScores := make(map[string]float64)
	textScores := make(map[string]float64)
	allIDs := make(map[string]bool)
	
	// Normalize vector scores
	maxVectorScore := 0.0
	for _, res := range vectorResults {
		if res.Score > maxVectorScore {
			maxVectorScore = res.Score
		}
	}
	
	for _, res := range vectorResults {
		normalizedScore := res.Score
		if maxVectorScore > 0 {
			normalizedScore = res.Score / maxVectorScore
		}
		vectorScores[res.ID] = normalizedScore
		allIDs[res.ID] = true
	}
	
	// Normalize text scores
	maxTextScore := 0.0
	for _, hit := range textSearchResults.Hits {
		if hit.Score > maxTextScore {
			maxTextScore = hit.Score
		}
	}
	
	for _, hit := range textSearchResults.Hits {
		normalizedScore := hit.Score
		if maxTextScore > 0 {
			normalizedScore = hit.Score / maxTextScore
		}
		textScores[hit.ID] = normalizedScore
		allIDs[hit.ID] = true
	}
	
	// Combine scores with weighting
	var hybridResults []HybridSearchResult
	for id := range allIDs {
		vectorScore := vectorScores[id]
		textScore := textScores[id]
		
		// If a document is only in one result set, give it a minimum score in the other
		if _, exists := vectorScores[id]; !exists {
			vectorScore = 0.01 // Minimum score to not completely eliminate
		}
		if _, exists := textScores[id]; !exists {
			textScore = 0.01 // Minimum score to not completely eliminate
		}
		
		// Weighted combined score
		combinedScore := (hs.WeightBM25 * textScore) + ((1 - hs.WeightBM25) * vectorScore)
		
		hybridResults = append(hybridResults, HybridSearchResult{
			ID:            id,
			VectorScore:   vectorScore,
			TextScore:     textScore,
			CombinedScore: combinedScore,
		})
	}
	
	// Sort by combined score in descending order
	SortHybridResults(hybridResults)
	
	// Limit results
	if len(hybridResults) > limit {
		hybridResults = hybridResults[:limit]
	}
	
	return hybridResults, nil
}

// Search implements the basic vector search interface
func (hs *EnhancedHybridStore) Search(query []float32, limit int) []SearchResult {
	return hs.VectorStore.Search(query, limit)
}

// Save saves both indexes
func (hs *EnhancedHybridStore) Save(vectorPath string) error {
	// Save vector store
	err := hs.VectorStore.Save(vectorPath)
	if err != nil {
		return fmt.Errorf("error saving vector store: %w", err)
	}
	
	// Bleve index is saved automatically if on disk
	
	return nil
}

// Load loads the store from a file
func (hs *EnhancedHybridStore) Load(path string) error {
	// Load vector store
	err := hs.VectorStore.Load(path)
	if err != nil {
		return fmt.Errorf("error loading vector store: %w", err)
	}
	
	// Bleve index is managed separately
	
	return nil
}

// Close properly closes the indexes
func (hs *EnhancedHybridStore) Close() error {
	return hs.TextIndex.Close()
}

// SortHybridResults trie les résultats par score combiné décroissant
func SortHybridResults(results []HybridSearchResult) {
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[j].CombinedScore > results[i].CombinedScore {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
} 