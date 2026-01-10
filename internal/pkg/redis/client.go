package redispkg

import (
	"fmt"

	"github.com/Soujuruya/01_SPEC/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewClient(cfg *config.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})
	return rdb
}
