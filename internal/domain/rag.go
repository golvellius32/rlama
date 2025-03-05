package domain

import (
	"time"

	"github.com/dontizi/rlama/pkg/vector"
)

// RagSystem représente un système RAG complet
type RagSystem struct {
	Name        string    `json:"name"`
	ModelName   string    `json:"model_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Description string    `json:"description"`
	VectorStore *vector.Store
	Documents   []*Document `json:"documents"`
}

// NewRagSystem crée une nouvelle instance de RagSystem
func NewRagSystem(name, modelName string) *RagSystem {
	now := time.Now()
	return &RagSystem{
		Name:        name,
		ModelName:   modelName,
		CreatedAt:   now,
		UpdatedAt:   now,
		VectorStore: vector.NewStore(),
		Documents:   []*Document{},
	}
}

// AddDocument ajoute un document au système RAG
func (r *RagSystem) AddDocument(doc *Document) {
	r.Documents = append(r.Documents, doc)
	if doc.Embedding != nil {
		r.VectorStore.Add(doc.ID, doc.Embedding)
	}
	r.UpdatedAt = time.Now()
}

// GetDocumentByID récupère un document par son ID
func (r *RagSystem) GetDocumentByID(id string) *Document {
	for _, doc := range r.Documents {
		if doc.ID == id {
			return doc
		}
	}
	return nil
} 