package client

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestPureGoTokenizer_Basic(t *testing.T) {
	// Path to BGE tokenizer.json
	tokenizerPath := filepath.Join("..", "..", "models", "bge-reranker-large-onnx", "tokenizer.json")
	
	tokenizer, err := NewPureGoTokenizer(tokenizerPath)
	if err != nil {
		t.Fatalf("Failed to create tokenizer: %v", err)
	}
	
	// Test basic encoding
	text := "Hello world"
	maxLength := 512
	
	tokenIDs, attentionMask, tokenTypeIDs := tokenizer.Encode(text, maxLength)
	
	// Verify output dimensions
	if len(tokenIDs) != maxLength {
		t.Errorf("Expected token IDs length %d, got %d", maxLength, len(tokenIDs))
	}
	
	if len(attentionMask) != maxLength {
		t.Errorf("Expected attention mask length %d, got %d", maxLength, len(attentionMask))
	}
	
	if len(tokenTypeIDs) != maxLength {
		t.Errorf("Expected token type IDs length %d, got %d", maxLength, len(tokenTypeIDs))
	}
	
	// Verify BOS token is first
	if tokenIDs[0] != tokenizer.bosTokenID {
		t.Errorf("Expected first token to be BOS (%d), got %d", tokenizer.bosTokenID, tokenIDs[0])
	}
	
	// Count non-padding tokens
	nonPaddingCount := 0
	for i, mask := range attentionMask {
		if mask == 1 {
			nonPaddingCount++
		} else if tokenIDs[i] != tokenizer.padTokenID {
			t.Errorf("Token at position %d should be padding token (%d), got %d", i, tokenizer.padTokenID, tokenIDs[i])
		}
	}
	
	t.Logf("Encoded text: %s", text)
	t.Logf("Non-padding tokens: %d", nonPaddingCount)
	t.Logf("First 10 token IDs: %v", tokenIDs[:10])
}

func TestPureGoTokenizer_QueryPassagePair(t *testing.T) {
	tokenizerPath := filepath.Join("..", "..", "models", "bge-reranker-large-onnx", "tokenizer.json")
	
	tokenizer, err := NewPureGoTokenizer(tokenizerPath)
	if err != nil {
		t.Fatalf("Failed to create tokenizer: %v", err)
	}
	
	query := "What is machine learning?"
	passage := "Machine learning is a subset of artificial intelligence that focuses on algorithms."
	maxLength := 512
	
	tokenIDs, attentionMask, _ := tokenizer.EncodeQueryPassagePair(query, passage, maxLength)
	
	// Verify output dimensions
	if len(tokenIDs) != maxLength {
		t.Errorf("Expected token IDs length %d, got %d", maxLength, len(tokenIDs))
	}
	
	// Count non-padding tokens
	nonPaddingCount := 0
	for _, mask := range attentionMask {
		if mask == 1 {
			nonPaddingCount++
		}
	}
	
	t.Logf("Query: %s", query)
	t.Logf("Passage: %s", passage)
	t.Logf("Non-padding tokens: %d", nonPaddingCount)
	t.Logf("First 15 token IDs: %v", tokenIDs[:15])
}

func TestPureGoTokenizer_SpecialTokens(t *testing.T) {
	tokenizerPath := filepath.Join("..", "..", "models", "bge-reranker-large-onnx", "tokenizer.json")
	
	tokenizer, err := NewPureGoTokenizer(tokenizerPath)
	if err != nil {
		t.Fatalf("Failed to create tokenizer: %v", err)
	}
	
	// Test that special tokens are properly identified
	expectedSpecialTokens := []string{"<s>", "</s>", "<pad>", "<unk>", "<mask>"}
	
	for _, token := range expectedSpecialTokens {
		if id, exists := tokenizer.specialTokens[token]; exists {
			t.Logf("Special token %s has ID: %d", token, id)
		} else {
			t.Logf("Warning: Special token %s not found", token)
		}
	}
	
	// Verify special token IDs are set
	t.Logf("BOS token ID: %d", tokenizer.bosTokenID)
	t.Logf("EOS token ID: %d", tokenizer.eosTokenID)
	t.Logf("PAD token ID: %d", tokenizer.padTokenID)
	t.Logf("UNK token ID: %d", tokenizer.unkTokenID)
	t.Logf("MASK token ID: %d", tokenizer.maskTokenID)
}

func TestPureGoTokenizer_VocabStats(t *testing.T) {
	tokenizerPath := filepath.Join("..", "..", "models", "bge-reranker-large-onnx", "tokenizer.json")
	
	tokenizer, err := NewPureGoTokenizer(tokenizerPath)
	if err != nil {
		t.Fatalf("Failed to create tokenizer: %v", err)
	}
	
	t.Logf("Vocabulary size: %d", len(tokenizer.vocab))
	t.Logf("Special tokens count: %d", len(tokenizer.specialTokens))
	t.Logf("Merge rules count: %d", len(tokenizer.merges))
	
	// Test some common tokens
	commonTokens := []string{"▁the", "▁and", "▁a", "▁to", "▁of"}
	
	for _, token := range commonTokens {
		if id, exists := tokenizer.vocab[token]; exists {
			t.Logf("Common token '%s' has ID: %d", token, id)
		} else {
			t.Logf("Common token '%s' not found in vocabulary", token)
		}
	}
}

func TestPureGoTokenizer_BPEEncoding(t *testing.T) {
	tokenizerPath := filepath.Join("..", "..", "models", "bge-reranker-large-onnx", "tokenizer.json")
	
	tokenizer, err := NewPureGoTokenizer(tokenizerPath)
	if err != nil {
		t.Fatalf("Failed to create tokenizer: %v", err)
	}
	
	// Test BPE encoding of known and unknown words
	testWords := []string{"▁hello", "▁world", "▁unknownword123", "▁test"}
	
	for _, word := range testWords {
		subTokens := tokenizer.bpeEncode(word)
		t.Logf("BPE encoding of '%s': %v", word, subTokens)
		
		// Verify all sub-tokens have IDs
		for _, subToken := range subTokens {
			tokenID := tokenizer.getTokenID(subToken)
			if tokenID == tokenizer.unkTokenID {
				t.Logf("  Sub-token '%s' -> UNK (%d)", subToken, tokenID)
			} else {
				t.Logf("  Sub-token '%s' -> %d", subToken, tokenID)
			}
		}
	}
}

func TestPureGoTokenizer_EdgeCases(t *testing.T) {
	tokenizerPath := filepath.Join("..", "..", "models", "bge-reranker-large-onnx", "tokenizer.json")
	
	tokenizer, err := NewPureGoTokenizer(tokenizerPath)
	if err != nil {
		t.Fatalf("Failed to create tokenizer: %v", err)
	}
	
	// Test edge cases
	testCases := []struct {
		name string
		text string
	}{
		{"Empty string", ""},
		{"Single character", "a"},
		{"Only spaces", "   "},
		{"Special characters", "!@#$%^&*()"},
		{"Numbers", "12345"},
		{"Mixed case", "Hello WORLD"},
		{"Very long text", strings.Repeat("This is a test sentence. ", 50)},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tokenIDs, attentionMask, _ := tokenizer.Encode(tc.text, 128)
			
			if len(tokenIDs) != 128 {
				t.Errorf("Expected 128 tokens, got %d", len(tokenIDs))
			}
			
			if len(attentionMask) != 128 {
				t.Errorf("Expected 128 attention mask values, got %d", len(attentionMask))
			}
			
			// Count actual tokens (non-padding)
			actualTokens := 0
			for _, mask := range attentionMask {
				if mask == 1 {
					actualTokens++
				}
			}
			
			t.Logf("Text: '%s' -> %d actual tokens", tc.text, actualTokens)
		})
	}
}