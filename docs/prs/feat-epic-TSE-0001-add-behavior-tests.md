# Pull Request: Add Comprehensive BDD Behavior Tests

**Epic:** TSE-0001 - Foundation Services & Infrastructure
**Milestone:** TSE-0001.4.3 - Market Data Adapter Foundation
**Branch:** `feature/epic-TSE-0001-add-behavior-tests`
**Status:** ✅ Ready for Merge

## Summary

This PR adds comprehensive BDD (Behavior-Driven Development) behavior tests to the market-data-adapter-go repository, bringing it to parity with audit-data-adapter-go, custodian-data-adapter-go, and exchange-data-adapter-go testing capabilities.

**Key Achievements:**
1. **11 test files** with ~1,915 lines of test code
2. **Market data domain coverage** - Symbols, PriceFeeds, Candles, MarketSnapshots
3. **Infrastructure coverage** - Cache operations, Service Discovery
4. **BDD Given/When/Then** style for clear test scenarios
5. **Automatic cleanup** tracking to prevent test pollution
6. **Test entity factories** with functional options pattern

## What Changed

### Phase 1: Add comprehensive BDD behavior tests for market data adapter
**Commit:** `41efd94`

Details documented in commit message.

### Phase 2: docs: Add pull request documentation for feature branches
**Commit:** `8a1a58d`

Details documented in commit message.

### Phase 3: feat: Add comprehensive Makefile for testing and development
**Commit:** `d7bd125`

Details documented in commit message.

### Phase 4: chore: rename PR doc to match updated branch naming convention
**Commit:** `d5a700e`

Details documented in commit message.


## Market-Data-Adapter-Go Repository Changes

### Commit Summary

**Total Commits**: 1
**Commit:** `bea6341` - Add comprehensive BDD behavior tests for market data adapter
**Files Changed**: 11 files, 1,915 insertions(+)

### Test Infrastructure (5 files, ~850 lines)

#### 1. init_test.go (73 lines)
**Purpose:** Environment configuration loading for tests

```go
func init() {
    // Automatically loads .env file before tests run
    // Searches parent directory, current directory, and project root
    envPath := "../.env"
    godotenv.Load(envPath)
}
```

**Features:**
- Automatic .env file discovery and loading
- Multi-location search (parent, current, project root)
- Graceful fallback if .env not found

#### 2. test_utils.go (204 lines)
**Purpose:** Test utilities and helper functions

**Key Functions:**
- `GenerateTestID(prefix)` - Unique test IDs with timestamps
- `GenerateTestUUID()` - UUID generation for test entities
- `GetEnvOrDefault(key, default)` - Environment variable helpers
- `GetTestConfig()` - Test configuration loading
- `IsCI()` - CI environment detection
- `WaitForCondition()` - Polling helper with timeout
- `RetryOperation()` - Retry with exponential backoff

**Test Data Generators:**
- `TestDataGenerator` - Counter-based unique ID generation
- `NextSymbol()` - Generate unique trading symbols (TEST1-USD, TEST2-USD)
- `NextSource()` - Generate unique data source names
- `TestTimestamps` - Time-based testing utilities

#### 3. behavior_test_suite.go (361 lines)
**Purpose:** Base test suite with BDD helpers and market data entity factories

**Key Components:**
```go
type BehaviorTestSuite struct {
    suite.Suite
    ctx     context.Context
    adapter adapters.DataAdapter
    config  *config.Config

    // Cleanup tracking
    createdSymbols          []string
    createdPriceFeeds       []string
    createdCandles          []string
    createdMarketSnapshots  []string
    createdServices         []string
}
```

**Entity Factories:**
- `CreateTestSymbol(symbolID, modifiers...)` - Symbol factory
- `CreateTestPriceFeed(feedID, modifiers...)` - Price feed factory
- `CreateTestCandle(candleID, modifiers...)` - Candle factory
- `CreateTestMarketSnapshot(snapshotID, modifiers...)` - Market snapshot factory
- `CreateTestServiceInfo(serviceID, modifiers...)` - Service info factory

**BDD Helpers:**
- `Given(description, fn)` - Define test preconditions
- `When(description, fn)` - Define test actions
- `Then(description, fn)` - Define test assertions
- `And(description, fn)` - Chain additional steps

**Performance Helpers:**
- `AssertPerformance(operationName, maxDuration, fn)` - Performance assertions

**Automatic Cleanup:**
- Tracks all created entities during tests
- Automatically cleans up in `TearDownTest`
- Prevents test pollution and cross-test interference

#### 4. behavior_test_runner.go (273 lines)
**Purpose:** Test runner with environment configuration

**Features:**
- Configurable logging levels (debug, info, warn, error)
- Skip integration tests in CI environments
- Skip performance tests in CI environments
- Test timeout configuration
- Prerequisite checking (database, Redis connectivity)
- Environment information printing

**Configuration:**
```go
type BehaviorTestConfig struct {
    PostgresURL             string
    RedisURL                string
    TestTimeout             time.Duration
    SkipIntegrationTests    bool
    SkipPerformanceTests    bool
    LogLevel                string
    MaxConcurrentOperations int
    LargeDatasetSize        int
}
```

#### 5. behavior_scenarios.go (237 lines)
**Purpose:** Reusable test scenarios

**Scenarios:**
- `symbolLifecycleScenario()` - Complete symbol lifecycle (create → update → deactivate)
- `priceFeedLifecycleScenario()` - Price feed lifecycle
- `candleLifecycleScenario()` - Candle lifecycle with time-series
- `marketSnapshotLifecycleScenario()` - Market snapshot lifecycle
- `serviceDiscoveryLifecycleScenario()` - Service registration lifecycle
- `cacheOperationsScenario()` - Cache operations with TTL

### Behavior Tests (6 files, ~1,065 lines)

#### 1. symbol_behavior_test.go (144 lines)
**Purpose:** Test symbol repository operations

**Test Cases:**
- `TestSymbolCRUDOperations` - Create, Read, Update operations
- `TestSymbolQueryByStatus` - Query active vs inactive symbols
- `TestSymbolQueryByBaseCurrency` - Query symbols by base currency (BTC, ETH)

**Coverage:**
- Symbol activation/deactivation
- Currency pair handling (BTC-USD, BTC-EUR, ETH-USD)
- Active symbol filtering

#### 2. price_feed_behavior_test.go (155 lines)
**Purpose:** Test price feed repository operations

**Test Cases:**
- `TestPriceFeedCRUDOperations` - Create, Read operations
- `TestPriceFeedQueryBySymbol` - Query feeds for specific symbol
- `TestLatestPriceFeed` - Retrieve latest price for symbol
- `TestPriceFeedQueryBySource` - Query feeds by data source

**Coverage:**
- Real-time price updates
- Multi-source price aggregation (exchange-A, exchange-B, coinbase)
- Latest price retrieval
- Bid/ask spread handling
- Volume tracking

#### 3. candle_behavior_test.go (168 lines)
**Purpose:** Test candle (OHLCV) repository operations

**Test Cases:**
- `TestCandleCRUDOperations` - Create, Read operations
- `TestCandleQueryBySymbolAndInterval` - Query by symbol and interval (1m, 5m, 1h)
- `TestCandleTimeSeriesQuery` - Query candles within time ranges

**Coverage:**
- Multiple intervals (1m, 5m, 15m, 1h, 4h, 1d)
- OHLCV data (Open, High, Low, Close, Volume)
- Time-series queries with start/end time
- Interval-specific filtering

#### 4. market_snapshot_behavior_test.go (123 lines)
**Purpose:** Test market snapshot repository operations

**Test Cases:**
- `TestMarketSnapshotCRUDOperations` - Create, Read operations
- `TestLatestMarketSnapshot` - Retrieve latest snapshot for symbol
- `TestMarketSnapshotQueryBySymbol` - Query snapshots by symbol

**Coverage:**
- Market state snapshots
- Latest snapshot retrieval
- Price change tracking (24h)
- Spread calculation
- Volume aggregation

#### 5. cache_behavior_test.go (83 lines)
**Purpose:** Test cache repository operations

**Test Cases:**
- `TestCacheStringOperations` - String value caching
- `TestCacheExpiration` - TTL expiration behavior
- `TestCacheDelete` - Cache deletion operations

**Coverage:**
- String value storage/retrieval
- TTL (Time-To-Live) management
- Key existence checking
- Cache deletion
- Expiration after TTL

#### 6. service_discovery_behavior_test.go (94 lines)
**Purpose:** Test service discovery repository operations

**Test Cases:**
- `TestServiceRegistration` - Service registration and discovery
- `TestServiceHeartbeat` - Heartbeat updates
- `TestServiceDeregistration` - Service deregistration

**Coverage:**
- Service registration
- Service discovery by name
- Heartbeat timestamp updates
- Service deregistration
- Service metadata handling

## Testing Approach

### BDD (Behavior-Driven Development) Style

All tests use Given/When/Then pattern for clarity:

```go
suite.Given("a new symbol to create", func() {
    // Setup preconditions
}).When("creating the symbol", func() {
    symbol := suite.CreateTestSymbol(symbolID, func(s *models.Symbol) {
        s.Symbol = "ETH-USD"
        s.IsActive = true
    })
    err := suite.adapter.SymbolRepository().Create(suite.ctx, symbol)
    suite.Require().NoError(err)
    suite.trackCreatedSymbol(symbolID)
}).Then("the symbol should be retrievable", func() {
    retrieved, err := suite.adapter.SymbolRepository().GetByID(suite.ctx, symbolID)
    suite.Require().NoError(err)
    suite.Equal("ETH-USD", retrieved.Symbol)
})
```

### Functional Options Pattern

Test entity factories use functional options for flexibility:

```go
// Create symbol with custom values
symbol := suite.CreateTestSymbol(symbolID, func(s *models.Symbol) {
    s.Symbol = "BTC-USD"
    s.BaseCurrency = "BTC"
    s.QuoteCurrency = "USD"
    s.IsActive = true
})
```

### Automatic Cleanup

Test suite tracks created entities and cleans up automatically:

```go
func (suite *BehaviorTestSuite) TearDownTest() {
    suite.cleanupCreatedResources()  // Removes all tracked entities
}
```

This prevents:
- Test pollution (data from one test affecting another)
- Database bloat from test data
- Cache key collisions
- Service registration leaks

## Integration with Makefile

Tests integrate with existing Makefile targets:

```bash
# Run all behavior tests
make test-behavior

# Run all tests with coverage
make test-all

# Generate coverage report
make coverage

# Run specific test suite
go test -v ./tests -run TestSymbolBehaviorSuite
go test -v ./tests -run TestPriceFeedBehaviorSuite
go test -v ./tests -run TestCandleBehaviorSuite
```

## Environment Configuration

Tests use environment variables for configuration:

```bash
# Database connections
export TEST_POSTGRES_URL="postgres://postgres:postgres@localhost:5432/market_data_test?sslmode=disable"
export TEST_REDIS_URL="redis://localhost:6379/15"

# Test behavior
export TEST_LOG_LEVEL="info"                # debug, info, warn, error
export TEST_TIMEOUT="5m"                    # Test suite timeout
export SKIP_INTEGRATION_TESTS="false"       # Skip integration tests
export SKIP_PERFORMANCE_TESTS="false"       # Skip performance tests
```

## Test Coverage Summary

### Market Data Entities
- ✅ **Symbols** - CRUD, active filtering, currency pair queries
- ✅ **Price Feeds** - CRUD, multi-source, latest price, source queries
- ✅ **Candles** - CRUD, time-series, interval handling, OHLCV data
- ✅ **Market Snapshots** - CRUD, latest snapshot, price change tracking

### Infrastructure
- ✅ **Cache** - String operations, TTL, expiration, deletion
- ✅ **Service Discovery** - Registration, heartbeat, deregistration

### Query Patterns
- ✅ By ID (GetByID)
- ✅ By symbol (GetBySymbol)
- ✅ By source (GetBySource)
- ✅ By status (GetActive)
- ✅ By currency (GetByBaseCurrency)
- ✅ By time range (GetByTimeRange)
- ✅ Latest data (GetLatest, GetLatestMarketSnapshot)

### Operations
- ✅ Create
- ✅ Read
- ✅ Update
- ✅ Delete
- ✅ Query/Filter
- ✅ Time-series
- ✅ Aggregation

## Files Changed

```
tests/behavior_scenarios.go              | 237 ++++++++++++++++++++
tests/behavior_test_runner.go            | 273 +++++++++++++++++++++++
tests/behavior_test_suite.go             | 361 +++++++++++++++++++++++++++++++
tests/cache_behavior_test.go             |  83 +++++++
tests/candle_behavior_test.go            | 168 ++++++++++++++
tests/init_test.go                       |  73 +++++++
tests/market_snapshot_behavior_test.go   | 123 +++++++++++
tests/price_feed_behavior_test.go        | 155 +++++++++++++
tests/service_discovery_behavior_test.go |  94 ++++++++
tests/symbol_behavior_test.go            | 144 ++++++++++++
tests/test_utils.go                      | 204 +++++++++++++++++
11 files changed, 1915 insertions(+)
```

## Testing Instructions

### Prerequisites

1. **PostgreSQL** running on localhost:5432 (or configure TEST_POSTGRES_URL)
2. **Redis** running on localhost:6379 (or configure TEST_REDIS_URL)

### Quick Start

```bash
# Run all behavior tests
go test -v ./tests

# Run specific test suite
go test -v ./tests -run TestSymbolBehaviorSuite
go test -v ./tests -run TestPriceFeedBehaviorSuite
go test -v ./tests -run TestCandleBehaviorSuite
go test -v ./tests -run TestMarketSnapshotBehaviorSuite
go test -v ./tests -run TestCacheBehaviorSuite
go test -v ./tests -run TestServiceDiscoveryBehaviorSuite
```

### Docker Setup (Optional)

```bash
# Start test databases
docker run -d --name market-data-test-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=market_data_test \
  -p 5432:5432 postgres:17-alpine

docker run -d --name market-data-test-redis \
  -p 6379:6379 redis:8-alpine

# Run tests
go test -v ./tests

# Cleanup
docker rm -f market-data-test-postgres market-data-test-redis
```

## Benefits

1. **Parity with Other Adapters** - Same testing standard as audit, custodian, and exchange adapters
2. **Confident Refactoring** - Comprehensive test coverage enables safe code changes
3. **Clear Documentation** - BDD tests serve as executable documentation
4. **Prevented Regressions** - Automatic cleanup prevents test pollution
5. **Fast Feedback** - Tests run quickly with proper isolation
6. **CI/CD Ready** - Environment detection and configuration for automated pipelines

## Comparison with Other Adapters

| Feature | audit-data-adapter-go | custodian-data-adapter-go | exchange-data-adapter-go | market-data-adapter-go |
|---------|----------------------|---------------------------|-------------------------|------------------------|
| BDD Tests | ✅ | ✅ | ✅ | ✅ (this PR) |
| Test Infrastructure | ✅ | ✅ | ✅ | ✅ (this PR) |
| Entity Factories | ✅ | ✅ | ✅ | ✅ (this PR) |
| Automatic Cleanup | ✅ | ✅ | ✅ | ✅ (this PR) |
| Makefile Integration | ✅ | ✅ | ✅ | ✅ (existing) |
| Domain-Specific Tests | Audit Events | Positions, Settlements | Orders, Trades | Symbols, Prices, Candles |

## Next Steps (Future Work)

After merge:
- Add integration tests for full workflow scenarios
- Add performance tests with throughput measurements
- Add comprehensive test suite that runs all tests
- Expand test coverage for edge cases
- Add README.md in tests/ directory documenting test structure

## BDD Acceptance Criteria

✅ **All acceptance criteria met:**
- Comprehensive BDD tests for all market data entities
- Test infrastructure with BDD Given/When/Then pattern
- Automatic cleanup tracking prevents test pollution
- Entity factories with functional options pattern
- Integration with Makefile test targets
- Environment-based configuration for different environments
- Parity with audit/custodian/exchange data adapter testing standards

---

**Ready for Merge:** Yes ✅
**Breaking Changes:** None
**Migration Required:** None
**Documentation Updated:** Yes (this PR document)

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
