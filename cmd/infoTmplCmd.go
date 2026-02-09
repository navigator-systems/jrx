package cmd

import (
	"fmt"

	"github.com/navigator-systems/jrx/internal/config"
	"github.com/navigator-systems/jrx/internal/templates"
)

func TmplInfoCmd(version string) {
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

	fmt.Println("Template version is: ", version)
	// Load templates
	if err := tm.LoadTemplates(version); err != nil {
		fmt.Printf("Error loading templates: %v\n", err)
		return
	}

	// List all templates
	tmplList, err := tm.ListAll()
	if err != nil {
		fmt.Printf("Error listing templates: %v\n", err)
		return
	}

	fmt.Println("Available templates:")
	for _, tmpl := range tmplList {
		fmt.Printf("\nName: %s\n", tmpl.Name)
		fmt.Printf("  Path: %s\n", tmpl.Path)
		fmt.Printf("  Description: %s\n", tmpl.Description)
		if len(tmpl.Tags) > 0 {
			fmt.Printf("  Tags: %v\n", tmpl.Tags)
		}
		if len(tmpl.Variables) > 0 {
			fmt.Println("  Variables:")
			for _, v := range tmpl.Variables {
				fmt.Printf("    - %s: %s (default: '%s')\n", v.Key, v.Description, v.Default)
			}
		}
	}
	fmt.Println("\nUse 'jrx project new <project_name> <template_name>' to create a new project from a template.")
}
