package main

import (
	"log"

	"github.com/golvellius32/rlama/api" // Update with your module name
)

func main() {
	router := api.SetupRouter()
	log.Println("Starting RLAMA API server on http://localhost:3001")
	router.Run(":3001")
}
