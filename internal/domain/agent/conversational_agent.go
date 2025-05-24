package agent

import (
	"context"
	"fmt"
	"strings"
)

// Debug flag
var Debug = false

// ConversationalAgent represents an agent that operates in a conversational manner
type ConversationalAgent struct {
	*baseAgent
	llmClient LLMClient
}

// debugPrint prints debug information if Debug is true
func debugPrint(format string, args ...interface{}) {
	if Debug {
		fmt.Printf("\n🔍 DEBUG: "+format+"\n", args...)
	}
}

// NewConversationalAgent creates a new conversational agent
func NewConversationalAgent(memory Memory, llmClient LLMClient) *ConversationalAgent {
	return &ConversationalAgent{
		baseAgent: NewBaseAgent(ConversationalMode, memory),
		llmClient: llmClient,
	}
}

// Run implements the conversational agent's main loop
func (a *ConversationalAgent) Run(ctx context.Context, input string) (string, error) {
	debugPrint("Starting agent run with input: %s", input)

	// Add user input to conversation history
	if err := a.memory.AddToHistory("User: " + input); err != nil {
		return "", fmt.Errorf("failed to add user input to history: %w", err)
	}

	// Build the prompt for the LLM
	prompt := a.buildPrompt(input)
	debugPrint("Built prompt for LLM:\n%s", prompt)

	// Get LLM's response
	response, err := a.llmClient.GenerateCompletion(ctx, prompt)
	if err != nil {
		debugPrint("Error getting LLM completion: %v", err)
		return "", fmt.Errorf("failed to generate LLM completion: %w", err)
	}
	debugPrint("Raw LLM response:\n%s", response)

	// Parse LLM's response to get action and input
	action, toolInput, err := parseResponse(response)
	if err != nil {
		debugPrint("Error parsing LLM response: %v", err)
		return "", fmt.Errorf("failed to parse LLM response: %w", err)
	}
	debugPrint("Parsed response - Action: %s, Input: %s", action, toolInput)

	// If no action needed, clean and return the response directly
	if action == "" {
		debugPrint("No action needed, returning direct response")
		cleanResponse := cleanThinkingProcess(response)
		return cleanResponse, nil
	}

	// Execute the tool
	tool := a.findTool(action)
	if tool == nil {
		debugPrint("Tool not found: %s", action)
		debugPrint("Available tools: %v", a.getToolNames())
		return "", fmt.Errorf("tool %s not found", action)
	}

	debugPrint("Executing tool %s with input: %s", action, toolInput)
	result, err := tool.Execute(ctx, toolInput)
	if err != nil {
		debugPrint("Error executing tool: %v", err)
		return "", fmt.Errorf("failed to execute tool %s: %w", action, err)
	}
	debugPrint("Tool execution result:\n%s", result)

	// Add tool execution to history
	if err := a.memory.AddToHistory(fmt.Sprintf("Tool %s executed with result: %s", action, result)); err != nil {
		return "", fmt.Errorf("failed to add tool result to history: %w", err)
	}

	// Loop to handle potential chained actions
	maxActions := 5 // Prevent infinite loops
	actionCount := 1

	for actionCount <= maxActions {
		// Get response from LLM with tool result
		finalPrompt := a.buildFinalPrompt(input, action, result)
		debugPrint("Built final prompt:\n%s", finalPrompt)

		finalResponse, err := a.llmClient.GenerateCompletion(ctx, finalPrompt)
		if err != nil {
			debugPrint("Error getting final LLM completion: %v", err)
			return "", fmt.Errorf("failed to generate final response: %w", err)
		}
		debugPrint("Raw final LLM response:\n%s", finalResponse)

		// Check if the response contains another action
		nextAction, nextInput, err := parseResponse(finalResponse)
		if err != nil {
			debugPrint("Error parsing final response: %v", err)
			return "", fmt.Errorf("failed to parse final response: %w", err)
		}

		// If no further action, we're done
		if nextAction == "" {
			cleanResponse := cleanThinkingProcess(finalResponse)
			debugPrint("No more actions needed, final response:\n%s", cleanResponse)

			if err := a.memory.AddToHistory("Agent: " + cleanResponse); err != nil {
				return "", fmt.Errorf("failed to add agent response to history: %w", err)
			}
			return cleanResponse, nil
		}

		// Execute the next tool
		debugPrint("Executing chained action %d: %s with input: %s", actionCount+1, nextAction, nextInput)
		nextTool := a.findTool(nextAction)
		if nextTool == nil {
			debugPrint("Chained tool not found: %s", nextAction)
			return "", fmt.Errorf("chained tool %s not found", nextAction)
		}

		nextResult, err := nextTool.Execute(ctx, nextInput)
		if err != nil {
			debugPrint("Error executing chained tool: %v", err)
			return "", fmt.Errorf("failed to execute chained tool %s: %w", nextAction, err)
		}
		debugPrint("Chained tool execution result:\n%s", nextResult)

		// Add chained tool execution to history
		if err := a.memory.AddToHistory(fmt.Sprintf("Tool %s executed with result: %s", nextAction, nextResult)); err != nil {
			return "", fmt.Errorf("failed to add chained tool result to history: %w", err)
		}

		// Update for next iteration
		action = nextAction
		result = nextResult
		actionCount++
	}

	return "", fmt.Errorf("reached maximum number of chained actions (%d)", maxActions)
}

// getToolNames returns a list of available tool names
func (a *ConversationalAgent) getToolNames() []string {
	tools := a.GetTools()
	names := make([]string, len(tools))
	for i, tool := range tools {
		names[i] = tool.Name()
	}
	return names
}

// buildPrompt constructs the prompt for the LLM
func (a *ConversationalAgent) buildPrompt(input string) string {
	history := a.memory.GetHistory()
	tools := a.getToolDescriptions()

	return fmt.Sprintf(`You are a helpful AI assistant with access to the following tools:

%s

Conversation history:
%s

Current user input: %s

IMPORTANT INSTRUCTIONS:
1. For any real-time information (weather, news, etc.), ALWAYS use the web_search tool first
2. For local document queries, use the rag_search tool
3. For file operations: ALWAYS use list_dir FIRST to see what files exist, then use read_file or file_write
4. NEVER try to read a file without first listing the directory to get the actual file names
5. NEVER say you can't do something without trying available tools first
6. If a tool fails, try another approach or explain why it's not possible

SEARCH STRATEGY (use the right tool for the task):
- list_dir: To see what files/folders exist in a directory
- file_search: To find files by name (fuzzy matching) - e.g., "security", "config"
- grep_search: To find exact text patterns in files - e.g., "function.*eval|*.py" 
- codebase_search: To find code semantically - e.g., "error handling", "database connection"
- read_file: To read specific file content after you know it exists

To use a tool, respond in the format:
ACTION: <tool_name>
INPUT: <tool_input>

To respond directly to the user, write your response.

You can use <think>...</think> tags to show your reasoning process, but this will be removed from the final response.

Example for weather query:
<think>User asks about weather, I should use web_search to get real-time data</think>
ACTION: web_search
INPUT: current weather in user's location

Example for document query:
<think>User asks about a document, I should search the local knowledge base</think>
ACTION: rag_search
INPUT: user's document query

Example for file analysis:
<think>User wants to analyze files, I need to list the directory first to see what files exist</think>
ACTION: list_dir
INPUT: .

Example for finding security-related code:
<think>User wants to find security code, I should use codebase_search for semantic search</think>
ACTION: codebase_search
INPUT: security vulnerabilities authentication

Remember: You MUST use web_search for real-time information like weather. Always list directories before reading files!`, tools, formatHistory(history), input)
}

// buildFinalPrompt constructs the final prompt with tool results
func (a *ConversationalAgent) buildFinalPrompt(input, action, result string) string {
	return fmt.Sprintf(`Based on the user's input and the tool result, determine if you need to use another tool or provide a final response.

User input: %s
Tool used: %s
Tool result: %s

ANALYSIS:
- If the user asked to "list files then read X" and you just listed files, you MUST now read the file
- If the user asked to analyze/read a file and you just listed the directory, you MUST now read the specific file
- If you have all the information needed to answer the user's question, provide your final response

DECISION RULES:
1. If user wants to read/analyze a specific file and you only listed the directory: USE read_file
2. If user wants multiple actions and you only completed the first: CONTINUE with the next action
3. If you have completed all requested actions: PROVIDE final response

EXAMPLE: If user said "list files then read test.py" and you just listed files, respond:
ACTION: read_file
INPUT: test.py

To use another tool, respond EXACTLY in this format (two separate lines):
ACTION: <tool_name>
INPUT: <tool_input>

To provide your final response, write your answer directly.

Your response:`, input, action, result)
}

// findTool finds a tool by name
func (a *ConversationalAgent) findTool(name string) Tool {
	for _, tool := range a.GetTools() {
		if tool.Name() == name {
			return tool
		}
	}
	return nil
}

// Helper functions

func formatHistory(history []string) string {
	var formatted string
	for _, entry := range history {
		formatted += entry + "\n"
	}
	return formatted
}

// parseResponse parses the LLM response to extract action and input
func parseResponse(response string) (action, input string, err error) {
	debugPrint("Parsing response:\n%s", response)

	// Remove thinking process first
	response = cleanThinkingProcess(response)
	debugPrint("Response after cleaning thinking process:\n%s", response)

	// Check for ACTION/INPUT format
	lines := strings.Split(response, "\n")
	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		debugPrint("Checking line: %s", line)
		if strings.HasPrefix(line, "ACTION:") {
			action = strings.TrimSpace(strings.TrimPrefix(line, "ACTION:"))
			debugPrint("Found action: %s", action)
			// Look for INPUT on next line
			if i+1 < len(lines) {
				nextLine := strings.TrimSpace(lines[i+1])
				if strings.HasPrefix(nextLine, "INPUT:") {
					input = strings.TrimSpace(strings.TrimPrefix(nextLine, "INPUT:"))
					debugPrint("Found input: %s", input)
					return action, input, nil
				}
			}
		}
	}

	// If no action/input found, treat as direct response
	debugPrint("No action/input found, treating as direct response")
	return "", response, nil
}

// cleanThinkingProcess removes the thinking process from the response
func cleanThinkingProcess(response string) string {
	debugPrint("Cleaning thinking process from:\n%s", response)

	// Remove all <think>...</think> blocks
	for {
		start := strings.Index(response, "<think>")
		if start == -1 {
			break
		}
		end := strings.Index(response, "</think>")
		if end == -1 {
			break
		}
		debugPrint("Removing think block from %d to %d", start, end+8)
		response = response[:start] + response[end+8:]
	}

	// Clean up any extra whitespace
	response = strings.TrimSpace(response)
	// Replace multiple newlines with a single newline
	response = strings.ReplaceAll(response, "\n\n", "\n")

	debugPrint("After cleaning:\n%s", response)
	return response
}
