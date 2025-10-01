package adapters

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/pkg/interfaces"
	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/pkg/models"
	"github.com/sirupsen/logrus"
)

type PostgresSymbolRepository struct {
	db     *sql.DB
	logger *logrus.Logger
}

func NewPostgresSymbolRepository(db *sql.DB, logger *logrus.Logger) interfaces.SymbolRepository {
	return &PostgresSymbolRepository{
		db:     db,
		logger: logger,
	}
}

func (r *PostgresSymbolRepository) Create(ctx context.Context, symbol *models.Symbol) error {
	return fmt.Errorf("not implemented: Create symbol")
}

func (r *PostgresSymbolRepository) GetByID(ctx context.Context, symbolID string) (*models.Symbol, error) {
	return nil, fmt.Errorf("not implemented: GetByID symbol")
}

func (r *PostgresSymbolRepository) GetBySymbol(ctx context.Context, symbol string) (*models.Symbol, error) {
	return nil, fmt.Errorf("not implemented: GetBySymbol")
}

func (r *PostgresSymbolRepository) Query(ctx context.Context, query *models.SymbolQuery) ([]*models.Symbol, error) {
	return nil, fmt.Errorf("not implemented: Query symbols")
}

func (r *PostgresSymbolRepository) Update(ctx context.Context, symbol *models.Symbol) error {
	return fmt.Errorf("not implemented: Update symbol")
}

func (r *PostgresSymbolRepository) UpdateActiveStatus(ctx context.Context, symbolID string, isActive bool) error {
	return fmt.Errorf("not implemented: UpdateActiveStatus")
}

func (r *PostgresSymbolRepository) GetActive(ctx context.Context) ([]*models.Symbol, error) {
	return nil, fmt.Errorf("not implemented: GetActive")
}

func (r *PostgresSymbolRepository) Delete(ctx context.Context, symbolID string) error {
	return fmt.Errorf("not implemented: Delete symbol")
}
