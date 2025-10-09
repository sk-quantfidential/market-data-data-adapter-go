package tests

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/suite"

	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/pkg/models"
)

// MarketSnapshotBehaviorTestSuite tests the behavior of market snapshot repository operations
type MarketSnapshotBehaviorTestSuite struct {
	BehaviorTestSuite
}

// TestMarketSnapshotBehaviorSuite runs the market snapshot behavior test suite
func TestMarketSnapshotBehaviorSuite(t *testing.T) {
	suite.Run(t, new(MarketSnapshotBehaviorTestSuite))
}

// TestMarketSnapshotCRUDOperations tests basic market snapshot CRUD operations
func (suite *MarketSnapshotBehaviorTestSuite) TestMarketSnapshotCRUDOperations() {
	var snapshotID = GenerateTestUUID()

	suite.Given("a new market snapshot to create", func() {
		// Snapshot defined below
	}).When("creating the market snapshot", func() {
		snapshot := suite.CreateTestMarketSnapshot(snapshotID, func(ms *models.MarketSnapshot) {
			ms.Symbol = "BTC-USD"
			ms.LastPrice = decimal.NewFromFloat(50000.00)
		})

		err := suite.adapter.MarketSnapshotRepository().Create(suite.ctx, snapshot)
		suite.Require().NoError(err)
		suite.trackCreatedMarketSnapshot(snapshotID)
	}).Then("the market snapshot should be retrievable", func() {
		retrieved, err := suite.adapter.MarketSnapshotRepository().GetByID(suite.ctx, snapshotID)
		suite.Require().NoError(err)
		suite.Equal(snapshotID, retrieved.SnapshotID)
		suite.Equal("BTC-USD", retrieved.Symbol)
	})
}

// TestLatestMarketSnapshot tests retrieving the latest market snapshot for a symbol
func (suite *MarketSnapshotBehaviorTestSuite) TestLatestMarketSnapshot() {
	var (
		snapshotID1 = GenerateTestUUID()
		snapshotID2 = GenerateTestUUID()
	)

	suite.Given("multiple snapshots with different timestamps", func() {
		snapshot1 := suite.CreateTestMarketSnapshot(snapshotID1, func(ms *models.MarketSnapshot) {
			ms.Symbol = "ETH-USD"
			ms.LastPrice = decimal.NewFromFloat(3000.00)
		})
		snapshot1.Timestamp = snapshot1.Timestamp.Add(-1 * time.Hour)
		err := suite.adapter.MarketSnapshotRepository().Create(suite.ctx, snapshot1)
		suite.Require().NoError(err)
		suite.trackCreatedMarketSnapshot(snapshotID1)

		snapshot2 := suite.CreateTestMarketSnapshot(snapshotID2, func(ms *models.MarketSnapshot) {
			ms.Symbol = "ETH-USD"
			ms.LastPrice = decimal.NewFromFloat(3050.00)
		})
		err = suite.adapter.MarketSnapshotRepository().Create(suite.ctx, snapshot2)
		suite.Require().NoError(err)
		suite.trackCreatedMarketSnapshot(snapshotID2)
	}).When("querying for the latest snapshot", func() {
		latest, err := suite.adapter.MarketSnapshotRepository().GetLatest(suite.ctx, "ETH-USD")
		suite.Require().NoError(err)

		suite.Then("the most recent snapshot should be returned", func() {
			suite.Equal(snapshotID2, latest.SnapshotID)
			expectedPrice := decimal.NewFromFloat(3050.00)
			suite.True(latest.LastPrice.Equal(expectedPrice))
		})
	})
}

// TestMarketSnapshotQueryBySymbol tests querying snapshots by symbol
func (suite *MarketSnapshotBehaviorTestSuite) TestMarketSnapshotQueryBySymbol() {
	var (
		snapshotID1 = GenerateTestUUID()
		snapshotID2 = GenerateTestUUID()
	)

	suite.Given("snapshots for different symbols", func() {
		snapshot1 := suite.CreateTestMarketSnapshot(snapshotID1, func(ms *models.MarketSnapshot) {
			ms.Symbol = "SOL-USD"
		})
		err := suite.adapter.MarketSnapshotRepository().Create(suite.ctx, snapshot1)
		suite.Require().NoError(err)
		suite.trackCreatedMarketSnapshot(snapshotID1)

		snapshot2 := suite.CreateTestMarketSnapshot(snapshotID2, func(ms *models.MarketSnapshot) {
			ms.Symbol = "ADA-USD"
		})
		err = suite.adapter.MarketSnapshotRepository().Create(suite.ctx, snapshot2)
		suite.Require().NoError(err)
		suite.trackCreatedMarketSnapshot(snapshotID2)
	}).When("querying snapshots for SOL-USD", func() {
		snapshots, err := suite.adapter.MarketSnapshotRepository().GetBySymbol(suite.ctx, "SOL-USD")
		suite.Require().NoError(err)

		suite.Then("only SOL-USD snapshots should be returned", func() {
			suite.GreaterOrEqual(len(snapshots), 1)

			foundSOL := false
			foundADA := false
			for _, snapshot := range snapshots {
				if snapshot.SnapshotID == snapshotID1 {
					foundSOL = true
				}
				if snapshot.SnapshotID == snapshotID2 {
					foundADA = true
				}
			}
			suite.True(foundSOL)
			suite.False(foundADA, "ADA snapshot should not be in SOL query")
		})
	})
}
