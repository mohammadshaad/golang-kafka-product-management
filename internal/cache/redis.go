package cache

import (
	"context"
	"time"


	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

var ctx = context.Background()

func InitRedis(addr, username, password string) {

	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       0,
	})
}

func GetCache(key string) (string, error) {
	return rdb.Get(ctx, key).Result()
}

func SetCache(key string, value string, ttl time.Duration) error {
	return rdb.Set(ctx, key, value, ttl).Err()
}

func InvalidateCache(key string) error {
	return rdb.Del(ctx, key).Err()
}
