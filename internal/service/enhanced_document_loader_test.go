package service

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnhancedDocumentLoader_NewEnhancedDocumentLoader(t *testing.T) {
	loader := NewEnhancedDocumentLoader()
	
	assert.NotNil(t, loader)
	assert.NotNil(t, loader.legacyStrategy)
	assert.NotNil(t, loader.langchainStrategy)
	assert.NotNil(t, loader.telemetry)
	assert.NotEmpty(t, loader.strategy)
}

func TestEnhancedDocumentLoader_StrategySelection(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected string
	}{
		{
			name:     "default_hybrid",
			envValue: "",
			expected: "hybrid",
		},
		{
			name:     "explicit_langchain",
			envValue: "langchain",
			expected: "langchain",
		},
		{
			name:     "explicit_legacy",
			envValue: "legacy",
			expected: "legacy",
		},
		{
			name:     "explicit_hybrid",
			envValue: "hybrid",
			expected: "hybrid",
		},
		{
			name:     "invalid_fallback_to_hybrid",
			envValue: "invalid",
			expected: "hybrid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			originalEnv := os.Getenv("RLAMA_LOADER_STRATEGY")
			defer func() {
				if originalEnv == "" {
					os.Unsetenv("RLAMA_LOADER_STRATEGY")
				} else {
					os.Setenv("RLAMA_LOADER_STRATEGY", originalEnv)
				}
			}()

			if tt.envValue != "" {
				os.Setenv("RLAMA_LOADER_STRATEGY", tt.envValue)
			} else {
				os.Unsetenv("RLAMA_LOADER_STRATEGY")
			}

			loader := NewEnhancedDocumentLoader()
			assert.Equal(t, tt.expected, loader.GetStrategy())
		})
	}
}

func TestEnhancedDocumentLoader_SetStrategy(t *testing.T) {
	loader := NewEnhancedDocumentLoader()
	
	// Valid strategies
	validStrategies := []string{"langchain", "legacy", "hybrid"}
	for _, strategy := range validStrategies {
		loader.SetStrategy(strategy)
		assert.Equal(t, strategy, loader.GetStrategy())
	}
	
	// Invalid strategy should not change current strategy
	currentStrategy := loader.GetStrategy()
	loader.SetStrategy("invalid")
	assert.Equal(t, currentStrategy, loader.GetStrategy())
}

func TestEnhancedDocumentLoader_LoadDocuments_Integration(t *testing.T) {
	// Skip this test if running in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create test directory with sample files
	testDir, err := os.MkdirTemp("", "enhanced_loader_test_")
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	// Create test files
	testFiles := map[string]string{
		"simple.txt":    "This is a simple text file for testing enhanced document loading capabilities.",
		"readme.md":     "# Test Document\n\nThis is a markdown file with **bold** and *italic* text.\n\n## Section\n\nSome content here.",
		"config.json":   `{"name": "test", "version": "1.0", "description": "Test configuration file"}`,
		"sample.py":     "#!/usr/bin/env python3\n\ndef hello():\n    print('Hello from Python!')\n\nif __name__ == '__main__':\n    hello()",
		"data.csv":      "name,age,city\nJohn,30,New York\nJane,25,Los Angeles\nBob,35,Chicago",
		"notes.log":     "[2024-01-01 10:00:00] INFO: Application started\n[2024-01-01 10:01:00] DEBUG: Processing documents\n[2024-01-01 10:02:00] INFO: Processing complete",
	}

	for filename, content := range testFiles {
		err := os.WriteFile(filepath.Join(testDir, filename), []byte(content), 0644)
		require.NoError(t, err)
	}

	// Test different strategies
	strategies := []string{"legacy", "langchain", "hybrid"}
	
	for _, strategy := range strategies {
		t.Run("strategy_"+strategy, func(t *testing.T) {
			loader := NewEnhancedDocumentLoader()
			loader.SetStrategy(strategy)
			
			options := NewDocumentLoaderOptions()
			options.ProcessExts = []string{".txt", ".md", ".json", ".py", ".csv", ".log"}
			
			start := time.Now()
			docs, err := loader.LoadDocumentsFromFolderWithOptions(testDir, options)
			duration := time.Since(start)
			
			assert.NoError(t, err, "Strategy %s should not fail", strategy)
			assert.Len(t, docs, 6, "Should load all 6 test files")
			
			// Verify document content
			docMap := make(map[string]string)
			for _, doc := range docs {
				docMap[doc.Name] = doc.Content
			}
			
			assert.Contains(t, docMap, "simple.txt")
			assert.Contains(t, docMap, "readme.md")
			assert.Contains(t, docMap, "config.json")
			assert.Contains(t, docMap, "sample.py")
			assert.Contains(t, docMap, "data.csv")
			assert.Contains(t, docMap, "notes.log")
			
			// Verify content integrity
			assert.Contains(t, docMap["simple.txt"], "enhanced document loading")
			assert.Contains(t, docMap["readme.md"], "# Test Document")
			assert.Contains(t, docMap["config.json"], "test configuration")
			assert.Contains(t, docMap["sample.py"], "def hello()")
			assert.Contains(t, docMap["data.csv"], "John,30,New York")
			assert.Contains(t, docMap["notes.log"], "Application started")
			
			t.Logf("Strategy %s: loaded %d documents in %v", strategy, len(docs), duration)
		})
	}
}

func TestEnhancedDocumentLoader_TelemetryTracking(t *testing.T) {
	loader := NewEnhancedDocumentLoader()
	
	// Initial telemetry should be empty
	telemetry := loader.GetTelemetry()
	assert.Equal(t, 0, telemetry.LangChainSuccesses)
	assert.Equal(t, 0, telemetry.LangChainFailures)
	assert.Equal(t, 0, telemetry.LegacySuccesses)
	assert.Equal(t, 0, telemetry.LegacyFailures)
	
	// Create test directory
	testDir, err := os.MkdirTemp("", "telemetry_test_")
	require.NoError(t, err)
	defer os.RemoveAll(testDir)
	
	// Create a simple test file
	err = os.WriteFile(filepath.Join(testDir, "test.txt"), []byte("test content"), 0644)
	require.NoError(t, err)
	
	// Test LangChain strategy
	loader.SetStrategy("langchain")
	options := NewDocumentLoaderOptions()
	docs, err := loader.LoadDocumentsFromFolderWithOptions(testDir, options)
	
	if err == nil && len(docs) > 0 {
		// If LangChain succeeded
		telemetry = loader.GetTelemetry()
		assert.Equal(t, 1, telemetry.LangChainSuccesses)
		assert.Equal(t, 0, telemetry.LangChainFailures)
	} else {
		// If LangChain failed
		telemetry = loader.GetTelemetry()
		assert.Equal(t, 0, telemetry.LangChainSuccesses)
		assert.Equal(t, 1, telemetry.LangChainFailures)
	}
	
	// Test legacy strategy
	loader.SetStrategy("legacy")
	docs, err = loader.LoadDocumentsFromFolderWithOptions(testDir, options)
	
	if err == nil && len(docs) > 0 {
		// Legacy should succeed
		telemetry = loader.GetTelemetry()
		assert.Equal(t, 1, telemetry.LegacySuccesses)
		assert.Equal(t, 0, telemetry.LegacyFailures)
	}
}

func TestEnhancedDocumentLoader_GetSupportedFileTypes(t *testing.T) {
	loader := NewEnhancedDocumentLoader()
	
	supportedTypes := loader.GetSupportedFileTypes()
	
	// Should have a reasonable number of supported types
	assert.GreaterOrEqual(t, len(supportedTypes), 20)
	
	// Check for common file types
	expectedTypes := []string{".txt", ".md", ".pdf", ".json", ".csv", ".html", ".go", ".py", ".js"}
	for _, expectedType := range expectedTypes {
		assert.Contains(t, supportedTypes, expectedType, "Should support %s files", expectedType)
	}
}

func TestEnhancedDocumentLoader_GetAvailableStrategies(t *testing.T) {
	loader := NewEnhancedDocumentLoader()
	
	strategies := loader.GetAvailableStrategies()
	
	// Should have all three strategies
	assert.Len(t, strategies, 3)
	assert.Contains(t, strategies, "langchain")
	assert.Contains(t, strategies, "legacy")
	assert.Contains(t, strategies, "hybrid")
	
	// Each strategy should have required fields
	for name, info := range strategies {
		assert.NotEmpty(t, info.Name, "Strategy %s should have a name", name)
		assert.NotEmpty(t, info.Description, "Strategy %s should have a description", name)
		assert.NotEmpty(t, info.FileTypes, "Strategy %s should have file types", name)
		// Available field can be true or false
	}
}

func TestEnhancedDocumentLoader_ErrorHandling(t *testing.T) {
	loader := NewEnhancedDocumentLoader()
	
	// Test with non-existent directory
	docs, err := loader.LoadDocumentsFromFolderWithOptions("/non/existent/path", NewDocumentLoaderOptions())
	assert.Error(t, err)
	assert.Nil(t, docs)
	
	// Test with empty directory
	emptyDir, err := os.MkdirTemp("", "empty_test_")
	require.NoError(t, err)
	defer os.RemoveAll(emptyDir)
	
	docs, err = loader.LoadDocumentsFromFolderWithOptions(emptyDir, NewDocumentLoaderOptions())
	assert.Error(t, err)
	assert.Nil(t, docs)
}

func TestEnhancedDocumentLoader_FilteringOptions(t *testing.T) {
	// Create test directory with various files
	testDir, err := os.MkdirTemp("", "filtering_test_")
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	// Create subdirectories
	nodeModulesDir := filepath.Join(testDir, "node_modules")
	err = os.MkdirAll(nodeModulesDir, 0755)
	require.NoError(t, err)

	// Create test files
	testFiles := map[string]string{
		"include.txt":                    "This should be included",
		"include.md":                     "# This should be included",
		"exclude.log":                    "This should be excluded",
		"script.js":                      "console.log('hello');",
		"node_modules/package.json":      `{"name": "test"}`,
		filepath.Join("src", "main.go"): "package main",
	}

	// Create src directory
	err = os.MkdirAll(filepath.Join(testDir, "src"), 0755)
	require.NoError(t, err)

	for filePath, content := range testFiles {
		fullPath := filepath.Join(testDir, filePath)
		err := os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(t, err)
	}

	loader := NewEnhancedDocumentLoader()
	loader.SetStrategy("hybrid") // Use hybrid for reliability

	// Test exclude directories
	t.Run("exclude_directories", func(t *testing.T) {
		options := NewDocumentLoaderOptions()
		options.ExcludeDirs = []string{"node_modules"}
		
		docs, err := loader.LoadDocumentsFromFolderWithOptions(testDir, options)
		require.NoError(t, err)
		
		// Should not include files from node_modules
		for _, doc := range docs {
			assert.NotContains(t, doc.Path, "node_modules")
		}
	})

	// Test exclude extensions
	t.Run("exclude_extensions", func(t *testing.T) {
		options := NewDocumentLoaderOptions()
		options.ExcludeExts = []string{".log", ".js"}
		
		docs, err := loader.LoadDocumentsFromFolderWithOptions(testDir, options)
		require.NoError(t, err)
		
		// Should not include .log or .js files
		for _, doc := range docs {
			ext := filepath.Ext(doc.Path)
			assert.NotEqual(t, ".log", ext)
			assert.NotEqual(t, ".js", ext)
		}
	})

	// Test process specific extensions
	t.Run("process_specific_extensions", func(t *testing.T) {
		options := NewDocumentLoaderOptions()
		options.ProcessExts = []string{".txt", ".md"}
		
		docs, err := loader.LoadDocumentsFromFolderWithOptions(testDir, options)
		require.NoError(t, err)
		
		// Should only include .txt and .md files
		for _, doc := range docs {
			ext := filepath.Ext(doc.Path)
			assert.True(t, ext == ".txt" || ext == ".md", "Should only process .txt and .md files, got %s", ext)
		}
	})
}

func BenchmarkEnhancedDocumentLoader_Strategies(b *testing.B) {
	// Create test directory with files
	testDir, err := os.MkdirTemp("", "benchmark_test_")
	require.NoError(b, err)
	defer os.RemoveAll(testDir)

	// Create multiple test files
	for i := 0; i < 10; i++ {
		content := "This is test file number " + string(rune(i+'0')) + " with some content for benchmarking purposes. " +
			"It contains multiple sentences to simulate real document content. " +
			"The content is long enough to provide meaningful processing time measurements."
		
		filename := filepath.Join(testDir, "test_"+string(rune(i+'0'))+".txt")
		err := os.WriteFile(filename, []byte(content), 0644)
		require.NoError(b, err)
	}

	strategies := []string{"legacy", "langchain", "hybrid"}
	
	for _, strategy := range strategies {
		b.Run("strategy_"+strategy, func(b *testing.B) {
			loader := NewEnhancedDocumentLoader()
			loader.SetStrategy(strategy)
			options := NewDocumentLoaderOptions()
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				docs, err := loader.LoadDocumentsFromFolderWithOptions(testDir, options)
				if err != nil {
					b.Fatalf("Benchmark failed for strategy %s: %v", strategy, err)
				}
				if len(docs) == 0 {
					b.Fatalf("No documents loaded for strategy %s", strategy)
				}
			}
		})
	}
}