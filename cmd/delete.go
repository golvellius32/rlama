package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/repository"
)

var forceDelete bool

var deleteCmd = &cobra.Command{
	Use:   "delete [rag-name]",
	Short: "Delete a RAG system",
	Long:  `Permanently delete a RAG system and all its indexed documents.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]
		repo := repository.NewRagRepository()

		// Check if the RAG exists
		if !repo.Exists(ragName) {
			return fmt.Errorf("the RAG system '%s' does not exist", ragName)
		}

		// Ask for confirmation unless --force is specified
		if !forceDelete {
			fmt.Printf("Are you sure you want to permanently delete the RAG system '%s'? (y/n): ", ragName)
			var response string
			fmt.Scanln(&response)
			
			response = strings.ToLower(strings.TrimSpace(response))
			if response != "y" && response != "yes" {
				fmt.Println("Deletion cancelled.")
				return nil
			}
		}

		// Delete the RAG
		err := repo.Delete(ragName)
		if err != nil {
			return err
		}

		fmt.Printf("The RAG system '%s' has been successfully deleted.\n", ragName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "Delete without asking for confirmation")
} 