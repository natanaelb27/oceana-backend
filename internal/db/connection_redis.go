package db

import (
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

func RedisConnect() {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal("Connecting to redis failed")
	}
	redisClient = redis.NewClient(opt)
}

func RedisGet() *redis.Client {
	return redisClient
}
