package templates

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/navigator-systems/jrx/internal/config"
	"github.com/navigator-systems/jrx/internal/errors"
)

// TemplateManager manages template operations
type TemplateManager struct {
	config         config.JRXConfig
	templateFile   TemplateFile
	funcMap        template.FuncMap
	loaded         bool
	currentVersion string
}

// NewTemplateManager creates a new TemplateManager instance
func NewTemplateManager(cfg config.JRXConfig) *TemplateManager {
	return &TemplateManager{
		config:  cfg,
		funcMap: buildFuncMap(),
		loaded:  false,
	}
}

// buildFuncMap creates the function map for template execution
func buildFuncMap() template.FuncMap {
	return template.FuncMap{
		"index": Index,
		"getVariable": func(key string, rt *RootTemplate) string {
			return rt.GetVariable(key)
		},
		"join":      strings.Join,
		"toLower":   strings.ToLower,
		"toUpper":   strings.ToUpper,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix,
	}
}

// GetFuncMap returns the template function map
func (tm *TemplateManager) GetFuncMap() template.FuncMap {
	return tm.funcMap
}

// Initialize clones the template repository
func (tm *TemplateManager) Initialize() error {
	log.Println("Downloading templates...")

	// Remove existing templates directory if it exists
	if _, err := os.Stat(tm.config.TemplatesCacheDir); err == nil {
		if err := os.RemoveAll(tm.config.TemplatesCacheDir); err != nil {
			return errors.NewError("remove templates directory", err)
		}
	}

	//Create cache directory if it doesn't exist
	if err := CreateCacheDir(tm.config.TemplatesCacheDir); err != nil {
		return errors.NewError("create cache directory", err)
	}

	// Create directories for branch versions
	if err := CreateDirsIfNotExist(tm.config.TemplatesCacheDir, tm.config.TemplatesBranch); err != nil {
		return errors.NewError("create version directories", err)
	}

	// Get versions of tags to clone based on pattern and max versions
	tagVersions, err := tm.GetVersionsTags()
	if err != nil {
		return errors.NewError("get tag versions", err)
	}
	//Create directories for tag versions
	if err := CreateDirsIfNotExist(tm.config.TemplatesCacheDir, tagVersions); err != nil {
		return errors.NewError("create version directories", err)
	}

	// Setup SSH authentication
	publicKeys, err := ssh.NewPublicKeysFromFile("git", tm.config.SshKeyPath, tm.config.SshKeyPassphrase)
	if err != nil {
		return errors.NewError("create SSH keys", err)
	}

	// Clone the repository for each branch
	for _, branch := range tm.config.TemplatesBranch {
		repoDest := filepath.Join(tm.config.TemplatesCacheDir, branch)
		_, err = git.PlainClone(repoDest, false, &git.CloneOptions{
			URL:           tm.config.TemplatesRepo,
			ReferenceName: plumbing.NewBranchReferenceName(branch),
			SingleBranch:  true,
			Auth:          publicKeys,
			Depth:         1,
		})
	}

	// Clone the repository for each tag
	for _, tag := range tagVersions {
		repoDest := filepath.Join(tm.config.TemplatesCacheDir, tag)
		_, err = git.PlainClone(repoDest, false, &git.CloneOptions{
			URL:           tm.config.TemplatesRepo,
			ReferenceName: plumbing.NewTagReferenceName(tag),
			SingleBranch:  true,
			Auth:          publicKeys,
			Depth:         1,
		})
	}

	if err != nil {
		return errors.NewError("clone template repository", err)
	}

	log.Printf("Successfully cloned templates from '%s'\n", tm.config.TemplatesRepo)
	return nil
}

// LoadTemplates loads all templates from the templates directory
func (tm *TemplateManager) LoadTemplates(templatesVersion string) error {
	log.Println("Loading templates...")

	if templatesVersion == "" {
		templatesVersion = tm.config.TemplatesDefault
	}
	log.Println("Template version is:", templatesVersion)

	if !tm.ValidateVersion(templatesVersion) {
		return errors.NewError("Load templates", fmt.Errorf("Version %s is not available", templatesVersion))
	}

	templatePath := filepath.Join(tm.config.TemplatesCacheDir, templatesVersion, "templates.toml")
	if _, err := os.Stat(templatePath); err != nil {
		return errors.NewError("find templates.toml", errors.ErrConfigNotFound)
	}

	// Decode the main template file
	if _, err := toml.DecodeFile(templatePath, &tm.templateFile); err != nil {
		return errors.NewError("decode templates.toml", err)
	}

	// Process each template to load additional configuration files
	for templateKey, tpl := range tm.templateFile.Templates {
		baseDir := filepath.Join(tm.config.TemplatesCacheDir, templatesVersion, tpl.Path)

		// Load project.toml if it exists
		projectPath := filepath.Join(baseDir, "project.toml")
		if _, err := os.Stat(projectPath); err == nil {
			if err := tm.loadProjectConfig(templateKey, projectPath, &tpl); err != nil {
				log.Printf("Warning: could not load project.toml for %s: %v", templateKey, err)
			}
		}

		// Load vars.toml if it exists
		varsPath := filepath.Join(baseDir, "vars.toml")
		if _, err := os.Stat(varsPath); err == nil {
			if err := tm.loadVarsConfig(templateKey, varsPath, &tpl); err != nil {
				log.Printf("Warning: could not load vars.toml for %s: %v", templateKey, err)
			}
		}

		// Update the template in the map
		tm.templateFile.Templates[templateKey] = tpl
	}

	tm.loaded = true
	log.Printf("Successfully loaded %d templates\n", len(tm.templateFile.Templates))
	return nil
}

// loadProjectConfig loads and decodes project.toml into the template's ProjectInfo
func (tm *TemplateManager) loadProjectConfig(templateKey, filePath string, tpl *RootTemplate) error {
	var projectInfo ProjectTemplate
	if _, err := toml.DecodeFile(filePath, &projectInfo); err != nil {
		return fmt.Errorf("error decoding project.toml: %w", err)
	}

	tpl.ProjectInfo = projectInfo
	return nil
}

// loadVarsConfig loads and decodes vars.toml into the template's Variables
func (tm *TemplateManager) loadVarsConfig(templateKey, filePath string, tpl *RootTemplate) error {
	var varsConfig struct {
		Variable map[string]struct {
			Default     string `toml:"default"`
			Description string `toml:"description"`
		} `toml:"variable"`
	}

	if _, err := toml.DecodeFile(filePath, &varsConfig); err != nil {
		return fmt.Errorf("error decoding vars.toml: %w", err)
	}

	tpl.Variables = make([]VariablesTemplate, 0, len(varsConfig.Variable))
	for key, varInfo := range varsConfig.Variable {
		tpl.Variables = append(tpl.Variables, VariablesTemplate{
			Key:         key,
			Description: varInfo.Description,
			Default:     varInfo.Default,
		})
	}

	return nil
}

// GetTemplate returns a specific template by name
func (tm *TemplateManager) GetTemplate(name string) (*RootTemplate, error) {
	if !tm.loaded {
		return nil, errors.NewError("get template", errors.ErrLoadTemplates)
	}

	tpl, exists := tm.templateFile.Templates[name]
	if !exists {
		return nil, errors.NewError("get template", errors.ErrTemplateNotFound)
	}

	if err := tpl.Validate(); err != nil {
		return nil, err
	}

	return &tpl, nil
}

// GetTemplatesMap returns the templates map
func (tm *TemplateManager) GetTemplatesMap() map[string]RootTemplate {
	return tm.templateFile.Templates
}

// GetTemplatesDir returns the templates directory path
func (tm *TemplateManager) GetTemplatesDir() string {
	return tm.config.TemplatesCacheDir
}

// IsLoaded returns whether templates have been loaded
func (tm *TemplateManager) IsLoaded() bool {
	return tm.loaded
}
