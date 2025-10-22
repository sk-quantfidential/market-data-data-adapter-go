package models

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

type Symbol struct {
	SymbolID         string           `json:"symbol_id" db:"symbol_id"`
	Symbol           string           `json:"symbol" db:"symbol"`
	BaseCurrency     string           `json:"base_currency" db:"base_currency"`
	QuoteCurrency    string           `json:"quote_currency" db:"quote_currency"`
	DisplayName      *string          `json:"display_name,omitempty" db:"display_name"`
	IsActive         bool             `json:"is_active" db:"is_active"`
	MinPriceMovement *decimal.Decimal `json:"min_price_movement,omitempty" db:"min_price_movement"`
	MinOrderSize     *decimal.Decimal `json:"min_order_size,omitempty" db:"min_order_size"`
	MaxOrderSize     *decimal.Decimal `json:"max_order_size,omitempty" db:"max_order_size"`
	CreatedAt        time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at" db:"updated_at"`
	Metadata         json.RawMessage  `json:"metadata,omitempty" db:"metadata"`
}

type SymbolQuery struct {
	Symbol        *string
	BaseCurrency  *string
	QuoteCurrency *string
	IsActive      *bool
	Limit         int
	Offset        int
	SortBy        string
	SortOrder     string
}
