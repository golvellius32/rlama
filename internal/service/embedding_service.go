package service

import (
	"fmt"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/domain"
)

// EmbeddingService gère la génération d'embeddings pour les documents
type EmbeddingService struct {
	ollamaClient *client.OllamaClient
}

// NewEmbeddingService crée une nouvelle instance de EmbeddingService
func NewEmbeddingService() *EmbeddingService {
	return &EmbeddingService{
		ollamaClient: client.NewOllamaClient(),
	}
}

// GenerateEmbeddings génère des embeddings pour une liste de documents
func (es *EmbeddingService) GenerateEmbeddings(docs []*domain.Document, modelName string) error {
	for _, doc := range docs {
		// On peut chunker ici si nécessaire

		// Génération de l'embedding
		embedding, err := es.ollamaClient.GenerateEmbedding(modelName, doc.Content)
		if err != nil {
			return fmt.Errorf("erreur lors de la génération de l'embedding pour %s: %w", doc.Path, err)
		}

		doc.Embedding = embedding
	}

	return nil
}

// GenerateQueryEmbedding génère un embedding pour une requête
func (es *EmbeddingService) GenerateQueryEmbedding(query string, modelName string) ([]float32, error) {
	embedding, err := es.ollamaClient.GenerateEmbedding(modelName, query)
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la génération de l'embedding pour la requête: %w", err)
	}

	return embedding, nil
} 