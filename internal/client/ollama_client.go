package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// Default connection settings for Ollama
const (
	DefaultOllamaHost = "localhost"
	DefaultOllamaPort = "11434"
)

// OllamaClient est un client pour l'API Ollama
type OllamaClient struct {
	BaseURL string
	Client  *http.Client
}

// EmbeddingRequest est la structure de la requête pour l'API /api/embeddings
type EmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// EmbeddingResponse est la structure de la réponse de l'API /api/embeddings
type EmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

// GenerationRequest est la structure de la requête pour l'API /api/generate
type GenerationRequest struct {
	Model    string   `json:"model"`
	Prompt   string   `json:"prompt"`
	Context  []int    `json:"context,omitempty"`
	Options  Options  `json:"options,omitempty"`
	Format   string   `json:"format,omitempty"`
	Template string   `json:"template,omitempty"`
	Stream   bool     `json:"stream"`
}

// Options pour l'API generate
type Options struct {
	Temperature float64 `json:"temperature,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
	TopK        int     `json:"top_k,omitempty"`
	NumPredict  int     `json:"num_predict,omitempty"`
}

// GenerationResponse est la structure de la réponse de l'API /api/generate
type GenerationResponse struct {
	Model     string `json:"model"`
	Response  string `json:"response"`
	Context   []int  `json:"context"`
	CreatedAt string `json:"created_at"`
	Done      bool   `json:"done"`
}

// NewOllamaClient crée un nouveau client Ollama
// Si host ou port sont vides, les valeurs par défaut sont utilisées
// Si OLLAMA_HOST est défini, il est utilisé comme valeur par défaut
func NewOllamaClient(host, port string) *OllamaClient {
	// Check for OLLAMA_HOST environment variable
	ollamaHostEnv := os.Getenv("OLLAMA_HOST")
	
	// Default values
	defaultHost := DefaultOllamaHost
	defaultPort := DefaultOllamaPort
	
	// If OLLAMA_HOST is set, parse it
	if ollamaHostEnv != "" {
		// OLLAMA_HOST could be in the form "host:port" or just "host"
		parts := strings.Split(ollamaHostEnv, ":")
		if len(parts) >= 1 {
			defaultHost = parts[0]
		}
		if len(parts) >= 2 {
			defaultPort = parts[1]
		}
	}
	
	// Command flags override environment variables
	if host == "" {
		host = defaultHost
	}
	if port == "" {
		port = defaultPort
	}
	
	baseURL := fmt.Sprintf("http://%s:%s", host, port)
	
	return &OllamaClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

// NewDefaultOllamaClient crée un nouveau client Ollama avec les valeurs par défaut
// Gardé pour compatibilité avec le code existant
func NewDefaultOllamaClient() *OllamaClient {
	return NewOllamaClient(DefaultOllamaHost, DefaultOllamaPort)
}

// GenerateEmbedding génère un embedding pour le texte donné
func (c *OllamaClient) GenerateEmbedding(model, text string) ([]float32, error) {
	reqBody := EmbeddingRequest{
		Model:  model,
		Prompt: text,
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Post(
		fmt.Sprintf("%s/api/embeddings", c.BaseURL),
		"application/json",
		bytes.NewBuffer(reqJSON),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to generate embedding: %s (status: %d)", string(bodyBytes), resp.StatusCode)
	}

	var embeddingResp EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embeddingResp); err != nil {
		return nil, err
	}

	return embeddingResp.Embedding, nil
}

// GenerateCompletion génère une réponse pour le prompt donné
func (c *OllamaClient) GenerateCompletion(model, prompt string) (string, error) {
	reqBody := GenerationRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
		Options: Options{
			Temperature: 0.7,
			TopP:        0.9,
			NumPredict:  1024,
		},
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := c.Client.Post(
		fmt.Sprintf("%s/api/generate", c.BaseURL),
		"application/json",
		bytes.NewBuffer(reqJSON),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to generate completion: %s (status: %d)", string(bodyBytes), resp.StatusCode)
	}

	var genResp GenerationResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		return "", err
	}

	return genResp.Response, nil
}

// IsOllamaRunning checks if Ollama is installed and running
func (c *OllamaClient) IsOllamaRunning() (bool, error) {
	resp, err := c.Client.Get(fmt.Sprintf("%s/api/version", c.BaseURL))
	if err != nil {
		return false, fmt.Errorf("Ollama is not accessible: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Ollama responded with error code: %d", resp.StatusCode)
	}
	
	return true, nil
}

// CheckOllamaAndModel verifies if Ollama is running and if the specified model is available
func (c *OllamaClient) CheckOllamaAndModel(modelName string) error {
	// Check if Ollama is running
	running, err := c.IsOllamaRunning()
	if err != nil {
		return fmt.Errorf("⚠️ Ollama is not installed or not running.\n"+
			"RLAMA requires Ollama to function.\n"+
			"Please install Ollama with: curl -fsSL https://ollama.com/install.sh | sh\n"+
			"Then start it before using RLAMA.")
	}
	
	if !running {
		return fmt.Errorf("⚠️ Ollama is not running.\n"+
			"Please start Ollama before using RLAMA.")
	}
	
	// Check if model is available (optional)
	// This check could be added here
	
	return nil
}

// RunHuggingFaceModel prepares a Hugging Face model for use with Ollama
func (c *OllamaClient) RunHuggingFaceModel(hfModelPath string, quantization string) error {
	modelRef := "hf.co/" + hfModelPath
	if quantization != "" {
		modelRef += ":" + quantization
	}
	
	cmd := exec.Command("ollama", "run", modelRef)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	return cmd.Run()
}

// PullHuggingFaceModel pulls a Hugging Face model into Ollama without running it
func (c *OllamaClient) PullHuggingFaceModel(hfModelPath string, quantization string) error {
	modelRef := "hf.co/" + hfModelPath
	if quantization != "" {
		modelRef += ":" + quantization
	}
	
	cmd := exec.Command("ollama", "pull", modelRef)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

// IsHuggingFaceModel checks if a model name is a Hugging Face model reference
func IsHuggingFaceModel(modelName string) bool {
	return strings.HasPrefix(modelName, "hf.co/") || 
		   strings.HasPrefix(modelName, "huggingface.co/")
}

// GetHuggingFaceModelName extracts the repository name from a Hugging Face model reference
func GetHuggingFaceModelName(modelRef string) string {
	// Strip any prefix
	modelName := modelRef
	if strings.HasPrefix(modelRef, "hf.co/") {
		modelName = strings.TrimPrefix(modelRef, "hf.co/")
	} else if strings.HasPrefix(modelRef, "huggingface.co/") {
		modelName = strings.TrimPrefix(modelRef, "huggingface.co/")
	}
	
	// Strip any quantization suffix
	if colonIdx := strings.Index(modelName, ":"); colonIdx != -1 {
		modelName = modelName[:colonIdx]
	}
	
	return modelName
}

// GetQuantizationFromModelRef extracts the quantization suffix from a model reference
func GetQuantizationFromModelRef(modelRef string) string {
	if colonIdx := strings.Index(modelRef, ":"); colonIdx != -1 {
		return modelRef[colonIdx+1:]
	}
	return ""
} 