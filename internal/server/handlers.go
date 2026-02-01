package server

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/navigator-systems/jrx/internal/generator"
	"github.com/navigator-systems/jrx/internal/templates"
)

// handleIndex handles the home page request
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("static/index.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading landing page: %v", err), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleTemplates displays the templates page
func (s *Server) handleTemplates(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/templates.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading template: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title     string
		IsLoaded  bool
		Templates []templates.RootTemplate
		Error     string
	}{
		Title:    "JRX Templates",
		IsLoaded: s.templateManager.IsLoaded(),
	}

	if s.templateManager.IsLoaded() {
		tmplList, err := s.templateManager.ListAll()
		if err != nil {
			data.Error = err.Error()
		} else {
			data.Templates = tmplList
		}
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleDownloadTemplates handles template download/initialization
func (s *Server) handleDownloadTemplates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Initialize (download) templates
	if err := s.templateManager.Initialize(); err != nil {
		http.Error(w, fmt.Sprintf("Error downloading templates: %v", err), http.StatusInternalServerError)
		return
	}

	// Reload templates
	if err := s.templateManager.LoadTemplates(); err != nil {
		http.Error(w, fmt.Sprintf("Error loading templates: %v", err), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/templates", http.StatusSeeOther)
}

// handleNewProject displays the new project page
func (s *Server) handleNewProject(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		s.handleCreateProject(w, r)
		return
	}

	tmpl, err := template.ParseFiles("static/project.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading template: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title         string
		Templates     map[string]templates.RootTemplate
		Organizations []string
		Error         string
	}{
		Title:         "Create New Project",
		Organizations: s.config.GitProvider.GithubOrganization,
	}

	if s.templateManager.IsLoaded() {
		tmplList := s.templateManager.GetTemplatesMap()

		data.Templates = tmplList
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleCreateProject handles the project creation form submission
func (s *Server) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	projectName := strings.TrimSpace(r.FormValue("projectName"))
	templateName := strings.TrimSpace(r.FormValue("templateName"))
	githubOrg := strings.TrimSpace(r.FormValue("githubOrg"))

	// Parse variables from form fields (var_keyname)
	vars := make(map[string]string)
	if err := r.ParseForm(); err == nil {
		for key, values := range r.Form {
			if strings.HasPrefix(key, "var_") && len(values) > 0 {
				varName := strings.TrimPrefix(key, "var_")
				vars[varName] = values[0]
			}
		}
	}

	data := struct {
		Title         string
		ProjectName   string
		TemplateName  string
		Variables     map[string]string
		Success       bool
		Message       string
		OutputDir     string
		GithubOrg     string
		GithubRepoURL string
	}{
		Title:        "Project Creation Result",
		ProjectName:  projectName,
		TemplateName: templateName,
		Variables:    vars,
		GithubOrg:    githubOrg,
		Success:      true,
		Message:      "Project created successfully!",
	}

	// Validation
	if projectName == "" {
		data.Success = false
		data.Message = "Error: Project name is required"
	} else if templateName == "" {
		data.Success = false
		data.Message = "Error: Template name is required"
	} else {
		// Create the project
		if err := s.createProject(projectName, templateName, vars, githubOrg, &data, w, r); err != nil {
			data.Success = false
			data.Message = fmt.Sprintf("Error: %v", err)
		} else if githubOrg == "" {
			// If no GitHub org, the project was downloaded as ZIP
			// Response already sent, so return early
			return
		}
	}

	tmpl, err := template.ParseFiles("static/project-result.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading template: %v", err), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleGithubOrgs displays the GitHub organizations page
func (s *Server) handleGithubOrgs(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("static/github-orgs.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading template: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title         string
		Organizations []string
		GithubURL     string
		HasToken      bool
	}{
		Title:         "GitHub Organizations",
		Organizations: s.config.GitProvider.GithubOrganization,
		GithubURL:     s.config.GitProvider.GithubURL,
		HasToken:      s.config.GitProvider.GithubToken != "",
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// createProject creates a project from template with optional GitHub push
func (s *Server) createProject(projectName, templateName string, vars map[string]string, githubOrg string, data *struct {
	Title         string
	ProjectName   string
	TemplateName  string
	Variables     map[string]string
	Success       bool
	Message       string
	OutputDir     string
	GithubOrg     string
	GithubRepoURL string
}, w http.ResponseWriter, r *http.Request) error {
	// Log for debugging
	log.Printf("Creating project: name=%s, template=%s, org=%s\n", projectName, templateName, githubOrg)

	// Verify templates are loaded
	if !s.templateManager.IsLoaded() {
		return fmt.Errorf("templates are not loaded, please wait for server initialization")
	}

	// Get the specific template
	tmpl, err := s.templateManager.GetTemplate(templateName)
	if err != nil {
		// Log available templates for debugging
		availableTemplates, _ := s.templateManager.ListAll()
		log.Printf("Template '%s' not found. Available templates: ", templateName)
		for _, t := range availableTemplates {
			log.Printf("  - %s\n", t.Name)
		}
		return fmt.Errorf("template '%s' not found", templateName)
	}

	// Apply user variables to template
	if len(vars) > 0 {
		for i := range tmpl.Variables {
			if userValue, exists := vars[tmpl.Variables[i].Key]; exists {
				tmpl.Variables[i].Default = userValue
				log.Printf("Variable '%s' set to: %s\n", tmpl.Variables[i].Key, userValue)
			}
		}
	}

	// Create project generator
	pg := generator.NewProjectGenerator(tmpl, projectName, s.templateManager.GetTemplatesDir(), s.templateManager.GetFuncMap(), s.config)

	// Generate the project
	if err := pg.Generate(); err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}

	data.OutputDir = pg.GetOutputDir()
	data.Message = fmt.Sprintf("Project '%s' created successfully at: %s", projectName, pg.GetOutputDir())

	// If no GitHub organization, create ZIP and serve as download
	if githubOrg == "" {
		zipPath := pg.GetOutputDir() + ".zip"
		if err := s.createZipArchive(pg.GetOutputDir(), zipPath); err != nil {
			return fmt.Errorf("failed to create zip archive: %w", err)
		}

		// Serve the ZIP file as download
		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", projectName))
		http.ServeFile(w, r, zipPath)

		// Clean up both the project directory and the ZIP file
		go func() {
			if err := os.RemoveAll(pg.GetOutputDir()); err != nil {
				log.Printf("Warning: Failed to cleanup project directory: %v\n", err)
			}
			if err := os.Remove(zipPath); err != nil {
				log.Printf("Warning: Failed to cleanup zip file: %v\n", err)
			}
			log.Printf("Cleaned up project directory and ZIP file for '%s'\n", projectName)
		}()

		log.Printf("Project '%s' created and downloaded as ZIP\n", projectName)
		return nil
	}

	// If GitHub organization is specified, create repo and push
	if githubOrg != "" {
		ctx := context.Background()
		if err := pg.CreateAndPushToGitHub(ctx, githubOrg); err != nil {
			data.Message = fmt.Sprintf("Project created locally at: %s\nWarning: Failed to push to GitHub: %v", pg.GetOutputDir(), err)
			log.Printf("Failed to create/push GitHub repository: %v\n", err)
		} else {
			// Get the GitHub URL (we need to construct it from the org and project name)
			githubURL := fmt.Sprintf("%s/%s/%s", s.config.GitProvider.GithubURL, githubOrg, projectName)
			data.GithubRepoURL = githubURL

			// Clean up local files since project is now on GitHub
			if err := pg.CleanupLocalFiles(); err != nil {
				log.Printf("Warning: Failed to cleanup local files: %v\n", err)
				data.Message = fmt.Sprintf("Project '%s' created and pushed to GitHub successfully!\nRepository: %s\nWarning: Could not cleanup local files at: %s", projectName, githubURL, pg.GetOutputDir())
			} else {
				data.OutputDir = "" // Clear output dir since files were cleaned up
				data.Message = fmt.Sprintf("Project '%s' created and pushed to GitHub successfully!\nRepository: %s\nLocal files have been cleaned up.", projectName, githubURL)
			}
		}
	}

	log.Printf("Project '%s' created successfully from template '%s'\n", projectName, templateName)
	return nil
}

