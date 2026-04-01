package config

import (
	"context"
	"crypto/tls"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const (
	ClientExpiry       = 2 * time.Hour
	HeadlessSessionTTL = 3 * time.Hour
)

// RedisClient wraps go-redis with health tracking and safe operations.
type RedisClient struct {
	client    *redis.Client
	connected bool
}

// NewRedisClient creates and verifies a Redis connection.
func NewRedisClient(cfg *Config) (*RedisClient, error) {
	if cfg.RedisURL == "" {
		return nil, nil
	}

	opts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, err
	}

	// Enable TLS for Upstash
	if strings.Contains(cfg.RedisURL, "upstash.io") {
		opts.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	opts.DialTimeout = 15 * time.Second
	opts.MaxRetries = 20
	opts.MinRetryBackoff = 1 * time.Second
	opts.MaxRetryBackoff = 10 * time.Second

	client := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Info().Msg("Redis connected")
	return &RedisClient{client: client, connected: true}, nil
}

// Client returns the underlying redis client. May be nil.
func (r *RedisClient) Client() *redis.Client {
	if r == nil {
		return nil
	}
	return r.client
}

// IsConnected returns whether Redis is available.
func (r *RedisClient) IsConnected() bool {
	return r != nil && r.connected
}

// Close gracefully closes the Redis connection.
func (r *RedisClient) Close() error {
	if r == nil || r.client == nil {
		return nil
	}
	r.connected = false
	return r.client.Close()
}

// SafeSet performs a SET with fallback on error.
func (r *RedisClient) SafeSet(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if !r.IsConnected() {
		return nil
	}
	return r.client.Set(ctx, key, value, expiration).Err()
}

// SafeGet performs a GET with empty string fallback on error.
func (r *RedisClient) SafeGet(ctx context.Context, key string) (string, error) {
	if !r.IsConnected() {
		return "", nil
	}
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

// SafeDel performs a DEL with fallback on error.
func (r *RedisClient) SafeDel(ctx context.Context, keys ...string) error {
	if !r.IsConnected() {
		return nil
	}
	return r.client.Del(ctx, keys...).Err()
}

// SafeSAdd performs an SADD with fallback on error.
func (r *RedisClient) SafeSAdd(ctx context.Context, key string, members ...interface{}) error {
	if !r.IsConnected() {
		return nil
	}
	return r.client.SAdd(ctx, key, members...).Err()
}

// SafeSRem performs an SREM with fallback on error.
func (r *RedisClient) SafeSRem(ctx context.Context, key string, members ...interface{}) error {
	if !r.IsConnected() {
		return nil
	}
	return r.client.SRem(ctx, key, members...).Err()
}

// SafeSMembers performs an SMEMBERS with empty slice fallback.
func (r *RedisClient) SafeSMembers(ctx context.Context, key string) ([]string, error) {
	if !r.IsConnected() {
		return nil, nil
	}
	return r.client.SMembers(ctx, key).Result()
}

// SafeSCard performs an SCARD with 0 fallback.
func (r *RedisClient) SafeSCard(ctx context.Context, key string) (int64, error) {
	if !r.IsConnected() {
		return 0, nil
	}
	return r.client.SCard(ctx, key).Result()
}

// SafeExpire refreshes TTL on a key.
func (r *RedisClient) SafeExpire(ctx context.Context, key string, expiration time.Duration) error {
	if !r.IsConnected() {
		return nil
	}
	return r.client.Expire(ctx, key, expiration).Err()
}

// CheckHealth returns the Redis health status.
func (r *RedisClient) CheckHealth(ctx context.Context) (string, error) {
	if !r.IsConnected() {
		return "disconnected", nil
	}
	if err := r.client.Ping(ctx).Err(); err != nil {
		return "unhealthy", err
	}
	return "healthy", nil
}
