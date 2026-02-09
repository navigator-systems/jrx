package templates

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/navigator-systems/jrx/internal/config"
	"github.com/navigator-systems/jrx/internal/errors"
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

func CreateCacheDir(versionDir string) error {
	log.Printf("Creating cache directory at %s\n", versionDir)
	if err := os.MkdirAll(versionDir, 0755); err != nil {

		return errors.NewError("", errors.ErrCannotCreateDirectory)
	}
	return nil
}

func CreateDirsIfNotExist(parentDir string, dirs []string) error {
	for _, dir := range dirs {
		fullPath := filepath.Join(parentDir, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return errors.NewError(fmt.Sprintf("create directory %s", fullPath), errors.ErrCannotCreateDirectory)
		}
	}
	return nil
}
