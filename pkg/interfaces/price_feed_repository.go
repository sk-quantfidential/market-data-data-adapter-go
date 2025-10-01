package interfaces

import (
	"context"
	"time"

	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/pkg/models"
)

type PriceFeedRepository interface {
	// Create a new price feed entry
	Create(ctx context.Context, feed *models.PriceFeed) error

	// Get price feed by ID
	GetByID(ctx context.Context, feedID string) (*models.PriceFeed, error)

	// Get latest price for a symbol
	GetLatestBySymbol(ctx context.Context, symbol string) (*models.PriceFeed, error)

	// Get price history for a symbol
	GetBySymbol(ctx context.Context, symbol string, limit int) ([]*models.PriceFeed, error)

	// Query price feeds with filters
	Query(ctx context.Context, query *models.PriceFeedQuery) ([]*models.PriceFeed, error)

	// Delete old price feeds (cleanup)
	DeleteOlderThan(ctx context.Context, timestamp time.Time) (int64, error)
}
