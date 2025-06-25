package templategit

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func CloneTemplate(repoURL string, branchName string) error {
	if _, err := os.Stat("templates"); err == nil {
		log.Println("Templates directory already exists, removing...")
		if err := os.RemoveAll("templates"); err != nil {
			log.Fatalf("Error removing templates directory: %v", err)
			return err
		}

	}

	_, err := git.PlainClone("templates", false, &git.CloneOptions{
		URL:           repoURL,
		ReferenceName: plumbing.NewBranchReferenceName(branchName),
		SingleBranch:  true,
		Depth:         1,
	})
	if err != nil {
		log.Printf("Failed to clone template repository: %v\n", err)
		return err
	}

	log.Printf("Cloning template from '%s' into templates...\n", repoURL)

	log.Println("Template cloned successfully!")
	return nil
}

func GetTemplateCtrl(filename string) (TemplateFile, error) {
	log.Println("Loading templates...")
	if _, err := os.Stat("templates/templates.toml"); err != nil {
		log.Println("Template file not found...")
	}
	var templateFile TemplateFile
	file := filepath.Join("templates", filename)
	if _, err := toml.DecodeFile(file, &templateFile); err != nil {
		log.Fatal("Error decoding template file:", err)
		return TemplateFile{}, err
	}
	return templateFile, nil
}
