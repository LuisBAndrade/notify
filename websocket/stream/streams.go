package redisstream 

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/LuisBAndrade/notify/redisConfig/redisclient"
	"github.com/LuisBAndrade/notify/utils"
	"github.com/redis/go-redis/v9"
)

var (
	GROUP_NAME = utils.GetGroupName
	MAX_STREAM_LENGTH int64 = 1000
	XREAD_TIMEOUT = 5 * time.Second
	XREAD_COUNT int64 = 1
	ERROR_SLEEP_SEC = 1 *time.Second
)

func PublishMessage(ctx context.Context, userID, message string) (string, error) {
	streamKey := utils.GetStreamKey(userID)
	payload := map[string]interface{}{
		"message": message,
		"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
	}
	client := redisclient.GetRedisClient()

	args := &redis.XAddArgs{
		Stream: streamKey,
		Values: payload,
		MaxLen: MAX_STREAM_LENGTH,
		Approx: true,
	}

	msgID, err := client.XAdd(ctx, args).Result()
	if err != nil {
		log.Printf("[Redis Publisher] Error adding message to %s: %v", streamKey, err)
		return "", err
	}

	log.Printf("[Redis Publisher] Added message to %s: %v (id: %s)", streamKey, payload, msgID)
	return msgID, nil  
}

func CreateConsumerGroup(ctx context.Context, userID string) error {
	streamKey := utils.GetStreamKey(userID)
	client := redisclient.GetRedisClient()

	err := client.XGroupCreateMkStream(ctx, streamKey, GROUP_NAME(), "0-0").Err()
	if err != nil {
		if !isBusyGroupErr(err) {
			log.Printf("Error creating consumer group on stream %s: %v", streamKey, err)
			return err
		}
		log.Printf("Consumer group %s already exists for stream %s", GROUP_NAME, streamKey)
	}
	return nil
}

func isBusyGroupErr(err error) bool {
	return strings.Contains(err.Error(), "BUSYGROUP")
}

func GetPendingNotifications(ctx context.Context, userID string) ([]redis.XMessage, error) {
	streamKey := utils.GetStreamKey(userID)
	client := redisclient.GetRedisClient()

	res, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group: GROUP_NAME(),
		Consumer: userID,
		Streams: []string{streamKey, "0"},
		Count: 100,
		Block: 0,
	}).Result()

	if err != nil && err != redis.Nil {
		log.Printf("Error getting pending notifications from %s: %v", streamKey, err)
		return nil, err
	}

	var notifications []redis.XMessage
	for _, strm := range res {
		notifications = append(notifications, strm.Messages...)
	}

	return notifications, nil 
}

func ListenForNotifications(ctx context.Context, userID string, sendFn func(msgID string, message string) error) {
	streamKey := utils.GetStreamKey(userID)
	client := redisclient.GetRedisClient()

	for {
		res, err := client.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group: GROUP_NAME(),
			Consumer: userID,
			Streams: []string{streamKey, ">"},
			Count: XREAD_COUNT,
			Block: XREAD_TIMEOUT,
		}).Result()

		if err != nil && err != redis.Nil {
			log.Printf("Error listening for notifications on %s: %v", streamKey, err)
			time.Sleep(ERROR_SLEEP_SEC)
			continue
		}

		if len(res) == 0 {
			time.Sleep(ERROR_SLEEP_SEC)
			continue
		}

		for _, strm := range res {
			for _, msg := range strm.Messages {
				if val, ok := msg.Values["message"].(string); ok {
					if sendErr := sendFn(msg.ID, val); sendErr != nil {
						log.Printf("Error sending message to websocket: %v", sendErr)
					}
				}
			}
		}
	}
}

func AcknowledgeNotifications(ctx context.Context, userID string, messageIDs []string) error {
	if len(messageIDs) == 0 {
		return nil
	}

	streamKey := utils.GetStreamKey(userID)
	client := redisclient.GetRedisClient()

	err := client.XAck(ctx, streamKey, GROUP_NAME(), messageIDs...).Err()
	if err != nil {
		log.Printf("Error acknowledging messages on %s: %v", streamKey, err)
		return fmt.Errorf("failed to ack messages: %w", err)
	}
	return nil 

}