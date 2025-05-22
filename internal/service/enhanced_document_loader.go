package service

import (
	"fmt"
	"os"
	"time"

	"github.com/dontizi/rlama/internal/domain"
)

// EnhancedDocumentLoader implements the strategy pattern for document loading
// with graceful fallback between LangChain and legacy loaders
type EnhancedDocumentLoader struct {
	legacyStrategy    DocumentLoaderStrategy
	langchainStrategy DocumentLoaderStrategy
	strategy          string
	telemetry         *LoaderTelemetry
}

// LoaderTelemetry tracks usage statistics for different loading strategies
type LoaderTelemetry struct {
	LangChainSuccesses int
	LangChainFailures  int
	LegacySuccesses    int
	LegacyFailures     int
	LastUpdated        time.Time
}

// NewEnhancedDocumentLoader creates a new enhanced document loader with strategy selection
func NewEnhancedDocumentLoader() *EnhancedDocumentLoader {
	return &EnhancedDocumentLoader{
		legacyStrategy:    NewLegacyDocumentLoaderStrategy(),
		langchainStrategy: NewLangChainDocumentLoaderStrategy(), // This will use the v2 implementation
		strategy:          getLoaderStrategy(),
		telemetry: &LoaderTelemetry{
			LastUpdated: time.Now(),
		},
	}
}

// LoadDocumentsFromFolderWithOptions loads documents using the configured strategy with fallback
func (edl *EnhancedDocumentLoader) LoadDocumentsFromFolderWithOptions(folderPath string, options DocumentLoaderOptions) ([]*domain.Document, error) {
	var docs []*domain.Document
	var err error

	switch edl.strategy {
	case "langchain":
		docs, err = edl.loadWithLangChain(folderPath, options)
		
	case "legacy":
		docs, err = edl.loadWithLegacy(folderPath, options)
		
	case "hybrid":
		docs, err = edl.loadWithHybrid(folderPath, options)
		
	default:
		// Default to hybrid for unknown strategies
		fmt.Printf("⚠️ Unknown strategy '%s', defaulting to hybrid\n", edl.strategy)
		docs, err = edl.loadWithHybrid(folderPath, options)
	}

	edl.telemetry.LastUpdated = time.Now()
	return docs, err
}

// loadWithLangChain attempts to load documents using only LangChain
func (edl *EnhancedDocumentLoader) loadWithLangChain(folderPath string, options DocumentLoaderOptions) ([]*domain.Document, error) {
	if !edl.langchainStrategy.IsAvailable() {
		edl.telemetry.LangChainFailures++
		return nil, fmt.Errorf("LangChain strategy is not available")
	}

	docs, err := edl.langchainStrategy.LoadDocuments(folderPath, options)
	if err != nil {
		edl.telemetry.LangChainFailures++
		return nil, fmt.Errorf("LangChain loading failed: %w", err)
	}

	edl.telemetry.LangChainSuccesses++
	return docs, nil
}

// loadWithLegacy attempts to load documents using only the legacy loader
func (edl *EnhancedDocumentLoader) loadWithLegacy(folderPath string, options DocumentLoaderOptions) ([]*domain.Document, error) {
	docs, err := edl.legacyStrategy.LoadDocuments(folderPath, options)
	if err != nil {
		edl.telemetry.LegacyFailures++
		return nil, fmt.Errorf("legacy loading failed: %w", err)
	}

	edl.telemetry.LegacySuccesses++
	return docs, nil
}

// loadWithHybrid attempts LangChain first, falls back to legacy on failure
func (edl *EnhancedDocumentLoader) loadWithHybrid(folderPath string, options DocumentLoaderOptions) ([]*domain.Document, error) {
	// Try LangChain first if available
	if edl.langchainStrategy.IsAvailable() {
		docs, err := edl.langchainStrategy.LoadDocuments(folderPath, options)
		if err == nil {
			edl.telemetry.LangChainSuccesses++
			return docs, nil
		}

		// Log the LangChain failure but continue with fallback
		fmt.Printf("⚠️ LangChain loading failed, falling back to legacy: %v\n", err)
		edl.telemetry.LangChainFailures++
	} else {
		fmt.Printf("⚠️ LangChain not available, using legacy loader\n")
	}

	// Fallback to legacy loader
	docs, err := edl.legacyStrategy.LoadDocuments(folderPath, options)
	if err != nil {
		edl.telemetry.LegacyFailures++
		return nil, fmt.Errorf("both LangChain and legacy loading failed. Legacy error: %w", err)
	}

	edl.telemetry.LegacySuccesses++
	fmt.Printf("✅ Successfully loaded %d documents using legacy fallback\n", len(docs))
	return docs, nil
}

// SetStrategy changes the loading strategy
func (edl *EnhancedDocumentLoader) SetStrategy(strategy string) {
	validStrategies := map[string]bool{
		"langchain": true,
		"legacy":    true,
		"hybrid":    true,
	}

	if !validStrategies[strategy] {
		fmt.Printf("⚠️ Invalid strategy '%s', keeping current strategy '%s'\n", strategy, edl.strategy)
		return
	}

	edl.strategy = strategy
	fmt.Printf("✅ Document loading strategy changed to: %s\n", strategy)
}

// GetStrategy returns the current loading strategy
func (edl *EnhancedDocumentLoader) GetStrategy() string {
	return edl.strategy
}

// GetTelemetry returns usage statistics
func (edl *EnhancedDocumentLoader) GetTelemetry() *LoaderTelemetry {
	return edl.telemetry
}

// GetSupportedFileTypes returns all supported file types across strategies
func (edl *EnhancedDocumentLoader) GetSupportedFileTypes() []string {
	// Combine supported types from both strategies
	langchainTypes := edl.langchainStrategy.GetSupportedFileTypes()
	legacyTypes := edl.legacyStrategy.GetSupportedFileTypes()

	// Create a set to avoid duplicates
	typeSet := make(map[string]bool)
	for _, t := range langchainTypes {
		typeSet[t] = true
	}
	for _, t := range legacyTypes {
		typeSet[t] = true
	}

	// Convert back to slice
	var allTypes []string
	for t := range typeSet {
		allTypes = append(allTypes, t)
	}

	return allTypes
}

// GetAvailableStrategies returns information about available strategies
func (edl *EnhancedDocumentLoader) GetAvailableStrategies() map[string]StrategyInfo {
	return map[string]StrategyInfo{
		"langchain": {
			Name:        edl.langchainStrategy.GetName(),
			Available:   edl.langchainStrategy.IsAvailable(),
			FileTypes:   edl.langchainStrategy.GetSupportedFileTypes(),
			Description: "Advanced document processing using LangChainGo with robust error handling",
		},
		"legacy": {
			Name:        edl.legacyStrategy.GetName(),
			Available:   edl.legacyStrategy.IsAvailable(),
			FileTypes:   edl.legacyStrategy.GetSupportedFileTypes(),
			Description: "Original RLAMA document processor with external tool support",
		},
		"hybrid": {
			Name:        "hybrid",
			Available:   true,
			FileTypes:   edl.GetSupportedFileTypes(),
			Description: "Try LangChain first, fallback to legacy on failure (recommended)",
		},
	}
}

// StrategyInfo contains information about a loading strategy
type StrategyInfo struct {
	Name        string
	Available   bool
	FileTypes   []string
	Description string
}

// PrintTelemetryReport prints a human-readable telemetry report
func (edl *EnhancedDocumentLoader) PrintTelemetryReport() {
	t := edl.telemetry
	total := t.LangChainSuccesses + t.LangChainFailures + t.LegacySuccesses + t.LegacyFailures

	if total == 0 {
		fmt.Println("📊 No document loading operations recorded yet")
		return
	}

	fmt.Println("📊 Document Loading Telemetry Report")
	fmt.Printf("   Current Strategy: %s\n", edl.strategy)
	fmt.Printf("   Last Updated: %s\n", t.LastUpdated.Format("2006-01-02 15:04:05"))
	fmt.Println()
	
	fmt.Printf("   LangChain: %d successes, %d failures", t.LangChainSuccesses, t.LangChainFailures)
	if t.LangChainSuccesses+t.LangChainFailures > 0 {
		successRate := float64(t.LangChainSuccesses) / float64(t.LangChainSuccesses+t.LangChainFailures) * 100
		fmt.Printf(" (%.1f%% success rate)", successRate)
	}
	fmt.Println()
	
	fmt.Printf("   Legacy: %d successes, %d failures", t.LegacySuccesses, t.LegacyFailures)
	if t.LegacySuccesses+t.LegacyFailures > 0 {
		successRate := float64(t.LegacySuccesses) / float64(t.LegacySuccesses+t.LegacyFailures) * 100
		fmt.Printf(" (%.1f%% success rate)", successRate)
	}
	fmt.Println()
	
	fmt.Printf("   Total Operations: %d\n", total)
}

// getLoaderStrategy determines the strategy from environment variables
func getLoaderStrategy() string {
	// Check environment variable first
	if strategy := os.Getenv("RLAMA_LOADER_STRATEGY"); strategy != "" {
		// Validate strategy
		validStrategies := map[string]bool{
			"langchain": true,
			"legacy":    true,
			"hybrid":    true,
		}
		
		if validStrategies[strategy] {
			return strategy
		}
		// Invalid strategy, fall back to default
	}

	// Check if LangChain should be disabled
	if os.Getenv("RLAMA_USE_LANGCHAIN_LOADER") == "false" {
		return "legacy"
	}

	// Default to hybrid (safe option)
	return "hybrid"
}

// Compatibility methods to match existing DocumentLoader interface

// LoadDocumentsFromFolder provides compatibility with the old interface
func (edl *EnhancedDocumentLoader) LoadDocumentsFromFolder(folderPath string) ([]*domain.Document, error) {
	return edl.LoadDocumentsFromFolderWithOptions(folderPath, NewDocumentLoaderOptions())
}

// CreateRagWithOptions provides compatibility for RAG creation
func (edl *EnhancedDocumentLoader) CreateRagWithOptions(options DocumentLoaderOptions) (*domain.RagSystem, error) {
	// This method might need folderPath - for now return error suggesting proper usage
	return nil, fmt.Errorf("use LoadDocumentsFromFolderWithOptions with a specific folder path instead")
}