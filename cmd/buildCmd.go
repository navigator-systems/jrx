package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func ModCmd(path string) error {

	if path == "" {
		fmt.Println("Please provide a project name")
		return nil
	}

	fmt.Println("Performing go mod")

	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("Failed to resolve absolute path: %v\n", err)
		return err
	}

	// Prepare the command
	cmd := exec.Command("go", "mod", "init", path)
	cmd.Dir = absPath

	// Capture stdout and stderr
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		fmt.Println(stderr.String())
		return err
	}

	fmt.Println("Mod completed successfully")
	return nil
}

func BuildCmd(path, goarch, goos string) error {
	if path == "" {
		fmt.Println("Please provide a project name")
		return nil
	}

	fmt.Println("Building go project")

	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("Failed to resolve absolute path: %v\n", err)
		return err
	}
	// Get project name from path
	projectName := filepath.Base(absPath)
	outputName := projectName

	// If GOOS or GOARCH are specified, customize the output name
	if goos != "" || goarch != "" {
		outputName = fmt.Sprintf("%s-%s-%s", projectName, goos, goarch)
		if goos == "windows" {
			outputName += ".exe"
		}
	}

	gomod := filepath.Join(absPath, "go.mod")
	if _, err := os.Stat(gomod); err == nil {
		fmt.Println("go.mod file already exists")
	} else {
		ModCmd(path)
	}

	filebin := filepath.Join(absPath, "bin")
	outputPath := filepath.Join(filebin, outputName)
	// Prepare the command

	cmd := exec.Command("go", "build", "-o", outputPath, ".")
	cmd.Dir = absPath
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%s", goos),
		fmt.Sprintf("GOARCH=%s", goarch))

	// Capture stdout and stderr
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		fmt.Printf("Build failed: %v\n", stderr.String())
		return err
	}

	fmt.Println("Build completed successfully")
	return nil

}
