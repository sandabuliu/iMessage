package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var (
	rdb    *redis.Client
	prefix = "iMessage:"
)

func InitRedisClient(addr, password string, db int) {
	if rdb != nil {
		rdb.Close()
	}
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
}

func Set(ctx context.Context, key string, value interface{}) error {
	return rdb.Set(ctx, prefix+key, value, 0).Err()
}

func Get(ctx context.Context, key string) (string, error) {
	return rdb.Get(ctx, prefix+key).Result()
}

func LPush(ctx context.Context, key string, values ...interface{}) error {
	return rdb.LPush(ctx, prefix+key, values...).Err()
}

func RPush(ctx context.Context, key string, values ...interface{}) error {
	return rdb.RPush(ctx, prefix+key, values...).Err()
}

func LPop(ctx context.Context, key string) (string, error) {
	return rdb.LPop(ctx, prefix+key).Result()
}

func RPop(ctx context.Context, key string) (string, error) {
	return rdb.RPop(ctx, prefix+key).Result()
}

func BRPop(ctx context.Context, key string, timeout time.Duration) ([]string, error) {
	result := rdb.BRPop(ctx, timeout, prefix+key)
	return result.Result()
}
