package cmd

import (
	"fmt"
	"log"

	"github.com/navigator-systems/jrx/internal/config"
	"github.com/navigator-systems/jrx/internal/errors"
	"github.com/navigator-systems/jrx/internal/generator"
	"github.com/navigator-systems/jrx/internal/templates"
)

func NewCmd(projectName, templateName, gitOrg string) {
	// Validate input
	if projectName == "" {
		fmt.Println("Error:", errors.ErrEmptyProjectName)
		return
	}
	if templateName == "" {
		fmt.Println("Error:", errors.ErrEmptyTemplateName)
		return
	}

	// Load JRX configuration
	jrxConfig, err := config.ReadJRXConfig()
	if err != nil {
		fmt.Printf("Error reading JRX config: %v\n", err)
		return
	}

	// Create template manager
	tm := templates.NewTemplateManager(jrxConfig)

	// Load templates
	if err := tm.LoadTemplates(); err != nil {
		fmt.Printf("Error loading templates: %v\n", err)
		return
	}

	// Get the specific template
	tmpl, err := tm.GetTemplate(templateName)
	if err != nil {
		if err.Error() == "get template: template not found" {
			fmt.Printf("Template '%s' not found\n", templateName)
			return
		}
		fmt.Printf("Error getting template: %v\n", err)
		return
	}

	// Create project generator
	pg := generator.NewProjectGenerator(tmpl, projectName, tm.GetTemplatesDir(), tm.GetFuncMap())
	if gitOrg != "" {
		pg.SetGitOrg(gitOrg)
	}

	// Generate the project
	if err := pg.Generate(); err != nil {
		fmt.Printf("Error generating project: %v\n", err)
		return
	}

	fmt.Printf("✓ Project '%s' created successfully from template '%s'\n", projectName, templateName)
	log.Printf("Project directory: %s\n", pg.GetOutputDir())
}
