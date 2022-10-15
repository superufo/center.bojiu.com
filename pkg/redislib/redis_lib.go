package redislib

import (
	"fmt"

	"center.bojiu.com/config"

	"context"

	"github.com/go-redis/redis/v8"
)

var (
	client *redis.Client
)

func newRedisClient() *redis.Client {
	cfg := config.GlobalCfg.Redis
	client = redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("初始化redis失败:", pong, err)
		return nil
	}
	fmt.Println("初始化redis成功")
	return client
}

func GetClient() (c *redis.Client) {
	if client == nil {
		return newRedisClient()
	}
	return client
}
