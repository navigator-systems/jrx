package generator

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/navigator-systems/jrx/internal/adapters/scm"
	"github.com/navigator-systems/jrx/internal/errors"
	"github.com/navigator-systems/jrx/internal/templates"
)

// ProjectGenerator handles project generation from templates
type ProjectGenerator struct {
	template     *templates.RootTemplate
	projectName  string
	outputDir    string
	gitOrg       string
	templatesDir string
	funcMap      template.FuncMap
}

// NewProjectGenerator creates a new ProjectGenerator instance
func NewProjectGenerator(tmpl *templates.RootTemplate, projectName string, templatesDir string, funcMap template.FuncMap) *ProjectGenerator {
	return &ProjectGenerator{
		template:     tmpl,
		projectName:  projectName,
		outputDir:    projectName,
		templatesDir: templatesDir,
		funcMap:      funcMap,
	}
}

// SetGitOrg sets the Git organization for the project
func (pg *ProjectGenerator) SetGitOrg(gitOrg string) {
	pg.gitOrg = gitOrg
}

// SetOutputDir sets a custom output directory (defaults to project name)
func (pg *ProjectGenerator) SetOutputDir(dir string) {
	pg.outputDir = dir
}

// Generate creates the project from the template
func (pg *ProjectGenerator) Generate() error {
	// Validate project
	if err := pg.validateProject(); err != nil {
		return err
	}

	// Copy and process template files
	if err := pg.copyFiles(); err != nil {
		return err
	}

	// Initialize Git repository
	if err := pg.initializeGit(); err != nil {
		return err
	}

	log.Printf("Successfully created project '%s' from template '%s'\n", pg.projectName, pg.template.Name)
	return nil
}

// validateProject checks if the project can be created
func (pg *ProjectGenerator) validateProject() error {
	if pg.projectName == "" {
		return errors.NewError("validate project", errors.ErrEmptyProjectName)
	}

	if pg.template == nil {
		return errors.NewError("validate project", errors.ErrInvalidTemplate)
	}

	// Check if template path exists
	templatePath := pg.template.GetFullPath(pg.templatesDir)
	if _, err := os.Stat(templatePath); err != nil {
		return errors.NewError("validate project", errors.ErrTemplatePathMissing)
	}

	// Check if project directory already exists
	if _, err := os.Stat(pg.outputDir); err == nil {
		return errors.NewError("validate project", errors.ErrProjectExists)
	}

	return nil
}

// copyFiles walks through the template directory and processes all files
func (pg *ProjectGenerator) copyFiles() error {
	templatePath := pg.template.GetFullPath(pg.templatesDir)

	// Set the project name in the template
	pg.template.ProjectName = pg.projectName

	err := filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get the relative path to maintain directory structure
		relPath, err := filepath.Rel(templatePath, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Create the destination path
		destPath := filepath.Join(pg.outputDir, relPath)
		if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Parse and execute the template file
		tmpl, err := template.New(info.Name()).Funcs(pg.funcMap).ParseFiles(path)
		if err != nil {
			return fmt.Errorf("error parsing template file %s: %w", path, err)
		}

		// Create destination file
		dstFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("failed to create destination file: %w", err)
		}
		defer dstFile.Close()

		// Execute template
		if err := tmpl.Execute(dstFile, pg.template); err != nil {
			return fmt.Errorf("error executing template for %s: %w", relPath, err)
		}

		log.Printf("Rendered: %s\n", relPath)
		return nil
	})

	if err != nil {
		return errors.NewError("copy template files", err)
	}

	return nil
}

// initializeGit initializes a Git repository for the project
func (pg *ProjectGenerator) initializeGit() error {
	log.Printf("Initializing Git repository for %s...\n", pg.projectName)

	scm.GitInit(pg.outputDir)
	scm.GitBranchMain(pg.outputDir)
	scm.GitAddCommmit(pg.outputDir)

	log.Println("Git repository initialized successfully")
	return nil
}

// GetOutputDir returns the output directory path
func (pg *ProjectGenerator) GetOutputDir() string {
	return pg.outputDir
}

// GetProjectName returns the project name
func (pg *ProjectGenerator) GetProjectName() string {
	return pg.projectName
}
