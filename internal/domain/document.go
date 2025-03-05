package domain

import (
	"path/filepath"
	"strings"
	"time"
)

// Document représente un document indexé dans le système RAG
type Document struct {
	ID          string    `json:"id"`
	Path        string    `json:"path"`
	Name        string    `json:"name"`
	Content     string    `json:"content"`
	Embedding   []float32 `json:"-"` // Ne pas sérialiser en JSON
	CreatedAt   time.Time `json:"created_at"`
	ContentType string    `json:"content_type"`
	Size        int64     `json:"size"`
}

// NewDocument crée une nouvelle instance de Document
func NewDocument(path string, content string) *Document {
	return &Document{
		ID:          filepath.Base(path),
		Path:        path,
		Name:        filepath.Base(path),
		Content:     content,
		Embedding:   nil,
		CreatedAt:   time.Now(),
		ContentType: guessContentType(path),
		Size:        int64(len(content)),
	}
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
	default:
		return "application/octet-stream"
	}
} 