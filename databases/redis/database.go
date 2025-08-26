package redis

import (
	"context"
	"log"
	"simple-crud-notes/configs"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func RedisInit(appConfig configs.AppConfig) {
	opts, err := redis.ParseURL(appConfig.REDIS_URL)
	if err != nil {
		log.Fatal("Failed to parse REDIS_URL: ", err)
	}

	RedisClient = redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis: ", err)
	}

	log.Println("Successfully connected to the redis")
}

func RedisClose() {
	if RedisClient != nil {
		RedisClient.Close()
	}
}
