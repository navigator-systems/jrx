package db

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/navigator-systems/jrx/internal/errors"
)

// sqlDialect identifies query differences between engines.
type sqlDialect string

const (
	dialectSQLite   sqlDialect = "sqlite"
	dialectPostgres sqlDialect = "postgres"
)

// sqlProjectRepository implements ProjectRepository for SQL engines.
// A single repository is used for both SQLite and PostgreSQL, while placeholder syntax
// is generated dynamically per dialect.
type sqlProjectRepository struct {
	db      *sql.DB
	dialect sqlDialect
}

func newSQLProjectRepository(dbConn *sql.DB, dialect sqlDialect) ProjectRepository {
	return &sqlProjectRepository{db: dbConn, dialect: dialect}
}

func (r *sqlProjectRepository) Create(ctx context.Context, project *Project) error {
	if project == nil {
		return fmt.Errorf("project cannot be nil")
	}

	if project.ID == "" {
		project.ID = generateID()
	}

	now := time.Now().UTC()
	if project.CreatedAt.IsZero() {
		project.CreatedAt = now
	}
	if project.UpdatedAt.IsZero() {
		project.UpdatedAt = now
	}
	if strings.TrimSpace(project.Status) == "" {
		project.Status = "active"
	}

	tagsJSON, err := json.Marshal(project.Tags)
	if err != nil {
		return fmt.Errorf("failed to serialize project tags: %w", err)
	}

	metadataJSON, err := json.Marshal(project.Metadata)
	if err != nil {
		return fmt.Errorf("failed to serialize project metadata: %w", err)
	}

	query := fmt.Sprintf(`
		INSERT INTO projects (
			id, name, template_name, template_version, repository_url,
			created_at, updated_at, created_by, team, status, tags, metadata
		) VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
	`, r.p(1), r.p(2), r.p(3), r.p(4), r.p(5), r.p(6), r.p(7), r.p(8), r.p(9), r.p(10), r.p(11), r.p(12))

	_, err = r.db.ExecContext(
		ctx,
		query,
		project.ID,
		project.Name,
		project.TemplateName,
		project.TemplateVersion,
		project.RepositoryURL,
		project.CreatedAt,
		project.UpdatedAt,
		project.Status,
		string(tagsJSON),
		string(metadataJSON),
	)
	if err != nil {
		return errors.NewError("create project", err)
	}

	return nil
}

func (r *sqlProjectRepository) Get(ctx context.Context, id string) (*Project, error) {
	query := fmt.Sprintf(`
		SELECT id, name, template_name, template_version, repository_url,
		       created_at, updated_at, created_by, team, status, tags, metadata
		FROM projects
		WHERE id = %s
	`, r.p(1))

	row := r.db.QueryRowContext(ctx, query, id)
	return scanProjectRow(row)
}

func (r *sqlProjectRepository) GetByName(ctx context.Context, name string) (*Project, error) {
	query := fmt.Sprintf(`
		SELECT id, name, template_name, template_version, repository_url,
		       created_at, updated_at, created_by, team, status, tags, metadata
		FROM projects
		WHERE name = %s
	`, r.p(1))

	row := r.db.QueryRowContext(ctx, query, name)
	return scanProjectRow(row)
}

func (r *sqlProjectRepository) List(ctx context.Context, filter *ProjectFilter) ([]*Project, error) {
	base := `
		SELECT id, name, template_name, template_version, repository_url,
		       created_at, updated_at, created_by, team, status, tags, metadata
		FROM projects
	`

	args := make([]interface{}, 0)
	where := make([]string, 0)
	argPos := 1

	if filter != nil {
		if filter.Status != "" {
			where = append(where, fmt.Sprintf("status = %s", r.p(argPos)))
			args = append(args, filter.Status)
			argPos++
		}
		if filter.TemplateName != "" {
			where = append(where, fmt.Sprintf("template_name = %s", r.p(argPos)))
			args = append(args, filter.TemplateName)
			argPos++
		}
		if filter.Team != "" {
			where = append(where, fmt.Sprintf("team = %s", r.p(argPos)))
			args = append(args, filter.Team)
			argPos++
		}
	}

	queryBuilder := strings.Builder{}
	queryBuilder.WriteString(base)
	if len(where) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(where, " AND "))
	}
	queryBuilder.WriteString(" ORDER BY created_at DESC")

	if filter != nil && filter.Limit > 0 {
		queryBuilder.WriteString(fmt.Sprintf(" LIMIT %s", r.p(argPos)))
		args = append(args, filter.Limit)
		argPos++
	}

	if filter != nil && filter.Offset > 0 {
		queryBuilder.WriteString(fmt.Sprintf(" OFFSET %s", r.p(argPos)))
		args = append(args, filter.Offset)
	}

	rows, err := r.db.QueryContext(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, errors.NewError("list projects", err)
	}
	defer rows.Close()

	projects := make([]*Project, 0)
	for rows.Next() {
		project, err := scanProjectRows(rows)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.NewError("list projects", err)
	}

	return projects, nil
}

func (r *sqlProjectRepository) Update(ctx context.Context, project *Project) error {
	if project == nil {
		return fmt.Errorf("project cannot be nil")
	}

	project.UpdatedAt = time.Now().UTC()

	tagsJSON, err := json.Marshal(project.Tags)
	if err != nil {
		return fmt.Errorf("failed to serialize project tags: %w", err)
	}

	metadataJSON, err := json.Marshal(project.Metadata)
	if err != nil {
		return fmt.Errorf("failed to serialize project metadata: %w", err)
	}

	query := fmt.Sprintf(`
		UPDATE projects
		SET name = %s,
		    template_name = %s,
		    template_version = %s,
		    repository_url = %s,
		    updated_at = %s,
		    created_by = %s,
		    team = %s,
		    status = %s,
		    tags = %s,
		    metadata = %s
		WHERE id = %s
	`, r.p(1), r.p(2), r.p(3), r.p(4), r.p(5), r.p(6), r.p(7), r.p(8), r.p(9), r.p(10), r.p(11))

	result, err := r.db.ExecContext(
		ctx,
		query,
		project.Name,
		project.TemplateName,
		project.TemplateVersion,
		project.RepositoryURL,
		project.UpdatedAt,

		project.Status,
		string(tagsJSON),
		string(metadataJSON),
		project.ID,
	)
	if err != nil {
		return errors.NewError("update project", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewError("update project", err)
	}
	if rowsAffected == 0 {
		return errors.NewError("update project", sql.ErrNoRows)
	}

	return nil
}

func (r *sqlProjectRepository) Delete(ctx context.Context, id string) error {
	query := fmt.Sprintf("DELETE FROM projects WHERE id = %s", r.p(1))
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.NewError("delete project", err)
	}

	return nil
}

func (r *sqlProjectRepository) AddHistory(ctx context.Context, history *ProjectHistory) error {
	if history == nil {
		return fmt.Errorf("history cannot be nil")
	}

	if history.ID == "" {
		history.ID = generateID()
	}
	if history.Timestamp.IsZero() {
		history.Timestamp = time.Now().UTC()
	}

	detailsJSON, err := json.Marshal(history.Details)
	if err != nil {
		return fmt.Errorf("failed to serialize history details: %w", err)
	}

	query := fmt.Sprintf(`
		INSERT INTO project_history (id, project_id, action, details, timestamp, user_name)
		VALUES (%s, %s, %s, %s, %s, %s)
	`, r.p(1), r.p(2), r.p(3), r.p(4), r.p(5), r.p(6))

	_, err = r.db.ExecContext(
		ctx,
		query,
		history.ID,
		history.ProjectID,
		history.Action,
		string(detailsJSON),
		history.Timestamp,
		history.User,
	)
	if err != nil {
		return errors.NewError("add project history", err)
	}

	return nil
}

func (r *sqlProjectRepository) GetHistory(ctx context.Context, projectID string) ([]*ProjectHistory, error) {
	query := fmt.Sprintf(`
		SELECT id, project_id, action, details, timestamp, user_name
		FROM project_history
		WHERE project_id = %s
		ORDER BY timestamp DESC
	`, r.p(1))

	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, errors.NewError("get project history", err)
	}
	defer rows.Close()

	historyItems := make([]*ProjectHistory, 0)
	for rows.Next() {
		var history ProjectHistory
		var detailsRaw sql.NullString

		if err := rows.Scan(
			&history.ID,
			&history.ProjectID,
			&history.Action,
			&detailsRaw,
			&history.Timestamp,
			&history.User,
		); err != nil {
			return nil, errors.NewError("get project history", err)
		}

		history.Details = make(map[string]interface{})
		if detailsRaw.Valid && strings.TrimSpace(detailsRaw.String) != "" {
			if err := json.Unmarshal([]byte(detailsRaw.String), &history.Details); err != nil {
				return nil, errors.NewError("get project history", err)
			}
		}

		historyItems = append(historyItems, &history)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.NewError("get project history", err)
	}

	return historyItems, nil
}

func (r *sqlProjectRepository) GetStats(ctx context.Context) (*ProjectStats, error) {
	stats := &ProjectStats{
		ProjectsByTemplate: make(map[string]int),
		ProjectsByTeam:     make(map[string]int),
	}

	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM projects`).Scan(&stats.TotalProjects); err != nil {
		return nil, errors.NewError("get project stats", err)
	}

	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM projects WHERE status = 'active'`).Scan(&stats.ActiveProjects); err != nil {
		return nil, errors.NewError("get project stats", err)
	}

	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM projects WHERE status = 'archived'`).Scan(&stats.ArchivedProjects); err != nil {
		return nil, errors.NewError("get project stats", err)
	}

	recentQuery := `SELECT COUNT(*) FROM projects WHERE created_at >= ?`
	recentArgs := []interface{}{time.Now().UTC().Add(-30 * 24 * time.Hour)}
	if r.dialect == dialectPostgres {
		recentQuery = `SELECT COUNT(*) FROM projects WHERE created_at >= $1`
	}
	if err := r.db.QueryRowContext(ctx, recentQuery, recentArgs...).Scan(&stats.RecentProjects); err != nil {
		return nil, errors.NewError("get project stats", err)
	}

	if err := r.loadTemplateStats(ctx, stats); err != nil {
		return nil, err
	}

	if err := r.loadTeamStats(ctx, stats); err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *sqlProjectRepository) loadTemplateStats(ctx context.Context, stats *ProjectStats) error {
	rows, err := r.db.QueryContext(ctx, `SELECT template_name, COUNT(*) FROM projects GROUP BY template_name`)
	if err != nil {
		return errors.NewError("get template stats", err)
	}
	defer rows.Close()

	for rows.Next() {
		var templateName string
		var count int
		if err := rows.Scan(&templateName, &count); err != nil {
			return errors.NewError("get template stats", err)
		}
		stats.ProjectsByTemplate[templateName] = count
	}

	if err := rows.Err(); err != nil {
		return errors.NewError("get template stats", err)
	}

	return nil
}

func (r *sqlProjectRepository) loadTeamStats(ctx context.Context, stats *ProjectStats) error {
	rows, err := r.db.QueryContext(ctx, `SELECT COALESCE(team, ''), COUNT(*) FROM projects GROUP BY COALESCE(team, '')`)
	if err != nil {
		return errors.NewError("get team stats", err)
	}
	defer rows.Close()

	for rows.Next() {
		var team string
		var count int
		if err := rows.Scan(&team, &count); err != nil {
			return errors.NewError("get team stats", err)
		}
		stats.ProjectsByTeam[team] = count
	}

	if err := rows.Err(); err != nil {
		return errors.NewError("get team stats", err)
	}

	return nil
}

func (r *sqlProjectRepository) p(index int) string {
	if r.dialect == dialectPostgres {
		return fmt.Sprintf("$%d", index)
	}
	return "?"
}

type projectRowScanner interface {
	Scan(dest ...interface{}) error
}

func scanProjectRow(scanner projectRowScanner) (*Project, error) {
	var project Project
	var tagsRaw sql.NullString
	var metadataRaw sql.NullString

	err := scanner.Scan(
		&project.ID,
		&project.Name,
		&project.TemplateName,
		&project.TemplateVersion,
		&project.RepositoryURL,
		&project.CreatedAt,
		&project.UpdatedAt,
		&project.Status,
		&tagsRaw,
		&metadataRaw,
	)
	if err != nil {
		return nil, errors.NewError("scan project", err)
	}

	project.Tags = make([]string, 0)
	if tagsRaw.Valid && strings.TrimSpace(tagsRaw.String) != "" {
		if err := json.Unmarshal([]byte(tagsRaw.String), &project.Tags); err != nil {
			return nil, errors.NewError("scan project", err)
		}
	}

	project.Metadata = make(map[string]interface{})
	if metadataRaw.Valid && strings.TrimSpace(metadataRaw.String) != "" {
		if err := json.Unmarshal([]byte(metadataRaw.String), &project.Metadata); err != nil {
			return nil, errors.NewError("scan project", err)
		}
	}

	return &project, nil
}

func scanProjectRows(rows *sql.Rows) (*Project, error) {
	return scanProjectRow(rows)
}

func generateID() string {
	// 16 random bytes provide enough entropy for unique record IDs in this context.
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		// Fallback to time-based bytes only if crypto random fails unexpectedly.
		now := time.Now().UnixNano()
		for i := range raw {
			raw[i] = byte((now >> (i % 8)) & 0xff)
		}
	}
	return hex.EncodeToString(raw)
}
