package service

import (
    "os"
    "path/filepath"
    "testing"

    "github.com/dontizi/rlama/internal/domain"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestNewSimpleDocumentProcessor(t *testing.T) {
    processor := NewSimpleDocumentProcessor()
    
    assert.NotNil(t, processor)
    assert.NotNil(t, processor.documentLoader)
    assert.NotNil(t, processor.chunkService)
}

func TestSimpleDocumentProcessor_ProcessContent(t *testing.T) {
    processor := NewSimpleDocumentProcessor()

    tests := []struct {
        name        string
        content     string
        config      ChunkingConfig
        expectError bool
        expectChunks bool
    }{
        {
            name:    "valid content with fixed chunking",
            content: "This is a test content that should be chunked properly. It has enough content to create multiple chunks when configured correctly.",
            config: ChunkingConfig{
                ChunkSize:        50,
                ChunkOverlap:     10,
                ChunkingStrategy: "fixed",
            },
            expectError:  false,
            expectChunks: true,
        },
        {
            name:    "empty content",
            content: "",
            config: ChunkingConfig{
                ChunkSize:        100,
                ChunkOverlap:     20,
                ChunkingStrategy: "fixed",
            },
            expectError:  true,
            expectChunks: false,
        },
        {
            name:    "whitespace only content",
            content: "   \n\t   ",
            config: ChunkingConfig{
                ChunkSize:        100,
                ChunkOverlap:     20,
                ChunkingStrategy: "fixed",
            },
            expectError:  true,
            expectChunks: false,
        },
        {
            name:    "invalid chunk size",
            content: "Valid content here",
            config: ChunkingConfig{
                ChunkSize:        0,
                ChunkOverlap:     10,
                ChunkingStrategy: "fixed",
            },
            expectError:  true,
            expectChunks: false,
        },
        {
            name:    "chunk overlap greater than size",
            content: "Valid content here",
            config: ChunkingConfig{
                ChunkSize:        50,
                ChunkOverlap:     60,
                ChunkingStrategy: "fixed",
            },
            expectError:  true,
            expectChunks: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            chunks, err := processor.ProcessContent(tt.content, tt.config)

            if tt.expectError {
                assert.Error(t, err)
                assert.Nil(t, chunks)
            } else {
                assert.NoError(t, err)
                if tt.expectChunks {
                    assert.NotEmpty(t, chunks)
                    for _, chunk := range chunks {
                        assert.NotEmpty(t, chunk.Content)
                        assert.GreaterOrEqual(t, chunk.ChunkNumber, 0)
                    }
                }
            }
        })
    }
}

func TestSimpleDocumentProcessor_EstimateChunkCount(t *testing.T) {
    processor := NewSimpleDocumentProcessor()

    tests := []struct {
        name           string
        content        string
        config         ChunkingConfig
        expectedMin    int
        expectedMax    int
    }{
        {
            name:    "empty content",
            content: "",
            config: ChunkingConfig{
                ChunkSize:    100,
                ChunkOverlap: 20,
            },
            expectedMin: 0,
            expectedMax: 0,
        },
        {
            name:    "small content",
            content: "Small content",
            config: ChunkingConfig{
                ChunkSize:    100,
                ChunkOverlap: 20,
            },
            expectedMin: 1,
            expectedMax: 1,
        },
        {
            name:    "large content",
            content: string(make([]rune, 500)), // 500 characters
            config: ChunkingConfig{
                ChunkSize:    100,
                ChunkOverlap: 20,
            },
            expectedMin: 6,
            expectedMax: 7,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            count := processor.EstimateChunkCount(tt.content, tt.config)
            assert.GreaterOrEqual(t, count, tt.expectedMin)
            assert.LessOrEqual(t, count, tt.expectedMax)
        })
    }
}

func TestSimpleDocumentProcessor_ValidateDocument(t *testing.T) {
    processor := NewSimpleDocumentProcessor()

    tests := []struct {
        name        string
        doc         *domain.Document
        expectError bool
    }{
        {
            name: "valid document",
            doc: &domain.Document{
                ID:      "test-id",
                Name:    "test.txt",
                Content: "This is valid content with sufficient length",
                Path:    "/path/to/test.txt",
            },
            expectError: false,
        },
        {
            name:        "nil document",
            doc:         nil,
            expectError: true,
        },
        {
            name: "empty ID",
            doc: &domain.Document{
                ID:      "",
                Name:    "test.txt",
                Content: "Valid content",
                Path:    "/path/to/test.txt",
            },
            expectError: true,
        },
        {
            name: "empty name",
            doc: &domain.Document{
                ID:      "test-id",
                Name:    "",
                Content: "Valid content",
                Path:    "/path/to/test.txt",
            },
            expectError: true,
        },
        {
            name: "empty content",
            doc: &domain.Document{
                ID:      "test-id",
                Name:    "test.txt",
                Content: "",
                Path:    "/path/to/test.txt",
            },
            expectError: true,
        },
        {
            name: "content too short",
            doc: &domain.Document{
                ID:      "test-id",
                Name:    "test.txt",
                Content: "short",
                Path:    "/path/to/test.txt",
            },
            expectError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := processor.validateDocument(tt.doc)
            if tt.expectError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

func TestSimpleDocumentProcessor_ValidateChunkingConfig(t *testing.T) {
    processor := NewSimpleDocumentProcessor()

    tests := []struct {
        name        string
        config      ChunkingConfig
        expectError bool
    }{
        {
            name: "valid config",
            config: ChunkingConfig{
                ChunkSize:        100,
                ChunkOverlap:     20,
                ChunkingStrategy: "fixed",
            },
            expectError: false,
        },
        {
            name: "zero chunk size",
            config: ChunkingConfig{
                ChunkSize:        0,
                ChunkOverlap:     20,
                ChunkingStrategy: "fixed",
            },
            expectError: true,
        },
        {
            name: "negative chunk size",
            config: ChunkingConfig{
                ChunkSize:        -100,
                ChunkOverlap:     20,
                ChunkingStrategy: "fixed",
            },
            expectError: true,
        },
        {
            name: "chunk size too small",
            config: ChunkingConfig{
                ChunkSize:        30,
                ChunkOverlap:     10,
                ChunkingStrategy: "fixed",
            },
            expectError: true,
        },
        {
            name: "negative overlap",
            config: ChunkingConfig{
                ChunkSize:        100,
                ChunkOverlap:     -10,
                ChunkingStrategy: "fixed",
            },
            expectError: true,
        },
        {
            name: "overlap greater than size",
            config: ChunkingConfig{
                ChunkSize:        100,
                ChunkOverlap:     150,
                ChunkingStrategy: "fixed",
            },
            expectError: true,
        },
        {
            name: "invalid strategy",
            config: ChunkingConfig{
                ChunkSize:        100,
                ChunkOverlap:     20,
                ChunkingStrategy: "invalid",
            },
            expectError: true,
        },
        {
            name: "empty strategy (should be valid)",
            config: ChunkingConfig{
                ChunkSize:        100,
                ChunkOverlap:     20,
                ChunkingStrategy: "",
            },
            expectError: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := processor.validateChunkingConfig(tt.config)
            if tt.expectError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

func TestSimpleDocumentProcessor_GetSupportedFileTypes(t *testing.T) {
    processor := NewSimpleDocumentProcessor()
    
    types := processor.GetSupportedFileTypes()
    
    assert.NotEmpty(t, types)
    assert.Contains(t, types, ".txt")
    assert.Contains(t, types, ".md")
    assert.Contains(t, types, ".pdf")
    assert.Contains(t, types, ".go")
}

func TestSimpleDocumentProcessor_GetProcessingStatistics(t *testing.T) {
    processor := NewSimpleDocumentProcessor()
    
    stats := processor.GetProcessingStatistics()
    
    assert.NotEmpty(t, stats)
    assert.Equal(t, "simple", stats["processor_type"])
    assert.NotZero(t, stats["supported_formats"])
    assert.NotNil(t, stats["last_updated"])
    assert.NotEmpty(t, stats["features"])
}

func TestSimpleDocumentProcessor_ProcessDocuments(t *testing.T) {
    processor := NewSimpleDocumentProcessor()

    docs := []*domain.Document{
        {
            ID:      "doc1",
            Name:    "test1.txt",
            Content: "This is the first document with sufficient content for chunking purposes",
            Path:    "/path/to/test1.txt",
        },
        {
            ID:      "doc2",
            Name:    "test2.txt",
            Content: "This is the second document with enough content to be processed properly",
            Path:    "/path/to/test2.txt",
        },
    }

    config := ChunkingConfig{
        ChunkSize:        50,
        ChunkOverlap:     10,
        ChunkingStrategy: "fixed",
    }

    result, err := processor.ProcessDocuments(docs, config)
    
    require.NoError(t, err)
    assert.Len(t, result, 2)
    assert.Contains(t, result, "doc1")
    assert.Contains(t, result, "doc2")
    assert.NotEmpty(t, result["doc1"])
    assert.NotEmpty(t, result["doc2"])
}

func TestSimpleDocumentProcessor_LoadDocuments_Integration(t *testing.T) {
    // Skip this test if we can't create temp files
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    processor := NewSimpleDocumentProcessor()

    // Create a temporary directory with test files
    tempDir, err := os.MkdirTemp("", "processor_test_")
    require.NoError(t, err)
    defer os.RemoveAll(tempDir)

    // Create test files
    testFiles := map[string]string{
        "test1.txt": "This is the content of the first test file with sufficient length for testing",
        "test2.md":  "# Test Markdown\n\nThis is a markdown file with enough content to be processed",
        "test3.log": "2023-01-01 10:00:00 INFO Test log entry with sufficient content for processing",
    }

    for filename, content := range testFiles {
        err := os.WriteFile(filepath.Join(tempDir, filename), []byte(content), 0644)
        require.NoError(t, err)
    }

    // Test loading documents
    options := DocumentLoaderOptions{
        ChunkSize:        100,
        ChunkOverlap:     20,
        ChunkingStrategy: "fixed",
        ProcessExts:      []string{".txt", ".md", ".log"},
    }

    docs, err := processor.LoadDocuments(tempDir, options)
    
    require.NoError(t, err)
    assert.Len(t, docs, 3)
    
    for _, doc := range docs {
        assert.NotEmpty(t, doc.ID)
        assert.NotEmpty(t, doc.Name)
        assert.NotEmpty(t, doc.Content)
        assert.Greater(t, len(doc.Content), 10)
    }
}