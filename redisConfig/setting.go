package redisConfig

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

type RedisSettings struct {
	Host string
	Port int
	DB int
}

var (
	settings *RedisSettings
	once sync.Once
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func GetRedisSettings() *RedisSettings {
	once.Do(func() {
		portStr := getEnv("REDIS_PORT", "6379")
		port, err := strconv.Atoi(portStr)
		if err != nil {
			port = 6379
		}
		dbStr := getEnv("REDIS_DB", "0")
		db, err := strconv.Atoi(dbStr)
		if err != nil {
			db = 0
		}
		settings = &RedisSettings{
			Host: getEnv("REDIS_HOST", "localhost"),
			Port: port,
			DB: db,
		}
	})
	return settings
}

func (r *RedisSettings) GetRedisURL() string {
	return fmt.Sprintf("redis://%s:%d/%d", r.Host, r.Port, r.DB)
}