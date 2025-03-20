package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dontizi/rlama/internal/domain"
)

// ProfileRepository gère le stockage des profils API
type ProfileRepository struct {
	basePath string
}

// NewProfileRepository crée une nouvelle instance de ProfileRepository
func NewProfileRepository() *ProfileRepository {
	// Utiliser ~/.rlama/profiles comme dossier par défaut
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	
	basePath := filepath.Join(homeDir, ".rlama", "profiles")
	
	// Créer le dossier s'il n'existe pas
	os.MkdirAll(basePath, 0755)
	
	return &ProfileRepository{
		basePath: basePath,
	}
}

// getProfilePath retourne le chemin complet pour un profil donné
func (r *ProfileRepository) getProfilePath(name string) string {
	return filepath.Join(r.basePath, name+".json")
}

// Exists vérifie si un profil existe
func (r *ProfileRepository) Exists(name string) bool {
	_, err := os.Stat(r.getProfilePath(name))
	return err == nil
}

// Save sauvegarde un profil
func (r *ProfileRepository) Save(profile *domain.APIProfile) error {
	profile.UpdatedAt = time.Now()
	
	data, err := json.MarshalIndent(profile, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling profile: %w", err)
	}
	
	err = os.WriteFile(r.getProfilePath(profile.Name), data, 0644)
	if err != nil {
		return fmt.Errorf("error writing profile file: %w", err)
	}
	
	return nil
}

// Load charge un profil
func (r *ProfileRepository) Load(name string) (*domain.APIProfile, error) {
	data, err := os.ReadFile(r.getProfilePath(name))
	if err != nil {
		return nil, fmt.Errorf("error reading profile '%s': %w", name, err)
	}
	
	var profile domain.APIProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, fmt.Errorf("error unmarshaling profile '%s': %w", name, err)
	}
	
	return &profile, nil
}

// Delete supprime un profil
func (r *ProfileRepository) Delete(name string) error {
	if !r.Exists(name) {
		return fmt.Errorf("profile '%s' does not exist", name)
	}
	
	err := os.Remove(r.getProfilePath(name))
	if err != nil {
		return fmt.Errorf("error deleting profile '%s': %w", name, err)
	}
	
	return nil
}

// ListAll retourne la liste de tous les profils
func (r *ProfileRepository) ListAll() ([]string, error) {
	files, err := os.ReadDir(r.basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("error reading profiles directory: %w", err)
	}
	
	var profileNames []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".json" {
			// Enlever l'extension .json
			name := file.Name()[:len(file.Name())-5]
			profileNames = append(profileNames, name)
		}
	}
	
	return profileNames, nil
} 