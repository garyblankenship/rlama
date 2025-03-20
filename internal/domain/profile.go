package domain

import (
	"time"
)

// APIProfile représente un profil pour les clés API
type APIProfile struct {
	Name       string    `json:"name"`
	Provider   string    `json:"provider"` // "openai", "anthropic", etc.
	APIKey     string    `json:"api_key"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	LastUsedAt time.Time `json:"last_used_at,omitempty"`
}

// NewAPIProfile crée un nouveau profil API
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