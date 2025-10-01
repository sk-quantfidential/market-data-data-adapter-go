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

type PostgresMarketSnapshotRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewPostgresMarketSnapshotRepository(db *sql.DB, logger *logrus.Logger) interfaces.MarketSnapshotRepository {
	return &PostgresMarketSnapshotRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PostgresMarketSnapshotRepository) Create(ctx context.Context, snapshot *models.MarketSnapshot) error {
	return fmt.Errorf("not implemented: Create snapshot")
}

func (r *PostgresMarketSnapshotRepository) GetByID(ctx context.Context, snapshotID string) (*models.MarketSnapshot, error) {
	return nil, fmt.Errorf("not implemented: GetByID snapshot")
}

func (r *PostgresMarketSnapshotRepository) GetLatestBySymbol(ctx context.Context, symbol string) (*models.MarketSnapshot, error) {
	return nil, fmt.Errorf("not implemented: GetLatestBySymbol")
}

func (r *PostgresMarketSnapshotRepository) GetBySymbol(ctx context.Context, symbol string, limit int) ([]*models.MarketSnapshot, error) {
	return nil, fmt.Errorf("not implemented: GetBySymbol")
}

func (r *PostgresMarketSnapshotRepository) Query(ctx context.Context, query *models.MarketSnapshotQuery) ([]*models.MarketSnapshot, error) {
	return nil, fmt.Errorf("not implemented: Query snapshots")
}

func (r *PostgresMarketSnapshotRepository) DeleteOlderThan(ctx context.Context, timestamp time.Time) (int64, error) {
	return 0, fmt.Errorf("not implemented: DeleteOlderThan")
}
