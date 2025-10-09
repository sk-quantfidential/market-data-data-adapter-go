package tests

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"

	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/pkg/models"
)

// CandleBehaviorTestSuite tests the behavior of candle repository operations
type CandleBehaviorTestSuite struct {
	BehaviorTestSuite
}

// TestCandleBehaviorSuite runs the candle behavior test suite
func TestCandleBehaviorSuite(t *testing.T) {
	suite.Run(t, new(CandleBehaviorTestSuite))
}

// TestCandleCRUDOperations tests basic candle CRUD operations
func (suite *CandleBehaviorTestSuite) TestCandleCRUDOperations() {
	var candleID = GenerateTestUUID()

	suite.Given("a new candle to create", func() {
		// Candle defined below
	}).When("creating the candle", func() {
		candle := suite.CreateTestCandle(candleID, func(c *models.Candle) {
			c.Symbol = "BTC-USD"
			c.Interval = models.Interval1m
			c.Open = decimal.NewFromFloat(50000.00)
			c.Close = decimal.NewFromFloat(50050.00)
		})

		err := suite.adapter.CandleRepository().Create(suite.ctx, candle)
		suite.Require().NoError(err)
		suite.trackCreatedCandle(candleID)
	}).Then("the candle should be retrievable", func() {
		retrieved, err := suite.adapter.CandleRepository().GetByID(suite.ctx, candleID)
		suite.Require().NoError(err)
		suite.Equal(candleID, retrieved.CandleID)
		suite.Equal("BTC-USD", retrieved.Symbol)
		suite.Equal(models.Interval1m, retrieved.Interval)
	})
}

// TestCandleQueryBySymbolAndInterval tests querying candles by symbol and interval
func (suite *CandleBehaviorTestSuite) TestCandleQueryBySymbolAndInterval() {
	var (
		candleID1 = GenerateTestUUID()
		candleID2 = GenerateTestUUID()
	)

	suite.Given("candles with different intervals", func() {
		candle1 := suite.CreateTestCandle(candleID1, func(c *models.Candle) {
			c.Symbol = "ETH-USD"
			c.Interval = models.Interval1m
		})
		err := suite.adapter.CandleRepository().Create(suite.ctx, candle1)
		suite.Require().NoError(err)
		suite.trackCreatedCandle(candleID1)

		candle2 := suite.CreateTestCandle(candleID2, func(c *models.Candle) {
			c.Symbol = "ETH-USD"
			c.Interval = models.Interval5m
		})
		err = suite.adapter.CandleRepository().Create(suite.ctx, candle2)
		suite.Require().NoError(err)
		suite.trackCreatedCandle(candleID2)
	}).When("querying candles by symbol and 1m interval", func() {
		candles, err := suite.adapter.CandleRepository().GetBySymbolAndInterval(suite.ctx, "ETH-USD", models.Interval1m)
		suite.Require().NoError(err)

		suite.Then("only 1m interval candles should be returned", func() {
			suite.GreaterOrEqual(len(candles), 1)

			for _, candle := range candles {
				if candle.CandleID == candleID1 || candle.CandleID == candleID2 {
					suite.Equal(models.Interval1m, candle.Interval)
				}
			}
		})
	})
}

// TestCandleTimeSeriesQuery tests querying candles by time range
func (suite *CandleBehaviorTestSuite) TestCandleTimeSeriesQuery() {
	var (
		candleID1 = GenerateTestUUID()
		candleID2 = GenerateTestUUID()
		candleID3 = GenerateTestUUID()
	)

	now := time.Now().Truncate(time.Minute)

	suite.Given("candles across different time periods", func() {
		// Old candle
		candle1 := suite.CreateTestCandle(candleID1, func(c *models.Candle) {
			c.Symbol = "BTC-USD"
			c.Interval = models.Interval1h
			c.StartTime = now.Add(-3 * time.Hour)
			c.EndTime = now.Add(-2 * time.Hour)
		})
		err := suite.adapter.CandleRepository().Create(suite.ctx, candle1)
		suite.Require().NoError(err)
		suite.trackCreatedCandle(candleID1)

		// Recent candle
		candle2 := suite.CreateTestCandle(candleID2, func(c *models.Candle) {
			c.Symbol = "BTC-USD"
			c.Interval = models.Interval1h
			c.StartTime = now.Add(-1 * time.Hour)
			c.EndTime = now
		})
		err = suite.adapter.CandleRepository().Create(suite.ctx, candle2)
		suite.Require().NoError(err)
		suite.trackCreatedCandle(candleID2)

		// Future candle
		candle3 := suite.CreateTestCandle(candleID3, func(c *models.Candle) {
			c.Symbol = "BTC-USD"
			c.Interval = models.Interval1h
			c.StartTime = now.Add(1 * time.Hour)
			c.EndTime = now.Add(2 * time.Hour)
		})
		err = suite.adapter.CandleRepository().Create(suite.ctx, candle3)
		suite.Require().NoError(err)
		suite.trackCreatedCandle(candleID3)
	}).When("querying candles within a time range", func() {
		startTime := now.Add(-2 * time.Hour)
		endTime := now.Add(30 * time.Minute)

		candles, err := suite.adapter.CandleRepository().GetByTimeRange(
			suite.ctx,
			"BTC-USD",
			models.Interval1h,
			startTime,
			endTime,
		)
		suite.Require().NoError(err)

		suite.Then("only candles within the time range should be returned", func() {
			suite.GreaterOrEqual(len(candles), 1)

			foundRecent := false
			foundOld := false
			foundFuture := false

			for _, candle := range candles {
				if candle.CandleID == candleID2 {
					foundRecent = true
				}
				if candle.CandleID == candleID1 {
					foundOld = true
				}
				if candle.CandleID == candleID3 {
					foundFuture = true
				}
			}

			suite.True(foundRecent, "Recent candle should be in range")
			suite.False(foundOld, "Old candle should not be in range")
			suite.False(foundFuture, "Future candle should not be in range")
		})
	})
}
