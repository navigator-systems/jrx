package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteDB provides a concrete Database implementation for local environments.
// It stores all records in a single file and auto-creates schema objects on startup.
type SQLiteDB struct {
	sqlDatabaseBase
	dbPath string
}

// NewSQLiteDB builds a SQLite database handler bound to the configured file path.
// Connection and schema setup happen through Connect + InitDB to keep lifecycle explicit.
func NewSQLiteDB(dbPath string) (*SQLiteDB, error) {
	return &SQLiteDB{dbPath: dbPath}, nil
}

func (s *SQLiteDB) Connect(ctx context.Context) error {
	if s.db != nil {
		return nil
	}

	dbConn, err := sql.Open("sqlite3", s.dbPath)
	if err != nil {
		return fmt.Errorf("failed to open sqlite database: %w", err)
	}

	if _, err := dbConn.ExecContext(ctx, `PRAGMA foreign_keys = ON`); err != nil {
		_ = dbConn.Close()
		return fmt.Errorf("failed to enable sqlite foreign keys: %w", err)
	}

	if err := dbConn.PingContext(ctx); err != nil {
		_ = dbConn.Close()
		return fmt.Errorf("failed to ping sqlite database: %w", err)
	}

	s.db = dbConn
	s.repo = newSQLProjectRepository(dbConn, dialectSQLite)
	log.Printf("SQLite database connected at '%s'", s.dbPath)
	return nil
}

func (s *SQLiteDB) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *SQLiteDB) Ping() error {
	if s.db == nil {
		return fmt.Errorf("sqlite database is not connected")
	}
	return s.db.Ping()
}

func (s *SQLiteDB) InitDB() error {
	if s.db == nil {
		return fmt.Errorf("sqlite database is not connected")
	}

	statements := []string{
		// projects keeps the canonical lifecycle state for generated projects.
		`CREATE TABLE IF NOT EXISTS projects (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			template_name TEXT NOT NULL,
			template_version TEXT,
			repository_url TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			created_by TEXT,
			team TEXT,
			status TEXT NOT NULL,
			tags TEXT,
			metadata TEXT
		)`,
		// project_history stores immutable audit records for each project action.
		`CREATE TABLE IF NOT EXISTS project_history (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL,
			action TEXT NOT NULL,
			details TEXT,
			timestamp TIMESTAMP NOT NULL,
			user_name TEXT,
			FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
		)`,
		// Indexes make common list/filter operations fast even with large project history.
		`CREATE INDEX IF NOT EXISTS idx_projects_template_name ON projects(template_name)`,
		`CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status)`,
		`CREATE INDEX IF NOT EXISTS idx_projects_created_at ON projects(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_projects_team ON projects(team)`,
		`CREATE INDEX IF NOT EXISTS idx_history_project_id ON project_history(project_id)`,
		`CREATE INDEX IF NOT EXISTS idx_history_timestamp ON project_history(timestamp)`,
	}

	for _, statement := range statements {
		if _, err := s.db.Exec(statement); err != nil {
			return fmt.Errorf("failed to initialize sqlite schema: %w", err)
		}
	}

	log.Printf("SQLite schema ensured for '%s'", s.dbPath)
	return nil
}

func (s *SQLiteDB) Project() ProjectRepository {
	return s.repo
}
