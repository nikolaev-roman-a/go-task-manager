package services

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nikolaev-roman-a/go-task-manager/internal/models"
	"github.com/nikolaev-roman-a/go-task-manager/internal/repository"
	"go.uber.org/zap"
)

type TaskService struct {
	mu      sync.RWMutex
	running map[uuid.UUID]context.CancelFunc
	store   repository.Repository
	logger  *zap.Logger
}

func NewTaskService(store repository.Repository, logger *zap.Logger) *TaskService {
	return &TaskService{
		store:   store,
		running: make(map[uuid.UUID]context.CancelFunc),
		logger:  logger,
	}
}

// crud

func (s *TaskService) Create(ctx context.Context, task *models.Task) (*models.Task, error) {

	task.ID = uuid.New()
	task.Status = models.StatusPending
	task.CreatedAt = time.Now().UTC()

	err := s.store.Save(task)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) Read(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	return s.store.Get(id)
}

func (s *TaskService) Search(ctx context.Context) ([]*models.Task, error) {
	return s.store.Search()
}

func (s *TaskService) Update(ctx context.Context, task *models.Task) (*models.Task, error) {

	s.store.Save(task)

	return task, nil
}

func (s *TaskService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.store.Delete(id)
}

// executor

func (s *TaskService) Run(ctx context.Context, task *models.Task) error {
	if _, ok := s.running[task.ID]; ok {
		return errors.New("already running")
	}

	taskCtx, cancel := context.WithCancel(context.Background())

	s.mu.Lock()
	defer s.mu.Unlock()
	s.running[task.ID] = cancel

	go func() {
		defer func() {
			if r := recover(); r != nil {
				return
			}
		}()

		s.setTaskRunning(taskCtx, task)
		s.logger.Info("task")
		result, err := s.process(taskCtx, task)
		if err != nil {
			s.setTaskResult(ctx, task, err.Error(), models.StatusFailed)
		}
		if result != "" {
			s.setTaskResult(ctx, task, result, models.StatusCompleted)
		}
	}()
	return nil
}

func (s *TaskService) setTaskRunning(ctx context.Context, task *models.Task) {
	task.Status = models.StatusRunning
	currentTime := time.Now().UTC()
	task.StartedAt = &currentTime
	s.Update(ctx, task)
}

func (s *TaskService) setTaskResult(ctx context.Context, task *models.Task, result string, status models.TaskStatus) {
	currentTime := time.Now().UTC()
	task.FinishedAt = &currentTime
	task.Status = status
	task.Result = result
	s.Update(ctx, task)
}

func (s *TaskService) Cancel(ctx context.Context, id uuid.UUID) error {
	task, err := s.Read(ctx, id)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	cancel, ok := s.running[id]
	if !ok {
		return nil
	}

	if task.Status == models.StatusRunning && ok {
		cancel()
		delete(s.running, task.ID)
	}

	s.Delete(ctx, task.ID)

	return nil
}

func (s *TaskService) process(ctx context.Context, task *models.Task) (string, error) {
	duration := time.Duration(10+rand.Intn(30)) * time.Second

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case t := <-ticker.C:
				s.logger.Info(fmt.Sprintf("Task %s - Current time: %s", task.ID, t.Format(time.RFC3339)))
			}
		}
	}()

	select {
	case <-ctx.Done():
		return "", nil
	case <-time.After(duration):
		if rand.Float32() < 0.1 {
			return "", errors.New("I/O operation failed")
		}
		return "Result", nil
	}
}

// manager

func (s *TaskService) CreateAndRun(ctx context.Context, task *models.Task) (*models.Task, error) {
	created, err := s.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	err = s.Run(ctx, created)
	if err != nil {
		return nil, err
	}

	return created, nil
}
