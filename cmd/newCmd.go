package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/navigator-systems/jrx/internal/config"
	"github.com/navigator-systems/jrx/internal/errors"
	"github.com/navigator-systems/jrx/internal/generator"
	"github.com/navigator-systems/jrx/internal/templates"
)

func parseVars(varsString string) map[string]string {
	vars := make(map[string]string)
	if varsString == "" {
		return vars
	}
	// Split por comas
	pairs := strings.Split(varsString, ",")
	for _, pair := range pairs {
		// Split por =
		kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.Trim(strings.TrimSpace(kv[1]), "\"'") // Remove quotes
			vars[key] = value
		}
	}
	return vars
}

func NewCmd(projectName, templateName, varsString, githubOrg, version string) {
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

	if version == "" {
		version = jrxConfig.TemplatesDefault
	}

	// Load templates
	if err := tm.LoadTemplates(version); err != nil {
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

	// Parse Variables
	userVars := parseVars(varsString)

	if len(userVars) > 0 {
		//
		for i := range tmpl.Variables {
			fmt.Println(i)
			if userValue, exists := userVars[tmpl.Variables[i].Key]; exists {
				tmpl.Variables[i].Default = userValue
				log.Printf("Variable '%s' set to: %s\n", tmpl.Variables[i].Key, userValue)
			}
		}
	}

	// Create project generator
	pg := generator.NewProjectGenerator(
		tmpl, projectName, tm.GetTemplatesDir(), version, tm.GetFuncMap(), jrxConfig)

	// Generate the project
	if err := pg.Generate(); err != nil {
		fmt.Printf("Error generating project: %v\n", err)
		return
	}

	fmt.Printf("Project '%s' created successfully from template '%s'\n", projectName, templateName)
	log.Printf("Project directory: %s\n", pg.GetOutputDir())

	if githubOrg != "" {
		// Create GitHub repository
		ctx := context.Background()
		if err := pg.CreateAndPushToGitHub(ctx, githubOrg); err != nil {
			fmt.Printf("Warning: Failed to create/push GitHub repository: %v\n", err)
			fmt.Printf("Project was created locally. You can push manually:\n")
			fmt.Printf("  cd %s\n", pg.GetOutputDir())
			fmt.Printf("  git remote add origin <repo-url>\n")
			fmt.Printf("  git push -u origin main\n")
			return
		}
	}
}
