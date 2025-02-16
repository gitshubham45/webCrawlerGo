package queue

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var redisClient *redis.Client

func init() {
	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", 
		Password: "",               
		DB:       0,             
	})
}

func ClearQueue(queueName string) {
	err := redisClient.Del(ctx, queueName).Err()
	if err != nil {
		log.Printf("Failed to clear queue %s: %v", queueName, err)
		return
	}
	log.Printf("Cleared queue %s", queueName)
}

func AddTask(queueName string, url string) error {
	err := redisClient.RPush(ctx, queueName, url).Err()
	if err != nil {
		log.Printf("Failed to add task %s to queue %s: %v", url, queueName, err)
		return err
	}
	log.Printf("Added task %s to queue %s", url, queueName)
	return nil
}

func GetTask(queueName string) (string, error) {
	result, err := redisClient.RPop(ctx, queueName).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		log.Printf("Error popping from queue %s: %v", queueName, err)
		return "", err
	}
	log.Printf("Retrieved task %s from queue %s", result, queueName)
	return result, nil
}