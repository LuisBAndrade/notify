package utils

import (
	"fmt"
	"os"
)


func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func GetStreamKey(userID string) string {
	env := getEnv("APP_ENV", "local")
	app := "notiq"
	streamPrefix := "notifications"
	return fmt.Sprintf("%s:%s:%s:%s", env, app, streamPrefix, userID)
}

func GetGroupName() string {
	env := getEnv("APP_ENV", "local")
	app := "notiq"
	groupPrefix := "notification_group"
	return fmt.Sprintf("%s:%s:%s", env,app, groupPrefix)
}