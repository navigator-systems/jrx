package cmd

import (
	"fmt"

	"github.com/navigator-systems/jrx/patterns"
)

func TmplInfoCmd() {
	jrxGit, err := patterns.GetTemplateCtrl()
	if err != nil {
		fmt.Println("Error getting templates:", err)
		return
	}
	fmt.Println("Available templates:")
	for name, tmpl := range jrxGit.Templates {
		fmt.Printf("Name: %s, Path: %s, Description: %s\n", name, tmpl.Path, tmpl.Description)
	}
	fmt.Println("Use 'jrx new <project_name> <template_name>' to create a new project from a template.")

}
