package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/dontizi/rlama/internal/repository"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Liste tous les systèmes RAG disponibles",
	Long:  `Affiche la liste de tous les systèmes RAG qui ont été créés.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := repository.NewRagRepository()
		ragNames, err := repo.ListAll()
		if err != nil {
			return err
		}

		if len(ragNames) == 0 {
			fmt.Println("Aucun système RAG n'a été trouvé.")
			return nil
		}

		fmt.Printf("Systèmes RAG disponibles (%d trouvés):\n\n", len(ragNames))
		
		// Utilisation de tabwriter pour un affichage aligné
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NOM\tMODÈLE\tCRÉÉ LE\tDOCUMENTS")
		
		for _, name := range ragNames {
			rag, err := repo.Load(name)
			if err != nil {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", name, "erreur", "erreur", "erreur")
				continue
			}
			
			// Formater la date
			createdAt := rag.CreatedAt.Format("2006-01-02 15:04:05")
			
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\n", rag.Name, rag.ModelName, createdAt, len(rag.Documents))
		}
		w.Flush()
		
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
} 