package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/internal/config"
	"github.com/sirupsen/logrus"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB     *sql.DB
	config *config.Config
	logger *logrus.Logger
}

func NewPostgresDB(cfg *config.Config, logger *logrus.Logger) (*PostgresDB, error) {
	if cfg.PostgresURL == "" {
		return nil, fmt.Errorf("PostgreSQL URL is required")
	}

	return &PostgresDB{
		config: cfg,
		logger: logger,
	}, nil
}

func (p *PostgresDB) Connect(ctx context.Context) error {
	db, err := sql.Open("postgres", p.config.PostgresURL)
	if err != nil {
		return fmt.Errorf("failed to open PostgreSQL connection: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(p.config.MaxConnections)
	db.SetMaxIdleConns(p.config.MaxIdleConnections)
	db.SetConnMaxLifetime(p.config.ConnectionMaxLifetime)
	db.SetConnMaxIdleTime(p.config.ConnectionMaxIdleTime)

	// Test connection
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	p.DB = db
	p.logger.Info("PostgreSQL connected successfully")
	return nil
}

func (p *PostgresDB) Disconnect(ctx context.Context) error {
	if p.DB != nil {
		if err := p.DB.Close(); err != nil {
			return fmt.Errorf("failed to close PostgreSQL connection: %w", err)
		}
		p.logger.Info("PostgreSQL disconnected")
	}
	return nil
}

func (p *PostgresDB) HealthCheck(ctx context.Context) error {
	if p.DB == nil {
		return fmt.Errorf("PostgreSQL not connected")
	}
	return p.DB.PingContext(ctx)
}
