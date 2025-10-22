package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/pkg/interfaces"
	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/pkg/models"
	"github.com/sirupsen/logrus"
)

type PostgresCandleRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewPostgresCandleRepository(db *sql.DB, logger *logrus.Logger) interfaces.CandleRepository {
	return &PostgresCandleRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PostgresCandleRepository) Upsert(ctx context.Context, candle *models.Candle) error {
	return fmt.Errorf("not implemented: Upsert candle")
}

func (r *PostgresCandleRepository) GetByID(ctx context.Context, candleID string) (*models.Candle, error) {
	return nil, fmt.Errorf("not implemented: GetByID candle")
}

func (r *PostgresCandleRepository) GetBySymbolAndInterval(ctx context.Context, symbol string, interval models.CandleInterval, limit int) ([]*models.Candle, error) {
	return nil, fmt.Errorf("not implemented: GetBySymbolAndInterval")
}

func (r *PostgresCandleRepository) Query(ctx context.Context, query *models.CandleQuery) ([]*models.Candle, error) {
	return nil, fmt.Errorf("not implemented: Query candles")
}

func (r *PostgresCandleRepository) GetLatest(ctx context.Context, symbol string, interval models.CandleInterval) (*models.Candle, error) {
	return nil, fmt.Errorf("not implemented: GetLatest candle")
}

func (r *PostgresCandleRepository) DeleteOlderThan(ctx context.Context, timestamp time.Time) (int64, error) {
	return 0, fmt.Errorf("not implemented: DeleteOlderThan")
}
