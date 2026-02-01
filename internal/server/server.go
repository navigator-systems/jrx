package server

import (
	"log"
	"net/http"

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
	// Download/Initialize templates on server start
	log.Println("Initializing templates...")
	if err := s.templateManager.Initialize(); err != nil {
		log.Printf("Warning: Could not initialize templates: %v\n", err)
	} else {
		log.Println("Templates initialized successfully")
	}

	// Load templates
	if err := s.templateManager.LoadTemplates(); err != nil {
		log.Printf("Warning: Could not load templates: %v\n", err)
	} else {
		log.Println("Templates loaded successfully")
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
	mux.HandleFunc("/github-orgs", s.handleGithubOrgs)
	mux.HandleFunc("/project", s.handleNewProject) // Keep /project as alias for backwards compatibility

	mux.HandleFunc("/", s.handleIndex)

	log.Printf("Server starting on http://localhost:%s\n", s.port)
	return http.ListenAndServe(":"+s.port, mux)
}
