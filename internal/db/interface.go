package db

import (
	"context"
)

type Database interface {
	Connect(ctx context.Context) error
	Close() error
	Ping() error
	InitDB() error
	Project() ProjectRepository
}

type ProjectRepository interface {
	// Create creates a new project record
	Create(ctx context.Context, project *Project) error

	// Get retrieves a project by ID
	Get(ctx context.Context, id string) (*Project, error)

	// GetByName retrieves a project by name
	GetByName(ctx context.Context, name string) (*Project, error)

	// List retrieves all projects with optional filters
	List(ctx context.Context, filter *ProjectFilter) ([]*Project, error)

	// Update updates an existing project
	Update(ctx context.Context, project *Project) error

	// Delete deletes a project by ID
	Delete(ctx context.Context, id string) error

	// AddHistory adds a history record for a project
	AddHistory(ctx context.Context, history *ProjectHistory) error

	// GetHistory retrieves history for a project
	GetHistory(ctx context.Context, projectID string) ([]*ProjectHistory, error)

	// GetStats retrieves project statistics
	GetStats(ctx context.Context) (*ProjectStats, error)
}
