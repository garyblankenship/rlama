package client

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// TokenizerConfig represents the tokenizer.json configuration
type TokenizerConfig struct {
	Version       string      `json:"version"`
	AddedTokens   []AddedToken `json:"added_tokens"`
	Normalizer    Normalizer  `json:"normalizer"`
	PreTokenizer  PreTokenizer `json:"pre_tokenizer"`
	PostProcessor PostProcessor `json:"post_processor"`
	Decoder       Decoder     `json:"decoder"`
	Model         Model       `json:"model"`
}

type AddedToken struct {
	ID         int    `json:"id"`
	Content    string `json:"content"`
	SingleWord bool   `json:"single_word"`
	LStrip     bool   `json:"lstrip"`
	RStrip     bool   `json:"rstrip"`
	Normalized bool   `json:"normalized"`
	Special    bool   `json:"special"`
}

type Normalizer struct {
	Type        string      `json:"type"`
	Normalizers []Normalizer `json:"normalizers,omitempty"`
	// Add specific fields as needed
}

type PreTokenizer struct {
	Type            string `json:"type"`
	Replacement     string `json:"replacement,omitempty"`
	PrependScheme   string `json:"prepend_scheme,omitempty"`
	Split           bool   `json:"split,omitempty"`
}

type PostProcessor struct {
	Type   string                 `json:"type"`
	Single []PostProcessorStep    `json:"single,omitempty"`
	Pair   []PostProcessorStep    `json:"pair,omitempty"`
}

type PostProcessorStep struct {
	SpecialToken *SpecialTokenStep `json:"SpecialToken,omitempty"`
	Sequence     *SequenceStep     `json:"Sequence,omitempty"`
}

type SpecialTokenStep struct {
	ID     string `json:"id"`
	TypeID int    `json:"type_id"`
}

type SequenceStep struct {
	ID     string `json:"id"`
	TypeID int    `json:"type_id"`
}

type Decoder struct {
	Type string `json:"type"`
}

type Model struct {
	Type   string        `json:"type"`
	UnkID  int          `json:"unk_id,omitempty"`
	Vocab  [][]interface{} `json:"vocab,omitempty"`  // For Unigram: [[token, score], ...]
	Merges []string      `json:"merges,omitempty"`   // For BPE
}

// PureGoTokenizer implements XLM-RoBERTa tokenization in pure Go
type PureGoTokenizer struct {
	config        *TokenizerConfig
	vocab         map[string]int64
	idToToken     map[int64]string
	specialTokens map[string]int64
	
	// XLM-RoBERTa specific
	bosTokenID int64
	eosTokenID int64
	padTokenID int64
	unkTokenID int64
	maskTokenID int64
	
	// Merges for BPE
	merges map[string]int
	
	// Regex patterns
	preTokenizePattern *regexp.Regexp
}

// NewPureGoTokenizer creates a new pure Go tokenizer from tokenizer.json
func NewPureGoTokenizer(tokenizerJSONPath string) (*PureGoTokenizer, error) {
	// Read tokenizer.json
	data, err := os.ReadFile(tokenizerJSONPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read tokenizer.json: %w", err)
	}
	
	var config TokenizerConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse tokenizer.json: %w", err)
	}
	
	tokenizer := &PureGoTokenizer{
		config:        &config,
		vocab:         make(map[string]int64),
		idToToken:     make(map[int64]string),
		specialTokens: make(map[string]int64),
		merges:        make(map[string]int),
	}
	
	// Load vocabulary and special tokens
	if err := tokenizer.loadVocabAndSpecialTokens(); err != nil {
		return nil, fmt.Errorf("failed to load vocabulary: %w", err)
	}
	
	// Load BPE merges
	if err := tokenizer.loadMerges(); err != nil {
		return nil, fmt.Errorf("failed to load merges: %w", err)
	}
	
	// Compile pre-tokenization regex
	tokenizer.preTokenizePattern = regexp.MustCompile(`\S+`)
	
	return tokenizer, nil
}

// loadVocabAndSpecialTokens loads the vocabulary and identifies special tokens
func (t *PureGoTokenizer) loadVocabAndSpecialTokens() error {
	// Load special tokens from added_tokens
	for _, token := range t.config.AddedTokens {
		tokenID := int64(token.ID)
		t.specialTokens[token.Content] = tokenID
		t.idToToken[tokenID] = token.Content
		
		// Identify standard special tokens
		switch token.Content {
		case "<s>":
			t.bosTokenID = tokenID
		case "</s>":
			t.eosTokenID = tokenID
		case "<pad>":
			t.padTokenID = tokenID
		case "<unk>":
			t.unkTokenID = tokenID
		case "<mask>":
			t.maskTokenID = tokenID
		}
	}
	
	// Load regular vocabulary from model
	if t.config.Model.Vocab != nil {
		for i, vocabEntry := range t.config.Model.Vocab {
			if len(vocabEntry) >= 2 {
				if token, ok := vocabEntry[0].(string); ok {
					tokenID := int64(i)
					t.vocab[token] = tokenID
					t.idToToken[tokenID] = token
				}
			}
		}
	}
	
	return nil
}

// loadMerges loads BPE merge rules
func (t *PureGoTokenizer) loadMerges() error {
	if t.config.Model.Merges == nil {
		return nil
	}
	
	for i, merge := range t.config.Model.Merges {
		t.merges[merge] = i
	}
	
	return nil
}

// normalize applies text normalization
func (t *PureGoTokenizer) normalize(text string) string {
	// Basic normalization - this is a simplified version
	// The real implementation would handle the complex normalizer chain
	
	// Replace multiple spaces with single space
	text = regexp.MustCompile(` {2,}`).ReplaceAllString(text, " ")
	
	// Trim spaces
	text = strings.TrimSpace(text)
	
	return text
}

// preTokenize splits text into preliminary tokens
func (t *PureGoTokenizer) preTokenize(text string) []string {
	// XLM-RoBERTa uses Metaspace pre-tokenization
	// This adds a special character (▁) at word boundaries
	
	// Split on whitespace first
	words := strings.Fields(text)
	var tokens []string
	
	for i, word := range words {
		if i == 0 {
			// First word gets prepended with ▁ (always scheme)
			tokens = append(tokens, "▁"+word)
		} else {
			tokens = append(tokens, "▁"+word)
		}
	}
	
	return tokens
}

// getTokenID returns the token ID for a given token
func (t *PureGoTokenizer) getTokenID(token string) int64 {
	// Check special tokens first
	if id, exists := t.specialTokens[token]; exists {
		return id
	}
	
	// Check regular vocabulary
	if id, exists := t.vocab[token]; exists {
		return id
	}
	
	// Return unknown token ID
	return t.unkTokenID
}

// Encode tokenizes a text string and returns token IDs, attention mask, and token type IDs
func (t *PureGoTokenizer) Encode(text string, maxLength int) ([]int64, []int64, []int64) {
	// Step 1: Normalize text
	normalizedText := t.normalize(text)
	
	// Step 2: Pre-tokenize
	tokens := t.preTokenize(normalizedText)
	
	// Step 3: Apply BPE (simplified - this is complex to implement fully)
	// For now, we'll use existing tokens and break unknown ones into characters
	var tokenIDs []int64
	
	// Add BOS token
	tokenIDs = append(tokenIDs, t.bosTokenID)
	
	// Process each token
	for _, token := range tokens {
		// Try to find the token directly
		if id, exists := t.vocab[token]; exists {
			tokenIDs = append(tokenIDs, id)
		} else {
			// For unknown tokens, try to break them down
			// This is a simplified approach - real BPE is more complex
			subTokens := t.bpeEncode(token)
			for _, subToken := range subTokens {
				tokenIDs = append(tokenIDs, t.getTokenID(subToken))
			}
		}
	}
	
	// Add EOS token
	tokenIDs = append(tokenIDs, t.eosTokenID)
	
	// Step 4: Handle edge cases and truncate or pad to max length
	if maxLength <= 0 {
		// Return empty arrays for invalid max length
		return []int64{}, []int64{}, []int64{}
	}
	
	if len(tokenIDs) > maxLength {
		tokenIDs = tokenIDs[:maxLength]
		// Ensure we end with EOS if truncated
		if len(tokenIDs) > 0 {
			tokenIDs[len(tokenIDs)-1] = t.eosTokenID
		}
	}
	
	// Create attention mask (1 for real tokens, 0 for padding)
	attentionMask := make([]int64, maxLength)
	tokenTypeIDs := make([]int64, maxLength)
	
	for i := 0; i < len(tokenIDs) && i < maxLength; i++ {
		attentionMask[i] = 1
		tokenTypeIDs[i] = 0 // All tokens are type 0 for single sequence
	}
	
	// Pad with pad tokens
	for len(tokenIDs) < maxLength {
		tokenIDs = append(tokenIDs, t.padTokenID)
	}
	
	return tokenIDs, attentionMask, tokenTypeIDs
}

// bpeEncode applies BPE encoding to a token using merge rules
func (t *PureGoTokenizer) bpeEncode(token string) []string {
	// If token exists in vocabulary, return it
	if _, exists := t.vocab[token]; exists {
		return []string{token}
	}
	
	// Initialize with individual characters
	word := make([]string, 0, len(token))
	for _, char := range token {
		word = append(word, string(char))
	}
	
	// If single character, return as is
	if len(word) <= 1 {
		return word
	}
	
	// Apply BPE merge rules iteratively
	for {
		pairs := t.getPairs(word)
		if len(pairs) == 0 {
			break
		}
		
		// Find the merge with highest priority (lowest rank)
		bestPair := ""
		bestRank := len(t.merges) + 1
		
		for pair := range pairs {
			if rank, exists := t.merges[pair]; exists && rank < bestRank {
				bestPair = pair
				bestRank = rank
			}
		}
		
		if bestPair == "" {
			break
		}
		
		// Apply the merge
		word = t.applyMerge(word, bestPair)
	}
	
	// Filter out tokens that don't exist in vocabulary
	var result []string
	for _, subToken := range word {
		if _, exists := t.vocab[subToken]; exists {
			result = append(result, subToken)
		} else {
			// Break down further or use unknown token
			for _, char := range subToken {
				charStr := string(char)
				if _, exists := t.vocab[charStr]; exists {
					result = append(result, charStr)
				} else {
					result = append(result, "<unk>")
				}
			}
		}
	}
	
	return result
}

// getPairs gets all adjacent pairs in the word
func (t *PureGoTokenizer) getPairs(word []string) map[string]bool {
	pairs := make(map[string]bool)
	
	for i := 0; i < len(word)-1; i++ {
		pair := word[i] + " " + word[i+1]
		pairs[pair] = true
	}
	
	return pairs
}

// applyMerge applies a single merge rule to the word
func (t *PureGoTokenizer) applyMerge(word []string, merge string) []string {
	parts := strings.Split(merge, " ")
	if len(parts) != 2 {
		return word
	}
	
	first, second := parts[0], parts[1]
	newWord := make([]string, 0, len(word))
	
	i := 0
	for i < len(word) {
		if i < len(word)-1 && word[i] == first && word[i+1] == second {
			// Apply merge
			newWord = append(newWord, first+second)
			i += 2
		} else {
			newWord = append(newWord, word[i])
			i++
		}
	}
	
	return newWord
}

// EncodeQueryPassagePair encodes a query-passage pair for BGE reranker
func (t *PureGoTokenizer) EncodeQueryPassagePair(query, passage string, maxLength int) ([]int64, []int64, []int64) {
	// BGE format: query + " </s> " + passage
	combined := query + " </s> " + passage
	return t.Encode(combined, maxLength)
}