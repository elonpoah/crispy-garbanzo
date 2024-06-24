package initialize

import (
	"crispy-garbanzo/global"

	"github.com/go-redis/redis/v8"
)

func Redis() *redis.Client {
	cfg := global.FPG_CONFIG.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Url,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	return rdb
}
