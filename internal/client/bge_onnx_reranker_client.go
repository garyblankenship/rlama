package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// BGEONNXRerankerClient handles BGE reranking using ONNX runtime via HTTP microservice
type BGEONNXRerankerClient struct {
	serverURL   string
	httpClient  *http.Client
	modelDir    string
	serverProc  *exec.Cmd
}

// NewBGEONNXRerankerClient creates a new ONNX-based BGE reranker client
func NewBGEONNXRerankerClient(modelDir string) (*BGEONNXRerankerClient, error) {
	// Find an available port
	port := findAvailablePort()
	
	client := &BGEONNXRerankerClient{
		serverURL:  fmt.Sprintf("http://localhost:%d", port),
		httpClient: &http.Client{Timeout: 30 * time.Second},
		modelDir:   modelDir,
	}
	
	// Start the Python ONNX server
	if err := client.startONNXServer(port); err != nil {
		return nil, fmt.Errorf("failed to start ONNX server: %w", err)
	}
	
	return client, nil
}

// findAvailablePort finds an available port for the ONNX server
func findAvailablePort() int {
	// Start from a base port and add a random offset to avoid conflicts
	basePort := 8765
	return basePort + rand.Intn(1000)
}

// startONNXServer starts a Python HTTP server that runs ONNX inference
func (c *BGEONNXRerankerClient) startONNXServer(port int) error {
	// Create the Python server script
	serverScript := fmt.Sprintf(`
import sys
import json
import warnings
from http.server import HTTPServer, BaseHTTPRequestHandler
from urllib.parse import urlparse
import onnxruntime as ort
from transformers import AutoTokenizer
import numpy as np

# Suppress warnings
warnings.filterwarnings("ignore")

class ONNXRerankerHandler(BaseHTTPRequestHandler):
    def __init__(self, *args, **kwargs):
        # Initialize model and tokenizer
        model_dir = "%s"
        self.session = ort.InferenceSession(f"{model_dir}/model.onnx")
        self.tokenizer = AutoTokenizer.from_pretrained("BAAI/bge-reranker-large")
        self.max_length = 512
        super().__init__(*args, **kwargs)
    
    def log_message(self, format, *args):
        # Suppress request logging
        pass
    
    def do_POST(self):
        if self.path == '/rerank':
            try:
                content_length = int(self.headers.get('Content-Length', 0))
                post_data = self.rfile.read(content_length)
                request_data = json.loads(post_data)
                
                pairs = request_data.get('pairs', [])
                normalize = request_data.get('normalize', True)
                
                scores = []
                for pair in pairs:
                    if len(pair) != 2:
                        raise ValueError(f"each pair must contain exactly 2 elements (query, passage), got {len(pair)}")
                    query, passage = pair[0], pair[1]
                    text = f"{query} </s> {passage}"
                    
                    # Tokenize
                    encoding = self.tokenizer(
                        text,
                        max_length=self.max_length,
                        padding='max_length',
                        truncation=True,
                        return_tensors='np'
                    )
                    
                    # Run inference
                    ort_inputs = {
                        'input_ids': encoding['input_ids'].astype(np.int64),
                        'attention_mask': encoding['attention_mask'].astype(np.int64)
                    }
                    
                    outputs = self.session.run(None, ort_inputs)
                    logits = outputs[0][0][0]  # Extract scalar logit
                    
                    # Apply normalization if requested
                    if normalize:
                        score = 1.0 / (1.0 + np.exp(-logits))
                    else:
                        score = float(logits)
                    
                    scores.append(score)
                
                response = {'scores': scores}
                self.send_response(200)
                self.send_header('Content-Type', 'application/json')
                self.end_headers()
                self.wfile.write(json.dumps(response).encode())
                
            except Exception as e:
                error_response = {'error': str(e)}
                self.send_response(500)
                self.send_header('Content-Type', 'application/json')
                self.end_headers()
                self.wfile.write(json.dumps(error_response).encode())
        else:
            self.send_response(404)
            self.end_headers()

if __name__ == '__main__':
    server = HTTPServer(('localhost', %d), ONNXRerankerHandler)
    print("ONNX Reranker server started on http://localhost:%d", file=sys.stderr)
    server.serve_forever()
`, c.modelDir, port, port)

	// Write server script to temporary file
	scriptPath := filepath.Join(c.modelDir, "onnx_server.py")
	if err := os.WriteFile(scriptPath, []byte(serverScript), 0644); err != nil {
		return fmt.Errorf("failed to write server script: %w", err)
	}

	// Start the server process
	c.serverProc = exec.Command("python", scriptPath)
	c.serverProc.Stdout = os.Stdout
	c.serverProc.Stderr = os.Stderr
	
	if err := c.serverProc.Start(); err != nil {
		return fmt.Errorf("failed to start Python server: %w", err)
	}

	// Wait for server to be ready
	for i := 0; i < 30; i++ {
		if c.isServerReady() {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("server failed to start within timeout")
}

// isServerReady checks if the server is responding
func (c *BGEONNXRerankerClient) isServerReady() bool {
	resp, err := c.httpClient.Get(c.serverURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return true
}

// ComputeScores calculates relevance scores between queries and passages using ONNX
func (c *BGEONNXRerankerClient) ComputeScores(pairs [][]string, normalize bool) ([]float64, error) {
	requestData := map[string]interface{}{
		"pairs":     pairs,
		"normalize": normalize,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.httpClient.Post(
		c.serverURL+"/rerank",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if errorMsg, ok := response["error"]; ok {
		return nil, fmt.Errorf("server error: %v", errorMsg)
	}

	scoresInterface, ok := response["scores"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid scores format in response")
	}

	scores := make([]float64, len(scoresInterface))
	for i, s := range scoresInterface {
		scores[i], ok = s.(float64)
		if !ok {
			return nil, fmt.Errorf("invalid score type at index %d", i)
		}
	}

	return scores, nil
}

// GetModelName returns the model identifier
func (c *BGEONNXRerankerClient) GetModelName() string {
	return "bge-reranker-large-onnx"
}

// Cleanup properly stops the server and frees resources
func (c *BGEONNXRerankerClient) Cleanup() error {
	if c.serverProc != nil {
		if err := c.serverProc.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill server process: %w", err)
		}
		c.serverProc.Wait()
	}
	return nil
}