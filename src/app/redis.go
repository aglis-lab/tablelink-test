package app

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedis() (*redis.Client, error) {
	// Redis Database
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	status := rdb.Ping(context.Background())
	if err := status.Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}

// TODO: Cache Token and access 1 Hour
func RedisCache[T any](ctx context.Context, redisClient *redis.Client, redisKey string, callback func(context.Context) (T, error)) (T, error) {
	val, err := redisClient.Get(ctx, redisKey).Bytes()
	if err == nil {
		var data T
		if err = json.Unmarshal(val, &data); err == nil {
			return data, nil
		}
	}

	data, err := callback(ctx)
	if err != nil {
		return data, err
	}

	if valBytes, err := json.Marshal(data); err == nil {
		redisClient.Set(ctx, redisKey, string(valBytes), time.Hour*24*7)
	}

	return data, nil
}
