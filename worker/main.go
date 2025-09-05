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

// TestMessagePayload must match the structure of the payload sent by the Nuxt API.
type TestMessagePayload struct {
	Message string `json:"message"`
	SentAt  string `json:"sentAt"`
}

type ExecuteWorkloadPayload struct {
	WorkflowID int    `json:"workflowId"`
	UserID     string `json:"userId"`
}

func handleExecuteWorkflowTask(_ context.Context, t *asynq.Task) error {
	var p ExecuteWorkloadPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		log.Printf("ERROR: Failed to unmarshal payload for task %s: %v", t.Type(), err)
		return err
	}

	log.Printf("Received a job to execute workflow with ID: %d for user %s", p.WorkflowID, p.UserID)

	// 1. Fetch data from the database
	workflowDef, err := fetchWorkflowDefinition(p.WorkflowID)
	if err != nil {
		log.Printf("ERROR: Failed to fetch workflow definition for ID %d: %v", p.WorkflowID, err)
		return err
	}

	slackCred, err := fetchUserCredential(p.UserID, "slack")
	if err != nil {
		log.Printf("ERROR: Failed to fetch Slack credential for user %s: %v", p.UserID, err)
		return err
	}

	// 2. Simple execution logic: Find the Slack node and execute it
	for _, node := range workflowDef.Nodes {
		if node.Type == "slackAction" {
			channel, _ := node.Data["channel"].(string)
			message, _ := node.Data["message"].(string)

			log.Printf("Found Slack action. Sending '%s' to user %s", message, channel)

			// 3. Make the API call to Slack
			err := postToSlack(slackCred.AccessToken, channel, message)
			if err != nil {
				log.Printf("ERROR: Failed to post to Slack: %v", err)
				return err
			}
		}
	}

	log.Printf("Finished executing Workflow ID: %d", p.WorkflowID)
	return nil
}

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
		return fmt.Errorf("Slack API return non-200 status code: %s", response.Status)
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
	if err := connectToDb(); err != nil {
		log.Fatalf("FATAL: Failed to connect to database: %v", err)
	}
	defer dbConn.Close(context.Background())

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
