package service

import (
	"os"
	"strings"
)

// Feature flags for controlling document loading behavior

// UseLangChainLoader returns whether to use LangChain loader
func UseLangChainLoader() bool {
	value := os.Getenv("RLAMA_USE_LANGCHAIN_LOADER")
	return value != "false" && value != "0"
}

// GetLoaderStrategy returns the document loading strategy from environment
func GetLoaderStrategy() string {
	strategy := strings.ToLower(os.Getenv("RLAMA_LOADER_STRATEGY"))
	
	validStrategies := map[string]bool{
		"langchain": true,
		"legacy":    true,
		"hybrid":    true,
	}
	
	if validStrategies[strategy] {
		return strategy
	}
	
	// Default to hybrid for safety
	return "hybrid"
}

// GetLoaderTimeout returns the timeout for document loading operations in minutes
func GetLoaderTimeout() int {
	timeoutStr := os.Getenv("RLAMA_LOADER_TIMEOUT_MINUTES")
	if timeoutStr == "" {
		return 5 // Default 5 minutes
	}
	
	// Simple parsing - if invalid, return default
	switch timeoutStr {
	case "1":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	case "5":
		return 5
	case "10":
		return 10
	case "15":
		return 15
	case "30":
		return 30
	default:
		return 5
	}
}

// IsDebugMode returns whether debug mode is enabled for document loading
func IsDebugMode() bool {
	value := os.Getenv("RLAMA_DEBUG_LOADER")
	return value == "true" || value == "1"
}

// GetMaxRetries returns the maximum number of retry attempts for document loading
func GetMaxRetries() int {
	retriesStr := os.Getenv("RLAMA_LOADER_MAX_RETRIES")
	switch retriesStr {
	case "1":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	case "5":
		return 5
	default:
		return 3 // Default to 3 retries
	}
}

// ShouldCollectTelemetry returns whether to collect and report telemetry
func ShouldCollectTelemetry() bool {
	value := os.Getenv("RLAMA_COLLECT_TELEMETRY")
	return value != "false" && value != "0"
}

// GetPreferredChunkSize returns the preferred chunk size from environment
func GetPreferredChunkSize() int {
	sizeStr := os.Getenv("RLAMA_PREFERRED_CHUNK_SIZE")
	switch sizeStr {
	case "500":
		return 500
	case "750":
		return 750
	case "1000":
		return 1000
	case "1500":
		return 1500
	case "2000":
		return 2000
	default:
		return 1000 // Default chunk size
	}
}

// GetPreferredChunkOverlap returns the preferred chunk overlap from environment
func GetPreferredChunkOverlap() int {
	overlapStr := os.Getenv("RLAMA_PREFERRED_CHUNK_OVERLAP")
	switch overlapStr {
	case "50":
		return 50
	case "100":
		return 100
	case "150":
		return 150
	case "200":
		return 200
	case "250":
		return 250
	case "300":
		return 300
	default:
		return 200 // Default overlap
	}
}

// FeatureFlags contains all feature flag values for easy access
type FeatureFlags struct {
	UseLangChain      bool
	LoaderStrategy    string
	TimeoutMinutes    int
	DebugMode         bool
	MaxRetries        int
	CollectTelemetry  bool
	ChunkSize         int
	ChunkOverlap      int
}

// GetAllFeatureFlags returns a struct with all current feature flag values
func GetAllFeatureFlags() FeatureFlags {
	return FeatureFlags{
		UseLangChain:     UseLangChainLoader(),
		LoaderStrategy:   GetLoaderStrategy(),
		TimeoutMinutes:   GetLoaderTimeout(),
		DebugMode:        IsDebugMode(),
		MaxRetries:       GetMaxRetries(),
		CollectTelemetry: ShouldCollectTelemetry(),
		ChunkSize:        GetPreferredChunkSize(),
		ChunkOverlap:     GetPreferredChunkOverlap(),
	}
}

// PrintFeatureFlags prints current feature flag values
func PrintFeatureFlags() {
	flags := GetAllFeatureFlags()
	
	println("🏁 RLAMA Feature Flags:")
	println("   RLAMA_USE_LANGCHAIN_LOADER:", boolToString(flags.UseLangChain))
	println("   RLAMA_LOADER_STRATEGY:", flags.LoaderStrategy)
	printf("   RLAMA_LOADER_TIMEOUT_MINUTES: %d\n", flags.TimeoutMinutes)
	println("   RLAMA_DEBUG_LOADER:", boolToString(flags.DebugMode))
	printf("   RLAMA_LOADER_MAX_RETRIES: %d\n", flags.MaxRetries)
	println("   RLAMA_COLLECT_TELEMETRY:", boolToString(flags.CollectTelemetry))
	printf("   RLAMA_PREFERRED_CHUNK_SIZE: %d\n", flags.ChunkSize)
	printf("   RLAMA_PREFERRED_CHUNK_OVERLAP: %d\n", flags.ChunkOverlap)
}

// SetDefaultEnvironmentForTesting sets reasonable defaults for testing
func SetDefaultEnvironmentForTesting() {
	os.Setenv("RLAMA_USE_LANGCHAIN_LOADER", "true")
	os.Setenv("RLAMA_LOADER_STRATEGY", "hybrid")
	os.Setenv("RLAMA_LOADER_TIMEOUT_MINUTES", "2")
	os.Setenv("RLAMA_DEBUG_LOADER", "false")
	os.Setenv("RLAMA_LOADER_MAX_RETRIES", "2")
	os.Setenv("RLAMA_COLLECT_TELEMETRY", "false")
	os.Setenv("RLAMA_PREFERRED_CHUNK_SIZE", "500")
	os.Setenv("RLAMA_PREFERRED_CHUNK_OVERLAP", "100")
}

// Helper functions
func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// Simple printf implementation to avoid importing fmt in feature flags
func printf(format string, args ...interface{}) {
	// This is a simplified implementation for feature flags
	// In a real implementation you might want to use fmt
	if len(args) == 1 {
		if format == "   RLAMA_LOADER_TIMEOUT_MINUTES: %d\n" {
			if val, ok := args[0].(int); ok {
				print("   RLAMA_LOADER_TIMEOUT_MINUTES: ")
				print(val)
				print("\n")
			}
		} else if format == "   RLAMA_LOADER_MAX_RETRIES: %d\n" {
			if val, ok := args[0].(int); ok {
				print("   RLAMA_LOADER_MAX_RETRIES: ")
				print(val)
				print("\n")
			}
		} else if format == "   RLAMA_PREFERRED_CHUNK_SIZE: %d\n" {
			if val, ok := args[0].(int); ok {
				print("   RLAMA_PREFERRED_CHUNK_SIZE: ")
				print(val)
				print("\n")
			}
		} else if format == "   RLAMA_PREFERRED_CHUNK_OVERLAP: %d\n" {
			if val, ok := args[0].(int); ok {
				print("   RLAMA_PREFERRED_CHUNK_OVERLAP: ")
				print(val)
				print("\n")
			}
		}
	}
}