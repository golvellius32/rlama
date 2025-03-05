package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	DefaultOllamaURL = "http://localhost:11434"
)

// OllamaClient est un client pour l'API Ollama
type OllamaClient struct {
	BaseURL string
	Client  *http.Client
}

// EmbeddingRequest est la structure de la requête pour l'API /api/embeddings
type EmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// EmbeddingResponse est la structure de la réponse de l'API /api/embeddings
type EmbeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

// GenerationRequest est la structure de la requête pour l'API /api/generate
type GenerationRequest struct {
	Model    string   `json:"model"`
	Prompt   string   `json:"prompt"`
	Context  []int    `json:"context,omitempty"`
	Options  Options  `json:"options,omitempty"`
	Format   string   `json:"format,omitempty"`
	Template string   `json:"template,omitempty"`
	Stream   bool     `json:"stream"`
}

// Options pour l'API generate
type Options struct {
	Temperature float64 `json:"temperature,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
	TopK        int     `json:"top_k,omitempty"`
	NumPredict  int     `json:"num_predict,omitempty"`
}

// GenerationResponse est la structure de la réponse de l'API /api/generate
type GenerationResponse struct {
	Model     string `json:"model"`
	Response  string `json:"response"`
	Context   []int  `json:"context"`
	CreatedAt string `json:"created_at"`
	Done      bool   `json:"done"`
}

// NewOllamaClient crée un nouveau client Ollama
func NewOllamaClient() *OllamaClient {
	return &OllamaClient{
		BaseURL: DefaultOllamaURL,
		Client:  &http.Client{},
	}
}

// GenerateEmbedding génère un embedding pour le texte donné
func (c *OllamaClient) GenerateEmbedding(model, text string) ([]float32, error) {
	reqBody := EmbeddingRequest{
		Model:  model,
		Prompt: text,
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Post(
		fmt.Sprintf("%s/api/embeddings", c.BaseURL),
		"application/json",
		bytes.NewBuffer(reqJSON),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to generate embedding: %s (status: %d)", string(bodyBytes), resp.StatusCode)
	}

	var embeddingResp EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embeddingResp); err != nil {
		return nil, err
	}

	return embeddingResp.Embedding, nil
}

// GenerateCompletion génère une réponse pour le prompt donné
func (c *OllamaClient) GenerateCompletion(model, prompt string) (string, error) {
	reqBody := GenerationRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
		Options: Options{
			Temperature: 0.7,
			TopP:        0.9,
			NumPredict:  1024,
		},
	}

	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	resp, err := c.Client.Post(
		fmt.Sprintf("%s/api/generate", c.BaseURL),
		"application/json",
		bytes.NewBuffer(reqJSON),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to generate completion: %s (status: %d)", string(bodyBytes), resp.StatusCode)
	}

	var genResp GenerationResponse
	if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
		return "", err
	}

	return genResp.Response, nil
} 