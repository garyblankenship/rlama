package client

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestPureGoTokenizer_ConcurrentAccess(t *testing.T) {
	tokenizerPath := createMinimalTokenizerConfig(t)
	defer func() {
		// Small delay to ensure file is not in use
		time.Sleep(10 * time.Millisecond)
		removeFile(tokenizerPath)
	}()

	tokenizer, err := NewPureGoTokenizer(tokenizerPath)
	if err != nil {
		t.Fatalf("Failed to create tokenizer: %v", err)
	}

	const numGoroutines = 50
	const numIterations = 100

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	// Test concurrent tokenization
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < numIterations; j++ {
				text := "Concurrent test message from goroutine"
				maxLength := 128

				tokenIDs, attentionMask, tokenTypeIDs := tokenizer.Encode(text, maxLength)

				// Validate results
				if len(tokenIDs) != maxLength {
					errors <- &ConcurrentTestError{
						GoroutineID: id,
						Iteration:   j,
						Message:     "Invalid token IDs length",
					}
					return
				}

				if len(attentionMask) != maxLength {
					errors <- &ConcurrentTestError{
						GoroutineID: id,
						Iteration:   j,
						Message:     "Invalid attention mask length",
					}
					return
				}

				if len(tokenTypeIDs) != maxLength {
					errors <- &ConcurrentTestError{
						GoroutineID: id,
						Iteration:   j,
						Message:     "Invalid token type IDs length",
					}
					return
				}

				// Verify data integrity
				if tokenIDs[0] != tokenizer.bosTokenID {
					errors <- &ConcurrentTestError{
						GoroutineID: id,
						Iteration:   j,
						Message:     "BOS token not found at beginning",
					}
					return
				}
			}
		}(i)
	}

	// Wait for completion or timeout
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Check for errors
		close(errors)
		for err := range errors {
			t.Error(err)
		}
	case <-time.After(30 * time.Second):
		t.Fatal("Concurrent test timed out")
	}
}

func TestPureGoTokenizer_RaceConditionDetection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping race condition test in short mode")
	}

	tokenizerPath := createMinimalTokenizerConfig(t)
	defer removeFile(tokenizerPath)

	tokenizer, err := NewPureGoTokenizer(tokenizerPath)
	if err != nil {
		t.Fatalf("Failed to create tokenizer: %v", err)
	}

	const numWorkers = 20
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	
	// Different operations to stress test
	operations := []func(){
		// Vocabulary access
		func() {
			for i := 0; i < 100; i++ {
				_ = tokenizer.getTokenID("test")
				_ = tokenizer.getTokenID("▁another")
				_ = tokenizer.getTokenID("<unk>")
			}
		},
		// BPE encoding
		func() {
			for i := 0; i < 100; i++ {
				_ = tokenizer.bpeEncode("testtoken")
				_ = tokenizer.bpeEncode("▁another")
				_ = tokenizer.bpeEncode("unknown123")
			}
		},
		// Full encoding
		func() {
			for i := 0; i < 50; i++ {
				_, _, _ = tokenizer.Encode("Test message for race detection", 128)
				_, _, _ = tokenizer.EncodeQueryPassagePair("query", "passage", 256)
			}
		},
		// Normalization
		func() {
			for i := 0; i < 100; i++ {
				_ = tokenizer.normalize("Test  text   with   spaces")
				_ = tokenizer.normalize("Another test string")
			}
		},
		// Pre-tokenization
		func() {
			for i := 0; i < 100; i++ {
				_ = tokenizer.preTokenize("Test tokenization string")
				_ = tokenizer.preTokenize("Multiple words here")
			}
		},
	}

	for _, op := range operations {
		for w := 0; w < numWorkers; w++ {
			wg.Add(1)
			go func(operation func()) {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done():
						return
					default:
						operation()
					}
				}
			}(op)
		}
	}

	// Let it run for a bit
	time.Sleep(5 * time.Second)
	cancel()
	wg.Wait()

	t.Log("Race condition test completed without deadlocks")
}

func TestPureGoTokenizer_MemoryStress(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory stress test in short mode")
	}

	tokenizerPath := createMinimalTokenizerConfig(t)
	defer removeFile(tokenizerPath)

	tokenizer, err := NewPureGoTokenizer(tokenizerPath)
	if err != nil {
		t.Fatalf("Failed to create tokenizer: %v", err)
	}

	// Get initial memory stats
	var initialStats runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&initialStats)

	const iterations = 1000
	const textSize = 1000

	// Generate large text
	largeText := make([]string, textSize)
	for i := range largeText {
		largeText[i] = "This is a memory stress test sentence that will be repeated many times."
	}
	text := "Memory stress test: " + joinStrings(largeText, " ")

	t.Logf("Testing with text of length: %d characters", len(text))

	// Perform many tokenizations
	for i := 0; i < iterations; i++ {
		tokenIDs, attentionMask, tokenTypeIDs := tokenizer.Encode(text, 512)
		
		// Verify basic properties
		if len(tokenIDs) != 512 || len(attentionMask) != 512 || len(tokenTypeIDs) != 512 {
			t.Errorf("Iteration %d: Invalid output lengths", i)
		}

		// Force garbage collection periodically
		if i%100 == 0 {
			runtime.GC()
		}
	}

	// Final memory check
	runtime.GC()
	var finalStats runtime.MemStats
	runtime.ReadMemStats(&finalStats)

	memoryIncrease := finalStats.Alloc - initialStats.Alloc
	t.Logf("Memory increase after %d iterations: %d bytes", iterations, memoryIncrease)

	// Check for reasonable memory usage (less than 100MB increase)
	if memoryIncrease > 100*1024*1024 {
		t.Errorf("Excessive memory usage detected: %d bytes", memoryIncrease)
	}
}

func TestPureGoTokenizer_ConcurrentCreation(t *testing.T) {
	const numCreations = 10

	var wg sync.WaitGroup
	errors := make(chan error, numCreations)

	for i := 0; i < numCreations; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			tokenizerPath := createMinimalTokenizerConfig(t)
			defer removeFile(tokenizerPath)

			tokenizer, err := NewPureGoTokenizer(tokenizerPath)
			if err != nil {
				errors <- &ConcurrentTestError{
					GoroutineID: id,
					Message:     "Failed to create tokenizer: " + err.Error(),
				}
				return
			}

			// Test basic functionality
			_, _, _ = tokenizer.Encode("test", 128)
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Error(err)
	}
}

func BenchmarkPureGoTokenizer_Concurrent(b *testing.B) {
	// Create a testing.T wrapper for createMinimalTokenizerConfig
	t := &testing.T{}
	tokenizerPath := createMinimalTokenizerConfig(t)
	defer removeFile(tokenizerPath)

	tokenizer, err := NewPureGoTokenizer(tokenizerPath)
	if err != nil {
		b.Fatalf("Failed to create tokenizer: %v", err)
	}

	text := "This is a benchmark test for concurrent tokenization performance"
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _, _ = tokenizer.Encode(text, 256)
		}
	})
}

// Helper types and functions

type ConcurrentTestError struct {
	GoroutineID int
	Iteration   int
	Message     string
}

func (e *ConcurrentTestError) Error() string {
	if e.Iteration >= 0 {
		return fmt.Sprintf("Goroutine %d, Iteration %d: %s", e.GoroutineID, e.Iteration, e.Message)
	}
	return fmt.Sprintf("Goroutine %d: %s", e.GoroutineID, e.Message)
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}
	
	var result strings.Builder
	result.WriteString(strs[0])
	for i := 1; i < len(strs); i++ {
		result.WriteString(sep)
		result.WriteString(strs[i])
	}
	return result.String()
}

func removeFile(path string) {
	// Small delay to ensure file is not in use
	time.Sleep(10 * time.Millisecond)
	os.Remove(path)
}