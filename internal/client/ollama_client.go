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

// OllamaClient is a client for the Ollama API
type OllamaClient struct {
	BaseURL string
	Client  *http.Client
}

// EmbeddingRequest is the structure of the request for the /api/embeddings API
type EmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// EmbeddingResponse is the structure of the response for the /api/embeddings API
type EmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

// GenerationRequest is the structure of the request for the /api/generate API
type GenerationRequest struct {
	Model    string  `json:"model"`
	Prompt   string  `json:"prompt"`
	Context  []int   `json:"context,omitempty"`
	Options  Options `json:"options,omitempty"`
	Format   string  `json:"format,omitempty"`
	Template string  `json:"template,omitempty"`
	Stream   bool    `json:"stream"`
}

// Options for the /api/generate API
type Options struct {
	Temperature float64 `json:"temperature,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
	TopK        int     `json:"top_k,omitempty"`
	NumPredict  int     `json:"num_predict,omitempty"`
}

// GenerationResponse is the structure of the response for the /api/generate API
type GenerationResponse struct {
	Model     string `json:"model"`
	Response  string `json:"response"`
	Context   []int  `json:"context"`
	CreatedAt string `json:"created_at"`
	Done      bool   `json:"done"`
}

// NewOllamaClient creates a new Ollama client
// If host or port are empty, the default values are used
// If OLLAMA_HOST is defined, it is used as the default value
func NewOllamaClient(host, port string) *OllamaClient {
	// Check for OLLAMA_HOST environment variable
	ollamaHostEnv := os.Getenv("OLLAMA_HOST")

	// Default values and protocol
	defaultHost := DefaultOllamaHost
	defaultPort := DefaultOllamaPort
	protocol := "http://"

	// If OLLAMA_HOST is set, parse it
	if ollamaHostEnv != "" {
		// Handle if OLLAMA_HOST includes protocol
		if strings.HasPrefix(ollamaHostEnv, "http://") || strings.HasPrefix(ollamaHostEnv, "https://") {
			// Extract protocol and host
			parts := strings.SplitN(ollamaHostEnv, "://", 2)
			protocol = parts[0] + "://"
			hostParts := strings.Split(parts[1], ":")
			if len(hostParts) >= 1 {
				defaultHost = hostParts[0]
			}
			if len(hostParts) >= 2 {
				defaultPort = hostParts[1]
			}
		} else {
			// No protocol specified, use standard pattern
			parts := strings.Split(ollamaHostEnv, ":")
			if len(parts) >= 1 {
				defaultHost = parts[0]
			}
			if len(parts) >= 2 {
				defaultPort = parts[1]
			}
		}
	}

	// Command flags override environment variables
	if host != "" {
		// Check if host includes protocol
		if strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://") {
			parts := strings.SplitN(host, "://", 2)
			protocol = parts[0] + "://"
			host = strings.Split(parts[1], ":")[0]
		}
		defaultHost = host
	} else {
		host = defaultHost
	}

	if port != "" {
		defaultPort = port
	}

	baseURL := fmt.Sprintf("%s%s:%s", protocol, host, defaultPort)

	return &OllamaClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

// NewDefaultOllamaClient creates a new Ollama client with the default values
// Kept for compatibility with existing code
func NewDefaultOllamaClient() *OllamaClient {
	return NewOllamaClient(DefaultOllamaHost, DefaultOllamaPort)
}

// GenerateEmbedding generates an embedding for the given text
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

// GenerateCompletion generates a response for the given prompt
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
		return fmt.Errorf("⚠️ Ollama is not installed or not running.\n" +
			"RLAMA requires Ollama to function.\n" +
			"Please install Ollama with: curl -fsSL https://ollama.com/install.sh | sh\n" +
			"Then start it before using RLAMA.")
	}

	if !running {
		return fmt.Errorf("⚠️ Ollama is not running.\n" +
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
