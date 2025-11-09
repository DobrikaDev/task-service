package task

import (
	"DobrikaDev/task-service/internal/domain"
	"DobrikaDev/task-service/internal/storage/sql"
	"context"
	"errors"

	"go.uber.org/zap"
)

func (s *TaskService) GetTaskByMaxID(ctx context.Context, maxID string) (*domain.Task, error) {
	task, err := s.storage.GetTaskByMaxID(ctx, maxID)
	if err != nil {
		if errors.Is(err, sql.ErrTaskNotFound) {
			return nil, ErrTaskNotFound
		}
		s.logger.Error("failed to get task by max id", zap.Error(err), zap.String("max_id", maxID))
		return nil, ErrTaskInternal
	}
	return task, nil
}

func (s *TaskService) GetTasks(ctx context.Context, maxID string, limit int, offset int) ([]*domain.Task, int, error) {
	opts := []sql.GetTasksOption{
		sql.WithTaskMaxID(maxID),
		sql.WithTaskLimit(limit),
		sql.WithTaskOffset(offset),
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
	return task, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, maxID string) error {
	err := s.storage.DeleteTask(ctx, maxID)
	if err != nil {
		if errors.Is(err, sql.ErrTaskNotFound) {
			return ErrTaskNotFound
		}
		s.logger.Error("failed to delete task", zap.Error(err), zap.String("max_id", maxID))
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
	userTask, err := s.storage.CreateUserTask(ctx, userTask)
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