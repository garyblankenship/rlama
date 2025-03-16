package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	runHfQuantization string
)

var runHfCmd = &cobra.Command{
	Use:   "run-hf [huggingface-model]",
	Short: "Run a Hugging Face GGUF model with Ollama",
	Long: `Run a Hugging Face GGUF model directly using Ollama.
This is convenient for testing models before creating a RAG system with them.

Examples:
  rlama run-hf bartowski/Llama-3.2-1B-Instruct-GGUF
  rlama run-hf mlabonne/Meta-Llama-3.1-8B-Instruct-abliterated-GGUF --quant Q5_K_M`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		modelPath := args[0]
		
		// Prepare the model reference
		if !strings.Contains(modelPath, "/") {
			return fmt.Errorf("invalid model format. Use 'username/repository' format")
		}
		
		// Get the Ollama client
		ollamaClient := GetOllamaClient()
		
		fmt.Printf("Running Hugging Face model: %s\n", modelPath)
		if runHfQuantization != "" {
			fmt.Printf("Using quantization: %s\n", runHfQuantization)
		}
		
		return ollamaClient.RunHuggingFaceModel(modelPath, runHfQuantization)
	},
}

func init() {
	rootCmd.AddCommand(runHfCmd)
	
	runHfCmd.Flags().StringVar(&runHfQuantization, "quant", "", "Quantization to use (e.g., Q4_K_M, Q5_K_M)")
} 