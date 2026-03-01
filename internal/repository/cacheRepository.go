package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheRepository interface {
	Get(key string) (string, error)
	Set(key string, value string, exp time.Duration) error
	Delete(key string) error
}

type cacheRepository struct {
	redis *redis.Client
}

func (c *cacheRepository) Delete(key string) error {
	err := c.redis.Del(context.Background(), key).Err()
	return err
}

func (c *cacheRepository) Get(key string) (string, error) {
	str, err := c.redis.Get(context.Background(), key).Result()
	return str, err
}

func (c *cacheRepository) Set(key string, value string, exp time.Duration) error {
	err := c.redis.Set(context.Background(), key, value, exp).Err()
	return err
}

func NewCacheRepository(redis *redis.Client) CacheRepository {
	return &cacheRepository{
		redis: redis,
	}
}
