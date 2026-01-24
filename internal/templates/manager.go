package templates

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/BurntSushi/toml"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/google/go-github/v58/github"
	"github.com/navigator-systems/jrx/internal/adapters/scm"
	"github.com/navigator-systems/jrx/internal/config"
	"github.com/navigator-systems/jrx/internal/errors"
)

// TemplateManager manages template operations
type TemplateManager struct {
	config       config.JRXConfig
	templateFile TemplateFile
	templatesDir string
	funcMap      template.FuncMap
	loaded       bool
}

// NewTemplateManager creates a new TemplateManager instance
func NewTemplateManager(cfg config.JRXConfig) *TemplateManager {
	return &TemplateManager{
		config:       cfg,
		templatesDir: "jrxTemplates",
		funcMap:      buildFuncMap(),
		loaded:       false,
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
	if _, err := os.Stat(tm.templatesDir); err == nil {
		if err := os.RemoveAll(tm.templatesDir); err != nil {
			return errors.NewError("remove templates directory", err)
		}
	}

	// Setup SSH authentication
	publicKeys, err := ssh.NewPublicKeysFromFile("git", tm.config.SshKeyPath, tm.config.SshKeyPassphrase)
	if err != nil {
		return errors.NewError("create SSH keys", err)
	}

	// Clone the repository
	_, err = git.PlainClone(tm.templatesDir, false, &git.CloneOptions{
		URL:           tm.config.TemplatesRepo,
		ReferenceName: plumbing.NewBranchReferenceName(tm.config.TemplatesBranch),
		SingleBranch:  true,
		Auth:          publicKeys,
		Depth:         1,
	})

	if err != nil {
		return errors.NewError("clone template repository", err)
	}

	log.Printf("Successfully cloned templates from '%s'\n", tm.config.TemplatesRepo)
	return nil
}

// LoadTemplates loads all templates from the templates directory
func (tm *TemplateManager) LoadTemplates() error {
	log.Println("Loading templates...")

	templatePath := filepath.Join(tm.templatesDir, "templates.toml")
	if _, err := os.Stat(templatePath); err != nil {
		return errors.NewError("find templates.toml", errors.ErrConfigNotFound)
	}

	// Decode the main template file
	if _, err := toml.DecodeFile(templatePath, &tm.templateFile); err != nil {
		return errors.NewError("decode templates.toml", err)
	}

	// Process each template to load additional configuration files
	for templateKey, tpl := range tm.templateFile.Templates {
		baseDir := filepath.Join(tm.templatesDir, tpl.Path)

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

// ListAll returns all available templates
func (tm *TemplateManager) ListAll() ([]RootTemplate, error) {
	if !tm.loaded {
		return nil, errors.NewError("list templates", errors.ErrLoadTemplates)
	}

	templates := make([]RootTemplate, 0, len(tm.templateFile.Templates))
	for _, tpl := range tm.templateFile.Templates {
		templates = append(templates, tpl)
	}

	return templates, nil
}

// GetTemplatesMap returns the templates map
func (tm *TemplateManager) GetTemplatesMap() map[string]RootTemplate {
	return tm.templateFile.Templates
}

// GetTemplatesDir returns the templates directory path
func (tm *TemplateManager) GetTemplatesDir() string {
	return tm.templatesDir
}

// IsLoaded returns whether templates have been loaded
func (tm *TemplateManager) IsLoaded() bool {
	return tm.loaded
}

// CreateGitHubRepository creates a GitHub repository for a project using a template
func (tm *TemplateManager) CreateGitHubRepository(ctx context.Context, projectName string, tmpl *RootTemplate, githubOrg string) (*github.Repository, error) {
	log.Println("Creating GitHub repository...")

	// Create GitHub client
	ghClient, err := scm.NewGitHubClient(tm.config, githubOrg)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub client: %w", err)
	}

	// Create repository description from template
	description := fmt.Sprintf("Project created from template: %s", tmpl.Name)
	if tmpl.Description != "" {
		description = tmpl.Description
	}

	// Create the repository (private by default)
	repo, err := ghClient.CreateRepository(ctx, projectName, description, true)
	if err != nil {
		return nil, err
	}

	log.Printf("Repository created: %s\n", repo.GetHTMLURL())
	return repo, nil
}
