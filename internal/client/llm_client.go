package client

// LLMClient est une interface commune pour les clients de modèles de langage
type LLMClient interface {
	GenerateCompletion(model, prompt string) (string, error)
	CheckLLMAndModel(modelName string) error
}

// Adapter les méthodes existantes de OllamaClient pour implémenter LLMClient
func (c *OllamaClient) CheckLLMAndModel(modelName string) error {
	return c.CheckOllamaAndModel(modelName)
}

// Adapter les méthodes de OpenAIClient pour implémenter LLMClient
func (c *OpenAIClient) CheckLLMAndModel(modelName string) error {
	return c.CheckOpenAIAndModel(modelName)
}

// IsOpenAIModel vérifie si un modèle est un modèle OpenAI
func IsOpenAIModel(modelName string) bool {
	// Les modèles OpenAI commencent généralement par "gpt-" ou "text-"
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

// StartsWith vérifie si une chaîne commence par un préfixe
func StartsWith(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return s[:len(prefix)] == prefix
}

// GetLLMClient retourne le client approprié selon le modèle
func GetLLMClient(modelName string, ollamaClient *OllamaClient) LLMClient {
	if IsOpenAIModel(modelName) {
		return NewOpenAIClient()
	}
	return ollamaClient
} 