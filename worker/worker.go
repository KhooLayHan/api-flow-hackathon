package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TaskProcessor is a struct that holds all the dependencies our task handlers will need.
type TaskProcessor struct {
	db  *pgxpool.Pool
	srv *asynq.Server
}

// TestMessagePayload must match the structure of the payload sent by the Nuxt API.
type TestMessagePayload struct {
	Message string `json:"message"`
	SentAt  string `json:"sentAt"`
}

// ExecuteWorkloadPayload must match the structure of the payload sent by the Nuxt API.
type ExecuteWorkloadPayload struct {
	WorkflowID int    `json:"workflowId"`
	UserID     string `json:"userId"`
}

// NewTaskProcessor creates a new TaskProcessor instance.
func NewTaskProcessor(redisOpt asynq.RedisClientOpt, dbPool *pgxpool.Pool) *TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{"workflows": 1},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Printf("Asynq Error: type=%s, payload=%s, err=%w", task.Type(), task.Payload(), err)
		}),
	})

	return &TaskProcessor{
		db:  dbPool,
		srv: server,
	}
}

func (p *TaskProcessor) fetchWorkflowDefinition(ctx context.Context, workflowID int) (*WorkflowDefinition, error) {
	var defJSON []byte

	err := p.db.QueryRow(ctx, "SELECT definition FROM workflow_definitions WHERE id = $1", workflowID).
		Scan(&defJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch workflow definition: %w", err)
	}

	var workflowDef WorkflowDefinition

	err = json.Unmarshal(defJSON, &workflowDef)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal workflow definition: %w", err)
	}
	return &workflowDef, nil
}

func (p *TaskProcessor) fetchUserCredential(ctx context.Context, userID string, service string) (*UserCredential, error) {
	var cred UserCredential

	err := p.db.QueryRow(ctx, "SELECT access_token FROM user_credentials WHERE user_id = $1 AND service = $2", userID, service).
		Scan(&cred.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user credential: %w", err)
	}
	return &cred, nil
}

// handleExecuteWorkflowTask handles the execution of a workflow task.
func (p *TaskProcessor) handleExecuteWorkflowTask(ctx context.Context, t *asynq.Task) error {
	var payload ExecuteWorkloadPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		log.Printf("ERROR: Failed to unmarshal payload for task %s: %v", t.Type(), err)
		return err
	}

	log.Printf("Received a job to execute workflow with ID: %d for user %s", payload.WorkflowID, payload.UserID)

	// 1. Fetch data from the database
	workflowDef, err := p.fetchWorkflowDefinition(ctx, payload.WorkflowID)
	if err != nil {
		log.Printf("ERROR: Failed to fetch workflow definition for ID %d: %v", payload.WorkflowID, err)
		return err
	}

	slackCred, err := p.fetchUserCredential(ctx, payload.UserID, "slack")
	if err != nil {
		log.Printf("ERROR: Failed to fetch Slack credential for user %s: %v", payload.UserID, err)
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

	log.Printf("Finished executing Workflow ID: %d", payload.WorkflowID)
	return nil
}

// Starts the worker server
func (p *TaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc("execute-workflow-v1", p.handleExecuteWorkflowTask)

	log.Println("Go Worker service started. Listening for jobs...")
	return p.srv.Run(mux)
}
