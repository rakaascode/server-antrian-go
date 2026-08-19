package redis

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	Ctx    = context.Background()
)

// InitRedis menginisialisasi koneksi Redis client
func InitRedis() (*redis.Client, error) {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}
	password := os.Getenv("REDIS_PASSWORD")

	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	})

	if err := Client.Ping(Ctx).Err(); err != nil {
		log.Printf("⚠️ Redis Ping gagal (%s:%s): %v", host, port, err)
		return nil, err
	}

	log.Printf("✅ Redis connected (%s:%s)", host, port)
	return Client, nil
}
