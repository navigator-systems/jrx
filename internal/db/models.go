package db

import (
	"time"
)

// Project represents a project created with JRX
type Project struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	TemplateName    string                 `json:"template_name"`
	TemplateVersion string                 `json:"template_version,omitempty"`
	RepositoryURL   string                 `json:"repository_url,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
	Status          string                 `json:"status"` // active, archived, deprecated
	Tags            []string               `json:"tags,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// ProjectHistory represents a historical event for a project
type ProjectHistory struct {
	ID        string                 `json:"id"`
	ProjectID string                 `json:"project_id"`
	Action    string                 `json:"action"` // created, updated, deleted, archived
	Details   map[string]interface{} `json:"details,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	User      string                 `json:"user,omitempty"`
}

// ProjectStats represents statistics about projects
type ProjectStats struct {
	TotalProjects      int            `json:"total_projects"`
	ActiveProjects     int            `json:"active_projects"`
	ArchivedProjects   int            `json:"archived_projects"`
	ProjectsByTemplate map[string]int `json:"projects_by_template"`
	ProjectsByTeam     map[string]int `json:"projects_by_team"`
	RecentProjects     int            `json:"recent_projects"` // Last 30 days
}

// ProjectFilter represents filters for querying projects
type ProjectFilter struct {
	Status       string
	TemplateName string
	Team         string
	Limit        int
	Offset       int
}
