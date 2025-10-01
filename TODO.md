# market-data-adapter-go - TSE-0001.4 Data Adapters and Orchestrator Integration

## Milestone: TSE-0001.4 - Data Adapters and Orchestrator Integration
**Status**: 📝 **PENDING** - Ready to Start
**Goal**: Create market data adapter following audit-data-adapter-go, custodian-data-adapter-go, and exchange-data-adapter-go proven pattern
**Components**: Market Data Adapter Go
**Dependencies**: TSE-0001.3a (Core Infrastructure Setup) ✅, previous data adapter patterns ✅
**Estimated Time**: 8-10 hours following established pattern

## 🎯 BDD Acceptance Criteria
> The market data adapter can connect to orchestrator PostgreSQL and Redis services, handle market-data-specific operations (price feeds, candles, snapshots, symbols), and pass comprehensive behavior tests with proper environment configuration management.

## 📋 Repository Creation and Setup

### Initial Repository Structure
```
market-data-adapter-go/
├── .gitignore
├── .env.example
├── go.mod
├── go.sum
├── README.md
├── TODO.md (this file)
├── cmd/
│   └── example/
│       └── main.go                    # Example usage
├── internal/
│   ├── config/
│   │   └── config.go                  # Environment configuration
│   ├── database/
│   │   └── postgres.go                # PostgreSQL connection
│   └── cache/
│       └── redis.go                   # Redis connection
├── pkg/
│   ├── adapters/
│   │   ├── factory.go                 # DataAdapter factory
│   │   ├── postgres_adapter.go        # PostgreSQL implementation
│   │   └── redis_adapter.go           # Redis implementation
│   ├── interfaces/
│   │   ├── price_feed_repository.go   # Price feed operations
│   │   ├── candle_repository.go       # Candle (OHLCV) operations
│   │   ├── market_snapshot_repository.go  # Snapshot operations
│   │   ├── symbol_repository.go       # Symbol metadata operations
│   │   ├── service_discovery.go       # Service discovery (shared)
│   │   └── cache.go                   # Cache operations (shared)
│   └── models/
│       ├── price_feed.go              # Price feed model
│       ├── candle.go                  # Candle model
│       ├── market_snapshot.go         # Market snapshot model
│       └── symbol.go                  # Symbol model
└── tests/
    ├── init_test.go                   # Test initialization with godotenv
    ├── behavior_test_suite.go         # BDD test framework
    ├── price_feed_behavior_test.go    # Price feed tests
    ├── candle_behavior_test.go        # Candle tests
    ├── snapshot_behavior_test.go      # Snapshot tests
    ├── symbol_behavior_test.go        # Symbol tests
    ├── service_discovery_behavior_test.go
    ├── cache_behavior_test.go
    ├── integration_behavior_test.go
    └── test_utils.go                  # Test utilities
```

## 📋 Task Checklist

### Task 0: Repository Creation and Foundation
**Goal**: Create repository structure and base configuration
**Estimated Time**: 1 hour

#### Steps:
- [ ] Create repository directory structure
- [ ] Initialize go.mod with dependencies:
  ```go
  module github.com/quantfidential/trading-ecosystem/market-data-adapter-go

  go 1.24

  require (
      github.com/lib/pq v1.10.9                    // PostgreSQL driver
      github.com/redis/go-redis/v9 v9.15.0        // Redis client
      github.com/sirupsen/logrus v1.9.3           // Logging
      github.com/joho/godotenv v1.5.1             // Environment loading
      github.com/stretchr/testify v1.8.4          // Testing framework
      github.com/shopspring/decimal v1.3.1        // Decimal precision for prices
      google.golang.org/grpc v1.58.3              // gRPC (for models compatibility)
      google.golang.org/protobuf v1.31.0          // Protobuf (for models)
  )
  ```
- [ ] Create .gitignore (copy from audit-data-adapter-go)
- [ ] Create README.md with overview and usage instructions
- [ ] Create .env.example (see configuration below)
- [ ] Create Makefile with test automation

**Evidence to Check**:
- Repository structure created
- go.mod initialized with correct dependencies
- .env.example ready for configuration
- Makefile with test targets

---

### Task 1: Environment Configuration System
**Goal**: Create production-ready .env configuration following 12-factor app principles
**Estimated Time**: 30 minutes

#### .env.example Template:
```bash
# Market Data Adapter Configuration
# Copy this to .env and update with your orchestrator credentials

# Service Identity
SERVICE_NAME=market-data-adapter
SERVICE_VERSION=1.0.0
ENVIRONMENT=development

# PostgreSQL Configuration (orchestrator credentials)
POSTGRES_URL=postgres://market_data_adapter:market-data-adapter-db-pass@localhost:5432/trading_ecosystem?sslmode=disable

# PostgreSQL Connection Pool
MAX_CONNECTIONS=25
MAX_IDLE_CONNECTIONS=10
CONNECTION_MAX_LIFETIME=300s
CONNECTION_MAX_IDLE_TIME=60s

# Redis Configuration (orchestrator credentials)
# Production: Use market-data-adapter user
# Testing: Use admin user for full access
REDIS_URL=redis://market-data-adapter:market-data-pass@localhost:6379/0

# Redis Connection Pool
REDIS_POOL_SIZE=10
REDIS_MIN_IDLE_CONNS=2
REDIS_MAX_RETRIES=3
REDIS_DIAL_TIMEOUT=5s
REDIS_READ_TIMEOUT=3s
REDIS_WRITE_TIMEOUT=3s

# Cache Configuration
CACHE_TTL=300s                          # 5 minutes default TTL
CACHE_NAMESPACE=market_data             # Redis key prefix

# Service Discovery
SERVICE_DISCOVERY_NAMESPACE=market_data # Service registry namespace
HEARTBEAT_INTERVAL=30s                  # Service heartbeat frequency
SERVICE_TTL=90s                         # Service registration TTL

# Test Environment (for integration tests)
TEST_POSTGRES_URL=postgres://market_data_adapter:market-data-adapter-db-pass@localhost:5432/trading_ecosystem?sslmode=disable
TEST_REDIS_URL=redis://admin:admin-secure-pass@localhost:6379/0

# Logging
LOG_LEVEL=info                          # debug, info, warn, error
LOG_FORMAT=json                         # json, text

# Performance Testing
PERF_TEST_SIZE=1000                     # Number of items for performance tests
PERF_THROUGHPUT_MIN=100                 # Minimum ops/second
PERF_LATENCY_MAX=100ms                  # Maximum average latency

# CI/CD
SKIP_INTEGRATION_TESTS=false            # Set to true in CI without infrastructure
```

#### Configuration Implementation (internal/config/config.go):

Follow audit-data-adapter-go pattern with:
- Environment variable loading with defaults
- godotenv integration for .env file loading
- Type-safe configuration struct
- Helper functions: `getEnv()`, `getEnvInt()`, `getEnvDuration()`, `getEnvBool()`

**Acceptance Criteria**:
- [ ] .env.example created with orchestrator credentials
- [ ] Configuration loading working with defaults
- [ ] godotenv integration for test environment
- [ ] All configuration values accessible via Config struct
- [ ] .gitignore includes .env for security

---

### Task 2: Database Schema and Models
**Goal**: Define market-data-specific database schema and Go models
**Estimated Time**: 2 hours

#### Database Schema (PostgreSQL)

**Schema**: `market_data` (to be created in orchestrator)

**Tables**:

```sql
-- price_feeds: Real-time price data
CREATE TABLE market_data.price_feeds (
    feed_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    symbol VARCHAR(50) NOT NULL,
    price DECIMAL(24, 8) NOT NULL,
    bid DECIMAL(24, 8),
    ask DECIMAL(24, 8),
    volume_24h DECIMAL(24, 8),
    source VARCHAR(100) NOT NULL DEFAULT 'simulator',
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB,

    CONSTRAINT positive_price CHECK (price > 0),
    CONSTRAINT positive_bid CHECK (bid IS NULL OR bid > 0),
    CONSTRAINT positive_ask CHECK (ask IS NULL OR ask > 0),
    CONSTRAINT positive_volume CHECK (volume_24h IS NULL OR volume_24h >= 0)
);

CREATE INDEX idx_price_feeds_symbol ON market_data.price_feeds(symbol);
CREATE INDEX idx_price_feeds_timestamp ON market_data.price_feeds(timestamp DESC);
CREATE INDEX idx_price_feeds_symbol_timestamp ON market_data.price_feeds(symbol, timestamp DESC);
CREATE INDEX idx_price_feeds_source ON market_data.price_feeds(source);

-- candles: OHLCV candle data
CREATE TABLE market_data.candles (
    candle_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    symbol VARCHAR(50) NOT NULL,
    interval VARCHAR(10) NOT NULL, -- '1m', '5m', '15m', '1h', '4h', '1d'
    open DECIMAL(24, 8) NOT NULL,
    high DECIMAL(24, 8) NOT NULL,
    low DECIMAL(24, 8) NOT NULL,
    close DECIMAL(24, 8) NOT NULL,
    volume DECIMAL(24, 8) NOT NULL DEFAULT 0,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    num_trades INTEGER DEFAULT 0,
    metadata JSONB,

    CONSTRAINT positive_ohlc CHECK (open > 0 AND high > 0 AND low > 0 AND close > 0),
    CONSTRAINT valid_high_low CHECK (high >= low),
    CONSTRAINT high_gte_open_close CHECK (high >= open AND high >= close),
    CONSTRAINT low_lte_open_close CHECK (low <= open AND low <= close),
    CONSTRAINT positive_volume CHECK (volume >= 0),
    CONSTRAINT non_negative_trades CHECK (num_trades >= 0),
    CONSTRAINT unique_symbol_interval_time UNIQUE (symbol, interval, start_time)
);

CREATE INDEX idx_candles_symbol ON market_data.candles(symbol);
CREATE INDEX idx_candles_interval ON market_data.candles(interval);
CREATE INDEX idx_candles_start_time ON market_data.candles(start_time DESC);
CREATE INDEX idx_candles_symbol_interval_time ON market_data.candles(symbol, interval, start_time DESC);

-- market_snapshots: Periodic market state snapshots
CREATE TABLE market_data.market_snapshots (
    snapshot_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    symbol VARCHAR(50) NOT NULL,
    last_price DECIMAL(24, 8) NOT NULL,
    bid DECIMAL(24, 8),
    ask DECIMAL(24, 8),
    spread DECIMAL(24, 8),
    volume_24h DECIMAL(24, 8),
    price_change_24h DECIMAL(24, 8),
    price_change_percent_24h DECIMAL(10, 4),
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB,

    CONSTRAINT positive_last_price CHECK (last_price > 0),
    CONSTRAINT positive_spread CHECK (spread IS NULL OR spread >= 0)
);

CREATE INDEX idx_snapshots_symbol ON market_data.market_snapshots(symbol);
CREATE INDEX idx_snapshots_timestamp ON market_data.market_snapshots(timestamp DESC);
CREATE INDEX idx_snapshots_symbol_timestamp ON market_data.market_snapshots(symbol, timestamp DESC);

-- symbols: Trading symbol metadata
CREATE TABLE market_data.symbols (
    symbol_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    symbol VARCHAR(50) NOT NULL UNIQUE,
    base_currency VARCHAR(10) NOT NULL,
    quote_currency VARCHAR(10) NOT NULL,
    display_name VARCHAR(100),
    is_active BOOLEAN NOT NULL DEFAULT true,
    min_price_movement DECIMAL(24, 8),
    min_order_size DECIMAL(24, 8),
    max_order_size DECIMAL(24, 8),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata JSONB
);

CREATE INDEX idx_symbols_active ON market_data.symbols(is_active);
CREATE INDEX idx_symbols_base_currency ON market_data.symbols(base_currency);
CREATE INDEX idx_symbols_quote_currency ON market_data.symbols(quote_currency);
```

#### Go Models (pkg/models/)

**pkg/models/price_feed.go**:
```go
package models

import (
    "encoding/json"
    "time"
    "github.com/shopspring/decimal"
)

type PriceFeed struct {
    FeedID    string          `json:"feed_id" db:"feed_id"`
    Symbol    string          `json:"symbol" db:"symbol"`
    Price     decimal.Decimal `json:"price" db:"price"`
    Bid       *decimal.Decimal `json:"bid,omitempty" db:"bid"`
    Ask       *decimal.Decimal `json:"ask,omitempty" db:"ask"`
    Volume24h *decimal.Decimal `json:"volume_24h,omitempty" db:"volume_24h"`
    Source    string          `json:"source" db:"source"`
    Timestamp time.Time       `json:"timestamp" db:"timestamp"`
    Metadata  json.RawMessage `json:"metadata,omitempty" db:"metadata"`
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
```

**pkg/models/candle.go**:
```go
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
```

**pkg/models/market_snapshot.go**:
```go
package models

import (
    "encoding/json"
    "time"
    "github.com/shopspring/decimal"
)

type MarketSnapshot struct {
    SnapshotID           string           `json:"snapshot_id" db:"snapshot_id"`
    Symbol               string           `json:"symbol" db:"symbol"`
    LastPrice            decimal.Decimal  `json:"last_price" db:"last_price"`
    Bid                  *decimal.Decimal `json:"bid,omitempty" db:"bid"`
    Ask                  *decimal.Decimal `json:"ask,omitempty" db:"ask"`
    Spread               *decimal.Decimal `json:"spread,omitempty" db:"spread"`
    Volume24h            *decimal.Decimal `json:"volume_24h,omitempty" db:"volume_24h"`
    PriceChange24h       *decimal.Decimal `json:"price_change_24h,omitempty" db:"price_change_24h"`
    PriceChangePercent24h *decimal.Decimal `json:"price_change_percent_24h,omitempty" db:"price_change_percent_24h"`
    Timestamp            time.Time        `json:"timestamp" db:"timestamp"`
    Metadata             json.RawMessage  `json:"metadata,omitempty" db:"metadata"`
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
```

**pkg/models/symbol.go**:
```go
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
```

**Acceptance Criteria**:
- [ ] Database schema defined for market_data domain (4 tables)
- [ ] Go models created with proper JSON tags
- [ ] Query models for flexible filtering
- [ ] Enums for candle intervals
- [ ] Proper use of decimal.Decimal for prices and volumes
- [ ] Proper use of json.RawMessage for metadata

---

### Task 3: Repository Interfaces
**Goal**: Define clean interfaces for all market data operations
**Estimated Time**: 1 hour

#### Price Feed Repository (pkg/interfaces/price_feed_repository.go):
```go
package interfaces

import (
    "context"
    "github.com/quantfidential/trading-ecosystem/market-data-adapter-go/pkg/models"
)

type PriceFeedRepository interface {
    // Create a new price feed entry
    Create(ctx context.Context, feed *models.PriceFeed) error

    // Get price feed by ID
    GetByID(ctx context.Context, feedID string) (*models.PriceFeed, error)

    // Get latest price for a symbol
    GetLatestBySymbol(ctx context.Context, symbol string) (*models.PriceFeed, error)

    // Get price history for a symbol
    GetBySymbol(ctx context.Context, symbol string, limit int) ([]*models.PriceFeed, error)

    // Query price feeds with filters
    Query(ctx context.Context, query *models.PriceFeedQuery) ([]*models.PriceFeed, error)

    // Delete old price feeds (cleanup)
    DeleteOlderThan(ctx context.Context, timestamp time.Time) (int64, error)
}
```

#### Candle Repository (pkg/interfaces/candle_repository.go):
```go
package interfaces

import (
    "context"
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
```

#### Market Snapshot Repository (pkg/interfaces/market_snapshot_repository.go):
```go
package interfaces

import (
    "context"
    "github.com/quantfidential/trading-ecosystem/market-data-adapter-go/pkg/models"
)

type MarketSnapshotRepository interface {
    // Create a new market snapshot
    Create(ctx context.Context, snapshot *models.MarketSnapshot) error

    // Get snapshot by ID
    GetByID(ctx context.Context, snapshotID string) (*models.MarketSnapshot, error)

    // Get latest snapshot for a symbol
    GetLatestBySymbol(ctx context.Context, symbol string) (*models.MarketSnapshot, error)

    // Get snapshot history for a symbol
    GetBySymbol(ctx context.Context, symbol string, limit int) ([]*models.MarketSnapshot, error)

    // Query snapshots with filters
    Query(ctx context.Context, query *models.MarketSnapshotQuery) ([]*models.MarketSnapshot, error)

    // Delete old snapshots (cleanup)
    DeleteOlderThan(ctx context.Context, timestamp time.Time) (int64, error)
}
```

#### Symbol Repository (pkg/interfaces/symbol_repository.go):
```go
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
```

#### Shared Interfaces (copy from audit-data-adapter-go):

**pkg/interfaces/service_discovery.go** - Same as audit-data-adapter-go
**pkg/interfaces/cache.go** - Same as audit-data-adapter-go

**Acceptance Criteria**:
- [ ] All repository interfaces defined
- [ ] Methods follow CRUD + domain-specific operations pattern
- [ ] Context passed to all methods
- [ ] Proper error handling signatures
- [ ] Query methods use query models for flexibility
- [ ] Cleanup methods for data retention management

---

### Task 4: PostgreSQL Implementation
**Goal**: Implement repository interfaces using PostgreSQL
**Estimated Time**: 3 hours

Follow audit-data-adapter-go pattern for:
- Connection management (internal/database/postgres.go)
- Repository implementations (pkg/adapters/postgres_*.go)
- Transaction support for atomic operations
- Error handling and logging
- Connection pooling

**Files to Create**:
- `internal/database/postgres.go` - Connection management
- `pkg/adapters/postgres_price_feed_repository.go` - Price feed operations
- `pkg/adapters/postgres_candle_repository.go` - Candle operations
- `pkg/adapters/postgres_market_snapshot_repository.go` - Snapshot operations
- `pkg/adapters/postgres_symbol_repository.go` - Symbol operations

**Key Implementation Notes**:
- Use prepared statements for performance
- Handle decimal.Decimal properly in PostgreSQL queries
- Implement upsert for candles (ON CONFLICT UPDATE)
- Efficient time-series queries with proper indexing
- Bulk deletion for data retention

**Acceptance Criteria**:
- [ ] PostgreSQL connection with pooling
- [ ] All repository interfaces implemented
- [ ] Proper error handling and logging
- [ ] Transaction support for atomic operations
- [ ] Decimal precision maintained
- [ ] Health check implementation
- [ ] Efficient time-series queries

---

### Task 5: Redis Implementation
**Goal**: Implement caching and service discovery using Redis
**Estimated Time**: 2 hours

Follow audit-data-adapter-go pattern for:
- Redis connection management (internal/cache/redis.go)
- Service discovery implementation (pkg/adapters/redis_service_discovery.go)
- Cache repository implementation (pkg/adapters/redis_cache_repository.go)

**Additional Features for Market Data**:
- Price feed caching with short TTL (1-5 seconds)
- Latest price caching by symbol
- Candle caching for recent intervals

**Acceptance Criteria**:
- [ ] Redis connection with pooling
- [ ] Service discovery working with market_data:* namespace
- [ ] Cache operations with TTL management
- [ ] Price feed caching for real-time access
- [ ] Health check implementation
- [ ] Graceful fallback when Redis unavailable

---

### Task 6: DataAdapter Factory
**Goal**: Create factory pattern for adapter initialization
**Estimated Time**: 1 hour

#### pkg/adapters/factory.go:
```go
package adapters

import (
    "context"
    "github.com/quantfidential/trading-ecosystem/market-data-adapter-go/internal/config"
    "github.com/quantfidential/trading-ecosystem/market-data-adapter-go/pkg/interfaces"
    "github.com/sirupsen/logrus"
)

type DataAdapter interface {
    // Repository access
    PriceFeedRepository() interfaces.PriceFeedRepository
    CandleRepository() interfaces.CandleRepository
    MarketSnapshotRepository() interfaces.MarketSnapshotRepository
    SymbolRepository() interfaces.SymbolRepository
    ServiceDiscoveryRepository() interfaces.ServiceDiscoveryRepository
    CacheRepository() interfaces.CacheRepository

    // Lifecycle
    Connect(ctx context.Context) error
    Disconnect(ctx context.Context) error
    HealthCheck(ctx context.Context) error
}

func NewMarketDataAdapter(cfg *config.Config, logger *logrus.Logger) (DataAdapter, error) {
    // Implementation following audit-data-adapter-go pattern
}

func NewMarketDataAdapterFromEnv(logger *logrus.Logger) (DataAdapter, error) {
    // Load config from environment and create adapter
}
```

**Acceptance Criteria**:
- [ ] Factory pattern implemented
- [ ] Environment-based initialization
- [ ] Proper lifecycle management
- [ ] Health check aggregation

---

### Task 7: BDD Behavior Testing Framework
**Goal**: Create comprehensive test suite following audit-data-adapter-go pattern
**Estimated Time**: 3 hours

#### Test Files to Create:
- `tests/init_test.go` - godotenv loading and test setup
- `tests/behavior_test_suite.go` - BDD framework with Given/When/Then
- `tests/price_feed_behavior_test.go` - Price feed CRUD and queries
- `tests/candle_behavior_test.go` - Candle upsert and time-series queries
- `tests/snapshot_behavior_test.go` - Snapshot creation and queries
- `tests/symbol_behavior_test.go` - Symbol CRUD and activation
- `tests/service_discovery_behavior_test.go` - Service registration tests
- `tests/cache_behavior_test.go` - Cache operations tests
- `tests/integration_behavior_test.go` - Cross-repository consistency tests
- `tests/test_utils.go` - Test utilities and factories

#### Makefile Test Automation:
```makefile
.PHONY: test test-quick test-price test-candle test-snapshot test-symbol test-service test-cache test-integration test-all test-coverage check-env

# Load .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

check-env:
	@if [ ! -f .env ]; then \
		echo "Warning: .env not found. Copy .env.example to .env"; \
		exit 1; \
	fi

test-quick:
	@if [ -f .env ]; then set -a && . ./.env && set +a; fi && \
	go test -v ./tests -run TestPriceFeedBehavior -timeout=2m

test-price: check-env
	@set -a && . ./.env && set +a && \
	go test -v ./tests -run TestPriceFeedBehaviorSuite -timeout=5m

test-candle: check-env
	@set -a && . ./.env && set +a && \
	go test -v ./tests -run TestCandleBehaviorSuite -timeout=5m

test-snapshot: check-env
	@set -a && . ./.env && set +a && \
	go test -v ./tests -run TestMarketSnapshotBehaviorSuite -timeout=5m

test-symbol: check-env
	@set -a && . ./.env && set +a && \
	go test -v ./tests -run TestSymbolBehaviorSuite -timeout=5m

test-service: check-env
	@set -a && . ./.env && set +a && \
	go test -v ./tests -run TestServiceDiscoveryBehaviorSuite -timeout=5m

test-cache: check-env
	@set -a && . ./.env && set +a && \
	go test -v ./tests -run TestCacheBehaviorSuite -timeout=5m

test-integration: check-env
	@set -a && . ./.env && set +a && \
	go test -v ./tests -run TestIntegrationBehaviorSuite -timeout=10m

test-all: check-env
	@set -a && . ./.env && set +a && \
	go test -v ./tests -timeout=15m

test-coverage: check-env
	@set -a && . ./.env && set +a && \
	go test -v ./tests -coverprofile=coverage.out -timeout=15m
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

build:
	go build -v ./...

clean:
	rm -f coverage.out coverage.html
	go clean -testcache
```

**Test Coverage Goals**:
- Price feed operations: Create, query, latest price, cleanup
- Candle operations: Upsert, time-series queries, interval handling
- Snapshot operations: Create, query, latest snapshot
- Symbol operations: CRUD, activation, active symbol queries
- Service discovery: Registration, heartbeat, cleanup
- Cache operations: Set/Get, TTL, pattern operations
- Integration: Price feed → candle aggregation → snapshot generation

**Acceptance Criteria**:
- [ ] BDD test framework established
- [ ] 25+ test scenarios covering all repositories
- [ ] Performance tests with configurable thresholds
- [ ] 80%+ average test pass rate
- [ ] CI/CD adaptation (SKIP_INTEGRATION_TESTS)
- [ ] Automatic .env loading in tests
- [ ] Time-series query testing

---

### Task 8: Documentation
**Goal**: Create comprehensive documentation for developers
**Estimated Time**: 1 hour

#### README.md:
- Overview of market data adapter
- Architecture and repository pattern
- Installation and setup instructions
- Usage examples for all repositories
- Testing guide
- Environment configuration reference
- Data retention and cleanup strategies

#### tests/README.md:
- Testing framework overview
- How to run different test suites
- Environment setup for tests
- CI/CD integration
- Performance testing configuration

**Acceptance Criteria**:
- [ ] README.md with complete usage guide
- [ ] tests/README.md with testing instructions
- [ ] Code examples for all repositories
- [ ] Environment configuration documented
- [ ] Data retention policy documented

---

## 📊 Success Metrics

| Metric | Target | Status |
|--------|--------|--------|
| Repository Interfaces | 6 | ⏳ Pending |
| PostgreSQL Tables | 4 | ⏳ Pending |
| Test Scenarios | 25+ | ⏳ Pending |
| Test Pass Rate | 80%+ | ⏳ Pending |
| Code Coverage | 70%+ | ⏳ Pending |
| Build Status | Pass | ⏳ Pending |
| Documentation | Complete | ⏳ Pending |

---

## 🔧 Validation Commands

### Environment Setup
```bash
# Copy environment template
cp .env.example .env

# Edit with orchestrator credentials
vim .env

# Validate environment
make check-env
```

### Testing
```bash
# Quick smoke test
make test-quick

# Individual test suites
make test-price
make test-candle
make test-snapshot
make test-symbol
make test-service
make test-cache

# All tests
make test-all

# With coverage
make test-coverage
```

### Build Validation
```bash
# Build all packages
make build

# Run example
go run cmd/example/main.go
```

---

## 🚀 Integration with market-data-simulator-go

Once complete, market-data-simulator-go will integrate by:
1. Adding dependency: `require github.com/quantfidential/trading-ecosystem/market-data-adapter-go v0.1.0`
2. Using `replace` directive for local development
3. Initializing adapter in config layer
4. Using repository interfaces in service layer
5. Following audit-correlator-go integration pattern

---

## ✅ Completion Checklist

- [ ] All 8 tasks completed
- [ ] Build passes without errors
- [ ] 25+ test scenarios passing (80%+ success rate)
- [ ] Documentation complete
- [ ] Example code working
- [ ] Ready for market-data-simulator-go integration

---

**Epic**: TSE-0001 Foundation Services & Infrastructure
**Milestone**: TSE-0001.4 Data Adapters & Orchestrator Integration
**Status**: 📝 READY TO START
**Pattern**: Following audit-data-adapter-go, custodian-data-adapter-go, and exchange-data-adapter-go proven approach
**Estimated Completion**: 8-10 hours following established pattern

**Last Updated**: 2025-09-30
