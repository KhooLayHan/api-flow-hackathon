package main

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
