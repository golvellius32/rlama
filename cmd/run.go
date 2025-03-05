package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/service"
)

var runCmd = &cobra.Command{
	Use:   "run [nom-rag]",
	Short: "Exécute un système RAG",
	Long: `Exécute un système RAG précédemment créé. 
Démarre une session interactive pour interagir avec le système RAG.
Exemple: rlama run rag1`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ragName := args[0]

		ragService := service.NewRagService()
		rag, err := ragService.LoadRag(ragName)
		if err != nil {
			return err
		}

		fmt.Printf("RAG '%s' chargé. Modèle: %s\n", rag.Name, rag.ModelName)
		fmt.Println("Tapez votre question (ou 'exit' pour quitter):")

		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("> ")
			if !scanner.Scan() {
				break
			}

			question := scanner.Text()
			if question == "exit" {
				break
			}

			if strings.TrimSpace(question) == "" {
				continue
			}

			answer, err := ragService.Query(rag, question)
			if err != nil {
				fmt.Printf("Erreur: %s\n", err)
				continue
			}

			fmt.Println(answer)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
} 