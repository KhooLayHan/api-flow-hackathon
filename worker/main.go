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

	req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewBuffer(jsonPayload))
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

// Handler function for the test message payload.
func handleTestMessageTask(_ context.Context, task *asynq.Task) error {
	var payload TestMessagePayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		log.Printf("ERROR: Failed to unmarshal payload for task %s: %v", task.Type(), err)
		return err
	}

	log.Printf("Received a job test message: %s, Sent at: %s", payload.Message, payload.SentAt)
	return nil
}

func main() {
	// 1. Create a new database connection pool.
	dbPool, err := newDBConnectionPool()
	if err != nil {
		log.Fatalf("FATAL: Failed to connect to database pool: %v", err)
	}
	defer dbPool.Close()

	// 2. Create a new Redis connection client.
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("FATAL: REDIS_URL environment variable is not set!")
	}

	redisConnection, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		log.Fatalf("FATAL: Failed to parse Redis URI: %v", err)
	}

	// 3. Create a new TaskProcessor.
	processor := NewTaskProcessor(redisConnection.(asynq.RedisClientOpt), dbPool)

	// 4. Start the processor.
	if err := processor.Start(); err != nil {
		log.Fatalf("FATAL: Could not run worker server: %v", err)
	}

	// srv := asynq.NewServer(redisConnection, asynq.Config{
	// 	// Queue must match the name from Nuxt
	// 	Queues: map[string]int{
	// 		"workflows": 1,
	// 	},
	// })

	// mux := asynq.NewServeMux()
	// mux.HandleFunc("test-message", handleTestMessageTask)
	// mux.HandleFunc("execute-workflow-v1", handleExecuteWorkflowTask)

	// log.Println("Worker service started. Listening for jobs...")
	// if srvErr := srv.Run(mux); srvErr != nil {
	// 	log.Fatalf("Failed to start server: %v", srvErr)
	// }
}
