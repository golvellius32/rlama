package service

import (
	"fmt"
	"strings"

	"github.com/dontizi/rlama/internal/client"
	"github.com/dontizi/rlama/internal/domain"
	"github.com/dontizi/rlama/internal/repository"
)

// RagService gère les opérations liées aux systèmes RAG
type RagService struct {
	documentLoader   *DocumentLoader
	embeddingService *EmbeddingService
	ragRepository    *repository.RagRepository
	ollamaClient     *client.OllamaClient
}

// NewRagService crée une nouvelle instance de RagService
func NewRagService() *RagService {
	return &RagService{
		documentLoader:   NewDocumentLoader(),
		embeddingService: NewEmbeddingService(),
		ragRepository:    repository.NewRagRepository(),
		ollamaClient:     client.NewOllamaClient(),
	}
}

// CreateRag crée un nouveau système RAG
func (rs *RagService) CreateRag(modelName, ragName, folderPath string) error {
	// Vérifier si le RAG existe déjà
	if rs.ragRepository.Exists(ragName) {
		return fmt.Errorf("un RAG avec le nom '%s' existe déjà", ragName)
	}

	// Créer le système RAG
	rag := domain.NewRagSystem(ragName, modelName)

	// Charger les documents
	docs, err := rs.documentLoader.LoadDocumentsFromFolder(folderPath)
	if err != nil {
		return fmt.Errorf("erreur lors du chargement des documents: %w", err)
	}

	if len(docs) == 0 {
		return fmt.Errorf("aucun document supporté trouvé dans le dossier %s", folderPath)
	}

	// Générer les embeddings pour tous les documents
	err = rs.embeddingService.GenerateEmbeddings(docs, modelName)
	if err != nil {
		return fmt.Errorf("erreur lors de la génération des embeddings: %w", err)
	}

	// Ajouter les documents au RAG
	for _, doc := range docs {
		rag.AddDocument(doc)
	}

	// Sauvegarder le RAG
	err = rs.ragRepository.Save(rag)
	if err != nil {
		return fmt.Errorf("erreur lors de la sauvegarde du RAG: %w", err)
	}

	return nil
}

// LoadRag charge un système RAG
func (rs *RagService) LoadRag(ragName string) (*domain.RagSystem, error) {
	rag, err := rs.ragRepository.Load(ragName)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du chargement du RAG '%s': %w", ragName, err)
	}

	return rag, nil
}

// Query effectue une requête sur un système RAG
func (rs *RagService) Query(rag *domain.RagSystem, query string) (string, error) {
	// Générer l'embedding pour la requête
	queryEmbedding, err := rs.embeddingService.GenerateQueryEmbedding(query, rag.ModelName)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la génération de l'embedding pour la requête: %w", err)
	}

	// Rechercher les documents les plus pertinents
	results := rag.VectorStore.Search(queryEmbedding, 3) // Top 3 documents

	// Construire le contexte
	var context strings.Builder
	context.WriteString("Informations pertinentes:\n\n")

	for _, result := range results {
		doc := rag.GetDocumentByID(result.ID)
		if doc != nil {
			// Limiter la taille du contenu pour éviter les prompts trop longs
			content := doc.Content
			if len(content) > 1000 {
				content = content[:1000] + "..."
			}
			context.WriteString(fmt.Sprintf("--- Document: %s ---\n%s\n\n", doc.Name, content))
		}
	}

	// Construire le prompt
	prompt := fmt.Sprintf(`Vous êtes un assistant IA utile. Utilisez les informations ci-dessous pour répondre à la question.

%s

Question: %s

Répondez de manière concise en vous basant uniquement sur les informations fournies ci-dessus:`, context.String(), query)

	// Générer la réponse
	response, err := rs.ollamaClient.GenerateCompletion(rag.ModelName, prompt)
	if err != nil {
		return "", fmt.Errorf("erreur lors de la génération de la réponse: %w", err)
	}

	return response, nil
} 