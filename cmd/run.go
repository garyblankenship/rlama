package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/service"
)

var runCmd = &cobra.Command{
	Use:   "run [rag-name]",
	Short: "Run a RAG system",
	Long: `Run a previously created RAG system. 
Starts an interactive session to interact with the RAG system.
Example: rlama run rag1`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]

		// Get Ollama client with configured host and port
		ollamaClient := GetOllamaClient()
		if err := ollamaClient.CheckOllamaAndModel(""); err != nil {
			return err
		}

		ragService := service.NewRagService(ollamaClient)
		rag, err := ragService.LoadRag(ragName)
		if err != nil {
			return err
		}

		fmt.Printf("RAG '%s' loaded. Model: %s\n", rag.Name, rag.ModelName)
		fmt.Println("Type your question (or 'exit' to quit):")

		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("> ")
			if !scanner.Scan() {
				break
			}

			question := scanner.Text()
			if question == "exit" {
				break
			}

			if strings.TrimSpace(question) == "" {
				continue
			}

			fmt.Println("\nSearching documents for relevant information...")

			answer, err := ragService.Query(rag, question)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				continue
			}

			fmt.Println("\n--- Answer ---")
			fmt.Println(answer)
			fmt.Println()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}