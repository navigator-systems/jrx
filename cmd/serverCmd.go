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

	// Create and start server
	srv := server.NewServer(jrxConfig)
	defer func() {
		if err := srv.Close(); err != nil {
			log.Printf("Warning: failed to close server resources: %v", err)
		}
	}()

	log.Printf("Starting JRX server on port %s...\n", jrxConfig.ServerPort)

	if err := srv.Start(); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
