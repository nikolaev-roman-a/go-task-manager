package repository

import (
	"github.com/google/uuid"
	"github.com/nikolaev-roman-a/go-task-manager/internal/models"
)

type Repository interface {
	Save(*models.Task) error
	Get(uuid.UUID) (*models.Task, error)
	Search() ([]*models.Task, error)
	Delete(uuid.UUID) error
}
