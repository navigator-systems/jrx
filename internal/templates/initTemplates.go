package templates

import (
	"fmt"
	"log"

	"github.com/navigator-systems/jrx/internal/config"
)

func InitTemplates() {
	fmt.Println("Downloading templates...")
	fmt.Println("Reading JRX config file...")

	jrxConfig, err := config.ReadJRXConfig()
	if err != nil {
		fmt.Printf("Error reading JRX config: %v\n", err)
		return
	}

	fmt.Println("JRX config loaded successfully")

	// Create template manager
	tm := NewTemplateManager(jrxConfig)

	// Initialize (clone) templates
	if err := tm.Initialize(); err != nil {
		fmt.Printf("Error initializing templates: %v\n", err)
		return
	}

	log.Println("Templates downloaded successfully")
}
