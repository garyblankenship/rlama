package domain

import (
	"time"

	"github.com/dontizi/rlama/pkg/vector"
)

// RagSystem represents a complete RAG system
type RagSystem struct {
	Name        string                      `json:"name"`
	ModelName   string                      `json:"model_name"`
	CreatedAt   time.Time                   `json:"created_at"`
	UpdatedAt   time.Time                   `json:"updated_at"`
	Description string                      `json:"description"`
	HybridStore *vector.EnhancedHybridStore // Use the hybrid store
	Documents   []*Document                 `json:"documents"`
	Chunks      []*DocumentChunk            `json:"chunks"`
	// Directory watching settings
	WatchedDir    string               `json:"watched_dir,omitempty"`
	WatchInterval int                  `json:"watch_interval,omitempty"` // In minutes, 0 means only check on use
	LastWatchedAt time.Time            `json:"last_watched_at,omitempty"`
	WatchEnabled  bool                 `json:"watch_enabled"`
	WatchOptions  DocumentWatchOptions `json:"watch_options,omitempty"`
	// Web watching settings
	WatchedURL       string          `json:"watched_url,omitempty"`
	WebWatchEnabled  bool            `json:"web_watch_enabled"`
	WebWatchInterval int             `json:"web_watch_interval,omitempty"` // In minutes
	LastWebWatchAt   time.Time       `json:"last_web_watched_at,omitempty"`
	WebWatchOptions  WebWatchOptions `json:"web_watch_options,omitempty"`
	APIProfileName   string          `json:"api_profile_name,omitempty"`  // Name of the API profile to use
	ChunkingStrategy string          `json:"chunking_strategy,omitempty"` // Type of chunking strategy used
	// Reranking settings
	RerankerEnabled   bool    `json:"reranker_enabled,omitempty"`   // Whether to use reranking
	RerankerModel     string  `json:"reranker_model,omitempty"`     // Model to use for reranking (if different from ModelName)
	RerankerWeight    float64 `json:"reranker_weight,omitempty"`    // Weight for reranker scores vs vector scores (0-1)
	RerankerThreshold float64 `json:"reranker_threshold,omitempty"` // Minimum score threshold for reranked results
	RerankerTopK      int     `json:"reranker_top_k,omitempty"`     // Default: return only top 5 results after reranking
	RerankerSilent    bool    `json:"reranker_silent,omitempty"`    // Whether to suppress warnings and output from the reranker
	// Embedding settings
	EmbeddingDimension int `json:"embedding_dimension,omitempty"` // Dimension of the embedding vectors
	
	// Vector Store Configuration
	VectorStoreType      string `json:"vector_store_type,omitempty"`      // e.g., "internal_hnsw", "qdrant"
	QdrantHost           string `json:"qdrant_host,omitempty"`
	QdrantPort           int    `json:"qdrant_port,omitempty"`            // e.g., 6333 for HTTP, 6334 for gRPC
	QdrantAPIKey         string `json:"qdrant_api_key,omitempty"`         // For Qdrant Cloud or secured instances
	QdrantCollectionName string `json:"qdrant_collection_name,omitempty"` // Typically derived from ragName
	QdrantGRPC           bool   `json:"qdrant_grpc,omitempty"`            // True to use gRPC, false for HTTP REST
}

// DocumentWatchOptions stores settings for directory watching
type DocumentWatchOptions struct {
	ExcludeDirs      []string `json:"exclude_dirs,omitempty"`
	ExcludeExts      []string `json:"exclude_exts,omitempty"`
	ProcessExts      []string `json:"process_exts,omitempty"`
	ChunkSize        int      `json:"chunk_size,omitempty"`
	ChunkOverlap     int      `json:"chunk_overlap,omitempty"`
	ChunkingStrategy string   `json:"chunking_strategy,omitempty"`
}

// WebWatchOptions stores settings for web watching
type WebWatchOptions struct {
	MaxDepth         int      `json:"max_depth,omitempty"`
	Concurrency      int      `json:"concurrency,omitempty"`
	ExcludePaths     []string `json:"exclude_paths,omitempty"`
	ChunkSize        int      `json:"chunk_size,omitempty"`
	ChunkOverlap     int      `json:"chunk_overlap,omitempty"`
	ChunkingStrategy string   `json:"chunking_strategy,omitempty"`
}

// NewRagSystem creates a new instance of RagSystem
func NewRagSystem(name, modelName string) *RagSystem {
	return NewRagSystemWithDimensions(name, modelName, 1536) // Default to 1536 dimensions
}

// NewRagSystemWithDimensions creates a new instance of RagSystem with specified embedding dimensions
func NewRagSystemWithDimensions(name, modelName string, dimensions int) *RagSystem {
	now := time.Now()
	hybridStore, err := vector.NewEnhancedHybridStore(":memory:", dimensions)
	if err != nil {
		// Handle error appropriately
		return nil
	}

	return &RagSystem{
		Name:               name,
		ModelName:          modelName,
		CreatedAt:          now,
		UpdatedAt:          now,
		HybridStore:        hybridStore,
		Documents:          []*Document{},
		Chunks:             []*DocumentChunk{},
		RerankerEnabled:    true,                      // Enable reranking by default
		RerankerModel:      "BAAI/bge-reranker-v2-m3", // Use BGE reranker by default
		RerankerWeight:     0.7,                       // Default: 70% reranker score, 30% vector similarity
		RerankerTopK:       5,                         // Default: return only top 5 results after reranking
		EmbeddingDimension: dimensions,                // Store the embedding dimension
		VectorStoreType:    "internal",                // Default to internal vector store
	}
}

// NewRagSystemWithVectorStore creates a new instance of RagSystem with vector store configuration
func NewRagSystemWithVectorStore(name, modelName string, dimensions int, vectorStoreType, qdrantHost string, qdrantPort int, qdrantAPIKey, qdrantCollection string, qdrantGRPC bool) *RagSystem {
	now := time.Now()
	
	// Create hybrid store config
	hybridConfig := vector.HybridStoreConfig{
		IndexPath:            ":memory:",
		Dimensions:           dimensions,
		VectorStoreType:      vectorStoreType,
		QdrantHost:           qdrantHost,
		QdrantPort:           qdrantPort,
		QdrantAPIKey:         qdrantAPIKey,
		QdrantCollectionName: qdrantCollection,
		QdrantGRPC:           qdrantGRPC,
	}
	
	hybridStore, err := vector.NewEnhancedHybridStoreWithConfig(hybridConfig)
	if err != nil {
		// Handle error appropriately
		return nil
	}

	return &RagSystem{
		Name:                 name,
		ModelName:            modelName,
		CreatedAt:            now,
		UpdatedAt:            now,
		HybridStore:          hybridStore,
		Documents:            []*Document{},
		Chunks:               []*DocumentChunk{},
		RerankerEnabled:      true,                      // Enable reranking by default
		RerankerModel:        "BAAI/bge-reranker-v2-m3", // Use BGE reranker by default
		RerankerWeight:       0.7,                       // Default: 70% reranker score, 30% vector similarity
		RerankerTopK:         5,                         // Default: return only top 5 results after reranking
		EmbeddingDimension:   dimensions,                // Store the embedding dimension
		VectorStoreType:      vectorStoreType,
		QdrantHost:           qdrantHost,
		QdrantPort:           qdrantPort,
		QdrantAPIKey:         qdrantAPIKey,
		QdrantCollectionName: qdrantCollection,
		QdrantGRPC:           qdrantGRPC,
	}
}

// AddDocument adds a document to the RAG system
func (r *RagSystem) AddDocument(doc *Document) {
	r.Documents = append(r.Documents, doc)
	if doc.Embedding != nil {
		// Don't use doc.Metadata if it doesn't exist
		r.HybridStore.Add(doc.ID, doc.Embedding)
	}
	r.UpdatedAt = time.Now()
}

// GetDocumentByID retrieves a document by its ID
func (r *RagSystem) GetDocumentByID(id string) *Document {
	for _, doc := range r.Documents {
		if doc.ID == id {
			return doc
		}
	}
	return nil
}

// RemoveDocument removes a document from the RAG system by its ID
func (r *RagSystem) RemoveDocument(id string) bool {
	// Find the document index
	var index = -1
	for i, doc := range r.Documents {
		if doc.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		return false
	}

	// Remove from the Documents slice
	r.Documents = append(r.Documents[:index], r.Documents[index+1:]...)

	// Remove from the HybridStore
	r.HybridStore.Remove(id)

	r.UpdatedAt = time.Now()
	return true
}

// AddChunk adds a chunk to the RAG system
func (r *RagSystem) AddChunk(chunk *DocumentChunk) {
	r.Chunks = append(r.Chunks, chunk)
	if chunk.Embedding != nil {
		r.HybridStore.Add(chunk.ID, chunk.Embedding)
	}
	r.UpdatedAt = time.Now()
}

// GetChunkByID retrieves a chunk by its ID
func (r *RagSystem) GetChunkByID(id string) *DocumentChunk {
	for _, chunk := range r.Chunks {
		if chunk.ID == id {
			return chunk
		}
	}
	return nil
}

// Search performs a hybrid search using the hybrid store
func (r *RagSystem) Search(queryVector []float32, queryText string, limit int) ([]vector.HybridSearchResult, error) {
	return r.HybridStore.HybridSearch(queryVector, queryText, limit)
}
