package client

// LLMClient is a common interface for language model clients
type LLMClient interface {
	GenerateCompletion(model, prompt string) (string, error)
	GenerateEmbedding(model, text string) ([]float32, error)
	CheckLLMAndModel(modelName string) error
}

// Adapt existing methods of OllamaClient to implement LLMClient
func (c *OllamaClient) CheckLLMAndModel(modelName string) error {
	return c.CheckOllamaAndModel(modelName)
}

// Adapt OpenAIClient methods to implement LLMClient
func (c *OpenAIClient) CheckLLMAndModel(modelName string) error {
	return c.CheckOpenAIAndModel(modelName)
}

// IsOpenAIModel checks if a model is an OpenAI model
func IsOpenAIModel(modelName string) bool {
	// OpenAI models generally start with "gpt-" or "text-"
	openAIModels := []string{
		"gpt-3.5-turbo", "gpt-4", "gpt-4-turbo", "gpt-4o",
		"text-davinci", "text-curie", "text-babbage", "text-ada",
	}

	for _, prefix := range openAIModels {
		if modelName == prefix || StartsWith(modelName, prefix+"-") {
			return true
		}
	}

	return false
}

// StartsWith checks if a string starts with a prefix
func StartsWith(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return s[:len(prefix)] == prefix
}

// GetLLMClient returns the appropriate client based on the model
func GetLLMClient(modelName string, ollamaClient *OllamaClient) LLMClient {
	if IsOpenAIModel(modelName) {
		return NewOpenAIClient()
	}
	return ollamaClient
}

// GetLLMClientWithProfile returns the appropriate client based on profile or model
func GetLLMClientWithProfile(modelName, profileName string, ollamaClient *OllamaClient) (LLMClient, error) {
	// If a profile is specified, use it
	if profileName != "" {
		return GetLLMClientFromProfile(profileName)
	}

	// Otherwise fall back to model-based selection
	if IsOpenAIModel(modelName) {
		return NewOpenAIClient(), nil
	}
	return ollamaClient, nil
}

// GetLLMClientFromProfile returns a client based on the specified profile
func GetLLMClientFromProfile(profileName string) (LLMClient, error) {
	client, err := NewOpenAIClientWithProfile(profileName)
	if err != nil {
		return nil, err
	}
	return client, nil
}
