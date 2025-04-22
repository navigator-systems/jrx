package cmd

import (
	"fmt"
	"os"
)

const (
	github = ".github/workflows"
)

func createCITree(name string, templates []string) error {

	for _, dir := range templates {
		err := os.MkdirAll(name+"/"+dir, 0755)
		if err != nil {
			fmt.Printf("Error creating directory '%s': %v\n", dir, err)
			return err
		}
	}
	fmt.Println("Project tree created successfully.")
	return nil
}

func AddCICmd(project, template string) {
	// Add command to initialize CI configuration
	if project == "" || template == "" {
		fmt.Println("Error: Project and template must be specified.")
		return
	}

	fmt.Println("Initializing CI configuration...")
	switch template {
	case "github":
		file := "main.yaml"
		createCITree(project, []string{github})
		createFilesFromTemplate(project, github, file)
	case "jenkins":
		file := "Jenkinsfile"
		createFilesFromTemplate(project, ".", file)
	default:
		fmt.Printf("Error: Unsupported CI template '%s'.\n", template)
	}

}
