# RLAMA Agents - Complete Guide

## Table of Contents
- [Introduction](#introduction)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [Available Tools](#available-tools)
- [Usage Examples](#usage-examples)
- [Configuration](#configuration)
- [Advanced Features](#advanced-features)
- [API Reference](#api-reference)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

## Introduction

RLAMA Agents are intelligent AI assistants powered by Large Language Models (LLMs) that can interact with your local environment using a comprehensive set of tools. Similar to Cursor AI, RLAMA agents can read files, search codebases, browse directories, and even search the web - all while maintaining context and chaining operations intelligently.

**🆕 New in this release**: RLAMA Agents now feature **automatic RAG detection and loading**. When you have a single knowledge base, it's loaded automatically. With multiple RAG systems, the agent intelligently suggests which one to use, making knowledge base interaction seamless and intuitive.

### Key Features

- 🤖 **Intelligent Tool Selection**: Automatically chooses the right tool for each task
- 🔗 **Tool Chaining**: Seamlessly chains multiple operations together
- 📁 **File System Integration**: Browse, read, and search your local files
- 🌐 **Web Search Capability**: Access real-time information from the internet
- 🧠 **RAG Integration**: Query your local knowledge bases with auto-detection
- 🔍 **Smart RAG Discovery**: Automatically detects and loads available RAG systems
- 🔄 **Structured Outputs**: JSON schema validation for reliable responses
- 💬 **Conversational Interface**: Natural language interaction with context awareness

## Architecture

### Agent Types

```
┌─────────────────────────────────────────────────────────────┐
│                         RLAMA Agent                        │
├─────────────────────────────────────────────────────────────┤
│                    Conversational Mode                     │
│  ┌─────────────────┐  ┌─────────────────┐  ┌──────────────┐ │
│  │   LLM Client    │  │     Memory      │  │    Tools     │ │
│  │                 │  │                 │  │              │ │
│  │ • qwen3:8b      │  │ • Conversation  │  │ • list_dir   │ │
│  │ • llama3        │  │   History       │  │ • read_file  │ │
│  │ • openai models │  │ • Context       │  │ • file_search│ │
│  │ • custom models │  │   Management    │  │ • grep_search│ │
│  └─────────────────┘  └─────────────────┘  │ • codebase   │ │
│                                           │ • web_search │ │
│                                           │ • rag_search │ │
│                                           └──────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

### Tool Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                         Tool Interface                       │
├──────────────────────────────────────────────────────────────┤
│  Name() string                                              │
│  Description() string                                       │
│  Execute(ctx, input) (string, error)                       │
│  Schema() map[string]interface{}                           │
│  ExecuteWithParams(ctx, params) (string, error)           │
└──────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────┬─────────────────┬─────────────────┬─────────┐
│   File Tools    │  Search Tools   │  Web Tools      │ RAG     │
├─────────────────┼─────────────────┼─────────────────┼─────────┤
│ • list_dir      │ • file_search   │ • web_search    │ • rag   │
│ • read_file     │ • grep_search   │                 │ _search │
│ • file_write    │ • codebase      │                 │         │
│                 │   _search       │                 │         │
└─────────────────┴─────────────────┴─────────────────┴─────────┘
```

## Quick Start

### Installation

Ensure you have RLAMA installed and configured:

```bash
# Install RLAMA
curl -fsSL https://raw.githubusercontent.com/dontizi/rlama/main/install.sh | bash

# Install dependencies for optimal performance
./scripts/install_deps.sh

# Configure Ollama (if using local models)
export OLLAMA_HOST=http://localhost:11434
```

### Basic Usage

#### Simple File Operations

```bash
# List files in current directory
rlama agent run -q "Show me all files in the current directory"

# Read a specific file
rlama agent run -q "Read the README.md file"

# Search for files by name
rlama agent run -q "Find all Python files in this project"
```

#### Code Analysis

```bash
# Analyze code structure
rlama agent run -q "Find all Go files with error handling"

# Search for specific patterns
rlama agent run -q "Find functions that contain 'database' or 'sql'"

# Get project overview
rlama agent run -q "Analyze the project structure and tell me what this codebase does"
```

#### With RAG Integration

```bash
# First create a RAG system
rlama rag qwen3:8b my-docs ./documents

# Run agent with explicit RAG specification
rlama agent run my-docs -q "What does the documentation say about installation?"

# Or let the agent auto-detect available RAG systems
rlama agent run -q "Search my documentation for installation instructions"

# With only one RAG system available, it will be loaded automatically
rlama agent run -q "What are the main features of this project?"
```

## Available Tools

### 1. File System Tools

#### `list_dir` - Directory Listing
**Purpose**: Explore directory contents and understand file structure

```bash
# Examples
rlama agent run -q "List all files in the src directory"
rlama agent run -q "Show me the project structure"
```

**Schema**:
```json
{
  "type": "object",
  "properties": {
    "explanation": {
      "type": "string",
      "description": "Why this tool is being used"
    },
    "relative_workspace_path": {
      "type": "string", 
      "description": "Path relative to workspace root"
    }
  },
  "required": ["relative_workspace_path"]
}
```

#### `read_file` - File Reading
**Purpose**: Read file contents with optional line range specification

```bash
# Examples
rlama agent run -q "Read the main.go file"
rlama agent run -q "Show me lines 1-50 of the config file"
rlama agent run -q "Read the error handling code in utils.py"
```

**Schema**:
```json
{
  "type": "object",
  "properties": {
    "target_file": {"type": "string"},
    "should_read_entire_file": {"type": "boolean"},
    "start_line_one_indexed": {"type": "integer"},
    "end_line_one_indexed_inclusive": {"type": "integer"}
  },
  "required": ["target_file", "should_read_entire_file", "start_line_one_indexed", "end_line_one_indexed_inclusive"]
}
```

### 2. Search Tools

#### `file_search` - Fuzzy File Name Search
**Purpose**: Find files by name using fuzzy matching

```bash
# Examples
rlama agent run -q "Find config files"
rlama agent run -q "Locate test files"
rlama agent run -q "Find all README files"
```

#### `grep_search` - Text Pattern Search
**Purpose**: Find exact text patterns or regex matches in files

```bash
# Examples
rlama agent run -q "Find all TODO comments in the codebase"
rlama agent run -q "Search for functions named 'validate'"
rlama agent run -q "Find error messages containing 'connection'"
```

**Advanced grep examples**:
```bash
# Case sensitive search in Go files
rlama agent run -q "Search for 'func main' in Go files, case sensitive"

# Search with file patterns
rlama agent run -q "Find 'import' statements in Python files only"
```

#### `codebase_search` - Semantic Code Search
**Purpose**: Find code semantically related to your query

```bash
# Examples
rlama agent run -q "Find authentication and security related code"
rlama agent run -q "Show me database connection handling"
rlama agent run -q "Find error handling patterns"
```

### 3. Web Tools

#### `web_search` - Real-time Web Search
**Purpose**: Access current information from the internet

```bash
# Setup (required for web search)
export GOOGLE_SEARCH_API_KEY="your_api_key"
export GOOGLE_SEARCH_ENGINE_ID="your_engine_id"

# Examples
rlama agent run -w -q "What's the current weather in Paris?"
rlama agent run -w -q "Latest Go 1.21 features"
rlama agent run -w -q "Best practices for Docker security 2024"
```

### 4. RAG Tools

#### `rag_search` - Knowledge Base Query & Auto-Detection
**Purpose**: Search your local document collections with intelligent RAG discovery

```bash
# Explicit RAG specification (traditional method)
rlama agent run my-docs -q "How do I configure the database?"
rlama agent run my-docs -q "What are the API endpoints?"

# Auto-detection mode (new feature)
rlama agent run -q "Search my knowledge base for database configuration"
rlama agent run -q "What RAG systems are available and what do they contain?"

# Automatic single RAG loading
# If only one RAG exists, it's loaded automatically:
rlama agent run -q "Query my documentation for API usage examples"
```

**Auto-Detection Features**:
- **Single RAG Auto-Load**: When only one RAG system exists, it's automatically loaded
- **Multi-RAG Discovery**: Lists all available RAG systems with metadata
- **Smart Suggestions**: Provides exact commands to query specific RAG systems
- **Seamless Integration**: Works transparently with other agent tools

## Usage Examples

### Example 1: Code Analysis Workflow

```bash
# Comprehensive code analysis
rlama agent run -q "Analyze this Go project: first show me the structure, then find the main entry points, and finally look for any error handling patterns"
```

**Expected workflow**:
1. Agent uses `list_dir` to explore project structure
2. Agent uses `file_search` to find main files
3. Agent uses `read_file` to examine entry points  
4. Agent uses `codebase_search` for error handling patterns
5. Agent provides comprehensive analysis

### Example 2: Documentation Research

```bash
# Research with web and local docs (explicit RAG)
rlama agent run my-docs -w -q "Compare the local API documentation with the latest online best practices for REST API design"

# Research with auto-detected RAG
rlama agent run -w -q "Compare my local documentation with online best practices for API design"
```

**Expected workflow**:
1. Agent auto-detects available RAG systems or uses specified one
2. Agent uses `rag_search` to query local documentation
3. Agent uses `web_search` to find current best practices
4. Agent provides comparative analysis

### Example 3: Debugging Session

```bash
# Debug investigation
rlama agent run -q "I'm getting a 'connection refused' error. Help me find where database connections are configured and check for any related error handling"
```

**Expected workflow**:
1. Agent uses `grep_search` to find "connection refused" messages
2. Agent uses `codebase_search` for database configuration
3. Agent uses `read_file` to examine specific config files
4. Agent provides debugging recommendations

### Example 4: Project Onboarding

```bash
# New team member onboarding
rlama agent run -q "I'm new to this project. Can you give me an overview of the codebase structure, main components, and point me to getting started documentation?"
```

## Configuration

### Agent Modes

```bash
# Conversational mode (default)
rlama agent run -q "your question"

# Autonomous mode (coming soon)
rlama agent run -m autonomous -q "your goal"
```

### Model Selection

```bash
# Using specific model
rlama agent run -l qwen3:8b -q "your question"
rlama agent run -l llama3 -q "your question"

# With OpenAI models (requires profile setup)
rlama profile add openai-profile --provider openai --api-key your_key
rlama agent run -l gpt-4 -q "your question" --profile openai-profile
```

### RAG Configuration

#### RAG Usage Modes

```bash
# Auto-detection mode (recommended for single RAG setups)
rlama agent run -q "your question"

# Explicit RAG specification (for multi-RAG environments)
rlama agent run my-docs -q "your question"

# List available RAG systems
rlama list
```

#### RAG Setup Workflow

1. **Create a RAG System**:
```bash
# Basic RAG creation
rlama rag qwen3:8b my-docs ./documents

# Advanced RAG with options
rlama rag llama3 research-rag ./papers --chunk-size=1500 --chunk-overlap=300
```

2. **Verify RAG Creation**:
```bash
# List all RAG systems
rlama list

# Check RAG contents
rlama list-docs my-docs
```

3. **Use with Agent**:
```bash
# Single RAG (auto-loads)
rlama agent run -q "What's in my documentation?"

# Multiple RAGs (shows options)
rlama agent run -q "Which knowledge bases are available?"

# Specific RAG
rlama agent run my-docs -q "Find API documentation"
```

### Web Search Setup

1. **Get Google Custom Search API Key**:
   - Visit [Google Cloud Console](https://console.cloud.google.com/)
   - Enable Custom Search API
   - Create credentials (API Key)

2. **Create Custom Search Engine**:
   - Go to [Google Custom Search](https://cse.google.com/)
   - Create a new search engine
   - Note the Search Engine ID

3. **Configure Environment**:
```bash
export GOOGLE_SEARCH_API_KEY="your_api_key_here"
export GOOGLE_SEARCH_ENGINE_ID="your_engine_id_here"
```

4. **Use Web Search**:
```bash
rlama agent run -w -q "your web query"
# OR
rlama agent run --web --search-api-key "key" --search-engine-id "id" -q "query"
```

## Advanced Features

### RAG Auto-Detection & Smart Loading

RLAMA Agents now feature intelligent RAG system discovery and automatic loading:

#### Automatic Single RAG Detection
```bash
# When you have only one RAG system, it's loaded automatically
rlama agent run -q "What's in my documentation?"
# Output: 📚 Auto-detected and loaded RAG system: my-docs
```

#### Multi-RAG Discovery
```bash
# When multiple RAG systems exist, agent shows all options
rlama agent run -q "What knowledge bases are available?"
```

**Example output**:
```
Found 2 available RAG system(s):

1. **project-docs**
   - Model: qwen3:8b
   - Documents: 15
   - Created: 2024-01-15 10:30:00

2. **research-papers**
   - Model: llama3.1
   - Documents: 8
   - Created: 2024-01-20 14:20:00

To query a specific RAG system, please run:
- rlama agent run project-docs -q "your question"
- rlama agent run research-papers -q "your question"
```

#### Smart RAG Suggestions
The agent provides context-aware suggestions:

```bash
rlama agent run -q "Find information about API design patterns"
```

**Agent response includes**:
- List of available RAG systems
- Suggested commands for each RAG
- Metadata about document collections
- Automatic command generation with your exact query

#### Fallback Behavior
- **No RAG systems**: Agent suggests creating one with `rlama rag` command
- **Single RAG**: Automatically loaded without user intervention
- **Multiple RAGs**: Interactive selection with detailed information
- **RAG specified**: Traditional explicit loading behavior

### Tool Chaining

The agent automatically chains tools based on context:

```bash
# Single request, multiple operations
rlama agent run -q "Find all configuration files, read the main config, and tell me what database settings are configured"
```

**Behind the scenes**:
1. `file_search` → finds config files
2. `read_file` → reads main configuration
3. `grep_search` → finds database settings
4. Synthesizes final answer

### Structured Outputs

All tools support JSON schema validation:

```go
// Example tool usage in Go code
type FileReadParams struct {
    TargetFile           string `json:"target_file"`
    ShouldReadEntireFile bool   `json:"should_read_entire_file"`
    StartLine           int    `json:"start_line_one_indexed"`
    EndLine             int    `json:"end_line_one_indexed_inclusive"`
}
```

### Context Management

The agent maintains conversation history:

```bash
# First query
rlama agent run -q "Show me the project structure"

# Follow-up query (remembers previous context)  
rlama agent run -q "Now read the main.go file"

# Another follow-up
rlama agent run -q "Are there any tests for the main package?"
```

### Error Handling

Intelligent error recovery:

```bash
# If a file doesn't exist
rlama agent run -q "Read nonexistent.txt"
# Agent will: list directory → suggest similar files → ask for clarification
```

## API Reference

### Agent Interface

```go
type Agent interface {
    Run(ctx context.Context, input string) (string, error)
    AddTool(tool Tool) error
    GetTools() []Tool
    GetMode() AgentMode
    GetMemory() Memory
}
```

### Tool Interface

```go
type Tool interface {
    Name() string
    Description() string
    Execute(ctx context.Context, input string) (string, error)
    Schema() map[string]interface{}
    ExecuteWithParams(ctx context.Context, params map[string]interface{}) (string, error)
}
```

### Memory Interface

```go
type Memory interface {
    Store(key string, value interface{}) error
    Retrieve(key string) (interface{}, error)
    GetHistory() []string
    AddToHistory(entry string) error
}
```

## Best Practices

### 1. Query Design

**Good queries**:
```bash
✅ "Find all error handling code in the user authentication module"
✅ "Show me the API configuration and explain how to modify the database settings"
✅ "List the test files and run the main test suite"
```

**Avoid**:
```bash
❌ "Fix my code" (too vague)
❌ "Do everything" (no clear goal)
❌ "Read all files" (inefficient)
```

### 2. Tool Selection

The agent follows this hierarchy:
1. **list_dir** → for exploration
2. **file_search** → for finding specific files
3. **read_file** → for reading known files
4. **grep_search** → for exact pattern matching
5. **codebase_search** → for semantic code search
6. **web_search** → for real-time information
7. **rag_search** → for local knowledge base

### 3. Workflow Optimization

**Efficient workflows**:
- Always list directories before reading files
- Use specific file patterns for grep search
- Combine local and web search for comprehensive research
- Leverage RAG for project-specific documentation
- Let the agent auto-detect RAG systems for simplified usage
- Use explicit RAG specification only when working with multiple knowledge bases

### 4. Security Considerations

- Files are accessed relative to the workspace root
- No access outside the base directory
- Web search requires explicit enabling
- API keys should be environment variables

## Troubleshooting

### Common Issues

#### 1. Tool Not Found
```
Error: chained tool file_read not found
```
**Solution**: Tool names have been updated. Use `read_file` instead of `file_read`.

#### 2. Web Search Fails
```
Error: Google Search API key not provided
```
**Solution**: 
```bash
export GOOGLE_SEARCH_API_KEY="your_key"
export GOOGLE_SEARCH_ENGINE_ID="your_id"
rlama agent run -w -q "your query"
```

#### 3. File Access Denied
```
Error: access to path outside base directory is not allowed
```
**Solution**: Ensure file paths are relative to workspace root.

#### 4. Model Connection Issues
```
Error: failed to generate LLM completion
```
**Solution**: 
- Check Ollama is running: `ollama serve`
- Verify model is available: `ollama list`
- Check host configuration: `export OLLAMA_HOST=http://localhost:11434`

#### 5. RAG Auto-Detection Issues
```
Error: No RAG systems found
```
**Solution**: 
```bash
# Check available RAG systems
rlama list

# Create a RAG system if none exist
rlama rag qwen3:8b my-docs ./documents

# Verify RAG creation
rlama list
```

#### 6. RAG Loading Issues
```
Error: failed to auto-load RAG system
```
**Solution**: 
- Ensure the RAG system exists: `rlama list`
- Check RAG permissions and file integrity
- Try loading explicitly: `rlama agent run [rag-name] -q "your question"`
- Recreate RAG if corrupted: `rlama delete [rag-name] --force` then recreate

### Debug Mode

Enable debug output for troubleshooting:

```go
// In your Go code
agent.Debug = true
```

Or set verbose mode:
```bash
rlama agent run -v -q "your query"
```

### Performance Tips

1. **Use specific queries** → faster tool selection
2. **Leverage file patterns** → reduce search scope
3. **Chain operations efficiently** → minimize LLM calls
4. **Use appropriate search tools** → grep for exact, codebase for semantic
5. **Configure local models** → faster response times

## Contributing

To add new tools:

1. Implement the `Tool` interface
2. Add JSON schema definition  
3. Register in `AgentService.setupTools()`
4. Update documentation
5. Add tests

Example tool implementation:

```go
type CustomTool struct {
    baseDir string
}

func (t *CustomTool) Name() string {
    return "custom_tool"
}

func (t *CustomTool) Description() string {
    return "Custom tool description"
}

func (t *CustomTool) Schema() map[string]interface{} {
    return map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "param": map[string]interface{}{
                "type": "string",
                "description": "Parameter description",
            },
        },
        "required": []string{"param"},
    }
}

func (t *CustomTool) Execute(ctx context.Context, input string) (string, error) {
    // Implementation
    return "result", nil
}

func (t *CustomTool) ExecuteWithParams(ctx context.Context, params map[string]interface{}) (string, error) {
    // Implementation with structured params
    return "result", nil
}
```

---

## License

RLAMA is licensed under the MIT License. See [LICENSE](LICENSE) for details.

## Support

- 📖 Documentation: [GitHub Wiki](https://github.com/dontizi/rlama/wiki)
- 🐛 Issues: [GitHub Issues](https://github.com/dontizi/rlama/issues)  
- 💬 Discussions: [GitHub Discussions](https://github.com/dontizi/rlama/discussions)
- 🌐 Website: [rlama.dev](https://rlama.dev)

---

**Happy coding with RLAMA Agents! 🚀** 