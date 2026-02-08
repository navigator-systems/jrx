package templates

import (
	"path/filepath"

	"github.com/navigator-systems/jrx/internal/errors"
)

// Root: Template definition
type RootTemplate struct {
	ProjectName string
	Name        string   `toml:"name"`
	Description string   `toml:"description"`
	Path        string   `toml:"path"`
	Tags        []string `toml:"tags"`
	ProjectInfo ProjectTemplate
	Variables   []VariablesTemplate `toml:"variables"`
}

// Validate checks if the template has all required fields
func (rt *RootTemplate) Validate() error {
	if rt.Name == "" {
		return errors.NewError("validate template", errors.ErrInvalidTemplate)
	}
	if rt.Path == "" {
		return errors.NewError("validate template", errors.ErrInvalidTemplate)
	}
	return nil
}

// HasTag checks if the template has a specific tag
func (rt *RootTemplate) HasTag(tag string) bool {
	for _, t := range rt.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

// GetFullPath returns the full path to the template directory
func (rt *RootTemplate) GetFullPath(baseDir, version string) string {
	fullPath := filepath.Join(baseDir, version, rt.Path)

	return fullPath
}

// GetVariableWithFallback returns the variable value or a fallback if not found
func (rt *RootTemplate) GetVariableWithFallback(key, fallback string) string {
	if val := rt.GetVariable(key); val != "" {
		return val
	}
	return fallback
}

type ProjectTemplate struct {
	Language        string `toml:"language"`
	LanguageVersion string `toml:"language_version,omitempty"`
	Entry           string `toml:"entry"`
	AppVersion      string `toml:"appversion,omitempty"`
}

// Metadata fields used to substitute inside template files.
type VariablesTemplate struct {
	Key         string `toml:"key"`
	Description string `toml:"description"`
	Default     string `toml:"default,omitempty"`
}

type TemplateFile struct {
	Templates map[string]RootTemplate `toml:"templates"`
}
