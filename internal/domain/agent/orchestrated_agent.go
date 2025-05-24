package agent

import (
	"context"
	"fmt"
	"strings"
)

// OrchestratedAgent represents an agent that uses task orchestration for complex queries
type OrchestratedAgent struct {
	*baseAgent
	llmClient    LLMClient
	orchestrator *Orchestrator
}

// NewOrchestratedAgent creates a new orchestrated agent
func NewOrchestratedAgent(memory Memory, llmClient LLMClient) *OrchestratedAgent {
	agent := &OrchestratedAgent{
		baseAgent: NewBaseAgent(AutonomousMode, memory),
		llmClient: llmClient,
	}

	// Create orchestrator with self-reference
	agent.orchestrator = NewOrchestrator(llmClient, agent)

	return agent
}

// Run implements the orchestrated agent's main execution logic
func (a *OrchestratedAgent) Run(ctx context.Context, input string) (string, error) {
	debugPrint("OrchestratedAgent: Starting run with input: %s", input)

	// Add user input to conversation history
	if err := a.memory.AddToHistory("User: " + input); err != nil {
		return "", fmt.Errorf("failed to add user input to history: %w", err)
	}

	// Determine if this is a complex query that needs orchestration
	isComplex, err := a.isComplexQuery(ctx, input)
	if err != nil {
		debugPrint("Error determining query complexity: %v", err)
		// Fallback to simple mode if we can't determine complexity
		isComplex = false
	}

	debugPrint("OrchestratedAgent: Query complexity assessment: %t", isComplex)

	if !isComplex {
		// Handle simple queries with direct tool execution
		return a.handleSimpleQuery(ctx, input)
	}

	// Handle complex queries with orchestration
	return a.handleComplexQuery(ctx, input)
}

// isComplexQuery determines if a query requires orchestration
func (a *OrchestratedAgent) isComplexQuery(ctx context.Context, query string) (bool, error) {
	prompt := fmt.Sprintf(`Analyze the following user query and determine if it requires multiple steps or tools to complete.

Query: %s

A query is considered COMPLEX if it:
1. Requires multiple pieces of information from different sources
2. Has dependencies between tasks (e.g., find event location, then find flights to that location)
3. Involves multiple steps that must be done in sequence
4. Combines real-time information with analysis or comparison

A query is considered SIMPLE if it:
1. Can be answered with a single tool or search
2. Requires only one piece of information
3. Is a straightforward question

Examples of COMPLEX queries:
- "When is the next Snowflake Summit and how much would it cost to attend from Montreal?"
- "Find the weather in Paris and suggest what to pack for a trip there"
- "Search for Python security vulnerabilities in my code and suggest fixes"

Examples of SIMPLE queries:
- "What's the weather today?"
- "List files in the current directory"
- "Search for information about machine learning"

Respond with only: COMPLEX or SIMPLE`, query)

	response, err := a.llmClient.GenerateCompletion(ctx, prompt)
	if err != nil {
		return false, err
	}

	response = strings.TrimSpace(strings.ToUpper(response))
	return strings.Contains(response, "COMPLEX"), nil
}

// handleSimpleQuery handles queries that don't need orchestration
func (a *OrchestratedAgent) handleSimpleQuery(ctx context.Context, input string) (string, error) {
	debugPrint("OrchestratedAgent: Handling simple query")

	// Build the prompt for the LLM
	prompt := a.buildSimplePrompt(input)

	// Get LLM's response
	response, err := a.llmClient.GenerateCompletion(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate LLM completion: %w", err)
	}

	// Parse LLM's response to get action and input
	action, toolInput, err := parseResponse(response)
	if err != nil {
		return "", fmt.Errorf("failed to parse LLM response: %w", err)
	}

	// If no action needed, return the response directly
	if action == "" {
		cleanResponse := cleanThinkingProcess(response)
		if err := a.memory.AddToHistory("Agent: " + cleanResponse); err != nil {
			return "", fmt.Errorf("failed to add agent response to history: %w", err)
		}
		return cleanResponse, nil
	}

	// Execute the tool
	tool := a.findTool(action)
	if tool == nil {
		return "", fmt.Errorf("tool %s not found", action)
	}

	result, err := tool.Execute(ctx, toolInput)
	if err != nil {
		return "", fmt.Errorf("failed to execute tool %s: %w", action, err)
	}

	// Generate final response with tool result
	finalPrompt := fmt.Sprintf(`User query: %s
Tool used: %s
Tool result: %s

Based on the tool result, provide a clear and helpful response to the user's query.`, input, action, result)

	finalResponse, err := a.llmClient.GenerateCompletion(ctx, finalPrompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate final response: %w", err)
	}

	cleanResponse := cleanThinkingProcess(finalResponse)
	if err := a.memory.AddToHistory("Agent: " + cleanResponse); err != nil {
		return "", fmt.Errorf("failed to add agent response to history: %w", err)
	}

	return cleanResponse, nil
}

// handleComplexQuery handles queries that need orchestration
func (a *OrchestratedAgent) handleComplexQuery(ctx context.Context, input string) (string, error) {
	if !Debug {
		fmt.Printf("\n🧠 Analyzing complex query and decomposing into tasks...\n")
	}
	debugPrint("OrchestratedAgent: Handling complex query with orchestration")

	// Step 1: Decompose the query into tasks
	tasks, err := a.orchestrator.DecomposeQuery(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to decompose query: %w", err)
	}

	if !Debug {
		fmt.Printf("📋 Decomposed into %d tasks\n", len(tasks))
	}
	debugPrint("OrchestratedAgent: Decomposed into %d tasks", len(tasks))

	// Step 2: Execute all tasks
	if err := a.orchestrator.ExecuteTasks(ctx, tasks); err != nil {
		return "", fmt.Errorf("failed to execute tasks: %w", err)
	}

	if !Debug {
		fmt.Printf("\n🔄 Synthesizing final response...\n")
	}

	// Step 3: Generate final response
	finalResponse, err := a.orchestrator.GenerateFinalResponse(ctx, input, tasks)
	if err != nil {
		return "", fmt.Errorf("failed to generate final response: %w", err)
	}

	// Add to conversation history
	if err := a.memory.AddToHistory("Agent: " + finalResponse); err != nil {
		return "", fmt.Errorf("failed to add agent response to history: %w", err)
	}

	return finalResponse, nil
}

// buildSimplePrompt constructs the prompt for simple queries
func (a *OrchestratedAgent) buildSimplePrompt(input string) string {
	history := a.memory.GetHistory()
	tools := a.getToolDescriptions()

	return fmt.Sprintf(`You are a helpful AI assistant with access to the following tools:

%s

Conversation history:
%s

Current user input: %s

IMPORTANT INSTRUCTIONS:
1. For any real-time information (weather, news, events, prices), ALWAYS use the web_search tool
2. For local document queries, use the rag_search tool
3. For file operations: use list_dir, read_file, or file_write as appropriate
4. NEVER say you can't do something without trying available tools first

To use a tool, respond in the format:
ACTION: <tool_name>
INPUT: <tool_input>

To respond directly to the user, write your response.

You can use <think>...</think> tags to show your reasoning process.`, tools, formatHistory(history), input)
}

// findTool finds a tool by name
func (a *OrchestratedAgent) findTool(name string) Tool {
	for _, tool := range a.GetTools() {
		if tool.Name() == name {
			return tool
		}
	}
	return nil
}

// getToolDescriptions returns formatted tool descriptions
func (a *OrchestratedAgent) getToolDescriptions() string {
	tools := a.GetTools()
	descriptions := make([]string, len(tools))
	for i, tool := range tools {
		descriptions[i] = fmt.Sprintf("- %s: %s", tool.Name(), tool.Description())
	}
	return strings.Join(descriptions, "\n")
}
