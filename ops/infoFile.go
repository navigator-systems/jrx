package ops

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

func GetFileInfo(path string) error {
	path = filepath.Join(path, "bin")
	err := filepath.WalkDir(path, func(path string, info fs.DirEntry, err error) error {

		if err != nil {
			return err
		}
		// Skip directories
		if info.IsDir() {
			return nil
		}
		// Print file size in MB
		fileInfo, err := info.Info()
		if err != nil {
			return err
		}

		fmt.Printf("Binary file found: %s\n", path)
		sizeMB := float64(fileInfo.Size()) / (1024 * 1024)
		fmt.Printf("File: %s, Size: %.2f MB\n", path, sizeMB)

		return nil

	})
	if err != nil {
		return err
	}
	return nil
}
