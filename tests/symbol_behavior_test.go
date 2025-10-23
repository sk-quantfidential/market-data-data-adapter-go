//go:build integration

package tests

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/pkg/models"
)

// SymbolBehaviorTestSuite tests the behavior of symbol repository operations
type SymbolBehaviorTestSuite struct {
	BehaviorTestSuite
}

// TestSymbolBehaviorSuite runs the symbol behavior test suite
func TestSymbolBehaviorSuite(t *testing.T) {
	suite.Run(t, new(SymbolBehaviorTestSuite))
}

// TestSymbolCRUDOperations tests basic symbol CRUD operations
func (suite *SymbolBehaviorTestSuite) TestSymbolCRUDOperations() {
	var symbolID = GenerateTestUUID()

	suite.Given("a new symbol to create", func() {
		// Symbol defined below
	}).When("creating the symbol", func() {
		symbol := suite.CreateTestSymbol(symbolID, func(s *models.Symbol) {
			s.Symbol = "ETH-USD"
			s.BaseCurrency = "ETH"
			s.QuoteCurrency = "USD"
			s.IsActive = true
		})

		err := suite.adapter.SymbolRepository().Create(suite.ctx, symbol)
		suite.Require().NoError(err)
		suite.trackCreatedSymbol(symbolID)
	}).Then("the symbol should be retrievable", func() {
		retrieved, err := suite.adapter.SymbolRepository().GetByID(suite.ctx, symbolID)
		suite.Require().NoError(err)
		suite.Equal(symbolID, retrieved.SymbolID)
		suite.Equal("ETH-USD", retrieved.Symbol)
		suite.True(retrieved.IsActive)
	}).And("the symbol can be updated", func() {
		err := suite.adapter.SymbolRepository().Update(suite.ctx, symbolID, func(s *models.Symbol) {
			s.IsActive = false
		})
		suite.Require().NoError(err)

		updated, err := suite.adapter.SymbolRepository().GetByID(suite.ctx, symbolID)
		suite.Require().NoError(err)
		suite.False(updated.IsActive)
	})
}

// TestSymbolQueryByStatus tests querying symbols by active status
func (suite *SymbolBehaviorTestSuite) TestSymbolQueryByStatus() {
	var (
		symbolID1 = GenerateTestUUID()
		symbolID2 = GenerateTestUUID()
	)

	suite.Given("multiple symbols with different status", func() {
		symbol1 := suite.CreateTestSymbol(symbolID1, func(s *models.Symbol) {
			s.Symbol = "BTC-USD"
			s.IsActive = true
		})
		err := suite.adapter.SymbolRepository().Create(suite.ctx, symbol1)
		suite.Require().NoError(err)
		suite.trackCreatedSymbol(symbolID1)

		symbol2 := suite.CreateTestSymbol(symbolID2, func(s *models.Symbol) {
			s.Symbol = "DOGE-USD"
			s.IsActive = false
		})
		err = suite.adapter.SymbolRepository().Create(suite.ctx, symbol2)
		suite.Require().NoError(err)
		suite.trackCreatedSymbol(symbolID2)
	}).When("querying active symbols", func() {
		activeSymbols, err := suite.adapter.SymbolRepository().GetActive(suite.ctx)
		suite.Require().NoError(err)

		suite.Then("only active symbols should be returned", func() {
			suite.GreaterOrEqual(len(activeSymbols), 1)

			foundBTC := false
			foundDOGE := false
			for _, symbol := range activeSymbols {
				if symbol.Symbol == "BTC-USD" {
					foundBTC = true
				}
				if symbol.Symbol == "DOGE-USD" {
					foundDOGE = true
				}
			}
			suite.True(foundBTC)
			suite.False(foundDOGE, "Inactive symbol should not be in active list")
		})
	})
}

// TestSymbolQueryByBaseCurrency tests querying symbols by base currency
func (suite *SymbolBehaviorTestSuite) TestSymbolQueryByBaseCurrency() {
	var (
		symbolID1 = GenerateTestUUID()
		symbolID2 = GenerateTestUUID()
	)

	suite.Given("symbols with different base currencies", func() {
		symbol1 := suite.CreateTestSymbol(symbolID1, func(s *models.Symbol) {
			s.Symbol = "BTC-USD"
			s.BaseCurrency = "BTC"
			s.QuoteCurrency = "USD"
		})
		err := suite.adapter.SymbolRepository().Create(suite.ctx, symbol1)
		suite.Require().NoError(err)
		suite.trackCreatedSymbol(symbolID1)

		symbol2 := suite.CreateTestSymbol(symbolID2, func(s *models.Symbol) {
			s.Symbol = "BTC-EUR"
			s.BaseCurrency = "BTC"
			s.QuoteCurrency = "EUR"
		})
		err = suite.adapter.SymbolRepository().Create(suite.ctx, symbol2)
		suite.Require().NoError(err)
		suite.trackCreatedSymbol(symbolID2)
	}).When("querying symbols by base currency BTC", func() {
		symbols, err := suite.adapter.SymbolRepository().GetByBaseCurrency(suite.ctx, "BTC")
		suite.Require().NoError(err)

		suite.Then("all BTC pairs should be returned", func() {
			suite.GreaterOrEqual(len(symbols), 2)

			quotes := make(map[string]bool)
			for _, symbol := range symbols {
				if symbol.BaseCurrency == "BTC" {
					quotes[symbol.QuoteCurrency] = true
				}
			}
			suite.True(quotes["USD"])
			suite.True(quotes["EUR"])
		})
	})
}
