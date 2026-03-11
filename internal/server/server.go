package server

import (
	"context"
	"log"
	"net/http"

	"github.com/navigator-systems/jrx/internal/config"
	"github.com/navigator-systems/jrx/internal/db"
	"github.com/navigator-systems/jrx/internal/templates"
)

// Server represents the web server
type Server struct {
	config          config.JRXConfig
	templateManager *templates.TemplateManager
	database        db.Database
	port            string
	currentVersion  string
}

// NewServer creates a new server instance
func NewServer(cfg config.JRXConfig) *Server {
	ctx := context.Background()
	database, err := db.InitDatabase(ctx, cfg)
	if err != nil {
		log.Printf("Warning: project tracking database unavailable: %v", err)
	}

	return &Server{
		config:          cfg,
		templateManager: templates.NewTemplateManager(cfg),
		database:        database,
		port:            cfg.ServerPort,
		currentVersion:  cfg.TemplatesDefault,
	}
}

// Close releases resources owned by the server, including optional database connections.
func (s *Server) Close() error {
	if s.database != nil {
		return s.database.Close()
	}
	return nil
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

	defaultVersion := s.config.TemplatesDefault
	// Load templates
	if err := s.templateManager.LoadTemplates(defaultVersion); err != nil {
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
	//mux.HandleFunc("/templates/switch-version", s.handleSwitchVersion)
	mux.HandleFunc("/templates", s.handleTemplates)

	mux.HandleFunc("/github-orgs", s.handleGithubOrgs)
	mux.HandleFunc("/project", s.handleNewProject) // New project creation

	mux.HandleFunc("/", s.handleIndex)

	log.Printf("Server starting on http://localhost:%s\n", s.port)
	return http.ListenAndServe(":"+s.port, mux)
}
