package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var dbConn *pgx.Conn

func connectToDb() error {
	var err error
	dbConn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("Failed to connect to database: %v", err)
	}
	return nil
}

func fetchWorkflowDefinition(workflowID int) (*WorkflowDefinition, error) {
	var defJSON []byte

	err := dbConn.QueryRow(context.Background(), "SELECT definition FROM workflow_definitions WHERE id = $1", workflowID).Scan(&defJSON)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch workflow definition: %v", err)
	}

	var workflowDef WorkflowDefinition

	err = json.Unmarshal(defJSON, &workflowDef)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal workflow definition: %v", err)
	}
	return &workflowDef, nil
}

func fetchUserCredential(userID string, service string) (*UserCredential, error) {
	var cred UserCredential

	err := dbConn.QueryRow(context.Background(), "SELECT access_token FROM user_credentials WHERE user_id = $1 AND service = $2", userID, service).Scan(&cred.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch user credential: %v", err)
	}
	return &cred, nil
}
