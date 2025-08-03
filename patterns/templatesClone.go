package patterns

import (
	"log"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func CloneTemplate(repoURL, branchName, key, passphrase string) error {

	if _, err := os.Stat("jrxTemplates"); err == nil {
		if err := os.RemoveAll("jrxTemplates"); err != nil {
			log.Fatalf("Error removing templates directory: %v", err)
			return err
		}
	}
	publicKeys, err := ssh.NewPublicKeysFromFile("git", key, passphrase)
	if err != nil {
		log.Fatalf("Error creating public keys: %v", err)
		return err
	}

	_, err = git.PlainClone("jrxTemplates", false, &git.CloneOptions{
		URL:           repoURL,
		ReferenceName: plumbing.NewBranchReferenceName(branchName),
		SingleBranch:  true,
		Auth:          publicKeys,
		Depth:         1,
	})

	if err != nil {
		log.Printf("Failed to clone template repository: %v\n", err)
		return err
	}
	log.Printf("Cloning template from '%s' into templates...\n", repoURL)

	return nil

}

func GetTemplateCtrl() (TemplateFile, error) {

	log.Println("Loading templates...")

	if _, err := os.Stat("jrxTemplates/templates.toml"); err != nil {

		log.Println("Template file not found...")

	}

	var templateFile TemplateFile

	file := filepath.Join("jrxTemplates", "templates.toml")

	if _, err := toml.DecodeFile(file, &templateFile); err != nil {

		log.Fatal("Error decoding template file:", err)

		return TemplateFile{}, err

	}

	return templateFile, nil

}
