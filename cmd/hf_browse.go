package cmd

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var (
	browseQuantization string
	browseLimit        int
	browseOpen         bool
)

var hfBrowseCmd = &cobra.Command{
	Use:   "hf-browse [search-term]",
	Short: "Browse GGUF models on Hugging Face",
	Long: `Search and browse GGUF models available on Hugging Face.
Optionally open the browser to view the model details.

Examples:
  rlama hf-browse llama3       # Search for Llama 3 models
  rlama hf-browse mistral --open # Search and open browser
  rlama hf-browse "open instruct" --limit 5  # Limit results`,
	RunE: func(cmd *cobra.Command, args []string) error {
		searchTerm := ""
		if len(args) > 0 {
			searchTerm = strings.Join(args, " ")
		}

		// Build the search URL
		searchURL := "https://huggingface.co/models?search=gguf"
		if searchTerm != "" {
			searchURL += "+" + strings.ReplaceAll(searchTerm, " ", "+")
		}

		if browseOpen {
			var err error
			switch runtime.GOOS {
			case "linux":
				err = exec.Command("xdg-open", searchURL).Start()
			case "windows":
				err = exec.Command("rundll32", "url.dll,FileProtocolHandler", searchURL).Start()
			case "darwin":
				err = exec.Command("open", searchURL).Start()
			default:
				err = fmt.Errorf("unsupported platform")
			}
			if err != nil {
				return fmt.Errorf("error opening browser: %w", err)
			}
			fmt.Printf("Opened browser to search for GGUF models: %s\n", searchURL)
		}

		fmt.Println("To use a Hugging Face model with RLAMA, use the format:")
		fmt.Println("  rlama rag hf.co/username/repository my-rag ./docs")
		
		if browseQuantization != "" {
			fmt.Println("\nSpecify quantization:")
			fmt.Println("  rlama rag hf.co/username/repository:" + browseQuantization + " my-rag ./docs")
		}
		
		fmt.Println("\nOr try out a model directly with:")
		fmt.Println("  rlama run-hf username/repository")
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(hfBrowseCmd)
	
	hfBrowseCmd.Flags().StringVar(&browseQuantization, "quant", "", "Specify quantization type (e.g., Q4_K_M, Q5_K_M)")
	hfBrowseCmd.Flags().IntVar(&browseLimit, "limit", 10, "Limit number of results")
	hfBrowseCmd.Flags().BoolVar(&browseOpen, "open", false, "Open browser with search results")
} 