package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis(opts ...opt) (*redis.Client, error) {
	redisCfg := defaultRedisConfig()

	for _, opt := range opts {
		opt(&redisCfg)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", redisCfg.Host, redisCfg.Port),
		Password: redisCfg.Password,
		DB:       redisCfg.Db,
	})

	err := rdb.Ping(context.Background()).Err()

	return rdb, err
}
