package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hibiken/asynq"
)

func postToSlack(token, channel, text string) error {
	payload := map[string]string{"channel": channel, "text": text}
	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"https://slack.com/api/chat.postMessage",
		bytes.NewBuffer(jsonPayload),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer"+token)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("slack api return non-200 status code: %s", response.Status)
	}

	return nil
}

func main() {
	// 1. Create a new database connection pool.
	dbPool, dbPoolErr := newDBConnectionPool()
	if dbPoolErr != nil {
		log.Printf("FATAL: Failed to connect to database pool: %v", dbPoolErr)
		return
	}
	defer dbPool.Close()

	// 2. Create a new Redis connection client.
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Printf("FATAL: REDIS_URL environment variable is not set!")
		return
	}

	redisConnection, redisErr := asynq.ParseRedisURI(redisURL)
	if redisErr != nil {
		log.Printf("FATAL: Failed to parse Redis URI: %v", redisErr)
		return
	}

	// 3. Create a new TaskProcessor.
	processor := NewTaskProcessor(redisConnection.(asynq.RedisClientOpt), dbPool)

	// 4. Start the processor.
	if err := processor.Start(); err != nil {
		log.Printf("FATAL: Could not run worker server: %v", err)
		return
	}
}
