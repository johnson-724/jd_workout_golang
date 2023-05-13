package redis

import (
	client "github.com/redis/go-redis/v9"
	"os"
)

var Connection *client.Client

func InitRedis() {
	connection := client.NewClient(&client.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})	

	Connection = connection
}