package model

import "time"

const (
	TaskStatusPending   = "pending"
	TaskStatusRunning   = "running"
	TaskStatusCompleted = "completed"
	TaskStatusFailed    = "failed"
)

type ResearchTask struct {
	ID          string         `json:"task_id"`
	Question    string         `json:"question"`
	Status      string         `json:"status"`
	CurrentStep string         `json:"current_step,omitempty"`
	Plan        map[string]any `json:"plan,omitempty"`
	Report      string         `json:"report,omitempty"`
	Error       string         `json:"error,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
