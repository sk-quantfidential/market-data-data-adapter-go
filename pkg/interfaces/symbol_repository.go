package interfaces

import (
	"context"

	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/pkg/models"
)

type SymbolRepository interface {
	// Create a new symbol
	Create(ctx context.Context, symbol *models.Symbol) error

	// Get symbol by ID
	GetByID(ctx context.Context, symbolID string) (*models.Symbol, error)

	// Get symbol by symbol string
	GetBySymbol(ctx context.Context, symbol string) (*models.Symbol, error)

	// Query symbols with filters
	Query(ctx context.Context, query *models.SymbolQuery) ([]*models.Symbol, error)

	// Update symbol
	Update(ctx context.Context, symbol *models.Symbol) error

	// Activate/deactivate symbol
	UpdateActiveStatus(ctx context.Context, symbolID string, isActive bool) error

	// Get all active symbols
	GetActive(ctx context.Context) ([]*models.Symbol, error)

	// Delete symbol
	Delete(ctx context.Context, symbolID string) error
}
