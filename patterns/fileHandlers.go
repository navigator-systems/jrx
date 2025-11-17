package patterns

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// loadProjectConfig loads and decodes project.toml into the template's ProjectInfo
func loadProjectConfig(templateKey, filePath string, tpl *RootTemplate) error {
	var projectInfo ProjectTemplate
	if _, err := toml.DecodeFile(filePath, &projectInfo); err != nil {
		return fmt.Errorf("error decoding project.toml: %w", err)
	}

	tpl.ProjectInfo = projectInfo

	return nil
}

// loadVarsConfig loads and decodes vars.toml into the template's Variables
func loadVarsConfig(templateKey, filePath string, tpl *RootTemplate) error {
	// Create a temporary struct to decode the vars.toml structure
	var varsConfig struct {
		Templates map[string]struct {
			Variables []VariablesTemplate `toml:"variables"`
		} `toml:"templates"`
	}

	if _, err := toml.DecodeFile(filePath, &varsConfig); err != nil {
		return fmt.Errorf("error decoding vars.toml: %w", err)
	}

	// Extract variables for the specific template
	if templateVars, exists := varsConfig.Templates[templateKey]; exists {
		tpl.Variables = templateVars.Variables
	}

	return nil
}
