package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	redis.Cmdable
}

func NewRedisClient(ctx context.Context, cfg Config) (Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       0,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

// Config 配置
type Config struct {
	Type     string
	Addr     string
	Password string
}
