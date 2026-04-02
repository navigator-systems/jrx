package db

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/navigator-systems/jrx/internal/config"
	"github.com/navigator-systems/jrx/internal/errors"
)

// PostgresConfig centralizes connection settings for PostgreSQL.
// Keeping this as a dedicated struct avoids scattering connection fields across the code.
type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewDatabase creates a new database instance based on configuration
func NewDatabase(cfg config.JRXConfig) (Database, error) {
	dbType := strings.ToLower(strings.TrimSpace(cfg.Database.Database))

	switch dbType {
	case "sqlite":
		dbPath := strings.TrimSpace(cfg.Database.DBPath)
		// Use a deterministic default path when no explicit SQLite file is configured.
		if dbPath == "" {
			dbPath = filepath.Join(".", "jrx_projects.db")
		}
		return NewSQLiteDB(dbPath)

	case "postgres", "postgresql":
		host := strings.TrimSpace(cfg.Database.DBHost)
		sslMode := ""

		if host != "" && host != "localhost" && host != "127.0.0.1" && host != "::1" {
			sslMode = "require"
		}

		pgConfig := PostgresConfig{
			Host:     host,
			Port:     cfg.Database.DBPort,
			User:     strings.TrimSpace(cfg.Database.DBUser),
			Password: cfg.Database.DBPassword,
			DBName:   strings.TrimSpace(cfg.Database.DBName),
			SSLMode:  sslMode,
		}
		return NewPostgresDB(pgConfig)

	default:
		return nil, errors.ErrInvalidDatabaseConfig
	}
}

// InitDatabase initializes database connectivity and creates the schema when needed.
// The operation is idempotent by design because InitDB internally uses IF NOT EXISTS.
func InitDatabase(ctx context.Context, cfg config.JRXConfig) (Database, error) {
	// If database is not configured, tracking remains disabled and the caller can continue normally.
	dbType := strings.TrimSpace(cfg.Database.Database)
	if dbType == "" {
		return nil, nil
	}

	db, err := NewDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}

	if err := db.Connect(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.InitDB(); err != nil {
		return nil, fmt.Errorf("failed to initialize database schema: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// sqlDatabaseBase contains shared members used by concrete DB implementations.
// This reduces duplicate fields and keeps constructors concise.
type sqlDatabaseBase struct {
	db   *sql.DB
	repo ProjectRepository
}
