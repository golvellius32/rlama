package cmd

import (
	"fmt"
	// Supprimez ou commentez ces deux lignes si elles ne sont pas utilisées
	// "fmt"
	// "os"

	"github.com/spf13/cobra"
)

const (
	Version = "0.1.0" // Version actuelle de RLAMA
)

var rootCmd = &cobra.Command{
	Use:   "rlama",
	Short: "RLAMA est un outil CLI pour créer et utiliser des systèmes RAG avec Ollama",
	Long: `RLAMA (Retrieval-Augmented Language Model Adapter) est un outil en ligne de commande 
qui simplifie la création et l'utilisation de systèmes RAG (Retrieval-Augmented Generation) 
basés sur les modèles Ollama.

Commandes principales:
  rag [modèle] [nom-rag] [chemin-dossier]    Crée un nouveau système RAG
  run [nom-rag]                              Exécute un système RAG existant
  list                                       Liste tous les systèmes RAG disponibles
  delete [nom-rag]                           Supprime un système RAG
  update                                     Vérifie et installe les mises à jour de RLAMA`,
}

// Variable pour stocker le flag de version
var versionFlag bool

// Execute exécute la commande racine
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Ajout du flag --version
	rootCmd.Flags().BoolVarP(&versionFlag, "version", "v", false, "Affiche la version de RLAMA")
	
	// Override de la fonction Run pour gérer le flag --version
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		if versionFlag {
			fmt.Printf("RLAMA version %s\n", Version)
			return
		}
		
		// Si aucun argument n'est fourni et --version n'est pas utilisé, afficher l'aide
		if len(args) == 0 {
			cmd.Help()
		}
	}
} 