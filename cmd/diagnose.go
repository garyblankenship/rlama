package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/dontizi/rlama/internal/service"
	"github.com/spf13/cobra"
)

var (
	diagnoseSaveReport   bool
	diagnoseTestDocs     bool
	diagnoseVerbose      bool
	diagnoseOutputFile   string
)

var diagnoseCmd = &cobra.Command{
	Use:   "diagnose",
	Short: "Diagnose RLAMA setup and document processing capabilities",
	Long: `Diagnose RLAMA installation and document processing setup.

This command performs comprehensive checks of your RLAMA installation,
including enhanced document processing capabilities, configuration,
and system dependencies.

Examples:
  rlama diagnose
  rlama diagnose --test-docs --verbose
  rlama diagnose --save-report --output=diagnosis.txt`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runDiagnosis()
	},
}

func init() {
	rootCmd.AddCommand(diagnoseCmd)

	diagnoseCmd.Flags().BoolVar(&diagnoseSaveReport, "save-report", false, "Save diagnosis report to file")
	diagnoseCmd.Flags().BoolVar(&diagnoseTestDocs, "test-docs", false, "Test document processing with sample files")
	diagnoseCmd.Flags().BoolVar(&diagnoseVerbose, "verbose", false, "Show detailed diagnostic information")
	diagnoseCmd.Flags().StringVar(&diagnoseOutputFile, "output", "", "Output file for diagnosis report")
}

func runDiagnosis() error {
	fmt.Println("🔍 RLAMA System Diagnosis")
	fmt.Println("=" + strings.Repeat("=", 50))

	diagnosis := &DiagnosisReport{
		Timestamp: time.Now(),
		Sections:  make(map[string]*DiagnosisSection),
	}

	// Run all diagnostic checks
	checkSystemInfo(diagnosis)
	checkRLAMAInstallation(diagnosis)
	checkEnhancedProcessing(diagnosis)
	checkLLMProviders(diagnosis)
	checkVectorStores(diagnosis)
	checkConfiguration(diagnosis)
	
	if diagnoseTestDocs {
		testDocumentProcessing(diagnosis)
	}

	// Print summary
	printDiagnosisSummary(diagnosis)

	// Save report if requested
	if diagnoseSaveReport || diagnoseOutputFile != "" {
		return saveDiagnosisReport(diagnosis)
	}

	return nil
}

type DiagnosisReport struct {
	Timestamp time.Time
	Sections  map[string]*DiagnosisSection
}

type DiagnosisSection struct {
	Name    string
	Status  DiagnosisStatus
	Items   []DiagnosisItem
	Summary string
}

type DiagnosisItem struct {
	Name    string
	Status  DiagnosisStatus
	Message string
	Details string
}

type DiagnosisStatus int

const (
	StatusOK DiagnosisStatus = iota
	StatusWarning
	StatusError
	StatusInfo
)

func (s DiagnosisStatus) String() string {
	switch s {
	case StatusOK:
		return "✅"
	case StatusWarning:
		return "⚠️"
	case StatusError:
		return "❌"
	case StatusInfo:
		return "ℹ️"
	default:
		return "❓"
	}
}

func (s DiagnosisStatus) ColorString() string {
	switch s {
	case StatusOK:
		return "\033[32m✅\033[0m" // Green
	case StatusWarning:
		return "\033[33m⚠️\033[0m"  // Yellow
	case StatusError:
		return "\033[31m❌\033[0m" // Red
	case StatusInfo:
		return "\033[34mℹ️\033[0m"  // Blue
	default:
		return "❓"
	}
}

func checkSystemInfo(diagnosis *DiagnosisReport) {
	section := &DiagnosisSection{
		Name:   "System Information",
		Status: StatusInfo,
		Items:  []DiagnosisItem{},
	}

	// Operating System
	section.Items = append(section.Items, DiagnosisItem{
		Name:    "Operating System",
		Status:  StatusInfo,
		Message: fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH),
	})

	// Go Version
	section.Items = append(section.Items, DiagnosisItem{
		Name:    "Go Runtime",
		Status:  StatusInfo,
		Message: runtime.Version(),
	})

	// CPU Count
	section.Items = append(section.Items, DiagnosisItem{
		Name:    "CPU Cores",
		Status:  StatusInfo,
		Message: fmt.Sprintf("%d cores", runtime.NumCPU()),
	})

	// Working Directory
	if wd, err := os.Getwd(); err == nil {
		section.Items = append(section.Items, DiagnosisItem{
			Name:    "Working Directory",
			Status:  StatusInfo,
			Message: wd,
		})
	}

	diagnosis.Sections["system"] = section
}

func checkRLAMAInstallation(diagnosis *DiagnosisReport) {
	section := &DiagnosisSection{
		Name:   "RLAMA Installation",
		Status: StatusOK,
		Items:  []DiagnosisItem{},
	}

	// Check RLAMA config directory
	configDir := filepath.Join(os.Getenv("HOME"), ".rlama")
	if _, err := os.Stat(configDir); err == nil {
		section.Items = append(section.Items, DiagnosisItem{
			Name:    "Configuration Directory",
			Status:  StatusOK,
			Message: "Found at " + configDir,
		})
	} else {
		section.Items = append(section.Items, DiagnosisItem{
			Name:    "Configuration Directory",
			Status:  StatusWarning,
			Message: "Not found, will be created when needed",
		})
		section.Status = StatusWarning
	}

	// Check for existing RAGs
	ragList, err := runRLAMACommand("list")
	if err == nil {
		ragCount := countRAGs(ragList)
		section.Items = append(section.Items, DiagnosisItem{
			Name:    "Existing RAG Systems",
			Status:  StatusInfo,
			Message: fmt.Sprintf("%d RAG systems found", ragCount),
		})
	} else {
		section.Items = append(section.Items, DiagnosisItem{
			Name:    "RAG Systems Check",
			Status:  StatusError,
			Message: "Could not list RAG systems: " + err.Error(),
		})
		section.Status = StatusError
	}

	diagnosis.Sections["installation"] = section
}

func checkEnhancedProcessing(diagnosis *DiagnosisReport) {
	section := &DiagnosisSection{
		Name:   "Enhanced Document Processing",
		Status: StatusOK,
		Items:  []DiagnosisItem{},
	}

	// Check feature flags
	flags := service.GetAllFeatureFlags()
	
	section.Items = append(section.Items, DiagnosisItem{
		Name:    "LangChain Integration",
		Status:  boolToStatus(flags.UseLangChain),
		Message: fmt.Sprintf("Enabled: %t", flags.UseLangChain),
	})

	section.Items = append(section.Items, DiagnosisItem{
		Name:    "Loading Strategy",
		Status:  StatusInfo,
		Message: flags.LoaderStrategy,
	})

	section.Items = append(section.Items, DiagnosisItem{
		Name:    "Debug Mode",
		Status:  StatusInfo,
		Message: fmt.Sprintf("Enabled: %t", flags.DebugMode),
	})

	section.Items = append(section.Items, DiagnosisItem{
		Name:    "Telemetry Collection",
		Status:  StatusInfo,
		Message: fmt.Sprintf("Enabled: %t", flags.CollectTelemetry),
	})

	// Test strategy availability
	loader := service.NewEnhancedDocumentLoader()
	strategies := loader.GetAvailableStrategies()
	
	for name, info := range strategies {
		status := StatusOK
		if !info.Available {
			status = StatusError
		}
		
		section.Items = append(section.Items, DiagnosisItem{
			Name:    fmt.Sprintf("Strategy: %s", name),
			Status:  status,
			Message: info.Description,
			Details: fmt.Sprintf("Available: %t, File types: %d", info.Available, len(info.FileTypes)),
		})
		
		if !info.Available {
			section.Status = StatusWarning
		}
	}

	diagnosis.Sections["enhanced_processing"] = section
}

func checkLLMProviders(diagnosis *DiagnosisReport) {
	section := &DiagnosisSection{
		Name:   "LLM Providers",
		Status: StatusOK,
		Items:  []DiagnosisItem{},
	}

	// Check Ollama
	ollamaHost := os.Getenv("OLLAMA_HOST")
	if ollamaHost == "" {
		ollamaHost = "http://localhost:11434"
	}
	
	section.Items = append(section.Items, DiagnosisItem{
		Name:    "Ollama Host",
		Status:  StatusInfo,
		Message: ollamaHost,
	})

	// Check OpenAI
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey != "" {
		section.Items = append(section.Items, DiagnosisItem{
			Name:    "OpenAI API Key",
			Status:  StatusOK,
			Message: "Configured (hidden for security)",
		})
	} else {
		section.Items = append(section.Items, DiagnosisItem{
			Name:    "OpenAI API Key",
			Status:  StatusInfo,
			Message: "Not configured",
		})
	}

	diagnosis.Sections["llm_providers"] = section
}

func checkVectorStores(diagnosis *DiagnosisReport) {
	section := &DiagnosisSection{
		Name:   "Vector Stores",
		Status: StatusOK,
		Items:  []DiagnosisItem{},
	}

	// Internal vector store
	section.Items = append(section.Items, DiagnosisItem{
		Name:    "Internal Vector Store",
		Status:  StatusOK,
		Message: "Available (built-in)",
	})

	// Qdrant configuration
	qdrantHost := os.Getenv("QDRANT_HOST")
	if qdrantHost == "" {
		qdrantHost = "localhost"
	}
	
	qdrantPort := os.Getenv("QDRANT_PORT")
	if qdrantPort == "" {
		qdrantPort = "6333"
	}

	section.Items = append(section.Items, DiagnosisItem{
		Name:    "Qdrant Configuration",
		Status:  StatusInfo,
		Message: fmt.Sprintf("Host: %s, Port: %s", qdrantHost, qdrantPort),
	})

	diagnosis.Sections["vector_stores"] = section
}

func checkConfiguration(diagnosis *DiagnosisReport) {
	section := &DiagnosisSection{
		Name:   "Configuration",
		Status: StatusOK,
		Items:  []DiagnosisItem{},
	}

	// Check environment variables
	envVars := []string{
		"RLAMA_LOADER_STRATEGY",
		"RLAMA_USE_LANGCHAIN_LOADER",
		"RLAMA_DEBUG_LOADER",
		"RLAMA_COLLECT_TELEMETRY",
		"RLAMA_PREFERRED_CHUNK_SIZE",
		"RLAMA_PREFERRED_CHUNK_OVERLAP",
	}

	configuredCount := 0
	for _, envVar := range envVars {
		value := os.Getenv(envVar)
		if value != "" {
			configuredCount++
			section.Items = append(section.Items, DiagnosisItem{
				Name:    envVar,
				Status:  StatusInfo,
				Message: value,
			})
		}
	}

	section.Items = append(section.Items, DiagnosisItem{
		Name:    "Environment Variables",
		Status:  StatusInfo,
		Message: fmt.Sprintf("%d of %d configured", configuredCount, len(envVars)),
	})

	diagnosis.Sections["configuration"] = section
}

func testDocumentProcessing(diagnosis *DiagnosisReport) {
	section := &DiagnosisSection{
		Name:   "Document Processing Test",
		Status: StatusOK,
		Items:  []DiagnosisItem{},
	}

	// Create test directory
	testDir, err := os.MkdirTemp("", "rlama_diagnose_")
	if err != nil {
		section.Items = append(section.Items, DiagnosisItem{
			Name:    "Test Setup",
			Status:  StatusError,
			Message: "Failed to create test directory: " + err.Error(),
		})
		section.Status = StatusError
		diagnosis.Sections["document_test"] = section
		return
	}
	defer os.RemoveAll(testDir)

	// Create test files
	testFiles := map[string]string{
		"test.txt": "This is a test document for diagnostic purposes.",
		"test.md":  "# Test\nThis is a **markdown** test document.",
		"test.json": `{"test": true, "purpose": "diagnostics"}`,
	}

	for filename, content := range testFiles {
		err := os.WriteFile(filepath.Join(testDir, filename), []byte(content), 0644)
		if err != nil {
			section.Items = append(section.Items, DiagnosisItem{
				Name:    "Test File Creation",
				Status:  StatusError,
				Message: "Failed to create " + filename + ": " + err.Error(),
			})
			section.Status = StatusError
			continue
		}
	}

	// Test different strategies
	loader := service.NewEnhancedDocumentLoader()
	strategies := []string{"legacy", "langchain", "hybrid"}

	for _, strategy := range strategies {
		loader.SetStrategy(strategy)
		options := service.NewDocumentLoaderOptions()
		
		start := time.Now()
		docs, err := loader.LoadDocumentsFromFolderWithOptions(testDir, options)
		duration := time.Since(start)

		if err != nil {
			section.Items = append(section.Items, DiagnosisItem{
				Name:    fmt.Sprintf("Strategy: %s", strategy),
				Status:  StatusError,
				Message: "Failed: " + err.Error(),
			})
			section.Status = StatusWarning
		} else {
			section.Items = append(section.Items, DiagnosisItem{
				Name:    fmt.Sprintf("Strategy: %s", strategy),
				Status:  StatusOK,
				Message: fmt.Sprintf("Loaded %d docs in %v", len(docs), duration),
			})
		}
	}

	diagnosis.Sections["document_test"] = section
}

func printDiagnosisSummary(diagnosis *DiagnosisReport) {
	fmt.Printf("\n📋 DIAGNOSIS SUMMARY\n")
	fmt.Println("=" + strings.Repeat("=", 50))

	overallStatus := StatusOK
	for _, section := range diagnosis.Sections {
		fmt.Printf("\n%s %s\n", section.Status.ColorString(), section.Name)
		
		if section.Status > overallStatus {
			overallStatus = section.Status
		}

		for _, item := range section.Items {
			if diagnoseVerbose {
				fmt.Printf("  %s %s: %s\n", item.Status.ColorString(), item.Name, item.Message)
				if item.Details != "" {
					fmt.Printf("     %s\n", item.Details)
				}
			} else if item.Status == StatusError || item.Status == StatusWarning {
				fmt.Printf("  %s %s: %s\n", item.Status.ColorString(), item.Name, item.Message)
			}
		}
	}

	// Overall status
	fmt.Printf("\n🎯 OVERALL STATUS: %s\n", overallStatus.ColorString())
	
	switch overallStatus {
	case StatusOK:
		fmt.Println("✨ RLAMA is properly configured and ready to use!")
	case StatusWarning:
		fmt.Println("⚠️ RLAMA is functional but has some warnings. Check the details above.")
	case StatusError:
		fmt.Println("❌ RLAMA has configuration issues that need attention.")
	}

	// Recommendations
	fmt.Printf("\n💡 RECOMMENDATIONS:\n")
	if overallStatus == StatusOK {
		fmt.Println("• Your RLAMA setup looks great!")
		fmt.Println("• Try creating a RAG system: rlama rag llama3.2 my-docs ./documents")
		fmt.Println("• Benchmark performance: rlama benchmark ./documents")
	} else {
		fmt.Println("• Run with --verbose to see all details")
		fmt.Println("• Check the documentation: docs/enhanced_document_processing.md")
		fmt.Println("• Run the migration script: scripts/migrate_to_enhanced.sh")
	}
}

func saveDiagnosisReport(diagnosis *DiagnosisReport) error {
	filename := diagnoseOutputFile
	if filename == "" {
		filename = fmt.Sprintf("rlama_diagnosis_%s.txt", 
			diagnosis.Timestamp.Format("20060102_150405"))
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write report header
	fmt.Fprintf(file, "RLAMA System Diagnosis Report\n")
	fmt.Fprintf(file, "Generated: %s\n", diagnosis.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(file, "System: %s %s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Fprintf(file, "%s\n\n", strings.Repeat("=", 60))

	// Write sections
	for _, section := range diagnosis.Sections {
		fmt.Fprintf(file, "%s %s\n", section.Status.String(), section.Name)
		fmt.Fprintf(file, "%s\n", strings.Repeat("-", 40))
		
		for _, item := range section.Items {
			fmt.Fprintf(file, "  %s %s: %s\n", item.Status.String(), item.Name, item.Message)
			if item.Details != "" {
				fmt.Fprintf(file, "     Details: %s\n", item.Details)
			}
		}
		fmt.Fprintf(file, "\n")
	}

	fmt.Printf("💾 Diagnosis report saved to: %s\n", filename)
	return nil
}

// Helper functions
func boolToStatus(b bool) DiagnosisStatus {
	if b {
		return StatusOK
	}
	return StatusWarning
}

func runRLAMACommand(args ...string) (string, error) {
	// This is a placeholder for running RLAMA commands
	// In a real implementation, you might execute the command and capture output
	return "", fmt.Errorf("command execution not implemented in this version")
}

func countRAGs(output string) int {
	// Simple implementation - count non-empty lines that aren't headers
	lines := strings.Split(output, "\n")
	count := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.Contains(line, "Available RAG systems:") && !strings.Contains(line, "No RAG systems") {
			count++
		}
	}
	return count
}