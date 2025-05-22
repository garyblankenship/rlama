package service

import (
	"fmt"
	"os"
	"strconv"
)

// ServiceConfig holds all configuration needed for service creation
type ServiceConfig struct {
	// Client Configuration
	OllamaHost     string
	OllamaPort     string
	OpenAIAPIKey   string
	DataDirectory  string
	
	// Profile Configuration
	APIProfileName string
	
	// Document Processing Configuration
	ChunkSize        int
	ChunkOverlap     int
	ChunkingStrategy string
	
	// Embedding Configuration
	EmbeddingModel string
	
	// Reranking Configuration
	RerankerEnabled   bool
	RerankerModel     string
	RerankerWeight    float64
	RerankerThreshold float64
	UseONNXReranker   bool
	ONNXModelDir      string
	
	// Vector Store Configuration
	VectorStoreType      string
	QdrantHost           string
	QdrantPort           int
	QdrantAPIKey         string
	QdrantCollectionName string
	QdrantGRPC           bool
	
	// Debugging and Logging
	Verbose bool
}

// NewServiceConfig creates a new service configuration with defaults
func NewServiceConfig() *ServiceConfig {
	return &ServiceConfig{
		// Client defaults
		OllamaHost:     getEnvWithDefault("OLLAMA_HOST", "localhost"),
		OllamaPort:     getEnvWithDefault("OLLAMA_PORT", "11434"),
		OpenAIAPIKey:   os.Getenv("OPENAI_API_KEY"),
		DataDirectory:  getEnvWithDefault("RLAMA_DATA_DIR", ""),
		
		// Document processing defaults
		ChunkSize:        1000,
		ChunkOverlap:     200,
		ChunkingStrategy: "hybrid",
		
		// Reranking defaults
		RerankerEnabled: true,
		RerankerWeight:  0.7,
		UseONNXReranker: false,
		ONNXModelDir:    getEnvWithDefault("RLAMA_ONNX_MODEL_DIR", "./models/bge-reranker-large-onnx"),
		
		// Vector store defaults
		VectorStoreType: "internal",
		QdrantHost:      "localhost",
		QdrantPort:      6333,
		QdrantGRPC:      false,
	}
}

// Validate checks that the configuration is valid
func (sc *ServiceConfig) Validate() error {
	// Validate Ollama configuration
	if sc.OllamaHost == "" {
		return fmt.Errorf("Ollama host cannot be empty")
	}
	
	if sc.OllamaPort == "" {
		return fmt.Errorf("Ollama port cannot be empty")
	}
	
	// Validate port is numeric
	if _, err := strconv.Atoi(sc.OllamaPort); err != nil {
		return fmt.Errorf("invalid Ollama port: %w", err)
	}
	
	// Validate chunk configuration
	if sc.ChunkSize <= 0 {
		return fmt.Errorf("chunk size must be positive")
	}
	
	if sc.ChunkOverlap < 0 {
		return fmt.Errorf("chunk overlap cannot be negative")
	}
	
	if sc.ChunkOverlap >= sc.ChunkSize {
		return fmt.Errorf("chunk overlap must be less than chunk size")
	}
	
	// Validate reranker weight
	if sc.RerankerWeight < 0 || sc.RerankerWeight > 1 {
		return fmt.Errorf("reranker weight must be between 0 and 1")
	}
	
	// Validate vector store configuration
	if sc.VectorStoreType != "internal" && sc.VectorStoreType != "qdrant" {
		return fmt.Errorf("vector store type must be 'internal' or 'qdrant'")
	}
	
	if sc.VectorStoreType == "qdrant" {
		if sc.QdrantHost == "" {
			return fmt.Errorf("Qdrant host cannot be empty when using Qdrant vector store")
		}
		if sc.QdrantPort <= 0 || sc.QdrantPort > 65535 {
			return fmt.Errorf("Qdrant port must be between 1 and 65535")
		}
	}
	
	return nil
}

// GetOllamaURL returns the full Ollama URL
func (sc *ServiceConfig) GetOllamaURL() string {
	return fmt.Sprintf("http://%s:%s", sc.OllamaHost, sc.OllamaPort)
}

// ToDocumentLoaderOptions converts the config to DocumentLoaderOptions
func (sc *ServiceConfig) ToDocumentLoaderOptions() DocumentLoaderOptions {
	return DocumentLoaderOptions{
		ChunkSize:            sc.ChunkSize,
		ChunkOverlap:         sc.ChunkOverlap,
		ChunkingStrategy:     sc.ChunkingStrategy,
		APIProfileName:       sc.APIProfileName,
		EmbeddingModel:       sc.EmbeddingModel,
		EnableReranker:       sc.RerankerEnabled,
		RerankerModel:        sc.RerankerModel,
		RerankerWeight:       sc.RerankerWeight,
		UseONNXReranker:      sc.UseONNXReranker,
		ONNXModelDir:         sc.ONNXModelDir,
		VectorStore:          sc.VectorStoreType,
		QdrantHost:           sc.QdrantHost,
		QdrantPort:           sc.QdrantPort,
		QdrantAPIKey:         sc.QdrantAPIKey,
		QdrantCollectionName: sc.QdrantCollectionName,
		QdrantGRPC:           sc.QdrantGRPC,
	}
}

// Clone creates a copy of the configuration
func (sc *ServiceConfig) Clone() *ServiceConfig {
	clone := *sc
	return &clone
}

// WithProfile returns a copy of the config with the specified profile
func (sc *ServiceConfig) WithProfile(profileName string) *ServiceConfig {
	clone := sc.Clone()
	clone.APIProfileName = profileName
	return clone
}

// WithVectorStore returns a copy of the config with the specified vector store settings
func (sc *ServiceConfig) WithVectorStore(storeType string, host string, port int, apiKey string) *ServiceConfig {
	clone := sc.Clone()
	clone.VectorStoreType = storeType
	if storeType == "qdrant" {
		clone.QdrantHost = host
		clone.QdrantPort = port
		clone.QdrantAPIKey = apiKey
	}
	return clone
}

// getEnvWithDefault gets an environment variable with a default value
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}