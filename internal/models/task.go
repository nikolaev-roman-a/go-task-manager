package models

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusCanceled  TaskStatus = "canceled"
	StatusFailed    TaskStatus = "failed"
)

type Task struct {
	ID         uuid.UUID  `json:"id"`
	Status     TaskStatus `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	StartedAt  *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
	Result     string     `json:"result,omitempty"`
}

type ErrorResponse struct {
	Error string
}
