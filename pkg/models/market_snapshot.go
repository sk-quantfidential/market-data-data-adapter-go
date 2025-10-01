package models

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

type MarketSnapshot struct {
	SnapshotID            string           `json:"snapshot_id" db:"snapshot_id"`
	Symbol                string           `json:"symbol" db:"symbol"`
	LastPrice             decimal.Decimal  `json:"last_price" db:"last_price"`
	Bid                   *decimal.Decimal `json:"bid,omitempty" db:"bid"`
	Ask                   *decimal.Decimal `json:"ask,omitempty" db:"ask"`
	Spread                *decimal.Decimal `json:"spread,omitempty" db:"spread"`
	Volume24h             *decimal.Decimal `json:"volume_24h,omitempty" db:"volume_24h"`
	PriceChange24h        *decimal.Decimal `json:"price_change_24h,omitempty" db:"price_change_24h"`
	PriceChangePercent24h *decimal.Decimal `json:"price_change_percent_24h,omitempty" db:"price_change_percent_24h"`
	Timestamp             time.Time        `json:"timestamp" db:"timestamp"`
	Metadata              json.RawMessage  `json:"metadata,omitempty" db:"metadata"`
}

type MarketSnapshotQuery struct {
	Symbol        *string
	TimestampFrom *time.Time
	TimestampTo   *time.Time
	Limit         int
	Offset        int
	SortBy        string
	SortOrder     string
}
