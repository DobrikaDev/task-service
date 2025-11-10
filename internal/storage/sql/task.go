package sql

import (
	"DobrikaDev/task-service/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

const (
	taskTableName            = "tasks"
	pgErrUniqueViolation     = "23505"
	pgErrForeignKeyViolation = "23503"
)

var taskSelectColumns = []string{
	"t.id",
	"t.customer_id",
	"t.name",
	"t.description",
	"t.verification_type",
	"t.cost",
	"t.members_count",
	"t.meta",
	"t.created_at",
	"t.updated_at",
}

type (
	taskOption interface {
		applySelect(sq.SelectBuilder) sq.SelectBuilder
		applyCount(sq.SelectBuilder) sq.SelectBuilder
	}

	taskOptionFunc struct {
		selectFn func(sq.SelectBuilder) sq.SelectBuilder
		countFn  func(sq.SelectBuilder) sq.SelectBuilder
	}
)

func (f taskOptionFunc) applySelect(sb sq.SelectBuilder) sq.SelectBuilder {
	if f.selectFn != nil {
		return f.selectFn(sb)
	}
	return sb
}

func (f taskOptionFunc) applyCount(sb sq.SelectBuilder) sq.SelectBuilder {
	if f.countFn != nil {
		return f.countFn(sb)
	}
	return sb
}

type GetTasksOption interface {
	taskOption
}

func WithTaskID(id string) GetTasksOption {
	return taskOptionFunc{
		selectFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			if id != "" {
				sb = sb.Where(sq.Eq{"t.id": id})
			}
			return sb
		},
		countFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			if id != "" {
				sb = sb.Where(sq.Eq{"t.id": id})
			}
			return sb
		},
	}
}

func WithTaskIDs(ids []string) GetTasksOption {
	return taskOptionFunc{
		selectFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			filtered := make([]string, 0, len(ids))
			for _, id := range ids {
				if id != "" {
					filtered = append(filtered, id)
				}
			}
			if len(filtered) > 0 {
				sb = sb.Where(sq.Eq{"t.id": filtered})
			}
			return sb
		},
		countFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			filtered := make([]string, 0, len(ids))
			for _, id := range ids {
				if id != "" {
					filtered = append(filtered, id)
				}
			}
			if len(filtered) > 0 {
				sb = sb.Where(sq.Eq{"t.id": filtered})
			}
			return sb
		},
	}
}

func WithTaskCustomerID(customerID string) GetTasksOption {
	return taskOptionFunc{
		selectFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			if customerID != "" {
				sb = sb.Where(sq.Eq{"t.customer_id": customerID})
			}
			return sb
		},
		countFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			if customerID != "" {
				sb = sb.Where(sq.Eq{"t.customer_id": customerID})
			}
			return sb
		},
	}
}

func WithTaskCustomerIDs(customerIDs []string) GetTasksOption {
	return taskOptionFunc{
		selectFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			filtered := make([]string, 0, len(customerIDs))
			for _, id := range customerIDs {
				if id != "" {
					filtered = append(filtered, id)
				}
			}
			if len(filtered) > 0 {
				sb = sb.Where(sq.Eq{"t.customer_id": filtered})
			}
			return sb
		},
		countFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			filtered := make([]string, 0, len(customerIDs))
			for _, id := range customerIDs {
				if id != "" {
					filtered = append(filtered, id)
				}
			}
			if len(filtered) > 0 {
				sb = sb.Where(sq.Eq{"t.customer_id": filtered})
			}
			return sb
		},
	}
}

func WithTaskName(name string) GetTasksOption {
	return taskOptionFunc{
		selectFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			if name != "" {
				sb = sb.Where(sq.Eq{"t.name": name})
			}
			return sb
		},
		countFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			if name != "" {
				sb = sb.Where(sq.Eq{"t.name": name})
			}
			return sb
		},
	}
}

func WithTaskNameLike(pattern string) GetTasksOption {
	return taskOptionFunc{
		selectFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			if pattern != "" {
				sb = sb.Where(sq.Expr("t.name ILIKE ?", fmt.Sprintf("%%%s%%", pattern)))
			}
			return sb
		},
		countFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			if pattern != "" {
				sb = sb.Where(sq.Expr("t.name ILIKE ?", fmt.Sprintf("%%%s%%", pattern)))
			}
			return sb
		},
	}
}

func WithTaskVerificationType(verificationType domain.VerificationType) GetTasksOption {
	return taskOptionFunc{
		selectFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			if verificationType != "" {
				sb = sb.Where(sq.Eq{"t.verification_type": verificationType})
			}
			return sb
		},
		countFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			if verificationType != "" {
				sb = sb.Where(sq.Eq{"t.verification_type": verificationType})
			}
			return sb
		},
	}
}

func WithTaskVerificationTypes(verificationTypes []domain.VerificationType) GetTasksOption {
	return taskOptionFunc{
		selectFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			filtered := make([]domain.VerificationType, 0, len(verificationTypes))
			for _, verificationType := range verificationTypes {
				if verificationType != "" {
					filtered = append(filtered, verificationType)
				}
			}
			if len(filtered) > 0 {
				sb = sb.Where(sq.Eq{"t.verification_type": filtered})
			}
			return sb
		},
		countFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			filtered := make([]domain.VerificationType, 0, len(verificationTypes))
			for _, verificationType := range verificationTypes {
				if verificationType != "" {
					filtered = append(filtered, verificationType)
				}
			}
			if len(filtered) > 0 {
				sb = sb.Where(sq.Eq{"t.verification_type": filtered})
			}
			return sb
		},
	}
}

func WithTaskLimit(limit int) GetTasksOption {
	return taskOptionFunc{
		selectFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			if limit > 0 {
				sb = sb.Limit(uint64(limit))
			}
			return sb
		},
	}
}

func WithTaskOffset(offset int) GetTasksOption {
	return taskOptionFunc{
		selectFn: func(sb sq.SelectBuilder) sq.SelectBuilder {
			if offset > 0 {
				sb = sb.Offset(uint64(offset))
			}
			return sb
		},
	}
}

func (s *SqlStorage) GetTaskByID(ctx context.Context, id string) (*domain.Task, error) {
	query, args := sq.Select(taskSelectColumns...).
		From(fmt.Sprintf("%s t", taskTableName)).
		Where(sq.Eq{"t.id": id}).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	var task domain.Task
	err := s.trf.Transaction(ctx).GetContext(ctx, &task, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		s.logger.Error("failed to get task by id", zap.Error(err), zap.String("task_id", id))
		return nil, ErrTaskInternal
	}

	return &task, nil
}

func (s *SqlStorage) GetTaskByMaxID(ctx context.Context, maxID string) (*domain.Task, error) {
	return s.GetTaskByID(ctx, maxID)
}

func (s *SqlStorage) GetTasks(ctx context.Context, opts ...GetTasksOption) ([]*domain.Task, int, error) {
	sb := sq.Select(taskSelectColumns...).
		From(fmt.Sprintf("%s t", taskTableName)).
		OrderBy("t.created_at DESC").
		PlaceholderFormat(sq.Dollar)

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		sb = opt.applySelect(sb)
	}

	query, args := sb.MustSql()

	tasks := make([]*domain.Task, 0)
	err := s.trf.Transaction(ctx).SelectContext(ctx, &tasks, query, args...)
	if err != nil {
		s.logger.Error("failed to get tasks", zap.Error(err))
		return nil, 0, ErrTaskInternal
	}

	count, err := s.CountTasks(ctx, opts...)
	if err != nil {
		s.logger.Error("failed to count tasks", zap.Error(err))
		return nil, 0, ErrTaskInternal
	}

	return tasks, count, nil
}

func (s *SqlStorage) CountTasks(ctx context.Context, opts ...GetTasksOption) (int, error) {
	sb := sq.Select("COUNT(*)").
		From(fmt.Sprintf("%s t", taskTableName)).
		PlaceholderFormat(sq.Dollar)

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		sb = opt.applyCount(sb)
	}

	query, args := sb.MustSql()

	var count int
	err := s.trf.Transaction(ctx).GetContext(ctx, &count, query, args...)
	if err != nil {
		return 0, ErrTaskInternal
	}

	return count, nil
}

func (s *SqlStorage) CreateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	query, args := sq.Insert(taskTableName).
		Columns(
			"id",
			"customer_id",
			"name",
			"description",
			"verification_type",
			"cost",
			"members_count",
			"meta",
		).
		Values(
			task.ID,
			task.CustomerID,
			task.Name,
			task.Description,
			task.VerificationType,
			task.Cost,
			task.MembersCount,
			task.Meta,
		).
		Suffix("RETURNING id, customer_id, name, description, verification_type, cost, members_count, meta, created_at, updated_at").
		PlaceholderFormat(sq.Dollar).
		MustSql()

	var created domain.Task
	err := s.trf.Transaction(ctx).GetContext(ctx, &created, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgErrUniqueViolation:
				return nil, ErrTaskAlreadyExists
			case pgErrForeignKeyViolation:
				return nil, ErrTaskInvalid
			}
		}
		s.logger.Error("failed to create task", zap.Error(err), zap.String("task_id", task.ID))
		return nil, ErrTaskInternal
	}

	return &created, nil
}

func (s *SqlStorage) UpdateTask(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	query, args := sq.Update(taskTableName).
		Set("name", task.Name).
		Set("description", task.Description).
		Set("verification_type", task.VerificationType).
		Set("cost", task.Cost).
		Set("members_count", task.MembersCount).
		Set("meta", task.Meta).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"id": task.ID}).
		Suffix("RETURNING id, customer_id, name, description, verification_type, cost, members_count, meta, created_at, updated_at").
		PlaceholderFormat(sq.Dollar).
		MustSql()

	var updated domain.Task
	err := s.trf.Transaction(ctx).GetContext(ctx, &updated, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgErrForeignKeyViolation {
			return nil, ErrTaskInvalid
		}

		s.logger.Error("failed to update task", zap.Error(err), zap.String("task_id", task.ID))
		return nil, ErrTaskInternal
	}

	return &updated, nil
}

func (s *SqlStorage) DeleteTask(ctx context.Context, id string) error {
	query, args := sq.Delete(taskTableName).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	result, err := s.trf.Transaction(ctx).ExecContext(ctx, query, args...)
	if err != nil {
		s.logger.Error("failed to delete task", zap.Error(err), zap.String("task_id", id))
		return ErrTaskInternal
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		s.logger.Error("failed to check affected rows when deleting task", zap.Error(err), zap.String("task_id", id))
		return ErrTaskInternal
	}

	if rowsAffected == 0 {
		return ErrTaskNotFound
	}

	return nil
}
