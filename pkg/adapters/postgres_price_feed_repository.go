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

type PostgresPriceFeedRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewPostgresPriceFeedRepository(db *sql.DB, logger *logrus.Logger) interfaces.PriceFeedRepository {
	return &PostgresPriceFeedRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PostgresPriceFeedRepository) Create(ctx context.Context, feed *models.PriceFeed) error {
	return fmt.Errorf("not implemented: Create price feed")
}

func (r *PostgresPriceFeedRepository) GetByID(ctx context.Context, feedID string) (*models.PriceFeed, error) {
	return nil, fmt.Errorf("not implemented: GetByID price feed")
}

func (r *PostgresPriceFeedRepository) GetLatestBySymbol(ctx context.Context, symbol string) (*models.PriceFeed, error) {
	return nil, fmt.Errorf("not implemented: GetLatestBySymbol")
}

func (r *PostgresPriceFeedRepository) GetBySymbol(ctx context.Context, symbol string, limit int) ([]*models.PriceFeed, error) {
	return nil, fmt.Errorf("not implemented: GetBySymbol")
}

func (r *PostgresPriceFeedRepository) Query(ctx context.Context, query *models.PriceFeedQuery) ([]*models.PriceFeed, error) {
	return nil, fmt.Errorf("not implemented: Query price feeds")
}

func (r *PostgresPriceFeedRepository) DeleteOlderThan(ctx context.Context, timestamp time.Time) (int64, error) {
	return 0, fmt.Errorf("not implemented: DeleteOlderThan")
}
