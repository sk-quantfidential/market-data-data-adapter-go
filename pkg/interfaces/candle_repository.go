package interfaces

import (
	"context"
	"time"

	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/pkg/models"
)

type CandleRepository interface {
	// Create or update a candle
	Upsert(ctx context.Context, candle *models.Candle) error

	// Get candle by ID
	GetByID(ctx context.Context, candleID string) (*models.Candle, error)

	// Get candles for symbol and interval
	GetBySymbolAndInterval(ctx context.Context, symbol string, interval models.CandleInterval, limit int) ([]*models.Candle, error)

	// Query candles with filters
	Query(ctx context.Context, query *models.CandleQuery) ([]*models.Candle, error)

	// Get latest candle for symbol and interval
	GetLatest(ctx context.Context, symbol string, interval models.CandleInterval) (*models.Candle, error)

	// Delete old candles (cleanup)
	DeleteOlderThan(ctx context.Context, timestamp time.Time) (int64, error)
}
