package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	
	"github.com/dontizi/rlama/internal/client"
)

const (
	Version = "0.1.21" // Current version of RLAMA
)

var rootCmd = &cobra.Command{
	Use:   "rlama",
	Short: "RLAMA is a CLI tool for creating and using RAG systems with Ollama",
	Long: `RLAMA (Retrieval-Augmented Language Model Adapter) is a command-line tool 
that simplifies the creation and use of RAG (Retrieval-Augmented Generation) systems 
based on Ollama models.

Main commands:
  rag [model] [rag-name] [folder-path]    Create a new RAG system
  run [rag-name]                          Run an existing RAG system
  list                                    List all available RAG systems
  delete [rag-name]                       Delete a RAG system
  update                                  Check and install RLAMA updates`,
}

// Variables to store command flags
var (
	versionFlag bool
	ollamaHost  string
	ollamaPort  string
)

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}

// GetOllamaClient returns an Ollama client configured with host and port from command flags
func GetOllamaClient() *client.OllamaClient {
	return client.NewOllamaClient(ollamaHost, ollamaPort)
}

func init() {
	// Add --version flag
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Display RLAMA version")
	
	// Add Ollama configuration flags
	rootCmd.PersistentFlags().StringVar(&ollamaHost, "host", "", "Ollama host (default: localhost)")
	rootCmd.PersistentFlags().StringVar(&ollamaPort, "port", "", "Ollama port (default: 11434)")
	
	// Override the Run function to handle the --version flag
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if versionFlag {
			fmt.Printf("RLAMA version %s\n", Version)
			return
		}
		
		// If no arguments are provided and --version is not used, display help
		if len(args) == 0 {
			cmd.Help()
		}
	}
} 