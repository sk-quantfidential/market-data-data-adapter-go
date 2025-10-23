package tests

import (
	"time"

	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/pkg/models"
)

// Common behavior scenarios that can be reused across test suites

// symbolLifecycleScenario tests the complete lifecycle of a symbol
func (suite *BehaviorTestSuite) symbolLifecycleScenario() {
	var (
		symbolID = GenerateTestUUID()
		symbol   *models.Symbol
		err      error
	)

	suite.Given("a symbol with active status", func() {
		symbol = suite.CreateTestSymbol(symbolID, func(s *models.Symbol) {
			s.IsActive = true
		})
	}).When("the symbol is created in the repository", func() {
		err = suite.adapter.SymbolRepository().Create(suite.ctx, symbol)
		suite.Require().NoError(err)
		suite.trackCreatedSymbol(symbolID)
	}).Then("the symbol should be retrievable", func() {
		retrievedSymbol, getErr := suite.adapter.SymbolRepository().GetByID(suite.ctx, symbolID)
		suite.Require().NoError(getErr)
		suite.Require().NotNil(retrievedSymbol)
		suite.Equal(symbolID, retrievedSymbol.SymbolID)
		suite.True(retrievedSymbol.IsActive)
	}).And("the symbol can be deactivated", func() {
		symbol.IsActive = false
		updateErr := suite.adapter.SymbolRepository().Update(suite.ctx, symbol)
		suite.Require().NoError(updateErr)
	}).And("the updated status should be persisted", func() {
		updatedSymbol, getErr := suite.adapter.SymbolRepository().GetByID(suite.ctx, symbolID)
		suite.Require().NoError(getErr)
		suite.False(updatedSymbol.IsActive)
	})
}

// priceFeedLifecycleScenario tests the complete lifecycle of price feeds
func (suite *BehaviorTestSuite) priceFeedLifecycleScenario() {
	var (
		feedID    = GenerateTestUUID()
		priceFeed *models.PriceFeed
		err       error
	)

	suite.Given("a price feed for BTC-USD", func() {
		priceFeed = suite.CreateTestPriceFeed(feedID, func(pf *models.PriceFeed) {
			pf.Symbol = "BTC-USD"
		})
	}).When("the price feed is created", func() {
		err = suite.adapter.PriceFeedRepository().Create(suite.ctx, priceFeed)
		suite.Require().NoError(err)
		suite.trackCreatedPriceFeed(feedID)
	}).Then("the price feed should be retrievable", func() {
		retrievedFeed, getErr := suite.adapter.PriceFeedRepository().GetByID(suite.ctx, feedID)
		suite.Require().NoError(getErr)
		suite.Require().NotNil(retrievedFeed)
		suite.Equal(feedID, retrievedFeed.FeedID)
		suite.Equal("BTC-USD", retrievedFeed.Symbol)
	}).And("the price feed should appear in symbol queries", func() {
		feeds, listErr := suite.adapter.PriceFeedRepository().GetBySymbol(suite.ctx, "BTC-USD", 100)
		suite.Require().NoError(listErr)
		suite.Require().NotEmpty(feeds)

		found := false
		for _, feed := range feeds {
			if feed.FeedID == feedID {
				found = true
				break
			}
		}
		suite.True(found, "Price feed should be found in symbol feeds list")
	})
}

// candleLifecycleScenario tests the complete lifecycle of candles
func (suite *BehaviorTestSuite) candleLifecycleScenario() {
	var (
		candleID = GenerateTestUUID()
		candle   *models.Candle
		err      error
	)

	suite.Given("a 1-minute candle for ETH-USD", func() {
		candle = suite.CreateTestCandle(candleID, func(c *models.Candle) {
			c.Symbol = "ETH-USD"
			c.Interval = models.Interval1m
		})
	}).When("the candle is created", func() {
		err = suite.adapter.CandleRepository().Upsert(suite.ctx, candle)
		suite.Require().NoError(err)
		suite.trackCreatedCandle(candleID)
	}).Then("the candle should be retrievable", func() {
		retrievedCandle, getErr := suite.adapter.CandleRepository().GetByID(suite.ctx, candleID)
		suite.Require().NoError(getErr)
		suite.Require().NotNil(retrievedCandle)
		suite.Equal(candleID, retrievedCandle.CandleID)
		suite.Equal("ETH-USD", retrievedCandle.Symbol)
		suite.Equal(models.Interval1m, retrievedCandle.Interval)
	}).And("the candle should appear in time-series queries", func() {
		interval := models.Interval1m
		query := &models.CandleQuery{
			Symbol:        stringPtr("ETH-USD"),
			Interval:      &interval,
			StartTimeFrom: &candle.StartTime,
			StartTimeTo:   &candle.EndTime,
		}
		candles, listErr := suite.adapter.CandleRepository().Query(suite.ctx, query)
		suite.Require().NoError(listErr)
		suite.Require().NotEmpty(candles)

		found := false
		for _, c := range candles {
			if c.CandleID == candleID {
				found = true
				break
			}
		}
		suite.True(found, "Candle should be found in time-series query")
	})
}

// marketSnapshotLifecycleScenario tests the complete lifecycle of market snapshots
func (suite *BehaviorTestSuite) marketSnapshotLifecycleScenario() {
	var (
		snapshotID = GenerateTestUUID()
		snapshot   *models.MarketSnapshot
		err        error
	)

	suite.Given("a market snapshot for SOL-USD", func() {
		snapshot = suite.CreateTestMarketSnapshot(snapshotID, func(ms *models.MarketSnapshot) {
			ms.Symbol = "SOL-USD"
		})
	}).When("the snapshot is created", func() {
		err = suite.adapter.MarketSnapshotRepository().Create(suite.ctx, snapshot)
		suite.Require().NoError(err)
		suite.trackCreatedMarketSnapshot(snapshotID)
	}).Then("the snapshot should be retrievable", func() {
		retrievedSnapshot, getErr := suite.adapter.MarketSnapshotRepository().GetByID(suite.ctx, snapshotID)
		suite.Require().NoError(getErr)
		suite.Require().NotNil(retrievedSnapshot)
		suite.Equal(snapshotID, retrievedSnapshot.SnapshotID)
		suite.Equal("SOL-USD", retrievedSnapshot.Symbol)
	}).And("the snapshot should be the latest for the symbol", func() {
		latestSnapshot, getErr := suite.adapter.MarketSnapshotRepository().GetLatestBySymbol(suite.ctx, "SOL-USD")
		suite.Require().NoError(getErr)
		suite.Equal(snapshotID, latestSnapshot.SnapshotID)
	})
}

// serviceDiscoveryLifecycleScenario tests the complete lifecycle of service registration
func (suite *BehaviorTestSuite) serviceDiscoveryLifecycleScenario() {
	var (
		serviceID = GenerateTestID("service")
		err       error
	)

	suite.Given("a service registration with healthy status", func() {
		// Service defined below
	}).When("the service is registered", func() {
		service := suite.CreateTestServiceRegistration(serviceID, func(s *models.ServiceRegistration) {
			s.Status = "healthy"
			s.Name = "test-lifecycle-service"
		})
		err = suite.adapter.RegisterService(suite.ctx, service)
		suite.Require().NoError(err)
		suite.trackCreatedService(serviceID)
	}).Then("the service should be discoverable", func() {
		retrievedService, getErr := suite.adapter.ServiceDiscoveryRepository().GetServiceInfo(suite.ctx, serviceID)
		suite.Require().NoError(getErr)
		suite.Require().NotNil(retrievedService)
		suite.Equal(serviceID, retrievedService.ID)
		suite.Equal("healthy", retrievedService.Status)
	}).And("the service should appear in service list by name", func() {
		services, listErr := suite.adapter.ServiceDiscoveryRepository().Discover(suite.ctx, "test-lifecycle-service")
		suite.Require().NoError(listErr)
		suite.Require().NotEmpty(services)

		found := false
		for _, s := range services {
			if s.ID == serviceID {
				found = true
				break
			}
		}
		suite.True(found, "Service should be found in services list")
	})
}

// cacheOperationsScenario tests cache functionality
func (suite *BehaviorTestSuite) cacheOperationsScenario() {
	var (
		key   = "test:cache:" + GenerateTestID("key")
		value = map[string]interface{}{
			"test_field": "test_value",
			"numeric":    42,
			"boolean":    true,
		}
		ttl = 30 * time.Second
	)

	suite.Given("a cache key-value pair", func() {
		// Key and value are defined above
	}).When("the value is stored in cache with TTL", func() {
		err := suite.adapter.CacheRepository().Set(suite.ctx, key, value, ttl)
		suite.Require().NoError(err)
	}).Then("the value should be retrievable from cache", func() {
		var retrieved map[string]interface{}
		err := suite.adapter.CacheRepository().Get(suite.ctx, key, &retrieved)
		suite.Require().NoError(err)
		suite.Equal(value["test_field"], retrieved["test_field"])
		suite.Equal(float64(42), retrieved["numeric"]) // JSON unmarshaling converts numbers to float64
		suite.Equal(true, retrieved["boolean"])
	}).And("the cache should confirm the key exists", func() {
		exists, err := suite.adapter.CacheRepository().Exists(suite.ctx, key)
		suite.Require().NoError(err)
		suite.True(exists)
	}).And("the key can be deleted from cache", func() {
		err := suite.adapter.CacheRepository().Delete(suite.ctx, key)
		suite.Require().NoError(err)
	}).And("the deleted key should not exist", func() {
		exists, err := suite.adapter.CacheRepository().Exists(suite.ctx, key)
		suite.Require().NoError(err)
		suite.False(exists)
	})
}

// Helper function to create string pointer for optional query fields
func stringPtr(s string) *string {
	return &s
}
