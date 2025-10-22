package interfaces

import (
	"context"
	"time"

	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/pkg/models"
)

type MarketSnapshotRepository interface {
	// Create a new market snapshot
	Create(ctx context.Context, snapshot *models.MarketSnapshot) error

	// Get snapshot by ID
	GetByID(ctx context.Context, snapshotID string) (*models.MarketSnapshot, error)

	// Get latest snapshot for a symbol
	GetLatestBySymbol(ctx context.Context, symbol string) (*models.MarketSnapshot, error)

	// Get snapshot history for a symbol
	GetBySymbol(ctx context.Context, symbol string, limit int) ([]*models.MarketSnapshot, error)

	// Query snapshots with filters
	Query(ctx context.Context, query *models.MarketSnapshotQuery) ([]*models.MarketSnapshot, error)

	// Delete old snapshots (cleanup)
	DeleteOlderThan(ctx context.Context, timestamp time.Time) (int64, error)
}
