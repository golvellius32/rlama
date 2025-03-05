package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/pkg/vector"
)

// RagRepository gère la persistance des systèmes RAG
type RagRepository struct {
	basePath string
}

// NewRagRepository crée une nouvelle instance de RagRepository
func NewRagRepository() *RagRepository {
	// Utilisez ~/.rlama comme dossier de données par défaut
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	
	basePath := filepath.Join(homeDir, ".rlama")
	
	// Créer le dossier s'il n'existe pas
	os.MkdirAll(basePath, 0755)
	
	return &RagRepository{
		basePath: basePath,
	}
}

// getRagPath renvoie le chemin complet pour un RAG donné
func (r *RagRepository) getRagPath(ragName string) string {
	return filepath.Join(r.basePath, ragName)
}

// getRagInfoPath renvoie le chemin du fichier d'information RAG
func (r *RagRepository) getRagInfoPath(ragName string) string {
	return filepath.Join(r.getRagPath(ragName), "info.json")
}

// getRagVectorStorePath renvoie le chemin du fichier de stockage des vecteurs
func (r *RagRepository) getRagVectorStorePath(ragName string) string {
	return filepath.Join(r.getRagPath(ragName), "vectors.json")
}

// Exists vérifie si un RAG existe
func (r *RagRepository) Exists(ragName string) bool {
	_, err := os.Stat(r.getRagInfoPath(ragName))
	return err == nil
}

// Save sauvegarde un système RAG
func (r *RagRepository) Save(rag *domain.RagSystem) error {
	ragPath := r.getRagPath(rag.Name)
	
	// Créer le dossier pour ce RAG
	err := os.MkdirAll(ragPath, 0755)
	if err != nil {
		return fmt.Errorf("impossible de créer le dossier pour le RAG: %w", err)
	}
	
	// Sauvegarder les informations du RAG
	ragInfo := *rag // Copie pour éviter de modifier l'original
	
	// Sérialiser et sauvegarder le fichier info.json
	infoJSON, err := json.MarshalIndent(ragInfo, "", "  ")
	if err != nil {
		return fmt.Errorf("impossible de sérialiser les informations du RAG: %w", err)
	}
	
	err = os.WriteFile(r.getRagInfoPath(rag.Name), infoJSON, 0644)
	if err != nil {
		return fmt.Errorf("impossible de sauvegarder les informations du RAG: %w", err)
	}
	
	// Sauvegarder le Vector Store
	err = rag.VectorStore.Save(r.getRagVectorStorePath(rag.Name))
	if err != nil {
		return fmt.Errorf("impossible de sauvegarder le Vector Store: %w", err)
	}
	
	return nil
}

// Load charge un système RAG
func (r *RagRepository) Load(ragName string) (*domain.RagSystem, error) {
	// Vérifier si le RAG existe
	if !r.Exists(ragName) {
		return nil, fmt.Errorf("le RAG '%s' n'existe pas", ragName)
	}
	
	// Charger les informations du RAG
	infoBytes, err := os.ReadFile(r.getRagInfoPath(ragName))
	if err != nil {
		return nil, fmt.Errorf("impossible de lire les informations du RAG: %w", err)
	}
	
	var ragInfo domain.RagSystem
	err = json.Unmarshal(infoBytes, &ragInfo)
	if err != nil {
		return nil, fmt.Errorf("impossible de désérialiser les informations du RAG: %w", err)
	}
	
	// Créer un nouveau Vector Store et le charger à partir du fichier
	ragInfo.VectorStore = vector.NewStore()
	err = ragInfo.VectorStore.Load(r.getRagVectorStorePath(ragName))
	if err != nil {
		return nil, fmt.Errorf("impossible de charger le Vector Store: %w", err)
	}
	
	return &ragInfo, nil
} 