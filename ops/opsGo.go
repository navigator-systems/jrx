package ops

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func Version() string {
	cmd := exec.Command("go", "version")
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}
	// Use regular expression to extract the version number
	re := regexp.MustCompile(`go(\d+\.\d+)`)
	match := re.FindStringSubmatch(string(out))

	return match[1]
}

func ReadGoSum(path string) (map[string]string, error) {
	pathSum := filepath.Join(path, "go.sum")
	if _, err := os.Stat(pathSum); os.IsNotExist(err) {
		fmt.Printf("File %s does not exist\n", pathSum)
		return nil, nil
	}
	file, err := os.Open(pathSum)
	if err != nil {
		fmt.Printf("Error opening go.sum: %v\n", err)
		return nil, err
	}
	defer file.Close()

	fmt.Println("Dependencies and versions from go.sum:")
	scanner := bufio.NewScanner(file)
	dependencies := make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			dependency := parts[0]
			version := parts[1]
			dependencies[dependency] = version
			fmt.Printf("Dependency: %s, Version: %s\n", dependency, version)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading go.sum: %v\n", err)
		return nil, err
	}
	return dependencies, nil
}
