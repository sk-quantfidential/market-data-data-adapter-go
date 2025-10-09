package tests

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"

	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/pkg/models"
)

// PriceFeedBehaviorTestSuite tests the behavior of price feed repository operations
type PriceFeedBehaviorTestSuite struct {
	BehaviorTestSuite
}

// TestPriceFeedBehaviorSuite runs the price feed behavior test suite
func TestPriceFeedBehaviorSuite(t *testing.T) {
	suite.Run(t, new(PriceFeedBehaviorTestSuite))
}

// TestPriceFeedCRUDOperations tests basic price feed CRUD operations
func (suite *PriceFeedBehaviorTestSuite) TestPriceFeedCRUDOperations() {
	var feedID = GenerateTestUUID()

	suite.Given("a new price feed to create", func() {
		// Price feed defined below
	}).When("creating the price feed", func() {
		feed := suite.CreateTestPriceFeed(feedID, func(pf *models.PriceFeed) {
			pf.Symbol = "BTC-USD"
			pf.Price = decimal.NewFromFloat(50000.00)
			pf.Source = "test-exchange"
		})

		err := suite.adapter.PriceFeedRepository().Create(suite.ctx, feed)
		suite.Require().NoError(err)
		suite.trackCreatedPriceFeed(feedID)
	}).Then("the price feed should be retrievable", func() {
		retrieved, err := suite.adapter.PriceFeedRepository().GetByID(suite.ctx, feedID)
		suite.Require().NoError(err)
		suite.Equal(feedID, retrieved.FeedID)
		suite.Equal("BTC-USD", retrieved.Symbol)
		suite.Equal("test-exchange", retrieved.Source)
	})
}

// TestPriceFeedQueryBySymbol tests querying price feeds by symbol
func (suite *PriceFeedBehaviorTestSuite) TestPriceFeedQueryBySymbol() {
	var (
		feedID1 = GenerateTestUUID()
		feedID2 = GenerateTestUUID()
	)

	suite.Given("multiple price feeds for a symbol", func() {
		feed1 := suite.CreateTestPriceFeed(feedID1, func(pf *models.PriceFeed) {
			pf.Symbol = "ETH-USD"
			pf.Source = "exchange-A"
			pf.Price = decimal.NewFromFloat(3000.00)
		})
		err := suite.adapter.PriceFeedRepository().Create(suite.ctx, feed1)
		suite.Require().NoError(err)
		suite.trackCreatedPriceFeed(feedID1)

		feed2 := suite.CreateTestPriceFeed(feedID2, func(pf *models.PriceFeed) {
			pf.Symbol = "ETH-USD"
			pf.Source = "exchange-B"
			pf.Price = decimal.NewFromFloat(3001.00)
		})
		err = suite.adapter.PriceFeedRepository().Create(suite.ctx, feed2)
		suite.Require().NoError(err)
		suite.trackCreatedPriceFeed(feedID2)
	}).When("querying price feeds by symbol", func() {
		feeds, err := suite.adapter.PriceFeedRepository().GetBySymbol(suite.ctx, "ETH-USD")
		suite.Require().NoError(err)

		suite.Then("all feeds for the symbol should be returned", func() {
			suite.GreaterOrEqual(len(feeds), 2)

			sources := make(map[string]bool)
			for _, feed := range feeds {
				if feed.Symbol == "ETH-USD" {
					sources[feed.Source] = true
				}
			}
			suite.True(sources["exchange-A"])
			suite.True(sources["exchange-B"])
		})
	})
}

// TestLatestPriceFeed tests retrieving the latest price feed for a symbol
func (suite *PriceFeedBehaviorTestSuite) TestLatestPriceFeed() {
	var (
		feedID1 = GenerateTestUUID()
		feedID2 = GenerateTestUUID()
	)

	suite.Given("multiple price feeds with different timestamps", func() {
		feed1 := suite.CreateTestPriceFeed(feedID1, func(pf *models.PriceFeed) {
			pf.Symbol = "SOL-USD"
			pf.Price = decimal.NewFromFloat(100.00)
		})
		feed1.Timestamp = feed1.Timestamp.Add(-1 * suite.Hour)
		err := suite.adapter.PriceFeedRepository().Create(suite.ctx, feed1)
		suite.Require().NoError(err)
		suite.trackCreatedPriceFeed(feedID1)

		feed2 := suite.CreateTestPriceFeed(feedID2, func(pf *models.PriceFeed) {
			pf.Symbol = "SOL-USD"
			pf.Price = decimal.NewFromFloat(105.00)
		})
		err = suite.adapter.PriceFeedRepository().Create(suite.ctx, feed2)
		suite.Require().NoError(err)
		suite.trackCreatedPriceFeed(feedID2)
	}).When("querying for the latest price feed", func() {
		latest, err := suite.adapter.PriceFeedRepository().GetLatest(suite.ctx, "SOL-USD")
		suite.Require().NoError(err)

		suite.Then("the most recent feed should be returned", func() {
			suite.Equal(feedID2, latest.FeedID)
			expectedPrice := decimal.NewFromFloat(105.00)
			suite.True(latest.Price.Equal(expectedPrice))
		})
	})
}

// TestPriceFeedQueryBySource tests querying price feeds by source
func (suite *PriceFeedBehaviorTestSuite) TestPriceFeedQueryBySource() {
	var feedID = GenerateTestUUID()

	suite.Given("a price feed from a specific source", func() {
		feed := suite.CreateTestPriceFeed(feedID, func(pf *models.PriceFeed) {
			pf.Source = "coinbase"
			pf.Symbol = "ADA-USD"
		})
		err := suite.adapter.PriceFeedRepository().Create(suite.ctx, feed)
		suite.Require().NoError(err)
		suite.trackCreatedPriceFeed(feedID)
	}).When("querying price feeds by source", func() {
		feeds, err := suite.adapter.PriceFeedRepository().GetBySource(suite.ctx, "coinbase")
		suite.Require().NoError(err)

		suite.Then("feeds from that source should be found", func() {
			suite.GreaterOrEqual(len(feeds), 1)
			var found bool
			for _, feed := range feeds {
				if feed.FeedID == feedID {
					found = true
					break
				}
			}
			suite.True(found)
		})
	})
}
