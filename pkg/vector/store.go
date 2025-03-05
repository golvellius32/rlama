package vector

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
)

// VectorItem représente un élément dans le stockage vectoriel
type VectorItem struct {
	ID      string    `json:"id"`
	Vector  []float32 `json:"vector"`
}

// SearchResult représente un résultat de recherche
type SearchResult struct {
	ID       string  `json:"id"`
	Score    float64 `json:"score"`
}

// Store est un stockage simple de vecteurs avec recherche par similarité cosinus
type Store struct {
	Items []VectorItem `json:"items"`
}

// NewStore crée un nouveau stockage de vecteurs
func NewStore() *Store {
	return &Store{
		Items: []VectorItem{},
	}
}

// Add ajoute un vecteur au stockage
func (s *Store) Add(id string, vector []float32) {
	// Vérifier si l'ID existe déjà
	for i, item := range s.Items {
		if item.ID == id {
			// Remplacer le vecteur existant
			s.Items[i].Vector = vector
			return
		}
	}
	
	// Ajouter un nouveau vecteur
	s.Items = append(s.Items, VectorItem{
		ID:     id,
		Vector: vector,
	})
}

// Search recherche les vecteurs les plus similaires
func (s *Store) Search(query []float32, limit int) []SearchResult {
	var results []SearchResult
	
	// Calculer la similarité cosinus pour chaque vecteur
	for _, item := range s.Items {
		score := cosineSimilarity(query, item.Vector)
		results = append(results, SearchResult{
			ID:    item.ID,
			Score: score,
		})
	}
	
	// Trier par score décroissant
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	
	// Limiter le nombre de résultats
	if limit > 0 && limit < len(results) {
		results = results[:limit]
	}
	
	return results
}

// cosineSimilarity calcule la similarité cosinus entre deux vecteurs
func cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}
	
	var dotProduct float64
	var normA, normB float64
	
	for i := range a {
		dotProduct += float64(a[i] * b[i])
		normA += float64(a[i] * a[i])
		normB += float64(b[i] * b[i])
	}
	
	if normA == 0 || normB == 0 {
		return 0
	}
	
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

// Save sauvegarde le stockage de vecteurs dans un fichier
func (s *Store) Save(path string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("impossible de sérialiser le stockage de vecteurs: %w", err)
	}
	
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("impossible de sauvegarder le stockage de vecteurs: %w", err)
	}
	
	return nil
}

// Load charge le stockage de vecteurs depuis un fichier
func (s *Store) Load(path string) error {
	// Vérifier si le fichier existe
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Le fichier n'existe pas, utiliser un stockage vide
			s.Items = []VectorItem{}
			return nil
		}
		return fmt.Errorf("impossible d'accéder au fichier de stockage de vecteurs: %w", err)
	}
	
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("impossible de lire le fichier de stockage de vecteurs: %w", err)
	}
	
	err = json.Unmarshal(data, s)
	if err != nil {
		return fmt.Errorf("impossible de désérialiser le stockage de vecteurs: %w", err)
	}
	
	return nil
} 