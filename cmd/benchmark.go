package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dontizi/rlama/internal/service"
	"github.com/spf13/cobra"
)

var (
	benchmarkStrategy    string
	benchmarkRuns        int
	benchmarkFileCount   int
	benchmarkFileSize    int
	benchmarkOutputFile  string
	benchmarkVerbose     bool
)

var benchmarkCmd = &cobra.Command{
	Use:   "benchmark [folder-path]",
	Short: "Benchmark document loading performance",
	Long: `Benchmark the performance of different document loading strategies.

This command helps you evaluate the performance characteristics of the enhanced 
document processing system by running controlled tests on document loading.

Examples:
  rlama benchmark ./docs
  rlama benchmark ./docs --strategy=all --runs=10
  rlama benchmark ./docs --strategy=langchain --verbose
  rlama benchmark --generate-test-data --file-count=100`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var testDir string
		var shouldCleanup bool

		// Determine test directory
		if len(args) > 0 {
			testDir = args[0]
		} else {
			// Generate test data
			var err error
			testDir, err = generateTestData()
			if err != nil {
				return fmt.Errorf("failed to generate test data: %w", err)
			}
			shouldCleanup = true
			defer func() {
				if shouldCleanup {
					os.RemoveAll(testDir)
				}
			}()
		}

		return runBenchmark(testDir)
	},
}

func init() {
	rootCmd.AddCommand(benchmarkCmd)

	benchmarkCmd.Flags().StringVar(&benchmarkStrategy, "strategy", "all", "Strategy to benchmark (all, langchain, legacy, hybrid)")
	benchmarkCmd.Flags().IntVar(&benchmarkRuns, "runs", 5, "Number of benchmark runs per strategy")
	benchmarkCmd.Flags().IntVar(&benchmarkFileCount, "file-count", 50, "Number of test files to generate")
	benchmarkCmd.Flags().IntVar(&benchmarkFileSize, "file-size", 1000, "Size of each test file in characters")
	benchmarkCmd.Flags().StringVar(&benchmarkOutputFile, "output", "", "Output benchmark results to file")
	benchmarkCmd.Flags().BoolVar(&benchmarkVerbose, "verbose", false, "Enable verbose output")
	benchmarkCmd.Flags().Bool("generate-test-data", false, "Generate test data instead of using existing folder")
}

func generateTestData() (string, error) {
	testDir, err := os.MkdirTemp("", "rlama_benchmark_")
	if err != nil {
		return "", err
	}

	fmt.Printf("📁 Generating %d test files in %s\n", benchmarkFileCount, testDir)

	// File templates for different types
	templates := map[string]string{
		".txt": "This is a sample text document for benchmarking purposes. " +
			"It contains multiple sentences and paragraphs to simulate real-world content. " +
			"The document discusses various topics including technology, science, and literature. " +
			"Content length: %d characters. File number: %d.",
		
		".md": "# Benchmark Document %d\n\n" +
			"This is a **markdown** document for testing purposes.\n\n" +
			"## Section 1\n\n" +
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. " +
			"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.\n\n" +
			"### Subsection\n\n" +
			"- Item 1\n- Item 2\n- Item 3\n\n" +
			"Content length: %d characters.",
		
		".json": `{
  "id": %d,
  "title": "Benchmark Document",
  "description": "This is a JSON document for benchmarking document loading performance.",
  "data": {
    "type": "test",
    "size": %d,
    "content": "Sample JSON content for performance testing."
  },
  "timestamp": "2024-01-01T00:00:00Z"
}`,
		
		".py": `#!/usr/bin/env python3
"""
Benchmark Python file %d
Content length: %d characters
"""

def benchmark_function():
    """Sample function for testing."""
    data = [i for i in range(100)]
    return sum(data)

class BenchmarkClass:
    """Sample class for testing."""
    
    def __init__(self, value):
        self.value = value
    
    def process(self):
        return self.value * 2

if __name__ == "__main__":
    bc = BenchmarkClass(42)
    result = bc.process()
    print(f"Result: {result}")
`,
		
		".go": `package main

import "fmt"

// BenchmarkStruct represents test data for file %d
// Content length: %d characters
type BenchmarkStruct struct {
	ID    int
	Value string
}

func main() {
	data := BenchmarkStruct{
		ID:    %d,
		Value: "benchmark test",
	}
	
	fmt.Printf("Processing: %%+v\n", data)
}

func processData(data []int) int {
	sum := 0
	for _, v := range data {
		sum += v
	}
	return sum
}
`,
	}

	extensions := []string{".txt", ".md", ".json", ".py", ".go"}
	
	for i := 0; i < benchmarkFileCount; i++ {
		ext := extensions[i%len(extensions)]
		template := templates[ext]
		
		var content string
		switch ext {
		case ".txt", ".md":
			content = fmt.Sprintf(template, i, benchmarkFileSize)
		case ".json":
			content = fmt.Sprintf(template, i, benchmarkFileSize)
		case ".py", ".go":
			content = fmt.Sprintf(template, i, benchmarkFileSize, i)
		}
		
		// Pad content to reach desired size
		for len(content) < benchmarkFileSize {
			content += " Additional content to reach the desired file size."
		}
		content = content[:benchmarkFileSize]
		
		filename := filepath.Join(testDir, fmt.Sprintf("benchmark_%03d%s", i, ext))
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
			return "", err
		}
	}

	fmt.Printf("✅ Generated %d test files\n", benchmarkFileCount)
	return testDir, nil
}

func runBenchmark(testDir string) error {
	fmt.Printf("🚀 Starting benchmark with %d runs per strategy\n", benchmarkRuns)
	fmt.Printf("📁 Test directory: %s\n", testDir)

	// Check directory exists and get file count
	files, err := filepath.Glob(filepath.Join(testDir, "*"))
	if err != nil {
		return err
	}
	fmt.Printf("📊 Testing with %d files\n", len(files))

	// Determine strategies to test
	var strategies []string
	if benchmarkStrategy == "all" {
		strategies = []string{"legacy", "langchain", "hybrid"}
	} else {
		strategies = []string{benchmarkStrategy}
	}

	results := make(map[string]BenchmarkResult)
	
	for _, strategy := range strategies {
		fmt.Printf("\n🔍 Benchmarking strategy: %s\n", strategy)
		
		result, err := benchmarkDocumentStrategy(testDir, strategy)
		if err != nil {
			fmt.Printf("❌ Strategy %s failed: %v\n", strategy, err)
			continue
		}
		
		results[strategy] = result
		printStrategyResults(strategy, result)
	}

	// Print summary
	printBenchmarkSummary(results)

	// Save results to file if requested
	if benchmarkOutputFile != "" {
		if err := saveBenchmarkResults(results, benchmarkOutputFile); err != nil {
			fmt.Printf("⚠️ Failed to save results to file: %v\n", err)
		} else {
			fmt.Printf("💾 Results saved to: %s\n", benchmarkOutputFile)
		}
	}

	return nil
}

func benchmarkDocumentStrategy(testDir, strategy string) (BenchmarkResult, error) {
	loader := service.NewEnhancedDocumentLoader()
	loader.SetStrategy(strategy)
	
	options := service.NewDocumentLoaderOptions()
	
	var totalDuration time.Duration
	var totalDocuments int
	var successfulRuns int
	var failures []string

	for run := 1; run <= benchmarkRuns; run++ {
		start := time.Now()
		docs, err := loader.LoadDocumentsFromFolderWithOptions(testDir, options)
		duration := time.Since(start)
		
		if err != nil {
			failures = append(failures, fmt.Sprintf("Run %d: %v", run, err))
			if benchmarkVerbose {
				fmt.Printf("  Run %d: FAILED (%v)\n", run, err)
			}
			continue
		}
		
		successfulRuns++
		totalDuration += duration
		totalDocuments = len(docs) // Should be same for all runs
		
		if benchmarkVerbose {
			fmt.Printf("  Run %d: %v (%d docs)\n", run, duration, len(docs))
		}
	}

	if successfulRuns == 0 {
		return BenchmarkResult{}, fmt.Errorf("all runs failed")
	}

	avgDuration := totalDuration / time.Duration(successfulRuns)
	
	return BenchmarkResult{
		Strategy:       strategy,
		TotalRuns:      benchmarkRuns,
		SuccessfulRuns: successfulRuns,
		FailedRuns:     benchmarkRuns - successfulRuns,
		TotalDocuments: totalDocuments,
		TotalDuration:  totalDuration,
		AverageDuration: avgDuration,
		DocsPerSecond:   float64(totalDocuments) / avgDuration.Seconds(),
		Failures:       failures,
	}, nil
}

type BenchmarkResult struct {
	Strategy        string
	TotalRuns       int
	SuccessfulRuns  int
	FailedRuns      int
	TotalDocuments  int
	TotalDuration   time.Duration
	AverageDuration time.Duration
	DocsPerSecond   float64
	Failures        []string
}

func printStrategyResults(strategy string, result BenchmarkResult) {
	fmt.Printf("  ✅ Success Rate: %d/%d (%.1f%%)\n", 
		result.SuccessfulRuns, result.TotalRuns, 
		float64(result.SuccessfulRuns)/float64(result.TotalRuns)*100)
	
	fmt.Printf("  ⏱️  Average Time: %v\n", result.AverageDuration)
	fmt.Printf("  📄 Documents: %d\n", result.TotalDocuments)
	fmt.Printf("  🚀 Throughput: %.1f docs/second\n", result.DocsPerSecond)
	
	if len(result.Failures) > 0 && benchmarkVerbose {
		fmt.Printf("  ❌ Failures:\n")
		for _, failure := range result.Failures {
			fmt.Printf("     %s\n", failure)
		}
	}
}

func printBenchmarkSummary(results map[string]BenchmarkResult) {
	if len(results) <= 1 {
		return
	}

	fmt.Printf("\n📊 BENCHMARK SUMMARY\n")
	fmt.Printf("═══════════════════════════════════════════════════════════\n")

	// Find fastest strategy
	var fastest string
	var fastestTime time.Duration = time.Hour // Start with a large value
	
	for strategy, result := range results {
		if result.SuccessfulRuns > 0 && result.AverageDuration < fastestTime {
			fastest = strategy
			fastestTime = result.AverageDuration
		}
	}

	// Print comparison table
	fmt.Printf("%-12s | %-12s | %-12s | %-15s | %-10s\n", 
		"Strategy", "Avg Time", "Success Rate", "Docs/Second", "Relative")
	fmt.Printf("─────────────┼──────────────┼──────────────┼─────────────────┼───────────\n")

	for _, strategy := range []string{"langchain", "legacy", "hybrid"} {
		result, exists := results[strategy]
		if !exists {
			continue
		}

		successRate := float64(result.SuccessfulRuns) / float64(result.TotalRuns) * 100
		relative := "baseline"
		if strategy != fastest && result.SuccessfulRuns > 0 {
			slowdown := float64(result.AverageDuration) / float64(fastestTime)
			relative = fmt.Sprintf("%.1fx slower", slowdown)
		} else if strategy == fastest {
			relative = "fastest ⚡"
		}

		fmt.Printf("%-12s | %-12v | %10.1f%% | %13.1f | %-10s\n",
			strategy, result.AverageDuration, successRate, result.DocsPerSecond, relative)
	}

	// Print recommendation
	fmt.Printf("\n💡 RECOMMENDATION\n")
	if fastest == "hybrid" {
		fmt.Printf("   Hybrid strategy provides the best balance of speed and reliability.\n")
	} else if fastest == "langchain" {
		fmt.Printf("   LangChain strategy is fastest. Consider using it if reliability is proven.\n")
	} else {
		fmt.Printf("   %s strategy performed best in this test.\n", fastest)
	}
}

func saveBenchmarkResults(results map[string]BenchmarkResult, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write CSV header
	fmt.Fprintf(file, "Strategy,TotalRuns,SuccessfulRuns,FailedRuns,TotalDocuments,AverageDurationMs,DocsPerSecond\n")

	// Write data
	for strategy, result := range results {
		fmt.Fprintf(file, "%s,%d,%d,%d,%d,%.2f,%.2f\n",
			strategy,
			result.TotalRuns,
			result.SuccessfulRuns,
			result.FailedRuns,
			result.TotalDocuments,
			float64(result.AverageDuration.Nanoseconds())/1e6, // Convert to milliseconds
			result.DocsPerSecond)
	}

	return nil
}