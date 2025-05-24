package service

import (
	"context"
	"fmt"
	"os"

	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/domain/agent"
)

// AgentService orchestrates agent operations
type AgentService struct {
	ragService   RagService
	ragSystem    *domain.RagSystem
	llmClient    agent.LLMClient
	baseDir      string
	webEnabled   bool
	searchApiKey string
}

// AgentServiceConfig holds configuration for AgentService
type AgentServiceConfig struct {
	BaseDir      string
	WebEnabled   bool
	SearchApiKey string
}

// NewAgentService creates a new AgentService
func NewAgentService(
	ragService RagService,
	ragSystem *domain.RagSystem,
	llmClient agent.LLMClient,
	config AgentServiceConfig,
) *AgentService {
	return &AgentService{
		ragService:   ragService,
		ragSystem:    ragSystem,
		llmClient:    llmClient,
		baseDir:      config.BaseDir,
		webEnabled:   config.WebEnabled,
		searchApiKey: config.SearchApiKey,
	}
}

// RunAgent runs an agent with the given parameters
func (s *AgentService) RunAgent(ctx context.Context, mode agent.AgentMode, input string) (string, error) {
	// Create agent memory
	memory := agent.NewSimpleMemory()

	// Create appropriate agent type based on mode
	var agentInstance agent.Agent
	switch mode {
	case agent.ConversationalMode:
		agentInstance = agent.NewConversationalAgent(memory, s.llmClient)
	case agent.AutonomousMode:
		return "", fmt.Errorf("autonomous mode not yet implemented")
	case agent.OrchestratedMode:
		agentInstance = agent.NewOrchestratedAgent(memory, s.llmClient)
	default:
		return "", fmt.Errorf("unsupported agent mode: %s", mode)
	}

	// Add tools to the agent
	if err := s.setupTools(agentInstance); err != nil {
		return "", fmt.Errorf("failed to setup tools: %w", err)
	}

	// Run the agent
	result, err := agentInstance.Run(ctx, input)
	if err != nil {
		return "", fmt.Errorf("agent execution failed: %w", err)
	}

	return result, nil
}

// setupTools configures the tools available to the agent
func (s *AgentService) setupTools(agentInstance agent.Agent) error {
	// Add RAG search tool only if a RAG system is available
	if s.ragSystem != nil {
		// Create RAG service adapter
		ragAdapter := NewRagServiceAdapter(s.ragService, s.ragSystem)

		// Add RAG search tool with adapter
		ragTool := agent.NewRagSearchTool(ragAdapter)
		if err := agentInstance.AddTool(ragTool); err != nil {
			return fmt.Errorf("failed to add RAG search tool: %w", err)
		}
	} else {
		// If no RAG system is specified, add a RAG auto-detection tool
		ragAutoTool := agent.NewRagAutoDetectionTool(s.ragService)
		if err := agentInstance.AddTool(ragAutoTool); err != nil {
			return fmt.Errorf("failed to add RAG auto-detection tool: %w", err)
		}
	}

	// Add file tools
	readTool := agent.NewFileTool(s.baseDir, "read")
	if err := agentInstance.AddTool(readTool); err != nil {
		return fmt.Errorf("failed to add file read tool: %w", err)
	}

	writeTool := agent.NewFileTool(s.baseDir, "write")
	if err := agentInstance.AddTool(writeTool); err != nil {
		return fmt.Errorf("failed to add file write tool: %w", err)
	}

	// Add directory listing tool
	listTool := agent.NewDirectoryListTool(s.baseDir)
	if err := agentInstance.AddTool(listTool); err != nil {
		return fmt.Errorf("failed to add directory list tool: %w", err)
	}

	// Add grep search tool for exact/regex search
	grepTool := agent.NewGrepSearchTool(s.baseDir)
	if err := agentInstance.AddTool(grepTool); err != nil {
		return fmt.Errorf("failed to add grep search tool: %w", err)
	}

	// Add file search tool for fuzzy file name search
	fileSearchTool := agent.NewFileSearchTool(s.baseDir)
	if err := agentInstance.AddTool(fileSearchTool); err != nil {
		return fmt.Errorf("failed to add file search tool: %w", err)
	}

	// Add codebase search tool for semantic code search
	codebaseSearchTool := agent.NewCodebaseSearchTool(s.baseDir)
	if err := agentInstance.AddTool(codebaseSearchTool); err != nil {
		return fmt.Errorf("failed to add codebase search tool: %w", err)
	}

	// Add web search tool if enabled
	if s.webEnabled {
		if s.searchApiKey == "" {
			// Try to get API key from environment
			s.searchApiKey = os.Getenv("GOOGLE_SEARCH_API_KEY")
			if s.searchApiKey == "" {
				return fmt.Errorf("web search enabled but no API key provided (set GOOGLE_SEARCH_API_KEY environment variable)")
			}
		}
		webTool := agent.NewWebSearchTool(s.searchApiKey)
		if err := agentInstance.AddTool(webTool); err != nil {
			return fmt.Errorf("failed to add web search tool: %w", err)
		}
	}

	return nil
}
