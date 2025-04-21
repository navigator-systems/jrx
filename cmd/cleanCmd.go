package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func CleanCmd(path string) error {

	if path == "" {
		fmt.Println("Please provide a project name")
		return nil
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Printf("Failed to resolve absolute path: %v\n", err)
		return err
	}
	// Get project name from path
	projectName := filepath.Base(absPath)

	fmt.Printf("Cleaning project: %s\n", projectName)
	binFolder := filepath.Join(absPath, "bin")

	// Read directory contents
	entries, err := os.ReadDir(binFolder)
	if err != nil {
		fmt.Printf("Failed to read directory: %v\n", err)
		return err
	}

	// Delete all files in the directory
	for _, entry := range entries {
		entryPath := filepath.Join(binFolder, entry.Name())
		if entry.IsDir() {
			fmt.Printf("Skipping directory: %s\n", entryPath)
			continue
		}

		err := os.Remove(entryPath)
		if err != nil {
			fmt.Printf("Failed to delete file %s: %v\n", entryPath, err)
			return err
		}
		fmt.Printf("Deleted file: %s\n", entryPath)
	}

	fmt.Println("Cleaned project successfully.")
	return nil
}
