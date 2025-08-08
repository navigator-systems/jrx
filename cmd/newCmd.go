package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/navigator-systems/jrx/giteng"
	"github.com/navigator-systems/jrx/patterns"
)

var TemplateFuncs = template.FuncMap{
	"index":     patterns.Index,
	"join":      strings.Join,
	"toLower":   strings.ToLower,
	"toUpper":   strings.ToUpper,
	"hasPrefix": strings.HasPrefix,
	"hasSuffix": strings.HasSuffix,
}

func NewCmd(ProjectName, templateName, gitOrg string) {
	if ProjectName == "" || templateName == "" {
		fmt.Println("Please provide a name for the project and a template name")
		return
	}
	jrxGit, err := patterns.GetTemplateCtrl()
	if err != nil {
		fmt.Println("Error getting templates:", err)
		return
	}
	if _, ok := jrxGit.Templates[templateName]; !ok {
		fmt.Printf("Template '%s' not found\n", templateName)
		return
	}

	project := jrxGit.Templates[templateName]
	project.Name = ProjectName

	templatePath := project.Path
	templatePath = "jrxTemplates/" + templatePath
	if _, err := os.Stat(templatePath); err != nil {
		fmt.Println("Template path does not exist")
		return
	}
	if _, err := os.Stat(ProjectName); err == nil {
		fmt.Printf("Project '%s' directory already exists\n", ProjectName)
		return
	}

	err = filepath.Walk(templatePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		//Get the relative path to maintain directory structure
		relPath, err := filepath.Rel(templatePath, path)
		if err != nil {
			return err
		}
		//Create the destination path
		destPath := filepath.Join(ProjectName, relPath)
		if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
			return err
		}

		tmpl, err := template.New(info.Name()).Funcs(TemplateFuncs).ParseFiles(path)
		//tmpl, err := template.ParseFiles(path)
		if err != nil {
			return fmt.Errorf("error parsing template file %s: %v", path, err)
		}

		dstFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		if err := tmpl.Execute(dstFile, project); err != nil {
			return fmt.Errorf("error executing template for %s: %v", relPath, err)
		}

		fmt.Println("Rendered:", relPath)
		return nil
	})

	if err != nil {
		fmt.Printf("Error copying template files: %v\n", err)
		return
	}

	giteng.GitInit(ProjectName)
	giteng.GitBranchMain(ProjectName)
	giteng.GitAddCommmit(ProjectName)

	fmt.Printf("Creating new project '%s' with template '%s'\n", ProjectName, templateName)

}
