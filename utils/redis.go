package utils

import (
	"context"
	"encoding/json"
	"simple-crud-notes/databases/redis"
	"time"
)

func SetCache(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return redis.RedisClient.Set(ctx, key, jsonValue, expiration).Err()
}

func GetCache(key string, dest interface{}) error {
	ctx := context.Background()
	val, _ := redis.RedisClient.Get(ctx, key).Result()

	if val == "" {
		return nil
	}

	return json.Unmarshal([]byte(val), dest)
}

func DeleteCache(key string) error {
	ctx := context.Background()
	return redis.RedisClient.Del(ctx, key).Err()
}

func DeleteTokensByPattern(pattern string) error {
	ctx := context.Background()
	iter := redis.RedisClient.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()

		err := redis.RedisClient.Del(ctx, key).Err()
		if err != nil {
			println("Failed to delete key %s: %v", key, err)
			return err
		}
	}

	if err := iter.Err(); err != nil {
		println("Error scanning keys: %v", err)
		return err
	}

	println("Successfully deleted keys matching pattern: %s", pattern)
	return nil
}
