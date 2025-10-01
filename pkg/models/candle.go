package models

import (
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

type CandleInterval string

const (
	Interval1m  CandleInterval = "1m"
	Interval5m  CandleInterval = "5m"
	Interval15m CandleInterval = "15m"
	Interval1h  CandleInterval = "1h"
	Interval4h  CandleInterval = "4h"
	Interval1d  CandleInterval = "1d"
)

type Candle struct {
	CandleID  string          `json:"candle_id" db:"candle_id"`
	Symbol    string          `json:"symbol" db:"symbol"`
	Interval  CandleInterval  `json:"interval" db:"interval"`
	Open      decimal.Decimal `json:"open" db:"open"`
	High      decimal.Decimal `json:"high" db:"high"`
	Low       decimal.Decimal `json:"low" db:"low"`
	Close     decimal.Decimal `json:"close" db:"close"`
	Volume    decimal.Decimal `json:"volume" db:"volume"`
	StartTime time.Time       `json:"start_time" db:"start_time"`
	EndTime   time.Time       `json:"end_time" db:"end_time"`
	NumTrades *int            `json:"num_trades,omitempty" db:"num_trades"`
	Metadata  json.RawMessage `json:"metadata,omitempty" db:"metadata"`
}

type CandleQuery struct {
	Symbol        *string
	Interval      *CandleInterval
	StartTimeFrom *time.Time
	StartTimeTo   *time.Time
	Limit         int
	Offset        int
	SortBy        string
	SortOrder     string
}
