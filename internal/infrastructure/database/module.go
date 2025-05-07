package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver
	"go.uber.org/fx"

	"ecosystem.garyle/service/internal/app/config"
	"ecosystem.garyle/service/pkg/logger"
)

// Module provides database dependencies
var Module = fx.Module("database",
	fx.Provide(
		NewDBConnection,
	),
)

// NewDBConnection creates a new database connection
func NewDBConnection(cfg *config.Config, log logger.Logger) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)

	log.Infof("Connected to database: %s", cfg.Database.Name)
	return db, nil
}

// RegisterHooks registers lifecycle hooks for the database
func RegisterHooks(lc fx.Lifecycle, db *sql.DB, log logger.Logger) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Info("Closing database connection")
			return db.Close()
		},
	})
}
