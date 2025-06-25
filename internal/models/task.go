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
	ID         uuid.UUID
	Status     TaskStatus
	CreatedAt  time.Time
	StartedAt  time.Time
	FinishedAt time.Time
	Result     []byte
}

type ErrorResponse struct {
	Error string
}
