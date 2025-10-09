package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// CacheBehaviorTestSuite tests the behavior of cache repository operations
type CacheBehaviorTestSuite struct {
	BehaviorTestSuite
}

// TestCacheBehaviorSuite runs the cache behavior test suite
func TestCacheBehaviorSuite(t *testing.T) {
	suite.Run(t, new(CacheBehaviorTestSuite))
}

// TestCacheStringOperations tests string value caching
func (suite *CacheBehaviorTestSuite) TestCacheStringOperations() {
	var (
		key   = "test:string:" + GenerateTestID("key")
		value = "Hello, Market Data Cache!"
		ttl   = 5 * time.Minute
	)

	suite.Given("a string value to cache", func() {
		// Value defined above
	}).When("storing the string in cache", func() {
		err := suite.adapter.CacheRepository().Set(suite.ctx, key, value, ttl)
		suite.Require().NoError(err)
	}).Then("the string should be retrievable", func() {
		retrieved, err := suite.adapter.CacheRepository().Get(suite.ctx, key)
		suite.Require().NoError(err)
		suite.Equal(value, retrieved)
	}).And("the key should exist", func() {
		exists, err := suite.adapter.CacheRepository().Exists(suite.ctx, key)
		suite.Require().NoError(err)
		suite.True(exists)
	})
}

// TestCacheExpiration tests cache TTL behavior
func (suite *CacheBehaviorTestSuite) TestCacheExpiration() {
	var (
		key   = "test:expire:" + GenerateTestID("key")
		value = "temporary-value"
		ttl   = 2 * time.Second
	)

	suite.Given("a cached value with short TTL", func() {
		err := suite.adapter.CacheRepository().Set(suite.ctx, key, value, ttl)
		suite.Require().NoError(err)
	}).When("waiting for TTL to expire", func() {
		time.Sleep(3 * time.Second)
	}).Then("the key should no longer exist", func() {
		exists, err := suite.adapter.CacheRepository().Exists(suite.ctx, key)
		suite.Require().NoError(err)
		suite.False(exists)
	})
}

// TestCacheDelete tests cache deletion
func (suite *CacheBehaviorTestSuite) TestCacheDelete() {
	var (
		key   = "test:delete:" + GenerateTestID("key")
		value = "to-be-deleted"
		ttl   = 10 * time.Minute
	)

	suite.Given("a cached value", func() {
		err := suite.adapter.CacheRepository().Set(suite.ctx, key, value, ttl)
		suite.Require().NoError(err)
	}).When("deleting the key", func() {
		err := suite.adapter.CacheRepository().Delete(suite.ctx, key)
		suite.Require().NoError(err)
	}).Then("the key should not exist", func() {
		exists, err := suite.adapter.CacheRepository().Exists(suite.ctx, key)
		suite.Require().NoError(err)
		suite.False(exists)
	})
}
