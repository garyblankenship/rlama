package cmd

import (
	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/server"
)

var (
	apiPort string
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the RLAMA API server",
	Long: `Start an HTTP API server for RLAMA, allowing RAG operations via RESTful endpoints.
	
Example: rlama api --port 11249

Available endpoints:
- POST /rag: Query a RAG system
  Request body: { "rag_name": "my-docs", "prompt": "How many documents do you have?", "context_size": 20 }
  
- GET /health: Check server health`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get Ollama client with configured host and port
		ollamaClient := GetOllamaClient()
		
		// Create and start the server
		srv := server.NewServer(apiPort, ollamaClient)
		return srv.Start()
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
	
	// Add port flag
	apiCmd.Flags().StringVar(&apiPort, "port", "11249", "Port to run the API server on")
} 