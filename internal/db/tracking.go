package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// TrackingInput contains the data required to persist a newly created project.
// This keeps call sites simple while still allowing rich metadata to be stored.
type TrackingInput struct {
	ProjectName     string
	TemplateName    string
	TemplateVersion string
	RepositoryURL   string

	Team     string
	Tags     []string
	Metadata map[string]interface{}
}

// TrackProjectCreation persists the project and a matching history event.
// It is safe to call from both CLI and web handlers and remains idempotent by project name.
func TrackProjectCreation(ctx context.Context, database Database, input TrackingInput) error {
	if database == nil {
		return nil
	}

	repo := database.Project()
	if repo == nil {
		return fmt.Errorf("project repository is not available")
	}

	// If the project already exists, we do not duplicate it.
	_, err := repo.GetByName(ctx, input.ProjectName)
	if err == nil {
		return nil
	}
	if !isNotFoundError(err) {
		return fmt.Errorf("failed to check existing project: %w", err)
	}

	now := time.Now().UTC()
	project := &Project{
		ID:              generateID(),
		Name:            input.ProjectName,
		TemplateName:    input.TemplateName,
		TemplateVersion: input.TemplateVersion,
		RepositoryURL:   input.RepositoryURL,
		CreatedAt:       now,
		UpdatedAt:       now,

		Status:   "active",
		Tags:     input.Tags,
		Metadata: input.Metadata,
	}

	if err := repo.Create(ctx, project); err != nil {
		return fmt.Errorf("failed to create project tracking record: %w", err)
	}

	history := &ProjectHistory{
		ID:        generateID(),
		ProjectID: project.ID,
		Action:    "created",
		Timestamp: now,
		Details: map[string]interface{}{
			"template_name":    input.TemplateName,
			"template_version": input.TemplateVersion,
			"repository_url":   input.RepositoryURL,
		},
	}

	if err := repo.AddHistory(ctx, history); err != nil {
		return fmt.Errorf("failed to create project history record: %w", err)
	}

	return nil
}

// UpdateProjectRepositoryURL updates the project repository URL and stores an audit event.
func UpdateProjectRepositoryURL(ctx context.Context, database Database, projectName, repositoryURL, updatedBy string) error {
	if database == nil {
		return nil
	}

	repo := database.Project()
	if repo == nil {
		return fmt.Errorf("project repository is not available")
	}

	project, err := repo.GetByName(ctx, projectName)
	if err != nil {
		return fmt.Errorf("failed to load project for repository update: %w", err)
	}

	project.RepositoryURL = repositoryURL
	project.UpdatedAt = time.Now().UTC()

	if err := repo.Update(ctx, project); err != nil {
		return fmt.Errorf("failed to update repository url: %w", err)
	}

	history := &ProjectHistory{
		ID:        generateID(),
		ProjectID: project.ID,
		Action:    "updated",
		Timestamp: time.Now().UTC(),
		User:      updatedBy,
		Details: map[string]interface{}{
			"repository_url": repositoryURL,
		},
	}

	if err := repo.AddHistory(ctx, history); err != nil {
		return fmt.Errorf("failed to add repository update history: %w", err)
	}

	return nil
}

func isNotFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
