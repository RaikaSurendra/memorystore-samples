package cache

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var Rdb *redis.Client

func InitRedis() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	if redisHost == "" {
		redisHost = "localhost"
	}
	if redisPort == "" {
		redisPort = "6379"
	}

	Rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})

	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("Could not connect to Redis: %v", err)
	} else {
		log.Println("Connected to Redis")
	}
}

func Get(key string) (string, error) {
	val, err := Rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Key does not exist
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func Set(key string, value string) error {
	// Set key with a TTL of 1 hour (simulating cache behavior)
	err := Rdb.Set(ctx, key, value, 1*time.Hour).Err()
	return err
}

func Delete(key string) error {
	return Rdb.Del(ctx, key).Err()
}
