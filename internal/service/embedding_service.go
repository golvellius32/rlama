package service

import (
	"fmt"

	"github.com/golvellius32/rlama/internal/client"
	"github.com/golvellius32/rlama/internal/domain"
)

// EmbeddingService manages the generation of embeddings for documents
type EmbeddingService struct {
	ollamaClient *client.OllamaClient
}

// NewEmbeddingService creates a new instance of EmbeddingService
func NewEmbeddingService() *EmbeddingService {
	return &EmbeddingService{
		ollamaClient: client.NewOllamaClient(),
	}
}

// GenerateEmbeddings generates embeddings for a list of documents
func (es *EmbeddingService) GenerateEmbeddings(docs []*domain.Document, modelName string) error {
	for _, doc := range docs {
		// We can chunk here if needed

		// Generate embedding
		embedding, err := es.ollamaClient.GenerateEmbedding(modelName, doc.Content)
		if err != nil {
			return fmt.Errorf("error generating embedding for %s: %w", doc.Path, err)
		}

		doc.Embedding = embedding
	}

	return nil
}

// GenerateQueryEmbedding generates an embedding for a query
func (es *EmbeddingService) GenerateQueryEmbedding(query string, modelName string) ([]float32, error) {
	embedding, err := es.ollamaClient.GenerateEmbedding(modelName, query)
	if err != nil {
		return nil, fmt.Errorf("error generating embedding for query: %w", err)
	}

	return embedding, nil
}
