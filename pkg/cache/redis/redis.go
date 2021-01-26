package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/tinysrc/z9go/pkg/conf"
	"github.com/tinysrc/z9go/pkg/log"
	"go.uber.org/zap"
)

// Pool redis pool
var Pool *redis.Pool

func initConfig() {
	conf.Global.SetDefault("redis.addr", "localhost:6379")
	conf.Global.SetDefault("redis.password", "")
	conf.Global.SetDefault("redis.dbIndex", 0)
	conf.Global.SetDefault("redis.maxIdle", 10)
	conf.Global.SetDefault("redis.maxActive", 100)
	conf.Global.SetDefault("redis.idleTimeout", 1800)
}

func init() {
	initConfig()
	addr := conf.Global.GetString("redis.addr")
	password := conf.Global.GetString("redis.password")
	dbIndex := conf.Global.GetInt("redis.dbIndex")
	maxIdle := conf.Global.GetInt("redis.maxIdle")
	maxActive := conf.Global.GetInt("redis.maxActive")
	idleTimeout := conf.Global.GetDuration("redis.IdleTimeout")
	Pool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
		Dial: func() (conn redis.Conn, err error) {
			conn, err = redis.Dial("tcp", addr)
			if err != nil {
				log.Error("redis dial failed", zap.String("addr", addr))
				return
			}
			if password != "" {
				if _, err = conn.Do("AUTH", password); err != nil {
					log.Error("redis auth failed", zap.String("addr", addr))
					conn.Close()
					return
				}
			}
			if _, err = conn.Do("SELECT", dbIndex); err != nil {
				log.Error("redis select failed", zap.String("addr", addr), zap.Int("dbIndex", dbIndex))
				conn.Close()
				return
			}
			return
		},
	}
	log.Info("redis init success")
}
