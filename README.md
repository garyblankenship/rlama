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
RLAMA is a powerful AI-driven question-answering tool for your documents that works with multiple LLM providers. It seamlessly integrates with Ollama, OpenAI, and any OpenAI-compatible endpoints (like LM Studio, VLLM, Text Generation Inference, etc.). RLAMA enables you to create, manage, and interact with Retrieval-Augmented Generation (RAG) systems tailored to your documentation needs.


[![RLAMA Demonstration](https://img.youtube.com/vi/EIsQnBqeQxQ/0.jpg)](https://www.youtube.com/watch?v=EIsQnBqeQxQ)

## Table of Contents
- [Vision & Roadmap](#vision--roadmap)
- [Installation](#installation)
- [Available Commands](#available-commands)
  - [rag - Create a RAG system](#rag---create-a-rag-system)
  - [crawl-rag - Create a RAG system from a website](#crawl-rag---create-a-rag-system-from-a-website)
  - [wizard - Create a RAG system with interactive setup](#wizard---create-a-rag-system-with-interactive-setup)
  - [watch - Set up directory watching for a RAG system](#watch---set-up-directory-watching-for-a-rag-system)
  - [watch-off - Disable directory watching for a RAG system](#watch-off---disable-directory-watching-for-a-rag-system)
  - [check-watched - Check a RAG's watched directory for new files](#check-watched---check-a-rags-watched-directory-for-new-files)
  - [web-watch - Set up website monitoring for a RAG system](#web-watch---set-up-website-monitoring-for-a-rag-system)
  - [web-watch-off - Disable website monitoring for a RAG system](#web-watch-off---disable-website-monitoring-for-a-rag-system)
  - [check-web-watched - Check a RAG's monitored website for updates](#check-web-watched---check-a-rags-monitored-website-for-updates)
  - [run - Use a RAG system](#run---use-a-rag-system)
  - [api - Start API server](#api---start-api-server)
  - [list - List RAG systems](#list---list-rag-systems)
  - [delete - Delete a RAG system](#delete---delete-a-rag-system)
  - [list-docs - List documents in a RAG](#list-docs---list-documents-in-a-rag)
  - [list-chunks - Inspect document chunks](#list-chunks---inspect-document-chunks)
  - [view-chunk - View chunk details](#view-chunk---view-chunk-details)
  - [add-docs - Add documents to RAG](#add-docs---add-documents-to-rag)
  - [crawl-add-docs - Add website content to RAG](#crawl-add-docs---add-website-content-to-rag)
  - [migrate-to-qdrant - Migrate RAG to Qdrant](#migration-between-vector-stores)
  - [migrate-to-internal - Migrate RAG to internal storage](#migration-between-vector-stores)
  - [migrate-batch - Batch migrate multiple RAGs](#migration-between-vector-stores)
  - [update-model - Change LLM model](#update-model---change-llm-model)
  - [profile - Manage API profiles](#profile---manage-api-profiles)
  - [update - Update RLAMA](#update---update-rlama)
  - [version - Display version](#version---display-version)
  - [hf-browse - Browse GGUF models on Hugging Face](#hf-browse---browse-gguf-models-on-hugging-face)
  - [run-hf - Run a Hugging Face GGUF model](#run-hf---run-a-hugging-face-gguf-model)
- [Uninstallation](#uninstallation)
- [Supported Document Formats](#supported-document-formats)
- [Enhanced Document Processing](docs/enhanced_document_processing.md)
- [Troubleshooting](#troubleshooting)
- [Model Support & LLM Providers](#model-support--llm-providers)
- [Managing API Profiles](#managing-api-profiles)

## Vision & Roadmap
RLAMA aims to become the definitive tool for creating local RAG systems that work seamlessly for everyone—from individual developers to large enterprises. Here's our strategic roadmap:

### Completed Features ✅
- ✅ **Basic RAG System Creation**: CLI tool for creating and managing RAG systems
- ✅ **Enhanced Document Processing**: LangChain-powered processing with 80% performance improvement and robust error handling
- ✅ **Document Chunking**: Advanced semantic chunking with multiple strategies (fixed, semantic, hierarchical, hybrid)
- ✅ **Vector Storage**: Local storage of document embeddings + Qdrant vector database integration
- ✅ **Production Vector Database**: Full Qdrant integration with gRPC/HTTP support, Qdrant Cloud compatibility
- ✅ **Seamless Migration Tools**: Complete migration system between internal and Qdrant storage with data integrity verification
- ✅ **Batch Operations**: Bulk migration of multiple RAGs with progress tracking and error recovery
- ✅ **Context Retrieval**: Basic semantic search with configurable context size
- ✅ **Ollama Integration**: Seamless connection to Ollama models
- ✅ **OpenAI Integration**: Full OpenAI API compatibility with profile management
- ✅ **Cross-Platform Support**: Works on Linux, macOS, and Windows
- ✅ **Easy Installation**: One-line installation script
- ✅ **API Server**: HTTP endpoints for integrating RAG capabilities in other applications
- ✅ **Web Crawling**: Create RAGs directly from websites
- ✅ **Guided RAG Setup Wizard**: Interactive interface for easy RAG creation
- ✅ **Hugging Face Integration**: Access to 45,000+ GGUF models from Hugging Face Hub
- ✅ **Advanced Reranking**: BGE reranker integration for improved search accuracy
- ✅ **Pure Go Tokenization**: 80% reduction in Python dependencies with 43K tokenizations/second

### 🔥 **Pure Go Evolution Roadmap (Q2 2025)**
Building on the breakthrough pure Go tokenization achievement:

- ✅ **Phase 1 Complete**: Pure Go tokenization with hybrid Python inference (80% dependency reduction)
- [ ] **Phase 2**: Full pure Go ONNX runtime integration (100% Python elimination)
- [ ] **Phase 3**: GPU-accelerated pure Go inference (CUDA/OpenCL support)
- [ ] **Phase 4**: Custom pure Go reranker models (domain-specific fine-tuning)

### Deployment & Performance Optimization (Q2 2025)
- ✅ **Pure Go Hybrid Architecture**: Achieved 80% Python dependency reduction with 43K tokenizations/sec
- [ ] **Full Pure Go ONNX Runtime**: Complete elimination of Python dependencies for 100% Go deployment
- [ ] **GPU Acceleration**: CUDA/OpenCL support for pure Go ONNX inference
- [ ] **Model Quantization**: INT8/FP16 quantization for faster inference
- [ ] **Prompt Compression**: Smart context summarization for limited context windows
- ✅ **Adaptive Chunking**: Dynamic content segmentation based on semantic boundaries and document structure
- ✅ **Minimal Context Retrieval**: Intelligent filtering to eliminate redundant content
- [ ] **Parameter Optimization**: Fine-tuned settings for different model sizes

### Advanced Search & Filtering (Q2 2025)
- [ ] **Enhanced Metadata Filtering**: Advanced search with document type, date, author, and custom metadata filters
- [ ] **Structured Query Language**: SQL-like queries for complex document retrieval
- [ ] **Faceted Search**: Multi-dimensional filtering with result counts
- [ ] **Similarity Thresholds**: Configurable relevance scoring and filtering

### Performance & Reliability (Q2-Q3 2025)
- ✅ **Thread-Safe Operations**: Concurrent tokenization and reranking with comprehensive testing
- ✅ **Memory Optimization**: Efficient Go memory management with stress testing
- ✅ **Error Recovery**: Robust error handling and fallback mechanisms
- [ ] **Connection Pooling**: Optimized Qdrant connections for high-throughput scenarios
- [ ] **Async Operations**: Non-blocking operations for large document imports
- [ ] **Caching Layer**: Smart caching for frequently accessed data and embeddings
- [ ] **Health Monitoring**: System health checks and performance metrics
- [ ] **Auto-retry Logic**: Exponential backoff for network failures

### Enhanced CLI & Developer Experience (Q2-Q3 2025)
- ✅ **Comprehensive Testing**: Unit tests, integration tests, and performance benchmarks for pure Go implementation
- ✅ **Advanced Error Handling**: Detailed error messages and troubleshooting for tokenization and reranking
- [ ] **RAG Status & Diagnostics**: `rag status`, `rag health-check`, `rag benchmark` commands
- [ ] **Performance Analytics**: Query performance metrics and optimization suggestions
- [ ] **Advanced Debugging**: Detailed logging, search result explanations, and troubleshooting tools

### Multi-Vector Store Ecosystem (Q3 2025)
- [ ] **Additional Vector Databases**: Support for Pinecone, Weaviate, Chroma
- [ ] **Pluggable Architecture**: Easy integration of new vector store backends
- [ ] **Performance Comparisons**: Built-in benchmarking between different vector stores
- [ ] **Cross-Store Migration**: Migration tools between any supported vector databases

### User Experience Enhancements (Q3-Q4 2025)
- [ ] **Lightweight Web Interface**: Simple browser-based UI for the existing CLI backend
- [ ] **Knowledge Graph Visualization**: Interactive exploration of document connections
- [ ] **Domain-Specific Templates**: Pre-configured settings for different domains

### Enterprise Features (Q4 2025)
- [ ] **Multi-User Access Control**: Role-based permissions for team environments
- [ ] **Integration with Enterprise Systems**: Connectors for SharePoint, Confluence, Google Workspace
- [ ] **Knowledge Quality Monitoring**: Detection of outdated or contradictory information
- [ ] **System Integration API**: Webhooks and APIs for embedding RLAMA in existing workflows
- [ ] **AI Agent Creation Framework**: Simplified system for building custom AI agents with RAG capabilities

### Next-Gen Retrieval Innovations (Q1 2026)
- [ ] **Custom Reranker Training**: Fine-tune reranking models on domain-specific data
- [ ] **Multi-Modal Reranking**: Support for image and text reranking  
- [ ] **Multi-Step Retrieval**: Using the LLM to refine search queries for complex questions
- [ ] **Cross-Modal Retrieval**: Support for image content understanding and retrieval
- [ ] **Feedback-Based Optimization**: Learning from user interactions to improve retrieval
- [ ] **Knowledge Graphs & Symbolic Reasoning**: Combining vector search with structured knowledge

### 🚀 **Current Status: Production-Ready with Simplified Deployment**

RLAMA has evolved from a simple local RAG tool to a comprehensive knowledge management platform that scales from individual developers to enterprise deployments. The recent breakthrough in **Pure Go Tokenization** represents a major achievement in deployment simplification, enabling users to:

- **Deploy Simply**: Single binary deployment with 80% fewer Python dependencies
- **Start Fast**: Pure Go tokenization at 43,000 tokens/second with zero setup
- **Scale Seamlessly**: Migrate to Qdrant for production workloads with zero data loss  
- **Enterprise Ready**: Deploy to Qdrant Cloud or self-hosted instances with full feature parity
- **Future Proof**: Built-in migration paths ensure no vendor lock-in

**Major Achievement: Pure Go Hybrid Architecture**
- ✅ **80% Python dependency reduction** - From complex Python stack to minimal inference-only requirements
- ✅ **43,000 tokenizations/second** - Blazing fast pure Go tokenization performance
- ✅ **Single binary deployment** - Most functionality requires only the `rlama` binary
- ✅ **Cross-platform compatibility** - Same performance benefits on all platforms
- ✅ **Production reliability** - Comprehensive testing with thread-safe operations

RLAMA's core philosophy remains unchanged: to provide a simple, powerful, local RAG solution that respects privacy, minimizes resource requirements, and works seamlessly across platforms. Now with the breakthrough pure Go implementation that makes deployment "more pleasing to use, less complicated for python execs."

## Installation

### Prerequisites
- **For Ollama models**: [Ollama](https://ollama.ai/) installed and running
- **For OpenAI models**: OpenAI API key or API profile configured
- **For OpenAI-compatible servers**: Local server running (e.g., LM Studio, VLLM, etc.)

### Installation from terminal

```bash
curl -fsSL https://raw.githubusercontent.com/dontizi/rlama/main/install.sh | sh
```

## Tech Stack

RLAMA is built with:

- **Core Language**: Go (chosen for performance, cross-platform compatibility, and single binary distribution)
- **CLI Framework**: Cobra (for command-line interface structure)
- **LLM Integration**: Multi-provider support (Ollama, OpenAI, OpenAI-compatible endpoints)
- **Reranking**: Pure Go BGE tokenization with hybrid Python inference (80% Python dependency reduction)
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

### Custom Data Directory

RLAMA stores data in `~/.rlama` by default. To use a different location:

1. **Command-line flag** (highest priority):
   ```bash
   # Use with any command
   rlama --data-dir /path/to/custom/directory run my-rag
   ```

2. **Environment variable**:
   ```bash
   # Set the environment variable
   export RLAMA_DATA_DIR=/path/to/custom/directory
   rlama run my-rag
   ```

The precedence order is: command-line flag > environment variable > default location.

### rag - Create a RAG system

Creates a new RAG system by indexing all documents in the specified folder.

```bash
rlama rag [model] [rag-name] [folder-path]
```

**Parameters:**
- `model`: Name of the Ollama model to use (e.g., llama3, mistral, gemma) or a Hugging Face model using the format `hf.co/username/repository[:quantization]`.
- `rag-name`: Unique name to identify your RAG system.
- `folder-path`: Path to the folder containing your documents.

**Example:**

```bash
# Using a standard Ollama model
rlama rag llama3 documentation ./docs

# Using a Hugging Face model
rlama rag hf.co/bartowski/Llama-3.2-1B-Instruct-GGUF my-rag ./docs

# Using a Hugging Face model with specific quantization
rlama rag hf.co/mlabonne/Meta-Llama-3.1-8B-Instruct-abliterated-GGUF:Q5_K_M my-rag ./docs
```

#### Enhanced Document Processing

RLAMA features advanced document processing powered by LangChainGo for improved reliability and performance:

**Document Loading Strategies:**
- **Hybrid (Default)**: Try LangChain first, fallback to legacy on failure - provides the best balance of performance and reliability
- **LangChain**: Advanced document processing with robust error handling and cross-platform consistency
- **Legacy**: Original RLAMA processor with external tool support

**Using Enhanced Document Processing:**

```bash
# Use default hybrid strategy (recommended)
rlama rag llama3 my-rag ./docs

# Force LangChain processing only
rlama rag llama3 my-rag ./docs --loader-strategy=langchain

# Use legacy processor
rlama rag llama3 my-rag ./docs --loader-strategy=legacy

# Disable advanced processing entirely
rlama rag llama3 my-rag ./docs --advanced-loader=false
```

**Environment Variables:**
Configure document processing behavior globally using environment variables:

```bash
# Set default loading strategy
export RLAMA_LOADER_STRATEGY=hybrid           # hybrid, langchain, or legacy

# Enable/disable LangChain processing
export RLAMA_USE_LANGCHAIN_LOADER=true        # true or false

# Debug output for document loading
export RLAMA_DEBUG_LOADER=true                # true or false

# Configure processing timeouts and retries
export RLAMA_LOADER_TIMEOUT_MINUTES=5         # 1, 2, 3, 5, 10, 15, 30
export RLAMA_LOADER_MAX_RETRIES=3             # 1, 2, 3, 5

# Optimize chunking settings
export RLAMA_PREFERRED_CHUNK_SIZE=1000        # 500, 750, 1000, 1500, 2000
export RLAMA_PREFERRED_CHUNK_OVERLAP=200      # 50, 100, 150, 200, 250, 300
```

**Performance Benefits:**
- **80% faster** document processing with LangChain
- **Robust error handling** eliminates most document processing failures
- **Cross-platform consistency** through standardized APIs
- **Better format support** for PDFs, Word docs, and other complex formats

#### Vector Store Options

RLAMA supports multiple vector storage backends to meet different performance and scaling needs:

**Internal Vector Store (Default)**
- File-based vector storage suitable for local development and small to medium datasets
- No external dependencies required

**Qdrant Vector Database**
- High-performance vector search engine optimized for large-scale semantic search
- Excellent for production environments and large document collections
- Supports advanced filtering and metadata search capabilities

**Using Qdrant Vector Store:**

```bash
# Create a RAG with Qdrant vector store
rlama rag llama3 docs-rag ./docs --vector-store=qdrant

# Customize Qdrant connection
rlama rag llama3 prod-rag ./docs \
  --vector-store=qdrant \
  --qdrant-host=localhost \
  --qdrant-port=6333 \
  --qdrant-collection=my-documents

# Using Qdrant Cloud with API key
rlama rag llama3 cloud-rag ./docs \
  --vector-store=qdrant \
  --qdrant-host=xyz.qdrant.cloud \
  --qdrant-port=6333 \
  --qdrant-apikey=your-api-key \
  --qdrant-grpc
```

**Qdrant Configuration Options:**
- `--vector-store`: Specify "qdrant" to use Qdrant vector database
- `--qdrant-host`: Qdrant server hostname (default: localhost)
- `--qdrant-port`: Qdrant server port (default: 6333)
- `--qdrant-apikey`: API key for Qdrant Cloud or secured instances
- `--qdrant-collection`: Collection name (defaults to RAG name)
- `--qdrant-grpc`: Use gRPC for communication (recommended for performance)

#### Migration Between Vector Stores

RLAMA provides seamless migration tools to move RAG systems between different vector storage backends without losing data.

**Individual RAG Migration:**

```bash
# Migrate from internal to Qdrant
rlama migrate-to-qdrant my-existing-rag \
  --qdrant-host=localhost \
  --qdrant-port=6333 \
  --backup

# Migrate back to internal storage
rlama migrate-to-internal my-qdrant-rag --backup

# Migrate to Qdrant Cloud
rlama migrate-to-qdrant prod-docs \
  --qdrant-host=xyz.qdrant.cloud \
  --qdrant-apikey=your-api-key \
  --qdrant-grpc \
  --backup \
  --verify
```

**Batch Migration:**

```bash
# Migrate all internal RAGs to Qdrant
rlama migrate-batch --from=internal --to=qdrant \
  --qdrant-host=production-server.com \
  --backup \
  --continue-on-error

# Migrate specific RAGs
rlama migrate-batch --from=internal --to=qdrant \
  --rags=docs,wiki,knowledge-base \
  --qdrant-host=localhost
```

**Migration Features:**
- ✅ **Data Integrity**: Automatic verification of migrated data
- ✅ **Backup Support**: Optional backup creation before migration
- ✅ **Progress Tracking**: Real-time progress for large migrations
- ✅ **Error Recovery**: Continue batch operations even if individual RAGs fail
- ✅ **Cleanup Options**: Automatic removal of old data after successful migration

### crawl-rag - Create a RAG system from a website

Creates a new RAG system by crawling a website and indexing its content.

```bash
rlama crawl-rag [model] [rag-name] [website-url]
```

**Parameters:**
- `model`: Name of the Ollama model to use (e.g., llama3, mistral, gemma).
- `rag-name`: Unique name to identify your RAG system.
- `website-url`: URL of the website to crawl and index.

**Options:**
- `--max-depth`: Maximum crawl depth (default: 2)
- `--concurrency`: Number of concurrent crawlers (default: 5)
- `--exclude-path`: Paths to exclude from crawling (comma-separated)
- `--chunk-size`: Character count per chunk (default: 1000)
- `--chunk-overlap`: Overlap between chunks in characters (default: 200)
- `--chunking-strategy`: Chunking strategy to use (options: "fixed", "semantic", "hybrid", "hierarchical", default: "hybrid")

#### Chunking Strategies

RLAMA offers multiple advanced chunking strategies to optimize document retrieval:

- **Fixed**: Traditional chunking with fixed size and overlap, respecting sentence boundaries when possible.
- **Semantic**: Intelligently splits documents based on semantic boundaries like headings, paragraphs, and natural topic shifts.
- **Hybrid**: Automatically selects the best strategy based on document type and content (markdown, HTML, code, or plain text).
- **Hierarchical**: For very long documents, creates a two-level chunking structure with major sections and sub-chunks.

The system automatically adapts to different document types:
- Markdown documents: Split by headers and sections
- HTML documents: Split by semantic HTML elements
- Code documents: Split by functions, classes, and logical blocks
- Plain text: Split by paragraphs with contextual overlap

#### Reranking Options

RLAMA includes advanced BGE-based reranking by default to improve result quality. Multiple implementations are available:

**🚀 Pure Go Hybrid Reranker (Recommended)**
- **Pure Go tokenization** with **43,000 tokenizations/second**
- **80% reduction** in Python dependency complexity
- **Single binary deployment** for most functionality
- Works out of the box with zero additional setup

**Python BGE Reranker (Legacy)**
- Uses the original Python FlagEmbedding library via subprocess calls
- Requires full Python environment

**ONNX BGE Reranker (Transitional)**
- Uses optimized ONNX models for **3.8x faster performance** over original Python
- Requires one-time setup but provides significant speed improvements

```bash
# Download ONNX model (one-time setup)
mkdir -p ./models
cd ./models
git clone https://huggingface.co/corto-ai/bge-reranker-large-onnx

# Use ONNX reranker for faster performance
rlama rag llama3.2 myrag ./docs --use-onnx-reranker

# Specify custom ONNX model directory
rlama rag llama3.2 myrag ./docs --use-onnx-reranker --onnx-model-dir ./models/bge-reranker-large-onnx
```

**ONNX Requirements:**
```bash
pip install onnxruntime transformers numpy
```

**Reranking Configuration Options:**
- `--disable-reranker`: Disable reranking (enabled by default)
- `--reranker-model`: Model to use for reranking (defaults to main model)
- `--reranker-weight`: Weight for reranker scores vs vector scores (0-1, default: 0.7)
- `--reranker-threshold`: Minimum score threshold for reranked results (default: 0.0)
- `--use-onnx-reranker`: Use ONNX reranker for faster performance
- `--onnx-model-dir`: Directory containing ONNX reranker model (default: ./models/bge-reranker-large-onnx)

**Performance Comparison:**
- **🚀 Pure Go Hybrid**: **43,000 tokenizations/second**, minimal Python dependencies
- **ONNX BGE**: ~2.0 seconds per query (**3.8x faster** than original Python)
- **Python BGE**: ~7.4 seconds per query

**Example:**

```bash
# Create a new RAG from a documentation website
rlama crawl-rag llama3 docs-rag https://docs.example.com

# Customize crawling behavior
rlama crawl-rag llama3 blog-rag https://blog.example.com --max-depth=3 --exclude-path=/archive,/tags

# Create a RAG with semantic chunking
rlama rag llama3 documentation ./docs --chunking-strategy=semantic

# Use hierarchical chunking for large documents
rlama rag llama3 book-rag ./books --chunking-strategy=hierarchical
```

### wizard - Create a RAG system with interactive setup

Provides an interactive step-by-step wizard for creating a new RAG system.

```bash
rlama wizard
```

The wizard guides you through:
- Naming your RAG
- Choosing an Ollama model
- Selecting document sources (local folder or website)
- Configuring chunking parameters
- Setting up file filtering

**Example:**

```bash
rlama wizard
# Follow the prompts to create your customized RAG
```

### watch - Set up directory watching for a RAG system

Configure a RAG system to automatically watch a directory for new files and add them to the RAG.

```bash
rlama watch [rag-name] [directory-path] [interval]
```

**Parameters:**
- `rag-name`: Name of the RAG system to watch.
- `directory-path`: Path to the directory to watch for new files.
- `interval`: Time in minutes to check for new files (use 0 to check only when the RAG is used).

**Example:**

```bash
# Set up directory watching to check every 60 minutes
rlama watch my-docs ./watched-folder 60

# Set up directory watching to only check when the RAG is used
rlama watch my-docs ./watched-folder 0

# Customize what files to watch
rlama watch my-docs ./watched-folder 30 --exclude-dir=node_modules,tmp --process-ext=.md,.txt
```

### watch-off - Disable directory watching for a RAG system

Disable automatic directory watching for a RAG system.

```bash
rlama watch-off [rag-name]
```

**Parameters:**
- `rag-name`: Name of the RAG system to disable watching.

**Example:**

```bash
rlama watch-off my-docs
```

### check-watched - Check a RAG's watched directory for new files

Manually check a RAG's watched directory for new files and add them to the RAG.

```bash
rlama check-watched [rag-name]
```

**Parameters:**
- `rag-name`: Name of the RAG system to check.

**Example:**

```bash
rlama check-watched my-docs
```

### web-watch - Set up website monitoring for a RAG system

Configure a RAG system to automatically monitor a website for updates and add new content to the RAG.

```bash
rlama web-watch [rag-name] [website-url] [interval]
```

**Parameters:**
- `rag-name`: Name of the RAG system to monitor.
- `website-url`: URL of the website to monitor.
- `interval`: Time in minutes between checks (use 0 to check only when the RAG is used).

**Example:**

```bash
# Set up website monitoring to check every 60 minutes
rlama web-watch my-docs https://example.com 60

# Set up website monitoring to only check when the RAG is used
rlama web-watch my-docs https://example.com 0

# Customize what content to monitor
rlama web-watch my-docs https://example.com 30 --exclude-path=/archive,/tags
```

### web-watch-off - Disable website monitoring for a RAG system

Disable automatic website monitoring for a RAG system.

```bash
rlama web-watch-off [rag-name]
```

**Parameters:**
- `rag-name`: Name of the RAG system to disable monitoring.

**Example:**

```bash
rlama web-watch-off my-docs
```

### check-web-watched - Check a RAG's monitored website for updates

Manually check a RAG's monitored website for new updates and add them to the RAG.

```bash
rlama check-web-watched [rag-name]
```

**Parameters:**
- `rag-name`: Name of the RAG system to check.

**Example:**

```bash
rlama check-web-watched my-docs
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

### api - Start API server

Starts an HTTP API server that exposes RLAMA's functionality through RESTful endpoints.

```bash
rlama api [--port PORT]
```

**Parameters:**
- `--port`: (Optional) Port number to run the API server on (default: 11249)

**Example:**

```bash
rlama api --port 8080
```

**Available Endpoints:**

1. **Query a RAG system** - `POST /rag`
   ```bash
   curl -X POST http://localhost:11249/rag \
     -H "Content-Type: application/json" \
     -d '{
       "rag_name": "documentation",
       "prompt": "How do I install the project?",
       "context_size": 20
     }'
   ```

   Request fields:
   - `rag_name` (required): Name of the RAG system to query
   - `prompt` (required): Question or prompt to send to the RAG
   - `context_size` (optional): Number of chunks to include in context
   - `model` (optional): Override the model used by the RAG

2. **Check server health** - `GET /health`
   ```bash
   curl http://localhost:11249/health
   ```

**Integration Example:**
```javascript
// Node.js example
const response = await fetch('http://localhost:11249/rag', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    rag_name: 'my-docs',
    prompt: 'Summarize the key features'
  })
});
const data = await response.json();
console.log(data.response);
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

### crawl-add-docs - Add website content to RAG

Add content from a website to an existing RAG system.

```bash
rlama crawl-add-docs [rag-name] [website-url]
```

**Parameters:**
- `rag-name`: Name of the RAG system
- `website-url`: URL of the website to crawl and add to the RAG

**Options:**
- `--max-depth`: Maximum crawl depth (default: 2)
- `--concurrency`: Number of concurrent crawlers (default: 5)
- `--exclude-path`: Paths to exclude from crawling (comma-separated)
- `--chunk-size`: Character count per chunk (default: 1000)
- `--chunk-overlap`: Overlap between chunks in characters (default: 200)

**Example:**

```bash
# Add blog content to an existing RAG
rlama crawl-add-docs my-docs https://blog.example.com

# Customize crawling behavior
rlama crawl-add-docs knowledge-base https://docs.example.com --max-depth=1 --exclude-path=/api
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

### profile - Manage API profiles

Manage API profiles for different LLM providers and endpoints.

#### profile add - Create a new profile

```bash
rlama profile add [name] [provider] [api-key] [flags]
```

**Parameters:**
- `name`: Unique name for the profile
- `provider`: Provider type (`openai` or `openai-api`)
- `api-key`: API key (use "none" for local servers without authentication)

**Flags:**
- `--base-url`: Custom base URL for OpenAI-compatible endpoints (required for `openai-api` provider)

**Examples:**

```bash
# Traditional OpenAI profile
rlama profile add openai-work openai sk-your-api-key

# LM Studio local server
rlama profile add lmstudio openai-api none --base-url http://localhost:1234/v1

# VLLM server with authentication
rlama profile add vllm openai-api your-token --base-url http://server:8000/v1
```

#### profile list - List all profiles

```bash
rlama profile list
```

#### profile delete - Delete a profile

```bash
rlama profile delete [name]
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

### hf-browse - Browse GGUF models on Hugging Face

Search and browse GGUF models available on Hugging Face.

```bash
rlama hf-browse [search-term] [flags]
```

**Parameters:**
- `search-term`: (Optional) Term to search for (e.g., "llama3", "mistral")

**Flags:**
- `--open`: Open the search results in your default web browser
- `--quant`: Specify quantization type to suggest (e.g., Q4_K_M, Q5_K_M)
- `--limit`: Limit number of results (default: 10)

**Examples:**

```bash
# Search for GGUF models and show command-line help
rlama hf-browse "llama 3"

# Open browser with search results
rlama hf-browse mistral --open

# Search with specific quantization suggestion
rlama hf-browse phi --quant Q4_K_M
```

### run-hf - Run a Hugging Face GGUF model

Run a Hugging Face GGUF model directly using Ollama. This is useful for testing models before creating a RAG system with them.

```bash
rlama run-hf [huggingface-model] [flags]
```

**Parameters:**
- `huggingface-model`: Hugging Face model path in the format `username/repository`

**Flags:**
- `--quant`: Quantization to use (e.g., Q4_K_M, Q5_K_M)

**Examples:**

```bash
# Try a model in chat mode
rlama run-hf bartowski/Llama-3.2-1B-Instruct-GGUF

# Specify quantization
rlama run-hf mlabonne/Meta-Llama-3.1-8B-Instruct-abliterated-GGUF --quant Q5_K_M
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

RLAMA supports many file formats with enhanced processing capabilities:

- **Text**: `.txt`, `.md`, `.html`, `.json`, `.csv`, `.yaml`, `.yml`, `.xml`, `.org`
- **Code**: `.go`, `.py`, `.js`, `.java`, `.c`, `.cpp`, `.cxx`, `.h`, `.rb`, `.php`, `.rs`, `.swift`, `.kt`, `.ts`, `.tsx`, `.f`, `.F`, `.F90`, `.el`, `.svelte`, `.jsx`, `.vue`
- **Documents**: `.pdf`, `.docx`, `.doc`, `.rtf`, `.odt`, `.pptx`, `.ppt`, `.xlsx`, `.xls`, `.epub`

**Enhanced Processing Features:**
- **LangChain Integration**: Advanced document processing with improved reliability and 80% faster performance
- **Robust Error Handling**: Graceful fallback mechanisms prevent processing failures
- **Cross-Platform Consistency**: Standardized document processing across Windows, Mac, and Linux
- **Automatic Format Detection**: Intelligent processing based on file content and metadata

**Processing Strategy Options:**
- Use `--loader-strategy=hybrid` for the best balance of reliability and performance (default)
- Use `--loader-strategy=langchain` for maximum performance with advanced features
- Use `--loader-strategy=legacy` for compatibility with external tools

Installing dependencies via `install_deps.sh` is recommended to improve support for certain formats when using the legacy processor.

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

## Hugging Face Integration

RLAMA now supports using GGUF models directly from Hugging Face through Ollama's native integration:

### Browsing Hugging Face Models

```bash
# Search for GGUF models on Hugging Face
rlama hf-browse "llama 3"

# Open browser with search results
rlama hf-browse mistral --open
```

### Testing a Model

Before creating a RAG, you can test a Hugging Face model directly:

```bash
# Try a model in chat mode
rlama run-hf bartowski/Llama-3.2-1B-Instruct-GGUF

# Specify quantization
rlama run-hf mlabonne/Meta-Llama-3.1-8B-Instruct-abliterated-GGUF --quant Q5_K_M
```

### Creating a RAG with Hugging Face Models

Use Hugging Face models when creating RAG systems:

```bash
# Create a RAG with a Hugging Face model
rlama rag hf.co/bartowski/Llama-3.2-1B-Instruct-GGUF my-rag ./docs

# Use specific quantization
rlama rag hf.co/mlabonne/Meta-Llama-3.1-8B-Instruct-abliterated-GGUF:Q5_K_M my-rag ./docs
```

## Model Support & LLM Providers

RLAMA supports multiple LLM providers for both **text generation** and **embeddings**:

### Supported Providers

1. **Ollama** (default): Local models via Ollama server
2. **OpenAI**: Official OpenAI API  
3. **OpenAI-Compatible**: Any server implementing OpenAI API (LM Studio, VLLM, TGI, etc.)

### How Models Are Used

RLAMA uses models for two distinct purposes:

- **Text Generation (Completions)**: Answering your questions using retrieved context
- **Embeddings**: Converting documents and queries into vectors for similarity search

### Model Selection Logic

When you specify a model name, RLAMA automatically determines which provider to use:

- **OpenAI models** (e.g., `gpt-4`, `gpt-3.5-turbo`): Uses OpenAI API for completions + embeddings
- **Hugging Face models** (e.g., `hf.co/username/model`): Downloads via Ollama
- **Other models** (e.g., `llama3`, `mistral`): Uses Ollama for completions + embeddings

### Using OpenAI Models

Set your OpenAI API key:
```bash
export OPENAI_API_KEY="your-api-key"
```

Create a RAG with OpenAI model:
```bash
rlama rag gpt-4 my-rag ./documents
```

Supported OpenAI models:
- `gpt-4`, `gpt-4-turbo`, `gpt-4o`
- `gpt-3.5-turbo`
- `o3-mini` and newer models

### Using OpenAI-Compatible Endpoints

RLAMA can connect to any server that implements the OpenAI API specification, including:

- **LM Studio**: Local model serving with OpenAI API
- **VLLM**: High-performance inference server  
- **Text Generation Inference (TGI)**: Hugging Face's inference server
- **Ollama's OpenAI compatibility mode**: `ollama serve` with OpenAI endpoints
- **Any custom OpenAI-compatible server**

#### Setting Up Profiles for Custom Endpoints

Create a profile for your OpenAI-compatible server:

```bash
# For LM Studio running locally
rlama profile add lmstudio openai-api none --base-url http://localhost:1234/v1

# For VLLM server (with authentication)
rlama profile add vllm openai-api your-api-key --base-url http://your-server:8000/v1

# For remote TGI server
rlama profile add tgi openai-api dummy --base-url https://tgi.example.com/v1
```

#### Using Custom Endpoints

Create a RAG with your custom endpoint:

```bash
# Use the profile when creating a RAG
rlama rag llama-3-8b my-rag ./documents --profile lmstudio

# The model name should match what your server expects
rlama rag custom-model-name knowledge-base ./docs --profile vllm
```

#### Common OpenAI-Compatible Servers

1. **LM Studio**:
   ```bash
   # Start LM Studio with OpenAI API on default port 1234
   rlama profile add lmstudio openai-api none --base-url http://localhost:1234/v1
   rlama rag llama-3-8b-instruct my-docs ./documents --profile lmstudio
   ```

2. **VLLM**:
   ```bash
   # VLLM typically runs on port 8000
   rlama profile add vllm openai-api none --base-url http://localhost:8000/v1
   rlama rag meta-llama/Llama-3-8B-Instruct my-rag ./docs --profile vllm
   ```

3. **Ollama OpenAI Mode**:
   ```bash
   # If using Ollama's experimental OpenAI endpoints
   rlama profile add ollama-openai openai-api none --base-url http://localhost:11434/v1
   rlama rag llama3 my-rag ./docs --profile ollama-openai
   ```

#### Benefits of OpenAI-Compatible Mode

- **Unified Interface**: Same API for different inference engines
- **Easy Migration**: Switch between providers without changing RAG structure  
- **Better Performance**: Use optimized inference servers (VLLM, TGI)
- **Model Flexibility**: Access models not available through Ollama
- **Embedding Support**: Full support for both completions and embeddings

## Managing API Profiles

RLAMA allows you to create API profiles to manage multiple API keys and endpoints for different providers:

### Profile Types

- **`openai`**: Official OpenAI API profiles
- **`openai-api`**: Generic OpenAI-compatible endpoints (LM Studio, VLLM, etc.)

### Creating Profiles

#### Traditional OpenAI Profiles
```bash
# Create a profile for your OpenAI account
rlama profile add openai-work openai "sk-your-api-key"

# Create another profile for a different account
rlama profile add openai-personal openai "sk-your-personal-api-key" 
```

#### OpenAI-Compatible Endpoint Profiles
```bash
# LM Studio local server (no API key needed)
rlama profile add lmstudio openai-api none --base-url http://localhost:1234/v1

# VLLM server with authentication
rlama profile add vllm-server openai-api your-token --base-url http://192.168.1.100:8000/v1

# Remote TGI deployment
rlama profile add tgi-prod openai-api api-key --base-url https://api.mycompany.com/v1
```

### Listing Profiles

```bash
# View all available profiles with their base URLs
rlama profile list
```

Output example:
```
NAME         PROVIDER    BASE URL                  CREATED ON           LAST USED
openai-work  openai      default                   2024-01-15 10:30:25  2024-01-16 14:22:10
lmstudio     openai-api  http://localhost:1234/v1  2024-01-16 09:15:33  never
vllm-server  openai-api  http://server:8000/v1     2024-01-16 11:45:12  2024-01-16 15:30:25
```

### Deleting Profiles

```bash
# Delete a profile
rlama profile delete openai-old
```

### Using Profiles with RAGs

When creating a new RAG:

```bash
# Create a RAG with an OpenAI profile
rlama rag gpt-4 my-rag ./documents --profile openai-work

# Create a RAG with a custom endpoint
rlama rag llama-3-8b local-rag ./docs --profile lmstudio
```

When running existing RAGs:

```bash
# RAGs remember their original configuration automatically
rlama run my-rag
```

### Profile Benefits

- **Multiple Endpoints**: Manage connections to different LLM servers
- **Easy Switching**: Change between local and remote inference
- **Secure Storage**: API keys stored safely in `~/.rlama/profiles`
- **Usage Tracking**: See when profiles were last used
- **Project Organization**: Use different profiles for different projects
- **Development Workflow**: Test locally (LM Studio) → deploy remotely (VLLM)
