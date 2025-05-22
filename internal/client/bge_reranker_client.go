package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// BGERerankerClient handles interactions with the BGE Reranker model via Python
type BGERerankerClient struct {
	modelName string
	useFP16   bool
	silent    bool
}

// NewBGERerankerClient creates a new instance of BGERerankerClient
func NewBGERerankerClient(modelName string) *BGERerankerClient {
	return NewBGERerankerClientWithOptions(modelName, false)
}

// NewBGERerankerClientWithOptions creates a new instance of BGERerankerClient with additional options
func NewBGERerankerClientWithOptions(modelName string, silent bool) *BGERerankerClient {
	client := &BGERerankerClient{
		modelName: modelName,
		useFP16:   true,
		silent:    silent,
	}

	// Check dependencies and model
	if !silent {
		if err := client.CheckDependencies(); err != nil {
			fmt.Printf("⚠️ Warning: %v\n", err)
			fmt.Println("To install dependencies, run: rlama install-dependencies")
		} else {
			// Only check model if dependencies are available
			if err := client.CheckModelExists(); err != nil {
				fmt.Printf("⚠️ Warning: %v\n", err)
				fmt.Println("The BGE Reranker model might not be accessible. Check your internet connection and model name.")
			}
		}
	}

	return client
}

// GetModelName returns the model name used by this client
func (c *BGERerankerClient) GetModelName() string {
	return c.modelName
}

// ComputeScores calculates relevance scores between queries and passages
func (c *BGERerankerClient) ComputeScores(pairs [][]string, normalize bool) ([]float64, error) {
	// Convert Go boolean to Python boolean
	normalizeStr := "False"
	if normalize {
		normalizeStr = "True"
	}

	useFP16Str := "False"
	if c.useFP16 {
		useFP16Str = "True"
	}

	pythonScript := fmt.Sprintf(`
import sys
import os
import json
import warnings

# Maximum warning suppression
warnings.filterwarnings("ignore")
os.environ["TOKENIZERS_PARALLELISM"] = "false"
os.environ["TF_CPP_MIN_LOG_LEVEL"] = "3"
os.environ["TRANSFORMERS_VERBOSITY"] = "error"
os.environ["PYTHONWARNINGS"] = "ignore"

# Specific patch for protobuf/MessageFactory errors
# Monkey patch the protobuf module to prevent MessageFactory errors
# This must be done before importing transformers or FlagEmbedding
try:
    import builtins
    original_import = builtins.__import__
    
    def custom_import(name, *args, **kwargs):
        module = original_import(name, *args, **kwargs)
        # Patch the protobuf module to suppress MessageFactory errors
        if name == 'google.protobuf.descriptor' or name.endswith('.descriptor'):
            if hasattr(module, 'MessageFactory'):
                original_get_prototype = getattr(module.MessageFactory, 'GetPrototype', None)
                if original_get_prototype:
                    def silent_get_prototype(*args, **kwargs):
                        try:
                            return original_get_prototype(*args, **kwargs)
                        except AttributeError:
                            return None
                    module.MessageFactory.GetPrototype = silent_get_prototype
        return module
    
    builtins.__import__ = custom_import
except Exception:
    pass  # If patching fails, continue with standard error suppression

# Redirect stderr to prevent protobuf errors
# This is a workaround for "MessageFactory has no attribute GetPrototype" errors
orig_stderr = sys.stderr
sys.stderr = open(os.devnull, 'w')

# Also patch sys.excepthook to avoid printing uncaught exceptions to stderr
orig_excepthook = sys.excepthook
def silent_excepthook(exctype, value, traceback):
    if orig_stderr != sys.stderr:  # Only silence if we're already redirecting stderr
        return
    # Still log critical errors to original stderr
    if exctype in (SystemExit, KeyboardInterrupt, MemoryError):
        orig_excepthook(exctype, value, traceback)
sys.excepthook = silent_excepthook

try:
    from FlagEmbedding import FlagReranker
    
    # Load the input data
    input_data = json.loads(sys.stdin.read())
    pairs = input_data["pairs"]
    normalize = input_data["normalize"]
    
    # Initialize the reranker
    reranker = FlagReranker('%s', use_fp16=%s)
    
    # Compute scores
    scores = reranker.compute_score(pairs, normalize=%s)
    
    # Convert scores to list if it's not already
    if not isinstance(scores, list):
        scores = [float(scores)]
    else:
        scores = [float(score) for score in scores]
    
    # Output ONLY the results as JSON and nothing else
    result = {"scores": scores}
    # Print to stdout and ensure it's flushed immediately
    sys.stdout.write(json.dumps(result) + "\n")
    sys.stdout.flush()
    
    # Don't print anything else after this
    sys.exit(0)
except Exception as e:
    # Restore stderr for error reporting
    sys.stderr = orig_stderr
    result = {"error": str(e)}
    sys.stdout.write(json.dumps(result) + "\n")
    sys.stdout.flush()
    sys.exit(1)
finally:
    # Restore stderr and other hooks
    sys.stderr = orig_stderr
    sys.excepthook = orig_excepthook
`, c.modelName, useFP16Str, normalizeStr)

	// Prepare input data
	inputData := map[string]interface{}{
		"pairs":     pairs,
		"normalize": normalize,
	}
	inputJSON, err := json.Marshal(inputData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling input data: %w", err)
	}

	// Execute the Python script
	cmd := exec.Command("python", "-c", pythonScript)
	cmd.Stdin = strings.NewReader(string(inputJSON))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("error executing Python script: %w, output: %s", err, string(output))
	}

	// Extract just the JSON part from the output
	// First try to find the first valid JSON object
	jsonStart := bytes.Index(output, []byte("{"))
	if jsonStart < 0 {
		return nil, fmt.Errorf("no JSON found in output: %s", string(output))
	}
	
	// Find where the JSON object ends
	bracketCount := 0
	jsonEnd := jsonStart
	for i := jsonStart; i < len(output); i++ {
		if output[i] == '{' {
			bracketCount++
		} else if output[i] == '}' {
			bracketCount--
			if bracketCount == 0 {
				jsonEnd = i + 1
				break
			}
		}
	}
	
	// Extract just the JSON portion
	jsonData := output[jsonStart:jsonEnd]
	
	// Parse the output
	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		// If parsing fails, try to extract the JSON using regex as a fallback
		fmt.Printf("Warning: Initial JSON parsing failed: %v\n", err)
		
		// Try one more approach - find the first line that looks like valid JSON
		lines := bytes.Split(output, []byte("\n"))
		for _, line := range lines {
			trimmed := bytes.TrimSpace(line)
			if len(trimmed) > 0 && trimmed[0] == '{' && trimmed[len(trimmed)-1] == '}' {
				if err := json.Unmarshal(trimmed, &result); err == nil {
					// Successfully parsed this line as JSON
					fmt.Println("Successfully parsed JSON from output line")
					break
				}
			}
		}
		
		// If we still don't have a result, return the error
		if result == nil {
			return nil, fmt.Errorf("error parsing Python output: %w, output: %s", err, string(output))
		}
	}

	// Check for error
	if errorMsg, ok := result["error"].(string); ok {
		return nil, fmt.Errorf("Python script error: %s", errorMsg)
	}

	// Extract scores
	scoresInterface, ok := result["scores"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid scores format in Python output")
	}

	scores := make([]float64, len(scoresInterface))
	for i, s := range scoresInterface {
		scores[i], ok = s.(float64)
		if !ok {
			return nil, fmt.Errorf("invalid score type at index %d", i)
		}
	}

	return scores, nil
}

// CheckDependencies checks if FlagEmbedding is installed
func (c *BGERerankerClient) CheckDependencies() error {
	checkScript := `
import importlib.util
import sys

# Check if FlagEmbedding is installed
flag_spec = importlib.util.find_spec("FlagEmbedding")
if flag_spec is None:
    print("not_installed")
    sys.exit(0)
else:
    print("installed")
    sys.exit(0)
`
	cmd := exec.Command("python", "-c", checkScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error checking Python dependencies: %w", err)
	}

	result := strings.TrimSpace(string(output))
	if result == "not_installed" {
		return fmt.Errorf("FlagEmbedding library is not installed. Run 'rlama install-dependencies' to install it")
	}

	return nil
}

// CheckModelExists verifies that the model exists and is accessible
func (c *BGERerankerClient) CheckModelExists() error {
	pythonScript := `
import sys
import os
import json
import warnings

# Maximum warning suppression
warnings.filterwarnings("ignore")
os.environ["TOKENIZERS_PARALLELISM"] = "false"
os.environ["TF_CPP_MIN_LOG_LEVEL"] = "3"
os.environ["TRANSFORMERS_VERBOSITY"] = "error"
os.environ["PYTHONWARNINGS"] = "ignore"

# Specific patch for protobuf/MessageFactory errors
# Monkey patch the protobuf module to prevent MessageFactory errors
# This must be done before importing transformers or FlagEmbedding
try:
    import builtins
    original_import = builtins.__import__
    
    def custom_import(name, *args, **kwargs):
        module = original_import(name, *args, **kwargs)
        # Patch the protobuf module to suppress MessageFactory errors
        if name == 'google.protobuf.descriptor' or name.endswith('.descriptor'):
            if hasattr(module, 'MessageFactory'):
                original_get_prototype = getattr(module.MessageFactory, 'GetPrototype', None)
                if original_get_prototype:
                    def silent_get_prototype(*args, **kwargs):
                        try:
                            return original_get_prototype(*args, **kwargs)
                        except AttributeError:
                            return None
                    module.MessageFactory.GetPrototype = silent_get_prototype
        return module
    
    builtins.__import__ = custom_import
except Exception:
    pass  # If patching fails, continue with standard error suppression

# Redirect stderr to prevent protobuf errors
orig_stderr = sys.stderr
sys.stderr = open(os.devnull, 'w')

# Also patch sys.excepthook to avoid printing uncaught exceptions to stderr
orig_excepthook = sys.excepthook
def silent_excepthook(exctype, value, traceback):
    if orig_stderr != sys.stderr:  # Only silence if we're already redirecting stderr
        return
    # Still log critical errors to original stderr
    if exctype in (SystemExit, KeyboardInterrupt, MemoryError):
        orig_excepthook(exctype, value, traceback)
sys.excepthook = silent_excepthook

try:
    from FlagEmbedding import FlagReranker
    
    # Just initialize the model to check if it exists
    model_name = "BAAI/bge-reranker-v2-m3"
    reranker = FlagReranker(model_name, use_fp16=True)
    
    # Output ONLY the results as JSON and nothing else
    result = {"success": True}
    sys.stdout.write(json.dumps(result) + "\n")
    sys.stdout.flush()
    sys.exit(0)
except Exception as e:
    # Restore stderr for error reporting
    sys.stderr = orig_stderr
    result = {"error": str(e)}
    sys.stdout.write(json.dumps(result) + "\n")
    sys.stdout.flush()
    sys.exit(1)
finally:
    # Restore stderr and other hooks
    sys.stderr = orig_stderr
    sys.excepthook = orig_excepthook
`
	cmd := exec.Command("python", "-c", pythonScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing Python script: %w, output: %s", err, string(output))
	}

	// Find the first line that starts with '{' (likely our JSON)
	lines := bytes.Split(output, []byte("\n"))
	var jsonLine []byte
	for _, line := range lines {
		trimmed := bytes.TrimSpace(line)
		if len(trimmed) > 0 && trimmed[0] == '{' {
			jsonLine = trimmed
			break
		}
	}

	// If no JSON line found, try to extract it from the whole output
	if jsonLine == nil {
		jsonStart := bytes.Index(output, []byte("{"))
		if jsonStart >= 0 {
			// Find where the JSON object ends
			bracketCount := 0
			jsonEnd := jsonStart
			for i := jsonStart; i < len(output); i++ {
				if output[i] == '{' {
					bracketCount++
				} else if output[i] == '}' {
					bracketCount--
					if bracketCount == 0 {
						jsonEnd = i + 1
						break
					}
				}
			}
			jsonLine = output[jsonStart:jsonEnd]
		}
	}

	// If we still couldn't find JSON data, return error
	if jsonLine == nil {
		return fmt.Errorf("no JSON found in output: %s", string(output))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonLine, &result); err != nil {
		return fmt.Errorf("error parsing Python output: %w, output: %s", err, string(output))
	}

	if _, ok := result["error"]; ok {
		return fmt.Errorf("model check failed: %s", result["error"])
	}

	return nil
}
