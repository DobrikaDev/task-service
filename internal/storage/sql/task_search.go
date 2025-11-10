package sql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"DobrikaDev/task-service/internal/domain"

	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

const (
	searchStateTable = "search_index_state"
	defaultCursorID  = 1
)

func (s *SqlStorage) GetTasksUpdatedAfter(ctx context.Context, after time.Time, limit int) ([]*domain.Task, error) {
	sb := sq.Select(taskSelectColumns...).
		From(fmt.Sprintf("%s t", taskTableName)).
		Where(sq.Gt{"t.updated_at": after}).
		OrderBy("t.updated_at ASC").
		PlaceholderFormat(sq.Dollar)

	if limit > 0 {
		sb = sb.Limit(uint64(limit))
	}

	query, args := sb.MustSql()

	tasks := make([]*domain.Task, 0)
	if err := s.trf.Transaction(ctx).SelectContext(ctx, &tasks, query, args...); err != nil {
		s.logger.Error("failed to get tasks updated after cursor", zap.Error(err), zap.Time("after", after))
		return nil, ErrTaskInternal
	}

	return tasks, nil
}

func (s *SqlStorage) GetTasksByIDs(ctx context.Context, ids []string) ([]*domain.Task, error) {
	filtered := make([]string, 0, len(ids))
	for _, id := range ids {
		if id != "" {
			filtered = append(filtered, id)
		}
	}
	if len(filtered) == 0 {
		return []*domain.Task{}, nil
	}

	query, args := sq.Select(taskSelectColumns...).
		From(fmt.Sprintf("%s t", taskTableName)).
		Where(sq.Eq{"t.id": filtered}).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	tasks := make([]*domain.Task, 0, len(filtered))
	if err := s.trf.Transaction(ctx).SelectContext(ctx, &tasks, query, args...); err != nil {
		s.logger.Error("failed to get tasks by ids", zap.Error(err), zap.Strings("task_ids", filtered))
		return nil, ErrTaskInternal
	}

	return tasks, nil
}

func (s *SqlStorage) LoadSearchCursor(ctx context.Context) (time.Time, error) {
	query, args := sq.Select("last_synced_at").
		From(searchStateTable).
		Where(sq.Eq{"id": defaultCursorID}).
		PlaceholderFormat(sq.Dollar).
		MustSql()

	var cursor time.Time
	err := s.trf.Transaction(ctx).GetContext(ctx, &cursor, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return time.Time{}, nil
		}
		s.logger.Error("failed to load search cursor", zap.Error(err))
		return time.Time{}, ErrTaskInternal
	}

	return cursor, nil
}

func (s *SqlStorage) SaveSearchCursor(ctx context.Context, cursor time.Time) error {
	query, args := sq.Insert(searchStateTable).
		Columns("id", "last_synced_at").
		Values(defaultCursorID, cursor).
		Suffix("ON CONFLICT (id) DO UPDATE SET last_synced_at = EXCLUDED.last_synced_at").
		PlaceholderFormat(sq.Dollar).
		MustSql()

	if _, err := s.trf.Transaction(ctx).ExecContext(ctx, query, args...); err != nil {
		s.logger.Error("failed to save search cursor", zap.Error(err), zap.Time("cursor", cursor))
		return ErrTaskInternal
	}

	return nil
}
