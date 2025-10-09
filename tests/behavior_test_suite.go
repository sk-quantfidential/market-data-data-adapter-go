package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"

	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/internal/config"
	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/pkg/adapters"
	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/pkg/interfaces"
	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/pkg/models"
)

// BehaviorTestSuite provides a base test suite with common functionality
type BehaviorTestSuite struct {
	suite.Suite
	ctx     context.Context
	adapter adapters.DataAdapter
	config  *config.Config
	logger  *logrus.Logger

	// Track created resources for cleanup
	createdSymbols          []string
	createdPriceFeeds       []string
	createdCandles          []string
	createdMarketSnapshots  []string
	createdServices         []string
}

// SetupSuite runs once before all tests in the suite
func (suite *BehaviorTestSuite) SetupSuite() {
	suite.logger = logrus.New()
	suite.logger.SetLevel(logrus.WarnLevel)

	// Load test configuration
	suite.config = &config.Config{
		ServiceName:                "market-data-test",
		ServiceInstanceName:        "market-data-test",
		PostgresURL:                GetEnvOrDefault("TEST_POSTGRES_URL", "postgres://postgres:postgres@localhost:5432/market_data_test?sslmode=disable"),
		RedisURL:                   GetEnvOrDefault("TEST_REDIS_URL", "redis://localhost:6379/15"),
		ServiceDiscoveryNamespace:  "test:market_data:discovery",
		CacheNamespace:             "test:market_data:cache",
	}

	// Create adapter
	adapter, err := adapters.NewMarketDataAdapter(suite.config, suite.logger)
	suite.Require().NoError(err, "Failed to create market data adapter")
	suite.adapter = adapter

	// Create context
	suite.ctx = context.Background()

	// Connect to infrastructure
	err = suite.adapter.Connect(suite.ctx)
	suite.Require().NoError(err, "Failed to connect to infrastructure")
}

// TearDownSuite runs once after all tests in the suite
func (suite *BehaviorTestSuite) TearDownSuite() {
	if suite.adapter != nil {
		_ = suite.adapter.Disconnect(suite.ctx)
	}
}

// TearDownTest runs after each test to clean up created resources
func (suite *BehaviorTestSuite) TearDownTest() {
	suite.cleanupCreatedResources()
}

// cleanupCreatedResources removes all tracked test resources
func (suite *BehaviorTestSuite) cleanupCreatedResources() {
	// Cleanup market snapshots
	for _, snapshotID := range suite.createdMarketSnapshots {
		_ = suite.adapter.MarketSnapshotRepository().Delete(suite.ctx, snapshotID)
	}
	suite.createdMarketSnapshots = nil

	// Cleanup candles
	for _, candleID := range suite.createdCandles {
		_ = suite.adapter.CandleRepository().Delete(suite.ctx, candleID)
	}
	suite.createdCandles = nil

	// Cleanup price feeds
	for _, feedID := range suite.createdPriceFeeds {
		_ = suite.adapter.PriceFeedRepository().Delete(suite.ctx, feedID)
	}
	suite.createdPriceFeeds = nil

	// Cleanup symbols
	for _, symbolID := range suite.createdSymbols {
		_ = suite.adapter.SymbolRepository().Delete(suite.ctx, symbolID)
	}
	suite.createdSymbols = nil

	// Cleanup services
	for _, serviceID := range suite.createdServices {
		_ = suite.adapter.ServiceDiscoveryRepository().Deregister(suite.ctx, serviceID)
	}
	suite.createdServices = nil
}

// Tracking methods for cleanup
func (suite *BehaviorTestSuite) trackCreatedSymbol(symbolID string) {
	suite.createdSymbols = append(suite.createdSymbols, symbolID)
}

func (suite *BehaviorTestSuite) trackCreatedPriceFeed(feedID string) {
	suite.createdPriceFeeds = append(suite.createdPriceFeeds, feedID)
}

func (suite *BehaviorTestSuite) trackCreatedCandle(candleID string) {
	suite.createdCandles = append(suite.createdCandles, candleID)
}

func (suite *BehaviorTestSuite) trackCreatedMarketSnapshot(snapshotID string) {
	suite.createdMarketSnapshots = append(suite.createdMarketSnapshots, snapshotID)
}

func (suite *BehaviorTestSuite) trackCreatedService(serviceID string) {
	suite.createdServices = append(suite.createdServices, serviceID)
}

// Test entity factories with functional options pattern

// CreateTestSymbol creates a test symbol with optional modifiers
func (suite *BehaviorTestSuite) CreateTestSymbol(symbolID string, modifiers ...func(*models.Symbol)) *models.Symbol {
	now := time.Now()
	symbol := &models.Symbol{
		SymbolID:      symbolID,
		Symbol:        "BTC-USD",
		BaseCurrency:  "BTC",
		QuoteCurrency: "USD",
		IsActive:      true,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	for _, modifier := range modifiers {
		modifier(symbol)
	}

	return symbol
}

// CreateTestPriceFeed creates a test price feed with optional modifiers
func (suite *BehaviorTestSuite) CreateTestPriceFeed(feedID string, modifiers ...func(*models.PriceFeed)) *models.PriceFeed {
	price := decimal.NewFromFloat(50000.00)
	bid := decimal.NewFromFloat(49999.00)
	ask := decimal.NewFromFloat(50001.00)
	volume := decimal.NewFromFloat(1000000.00)

	feed := &models.PriceFeed{
		FeedID:    feedID,
		Symbol:    "BTC-USD",
		Price:     price,
		Bid:       &bid,
		Ask:       &ask,
		Volume24h: &volume,
		Source:    "test-source",
		Timestamp: time.Now(),
	}

	for _, modifier := range modifiers {
		modifier(feed)
	}

	return feed
}

// CreateTestCandle creates a test candle with optional modifiers
func (suite *BehaviorTestSuite) CreateTestCandle(candleID string, modifiers ...func(*models.Candle)) *models.Candle {
	now := time.Now()
	startTime := now.Truncate(time.Minute)
	endTime := startTime.Add(time.Minute)

	candle := &models.Candle{
		CandleID:  candleID,
		Symbol:    "BTC-USD",
		Interval:  models.Interval1m,
		Open:      decimal.NewFromFloat(50000.00),
		High:      decimal.NewFromFloat(50100.00),
		Low:       decimal.NewFromFloat(49900.00),
		Close:     decimal.NewFromFloat(50050.00),
		Volume:    decimal.NewFromFloat(100.5),
		StartTime: startTime,
		EndTime:   endTime,
	}

	for _, modifier := range modifiers {
		modifier(candle)
	}

	return candle
}

// CreateTestMarketSnapshot creates a test market snapshot with optional modifiers
func (suite *BehaviorTestSuite) CreateTestMarketSnapshot(snapshotID string, modifiers ...func(*models.MarketSnapshot)) *models.MarketSnapshot {
	lastPrice := decimal.NewFromFloat(50000.00)
	bid := decimal.NewFromFloat(49999.00)
	ask := decimal.NewFromFloat(50001.00)
	spread := decimal.NewFromFloat(2.00)
	volume := decimal.NewFromFloat(1000000.00)

	snapshot := &models.MarketSnapshot{
		SnapshotID: snapshotID,
		Symbol:     "BTC-USD",
		LastPrice:  lastPrice,
		Bid:        &bid,
		Ask:        &ask,
		Spread:     &spread,
		Volume24h:  &volume,
		Timestamp:  time.Now(),
	}

	for _, modifier := range modifiers {
		modifier(snapshot)
	}

	return snapshot
}

// CreateTestServiceInfo creates a test service info with optional modifiers
func (suite *BehaviorTestSuite) CreateTestServiceInfo(serviceID string, modifiers ...func(*interfaces.ServiceInfo)) *interfaces.ServiceInfo {
	now := time.Now()
	service := &interfaces.ServiceInfo{
		ServiceName:   "test-service",
		ServiceID:     serviceID,
		Address:       "localhost",
		Port:          8080,
		Version:       "1.0.0",
		Metadata:      map[string]string{"environment": "test"},
		RegisteredAt:  now,
		LastHeartbeat: now,
	}

	for _, modifier := range modifiers {
		modifier(service)
	}

	return service
}

// BDD-style scenario helpers

// BDDScenario represents a behavior-driven test scenario
type BDDScenario struct {
	suite  *BehaviorTestSuite
	givens []BDDStep
	whens  []BDDStep
	thens  []BDDStep
}

// BDDStep represents a single step in a BDD scenario
type BDDStep struct {
	description string
	action      func()
}

// Given starts a new BDD scenario with preconditions
func (suite *BehaviorTestSuite) Given(description string, fn func()) *BDDScenario {
	scenario := &BDDScenario{suite: suite}
	scenario.givens = append(scenario.givens, BDDStep{description, fn})
	return scenario
}

// When adds an action to the scenario
func (scenario *BDDScenario) When(description string, fn func()) *BDDScenario {
	scenario.whens = append(scenario.whens, BDDStep{description, fn})
	return scenario
}

// Then adds an assertion to the scenario
func (scenario *BDDScenario) Then(description string, fn func()) *BDDScenario {
	scenario.thens = append(scenario.thens, BDDStep{description, fn})
	return scenario
}

// And adds an additional step to the current phase
func (scenario *BDDScenario) And(description string, fn func()) *BDDScenario {
	if len(scenario.thens) > 0 {
		scenario.thens = append(scenario.thens, BDDStep{description, fn})
	} else if len(scenario.whens) > 0 {
		scenario.whens = append(scenario.whens, BDDStep{description, fn})
	} else {
		scenario.givens = append(scenario.givens, BDDStep{description, fn})
	}
	return scenario
}

// Execute runs all steps in the scenario
func (scenario *BDDScenario) Execute() {
	// Execute givens
	for _, step := range scenario.givens {
		step.action()
	}

	// Execute whens
	for _, step := range scenario.whens {
		step.action()
	}

	// Execute thens
	for _, step := range scenario.thens {
		step.action()
	}
}

// Performance assertion helpers

// AssertPerformance asserts that an operation completes within a time limit
func (suite *BehaviorTestSuite) AssertPerformance(operationName string, maxDuration time.Duration, fn func()) {
	start := time.Now()
	fn()
	duration := time.Since(start)

	suite.logger.WithFields(logrus.Fields{
		"operation": operationName,
		"duration":  duration,
		"limit":     maxDuration,
	}).Debug("Performance measurement")

	if !GetEnvAsBool("SKIP_PERFORMANCE_TESTS", IsCI()) {
		suite.LessOrEqual(duration, maxDuration, fmt.Sprintf("Operation '%s' took too long: %v > %v", operationName, duration, maxDuration))
	}
}

// Cache operation helpers

// Set stores a value in cache
func (suite *BehaviorTestSuite) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return suite.adapter.CacheRepository().Set(ctx, key, string(data), ttl)
}

// Get retrieves a value from cache
func (suite *BehaviorTestSuite) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := suite.adapter.CacheRepository().Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

// Exists checks if a key exists in cache
func (suite *BehaviorTestSuite) Exists(ctx context.Context, key string) (bool, error) {
	return suite.adapter.CacheRepository().Exists(ctx, key)
}

// Delete removes a key from cache
func (suite *BehaviorTestSuite) Delete(ctx context.Context, key string) error {
	return suite.adapter.CacheRepository().Delete(ctx, key)
}
