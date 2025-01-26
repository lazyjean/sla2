package redis

import (
	"context"
	"fmt"

	"github.com/lazyjean/sla2/config"
	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func InitRedis(cfg *config.RedisConfig) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	// 测试连接
	ctx := context.Background()
	_, err := Client.Ping(ctx).Result()
	return err
}
