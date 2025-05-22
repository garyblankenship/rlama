package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dontizi/rlama/internal/domain"
)

func TestCreateTempDirForDocuments(t *testing.T) {
	// Create test documents
	documents := []*domain.Document{
		{
			ID:      "doc1",
			Content: "This is test document 1",
			Path:    "/test/path1",
		},
		{
			ID:      "doc2", 
			Content: "This is test document 2",
			Path:    "/test/path2",
		},
	}

	// Create temporary directory
	tempDir := CreateTempDirForDocuments(documents)
	if tempDir == "" {
		t.Fatal("CreateTempDirForDocuments returned empty string")
	}

	// Verify directory exists
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Fatalf("Temporary directory %s does not exist", tempDir)
	}

	// Verify files were created
	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatalf("Error reading temp directory: %v", err)
	}

	if len(files) != len(documents) {
		t.Fatalf("Expected %d files, got %d", len(documents), len(files))
	}

	// Verify file contents
	for i, file := range files {
		filePath := filepath.Join(tempDir, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Error reading file %s: %v", filePath, err)
		}

		expectedContent := documents[i].Content
		if string(content) != expectedContent {
			t.Errorf("File content mismatch. Expected: %s, Got: %s", expectedContent, string(content))
		}
	}

	// Clean up
	CleanupTempDir(tempDir)

	// Verify directory was cleaned up
	if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
		t.Fatalf("Temporary directory %s was not cleaned up", tempDir)
	}
}

func TestCleanupTempDirWithEmpty(t *testing.T) {
	// Test cleanup with empty path (should not panic)
	CleanupTempDir("")
}