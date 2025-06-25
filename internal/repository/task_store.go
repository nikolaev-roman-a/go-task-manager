package repository

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/nikolaev-roman-a/go-task-manager/internal/models"
)

type TaskStore struct {
	mu    sync.RWMutex
	tasks map[uuid.UUID]*models.Task
}

func NewTaskStore() *TaskStore {
	return &TaskStore{
		tasks: make(map[uuid.UUID]*models.Task),
	}
}

func (s *TaskStore) Save(task *models.Task) error {
	taskCopy := &models.Task{}
	*taskCopy = *task

	s.mu.Lock()
	defer s.mu.Unlock()

	s.tasks[task.ID] = taskCopy
	return nil
}

func (s *TaskStore) Get(id uuid.UUID) (*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	if !ok {
		return nil, errors.New("not found")
	}

	taskCopy := &models.Task{}
	*taskCopy = *task

	return taskCopy, nil
}

func (s *TaskStore) Search() ([]*models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]*models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *TaskStore) Delete(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.tasks, id)
	return nil
}
