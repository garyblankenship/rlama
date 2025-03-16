package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/service"
)

// Server represents the API server
type Server struct {
	port        string
	ragService  service.RagService
	ollamaClient *client.OllamaClient
}

// NewServer creates a new API server
func NewServer(port string, ollamaClient *client.OllamaClient) *Server {
	if port == "" {
		port = "11249" // Default port
	}
	
	return &Server{
		port:        port,
		ragService:  service.NewRagService(ollamaClient),
		ollamaClient: ollamaClient,
	}
}

// Start starts the API server
func (s *Server) Start() error {
	// Register routes
	http.HandleFunc("/rag", s.handleRagQuery)
	http.HandleFunc("/health", s.handleHealthCheck)
	
	// Start the server
	addr := fmt.Sprintf(":%s", s.port)
	log.Printf("Starting API server on http://localhost%s", addr)
	
	server := &http.Server{
		Addr:         addr,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	
	return server.ListenAndServe()
}

// RagQueryRequest represents the request body for RAG queries
type RagQueryRequest struct {
	RagName       string `json:"rag_name"`
	Model         string `json:"model,omitempty"`
	Prompt        string `json:"prompt"`
	ContextSize   int    `json:"context_size,omitempty"`
	MaxWorkers    int    `json:"max_workers,omitempty"` // Added for parallel processing
}

// RagQueryResponse represents the response for RAG queries
type RagQueryResponse struct {
	Response string `json:"response"`
	Success  bool   `json:"success"`
	Error    string `json:"error,omitempty"`
}

// Handle RAG queries
func (s *Server) handleRagQuery(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		sendErrorResponse(w, "Error reading request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	
	// Parse request
	var req RagQueryRequest
	if err := json.Unmarshal(body, &req); err != nil {
		sendErrorResponse(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}
	
	// Validate request
	if req.RagName == "" {
		sendErrorResponse(w, "Missing 'rag_name' field", http.StatusBadRequest)
		return
	}
	if req.Prompt == "" {
		sendErrorResponse(w, "Missing 'prompt' field", http.StatusBadRequest)
		return
	}
	
	// Set default context size if not provided
	if req.ContextSize <= 0 {
		req.ContextSize = 20
	}
	
	// Load the RAG system
	rag, err := s.ragService.LoadRag(req.RagName)
	if err != nil {
		sendErrorResponse(w, fmt.Sprintf("Error loading RAG: %v", err), http.StatusNotFound)
		return
	}
	
	// If model is specified and different from RAG's model, update it temporarily
	modelToUse := rag.ModelName
	if req.Model != "" && req.Model != rag.ModelName {
		modelToUse = req.Model
		log.Printf("Using specified model %s instead of RAG's model %s", req.Model, rag.ModelName)
	}
	
	// Check if Ollama model is available
	if err := s.ollamaClient.CheckOllamaAndModel(modelToUse); err != nil {
		sendErrorResponse(w, fmt.Sprintf("Error with Ollama model: %v", err), http.StatusInternalServerError)
		return
	}
	
	// Use the original model of the RAG
	originalModel := rag.ModelName
	
	// Temporarily update the model if needed
	if req.Model != "" && req.Model != originalModel {
		rag.ModelName = req.Model
	}
	
	// Set parallel workers if specified
	if req.MaxWorkers > 0 {
		embeddingService := service.NewEmbeddingService(s.ollamaClient)
		embeddingService.SetMaxWorkers(req.MaxWorkers)
		// Update the RAG service with the new embedding service
		s.ragService = service.NewRagServiceWithEmbedding(s.ollamaClient, embeddingService)
	}
	
	// Query the RAG system
	response, err := s.ragService.Query(rag, req.Prompt, req.ContextSize)
	
	// Restore original model
	rag.ModelName = originalModel
	
	if err != nil {
		sendErrorResponse(w, fmt.Sprintf("Error querying RAG: %v", err), http.StatusInternalServerError)
		return
	}
	
	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	resp := RagQueryResponse{
		Response: response,
		Success:  true,
	}
	
	json.NewEncoder(w).Encode(resp)
}

// Handle health check requests
func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]string{
		"status": "healthy",
		"time":   time.Now().Format(time.RFC3339),
	}
	
	json.NewEncoder(w).Encode(response)
}

// Helper function to send error responses
func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	resp := RagQueryResponse{
		Success: false,
		Error:   message,
	}
	
	json.NewEncoder(w).Encode(resp)
} 