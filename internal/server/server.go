package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/navigator-systems/jrx/internal/config"
	"github.com/navigator-systems/jrx/internal/templates"
)

// Server represents the web server
type Server struct {
	config          config.JRXConfig
	templateManager *templates.TemplateManager
	port            string
}

// NewServer creates a new server instance
func NewServer(cfg config.JRXConfig) *Server {
	return &Server{
		config:          cfg,
		templateManager: templates.NewTemplateManager(cfg),
		port:            cfg.ServerPort,
	}
}

// Start initializes and starts the web server
func (s *Server) Start() error {
	// Load templates
	if err := s.templateManager.LoadTemplates(); err != nil {
		log.Printf("Warning: Could not load templates: %v\n", err)
	}

	// Create a new ServeMux to properly handle routes
	mux := http.NewServeMux()

	// Serve static files (images, css, js, etc.)
	fs := http.FileServer(http.Dir("images"))
	mux.Handle("/images/", http.StripPrefix("/images/", fs))

	// Setup routes - specific routes must be registered before the root
	mux.HandleFunc("/templates/download", s.handleDownloadTemplates)
	mux.HandleFunc("/templates", s.handleTemplates)
	mux.HandleFunc("/new-project", s.handleNewProject)
	mux.HandleFunc("/api/templates", s.handleAPITemplates)
	mux.HandleFunc("/", s.handleHome)

	log.Printf("Server starting on http://localhost:%s\n", s.port)
	return http.ListenAndServe(":"+s.port, mux)
}

// handleHome handles the root endpoint
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
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

	tmpl, err := template.ParseFiles("static/new-project.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading template: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title     string
		Templates []templates.RootTemplate
		Error     string
	}{
		Title: "Create New Project",
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

// handleCreateProject handles the project creation form submission
func (s *Server) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	projectName := strings.TrimSpace(r.FormValue("projectName"))
	templateName := strings.TrimSpace(r.FormValue("templateName"))

	// Parse variables
	varsString := r.FormValue("variables")
	vars := parseVars(varsString)

	tmpl, err := template.ParseFiles("static/project-result.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error loading template: %v", err), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title        string
		ProjectName  string
		TemplateName string
		Variables    map[string]string
		Success      bool
		Message      string
	}{
		Title:        "Project Creation Result",
		ProjectName:  projectName,
		TemplateName: templateName,
		Variables:    vars,
		Success:      true,
		Message:      "Project Created",
	}

	// Validation
	if projectName == "" {
		data.Success = false
		data.Message = "Error: Project name is required"
	} else if templateName == "" {
		data.Success = false
		data.Message = "Error: Template name is required"
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handleAPITemplates returns templates as JSON
func (s *Server) handleAPITemplates(w http.ResponseWriter, r *http.Request) {
	if !s.templateManager.IsLoaded() {
		http.Error(w, "Templates not loaded", http.StatusServiceUnavailable)
		return
	}

	tmplList, err := s.templateManager.ListAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tmplList)
}

// parseVars parses variables from a string format "key1=value1,key2=value2"
func parseVars(varsString string) map[string]string {
	vars := make(map[string]string)
	if varsString == "" {
		return vars
	}

	pairs := strings.Split(varsString, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.Trim(strings.TrimSpace(kv[1]), "\"'")
			vars[key] = value
		}
	}
	return vars
}
