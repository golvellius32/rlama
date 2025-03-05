package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rlama",
	Short: "RLAMA est un outil CLI pour créer et utiliser des systèmes RAG avec Ollama",
	Long: `RLAMA (Retrieval-Augmented Language Model Adapter) est un outil en ligne de commande 
qui simplifie la création et l'utilisation de systèmes RAG (Retrieval-Augmented Generation) 
basés sur les modèles Ollama.`,
}

// Execute exécute la commande racine
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Ici vous pouvez définir des flags globaux
} 