# BGE ONNX Reranker Implementation

This document describes the Go-native BGE reranker implementation using ONNX runtime.

## Overview

The BGE ONNX reranker provides a faster alternative to the original Python subprocess-based implementation by using:

1. **Pre-exported ONNX models** - No need to export models yourself
2. **Python ONNX microservice** - Runs ONNX inference in a dedicated HTTP server
3. **Go HTTP client** - Communicates with the microservice for reranking

## Architecture

```
┌─────────────────┐    HTTP     ┌──────────────────────┐
│ Go Application  │ ──────────► │ Python ONNX Server   │
│                 │             │                      │
│ BGEONNXReranker │             │ - onnxruntime        │
│ Client          │             │ - transformers       │
│                 │             │ - model.onnx         │
└─────────────────┘             └──────────────────────┘
```

## Performance Benefits

The ONNX implementation provides significant performance improvements:

- **8-15 seconds** vs 20-30 seconds for the original PyTorch models
- **Persistent server** - No subprocess startup overhead
- **Optimized inference** - ONNX runtime optimizations
- **Batch processing** - Multiple pairs in single request

## Setup Requirements

### 1. Download Pre-exported ONNX Model

```bash
mkdir -p ./models
cd ./models
git clone https://huggingface.co/corto-ai/bge-reranker-large-onnx
```

### 2. Install Python Dependencies

```bash
pip install onnxruntime transformers numpy
```

### 3. Verify Installation

```bash
go test ./internal/client -v -run TestBGEONNXRerankerClient
```

## Usage

### Basic Usage

```go
import "github.com/dontizi/rlama/internal/client"

// Create ONNX reranker client
modelDir := "./models/bge-reranker-large-onnx"
client, err := client.NewBGEONNXRerankerClient(modelDir)
if err != nil {
    log.Fatal(err)
}
defer client.Cleanup() // Important: stops the Python server

// Rerank query-passage pairs
pairs := [][]string{
    {"What is a cat?", "A cat is a small domesticated carnivorous mammal."},
    {"What is a cat?", "The weather is nice today."},
}

scores, err := client.ComputeScores(pairs, true) // true = normalize scores
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Scores: %v\n", scores) // [0.95, 0.05] (first pair more relevant)
```

### Integration with Existing Reranker Service

The ONNX client implements the same interface as the original BGE client:

```go
type RerankerClient interface {
    ComputeScores(pairs [][]string, normalize bool) ([]float64, error)
    GetModelName() string
}
```

To integrate, modify the reranker service to choose between implementations:

```go
func NewRerankerClient(modelName string, useONNX bool) RerankerClient {
    if useONNX {
        modelDir := "./models/bge-reranker-large-onnx"
        return client.NewBGEONNXRerankerClient(modelDir)
    }
    return client.NewBGERerankerClient(modelName)
}
```

## Available ONNX Models

Several pre-exported ONNX models are available on Hugging Face:

- `corto-ai/bge-reranker-large-onnx` - Standard ONNX version (recommended)
- `swulling/bge-reranker-large-onnx-o4` - O4 optimized version
- `swulling/bge-reranker-base-onnx-o4` - Base model, O4 optimized  
- `EmbeddedLLM/bge-reranker-base-onnx-o4-o2-gpu` - GPU optimized

## Implementation Details

### Microservice Approach

The implementation uses a Python HTTP server that:

1. **Loads ONNX model** once at startup
2. **Tokenizes input** using HuggingFace transformers
3. **Runs ONNX inference** with optimized runtime
4. **Returns scores** via JSON API

### Input Format

The BGE reranker expects input in the format:
```
query + " </s> " + passage
```

### Output Format

- **Normalized scores**: Sigmoid applied to logits (0.0 to 1.0)
- **Raw scores**: Direct logits output (any real number)

### Error Handling

The client handles common errors:
- Invalid pair format (not exactly 2 elements)
- Server connection failures
- ONNX runtime errors
- Tokenization errors

## Testing

Run the test suite to verify functionality:

```bash
# Basic functionality tests
go test ./internal/client -v -run TestBGEONNXRerankerClient

# Performance tests  
go test ./internal/client -v -run TestBGEONNXRerankerClient_Performance

# Benchmark against original implementation
go test ./internal/client -bench=BenchmarkBGEReranker
```

## Troubleshooting

### Common Issues

1. **"Model directory not found"**
   - Ensure ONNX model is downloaded to correct path
   - Check file permissions

2. **"Failed to start Python server"**
   - Verify Python dependencies are installed
   - Check port 8765 is available
   - Ensure Python is in PATH

3. **"Invalid input name: token_type_ids"**
   - This indicates ONNX model doesn't expect token_type_ids
   - Fixed in current implementation

### Performance Tuning

1. **Batch Size**: Process multiple pairs in single request
2. **Server Persistence**: Keep server running between requests
3. **Model Selection**: Use base model for faster inference if acceptable

## Future Improvements

1. **Pure Go Implementation**: Direct ONNX runtime without Python
2. **GPU Acceleration**: Use CUDA-enabled ONNX models
3. **Model Caching**: Cache tokenizer and model in memory
4. **Connection Pooling**: Reuse HTTP connections

## References

- [ONNX Runtime Go bindings](https://github.com/yalue/onnxruntime_go)
- [BGE Reranker Paper](https://arxiv.org/abs/2309.07597)
- [ONNX Model Hub](https://huggingface.co/models?library=onnx)