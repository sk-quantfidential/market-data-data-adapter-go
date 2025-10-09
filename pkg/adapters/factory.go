package adapters

import (
	"context"
	"fmt"
	"strings"

	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/internal/cache"
	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/internal/config"
	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/internal/database"
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

type MarketDataAdapter struct {
	config *config.Config
	logger *logrus.Logger

	// Infrastructure
	postgresDB  *database.PostgresDB
	redisClient *cache.RedisClient

	// Repositories
	priceFeedRepo        interfaces.PriceFeedRepository
	candleRepo           interfaces.CandleRepository
	marketSnapshotRepo   interfaces.MarketSnapshotRepository
	symbolRepo           interfaces.SymbolRepository
	serviceDiscoveryRepo interfaces.ServiceDiscoveryRepository
	cacheRepo            interfaces.CacheRepository
}

func NewMarketDataAdapter(cfg *config.Config, logger *logrus.Logger) (DataAdapter, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}
	if logger == nil {
		return nil, fmt.Errorf("logger is required")
	}

	// Apply derivation if schema name not explicitly provided
	if cfg.SchemaName == "" {
		cfg.SchemaName = deriveSchemaName(cfg.ServiceName, cfg.ServiceInstanceName)
	}

	// Apply derivation if Redis namespace not explicitly provided
	if cfg.RedisNamespace == "" {
		cfg.RedisNamespace = deriveRedisNamespace(cfg.ServiceName, cfg.ServiceInstanceName)
	}

	logger.WithFields(logrus.Fields{
		"service_name":    cfg.ServiceName,
		"instance_name":   cfg.ServiceInstanceName,
		"schema_name":     cfg.SchemaName,
		"redis_namespace": cfg.RedisNamespace,
	}).Info("DataAdapter configuration resolved")

	adapter := &MarketDataAdapter{
		config: cfg,
		logger: logger,
	}

	// Initialize PostgreSQL
	if cfg.PostgresURL != "" {
		postgresDB, err := database.NewPostgresDB(cfg, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create PostgreSQL client: %w", err)
		}
		adapter.postgresDB = postgresDB

		// Initialize PostgreSQL repositories
		adapter.priceFeedRepo = NewPostgresPriceFeedRepository(postgresDB.DB, logger)
		adapter.candleRepo = NewPostgresCandleRepository(postgresDB.DB, logger)
		adapter.marketSnapshotRepo = NewPostgresMarketSnapshotRepository(postgresDB.DB, logger)
		adapter.symbolRepo = NewPostgresSymbolRepository(postgresDB.DB, logger)
	} else {
		logger.Warn("PostgreSQL URL not configured, repositories will not be available")
	}

	// Initialize Redis
	if cfg.RedisURL != "" {
		redisClient, err := cache.NewRedisClient(cfg, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to create Redis client: %w", err)
		}
		adapter.redisClient = redisClient

		// Initialize Redis repositories
		adapter.serviceDiscoveryRepo = NewRedisServiceDiscovery(redisClient.Client, cfg.ServiceDiscoveryNamespace, logger)
		adapter.cacheRepo = NewRedisCacheRepository(redisClient.Client, cfg.CacheNamespace, logger)
	} else {
		logger.Warn("Redis URL not configured, cache and service discovery will not be available")
	}

	return adapter, nil
}

func NewMarketDataAdapterFromEnv(logger *logrus.Logger) (DataAdapter, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return NewMarketDataAdapter(cfg, logger)
}

func (a *MarketDataAdapter) Connect(ctx context.Context) error {
	// Connect to PostgreSQL
	if a.postgresDB != nil {
		if err := a.postgresDB.Connect(ctx); err != nil {
			a.logger.WithError(err).Warn("Failed to connect to PostgreSQL (stub mode)")
		}
	}

	// Connect to Redis
	if a.redisClient != nil {
		if err := a.redisClient.Connect(ctx); err != nil {
			a.logger.WithError(err).Warn("Failed to connect to Redis (stub mode)")
		}
	}

	a.logger.Info("Market data adapter connected")
	return nil
}

func (a *MarketDataAdapter) Disconnect(ctx context.Context) error {
	var errors []error

	// Disconnect from PostgreSQL
	if a.postgresDB != nil {
		if err := a.postgresDB.Disconnect(ctx); err != nil {
			errors = append(errors, fmt.Errorf("PostgreSQL disconnect error: %w", err))
		}
	}

	// Disconnect from Redis
	if a.redisClient != nil {
		if err := a.redisClient.Disconnect(ctx); err != nil {
			errors = append(errors, fmt.Errorf("Redis disconnect error: %w", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("disconnect errors: %v", errors)
	}

	a.logger.Info("Market data adapter disconnected")
	return nil
}

func (a *MarketDataAdapter) HealthCheck(ctx context.Context) error {
	// Check PostgreSQL health
	if a.postgresDB != nil {
		if err := a.postgresDB.HealthCheck(ctx); err != nil {
			return fmt.Errorf("PostgreSQL health check failed: %w", err)
		}
	}

	// Check Redis health
	if a.redisClient != nil {
		if err := a.redisClient.HealthCheck(ctx); err != nil {
			return fmt.Errorf("Redis health check failed: %w", err)
		}
	}

	return nil
}

// Repository accessors
func (a *MarketDataAdapter) PriceFeedRepository() interfaces.PriceFeedRepository {
	return a.priceFeedRepo
}

func (a *MarketDataAdapter) CandleRepository() interfaces.CandleRepository {
	return a.candleRepo
}

func (a *MarketDataAdapter) MarketSnapshotRepository() interfaces.MarketSnapshotRepository {
	return a.marketSnapshotRepo
}

func (a *MarketDataAdapter) SymbolRepository() interfaces.SymbolRepository {
	return a.symbolRepo
}

func (a *MarketDataAdapter) ServiceDiscoveryRepository() interfaces.ServiceDiscoveryRepository {
	return a.serviceDiscoveryRepo
}

func (a *MarketDataAdapter) CacheRepository() interfaces.CacheRepository {
	return a.cacheRepo
}

// deriveSchemaName determines PostgreSQL schema based on service instance pattern
// Singleton: market-data-simulator == market-data-simulator → "market_data"
// Multi-instance: market-data-Coinmetrics → "market_data_coinmetrics"
func deriveSchemaName(serviceName, instanceName string) string {
	if serviceName == instanceName {
		// Singleton service pattern
		// Example: "market-data-simulator" -> "market_data"
		parts := strings.Split(serviceName, "-")
		if len(parts) >= 2 {
			return parts[0] + "_" + parts[1]
		}
		return serviceName
	}

	// Multi-instance service pattern
	// Example: "market-data-Coinmetrics" -> "market_data_coinmetrics"
	parts := strings.Split(instanceName, "-")
	if len(parts) >= 3 {
		// Has entity identifier: market-data-Coinmetrics -> market_data_coinmetrics
		return strings.ToLower(parts[0] + "_" + parts[1] + "_" + parts[2])
	} else if len(parts) >= 2 {
		// Only two parts: market-data -> market_data
		return strings.ToLower(parts[0] + "_" + parts[1])
	}
	return strings.ToLower(instanceName)
}

// deriveRedisNamespace determines Redis key prefix based on service instance pattern
// Singleton: market-data-simulator == market-data-simulator → "market_data"
// Multi-instance: market-data-Coinmetrics → "market_data:Coinmetrics"
func deriveRedisNamespace(serviceName, instanceName string) string {
	if serviceName == instanceName {
		// Singleton service pattern
		// Example: "market-data-simulator" -> "market_data"
		parts := strings.Split(serviceName, "-")
		if len(parts) >= 2 {
			return parts[0] + "_" + parts[1]
		}
		return serviceName
	}

	// Multi-instance service pattern
	// Example: "market-data-Coinmetrics" -> "market_data:Coinmetrics"
	parts := strings.Split(instanceName, "-")
	if len(parts) >= 2 {
		return parts[0] + "_" + parts[1] + ":" + strings.Join(parts[2:], "-")
	}
	return instanceName
}
