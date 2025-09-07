package main

import (
	"log/slog"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
)

// WorkflowDefinition represents a workflow definition; mirrors the structure of the API.
type WorkflowDefinition struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

type Node struct {
	ID       string             `json:"id"`
	Type     string             `json:"type"`
	Data     map[string]any     `json:"data"`
	Position map[string]float64 `json:"position"`
}

type Edge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

// UserCredential represents to hold the saved credentials.
type UserCredential struct {
	ID          string
	UserID      string
	Service     string
	AccessToken string
}

// TaskProcessor is a struct that holds all the dependencies our task handlers will need.
type TaskProcessor struct {
	db     *pgxpool.Pool
	srv    *asynq.Server
	logger *slog.Logger
}

// TestMessagePayload must match the structure of the payload sent by the Nuxt API.
type TestMessagePayload struct {
	Message string `json:"message"`
	SentAt  string `json:"sentAt"`
}

type ExecuteWorkloadPayload struct {
	WorkflowID     string            `json:"workflow_id"`
	UserID         string            `json:"user_id"`
	TriggerPayload GitHubPushPayload `json:"trigger_payload"`
}

type GitHubPushPayload struct {
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
	Pusher struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"pusher"`
	HeadCommit struct {
		Message string `json:"message"`
		URL     string `json:"url"`
	} `json:"head_commit"`
}
