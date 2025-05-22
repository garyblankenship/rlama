package client

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestPureGoTokenizer_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func() (string, func()) // returns path and cleanup
		expectError bool
		errorMsg    string
	}{
		{
			name: "NonexistentFile",
			setupFunc: func() (string, func()) {
				return "/nonexistent/path/tokenizer.json", func() {}
			},
			expectError: true,
			errorMsg:    "failed to read tokenizer.json",
		},
		{
			name: "InvalidJSON",
			setupFunc: func() (string, func()) {
				tmpFile, err := os.CreateTemp("", "invalid_tokenizer_*.json")
				if err != nil {
					t.Fatal(err)
				}
				tmpFile.WriteString("{ invalid json content")
				tmpFile.Close()
				return tmpFile.Name(), func() { os.Remove(tmpFile.Name()) }
			},
			expectError: true,
			errorMsg:    "failed to parse tokenizer.json",
		},
		{
			name: "EmptyFile",
			setupFunc: func() (string, func()) {
				tmpFile, err := os.CreateTemp("", "empty_tokenizer_*.json")
				if err != nil {
					t.Fatal(err)
				}
				tmpFile.Close()
				return tmpFile.Name(), func() { os.Remove(tmpFile.Name()) }
			},
			expectError: true,
			errorMsg:    "failed to parse tokenizer.json",
		},
		{
			name: "MissingRequiredFields",
			setupFunc: func() (string, func()) {
				tmpFile, err := os.CreateTemp("", "incomplete_tokenizer_*.json")
				if err != nil {
					t.Fatal(err)
				}
				
				incompleteConfig := map[string]interface{}{
					"version": "1.0",
					// Missing other required fields
				}
				
				encoder := json.NewEncoder(tmpFile)
				encoder.Encode(incompleteConfig)
				tmpFile.Close()
				
				return tmpFile.Name(), func() { os.Remove(tmpFile.Name()) }
			},
			expectError: false, // Should not error, but will have empty structures
		},
		{
			name: "PermissionDenied",
			setupFunc: func() (string, func()) {
				tmpFile, err := os.CreateTemp("", "permission_test_*.json")
				if err != nil {
					t.Fatal(err)
				}
				tmpFile.WriteString("{}")
				tmpFile.Close()
				
				// Remove read permissions
				os.Chmod(tmpFile.Name(), 0000)
				
				return tmpFile.Name(), func() {
					os.Chmod(tmpFile.Name(), 0644)
					os.Remove(tmpFile.Name())
				}
			},
			expectError: true,
			errorMsg:    "failed to read tokenizer.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenizerPath, cleanup := tt.setupFunc()
			defer cleanup()

			tokenizer, err := NewPureGoTokenizer(tokenizerPath)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain '%s', got: %s", tt.errorMsg, err.Error())
				}
				if tokenizer != nil {
					t.Errorf("Expected nil tokenizer on error, got non-nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestPureGoTokenizer_MalformedInput(t *testing.T) {
	// Use a minimal valid tokenizer config for testing
	tokenizerPath := createMinimalTokenizerConfig(t)
	defer os.Remove(tokenizerPath)

	tokenizer, err := NewPureGoTokenizer(tokenizerPath)
	if err != nil {
		t.Fatalf("Failed to create tokenizer: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		maxLen   int
		shouldPanic bool
	}{
		{
			name:   "InvalidUTF8",
			input:  string([]byte{0xff, 0xfe, 0xfd}),
			maxLen: 128,
		},
		{
			name:   "VeryLongString",
			input:  strings.Repeat("This is a very long test string. ", 1000),
			maxLen: 512,
		},
		{
			name:   "EmptyString",
			input:  "",
			maxLen: 128,
		},
		{
			name:   "OnlySpaces",
			input:  "   \t\n   ",
			maxLen: 128,
		},
		{
			name:   "SpecialCharacters",
			input:  "🚀🔥💯🎉\u0000\u001f\u007f",
			maxLen: 128,
		},
		{
			name:   "NullBytes",
			input:  "text\x00with\x00nulls",
			maxLen: 128,
		},
		{
			name:   "ZeroMaxLength",
			input:  "test",
			maxLen: 0,
		},
		{
			name:   "NegativeMaxLength",
			input:  "test",
			maxLen: -1,
		},
		{
			name:   "ExtremelyLargeMaxLength",
			input:  "test",
			maxLen: 1000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.shouldPanic {
						t.Errorf("Unexpected panic: %v", r)
					}
				}
			}()

			tokenIDs, attentionMask, tokenTypeIDs := tokenizer.Encode(tt.input, tt.maxLen)

			// Basic sanity checks for non-panicking cases
			if !tt.shouldPanic {
				expectedLen := tt.maxLen
				if tt.maxLen <= 0 {
					expectedLen = 0
				}
				
				if len(tokenIDs) != expectedLen {
					t.Errorf("Expected %d token IDs, got %d", expectedLen, len(tokenIDs))
				}
				if len(attentionMask) != expectedLen {
					t.Errorf("Expected %d attention mask values, got %d", expectedLen, len(attentionMask))
				}
				if len(tokenTypeIDs) != expectedLen {
					t.Errorf("Expected %d token type IDs, got %d", expectedLen, len(tokenTypeIDs))
				}

				t.Logf("Input: %q -> %d tokens", tt.input, len(tokenIDs))
			}
		})
	}
}

func TestPureGoTokenizer_BPEEdgeCases(t *testing.T) {
	tokenizerPath := createMinimalTokenizerConfig(t)
	defer os.Remove(tokenizerPath)

	tokenizer, err := NewPureGoTokenizer(tokenizerPath)
	if err != nil {
		t.Fatalf("Failed to create tokenizer: %v", err)
	}

	tests := []struct {
		name  string
		token string
	}{
		{"EmptyToken", ""},
		{"SingleChar", "a"},
		{"UnknownToken", "xyz123unknown"},
		{"SpecialCharsOnly", "!@#$%"},
		{"NumbersOnly", "123456"},
		{"MixedCase", "AbCdEf"},
		{"UnicodeEmoji", "🔥💯"},
		{"LongToken", strings.Repeat("verylongtoken", 10)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tokenizer.bpeEncode(tt.token)
			
			// Should always return some result, even if it's fallback
			if result == nil {
				t.Errorf("BPE encode returned nil for token: %s", tt.token)
			}
			
			t.Logf("Token: %q -> BPE: %v", tt.token, result)
		})
	}
}

func TestPureGoTokenizer_VocabularyEdgeCases(t *testing.T) {
	tokenizerPath := createMinimalTokenizerConfig(t)
	defer os.Remove(tokenizerPath)

	tokenizer, err := NewPureGoTokenizer(tokenizerPath)
	if err != nil {
		t.Fatalf("Failed to create tokenizer: %v", err)
	}

	tests := []struct {
		name  string
		token string
	}{
		{"EmptyString", ""},
		{"Whitespace", " "},
		{"Tab", "\t"},
		{"Newline", "\n"},
		{"MetaspacePrefix", "▁"},
		{"SpecialTokenLike", "<special>"},
		{"VeryLongToken", strings.Repeat("a", 1000)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenID := tokenizer.getTokenID(tt.token)
			
			// Should always return a valid ID (possibly UNK)
			if tokenID < 0 {
				t.Errorf("getTokenID returned negative ID for token: %s", tt.token)
			}
			
			t.Logf("Token: %q -> ID: %d", tt.token, tokenID)
		})
	}
}

// Helper function to create a minimal valid tokenizer config for testing
func createMinimalTokenizerConfig(t *testing.T) string {
	tmpFile, err := os.CreateTemp("", "minimal_tokenizer_*.json")
	if err != nil {
		t.Fatal(err)
	}

	minimalConfig := TokenizerConfig{
		Version: "1.0",
		AddedTokens: []AddedToken{
			{ID: 0, Content: "<s>", Special: true},
			{ID: 1, Content: "<pad>", Special: true},
			{ID: 2, Content: "</s>", Special: true},
			{ID: 3, Content: "<unk>", Special: true},
			{ID: 4, Content: "<mask>", Special: true},
		},
		Model: Model{
			Type:  "Unigram",
			UnkID: 3,
			Vocab: [][]interface{}{
				{"<s>", 0.0},
				{"<pad>", 0.0},
				{"</s>", 0.0},
				{"<unk>", 0.0},
				{"<mask>", 0.0},
				{"▁a", -1.0},
				{"▁the", -1.0},
				{"▁and", -1.0},
				{"a", -2.0},
				{"e", -2.0},
				{"i", -2.0},
				{"o", -2.0},
				{"u", -2.0},
			},
		},
		PreTokenizer: PreTokenizer{
			Type:          "Metaspace",
			Replacement:   "▁",
			PrependScheme: "always",
		},
	}

	encoder := json.NewEncoder(tmpFile)
	if err := encoder.Encode(minimalConfig); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		t.Fatal(err)
	}

	tmpFile.Close()
	return tmpFile.Name()
}