package cmd

import (
	"fmt"
	"log"

	"github.com/navigator-systems/jrx/internal/config"
	"github.com/navigator-systems/jrx/internal/server"
)

func ServerCmd(port string) {
	// Load JRX configuration
	jrxConfig, err := config.ReadJRXConfig()
	if err != nil {
		fmt.Printf("Error reading JRX config: %v\n", err)
		return
	}

	// Default port if not specified
	if port == "" {
		port = "8080"
	}

	// Create and start server
	srv := server.NewServer(jrxConfig, port)
	log.Printf("Starting JRX server on port %s...\n", port)

	if err := srv.Start(); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
