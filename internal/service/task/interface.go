package task

import (
	"context"
	"time"

	"DobrikaDev/task-service/internal/domain"
	searchintegration "DobrikaDev/task-service/internal/integration/search"
	"DobrikaDev/task-service/internal/storage/sql"
	"DobrikaDev/task-service/utils/config"

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

	GetTasksByIDs(ctx context.Context, ids []string) ([]*domain.Task, error)
	GetTasksUpdatedAfter(ctx context.Context, after time.Time, limit int) ([]*domain.Task, error)
	LoadSearchCursor(ctx context.Context) (time.Time, error)
	SaveSearchCursor(ctx context.Context, cursor time.Time) error
}

type indexer interface {
	NotifyTaskChanged(taskID string)
}

type searchClient interface {
	Search(ctx context.Context, req searchintegration.SearchRequest) (*searchintegration.SearchResponse, error)
}

type TaskService struct {
	storage storage
	cfg     *config.Config
	logger  *zap.Logger
	indexer indexer
	search  searchClient
}

func NewTaskService(storage storage, cfg *config.Config, logger *zap.Logger, indexer indexer, search searchClient) *TaskService {
	return &TaskService{
		storage: storage,
		cfg:     cfg,
		logger:  logger,
		indexer: indexer,
		search:  search,
	}
}

type SearchOptions struct {
	Query     string
	QueryType string
	GeoData   string
	Tags      []string
}
