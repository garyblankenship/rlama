package domain

import (
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Document représente un document indexé dans le système RAG
type Document struct {
	ID          string    `json:"id"`
	Path        string    `json:"path"`
	Name        string    `json:"name"`
	Content     string    `json:"content"`
	Metadata    string    `json:"metadata"`
	Embedding   []float32 `json:"-"` // Ne pas sérialiser en JSON
	CreatedAt   time.Time `json:"created_at"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
	URL         string    `json:"url,omitempty"` // URL source pour les documents web
}

// NewDocument crée une nouvelle instance de Document
func NewDocument(path string, content string) *Document {
	// Nettoyer le contenu extrait
	cleanedContent := cleanExtractedText(content)
	
	return &Document{
		ID:          filepath.Base(path),
		Path:        path,
		Name:        filepath.Base(path),
		Content:     cleanedContent,
		Metadata:    "",
		Embedding:   nil,
		CreatedAt:   time.Now(),
		ContentType: guessContentType(path),
		Size:        int64(len(cleanedContent)),
	}
}

// cleanExtractedText nettoie le texte extrait pour améliorer sa qualité
func cleanExtractedText(text string) string {
	// Remplacer les séquences de caractères non imprimables par des espaces
	re := regexp.MustCompile(`[\x00-\x09\x0B\x0C\x0E-\x1F\x7F]+`)
	text = re.ReplaceAllString(text, " ")
	
	// Remplacer les séquences de plus de 2 sauts de ligne par 2 sauts de ligne
	re = regexp.MustCompile(`\n{3,}`)
	text = re.ReplaceAllString(text, "\n\n")
	
	// Remplacer les séquences de plus de 2 espaces par 1 espace
	re = regexp.MustCompile(`[ \t]{2,}`)
	text = re.ReplaceAllString(text, " ")
	
	// Supprimer les lignes qui ne contiennent que des caractères spéciaux ou des chiffres
	lines := strings.Split(text, "\n")
	var cleanedLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > 0 {
			// Vérifier si la ligne contient au moins quelques lettres
			re = regexp.MustCompile(`[a-zA-Z]{2,}`)
			if re.MatchString(trimmed) || len(trimmed) > 20 {
				cleanedLines = append(cleanedLines, line)
			}
		}
	}
	
	return strings.Join(cleanedLines, "\n")
}

// guessContentType essaie de déterminer le type de contenu basé sur l'extension du fichier
func guessContentType(path string) string {
	ext := filepath.Ext(path)
	switch strings.ToLower(ext) {
	case ".txt":
		return "text/plain"
	case ".md", ".markdown":
		return "text/markdown"
	case ".html", ".htm":
		return "text/html"
	case ".pdf":
		return "application/pdf"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".doc":
		return "application/msword"
	case ".csv":
		return "text/csv"
	case ".json":
		return "application/json"
	case ".pptx":
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case ".ppt":
		return "application/vnd.ms-powerpoint"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".xls":
		return "application/vnd.ms-excel"
	case ".rtf":
		return "application/rtf"
	case ".odt":
		return "application/vnd.oasis.opendocument.text"
	default:
		return "application/octet-stream"
	}
} 