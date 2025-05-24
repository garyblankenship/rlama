package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/dontizi/rlama/internal/domain"
)

// WebSearchTool provides web search capability
type WebSearchTool struct {
	client         *http.Client
	apiKey         string
	searchEngineID string
}

// NewWebSearchTool creates a new WebSearchTool
func NewWebSearchTool(apiKey string) *WebSearchTool {
	// Get Search Engine ID from environment if not provided
	searchEngineID := os.Getenv("GOOGLE_SEARCH_ENGINE_ID")
	if searchEngineID == "" {
		searchEngineID = "YOUR_SEARCH_ENGINE_ID" // Default fallback
	}

	return &WebSearchTool{
		client:         &http.Client{},
		apiKey:         apiKey,
		searchEngineID: searchEngineID,
	}
}

func (t *WebSearchTool) Name() string {
	return "web_search"
}

func (t *WebSearchTool) Description() string {
	return "Search the web for real-time information about any topic. Use this tool when you need up-to-date information that might not be available in your training data, or when you need to verify current facts. The search results will include relevant snippets and URLs from web pages. This is particularly useful for questions about current events, technology updates, or any topic that requires recent information."
}

func (t *WebSearchTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"explanation": map[string]interface{}{
				"description": "One sentence explanation as to why this tool is being used, and how it contributes to the goal.",
				"type":        "string",
			},
			"search_term": map[string]interface{}{
				"description": "The search term to look up on the web. Be specific and include relevant keywords for better results. For technical queries, include version numbers or dates if relevant.",
				"type":        "string",
			},
		},
		"required": []string{"search_term"},
	}
}

func (t *WebSearchTool) ExecuteWithParams(ctx context.Context, params map[string]interface{}) (string, error) {
	searchTerm, ok := params["search_term"].(string)
	if !ok {
		return "", fmt.Errorf("search_term parameter is required and must be a string")
	}
	return t.Execute(ctx, searchTerm)
}

func (t *WebSearchTool) Execute(ctx context.Context, input string) (string, error) {
	if t.apiKey == "" {
		return "", fmt.Errorf("Google Search API key not provided")
	}

	if t.searchEngineID == "" || t.searchEngineID == "YOUR_SEARCH_ENGINE_ID" {
		return "", fmt.Errorf("Google Search Engine ID not provided (set GOOGLE_SEARCH_ENGINE_ID environment variable)")
	}

	// Use Google Custom Search API
	searchURL := "https://www.googleapis.com/customsearch/v1"

	// Build query parameters
	params := url.Values{}
	params.Add("key", t.apiKey)
	params.Add("cx", t.searchEngineID)
	params.Add("q", input)
	params.Add("num", "5") // Request 5 results

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL+"?"+params.Encode(), nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request
	resp, err := t.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	// Check for error response
	if resp.StatusCode != http.StatusOK {
		var errorResp map[string]interface{}
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
		}
		if errObj, ok := errorResp["error"].(map[string]interface{}); ok {
			return "", fmt.Errorf("API error: %v", errObj["message"])
		}
		return "", fmt.Errorf("API error (status %d)", resp.StatusCode)
	}

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Extract search results
	items, ok := result["items"].([]interface{})
	if !ok || len(items) == 0 {
		return "", fmt.Errorf("no results found")
	}

	// Build response text from top results
	var response strings.Builder
	response.WriteString("Search results:\n\n")

	maxResults := 3
	if len(items) < maxResults {
		maxResults = len(items)
	}

	for i := 0; i < maxResults; i++ {
		item := items[i].(map[string]interface{})
		title := item["title"].(string)
		snippet := item["snippet"].(string)
		link := item["link"].(string)
		response.WriteString(fmt.Sprintf("- %s\n%s\nSource: %s\n\n", title, snippet, link))
	}

	return response.String(), nil
}

// FileTool provides file operation capability
type FileTool struct {
	baseDir string
	mode    string // "read" or "write"
}

// NewFileTool creates a new FileTool
func NewFileTool(baseDir string, mode string) *FileTool {
	return &FileTool{
		baseDir: baseDir,
		mode:    mode,
	}
}

func (t *FileTool) Name() string {
	if t.mode == "read" {
		return "read_file"
	}
	return fmt.Sprintf("file_%s", t.mode)
}

func (t *FileTool) Description() string {
	if t.mode == "read" {
		return "Read the contents of a file. the output of this tool call will be the 1-indexed file contents from start_line_one_indexed to end_line_one_indexed_inclusive, together with a summary of the lines outside start_line_one_indexed and end_line_one_indexed_inclusive. Note that this call can view at most 1500 lines at a time and 500 lines minimum. When using this tool to gather information, it's your responsibility to ensure you have the COMPLETE context. Specifically, each time you call this command you should: 1) Assess if the contents you viewed are sufficient to proceed with your task. 2) Take note of where there are lines not shown. 3) If the file contents you have viewed are insufficient, and you suspect they may be in lines not shown, proactively call the tool again to view those lines. 4) When in doubt, call this tool again to gather more information. Remember that partial file views may miss critical dependencies, imports, or functionality. In some cases, if reading a range of lines is not enough, you may choose to read the entire file. Reading entire files is often wasteful and slow, especially for large files (i.e. more than a few hundred lines). So you should use this option sparingly. Reading the entire file is not allowed in most cases. You are only allowed to read the entire file if it has been edited or manually attached to the conversation by the user."
	}
	return "Writes content to a file. Input should be in format: 'path:content'"
}

func (t *FileTool) Schema() map[string]interface{} {
	if t.mode == "read" {
		return map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"explanation": map[string]interface{}{
					"description": "One sentence explanation as to why this tool is being used, and how it contributes to the goal.",
					"type":        "string",
				},
				"target_file": map[string]interface{}{
					"description": "The path of the file to read. You can use either a relative path in the workspace or an absolute path. If an absolute path is provided, it will be preserved as is.",
					"type":        "string",
				},
				"should_read_entire_file": map[string]interface{}{
					"description": "Whether to read the entire file. Defaults to false.",
					"type":        "boolean",
				},
				"start_line_one_indexed": map[string]interface{}{
					"description": "The one-indexed line number to start reading from (inclusive).",
					"type":        "integer",
				},
				"end_line_one_indexed_inclusive": map[string]interface{}{
					"description": "The one-indexed line number to end reading at (inclusive).",
					"type":        "integer",
				},
			},
			"required": []string{"target_file", "should_read_entire_file", "start_line_one_indexed", "end_line_one_indexed_inclusive"},
		}
	}
	// Schema for write mode
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"target_file": map[string]interface{}{
				"description": "The path of the file to write to.",
				"type":        "string",
			},
			"content": map[string]interface{}{
				"description": "The content to write to the file.",
				"type":        "string",
			},
		},
		"required": []string{"target_file", "content"},
	}
}

func (t *FileTool) ExecuteWithParams(ctx context.Context, params map[string]interface{}) (string, error) {
	if t.mode == "read" {
		targetFile, ok := params["target_file"].(string)
		if !ok {
			return "", fmt.Errorf("target_file parameter is required and must be a string")
		}

		shouldReadEntire, ok := params["should_read_entire_file"].(bool)
		if !ok {
			shouldReadEntire = false
		}

		if shouldReadEntire {
			return t.Execute(ctx, targetFile)
		}

		startLine, ok := params["start_line_one_indexed"].(float64)
		if !ok {
			return "", fmt.Errorf("start_line_one_indexed parameter is required and must be an integer")
		}

		endLine, ok := params["end_line_one_indexed_inclusive"].(float64)
		if !ok {
			return "", fmt.Errorf("end_line_one_indexed_inclusive parameter is required and must be an integer")
		}

		return t.readFileRange(ctx, targetFile, int(startLine), int(endLine))
	}

	// Write mode
	targetFile, ok := params["target_file"].(string)
	if !ok {
		return "", fmt.Errorf("target_file parameter is required and must be a string")
	}

	content, ok := params["content"].(string)
	if !ok {
		return "", fmt.Errorf("content parameter is required and must be a string")
	}

	input := fmt.Sprintf("%s:%s", targetFile, content)
	return t.Execute(ctx, input)
}

// readFileRange reads a specific range of lines from a file
func (t *FileTool) readFileRange(ctx context.Context, targetFile string, startLine, endLine int) (string, error) {
	// Ensure the path is within the allowed base directory
	path := filepath.Join(t.baseDir, targetFile)
	if !filepath.HasPrefix(path, t.baseDir) {
		return "", fmt.Errorf("access to path outside base directory is not allowed")
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	totalLines := len(lines)

	// Validate line numbers
	if startLine < 1 {
		startLine = 1
	}
	if endLine > totalLines {
		endLine = totalLines
	}
	if startLine > endLine {
		return "", fmt.Errorf("start_line (%d) cannot be greater than end_line (%d)", startLine, endLine)
	}

	// Build response with context
	var response strings.Builder
	response.WriteString(fmt.Sprintf("Contents of %s, lines %d-%d (total %d lines):\n\n", targetFile, startLine, endLine, totalLines))

	// Add summary of lines before if any
	if startLine > 1 {
		response.WriteString(fmt.Sprintf("... %d lines before ...\n\n", startLine-1))
	}

	// Add the requested lines with line numbers
	for i := startLine - 1; i < endLine; i++ {
		if i < len(lines) {
			response.WriteString(fmt.Sprintf("%d: %s\n", i+1, lines[i]))
		}
	}

	// Add summary of lines after if any
	if endLine < totalLines {
		response.WriteString(fmt.Sprintf("\n... %d lines after ...\n", totalLines-endLine))
	}

	return response.String(), nil
}

func (t *FileTool) Execute(ctx context.Context, input string) (string, error) {
	// Ensure the path is within the allowed base directory
	path := filepath.Join(t.baseDir, input)
	if !filepath.HasPrefix(path, t.baseDir) {
		return "", fmt.Errorf("access to path outside base directory is not allowed")
	}

	switch t.mode {
	case "read":
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("failed to read file: %w", err)
		}
		return string(content), nil

	case "write":
		// Parse path:content format
		parts := strings.SplitN(input, ":", 2)
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid input format, expected 'path:content'")
		}

		filePath := filepath.Join(t.baseDir, parts[0])
		content := parts[1]

		err := ioutil.WriteFile(filePath, []byte(content), 0644)
		if err != nil {
			return "", fmt.Errorf("failed to write file: %w", err)
		}
		return "File written successfully", nil

	default:
		return "", fmt.Errorf("unsupported file operation mode: %s", t.mode)
	}
}

// RagSearchTool provides RAG search capability
type RagSearchTool struct {
	ragService RagService
}

// NewRagSearchTool creates a new RagSearchTool
func NewRagSearchTool(ragService RagService) *RagSearchTool {
	return &RagSearchTool{
		ragService: ragService,
	}
}

func (t *RagSearchTool) Name() string {
	return "rag_search"
}

func (t *RagSearchTool) Description() string {
	return "Searches the local document base using RAG. Input should be a search query."
}

func (t *RagSearchTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"explanation": map[string]interface{}{
				"description": "One sentence explanation as to why this tool is being used, and how it contributes to the goal.",
				"type":        "string",
			},
			"query": map[string]interface{}{
				"description": "The search query for the local document base.",
				"type":        "string",
			},
		},
		"required": []string{"query"},
	}
}

func (t *RagSearchTool) ExecuteWithParams(ctx context.Context, params map[string]interface{}) (string, error) {
	query, ok := params["query"].(string)
	if !ok {
		return "", fmt.Errorf("query parameter is required and must be a string")
	}
	return t.Execute(ctx, query)
}

func (t *RagSearchTool) Execute(ctx context.Context, input string) (string, error) {
	result, err := t.ragService.Query(ctx, input)
	if err != nil {
		return "", fmt.Errorf("RAG search failed: %w", err)
	}
	return result, nil
}

// DirectoryListTool provides directory listing capability (list_dir in Cursor)
type DirectoryListTool struct {
	baseDir string
}

// NewDirectoryListTool creates a new DirectoryListTool
func NewDirectoryListTool(baseDir string) *DirectoryListTool {
	return &DirectoryListTool{
		baseDir: baseDir,
	}
}

func (t *DirectoryListTool) Name() string {
	return "list_dir"
}

func (t *DirectoryListTool) Description() string {
	return "List the contents of a directory. The quick tool to use for discovery, before using more targeted tools like semantic search or file reading. Useful to try to understand the file structure before diving deeper into specific files. Can be used to explore the codebase."
}

func (t *DirectoryListTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"explanation": map[string]interface{}{
				"description": "One sentence explanation as to why this tool is being used, and how it contributes to the goal.",
				"type":        "string",
			},
			"relative_workspace_path": map[string]interface{}{
				"description": "Path to list contents of, relative to the workspace root.",
				"type":        "string",
			},
		},
		"required": []string{"relative_workspace_path"},
	}
}

func (t *DirectoryListTool) ExecuteWithParams(ctx context.Context, params map[string]interface{}) (string, error) {
	path, ok := params["relative_workspace_path"].(string)
	if !ok {
		return "", fmt.Errorf("relative_workspace_path parameter is required and must be a string")
	}
	return t.Execute(ctx, path)
}

func (t *DirectoryListTool) Execute(ctx context.Context, input string) (string, error) {
	// Handle current directory
	if input == "" || input == "." {
		input = "."
	}

	// Ensure the path is within the allowed base directory
	dirPath := filepath.Join(t.baseDir, input)
	if !filepath.HasPrefix(dirPath, t.baseDir) {
		return "", fmt.Errorf("access to path outside base directory is not allowed")
	}

	// Check if directory exists
	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		return "", fmt.Errorf("failed to access directory: %w", err)
	}

	if !dirInfo.IsDir() {
		return "", fmt.Errorf("path is not a directory")
	}

	// Read directory contents
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %w", err)
	}

	// Build formatted response
	var response strings.Builder
	response.WriteString(fmt.Sprintf("Contents of directory: %s\n\n", input))

	if len(entries) == 0 {
		response.WriteString("Directory is empty.\n")
		return response.String(), nil
	}

	// Group by type and format output
	var files []os.DirEntry
	var dirs []os.DirEntry

	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry)
		} else {
			files = append(files, entry)
		}
	}

	// List directories first
	for _, dir := range dirs {
		response.WriteString(fmt.Sprintf("[dir]  %s/\n", dir.Name()))
	}

	// Then list files with sizes
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			response.WriteString(fmt.Sprintf("[file] %s (size unknown)\n", file.Name()))
			continue
		}

		size := info.Size()
		var sizeStr string
		if size >= 1024*1024 {
			sizeStr = fmt.Sprintf("%.1fMB", float64(size)/(1024*1024))
		} else if size >= 1024 {
			sizeStr = fmt.Sprintf("%.1fKB", float64(size)/1024)
		} else {
			sizeStr = fmt.Sprintf("%dB", size)
		}

		response.WriteString(fmt.Sprintf("[file] %s (%s)\n", file.Name(), sizeStr))
	}

	return response.String(), nil
}

// GrepSearchTool provides exact/regex search capability like ripgrep
type GrepSearchTool struct {
	baseDir string
}

// NewGrepSearchTool creates a new GrepSearchTool
func NewGrepSearchTool(baseDir string) *GrepSearchTool {
	return &GrepSearchTool{
		baseDir: baseDir,
	}
}

func (t *GrepSearchTool) Name() string {
	return "grep_search"
}

func (t *GrepSearchTool) Description() string {
	return "### Instructions: This is best for finding exact text matches or regex patterns. This is preferred over semantic search when we know the exact symbol/function name/etc. to search in some set of directories/file types. Use this tool to run fast, exact regex searches over text files using the `ripgrep` engine. To avoid overwhelming output, the results are capped at 50 matches. Use the include or exclude patterns to filter the search scope by file type or specific paths. - Always escape special regex characters: ( ) [ ] { } + * ? ^ $ | \\ - Use `\\` to escape any of these characters when they appear in your search string. - Do NOT perform fuzzy or semantic matches. - Return only a valid regex pattern string."
}

func (t *GrepSearchTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"explanation": map[string]interface{}{
				"description": "One sentence explanation as to why this tool is being used, and how it contributes to the goal.",
				"type":        "string",
			},
			"query": map[string]interface{}{
				"description": "The regex pattern to search for",
				"type":        "string",
			},
			"case_sensitive": map[string]interface{}{
				"description": "Whether the search should be case sensitive",
				"type":        "boolean",
			},
			"include_pattern": map[string]interface{}{
				"description": "Glob pattern for files to include (e.g. '*.ts' for TypeScript files)",
				"type":        "string",
			},
			"exclude_pattern": map[string]interface{}{
				"description": "Glob pattern for files to exclude",
				"type":        "string",
			},
		},
		"required": []string{"query"},
	}
}

func (t *GrepSearchTool) ExecuteWithParams(ctx context.Context, params map[string]interface{}) (string, error) {
	query, ok := params["query"].(string)
	if !ok {
		return "", fmt.Errorf("query parameter is required and must be a string")
	}

	caseSensitive, _ := params["case_sensitive"].(bool)
	includePattern, _ := params["include_pattern"].(string)
	excludePattern, _ := params["exclude_pattern"].(string)

	if includePattern == "" {
		includePattern = "*"
	}

	return t.executeWithPatterns(ctx, query, includePattern, excludePattern, caseSensitive)
}

// executeWithPatterns implements the enhanced grep search with include/exclude patterns
func (t *GrepSearchTool) executeWithPatterns(ctx context.Context, pattern, includePattern, excludePattern string, caseSensitive bool) (string, error) {
	var results []string
	maxResults := 50

	err := filepath.Walk(t.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue on errors
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check include pattern
		if includePattern != "" && includePattern != "*" {
			matched, err := filepath.Match(includePattern, info.Name())
			if err != nil || !matched {
				return nil
			}
		}

		// Check exclude pattern
		if excludePattern != "" {
			matched, err := filepath.Match(excludePattern, info.Name())
			if err == nil && matched {
				return nil // Skip excluded files
			}
		}

		// Read file content
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return nil // Continue on read errors
		}

		// Search for pattern
		lines := strings.Split(string(content), "\n")
		for lineNum, line := range lines {
			var found bool
			if caseSensitive {
				found = strings.Contains(line, pattern)
			} else {
				found = strings.Contains(strings.ToLower(line), strings.ToLower(pattern))
			}

			if found {
				relPath, _ := filepath.Rel(t.baseDir, path)
				results = append(results, fmt.Sprintf("%s:%d:%s", relPath, lineNum+1, strings.TrimSpace(line)))
				if len(results) >= maxResults {
					return filepath.SkipDir
				}
			}
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("search failed: %w", err)
	}

	if len(results) == 0 {
		return "No matches found.", nil
	}

	response := fmt.Sprintf("Found %d matches (showing up to %d):\n\n", len(results), maxResults)
	response += strings.Join(results, "\n")

	return response, nil
}

func (t *GrepSearchTool) Execute(ctx context.Context, input string) (string, error) {
	parts := strings.Split(input, "|")
	if len(parts) < 1 || len(parts) > 3 {
		return "", fmt.Errorf("invalid input format. Use: 'pattern|file_pattern|case_sensitive' or just 'pattern'")
	}

	pattern := parts[0]
	filePattern := "*"
	caseSensitive := false

	if len(parts) > 1 && parts[1] != "" {
		filePattern = parts[1]
	}
	if len(parts) > 2 && parts[2] == "true" {
		caseSensitive = true
	}

	var results []string
	maxResults := 50

	err := filepath.Walk(t.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue on errors
		}

		// Skip directories and check file pattern
		if info.IsDir() {
			return nil
		}

		matched, err := filepath.Match(filePattern, info.Name())
		if err != nil || !matched {
			return nil
		}

		// Read file content
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return nil // Continue on read errors
		}

		// Search for pattern
		lines := strings.Split(string(content), "\n")
		for lineNum, line := range lines {
			var found bool
			if caseSensitive {
				found = strings.Contains(line, pattern)
			} else {
				found = strings.Contains(strings.ToLower(line), strings.ToLower(pattern))
			}

			if found {
				relPath, _ := filepath.Rel(t.baseDir, path)
				results = append(results, fmt.Sprintf("%s:%d:%s", relPath, lineNum+1, strings.TrimSpace(line)))
				if len(results) >= maxResults {
					return filepath.SkipDir
				}
			}
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("search failed: %w", err)
	}

	if len(results) == 0 {
		return "No matches found.", nil
	}

	response := fmt.Sprintf("Found %d matches (showing up to %d):\n\n", len(results), maxResults)
	response += strings.Join(results, "\n")

	return response, nil
}

// FileSearchTool provides fuzzy file name search
type FileSearchTool struct {
	baseDir string
}

// NewFileSearchTool creates a new FileSearchTool
func NewFileSearchTool(baseDir string) *FileSearchTool {
	return &FileSearchTool{
		baseDir: baseDir,
	}
}

func (t *FileSearchTool) Name() string {
	return "file_search"
}

func (t *FileSearchTool) Description() string {
	return "Fast file search based on fuzzy matching against file path. Use if you know part of the file path but don't know where it's located exactly. Response will be capped to 10 results. Make your query more specific if need to filter results further."
}

func (t *FileSearchTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"explanation": map[string]interface{}{
				"description": "One sentence explanation as to why this tool is being used, and how it contributes to the goal.",
				"type":        "string",
			},
			"query": map[string]interface{}{
				"description": "Fuzzy filename to search for",
				"type":        "string",
			},
		},
		"required": []string{"query", "explanation"},
	}
}

func (t *FileSearchTool) ExecuteWithParams(ctx context.Context, params map[string]interface{}) (string, error) {
	query, ok := params["query"].(string)
	if !ok {
		return "", fmt.Errorf("query parameter is required and must be a string")
	}
	return t.Execute(ctx, query)
}

func (t *FileSearchTool) Execute(ctx context.Context, input string) (string, error) {
	if strings.TrimSpace(input) == "" {
		return "", fmt.Errorf("search query cannot be empty")
	}

	query := strings.ToLower(input)
	var matches []string
	maxResults := 10

	err := filepath.Walk(t.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue on errors
		}

		// Skip hidden directories like .git
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		// Only match files, not directories
		if !info.IsDir() {
			fileName := strings.ToLower(info.Name())

			// Simple fuzzy matching: check if all characters of query appear in order
			if fuzzyMatch(fileName, query) {
				relPath, _ := filepath.Rel(t.baseDir, path)
				size := info.Size()
				var sizeStr string
				if size >= 1024*1024 {
					sizeStr = fmt.Sprintf("%.1fMB", float64(size)/(1024*1024))
				} else if size >= 1024 {
					sizeStr = fmt.Sprintf("%.1fKB", float64(size)/1024)
				} else {
					sizeStr = fmt.Sprintf("%dB", size)
				}

				matches = append(matches, fmt.Sprintf("%s (%s)", relPath, sizeStr))
				if len(matches) >= maxResults {
					return filepath.SkipDir
				}
			}
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("file search failed: %w", err)
	}

	if len(matches) == 0 {
		return fmt.Sprintf("No files found matching '%s'", input), nil
	}

	response := fmt.Sprintf("Found %d files matching '%s':\n\n", len(matches), input)
	response += strings.Join(matches, "\n")

	return response, nil
}

// fuzzyMatch checks if all characters in query appear in target in order
func fuzzyMatch(target, query string) bool {
	targetIndex := 0
	for _, char := range query {
		found := false
		for targetIndex < len(target) {
			if rune(target[targetIndex]) == char {
				found = true
				targetIndex++
				break
			}
			targetIndex++
		}
		if !found {
			return false
		}
	}
	return true
}

// CodebaseSearchTool provides semantic search capability (simplified version)
type CodebaseSearchTool struct {
	baseDir string
}

// NewCodebaseSearchTool creates a new CodebaseSearchTool
func NewCodebaseSearchTool(baseDir string) *CodebaseSearchTool {
	return &CodebaseSearchTool{
		baseDir: baseDir,
	}
}

func (t *CodebaseSearchTool) Name() string {
	return "codebase_search"
}

func (t *CodebaseSearchTool) Description() string {
	return "Find snippets of code from the codebase most relevant to the search query. This is a semantic search tool, so the query should ask for something semantically matching what is needed. If it makes sense to only search in particular directories, please specify them in the target_directories field. Unless there is a clear reason to use your own search query, please just reuse the user's exact query with their wording. Their exact wording/phrasing can often be helpful for the semantic search query. Keeping the same exact question format can also be helpful."
}

func (t *CodebaseSearchTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"explanation": map[string]interface{}{
				"description": "One sentence explanation as to why this tool is being used, and how it contributes to the goal.",
				"type":        "string",
			},
			"query": map[string]interface{}{
				"description": "The search query to find relevant code. You should reuse the user's exact query/most recent message with their wording unless there is a clear reason not to.",
				"type":        "string",
			},
			"target_directories": map[string]interface{}{
				"description": "Glob patterns for directories to search over",
				"type":        "array",
				"items": map[string]interface{}{
					"type": "string",
				},
			},
		},
		"required": []string{"query"},
	}
}

func (t *CodebaseSearchTool) ExecuteWithParams(ctx context.Context, params map[string]interface{}) (string, error) {
	query, ok := params["query"].(string)
	if !ok {
		return "", fmt.Errorf("query parameter is required and must be a string")
	}

	// Handle target_directories if provided
	targetDirs, ok := params["target_directories"].([]interface{})
	if ok && len(targetDirs) > 0 {
		// For now, we'll use the first directory as a filter
		// In a more advanced implementation, we would search only in those directories
		firstDir, ok := targetDirs[0].(string)
		if ok {
			return t.executeInDirectory(ctx, query, firstDir)
		}
	}

	return t.Execute(ctx, query)
}

// executeInDirectory searches for code only in a specific directory
func (t *CodebaseSearchTool) executeInDirectory(ctx context.Context, query, targetDir string) (string, error) {
	// This is a simplified implementation that filters results by directory
	// A more advanced version would modify the search logic itself
	result, err := t.Execute(ctx, query)
	if err != nil {
		return "", err
	}

	// Filter results to only include files from the target directory
	lines := strings.Split(result, "\n")
	var filteredLines []string

	for _, line := range lines {
		if strings.Contains(line, targetDir) || !strings.Contains(line, "/") {
			filteredLines = append(filteredLines, line)
		}
	}

	if len(filteredLines) == 0 {
		return fmt.Sprintf("No relevant code found for '%s' in directory '%s'.", query, targetDir), nil
	}

	return strings.Join(filteredLines, "\n"), nil
}

func (t *CodebaseSearchTool) Execute(ctx context.Context, input string) (string, error) {
	if strings.TrimSpace(input) == "" {
		return "", fmt.Errorf("search query cannot be empty")
	}

	// Extract keywords from the query
	keywords := extractKeywords(input)
	var results []string
	maxResults := 20

	err := filepath.Walk(t.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip directories and non-code files
		if info.IsDir() {
			// Skip hidden directories
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		// Only search in code files
		ext := strings.ToLower(filepath.Ext(path))
		if !isCodeFile(ext) {
			return nil
		}

		// Read and search file content
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}

		lines := strings.Split(string(content), "\n")
		contentLower := strings.ToLower(string(content))

		// Check if file is relevant (contains keywords)
		relevantLines := []string{}
		for lineNum, line := range lines {
			lineLower := strings.ToLower(line)
			matchCount := 0

			for _, keyword := range keywords {
				if strings.Contains(lineLower, keyword) {
					matchCount++
				}
			}

			// If line contains multiple keywords or important patterns, include it
			if matchCount > 0 || containsImportantPatterns(lineLower, keywords) {
				relPath, _ := filepath.Rel(t.baseDir, path)
				relevantLines = append(relevantLines, fmt.Sprintf("%s:%d:%s", relPath, lineNum+1, strings.TrimSpace(line)))
			}
		}

		// Add file to results if it contains relevant content
		if len(relevantLines) > 0 {
			// Calculate relevance score
			score := calculateRelevanceScore(contentLower, keywords)
			if score > 0 {
				relPath, _ := filepath.Rel(t.baseDir, path)
				fileResult := fmt.Sprintf("\n📁 %s (relevance: %.2f)\n", relPath, score)

				// Add most relevant lines (max 5 per file)
				maxLines := 5
				if len(relevantLines) > maxLines {
					relevantLines = relevantLines[:maxLines]
				}

				for _, line := range relevantLines {
					fileResult += "  " + line + "\n"
				}

				results = append(results, fileResult)
				if len(results) >= maxResults {
					return filepath.SkipDir
				}
			}
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("codebase search failed: %w", err)
	}

	if len(results) == 0 {
		return fmt.Sprintf("No relevant code found for '%s'. Try different keywords or check if files exist.", input), nil
	}

	response := fmt.Sprintf("🔍 Found %d relevant code sections for '%s':\n", len(results), input)
	response += strings.Join(results, "\n")

	return response, nil
}

// extractKeywords extracts searchable keywords from a query
func extractKeywords(query string) []string {
	// Convert to lowercase and split by common separators
	words := strings.FieldsFunc(strings.ToLower(query), func(r rune) bool {
		return r == ' ' || r == ',' || r == ';' || r == '.'
	})

	// Filter out common stop words and short words
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "or": true, "but": true,
		"in": true, "on": true, "at": true, "to": true, "for": true, "of": true,
		"with": true, "by": true, "is": true, "are": true, "was": true, "were": true,
		"be": true, "been": true, "have": true, "has": true, "had": true, "do": true,
		"does": true, "did": true, "will": true, "would": true, "could": true, "should": true,
	}

	var keywords []string
	for _, word := range words {
		word = strings.TrimSpace(word)
		if len(word) > 2 && !stopWords[word] {
			keywords = append(keywords, word)
		}
	}

	return keywords
}

// isCodeFile checks if a file extension indicates a code file
func isCodeFile(ext string) bool {
	codeExts := map[string]bool{
		".go": true, ".py": true, ".js": true, ".ts": true, ".java": true,
		".c": true, ".cpp": true, ".h": true, ".hpp": true, ".cs": true,
		".php": true, ".rb": true, ".rs": true, ".swift": true, ".kt": true,
		".scala": true, ".sh": true, ".bash": true, ".zsh": true, ".fish": true,
		".sql": true, ".html": true, ".css": true, ".scss": true, ".less": true,
		".json": true, ".yaml": true, ".yml": true, ".xml": true, ".toml": true,
		".md": true, ".txt": true, ".cfg": true, ".conf": true, ".ini": true,
	}
	return codeExts[ext]
}

// containsImportantPatterns checks for important code patterns
func containsImportantPatterns(line string, keywords []string) bool {
	// Look for function definitions, class definitions, error handling, etc.
	patterns := []string{
		"func ", "function ", "def ", "class ", "interface ", "struct ",
		"error", "exception", "catch", "try", "throw", "panic", "fatal",
		"security", "auth", "login", "password", "token", "encrypt", "decrypt",
		"database", "db", "sql", "query", "connection", "transaction",
		"import ", "require", "include", "#include", "from ", "export ",
	}

	for _, pattern := range patterns {
		if strings.Contains(line, pattern) {
			return true
		}
	}

	return false
}

// calculateRelevanceScore calculates how relevant a file content is to the keywords
func calculateRelevanceScore(content string, keywords []string) float64 {
	if len(keywords) == 0 {
		return 0
	}

	score := 0.0
	contentWords := strings.Fields(content)
	totalWords := float64(len(contentWords))

	for _, keyword := range keywords {
		count := float64(strings.Count(content, keyword))
		// TF-IDF-like scoring: term frequency normalized by content length
		if totalWords > 0 {
			tf := count / totalWords
			score += tf * 100 // Scale up for readability
		}
	}

	return score / float64(len(keywords)) // Average across keywords
}

// RagAutoDetectionTool provides RAG system detection and suggestion capability
type RagAutoDetectionTool struct {
	ragService interface {
		LoadRag(ragName string) (*domain.RagSystem, error)
		ListAllRags() ([]string, error)
	}
}

// NewRagAutoDetectionTool creates a new RagAutoDetectionTool
func NewRagAutoDetectionTool(ragService interface {
	LoadRag(ragName string) (*domain.RagSystem, error)
	ListAllRags() ([]string, error)
}) *RagAutoDetectionTool {
	return &RagAutoDetectionTool{
		ragService: ragService,
	}
}

func (t *RagAutoDetectionTool) Name() string {
	return "rag_search"
}

func (t *RagAutoDetectionTool) Description() string {
	return "Search and query RAG (Retrieval-Augmented Generation) knowledge bases. This tool can automatically detect available RAG systems and help you query them. Use this tool when you need to search through local knowledge bases or documents."
}

func (t *RagAutoDetectionTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]interface{}{
				"description": "The search query to find relevant information in the RAG knowledge bases",
				"type":        "string",
			},
			"rag_name": map[string]interface{}{
				"description": "Optional: specify which RAG system to use. If not provided, the tool will suggest available options",
				"type":        "string",
			},
		},
		"required": []string{"query"},
	}
}

func (t *RagAutoDetectionTool) ExecuteWithParams(ctx context.Context, params map[string]interface{}) (string, error) {
	query, ok := params["query"].(string)
	if !ok {
		return "", fmt.Errorf("query parameter is required and must be a string")
	}

	ragName, _ := params["rag_name"].(string)

	if ragName == "" {
		return t.Execute(ctx, query)
	} else {
		return t.executeWithSpecificRag(ctx, query, ragName)
	}
}

func (t *RagAutoDetectionTool) Execute(ctx context.Context, input string) (string, error) {
	// List available RAG systems
	ragNames, err := t.ragService.ListAllRags()
	if err != nil {
		return "", fmt.Errorf("failed to list available RAG systems: %w", err)
	}

	if len(ragNames) == 0 {
		return "No RAG systems found. Please create a RAG system first using 'rlama rag [model] [rag-name] [folder-path]' or specify a RAG system when running the agent with 'rlama agent run [rag-name] -q \"your question\"'.", nil
	}

	// If there's only one RAG system, suggest using it directly
	if len(ragNames) == 1 {
		ragName := ragNames[0]
		rag, err := t.ragService.LoadRag(ragName)
		if err != nil {
			return "", fmt.Errorf("failed to load RAG system '%s': %w", ragName, err)
		}

		var response strings.Builder
		response.WriteString(fmt.Sprintf("Found one RAG system available: **%s**\n", ragName))
		response.WriteString(fmt.Sprintf("- Model: %s\n", rag.ModelName))
		response.WriteString(fmt.Sprintf("- Documents: %d\n", len(rag.Documents)))
		response.WriteString(fmt.Sprintf("- Created: %s\n\n", rag.CreatedAt.Format("2006-01-02 15:04:05")))

		response.WriteString("To query this RAG system with your question, please run:\n")
		response.WriteString(fmt.Sprintf("```\nrlama agent run %s -q \"%s\"\n```\n\n", ragName, input))
		response.WriteString("This will give you access to the knowledge base and allow the agent to answer your question using the documents in the RAG system.")

		return response.String(), nil
	}

	// If there are multiple RAG systems, show all options
	var ragInfo strings.Builder
	ragInfo.WriteString(fmt.Sprintf("Found %d available RAG system(s):\n\n", len(ragNames)))

	for i, ragName := range ragNames {
		rag, err := t.ragService.LoadRag(ragName)
		if err != nil {
			ragInfo.WriteString(fmt.Sprintf("%d. %s (failed to load details)\n", i+1, ragName))
			continue
		}

		ragInfo.WriteString(fmt.Sprintf("%d. **%s**\n", i+1, ragName))
		ragInfo.WriteString(fmt.Sprintf("   - Model: %s\n", rag.ModelName))
		ragInfo.WriteString(fmt.Sprintf("   - Documents: %d\n", len(rag.Documents)))
		ragInfo.WriteString(fmt.Sprintf("   - Created: %s\n", rag.CreatedAt.Format("2006-01-02 15:04:05")))
		ragInfo.WriteString("\n")
	}

	ragInfo.WriteString("To query a specific RAG system, please run the agent with a RAG system specified:\n")
	ragInfo.WriteString("Example: rlama agent run [rag-name] -q \"your question\"\n\n")

	ragInfo.WriteString("Available RAG systems to choose from:\n")
	for _, ragName := range ragNames {
		ragInfo.WriteString(fmt.Sprintf("- rlama agent run %s -q \"%s\"\n", ragName, input))
	}

	return ragInfo.String(), nil
}

func (t *RagAutoDetectionTool) executeWithSpecificRag(ctx context.Context, query, ragName string) (string, error) {
	// Load the specific RAG system
	rag, err := t.ragService.LoadRag(ragName)
	if err != nil {
		return "", fmt.Errorf("failed to load RAG system '%s': %w", ragName, err)
	}

	// Since we don't have access to the full RAG service here, we need to suggest using the correct command
	return fmt.Sprintf("RAG system '%s' is available with %d documents. However, this tool cannot directly query RAG systems. Please run:\n\nrlama agent run %s -q \"%s\"\n\nThis will give you access to the RAG search capabilities you're looking for.", ragName, len(rag.Documents), ragName, query), nil
}

// FlightSearchTool provides flight search capability using web search
type FlightSearchTool struct {
	webSearchTool *WebSearchTool
}

// NewFlightSearchTool creates a new FlightSearchTool
func NewFlightSearchTool(apiKey string) *FlightSearchTool {
	return &FlightSearchTool{
		webSearchTool: NewWebSearchTool(apiKey),
	}
}

func (t *FlightSearchTool) Name() string {
	return "flight_search"
}

func (t *FlightSearchTool) Description() string {
	return "Search for flight information between two cities. Input should be in format: 'from_city to destination_city' or 'from_city to destination_city on date'. This tool will search for flight options, prices, and availability."
}

func (t *FlightSearchTool) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"explanation": map[string]interface{}{
				"description": "One sentence explanation as to why this tool is being used, and how it contributes to the goal.",
				"type":        "string",
			},
			"from_city": map[string]interface{}{
				"description": "The departure city or airport code",
				"type":        "string",
			},
			"to_city": map[string]interface{}{
				"description": "The destination city or airport code",
				"type":        "string",
			},
			"date": map[string]interface{}{
				"description": "Optional: travel date in YYYY-MM-DD format or relative format like 'next week'",
				"type":        "string",
			},
		},
		"required": []string{"from_city", "to_city"},
	}
}

func (t *FlightSearchTool) ExecuteWithParams(ctx context.Context, params map[string]interface{}) (string, error) {
	fromCity, ok := params["from_city"].(string)
	if !ok {
		return "", fmt.Errorf("from_city parameter is required and must be a string")
	}

	toCity, ok := params["to_city"].(string)
	if !ok {
		return "", fmt.Errorf("to_city parameter is required and must be a string")
	}

	date, _ := params["date"].(string)

	var input string
	if date != "" {
		input = fmt.Sprintf("%s to %s on %s", fromCity, toCity, date)
	} else {
		input = fmt.Sprintf("%s to %s", fromCity, toCity)
	}

	return t.Execute(ctx, input)
}

func (t *FlightSearchTool) Execute(ctx context.Context, input string) (string, error) {
	// Parse the input to extract cities and optional date
	parts := strings.Split(strings.ToLower(input), " to ")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid input format. Expected 'from_city to destination_city' or 'from_city to destination_city on date'")
	}

	fromCity := strings.TrimSpace(parts[0])

	// Check if there's a date specified
	toParts := strings.Split(parts[1], " on ")
	toCity := strings.TrimSpace(toParts[0])
	var date string
	if len(toParts) > 1 {
		date = strings.TrimSpace(toParts[1])
	}

	// Build search query for flights
	var searchQuery string
	if date != "" {
		searchQuery = fmt.Sprintf("flights from %s to %s %s price booking", fromCity, toCity, date)
	} else {
		searchQuery = fmt.Sprintf("flights from %s to %s price booking cheap tickets", fromCity, toCity)
	}

	// Use web search to find flight information
	result, err := t.webSearchTool.Execute(ctx, searchQuery)
	if err != nil {
		return "", fmt.Errorf("failed to search for flights: %w", err)
	}

	// Format the response specifically for flights
	response := fmt.Sprintf("🛫 Flight search results for %s to %s", fromCity, toCity)
	if date != "" {
		response += fmt.Sprintf(" on %s", date)
	}
	response += ":\n\n"
	response += result
	response += "\n\n💡 Tip: For the best prices and to book flights, visit the airline websites or travel booking sites mentioned above."

	return response, nil
}
