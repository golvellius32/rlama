package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/repository"
)

var forceDelete bool

var deleteCmd = &cobra.Command{
	Use:   "delete [nom-rag]",
	Short: "Supprime un système RAG",
	Long:  `Supprime définitivement un système RAG et tous ses documents indexés.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]
		repo := repository.NewRagRepository()

		// Vérifier si le RAG existe
		if !repo.Exists(ragName) {
			return fmt.Errorf("le système RAG '%s' n'existe pas", ragName)
		}

		// Demander confirmation sauf si --force est spécifié
		if !forceDelete {
			fmt.Printf("Êtes-vous sûr de vouloir supprimer définitivement le système RAG '%s'? (o/n): ", ragName)
			var response string
			fmt.Scanln(&response)
			
			response = strings.ToLower(strings.TrimSpace(response))
			if response != "o" && response != "oui" {
				fmt.Println("Suppression annulée.")
				return nil
			}
		}

		// Supprimer le RAG
		err := repo.Delete(ragName)
		if err != nil {
			return err
		}

		fmt.Printf("Le système RAG '%s' a été supprimé avec succès.\n", ragName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "Supprimer sans demander de confirmation")
} 