package agent

import (
	"context"
	"fmt"
	"strings"
)

// TaskType represents different types of tasks the orchestrator can handle
type TaskType string

const (
	TaskTypeInformationRetrieval TaskType = "information_retrieval"
	TaskTypeCostLookup           TaskType = "cost_lookup"
	TaskTypeResponseGeneration   TaskType = "response_generation"
	TaskTypeWebSearch            TaskType = "web_search"
	TaskTypeFileOperation        TaskType = "file_operation"
	TaskTypeCodeAnalysis         TaskType = "code_analysis"
)

// Task represents a single unit of work
type Task struct {
	ID           string
	Type         TaskType
	Description  string
	Input        string
	Dependencies []string // IDs of tasks that must complete before this one
	Result       string
	Error        error
	Status       TaskStatus
	Tool         string // The tool to use for this task
}

// TaskStatus represents the current status of a task
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)

// Orchestrator manages complex task decomposition and execution
type Orchestrator struct {
	llmClient LLMClient
	agent     Agent
	progress  *ProgressDisplay
}

// NewOrchestrator creates a new task orchestrator
func NewOrchestrator(llmClient LLMClient, agent Agent) *Orchestrator {
	return &Orchestrator{
		llmClient: llmClient,
		agent:     agent,
		progress:  NewProgressDisplay(Debug),
	}
}

// DecomposeQuery breaks down a complex query into smaller tasks
func (o *Orchestrator) DecomposeQuery(ctx context.Context, query string) ([]*Task, error) {
	o.progress.DebugPrint("Orchestrator: Decomposing query: %s", query)

	// Use LLM to analyze and decompose the query
	prompt := o.buildDecompositionPrompt(query)
	response, err := o.llmClient.GenerateCompletion(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to decompose query: %w", err)
	}

	// Parse the decomposition response
	tasks, err := o.parseDecompositionResponse(response)
	if err != nil {
		return nil, fmt.Errorf("failed to parse decomposition: %w", err)
	}

	o.progress.DebugPrint("Orchestrator: Decomposed into %d tasks", len(tasks))
	for i, task := range tasks {
		o.progress.DebugPrint("  Task %d: %s (%s) - Tool: %s", i+1, task.Description, task.Type, task.Tool)
	}

	return tasks, nil
}

// ExecuteTasks executes a list of tasks in the correct order
func (o *Orchestrator) ExecuteTasks(ctx context.Context, tasks []*Task) error {
	// Create progress display
	progress := NewProgressDisplay(Debug)

	// Initialize all tasks in progress display
	for _, task := range tasks {
		progress.StartTask(task.ID, task.Description)
	}

	// Execute tasks respecting dependencies
	completed := make(map[string]bool)

	for {
		// Find tasks that can be executed
		readyTasks := o.findReadyTasks(tasks, completed)
		if len(readyTasks) == 0 {
			// Check if all tasks are completed
			if len(completed) == len(tasks) {
				if !Debug {
					progress.DebugPrint("Orchestrator: All tasks completed")
				}
				break
			}
			// No tasks ready and not all completed - we have a problem
			return fmt.Errorf("no tasks ready to execute, possible circular dependency")
		}

		// Execute ready tasks
		for _, task := range readyTasks {
			if !Debug {
				progress.DebugPrint("Orchestrator: Executing task %s: %s", task.ID, task.Description)
			}

			task.Status = TaskStatusRunning
			progress.UpdateTaskStatus(task.ID, TaskProgressRunning)

			// Build context with results from dependencies
			taskContext := o.buildTaskContext(task, tasks)

			// Execute the task using the appropriate tool
			result, err := o.executeTask(ctx, task, taskContext)
			if err != nil {
				task.Status = TaskStatusFailed
				task.Error = err
				progress.CompleteTask(task.ID, "", err)
				if !Debug {
					progress.DebugPrint("Orchestrator: Task %s failed: %v", task.ID, err)
				}
				// Continue with other tasks that don't depend on this one
			} else {
				task.Status = TaskStatusCompleted
				task.Result = result
				completed[task.ID] = true
				progress.CompleteTask(task.ID, result, nil)
				if !Debug {
					progress.DebugPrint("Orchestrator: Task %s completed successfully", task.ID)
				}
			}
		}
	}

	return nil
}

// GenerateFinalResponse creates a comprehensive response from all task results
func (o *Orchestrator) GenerateFinalResponse(ctx context.Context, originalQuery string, tasks []*Task) (string, error) {
	o.progress.DebugPrint("Orchestrator: Generating final response")

	// Build a prompt with all task results
	prompt := o.buildFinalResponsePrompt(originalQuery, tasks)

	response, err := o.llmClient.GenerateCompletion(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("failed to generate final response: %w", err)
	}

	// Clean the response
	finalResponse := cleanThinkingProcess(response)
	o.progress.DebugPrint("Orchestrator: Final response generated")

	return finalResponse, nil
}

// buildDecompositionPrompt creates the prompt for task decomposition
func (o *Orchestrator) buildDecompositionPrompt(query string) string {
	tools := o.agent.GetTools()
	toolDescriptions := make([]string, 0, len(tools))
	availableToolNames := make([]string, 0, len(tools))

	for _, tool := range tools {
		// Skip rag_search tool from the available tools list
		if tool.Name() == "rag_search" {
			continue
		}
		toolDescriptions = append(toolDescriptions, fmt.Sprintf("- %s: %s", tool.Name(), tool.Description()))
		availableToolNames = append(availableToolNames, tool.Name())
	}

	// Check if web search is available
	hasWebSearch := false
	for _, toolName := range availableToolNames {
		if toolName == "web_search" {
			hasWebSearch = true
			break
		}
	}

	// Analyze the query to determine if web search is needed
	queryLower := strings.ToLower(query)
	needsWebSearch := strings.Contains(queryLower, "site web") ||
		strings.Contains(queryLower, "website") ||
		strings.Contains(queryLower, ".com") ||
		strings.Contains(queryLower, ".sh") ||
		strings.Contains(queryLower, ".org") ||
		strings.Contains(queryLower, "internet") ||
		strings.Contains(queryLower, "recherche en ligne") ||
		strings.Contains(queryLower, "online search") ||
		strings.Contains(queryLower, "va sur") ||
		strings.Contains(queryLower, "go to") ||
		strings.Contains(queryLower, "search") ||
		strings.Contains(queryLower, "look up") ||
		strings.Contains(queryLower, "find") ||
		strings.Contains(queryLower, "ollama models") ||
		strings.Contains(queryLower, "benchmarks") ||
		strings.Contains(queryLower, "compare")

	if needsWebSearch && !hasWebSearch {
		return fmt.Sprintf(`ERROR: This query requires web search functionality but web search is not enabled.

Query: %s

This query requires online information. To enable web search:
1. Add the -w flag: rlama agent run -w -q "your query"  
2. Set up Google API credentials (GOOGLE_SEARCH_API_KEY and GOOGLE_SEARCH_ENGINE_ID)

Please retry with web search enabled.`, query)
	}

	return fmt.Sprintf(`You are a task orchestrator. Your job is to decompose complex queries into smaller, manageable tasks.

Available tools:
%s

User query: %s

IMPORTANT RULES FOR TOOL SELECTION:
- web_search: Use for ALL information gathering, research, and lookup tasks
- list_dir/read_file: Use ONLY when explicitly working with local files in current directory
- file_search/grep_search: Use ONLY for searching within known local files
- NEVER create tasks for non-existent local files or directories

CONTEXT ANALYSIS: For queries about external information (models, benchmarks, comparisons, etc.), ALWAYS use web_search.

Analyze this query and break it down into specific tasks. Each task should:
1. Have a unique ID (task1, task2, etc.) - NO DUPLICATES!
2. Have a clear type (information_retrieval, web_search, etc.)
3. Have a specific description  
4. Specify which tool to use BASED ON CONTEXT
5. List any dependencies on other tasks

CRITICAL RULES:
- Each task must have a UNIQUE ID and PURPOSE
- Do NOT create duplicate tasks
- Maximum 2-3 tasks per query (keep it simple!)
- Always end with ONE response_generation task
- For information gathering: ALWAYS use web_search
- Do NOT assume local files exist

Respond in this EXACT format:
TASK: task1
TYPE: web_search
DESCRIPTION: Search for Ollama models with tool support and benchmarks
TOOL: web_search
INPUT: Ollama models tool support benchmarks performance comparison
DEPENDENCIES: none

TASK: task2
TYPE: response_generation
DESCRIPTION: Create comparison table of models with tool efficiency
TOOL: none
INPUT: Create a table comparing Ollama models with their tool efficiency and benchmarks
DEPENDENCIES: task1

Important: 
- For information gathering: ALWAYS use web_search
- Keep task count low (2-3 tasks maximum)
- The final task should always be TYPE: response_generation
- Do NOT use file operations unless explicitly working with local files`, strings.Join(toolDescriptions, "\n"), query)
}

// parseDecompositionResponse parses the LLM's task decomposition
func (o *Orchestrator) parseDecompositionResponse(response string) ([]*Task, error) {
	tasks := make([]*Task, 0)
	lines := strings.Split(response, "\n")
	seenTasks := make(map[string]bool) // Track duplicate task IDs

	var currentTask *Task
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "TASK:") {
			// Save previous task if exists and is valid
			if currentTask != nil && o.isValidTask(currentTask) {
				// Check for duplicates
				if !seenTasks[currentTask.ID] {
					tasks = append(tasks, currentTask)
					seenTasks[currentTask.ID] = true
				} else {
					o.progress.DebugPrint("Orchestrator: Skipping duplicate task %s", currentTask.ID)
				}
			}
			// Start new task
			taskID := strings.TrimSpace(strings.TrimPrefix(line, "TASK:"))
			if taskID != "" {
				currentTask = &Task{
					ID:           taskID,
					Status:       TaskStatusPending,
					Dependencies: []string{},
				}
			}
		} else if currentTask != nil {
			if strings.HasPrefix(line, "TYPE:") {
				typeStr := strings.TrimSpace(strings.TrimPrefix(line, "TYPE:"))
				currentTask.Type = TaskType(typeStr)
			} else if strings.HasPrefix(line, "DESCRIPTION:") {
				currentTask.Description = strings.TrimSpace(strings.TrimPrefix(line, "DESCRIPTION:"))
			} else if strings.HasPrefix(line, "TOOL:") {
				toolName := strings.TrimSpace(strings.TrimPrefix(line, "TOOL:"))
				// Skip rag_search tool
				if toolName == "rag_search" {
					toolName = "web_search" // Replace with web_search
				}
				currentTask.Tool = toolName
			} else if strings.HasPrefix(line, "INPUT:") {
				currentTask.Input = strings.TrimSpace(strings.TrimPrefix(line, "INPUT:"))
			} else if strings.HasPrefix(line, "DEPENDENCIES:") {
				deps := strings.TrimSpace(strings.TrimPrefix(line, "DEPENDENCIES:"))
				if deps != "none" && deps != "" {
					currentTask.Dependencies = strings.Split(deps, ",")
					for i := range currentTask.Dependencies {
						currentTask.Dependencies[i] = strings.TrimSpace(currentTask.Dependencies[i])
					}
				}
			}
		}
	}

	// Don't forget the last task if it's valid
	if currentTask != nil && o.isValidTask(currentTask) {
		if !seenTasks[currentTask.ID] {
			tasks = append(tasks, currentTask)
			seenTasks[currentTask.ID] = true
		}
	}

	if len(tasks) == 0 {
		return nil, fmt.Errorf("no valid tasks found in decomposition")
	}

	// Clean up dependencies to only reference existing tasks
	existingTaskIDs := make(map[string]bool)
	for _, task := range tasks {
		existingTaskIDs[task.ID] = true
	}

	for _, task := range tasks {
		validDeps := []string{}
		for _, dep := range task.Dependencies {
			if existingTaskIDs[dep] {
				validDeps = append(validDeps, dep)
			} else {
				o.progress.DebugPrint("Orchestrator: Removing invalid dependency %s from task %s", dep, task.ID)
			}
		}
		task.Dependencies = validDeps
	}

	o.progress.DebugPrint("Orchestrator: Parsed %d valid unique tasks", len(tasks))
	return tasks, nil
}

// isValidTask checks if a task has the minimum required fields
func (o *Orchestrator) isValidTask(task *Task) bool {
	if task == nil {
		return false
	}
	if task.ID == "" {
		o.progress.DebugPrint("Orchestrator: Invalid task - missing ID")
		return false
	}
	if task.Description == "" {
		o.progress.DebugPrint("Orchestrator: Invalid task %s - missing description", task.ID)
		return false
	}
	if task.Tool == "" && task.Type != TaskTypeResponseGeneration {
		o.progress.DebugPrint("Orchestrator: Invalid task %s - missing tool", task.ID)
		return false
	}
	// Skip validation for rag_search since we're removing it
	if task.Tool == "rag_search" {
		o.progress.DebugPrint("Orchestrator: Invalid task %s - rag_search not available", task.ID)
		return false
	}
	return true
}

// findReadyTasks finds tasks that are ready to execute
func (o *Orchestrator) findReadyTasks(tasks []*Task, completed map[string]bool) []*Task {
	ready := make([]*Task, 0)

	for _, task := range tasks {
		if task.Status != TaskStatusPending {
			continue
		}

		// Check if all dependencies are completed
		allDepsCompleted := true
		for _, dep := range task.Dependencies {
			if !completed[dep] {
				allDepsCompleted = false
				break
			}
		}

		if allDepsCompleted {
			ready = append(ready, task)
		}
	}

	return ready
}

// buildTaskContext builds context for a task including results from dependencies
func (o *Orchestrator) buildTaskContext(task *Task, allTasks []*Task) string {
	context := fmt.Sprintf("Task: %s\n", task.Description)

	if len(task.Dependencies) > 0 {
		context += "\nContext from previous tasks:\n"
		for _, depID := range task.Dependencies {
			for _, t := range allTasks {
				if t.ID == depID && t.Status == TaskStatusCompleted {
					context += fmt.Sprintf("- %s: %s\n", t.Description, t.Result)
				}
			}
		}
	}

	return context
}

// executeTask executes a single task using the appropriate tool
func (o *Orchestrator) executeTask(ctx context.Context, task *Task, taskContext string) (string, error) {
	// For response generation tasks, we need to use the LLM directly
	if task.Type == TaskTypeResponseGeneration || task.Tool == "llm" || task.Tool == "none" || task.Tool == "" {
		prompt := fmt.Sprintf("%s\n\nBased on the above context, %s", taskContext, task.Input)
		response, err := o.llmClient.GenerateCompletion(ctx, prompt)
		if err != nil {
			return "", err
		}
		return cleanThinkingProcess(response), nil
	}

	// Find the specified tool
	var tool Tool
	for _, t := range o.agent.GetTools() {
		if t.Name() == task.Tool {
			tool = t
			break
		}
	}

	if tool == nil {
		// Handle missing tools gracefully
		if task.Tool == "web_search" {
			return "Web search is not available. Please enable web search with the -w flag and provide your Google API credentials to search for real-time information.", nil
		}
		if task.Tool == "rag_search" {
			return "RAG search is not available in this context. Please use web search for external information.", nil
		}
		return "", fmt.Errorf("tool %s not found", task.Tool)
	}

	// Execute the tool with the task input
	result, err := tool.Execute(ctx, task.Input)
	if err != nil {
		return "", fmt.Errorf("tool execution failed: %w", err)
	}

	return result, nil
}

// buildFinalResponsePrompt creates the prompt for generating the final response
func (o *Orchestrator) buildFinalResponsePrompt(originalQuery string, tasks []*Task) string {
	taskResults := ""
	for _, task := range tasks {
		if task.Status == TaskStatusCompleted {
			taskResults += fmt.Sprintf("\n%s:\n%s\n", task.Description, task.Result)
		} else if task.Status == TaskStatusFailed {
			taskResults += fmt.Sprintf("\n%s: FAILED - %v\n", task.Description, task.Error)
		}
	}

	return fmt.Sprintf(`Based on the following task results, provide a comprehensive answer to the user's original query.

Original query: %s

Task results:%s

Please provide a clear, well-structured response that:
1. Directly answers all parts of the user's question
2. Includes all relevant information gathered
3. Is formatted in a user-friendly way
4. Mentions if any information couldn't be found

Your response:`, originalQuery, taskResults)
}
