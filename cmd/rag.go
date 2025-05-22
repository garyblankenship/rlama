package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/service"
	"github.com/spf13/cobra"
)

var (
	excludeDirs          []string
	excludeExts          []string
	processExts          []string
	chunkSize            int
	chunkOverlap         int
	chunkingStrategy     string
	profileName          string
	embeddingModel       string
	ragDisableReranker   bool
	ragRerankerModel     string
	ragRerankerWeight    float64
	ragRerankerThreshold float64
	ragUseONNXReranker   bool
	ragONNXModelDir      string
	// Qdrant configuration flags
	vectorStore          string
	qdrantHost           string
	qdrantPort           int
	qdrantAPIKey         string
	qdrantCollection     string
	qdrantUseGRPC        bool
	// Enhanced document loader flags
	useAdvancedLoader    bool
	loaderStrategy       string
	testService          interface{} // Pour les tests
)

var ragCmd = &cobra.Command{
	Use:   "rag [model] [rag-name] [folder-path]",
	Short: "Create a new RAG system",
	Long: `Create a new RAG system by indexing all documents in the specified folder.
Example: rlama rag llama3.2 rag1 ./documents

The folder will be created if it doesn't exist yet.
Supported formats include: .txt, .md, .html, .json, .csv, and various source code files.

You can exclude directories or file types:
  rlama rag llama3 myproject ./code --excludedir=node_modules,dist,.git
  rlama rag llama3 mydocs ./docs --excludeext=.log,.tmp
  rlama rag llama3 specific ./mixed --processext=.md,.py,.js

Hugging Face Models:
  You can use Hugging Face GGUF models with the format:
  rlama rag hf.co/username/repository my-rag ./docs
  
  Specify quantization with:
  rlama rag hf.co/username/repository:Q4_K_M my-rag ./docs
  
OpenAI Models:
  You can use OpenAI models by setting the OPENAI_API_KEY environment variable:
  export OPENAI_API_KEY="your-api-key"
  
  Then use any OpenAI model:
  rlama rag gpt-4-turbo my-openai-rag ./docs
  rlama rag gpt-3.5-turbo my-openai-rag ./docs`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		modelName := args[0]
		ragName := args[1]
		folderPath := args[2]

		// Setup client
		llmClient, ollamaClient, err := setupLLMClient(modelName)
		if err != nil {
			return err
		}

		// Setup configuration
		loaderOptions, err := setupLoaderOptions(ragName)
		if err != nil {
			return err
		}

		// Apply enhanced loader configuration
		if err := applyEnhancedLoaderConfig(); err != nil {
			return fmt.Errorf("failed to configure enhanced document loader: %w", err)
		}

		// Create RAG
		ragService, err := createRagSystem(modelName, ragName, folderPath, llmClient, ollamaClient, loaderOptions)
		if err != nil {
			return handleRagCreationError(err)
		}

		// Post-creation configuration
		if err := configurePostCreation(cmd, ragService, ragName); err != nil {
			return err
		}

		fmt.Printf("RAG '%s' created successfully.\n", ragName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(ragCmd)

	// Add exclusion and processing flags
	ragCmd.Flags().StringSliceVar(&excludeDirs, "exclude-dir", []string{}, "Directories to exclude (comma-separated)")
	ragCmd.Flags().StringSliceVar(&excludeExts, "exclude-ext", []string{}, "File extensions to exclude (comma-separated)")
	ragCmd.Flags().StringSliceVar(&processExts, "process-ext", []string{}, "File extensions to process (others will be ignored)")

	// Add flags for chunking options
	ragCmd.Flags().IntVar(&chunkSize, "chunk-size", 1000, "Character count per chunk")
	ragCmd.Flags().IntVar(&chunkOverlap, "chunk-overlap", 200, "Overlap between chunks in characters")
	ragCmd.Flags().StringVar(&chunkingStrategy, "chunking", "hybrid", "Chunking strategy (options: fixed, semantic, hybrid, hierarchical)")
	ragCmd.Flags().StringVar(&chunkingStrategy, "chunking-strategy", "hybrid", "Chunking strategy (options: fixed, semantic, hybrid, hierarchical)")

	// Add reranking options - now with a flag to disable it instead
	ragCmd.Flags().BoolVar(&ragDisableReranker, "disable-reranker", false, "Disable reranking (enabled by default)")
	ragCmd.Flags().StringVar(&ragRerankerModel, "reranker-model", "", "Model to use for reranking (defaults to main model)")
	ragCmd.Flags().Float64Var(&ragRerankerWeight, "reranker-weight", 0.7, "Weight for reranker scores vs vector scores (0-1)")
	ragCmd.Flags().Float64Var(&ragRerankerThreshold, "reranker-threshold", 0.0, "Minimum score threshold for reranked results")
	ragCmd.Flags().BoolVar(&ragUseONNXReranker, "use-onnx-reranker", false, "Use ONNX reranker for faster performance")
	ragCmd.Flags().StringVar(&ragONNXModelDir, "onnx-model-dir", "./models/bge-reranker-large-onnx", "Directory containing ONNX reranker model")

	// Add profile option
	ragCmd.Flags().StringVar(&profileName, "profile", "", "API profile name for OpenAI models")
	
	// Add embedding model option
	ragCmd.Flags().StringVar(&embeddingModel, "embedding-model", "", "Model to use for embeddings (defaults to snowflake-arctic-embed2, then falls back to main model)")

	// Add Qdrant configuration flags
	ragCmd.Flags().StringVar(&vectorStore, "vector-store", "internal", "Vector store type (internal, qdrant)")
	ragCmd.Flags().StringVar(&qdrantHost, "qdrant-host", "localhost", "Qdrant server host")
	ragCmd.Flags().IntVar(&qdrantPort, "qdrant-port", 6333, "Qdrant server port")
	ragCmd.Flags().StringVar(&qdrantAPIKey, "qdrant-apikey", "", "Qdrant API key for secured instances")
	ragCmd.Flags().StringVar(&qdrantCollection, "qdrant-collection", "", "Qdrant collection name (defaults to RAG name)")
	ragCmd.Flags().BoolVar(&qdrantUseGRPC, "qdrant-grpc", false, "Use gRPC for Qdrant communication")

	// Add enhanced document loader flags
	ragCmd.Flags().BoolVar(&useAdvancedLoader, "advanced-loader", true, "Use enhanced document processing with LangChain (default: true)")
	ragCmd.Flags().StringVar(&loaderStrategy, "loader-strategy", "", "Document loading strategy: langchain, legacy, or hybrid (default: hybrid)")

	// Add help text for the new flags
	ragCmd.Long += `

Enhanced Document Processing:
  --advanced-loader=false        Disable advanced document processing
  --loader-strategy=hybrid       Choose loading strategy (langchain, legacy, hybrid)
  
Environment Variables:
  RLAMA_LOADER_STRATEGY          Set default loading strategy
  RLAMA_USE_LANGCHAIN_LOADER     Enable/disable LangChain loader (true/false)
  RLAMA_DEBUG_LOADER             Enable debug output for document loading`

	// Add logic to use the test service if available
	if testService != nil {
		// Use the test service
	}
}

// NewRagCommand returns the rag command
func NewRagCommand() *cobra.Command {
	return ragCmd
}

// InjectTestService injects a test service
func InjectTestService(service interface{}) {
	testService = service
}

// setupLLMClient creates and configures the appropriate LLM client using ServiceProvider
func setupLLMClient(modelName string) (client.LLMClient, *client.OllamaClient, error) {
	provider := GetServiceProvider()
	
	// Update provider config with profile if specified
	if profileName != "" {
		config := provider.GetConfig().WithProfile(profileName)
		if err := provider.UpdateConfig(config); err != nil {
			return nil, nil, fmt.Errorf("error updating service provider config: %w", err)
		}
	}
	
	// Get clients from provider
	llmClient, err := provider.GetLLMClient(modelName)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting LLM client: %w", err)
	}
	
	ollamaClient := provider.GetOllamaClient()

	// Debug: Check what type of client we got
	if profileName != "" {
		if _, isOpenAI := llmClient.(*client.OpenAIClient); isOpenAI {
			fmt.Printf("✓ Successfully created OpenAI-compatible client for profile '%s'\n", profileName)
		} else {
			fmt.Printf("⚠️ Warning: Expected OpenAI client but got different type for profile '%s'\n", profileName)
		}
	}

	// Check the client and model
	if err := llmClient.CheckLLMAndModel(modelName); err != nil {
		return nil, nil, err
	}

	// Display which client/provider is being used
	if profileName != "" {
		fmt.Printf("Using model '%s' with profile '%s'.\n", modelName, profileName)
	} else if client.IsOpenAIModel(modelName) {
		fmt.Printf("Using OpenAI model '%s' (no profile specified, using environment variable).\n", modelName)
	} else {
		fmt.Printf("Using model '%s' with Ollama.\n", modelName)
	}

	// Handle Hugging Face models (Ollama-specific)
	if client.IsHuggingFaceModel(modelName) {
		// Extract quantization if specified
		hfModelName := client.GetHuggingFaceModelName(modelName)
		quantization := client.GetQuantizationFromModelRef(modelName)

		// Pull the model from Hugging Face
		fmt.Printf("Pulling Hugging Face model %s...\n", hfModelName)
		if err := ollamaClient.PullHuggingFaceModel(hfModelName, quantization); err != nil {
			return nil, nil, fmt.Errorf("error pulling Hugging Face model: %w", err)
		}
	}

	return llmClient, ollamaClient, nil
}

// setupLoaderOptions creates and validates document loader options using ServiceProvider
func setupLoaderOptions(ragName string) (service.DocumentLoaderOptions, error) {
	provider := GetServiceProvider()
	config := provider.GetConfig()
	
	// Override ONNX configuration if specified via flags
	if ragUseONNXReranker {
		config.UseONNXReranker = true
		config.ONNXModelDir = ragONNXModelDir
		
		// Update the provider configuration
		if err := provider.UpdateConfig(config); err != nil {
			return service.DocumentLoaderOptions{}, fmt.Errorf("error updating service provider config: %w", err)
		}
	}
	
	// Create base options from configuration
	loaderOptions := config.ToDocumentLoaderOptions()
	
	// Override with command-specific flags
	loaderOptions.ExcludeDirs = excludeDirs
	loaderOptions.ExcludeExts = excludeExts
	loaderOptions.ProcessExts = processExts
	
	// Override chunk settings if provided via flags
	if chunkSize != 1000 { // Default value check
		loaderOptions.ChunkSize = chunkSize
	}
	if chunkOverlap != 200 { // Default value check
		loaderOptions.ChunkOverlap = chunkOverlap
	}
	if chunkingStrategy != "hybrid" { // Default value check
		loaderOptions.ChunkingStrategy = chunkingStrategy
	}
	
	// Override profile and embedding settings if provided
	if profileName != "" {
		loaderOptions.APIProfileName = profileName
	}
	if embeddingModel != "" {
		loaderOptions.EmbeddingModel = embeddingModel
	}
	
	// Override reranker settings if provided
	loaderOptions.EnableReranker = !ragDisableReranker
	if ragRerankerModel != "" {
		loaderOptions.RerankerModel = ragRerankerModel
	}
	if ragRerankerWeight != 0.7 { // Default value check
		loaderOptions.RerankerWeight = ragRerankerWeight
	}
	loaderOptions.UseONNXReranker = ragUseONNXReranker
	if ragONNXModelDir != "./models/bge-reranker-large-onnx" { // Default value check
		loaderOptions.ONNXModelDir = ragONNXModelDir
	}
	
	// Override vector store settings if provided
	if vectorStore != "internal" { // Default value check
		loaderOptions.VectorStore = vectorStore
	}
	if qdrantHost != "localhost" { // Default value check
		loaderOptions.QdrantHost = qdrantHost
	}
	if qdrantPort != 6333 { // Default value check
		loaderOptions.QdrantPort = qdrantPort
	}
	if qdrantAPIKey != "" {
		loaderOptions.QdrantAPIKey = qdrantAPIKey
	}
	if qdrantCollection != "" {
		loaderOptions.QdrantCollectionName = qdrantCollection
	}
	loaderOptions.QdrantGRPC = qdrantUseGRPC

	// Set default collection name if not provided
	if loaderOptions.QdrantCollectionName == "" && loaderOptions.VectorStore == "qdrant" {
		loaderOptions.QdrantCollectionName = ragName
	}

	// Validate Qdrant configuration if using Qdrant
	if loaderOptions.VectorStore == "qdrant" {
		if loaderOptions.QdrantHost == "" {
			return loaderOptions, fmt.Errorf("Qdrant host cannot be empty when using --vector-store=qdrant")
		}
		if loaderOptions.QdrantPort <= 0 || loaderOptions.QdrantPort > 65535 {
			return loaderOptions, fmt.Errorf("Qdrant port must be between 1 and 65535")
		}
		fmt.Printf("Using Qdrant vector store at %s:%d with collection '%s'\n", 
			loaderOptions.QdrantHost, loaderOptions.QdrantPort, loaderOptions.QdrantCollectionName)
	}

	return loaderOptions, nil
}

// createRagSystem creates the RAG system with the specified configuration using ServiceProvider
func createRagSystem(modelName, ragName, folderPath string, llmClient client.LLMClient, ollamaClient *client.OllamaClient, loaderOptions service.DocumentLoaderOptions) (service.RagService, error) {
	// Display a message to indicate that the process has started
	fmt.Printf("Creating RAG '%s' with model '%s' from folder '%s'...\n",
		ragName, modelName, folderPath)

	// Use the service provider to create a RAG service for the specific model
	provider := GetServiceProvider()
	ragService, err := provider.CreateRagServiceForModel(modelName)
	if err != nil {
		return nil, fmt.Errorf("failed to create RAG service: %w", err)
	}
	
	err = ragService.CreateRagWithOptions(modelName, ragName, folderPath, loaderOptions)
	if err != nil {
		return nil, err
	}

	return ragService, nil
}

// handleRagCreationError provides improved error messages for RAG creation failures
func handleRagCreationError(err error) error {
	// Improve error messages related to Ollama
	if strings.Contains(err.Error(), "connection refused") {
		return fmt.Errorf("⚠️ Unable to connect to Ollama.\n" +
			"Make sure Ollama is installed and running.\n")
	}
	return err
}

// configurePostCreation handles post-creation configuration like reranker threshold
func configurePostCreation(cmd *cobra.Command, ragService service.RagService, ragName string) error {
	// Set reranker threshold if specified
	if cmd.Flags().Changed("reranker-threshold") {
		// Load the RAG that was just created
		rag, err := ragService.LoadRag(ragName)
		if err != nil {
			return fmt.Errorf("error setting reranker threshold: %w", err)
		}

		// Set the threshold
		rag.RerankerThreshold = ragRerankerThreshold

		// Save the updated RAG
		err = ragService.UpdateRag(rag)
		if err != nil {
			return fmt.Errorf("error updating reranker threshold: %w", err)
		}
	}
	return nil
}

// applyEnhancedLoaderConfig configures the enhanced document loader based on command flags
func applyEnhancedLoaderConfig() error {
	// Apply loader strategy from command line flag
	if loaderStrategy != "" {
		os.Setenv("RLAMA_LOADER_STRATEGY", loaderStrategy)
	}
	
	// Apply advanced loader flag
	if !useAdvancedLoader {
		os.Setenv("RLAMA_USE_LANGCHAIN_LOADER", "false")
	}
	
	// Debug output if requested
	if service.IsDebugMode() {
		fmt.Printf("📋 Enhanced Document Loader Configuration:\n")
		fmt.Printf("   Strategy: %s\n", service.GetLoaderStrategy())
		fmt.Printf("   LangChain Enabled: %t\n", service.UseLangChainLoader())
		fmt.Printf("   Debug Mode: %t\n", service.IsDebugMode())
	}
	
	return nil
}
