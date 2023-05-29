package redisrepo

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitializeRedis() *redis.Client {
	conn := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_CONNECTION_STRING"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	pong, err := conn.Ping(context.Background()).Result()
	if err != nil {
		log.Println("redis connection failed ", err)
		return nil
	}

	log.Println("redis client connected ping : ", pong)
	redisClient = conn
	return redisClient
}
