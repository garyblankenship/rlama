package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/dontizi/rlama/internal/repository"
)

// OpenAIClient est un client pour l'API OpenAI
type OpenAIClient struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

// OpenAICompletionRequest est la structure pour les requêtes de complétion à OpenAI
type OpenAICompletionRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	Temperature float64         `json:"temperature,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
}

// OpenAIMessage représente un message dans le format attendu par OpenAI
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAICompletionResponse est la structure de la réponse de l'API OpenAI
type OpenAICompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int           `json:"index"`
		Message      OpenAIMessage `json:"message"`
		FinishReason string        `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// NewOpenAIClient crée un nouveau client OpenAI
func NewOpenAIClient() *OpenAIClient {
	// Utiliser la clé API depuis l'environnement
	apiKey := os.Getenv("OPENAI_API_KEY")
	
	return &OpenAIClient{
		BaseURL: "https://api.openai.com/v1",
		APIKey:  apiKey,
		Client:  &http.Client{},
	}
}

// NewOpenAIClientWithProfile crée un nouveau client OpenAI avec un profil spécifique
func NewOpenAIClientWithProfile(profileName string) (*OpenAIClient, error) {
	// Créer un repository de profils
	profileRepo := repository.NewProfileRepository()
	
	// Si aucun profil n'est spécifié, utiliser la variable d'environnement
	if profileName == "" {
		apiKey := os.Getenv("OPENAI_API_KEY")
		return &OpenAIClient{
			BaseURL: "https://api.openai.com/v1",
			APIKey:  apiKey,
			Client:  &http.Client{},
		}, nil
	}
	
	// Charger le profil spécifié
	profile, err := profileRepo.Load(profileName)
	if err != nil {
		return nil, fmt.Errorf("error loading profile '%s': %w", profileName, err)
	}
	
	// Vérifier que c'est un profil OpenAI
	if profile.Provider != "openai" {
		return nil, fmt.Errorf("profile '%s' is not an OpenAI profile", profileName)
	}
	
	// Mettre à jour la date de dernière utilisation
	profile.LastUsedAt = time.Now()
	profileRepo.Save(profile)
	
	return &OpenAIClient{
		BaseURL: "https://api.openai.com/v1",
		APIKey:  profile.APIKey,
		Client:  &http.Client{},
	}, nil
}

// GenerateCompletion génère une réponse à partir d'un prompt avec OpenAI
func (c *OpenAIClient) GenerateCompletion(model, prompt string) (string, error) {
	if c.APIKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}

	// Construire la requête
	reqBody := OpenAICompletionRequest{
		Model: model,
		Messages: []OpenAIMessage{
			{
				Role:    "system",
				Content: "You are a helpful assistant that answers questions based on the provided context.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.7, // Valeur par défaut, peut être configurée
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	// Créer la requête HTTP
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/chat/completions", c.BaseURL), bytes.NewBuffer(reqJSON))
	if err != nil {
		return "", err
	}

	// Ajouter les en-têtes nécessaires
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))

	// Envoyer la requête
	resp, err := c.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Vérifier le code de statut
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to generate completion: %s (status: %d)", string(bodyBytes), resp.StatusCode)
	}

	// Décoder la réponse
	var completionResp OpenAICompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&completionResp); err != nil {
		return "", err
	}

	// Vérifier qu'il y a au moins un choix
	if len(completionResp.Choices) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	// Retourner le contenu de la réponse
	return completionResp.Choices[0].Message.Content, nil
}

// CheckOpenAIAndModel vérifie si OpenAI est accessible et si le modèle est disponible
func (c *OpenAIClient) CheckOpenAIAndModel(modelName string) error {
	if c.APIKey == "" {
		return fmt.Errorf("⚠️ OPENAI_API_KEY environment variable not set.\n" +
			"Please set your OpenAI API key before using OpenAI models.")
	}

	// On pourrait faire une vérification de la validité du modèle ici
	// mais pour l'instant, on suppose que le modèle est valide si la clé API est définie
	
	return nil
} 