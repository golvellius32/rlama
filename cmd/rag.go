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
Exemple: rlama rag llama3.2 rag1 ./documents`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		modelName := args[0]
		ragName := args[1]
		folderPath := args[2]

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