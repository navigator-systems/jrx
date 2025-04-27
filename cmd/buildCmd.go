package cmd

import (
	"bytes"
	"fmt"
	"jrx/ops"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	projectName := filepath.Base(absPath)
	// Prepare the command
	cmd := exec.Command("go", "mod", "init", projectName)
	cmd.Dir = absPath

	// Capture stdout and stderr
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Run the command
	if err := cmd.Run(); err != nil {
		fmt.Println(stderr.String())
		return err
	}

	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = absPath
	var tidyStderr bytes.Buffer
	tidy.Stderr = &tidyStderr
	if err := tidy.Run(); err != nil {
		fmt.Println(tidyStderr.String())
		return err
	}

	return nil
}

func BuildCmd(path, goarch, goos string) error {
	if path == "" {
		fmt.Println("Please provide a project name")
		return nil
	}

	absPath, _ := filepath.Abs(path)
	// Check if the path exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		fmt.Printf("Project path does not exist: %s\n", absPath)
		return err
	}

	gomod := filepath.Join(absPath, "go.mod")
	if _, err := os.Stat(gomod); err == nil {
		fmt.Println("go.mod file already exists")
	} else {
		ModCmd(path)
	}

	if ops.CheckIfCfgExists(absPath) {
		buildProjectConfig(absPath)
	} else {
		buildGoCmd(absPath, "", goos, goarch)
	}

	fmt.Println("Build completed successfully")
	return nil
}

func buildProjectConfig(absPath string) error {
	jrxConfig, err := ops.ReadCfgFile(absPath)
	if err != nil {
		return err
	}
	if len(jrxConfig.Builds) > 0 {
		for _, build := range jrxConfig.Builds {
			buildGoCmd(absPath, build.Flags, build.OS, build.Arch)

		}
	} else {
		fmt.Println("No builds found in config")
	}

	return nil
}

func buildGoCmd(absPath, flags, goos, goarch string) error {

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

	filebin := filepath.Join(absPath, "bin")
	outputPath := filepath.Join(filebin, outputName)
	// Prepare the command
	command := []string{"go", "build", flags, "-o", outputPath, "."}
	if flags == "" {
		command = []string{"go", "build", "-o", outputPath, "."}
	}

	fmt.Println("Executing build:", strings.Join(command, " "))
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = absPath

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%s", goos),
		fmt.Sprintf("GOARCH=%s", goarch))

	// Capture stdout and stderr
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Build failed: %v\n", stderr.String())
		return err
	}

	return nil
}
