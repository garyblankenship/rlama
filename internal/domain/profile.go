package domain

import (
	"time"
)

// APIProfile represents a profile for API keys
type APIProfile struct {
	Name       string    `json:"name"`
	Provider   string    `json:"provider"` // "openai", "openai-api", "anthropic", etc.
	APIKey     string    `json:"api_key"`
	BaseURL    string    `json:"base_url,omitempty"` // For custom OpenAI-compatible endpoints
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	LastUsedAt time.Time `json:"last_used_at,omitempty"`
}

// NewAPIProfile creates a new API profile
func NewAPIProfile(name, provider, apiKey string) *APIProfile {
	now := time.Now()
	return &APIProfile{
		Name:      name,
		Provider:  provider,
		APIKey:    apiKey,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
