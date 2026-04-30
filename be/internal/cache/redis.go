package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"server/config"
)

var client *redis.Client

func InitRedis() error {
	cfg := config.GlobalConfig.Redis
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if _, err := client.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}
	return nil
}

func Client() *redis.Client {
	return client
}
