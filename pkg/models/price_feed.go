package models

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

type PriceFeed struct {
	FeedID    string           `json:"feed_id" db:"feed_id"`
	Symbol    string           `json:"symbol" db:"symbol"`
	Price     decimal.Decimal  `json:"price" db:"price"`
	Bid       *decimal.Decimal `json:"bid,omitempty" db:"bid"`
	Ask       *decimal.Decimal `json:"ask,omitempty" db:"ask"`
	Volume24h *decimal.Decimal `json:"volume_24h,omitempty" db:"volume_24h"`
	Source    string           `json:"source" db:"source"`
	Timestamp time.Time        `json:"timestamp" db:"timestamp"`
	Metadata  json.RawMessage  `json:"metadata,omitempty" db:"metadata"`
}

type PriceFeedQuery struct {
	Symbol        *string
	Source        *string
	TimestampFrom *time.Time
	TimestampTo   *time.Time
	Limit         int
	Offset        int
	SortBy        string
	SortOrder     string
}
