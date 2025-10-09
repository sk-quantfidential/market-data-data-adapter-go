package adapters

import (
	"testing"

	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Schema Derivation Tests
// =============================================================================

func TestDeriveSchemaName_Singleton_MarketDataSimulator(t *testing.T) {
	serviceName := "market-data-simulator"
	instanceName := "market-data-simulator"

	result := deriveSchemaName(serviceName, instanceName)

	assert.Equal(t, "market_data", result,
		"Singleton service should derive schema as 'market_data'")
}

func TestDeriveSchemaName_Singleton_MarketDataAdapter(t *testing.T) {
	serviceName := "market-data-adapter"
	instanceName := "market-data-adapter"

	result := deriveSchemaName(serviceName, instanceName)

	assert.Equal(t, "market_data", result,
		"Adapter singleton should also derive as 'market_data'")
}

func TestDeriveSchemaName_MultiInstance_Coinmetrics(t *testing.T) {
	serviceName := "market-data-simulator"
	instanceName := "market-data-Coinmetrics"

	result := deriveSchemaName(serviceName, instanceName)

	assert.Equal(t, "market_data_coinmetrics", result,
		"Multi-instance Coinmetrics should derive as 'market_data_coinmetrics' (lowercase)")
}

func TestDeriveSchemaName_MultiInstance_Bloomberg(t *testing.T) {
	serviceName := "market-data-simulator"
	instanceName := "market-data-Bloomberg"

	result := deriveSchemaName(serviceName, instanceName)

	assert.Equal(t, "market_data_bloomberg", result,
		"Multi-instance Bloomberg should derive as 'market_data_bloomberg' (lowercase)")
}

func TestDeriveSchemaName_MultiInstance_Glassnode(t *testing.T) {
	serviceName := "market-data-simulator"
	instanceName := "market-data-Glassnode"

	result := deriveSchemaName(serviceName, instanceName)

	assert.Equal(t, "market_data_glassnode", result,
		"Multi-instance Glassnode should derive as 'market_data_glassnode'")
}

func TestDeriveSchemaName_EdgeCase_SinglePart(t *testing.T) {
	serviceName := "marketdata"
	instanceName := "marketdata"

	result := deriveSchemaName(serviceName, instanceName)

	assert.Equal(t, "marketdata", result,
		"Single-part names should be returned as-is")
}

func TestDeriveSchemaName_EdgeCase_Empty(t *testing.T) {
	serviceName := ""
	instanceName := ""

	result := deriveSchemaName(serviceName, instanceName)

	assert.Equal(t, "", result,
		"Empty service names should return empty string")
}

func TestDeriveSchemaName_EdgeCase_MultipartInstance(t *testing.T) {
	serviceName := "market-data-simulator"
	instanceName := "market-data-Coinmetrics-Pro"

	result := deriveSchemaName(serviceName, instanceName)

	// Should take first two parts only
	assert.Equal(t, "market_data_coinmetrics", result,
		"Multi-part instance names should use first two parts only")
}

// =============================================================================
// Redis Namespace Derivation Tests
// =============================================================================

func TestDeriveRedisNamespace_Singleton_MarketDataSimulator(t *testing.T) {
	serviceName := "market-data-simulator"
	instanceName := "market-data-simulator"

	result := deriveRedisNamespace(serviceName, instanceName)

	assert.Equal(t, "market_data", result,
		"Singleton service should derive namespace as 'market_data'")
}

func TestDeriveRedisNamespace_Singleton_MarketDataAdapter(t *testing.T) {
	serviceName := "market-data-adapter"
	instanceName := "market-data-adapter"

	result := deriveRedisNamespace(serviceName, instanceName)

	assert.Equal(t, "market_data", result,
		"Adapter singleton should also derive as 'market_data'")
}

func TestDeriveRedisNamespace_MultiInstance_Coinmetrics(t *testing.T) {
	serviceName := "market-data-simulator"
	instanceName := "market-data-Coinmetrics"

	result := deriveRedisNamespace(serviceName, instanceName)

	assert.Equal(t, "market_data:Coinmetrics", result,
		"Multi-instance Coinmetrics should derive as 'market_data:Coinmetrics' (preserve case)")
}

func TestDeriveRedisNamespace_MultiInstance_Bloomberg(t *testing.T) {
	serviceName := "market-data-simulator"
	instanceName := "market-data-Bloomberg"

	result := deriveRedisNamespace(serviceName, instanceName)

	assert.Equal(t, "market_data:Bloomberg", result,
		"Multi-instance Bloomberg should derive as 'market_data:Bloomberg'")
}

func TestDeriveRedisNamespace_MultiInstance_Glassnode(t *testing.T) {
	serviceName := "market-data-simulator"
	instanceName := "market-data-Glassnode"

	result := deriveRedisNamespace(serviceName, instanceName)

	assert.Equal(t, "market_data:Glassnode", result,
		"Multi-instance Glassnode should derive as 'market_data:Glassnode'")
}

func TestDeriveRedisNamespace_EdgeCase_SinglePart(t *testing.T) {
	serviceName := "marketdata"
	instanceName := "marketdata"

	result := deriveRedisNamespace(serviceName, instanceName)

	assert.Equal(t, "marketdata", result,
		"Single-part names should be returned as-is")
}

func TestDeriveRedisNamespace_EdgeCase_Empty(t *testing.T) {
	serviceName := ""
	instanceName := ""

	result := deriveRedisNamespace(serviceName, instanceName)

	assert.Equal(t, "", result,
		"Empty service names should return empty string")
}

func TestDeriveRedisNamespace_EdgeCase_MultipartInstance(t *testing.T) {
	serviceName := "market-data-simulator"
	instanceName := "market-data-Coinmetrics-Pro"

	result := deriveRedisNamespace(serviceName, instanceName)

	// Should preserve all parts after the second hyphen in the entity section
	assert.Equal(t, "market_data:Coinmetrics-Pro", result,
		"Multi-part entity names should preserve full entity identifier after colon")
}

// =============================================================================
// Factory Integration Tests
// =============================================================================

func TestNewMarketDataAdapter_DerivesSchemaWhenNotProvided(t *testing.T) {
	cfg := &config.Config{
		ServiceName:         "market-data-simulator",
		ServiceInstanceName: "market-data-simulator",
		// SchemaName intentionally not set
		PostgresURL: "", // No actual connection needed for this test
		RedisURL:    "",
	}

	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel) // Suppress logs during test

	adapter, err := NewMarketDataAdapter(cfg, logger)

	require.NoError(t, err, "Should create adapter without error")
	require.NotNil(t, adapter, "Adapter should not be nil")

	// Verify schema was derived
	assert.Equal(t, "market_data", cfg.SchemaName,
		"SchemaName should be automatically derived to 'market_data'")
}

func TestNewMarketDataAdapter_DerivesNamespaceWhenNotProvided(t *testing.T) {
	cfg := &config.Config{
		ServiceName:         "market-data-simulator",
		ServiceInstanceName: "market-data-Coinmetrics",
		// RedisNamespace intentionally not set
		PostgresURL: "",
		RedisURL:    "",
	}

	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	adapter, err := NewMarketDataAdapter(cfg, logger)

	require.NoError(t, err)
	require.NotNil(t, adapter)

	// Verify namespace was derived
	assert.Equal(t, "market_data:Coinmetrics", cfg.RedisNamespace,
		"RedisNamespace should be automatically derived to 'market_data:Coinmetrics'")
}

func TestNewMarketDataAdapter_HonorsExplicitSchemaAndNamespace(t *testing.T) {
	explicitSchema := "custom_schema"
	explicitNamespace := "custom:namespace"

	cfg := &config.Config{
		ServiceName:         "market-data-simulator",
		ServiceInstanceName: "market-data-simulator",
		SchemaName:          explicitSchema,
		RedisNamespace:      explicitNamespace,
		PostgresURL:         "",
		RedisURL:            "",
	}

	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	adapter, err := NewMarketDataAdapter(cfg, logger)

	require.NoError(t, err)
	require.NotNil(t, adapter)

	// Verify explicit values are honored (not overridden by derivation)
	assert.Equal(t, explicitSchema, cfg.SchemaName,
		"Explicit SchemaName should not be overridden")
	assert.Equal(t, explicitNamespace, cfg.RedisNamespace,
		"Explicit RedisNamespace should not be overridden")
}
