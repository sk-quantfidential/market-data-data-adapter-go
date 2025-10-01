package cache

import (
	"context"
	"fmt"

	"github.com/quantfidential/trading-ecosystem/market-data-adapter-go/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisClient struct {
	Client *redis.Client
	config *config.Config
	logger *logrus.Logger
}

func NewRedisClient(cfg *config.Config, logger *logrus.Logger) (*RedisClient, error) {
	if cfg.RedisURL == "" {
		return nil, fmt.Errorf("Redis URL is required")
	}

	opts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	// Configure pool settings
	opts.PoolSize = cfg.RedisPoolSize
	opts.MinIdleConns = cfg.RedisMinIdleConns
	opts.MaxRetries = cfg.RedisMaxRetries
	opts.DialTimeout = cfg.RedisDialTimeout
	opts.ReadTimeout = cfg.RedisReadTimeout
	opts.WriteTimeout = cfg.RedisWriteTimeout

	return &RedisClient{
		Client: redis.NewClient(opts),
		config: cfg,
		logger: logger,
	}, nil
}

func (r *RedisClient) Connect(ctx context.Context) error {
	// Test connection
	if err := r.Client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	r.logger.Info("Redis connected successfully")
	return nil
}

func (r *RedisClient) Disconnect(ctx context.Context) error {
	if r.Client != nil {
		if err := r.Client.Close(); err != nil {
			return fmt.Errorf("failed to close Redis connection: %w", err)
		}
		r.logger.Info("Redis disconnected")
	}
	return nil
}

func (r *RedisClient) HealthCheck(ctx context.Context) error {
	if r.Client == nil {
		return fmt.Errorf("Redis not connected")
	}
	return r.Client.Ping(ctx).Err()
}
