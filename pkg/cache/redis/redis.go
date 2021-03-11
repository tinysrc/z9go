package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/tinysrc/z9go/pkg/conf"
	"github.com/tinysrc/z9go/pkg/log"
	"go.uber.org/zap"
)

// Client instance
var Client *redis.Client

func initConfig() {
	conf.Global.SetDefault("redis.addr", "localhost:6379")
	conf.Global.SetDefault("redis.password", "")
	conf.Global.SetDefault("redis.db", 0)
}

func init() {
	initConfig()
	Client = redis.NewClient(&redis.Options{
		Addr:     conf.Global.GetString("redis.addr"),
		Password: conf.Global.GetString("redis.password"),
		DB:       conf.Global.GetInt("redis.db"),
	})
	ctx := context.Background()
	_, err := Client.Ping(ctx).Result()
	if err != nil {
		log.Error("redis init failed", zap.Error(err))
	} else {
		log.Info("redis init success")
	}
}
