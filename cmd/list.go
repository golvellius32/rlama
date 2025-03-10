package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/golvellius32/rlama/internal/repository"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available RAG systems",
	Long:  `Display a list of all RAG systems that have been created.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := repository.NewRagRepository()
		ragNames, err := repo.ListAll()
		if err != nil {
			return err
		}

		if len(ragNames) == 0 {
			fmt.Println("No RAG systems found.")
			return nil
		}

		fmt.Printf("Available RAG systems (%d found):\n\n", len(ragNames))

		// Use tabwriter for aligned display
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "NAME\tMODEL\tCREATED ON\tDOCUMENTS")

		for _, name := range ragNames {
			rag, err := repo.Load(name)
			if err != nil {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", name, "error", "error", "error")
				continue
			}

			// Format the date
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
