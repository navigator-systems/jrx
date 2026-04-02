package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// PostgresDB provides a production-ready Database implementation.
// It supports the same repository contract as SQLite, enabling seamless env switching.
type PostgresDB struct {
	sqlDatabaseBase
	config PostgresConfig
}

// NewPostgresDB builds a PostgreSQL database handler from typed configuration.
func NewPostgresDB(cfg PostgresConfig) (*PostgresDB, error) {
	return &PostgresDB{config: cfg}, nil
}

func (p *PostgresDB) Connect(ctx context.Context) error {
	if p.db != nil {
		return nil
	}

	if p.config.Port == 0 {
		p.config.Port = 5432
	}
	if p.config.SSLMode == "" {
		p.config.SSLMode = "disable"
	}

	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.config.Host,
		p.config.Port,
		p.config.User,
		p.config.Password,
		p.config.DBName,
		p.config.SSLMode,
	)

	dbConn, err := sql.Open("postgres", connString)
	if err != nil {
		return fmt.Errorf("failed to open postgres database: %w", err)
	}

	if err := dbConn.PingContext(ctx); err != nil {
		_ = dbConn.Close()
		return fmt.Errorf("failed to ping postgres database: %w", err)
	}

	p.db = dbConn
	p.repo = newSQLProjectRepository(dbConn, dialectPostgres)
	log.Printf("PostgreSQL database connected at '%s:%d/%s'", p.config.Host, p.config.Port, p.config.DBName)
	return nil
}

func (p *PostgresDB) Close() error {
	if p.db == nil {
		return nil
	}
	return p.db.Close()
}

func (p *PostgresDB) Ping() error {
	if p.db == nil {
		return fmt.Errorf("postgres database is not connected")
	}
	return p.db.Ping()
}

func (p *PostgresDB) InitDB() error {
	if p.db == nil {
		return fmt.Errorf("postgres database is not connected")
	}

	statements := []string{
		// projects is the source of truth for project records created via CLI or web UI.
		`CREATE TABLE IF NOT EXISTS projects (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			template_name TEXT NOT NULL,
			template_version TEXT,
			repository_url TEXT,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL,
			created_by TEXT,
			team TEXT,
			status TEXT NOT NULL,
			tags JSONB,
			metadata JSONB
		)`,
		// project_history tracks immutable events for auditing and activity timelines.
		`CREATE TABLE IF NOT EXISTS project_history (
			id TEXT PRIMARY KEY,
			project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
			action TEXT NOT NULL,
			details JSONB,
			timestamp TIMESTAMPTZ NOT NULL,
			user_name TEXT
		)`,
		// Indexes optimize filtering and timeline access patterns.
		`CREATE INDEX IF NOT EXISTS idx_projects_template_name ON projects(template_name)`,
		`CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status)`,
		`CREATE INDEX IF NOT EXISTS idx_projects_created_at ON projects(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_projects_team ON projects(team)`,
		`CREATE INDEX IF NOT EXISTS idx_history_project_id ON project_history(project_id)`,
		`CREATE INDEX IF NOT EXISTS idx_history_timestamp ON project_history(timestamp)`,
	}

	for _, statement := range statements {
		if _, err := p.db.Exec(statement); err != nil {
			return fmt.Errorf("failed to initialize postgres schema: %w", err)
		}
	}

	log.Printf("PostgreSQL schema ensured for '%s:%d/%s'", p.config.Host, p.config.Port, p.config.DBName)
	return nil
}

func (p *PostgresDB) Project() ProjectRepository {
	return p.repo
}
