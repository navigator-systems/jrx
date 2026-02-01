package server

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// createZipArchive creates a ZIP archive of the specified directory
func (s *Server) createZipArchive(sourceDir, zipPath string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk through the directory and add files to ZIP
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the .git directory if it exists
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		// Get the relative path for the ZIP entry
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// Create ZIP entry header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath
		header.Method = zip.Deflate

		// Handle directories
		if info.IsDir() {
			header.Name += "/"
			_, err := zipWriter.CreateHeader(header)
			return err
		}

		// Write file to ZIP
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
}

// parseVars parses variables from a string format "key1=value1,key2=value2"
func parseVars(varsString string) map[string]string {
	vars := make(map[string]string)
	if varsString == "" {
		return vars
	}

	pairs := strings.Split(varsString, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.Trim(strings.TrimSpace(kv[1]), "\"'")
			vars[key] = value
		}
	}
	return vars
}
