package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/hibiken/asynq"
)

// TestMessagePayload must match the structure of the payload sent by the Nuxt API.
type TestMessagePayload struct {
	Message string `json:"message"`
	SentAt  string `json:"sentAt"`
}

type ExecuteWorkloadPayload struct {
	WorkflowID int `json:"workflowId"`
}

func handleExecuteWorkflowTask(_ context.Context, t *asynq.Task) error {
	var p ExecuteWorkloadPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Printf("ERROR: Failed to unmarshal payload for task %s: %v", t.Type(), err)
		return err
	}

	log.Printf("Received a job to execute workflow with ID: %d", p.WorkflowID)
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
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL environment variable is not set!")
	}

	redisConnection, err := asynq.ParseRedisURI(redisURL)
	if err != nil {
		log.Fatalf("Failed to parse Redis URI: %v", err)
	}

	srv := asynq.NewServer(redisConnection, asynq.Config{
		// Queue must match the name from Nuxt
		Queues: map[string]int{
			"workflows": 1,
		},
	})

	mux := asynq.NewServeMux()
	mux.HandleFunc("test-message", handleTestMessageTask)
	mux.HandleFunc("execute-workflow-v1", handleExecuteWorkflowTask)

	log.Println("Worker service started. Listening for jobs...")
	if srvErr := srv.Run(mux); srvErr != nil {
		log.Fatalf("Failed to start server: %v", srvErr)
	}
}
