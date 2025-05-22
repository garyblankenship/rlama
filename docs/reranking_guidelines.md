# RLAMA Reranking Documentation

## Overview
Reranking in RLAMA is a feature that improves retrieval accuracy by applying a second-stage ranking to initial search results using a cross-encoder approach. RLAMA now includes a **pure Go implementation** that significantly reduces deployment complexity while maintaining excellent performance.

## Features
- **🚀 Pure Go tokenization** with 43,000 tokenizations/second
- **80% reduction** in Python dependency complexity  
- **Single binary deployment** for most functionality
- Enabled by default for all RAG systems
- Configurable weights between vector similarity and reranking scores
- Adjustable result limits and thresholds
- Custom model support for reranking

## Implementation Options

### 1. Pure Go Hybrid Reranker (Recommended)
- **Zero setup required** - works out of the box
- **Pure Go tokenization** - no Python dependencies for text processing
- **Minimal Python inference** - only ONNX model execution requires Python
- **43,000 tokenizations/second** performance
- **Cross-platform compatibility**

### 2. Legacy Python Implementation
- Original FlagEmbedding library approach
- Full Python dependency stack required
- Available for compatibility

## Default Configuration
- TopK: 5 results (maximum number of results after reranking)
- Initial retrieval: 20 documents
- Reranker weight: 0.7 (70% reranker score, 30% vector similarity)
- Score threshold: 0.0 (no minimum score requirement)
- Model: Uses the same model as the RAG system by default

## Usage

### Command Line Interface

1. **Configure Reranking for a RAG System**
```bash
rlama add-reranker my-rag [options]
```

Available options:
- `--model`: Specify a custom model for reranking (defaults to RAG model)
- `--weight`: Set the weight for reranker scores (0-1)
- `--threshold`: Set minimum score threshold for results
- `--topk`: Set maximum number of results to return
- `--disable`: Disable reranking for this RAG

Examples:
```bash
# Configure with custom model
rlama add-reranker my-rag --model reranker-model

# Adjust weights and limits
rlama add-reranker my-rag --weight 0.8 --topk 10

# Disable reranking
rlama add-reranker my-rag --disable
```

### Programmatic Usage

1. **Creating a RAG with Reranking**
```go
err := ragService.CreateRagWithOptions("llama3.2", "my-rag", documentPath, service.DocumentLoaderOptions{
    ChunkSize: 200,
    ChunkOverlap: 50,
    EnableReranker: true,  // Reranking is enabled by default
})
```

2. **Customizing Reranking Options**
```go
options := service.RerankerOptions{
    TopK: 10,                // Return top 10 results
    InitialK: 30,           // Retrieve 30 initial results
    RerankerModel: "custom-model",  // Use custom model
    ScoreThreshold: 0.5,    // Minimum relevance score
    RerankerWeight: 0.8,    // 80% reranker, 20% vector similarity
}
```

## How It Works

1. **Initial Retrieval**: The system first retrieves an initial set of documents using vector similarity search (default: top 20 documents).

2. **Reranking Process**:
   - Each retrieved document is evaluated using a cross-encoder model
   - The model scores document relevance on a scale of 0 to 1
   - Final scores combine vector similarity and reranking scores based on weights
   - Results are sorted by final score and limited to TopK

3. **Scoring Formula**:
finalScore = (rerankerWeight × rerankerScore) + ((1 - rerankerWeight) × vectorScore) 

## Performance Considerations

### Pure Go Implementation Benefits
- **Tokenization**: 43,000 tokenizations/second (zero Python overhead)
- **Startup time**: Instant - no subprocess initialization
- **Memory usage**: Optimized Go memory management
- **Deployment**: Single binary - no Python environment setup

### General Performance
- Reranking adds additional processing time as each document needs to be evaluated
- The InitialK parameter affects both accuracy and performance  
- Larger TopK values increase processing time
- Pure Go tokenization eliminates most Python-related latency
- Only ONNX model inference still requires minimal Python service

## Best Practices

1. **Model Selection**
   - Use the same model as your RAG system for consistency
   - Choose models that excel at cross-encoding for better results

2. **Parameter Tuning**
   - Start with default weights (0.7) and adjust based on results
   - Increase InitialK for better recall at the cost of performance
   - Set appropriate thresholds based on your use case

3. **Performance Optimization**
   - Limit TopK to necessary minimum
   - Consider chunk size impact on reranking performance
   - Monitor and adjust InitialK based on result quality

## Troubleshooting

Common issues and solutions:

1. **Slow Response Times**
   - Reduce InitialK or TopK values
   - Consider using a lighter reranking model
   - Check if chunk sizes are appropriate

2. **Poor Result Quality**
   - Increase reranker weight
   - Adjust score threshold
   - Increase InitialK for more candidate documents

3. **Model Compatibility**
   - Ensure the reranking model supports the required operations
   - Check model availability in your Ollama installation

## Advanced Configuration

For specific use cases, you can fine-tune the reranking system by:

1. **Custom Scoring**
   - Adjust weights based on document types
   - Implement custom thresholds for different queries

2. **Model Chaining**
   - Use different models for initial retrieval and reranking
   - Combine multiple reranking passes with different models

## Examples

### Basic Usage Examples

1. **Create RAG with Pure Go Reranking (Default)**
```bash
# Pure Go implementation is used automatically
rlama rag llama3.2 my-documents ./docs
```

2. **Query with Advanced Reranking**
```bash
rlama run my-documents
> What are the key features of machine learning?
# Uses pure Go tokenization + minimal Python inference
```

3. **Configure Reranking Parameters**
```bash
# High-precision configuration  
rlama add-reranker research-papers --weight 0.9 --threshold 0.7 --topk 3

# Large-scale configuration
rlama add-reranker large-corpus --topk 20 --weight 0.6
```

### Programmatic Usage Examples

1. **Pure Go BGE Client**
```go
package main

import (
    "context"
    "log"
    "github.com/dontizi/rlama/internal/client"
)

func main() {
    // Create pure Go BGE client (recommended)
    modelPath := "./models/bge-reranker-large-onnx"
    client, err := client.NewPureGoBGEClient(modelPath, true, "http://localhost:8000")
    if err != nil {
        log.Fatal(err)
    }

    // Rerank documents with pure Go tokenization
    ctx := context.Background()
    query := "What is artificial intelligence?"
    documents := []string{
        "AI is a branch of computer science.",
        "Machine learning is a subset of AI.", 
        "Deep learning uses neural networks.",
    }

    results, err := client.Rerank(ctx, query, documents, 3)
    if err != nil {
        log.Fatal(err)
    }

    for i, result := range results {
        log.Printf("Rank %d: Score=%.3f, Doc=%s", 
            i+1, result.Score, result.Document[:30]+"...")
    }
}
```

2. **Performance Testing**
```go
// Benchmark pure Go tokenization
func BenchmarkPureGoTokenization(b *testing.B) {
    client, _ := client.NewPureGoBGEClient(modelPath, true, fallbackURL)
    query := "Test query"
    passage := "Test passage for benchmarking"
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        tokenizer := client.GetTokenizer()
        _, _, _ = tokenizer.EncodeQueryPassagePair(query, passage, 512)
    }
    // Result: ~23 microseconds per tokenization (43K/sec)
}
```

3. **Configuration Options**
```go
// Configure pure Go client
client, err := client.NewPureGoBGEClient(modelPath, true, fallbackURL)

// Set maximum token length for optimization
client.SetMaxLength(256)  // Shorter sequences for speed

// Health check
if err := client.Health(ctx); err != nil {
    log.Printf("Client health check failed: %v", err)
}
```

### Integration Examples

1. **API Server Integration**
```bash
# Start RLAMA API server with pure Go reranking
rlama api --port 8080

# Query via HTTP API
curl -X POST http://localhost:8080/rag \
  -H "Content-Type: application/json" \
  -d '{
    "rag_name": "my-documents",
    "prompt": "Explain machine learning concepts",
    "context_size": 20
  }'
```

2. **Performance Monitoring**
```bash
# Run comprehensive tests
go test ./internal/client -v -run TestPureGoBGEClient_EndToEndReranking

# Benchmark performance
go test ./internal/client -bench=BenchmarkPureGoBGEClient_Tokenization -benchmem

# Test concurrent usage
go test ./internal/client -v -run TestPureGoTokenizer_ConcurrentAccess
```

### Migration Examples

1. **From Legacy Python to Pure Go**
```bash
# No migration needed - pure Go is used automatically in RLAMA v0.1.36+
rlama run existing-rag  # Uses pure Go tokenization automatically
```

2. **Custom Deployment**
```bash
# Simple deployment with pure Go benefits
./rlama run my-rag  # Single binary, minimal Python dependencies

# Compare with legacy approach (much more complex):
# pip install torch transformers numpy onnxruntime
# python bge_server.py &
# ./rlama run my-rag
```

## Testing and Validation

### Test Coverage Examples

```bash
# Test pure Go tokenizer
go test ./internal/client -v -run TestPureGoTokenizer_Basic

# Test error handling
go test ./internal/client -v -run TestPureGoTokenizer_ErrorHandling

# Test concurrent access
go test ./internal/client -v -run TestPureGoTokenizer_ConcurrentAccess

# Test BGE client integration
go test ./internal/client -v -run TestPureGoBGEClient_EndToEndReranking

# Memory stress testing
go test ./internal/client -v -run TestPureGoTokenizer_MemoryStress
```

### Performance Validation

```bash
# Benchmark tokenization speed
go test ./internal/client -bench=BenchmarkPureGoBGEClient_Tokenization
# Expected: ~23,000 ns/op (43K tokenizations/second)

# Benchmark concurrent performance  
go test ./internal/client -bench=BenchmarkPureGoBGEClient_Concurrent

# Test batch processing
go test ./internal/client -v -run TestPureGoBGEClient_BatchProcessing
```

This comprehensive guide demonstrates RLAMA's advanced reranking capabilities with the new pure Go implementation, showcasing both ease of use and high performance.