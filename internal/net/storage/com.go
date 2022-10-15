package storage

import (
	"context"
	"time"

	"center.bojiu.com/pkg/redislib"
)

//为了保持数据一致性, 采用延迟双删策略
var retrys = []int{1, 1, 2, 3, 5, 8, 13, 21, 34, 55}

func RedisDel(key string) error {
	redis := redislib.GetClient()
	ctx := context.Background()
	_, err := redis.Del(ctx, key).Result()
	return err
}

func RedisDelayDel(key string) {
	go func() {
		redis := redislib.GetClient()
		ctx := context.Background()
		for _, v := range retrys {
			time.Sleep(time.Duration(v) * time.Second)
			if _, err := redis.Del(ctx, key).Result(); err == nil {
				return
			}
		}
	}()
}
