package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/service"
)

var ragCmd = &cobra.Command{
	Use:   "rag [modèle] [nom-rag] [chemin-dossier]",
	Short: "Crée un nouveau système RAG",
	Long: `Crée un nouveau système RAG en indexant tous les documents du dossier spécifié.
Exemple: rlama rag llama3.2 rag1 ./documents

Le dossier sera créé s'il n'existe pas encore. 
Les formats supportés incluent: .txt, .md, .html, .json, .csv, et divers fichiers de code source.`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		modelName := args[0]
		ragName := args[1]
		folderPath := args[2]

		// Afficher un message pour indiquer que le processus a commencé
		fmt.Printf("Création du RAG '%s' avec le modèle '%s' à partir du dossier '%s'...\n", 
			ragName, modelName, folderPath)

		ragService := service.NewRagService()
		err := ragService.CreateRag(modelName, ragName, folderPath)
		if err != nil {
			return err
		}

		fmt.Printf("RAG '%s' créé avec succès.\n", ragName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(ragCmd)
} 