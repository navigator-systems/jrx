package cmd

import (
	"embed"
	"fmt"
	"jrx/ops"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed templates/*
var templates embed.FS

const (
	bin string = "bin"
	lib string = "lib"
	doc string = "doc"
)

func createFilesFromTemplate(name, path, templateFile string) error {

	mainPath := filepath.Join(name, path, templateFile)
	templateFilePath := filepath.Join("templates", templateFile)
	tmplBytes, err := templates.ReadFile(templateFilePath)
	if err != nil {
		fmt.Printf("Error reading template file: %v\n", err)
		return err
	}
	tmpl, err := template.New("main").Parse(string(tmplBytes))
	if err != nil {
		fmt.Println("Failed to parse template:", err)
		return err
	}
	file, err := os.Create(mainPath)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return err
	}
	defer file.Close()
	context := struct {
		ProjectName string
		GoVersion   string
	}{
		ProjectName: name,
		GoVersion:   ops.Version(), // Example base image

	}

	tmpl.Execute(file, context)
	//tmpl.Execute(file, struct{ ProjectName string }{ProjectName: name})

	return nil

}

func createTree(name string) error {
	fmt.Println("Creating project tree...")
	directories := []string{bin, lib, doc}
	for _, dir := range directories {
		err := os.MkdirAll(name+"/"+dir, 0755)
		if err != nil {
			fmt.Printf("Error creating directory '%s': %v\n", dir, err)
			return err
		}
	}
	fmt.Println("Project tree created successfully.")
	return nil
}

func createDirectory(name string) error {
	fmt.Println("Creating directory...", name)
	err := os.Mkdir(name, 0755)
	if err != nil {
		if os.IsExist(err) {
			fmt.Printf("Error: Directory '%s' already exists\n", name)
		} else {
			fmt.Printf("Error creating directory '%s': %v\n", name, err)
		}
		return err
	}
	return nil

}

func NewCmd(name string) {
	if name == "" {
		fmt.Println("Please provide a name for the project")
		return
	}
	err := createDirectory(name)
	if err != nil {
		os.Exit(1)
	}

	err = createTree(name)
	if err != nil {
		os.Exit(1)
	}

	// Create main.go file from template
	err = createFilesFromTemplate(name, ".", "main.go")
	if err != nil {
		os.Exit(1)
	}
	// Create Makefile file from template
	err = createFilesFromTemplate(name, ".", "Makefile")
	if err != nil {
		os.Exit(1)
	}

	// Create Dockerfile file from template
	err = createFilesFromTemplate(name, ".", "Dockerfile")
	if err != nil {
		os.Exit(1)
	}

	err = createFilesFromTemplate(name, ".", "jrx.toml")
	if err != nil {
		os.Exit(1)
	}

}
