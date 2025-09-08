package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
)

// NewTaskProcessor creates a new TaskProcessor instance.
func NewTaskProcessor(redisOpt asynq.RedisClientOpt, dbPool *pgxpool.Pool, logger *slog.Logger) *TaskProcessor {
	server := asynq.NewServer(redisOpt, asynq.Config{
		Queues: map[string]int{"workflows": 1},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			logger.ErrorContext(ctx, "Asynq task processing failed.",
				"error", err.Error(),
				"task_type", task.Type(),
				"task_payload", string(task.Payload()))
		}),
	})

	return &TaskProcessor{
		db:     dbPool,
		srv:    server,
		logger: logger,
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

func (p *TaskProcessor) fetchUserCredential(
	ctx context.Context,
	userID string,
	service string,
) (*UserCredential, error) {
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
	if payloadErr := json.Unmarshal(t.Payload(), &payload); payloadErr != nil {
		p.logger.ErrorContext(ctx, "Failed to unmarshal payload", "error", payloadErr)
		return payloadErr
	}

	p.logger.InfoContext(
		ctx,
		"Received a job to execute workflow",
		"workflowID",
		payload.WorkflowID,
		"userID",
		payload.UserID,
	)

	// 1. Fetch data from the database.
	workflowDef, workflowDefErr := p.fetchWorkflowDefinition(ctx, payload.WorkflowID)
	if workflowDefErr != nil {
		p.logger.ErrorContext(
			ctx,
			"Failed to fetch workflow definition",
			"workflowID",
			payload.WorkflowID,
			"error",
			workflowDefErr,
		)
		return workflowDefErr
	}

	slackCred, slackCredErr := p.fetchUserCredential(ctx, payload.UserID, "slack")
	if slackCredErr != nil {
		p.logger.ErrorContext(ctx, "Failed to fetch Slack credential", "userID", payload.UserID, "error", slackCredErr)
		return slackCredErr
	}

	// 2. Simple execution logic: Find the Slack node and execute it.
	for _, node := range workflowDef.Nodes {
		if node.Type == "slackAction" {
			channel, _ := node.Data["channel"].(string)
			messageTemplate, _ := node.Data["message"].(string)

			// TODO: Uses simple string replacement for now, will be updated with a proper templating engine.
			finalMessage := strings.ReplaceAll(messageTemplate, "{commit.message}", payload.TriggerPayload.HeadCommit.Message)
			finalMessage = strings.ReplaceAll(finalMessage, "{pusher.name}", payload.TriggerPayload.Pusher.Name)
			finalMessage = strings.ReplaceAll(finalMessage, "{repo.name}", payload.TriggerPayload.Repository.FullName)
			finalMessage = fmt.Sprintf("%s\nView Commit: %s", finalMessage, payload.TriggerPayload.HeadCommit.URL)

			p.logger.InfoContext(
				ctx,
				"Found Slack action and sending",
				"workflowID",
				payload.WorkflowID,
				"message",
				messageTemplate,
				"channel",
				channel,
			)

			// 3. Make the API call to Slack.
			err := postToSlack(slackCred.AccessToken, channel, finalMessage)
			if err != nil {
				p.logger.ErrorContext(ctx, "Failed to post to Slack", "workflowID", payload.WorkflowID, "error", err)
				return err
			}
		}
	}

	p.logger.InfoContext(ctx, "Finished executing dynamic workflow successfully", "workflowID", payload.WorkflowID)
	return nil
}

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

// Start starts the worker server.
func (p *TaskProcessor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc("execute-workflow-v1", p.handleExecuteWorkflowTask)

	p.logger.Info("Go Worker service started. Listening for jobs...")
	return p.srv.Run(mux)
}
