package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/domain/agent"
	"github.com/dontizi/rlama/internal/service"
	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Run an AI agent",
	Long: `Run an AI agent that can perform various tasks using tools.
The agent can operate in conversational or autonomous mode.`,
}

var agentRunCmd = &cobra.Command{
	Use:   "run [rag-name]",
	Short: "Run an agent with optional RAG system",
	Long: `Run an agent that can use various tools including RAG search.
If a RAG name is provided, the agent will have access to that knowledge base.

Agent Modes:
- conversation: Simple conversational agent for basic queries
- autonomous: Autonomous agent (not yet implemented)
- orchestrated: Advanced agent that decomposes complex queries into tasks (DEFAULT)

The orchestrated mode is perfect for complex queries like:
"When is the next Snowflake Summit and how much would it cost to attend from Montreal?"

For web search functionality:
1. Enable with -w or --web flag
2. Set GOOGLE_SEARCH_API_KEY environment variable or use --search-api-key flag
3. Set GOOGLE_SEARCH_ENGINE_ID environment variable or use --search-engine-id flag

Examples:
  # Simple query with orchestrated mode (default)
  rlama agent run -w -q "When is the next Snowflake Summit and how much to attend from Montreal?"

  # Complex query with explicit orchestrated mode
  rlama agent run -w -m orchestrated -q "Find Python security issues in my code and suggest fixes"

  # Simple conversational mode
  rlama agent run -w -m conversation -q "What's the weather today?"`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get flags
		mode, _ := cmd.Flags().GetString("mode")
		webEnabled, _ := cmd.Flags().GetBool("web")
		query, _ := cmd.Flags().GetString("query")
		searchApiKey, _ := cmd.Flags().GetString("search-api-key")
		searchEngineID, _ := cmd.Flags().GetString("search-engine-id")

		// Use the model specified for the agent command, fallback to global model if not specified
		model, _ := cmd.Flags().GetString("model")
		if model == "" {
			model = modelName // Use the global model name
		}

		// Convert mode string to AgentMode
		var agentMode agent.AgentMode
		switch mode {
		case "conversation":
			agentMode = agent.ConversationalMode
		case "autonomous":
			agentMode = agent.AutonomousMode
		case "orchestrated":
			agentMode = agent.OrchestratedMode
		default:
			return fmt.Errorf("invalid mode: %s", mode)
		}

		// Get working directory for file operations
		workDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		// If web search is enabled, ensure we have the necessary credentials
		if webEnabled {
			if searchApiKey == "" {
				searchApiKey = os.Getenv("GOOGLE_SEARCH_API_KEY")
				if searchApiKey == "" {
					return fmt.Errorf("web search enabled but no API key provided (use --search-api-key flag or set GOOGLE_SEARCH_API_KEY environment variable)")
				}
			}

			if searchEngineID == "" {
				searchEngineID = os.Getenv("GOOGLE_SEARCH_ENGINE_ID")
				if searchEngineID == "" {
					return fmt.Errorf("web search enabled but no Search Engine ID provided (use --search-engine-id flag or set GOOGLE_SEARCH_ENGINE_ID environment variable)")
				}
			}

			// Set the Search Engine ID in environment for the tool to use
			os.Setenv("GOOGLE_SEARCH_ENGINE_ID", searchEngineID)
		}

		// Create agent service config
		config := service.AgentServiceConfig{
			BaseDir:      workDir,
			WebEnabled:   webEnabled,
			SearchApiKey: searchApiKey,
		}

		// Create agent service
		var ragSystem *domain.RagSystem
		if len(args) > 0 {
			ragName := args[0]
			ragSystem, err = Services.RagService.LoadRag(ragName)
			if err != nil {
				return fmt.Errorf("failed to load RAG system '%s': %w", ragName, err)
			}
			fmt.Printf("📚 Loaded RAG system: %s\n", ragName)
		} else {
			// No RAG specified - check if we can auto-detect one
			availableRags, err := Services.RagService.ListAllRags()
			if err == nil && len(availableRags) == 1 {
				// Only one RAG available - use it automatically
				ragName := availableRags[0]
				ragSystem, err = Services.RagService.LoadRag(ragName)
				if err != nil {
					return fmt.Errorf("failed to auto-load RAG system '%s': %w", ragName, err)
				}
				fmt.Printf("📚 Auto-detected and loaded RAG system: %s\n", ragName)
			}
			// If no RAGs or multiple RAGs, continue without one (will use auto-detection tool)
		}

		// Enable debug mode if verbose is set
		if verbose {
			agent.Debug = true
		}

		// Create Ollama adapter
		llmClient := client.NewOllamaAdapter(Services.OllamaClient, model)

		agentService := service.NewAgentService(
			Services.RagService,
			ragSystem,
			llmClient,
			config,
		)

		// Run the agent
		fmt.Printf("🤖 Starting agent in %s mode with model %s...\n", mode, model)
		result, err := agentService.RunAgent(context.Background(), agentMode, query)
		if err != nil {
			return fmt.Errorf("agent execution failed: %w", err)
		}

		fmt.Printf("\n🎯 Agent response:\n%s\n", result)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(agentCmd)
	agentCmd.AddCommand(agentRunCmd)

	// Add flags
	agentRunCmd.Flags().StringP("mode", "m", "orchestrated", "Agent mode (conversation, autonomous, or orchestrated)")
	agentRunCmd.Flags().BoolP("web", "w", false, "Enable web search capability")
	agentRunCmd.Flags().StringP("query", "q", "", "Query or goal for the agent")
	agentRunCmd.Flags().StringP("model", "l", "", "Model to use for the agent (overrides global model)")
	agentRunCmd.Flags().String("search-api-key", "", "Google Custom Search API key (overrides GOOGLE_SEARCH_API_KEY env var)")
	agentRunCmd.Flags().String("search-engine-id", "", "Google Custom Search Engine ID (overrides GOOGLE_SEARCH_ENGINE_ID env var)")
	agentRunCmd.MarkFlagRequired("query")
}
