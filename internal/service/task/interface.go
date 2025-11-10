package task

import (
	"DobrikaDev/task-service/internal/domain"
	"DobrikaDev/task-service/internal/storage/sql"
	"DobrikaDev/task-service/utils/config"
	"context"

	"go.uber.org/zap"
)

type storage interface {
	GetTaskByID(ctx context.Context, id string) (*domain.Task, error)
	GetTasks(ctx context.Context, opts ...sql.GetTasksOption) ([]*domain.Task, int, error)
	CountTasks(ctx context.Context, opts ...sql.GetTasksOption) (int, error)
	CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error)
	UpdateTask(ctx context.Context, task *domain.Task) (*domain.Task, error)
	DeleteTask(ctx context.Context, maxID string) error

	CreateUserTask(ctx context.Context, userTask *domain.UserTask) (*domain.UserTask, error)
	UpdateUserTaskStatus(ctx context.Context, userID, taskID string, status domain.Status) (*domain.UserTask, error)
}

type TaskService struct {
	storage storage
	cfg     *config.Config
	logger  *zap.Logger
}

func NewTaskService(storage storage, cfg *config.Config, logger *zap.Logger) *TaskService {
	return &TaskService{storage: storage, cfg: cfg, logger: logger}
}
