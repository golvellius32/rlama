package main

import (
	"fmt"
	"os"

	"github.com/yourusername/rlama/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Erreur: %s\n", err)
		os.Exit(1)
	}
} 