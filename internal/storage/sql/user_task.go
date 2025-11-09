package sql

import (
	"DobrikaDev/task-service/internal/domain"
	"context"
	"database/sql"
	"errors"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

const userTaskTableName = "user_tasks"

var userTaskSelectColumns = []string{
	"user_id",
	"task_id",
	"status",
	"created_at",
	"updated_at",
}

func (s *SqlStorage) CreateUserTask(ctx context.Context, userTask *domain.UserTask) (*domain.UserTask, error) {
	query, args := sq.Insert(userTaskTableName).
		Columns("user_id", "task_id", "status").
		Values(userTask.UserID, userTask.TaskID, userTask.Status).
		Suffix("RETURNING " + strings.Join(userTaskSelectColumns, ", ")).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	var created domain.UserTask
	err := s.trf.Transaction(ctx).GetContext(ctx, &created, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgErrUniqueViolation:
				return nil, ErrUserTaskAlreadyExists
			case pgErrForeignKeyViolation:
				return nil, ErrUserTaskInvalid
			}
		}

		s.logger.Error(
			"failed to create user task",
			zap.Error(err),
			zap.String("user_id", userTask.UserID),
			zap.String("task_id", userTask.TaskID),
		)

		return nil, ErrUserTaskInternal
	}

	return &created, nil
}

func (s *SqlStorage) UpdateUserTaskStatus(ctx context.Context, userID, taskID string, status domain.Status) (*domain.UserTask, error) {
	query, args := sq.Update(userTaskTableName).
		Set("status", status).
		Set("updated_at", sq.Expr("NOW()")).
		Where(sq.Eq{"user_id": userID, "task_id": taskID}).
		Suffix("RETURNING " + strings.Join(userTaskSelectColumns, ", ")).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	var updated domain.UserTask
	err := s.trf.Transaction(ctx).GetContext(ctx, &updated, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserTaskNotFound
		}

		s.logger.Error(
			"failed to update user task status",
			zap.Error(err),
			zap.String("user_id", userID),
			zap.String("task_id", taskID),
			zap.String("status", status.String()),
		)

		return nil, ErrUserTaskInternal
	}

	return &updated, nil
}
