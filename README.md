<!-- Social Links Navigation Bar -->
<div align="center">
  <a href="https://x.com/LeDonTizi" target="_blank">
    <img src="https://img.shields.io/badge/Twitter-1DA1F2?style=for-the-badge&logo=twitter&logoColor=white" alt="Twitter">
  </a>
  <a href="https://discord.gg/tP5JB9DR" target="_blank">
    <img src="https://img.shields.io/badge/Discord-5865F2?style=for-the-badge&logo=discord&logoColor=white" alt="Discord">
  </a>
  <a href="https://www.youtube.com/@Dontizi" target="_blank">
    <img src="https://img.shields.io/badge/YouTube-FF0000?style=for-the-badge&logo=youtube&logoColor=white" alt="YouTube">
  </a>
</div>

<br>

# RLAMA - User Guide
RLAMA is a powerful AI-driven question-answering tool for your documents, seamlessly integrating with your local Ollama models. It enables you to create, manage, and interact with Retrieval-Augmented Generation (RAG) systems tailored to your documentation needs.


[![RLAMA Demonstration](https://img.youtube.com/vi/EIsQnBqeQxQ/0.jpg)](https://www.youtube.com/watch?v=EIsQnBqeQxQ)

## Table of Contents
- [Vision & Roadmap](#vision--roadmap)
- [Installation](#installation)
- [Available Commands](#available-commands)
  - [rag - Create a RAG system](#rag---create-a-rag-system)
  - [run - Use a RAG system](#run---use-a-rag-system)
  - [list - List RAG systems](#list---list-rag-systems)
  - [delete - Delete a RAG system](#delete---delete-a-rag-system)
  - [list-docs - List documents in a RAG](#list-docs---list-documents-in-a-rag)
  - [list-chunks - Inspect document chunks](#list-chunks---inspect-document-chunks)
  - [view-chunk - View chunk details](#view-chunk---view-chunk-details)
  - [add-docs - Add documents to RAG](#add-docs---add-documents-to-rag)
  - [update-model - Change LLM model](#update-model---change-llm-model)
  - [update - Update RLAMA](#update---update-rlama)
  - [version - Display version](#version---display-version)
- [Uninstallation](#uninstallation)
- [Supported Document Formats](#supported-document-formats)
- [Troubleshooting](#troubleshooting)

## Vision & Roadmap

RLAMA aims to become the definitive tool for creating local RAG systems that work seamlessly for everyone—from individual developers to large enterprises. Here's our strategic roadmap:

### Completed Features ✅
- ✅ **Basic RAG System Creation**: CLI tool for creating and managing RAG systems
- ✅ **Document Processing**: Support for multiple document formats (.txt, .md, .pdf, etc.)
- ✅ **Document Chunking**: Basic text splitting with configurable size and overlap
- ✅ **Vector Storage**: Local storage of document embeddings
- ✅ **Context Retrieval**: Basic semantic search with configurable context size
- ✅ **Ollama Integration**: Seamless connection to Ollama models
- ✅ **Cross-Platform Support**: Works on Linux, macOS, and Windows
- ✅ **Easy Installation**: One-line installation script

### Small LLM Optimization (Q2 2025)
- [ ] **Prompt Compression**: Smart context summarization for limited context windows
- [ ] **Adaptive Chunking**: Dynamic content segmentation based on semantic boundaries
- [ ] **Minimal Context Retrieval**: Intelligent filtering to eliminate redundant content
- [ ] **Parameter Optimization**: Fine-tuned settings for different model sizes

### Advanced Embedding Pipeline (Q2-Q3 2025)
- [ ] **Multi-Model Embedding Support**: Integration with various embedding models
- [ ] **Hybrid Retrieval Techniques**: Combining sparse and dense retrievers for better accuracy
- [ ] **Embedding Evaluation Tools**: Built-in metrics to measure retrieval quality
- [ ] **Automated Embedding Cache**: Smart caching to reduce computation for similar queries

### User Experience Enhancements (Q3 2025)
- [ ] **Lightweight Web Interface**: Simple browser-based UI for the existing CLI backend
- [ ] **Knowledge Graph Visualization**: Interactive exploration of document connections
- [ ] **Guided RAG Setup Wizard**: Step-by-step interface for non-technical users
- [ ] **Domain-Specific Templates**: Pre-configured settings for different domains

### Enterprise Features (Q4 2025)
- [ ] **Multi-User Access Control**: Role-based permissions for team environments
- [ ] **Integration with Enterprise Systems**: Connectors for SharePoint, Confluence, Google Workspace
- [ ] **Knowledge Quality Monitoring**: Detection of outdated or contradictory information
- [ ] **System Integration API**: Webhooks and APIs for embedding RLAMA in existing workflows
- [ ] **AI Agent Creation Framework**: Simplified system for building custom AI agents with RAG capabilities

### Next-Gen Retrieval Innovations (Q1 2026)
- [ ] **Multi-Step Retrieval**: Using the LLM to refine search queries for complex questions
- [ ] **Cross-Modal Retrieval**: Support for image content understanding and retrieval
- [ ] **Feedback-Based Optimization**: Learning from user interactions to improve retrieval
- [ ] **Knowledge Graphs & Symbolic Reasoning**: Combining vector search with structured knowledge

RLAMA's core philosophy remains unchanged: to provide a simple, powerful, local RAG solution that respects privacy, minimizes resource requirements, and works seamlessly across platforms.

## Installation

### Prerequisites
- [Ollama](https://ollama.ai/) installed and running

### Installation from terminal

```bash
curl -fsSL https://raw.githubusercontent.com/dontizi/rlama/main/install.sh | sh
```

## Tech Stack

RLAMA is built with:

- **Core Language**: Go (chosen for performance, cross-platform compatibility, and single binary distribution)
- **CLI Framework**: Cobra (for command-line interface structure)
- **LLM Integration**: Ollama API (for embeddings and completions)
- **Storage**: Local filesystem-based storage (JSON files for simplicity and portability)
- **Vector Search**: Custom implementation of cosine similarity for embedding retrieval

## Architecture

RLAMA follows a clean architecture pattern with clear separation of concerns:

```
rlama/
├── cmd/                  # CLI commands (using Cobra)
│   ├── root.go           # Base command
│   ├── rag.go            # Create RAG systems
│   ├── run.go            # Query RAG systems
│   └── ...
├── internal/
│   ├── client/           # External API clients
│   │   └── ollama_client.go # Ollama API integration
│   ├── domain/           # Core domain models
│   │   ├── rag.go        # RAG system entity
│   │   └── document.go   # Document entity
│   ├── repository/       # Data persistence
│   │   └── rag_repository.go # Handles saving/loading RAGs
│   └── service/          # Business logic
│       ├── rag_service.go      # RAG operations
│       ├── document_loader.go  # Document processing
│       └── embedding_service.go # Vector embeddings
└── pkg/                  # Shared utilities
    └── vector/           # Vector operations
```

## Data Flow

1. **Document Processing**: Documents are loaded from the file system, parsed based on their type, and converted to plain text.
2. **Embedding Generation**: Document text is sent to Ollama to generate vector embeddings.
3. **Storage**: The RAG system (documents + embeddings) is stored in the user's home directory (~/.rlama).
4. **Query Process**: When a user asks a question, it's converted to an embedding, compared against stored document embeddings, and relevant content is retrieved.
5. **Response Generation**: Retrieved content and the question are sent to Ollama to generate a contextually-informed response.

## Visual Representation

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Documents  │────>│  Document   │────>│  Embedding  │
│  (Input)    │     │  Processing │     │  Generation │
└─────────────┘     └─────────────┘     └─────────────┘
                                              │
                                              ▼
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Query     │────>│  Vector     │<────│ Vector Store│
│  Response   │     │  Search     │     │ (RAG System)│
└─────────────┘     └─────────────┘     └─────────────┘
       ▲                   │
       │                   ▼
┌─────────────┐     ┌─────────────┐
│   Ollama    │<────│   Context   │
│    LLM      │     │  Building   │
└─────────────┘     └─────────────┘
```

RLAMA is designed to be lightweight and portable, focusing on providing RAG capabilities with minimal dependencies. The entire system runs locally, with the only external dependency being Ollama for LLM capabilities.

## Available Commands

You can get help on all commands by using:

```bash
rlama --help
```

### Global Flags

These flags can be used with any command:

```bash
--host string   Ollama host (default: localhost)
--port string   Ollama port (default: 11434)
```

### rag - Create a RAG system

Creates a new RAG system by indexing all documents in the specified folder.

```bash
rlama rag [model] [rag-name] [folder-path]
```

**Parameters:**
- `model`: Name of the Ollama model to use (e.g., llama3, mistral, gemma).
- `rag-name`: Unique name to identify your RAG system.
- `folder-path`: Path to the folder containing your documents.

**Example:**

```bash
rlama rag llama3 documentation ./docs
```

### run - Use a RAG system

Starts an interactive session to interact with an existing RAG system.

```bash
rlama run [rag-name]
```

**Parameters:**
- `rag-name`: Name of the RAG system to use.
- `--context-size`: (Optional) Number of context chunks to retrieve (default: 20)

**Example:**

```bash
rlama run documentation
> How do I install the project?
> What are the main features?
> exit
```

**Context Size Tips:**
- Smaller values (5-15) for faster responses with key information
- Medium values (20-40) for balanced performance
- Larger values (50+) for complex questions needing broad context
- Consider your model's context window limits

```bash
rlama run documentation --context-size=50  # Use 50 context chunks
```

### list - List RAG systems

Displays a list of all available RAG systems.

```bash
rlama list
```

### delete - Delete a RAG system

Permanently deletes a RAG system and all its indexed documents.

```bash
rlama delete [rag-name] [--force/-f]
```

**Parameters:**
- `rag-name`: Name of the RAG system to delete.
- `--force` or `-f`: (Optional) Delete without asking for confirmation.

**Example:**

```bash
rlama delete old-project
```

Or to delete without confirmation:

```bash
rlama delete old-project --force
```

### list-docs - List documents in a RAG

Displays all documents in a RAG system with metadata.

```bash
rlama list-docs [rag-name]
```

**Parameters:**
- `rag-name`: Name of the RAG system

**Example:**

```bash
rlama list-docs documentation
```

### list-chunks - Inspect document chunks

List and filter document chunks in a RAG system with various options:

```bash
# Basic chunk listing
rlama list-chunks [rag-name]

# With content preview (shows first 100 characters)
rlama list-chunks [rag-name] --show-content

# Filter by document name/ID substring
rlama list-chunks [rag-name] --document=readme

# Combine options
rlama list-chunks [rag-name] --document=api --show-content
```

**Options:**
- `--show-content`: Display chunk content preview
- `--document`: Filter by document name/ID substring

**Output columns:**
- Chunk ID (use with view-chunk command)
- Document Source
- Chunk Position (e.g., "2/5" for second of five chunks)
- Content Preview (if enabled)
- Created Date

### view-chunk - View chunk details

Display detailed information about a specific chunk.

```bash
rlama view-chunk [rag-name] [chunk-id]
```

**Parameters:**
- `rag-name`: Name of the RAG system
- `chunk-id`: Chunk identifier from list-chunks

**Example:**

```bash
rlama view-chunk documentation doc123_chunk_0
```

### add-docs - Add documents to RAG

Add new documents to an existing RAG system.

```bash
rlama add-docs [rag-name] [folder-path] [flags]
```

**Parameters:**
- `rag-name`: Name of the RAG system
- `folder-path`: Path to documents folder

**Example:**

```bash
rlama add-docs documentation ./new-docs --exclude-ext=.tmp
```

### update-model - Change LLM model

Update the LLM model used by a RAG system.

```bash
rlama update-model [rag-name] [new-model]
```

**Parameters:**
- `rag-name`: Name of the RAG system
- `new-model`: New Ollama model name

**Example:**

```bash
rlama update-model documentation deepseek-r1:7b-instruct
```

### update - Update RLAMA

Checks if a new version of RLAMA is available and installs it.

```bash
rlama update [--force/-f]
```

**Options:**
- `--force` or `-f`: (Optional) Update without asking for confirmation.

### version - Display version

Displays the current version of RLAMA.

```bash
rlama --version
```

or

```bash
rlama -v
```

## Uninstallation

To uninstall RLAMA:

### Removing the binary

If you installed via `go install`:

```bash
rlama uninstall
```

### Removing data

RLAMA stores its data in `~/.rlama`. To remove it:

```bash
rm -rf ~/.rlama
```

## Supported Document Formats

RLAMA supports many file formats:

- **Text**: `.txt`, `.md`, `.html`, `.json`, `.csv`, `.yaml`, `.yml`, `.xml`, `.org`
- **Code**: `.go`, `.py`, `.js`, `.java`, `.c`, `.cpp`, `.cxx`, `.h`, `.rb`, `.php`, `.rs`, `.swift`, `.kt`, `.ts`, `.f`, `.F`, `.F90`, `.el`, `.svelte`
- **Documents**: `.pdf`, `.docx`, `.doc`, `.rtf`, `.odt`, `.pptx`, `.ppt`, `.xlsx`, `.xls`, `.epub`

Installing dependencies via `install_deps.sh` is recommended to improve support for certain formats.

## Troubleshooting

### Ollama is not accessible

If you encounter connection errors to Ollama:
1. Check that Ollama is running.
2. By default, Ollama must be accessible at `http://localhost:11434` or the host and port specified by the OLLAMA_HOST environment variable.
3. If your Ollama instance is running on a different host or port, use the `--host` and `--port` flags:
   ```bash
   rlama --host 192.168.1.100 --port 8000 list
   rlama --host my-ollama-server --port 11434 run my-rag
   ```
4. Check Ollama logs for potential errors.

### Text extraction issues

If you encounter problems with certain formats:
1. Install dependencies via `./scripts/install_deps.sh`.
2. Verify that your system has the required tools (`pdftotext`, `tesseract`, etc.).

### The RAG doesn't find relevant information

If the answers are not relevant:
1. Check that the documents are properly indexed with `rlama list`.
2. Make sure the content of the documents is properly extracted.
3. Try rephrasing your question more precisely.
4. Consider adjusting chunking parameters during RAG creation

### Other issues

For any other issues, please open an issue on the [GitHub repository](https://github.com/dontizi/rlama/issues) providing:
1. The exact command used.
2. The complete output of the command.
3. Your operating system and architecture.
4. The RLAMA version (`rlama --version`).

### Configuring Ollama Connection

RLAMA provides multiple ways to connect to your Ollama instance:

1. **Command-line flags** (highest priority):
   ```bash
   rlama --host 192.168.1.100 --port 8080 run my-rag
   ```

2. **Environment variable**:
   ```bash
   # Format: "host:port" or just "host"
   export OLLAMA_HOST=remote-server:8080
   rlama run my-rag
   ```

3. **Default values** (used if no other method is specified):
   - Host: `localhost`
   - Port: `11434`

The precedence order is: command-line flags > environment variable > default values.

## Advanced Usage

### Context Size Management

```bash
# Quick answers with minimal context
rlama run my-docs --context-size=10

# Deep analysis with maximum context
rlama run my-docs --context-size=50

# Balance between speed and depth
rlama run my-docs --context-size=30
```

### RAG Creation with Filtering
```bash
rlama rag llama3 my-project ./code \
  --exclude-dir=node_modules,dist \
  --process-ext=.go,.ts \
  --exclude-ext=.spec.ts
```

### Chunk Inspection
```bash
# List chunks with content preview
rlama list-chunks my-project --show-content

# Filter chunks from specific document
rlama list-chunks my-project --document=architecture
```

## Help System

Get full command help:
```bash
rlama --help
```

Command-specific help:
```bash
rlama rag --help
rlama list-chunks --help
rlama update-model --help
```

All commands support the global `--host` and `--port` flags for custom Ollama connections.

The precedence order is: command-line flags > environment variable > default values.