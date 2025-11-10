package task

import (
	"context"
	"errors"
	"strings"

	"DobrikaDev/task-service/internal/domain"
	searchintegration "DobrikaDev/task-service/internal/integration/search"
	"DobrikaDev/task-service/internal/storage/sql"

	"go.uber.org/zap"
)

func (s *TaskService) GetTaskByID(ctx context.Context, id string) (*domain.Task, error) {
	task, err := s.storage.GetTaskByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrTaskNotFound) {
			return nil, ErrTaskNotFound
		}
		s.logger.Error("failed to get task by id", zap.Error(err), zap.String("id", id))
		return nil, ErrTaskInternal
	}
	return task, nil
}

func (s *TaskService) GetTasks(ctx context.Context, customerID string, limit int, offset int) ([]*domain.Task, int, error) {
	opts := []sql.GetTasksOption{
		sql.WithTaskLimit(limit),
		sql.WithTaskOffset(offset),
		sql.WithTaskCustomerID(customerID),
	}
	tasks, count, err := s.storage.GetTasks(ctx, opts...)
	if err != nil {
		return nil, 0, ErrTaskInternal
	}
	return tasks, count, nil
}

func (s *TaskService) CountTasks(ctx context.Context, opts ...sql.GetTasksOption) (int, error) {
	count, err := s.storage.CountTasks(ctx, opts...)
	if err != nil {
		return 0, ErrTaskInternal
	}
	return count, nil
}

func (s *TaskService) CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	task, err := s.storage.CreateTask(ctx, task)
	if err != nil {
		if errors.Is(err, sql.ErrTaskAlreadyExists) {
			return nil, ErrTaskAlreadyExists
		}
		s.logger.Error("failed to create task", zap.Error(err), zap.Any("task", task))
		return nil, ErrTaskInternal
	}
	if s.indexer != nil {
		s.indexer.NotifyTaskChanged(task.ID)
	}
	return task, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	task, err := s.storage.UpdateTask(ctx, task)
	if err != nil {
		if errors.Is(err, sql.ErrTaskNotFound) {
			return nil, ErrTaskNotFound
		}
		s.logger.Error("failed to update task", zap.Error(err), zap.Any("task", task))
		return nil, ErrTaskInternal
	}
	if s.indexer != nil {
		s.indexer.NotifyTaskChanged(task.ID)
	}
	return task, nil
}

func (s *TaskService) SearchTasks(ctx context.Context, opts SearchOptions) ([]*domain.Task, error) {
	if strings.TrimSpace(opts.Query) == "" {
		return nil, ErrTaskInvalid
	}
	if s.search == nil {
		return nil, ErrTaskSearchUnavailable
	}

	req := searchintegration.SearchRequest{
		UserQuery: strings.TrimSpace(opts.Query),
		GeoData:   strings.TrimSpace(opts.GeoData),
		QueryType: strings.TrimSpace(opts.QueryType),
		UserTags:  make([]string, 0, len(opts.Tags)),
	}

	for _, tag := range opts.Tags {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			req.UserTags = append(req.UserTags, tag)
		}
	}

	resp, err := s.search.Search(ctx, req)
	if err != nil {
		s.logger.Error("search request failed", zap.Error(err), zap.String("query", req.UserQuery))
		return nil, ErrTaskInternal
	}

	if resp == nil || len(resp.TaskIDs) == 0 {
		return []*domain.Task{}, nil
	}

	tasks, err := s.storage.GetTasksByIDs(ctx, resp.TaskIDs)
	if err != nil {
		return nil, err
	}

	taskByID := make(map[string]*domain.Task, len(tasks))
	for _, task := range tasks {
		if task != nil && task.ID != "" {
			taskByID[task.ID] = task
		}
	}

	result := make([]*domain.Task, 0, len(resp.TaskIDs))
	for _, id := range resp.TaskIDs {
		if task, ok := taskByID[id]; ok {
			result = append(result, task)
		}
	}

	return result, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	err := s.storage.DeleteTask(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrTaskNotFound) {
			return ErrTaskNotFound
		}
		s.logger.Error("failed to delete task", zap.Error(err), zap.String("id", id))
		return ErrTaskInternal
	}
	return nil
}

func (s *TaskService) UserJoinTask(ctx context.Context, userID, taskID string) (*domain.UserTask, error) {
	userTask := &domain.UserTask{
		UserID: userID,
		TaskID: taskID,
		Status: domain.StatusInProgress,
	}

	task, err := s.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if task.CustomerID != userID {
		return nil, ErrUserTaskInvalid
	}

	userTask, err = s.storage.CreateUserTask(ctx, userTask)
	if err != nil {
		if errors.Is(err, sql.ErrUserTaskAlreadyExists) {
			return nil, ErrUserTaskAlreadyExists
		}
		if errors.Is(err, sql.ErrUserTaskInvalid) {
			return nil, ErrTaskNotFound
		}
		s.logger.Error("failed to create user task", zap.Error(err), zap.String("user_id", userID), zap.String("task_id", taskID))
		return nil, ErrTaskInternal
	}
	return userTask, nil
}

func (s *TaskService) UpdateUserTaskStatus(ctx context.Context, userID, taskID string, status domain.Status) (*domain.UserTask, error) {
	userTask, err := s.storage.UpdateUserTaskStatus(ctx, userID, taskID, status)
	if err != nil {
		if errors.Is(err, sql.ErrUserTaskNotFound) {
			return nil, ErrUserTaskNotFound
		}
		if errors.Is(err, sql.ErrUserTaskInvalid) {
			return nil, ErrTaskNotFound
		}
		s.logger.Error("failed to update user task status", zap.Error(err), zap.String("user_id", userID), zap.String("task_id", taskID), zap.Any("status", status))
		return nil, ErrTaskInternal
	}
	return userTask, nil
}
