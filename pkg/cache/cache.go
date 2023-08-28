package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrCacheMiss = errors.New("cache miss")

type Cache interface {
	Set(ctx context.Context, key string, data interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, data interface{}) error
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{
		client: client,
	}
}

func (rc *RedisCache) Set(ctx context.Context, key string, data interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data for cache: %w", err)
	}

	_, err = rc.client.Set(ctx, key, jsonData, expiration).Result()
	if err != nil {
		return fmt.Errorf("failed to set data in cache: %w", err)
	}

	return nil
}

func (rc *RedisCache) Get(ctx context.Context, key string, data interface{}) error {
	val, err := rc.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return fmt.Errorf("failed to get data from cache: %w", err)
	}

	err = json.Unmarshal([]byte(val), data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal data from cache: %w", err)
	}

	return nil
}

func (rc *RedisCache) Invalidate(ctx context.Context) error {
	return rc.client.FlushAll(ctx).Err()
}
