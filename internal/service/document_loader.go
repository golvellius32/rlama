package service

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dontizi/rlama/internal/domain"
)

// DocumentLoader se charge de charger les documents depuis le système de fichiers
type DocumentLoader struct {
	supportedExtensions map[string]bool
}

// NewDocumentLoader crée une nouvelle instance de DocumentLoader
func NewDocumentLoader() *DocumentLoader {
	return &DocumentLoader{
		supportedExtensions: map[string]bool{
			".txt":   true,
			".md":    true,
			".html":  true,
			".htm":   true,
			".json":  true,
			".csv":   true,
			".log":   true,
			".xml":   true,
			".yaml":  true,
			".yml":   true,
			".go":    true,
			".py":    true,
			".js":    true,
			".java":  true,
			".c":     true,
			".cpp":   true,
			".h":     true,
			".rb":    true,
			".php":   true,
			".rs":    true,
			".swift": true,
			".kt":    true,
		},
	}
}

// LoadDocumentsFromFolder charge tous les documents supportés du dossier spécifié
func (dl *DocumentLoader) LoadDocumentsFromFolder(folderPath string) ([]*domain.Document, error) {
	var documents []*domain.Document

	// Vérifie si le dossier existe
	info, err := os.Stat(folderPath)
	if err != nil {
		return nil, fmt.Errorf("impossible d'accéder au dossier: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("le chemin spécifié n'est pas un dossier: %s", folderPath)
	}

	// Parcourt le dossier de manière récursive
	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ignore les dossiers
		if info.IsDir() {
			return nil
		}

		// Vérifie si l'extension est supportée
		ext := strings.ToLower(filepath.Ext(path))
		if !dl.supportedExtensions[ext] {
			return nil
		}

		// Charge le contenu du fichier
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("impossible de lire le fichier %s: %w", path, err)
		}

		// Crée un document
		doc := domain.NewDocument(path, string(content))
		documents = append(documents, doc)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("erreur lors du parcours du dossier: %w", err)
	}

	return documents, nil
} 