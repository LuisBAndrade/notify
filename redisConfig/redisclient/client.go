package redisclient

import (
	"context"
	"strconv"
	"sync"

	"github.com/LuisBAndrade/notify/redisConfig"
	"github.com/redis/go-redis/v9"
)

var (
	client *redis.Client
	once sync.Once
)

func GetRedisClient() *redis.Client {
	once.Do(func() {
		settings := redisConfig.GetRedisSettings()
		client = redis.NewClient(&redis.Options{
			Addr: settings.Host + ":" + strconv.Itoa(settings.Port),
			DB: settings.DB,
			Password: "",
		})

		if err := client.Ping(context.Background()).Err(); err != nil {
			panic("failed to connect to Redis: " + err.Error())
		}
	})

	return client
}