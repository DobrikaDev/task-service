package indexer

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"time"

	"DobrikaDev/task-service/internal/domain"
	searchintegration "DobrikaDev/task-service/internal/integration/search"
	"DobrikaDev/task-service/utils/config"

	"go.uber.org/zap"
)

type Storage interface {
	GetTasksUpdatedAfter(ctx context.Context, after time.Time, limit int) ([]*domain.Task, error)
	GetTasksByIDs(ctx context.Context, ids []string) ([]*domain.Task, error)
	LoadSearchCursor(ctx context.Context) (time.Time, error)
	SaveSearchCursor(ctx context.Context, cursor time.Time) error
}

type searchClient interface {
	IndexTask(ctx context.Context, task searchintegration.IndexTask) error
}

type Scheduler struct {
	storage Storage
	client  searchClient
	cfg     config.SearchConfig
	logger  *zap.Logger

	updates   chan string
	startOnce sync.Once
	stopOnce  sync.Once

	ctx    context.Context
	cancel context.CancelFunc
}

func NewScheduler(storage Storage, client searchClient, cfg config.SearchConfig, logger *zap.Logger) *Scheduler {
	if logger == nil {
		logger = zap.NewNop()
	}

	buffer := cfg.SchedulerBatchSize
	if buffer <= 0 {
		buffer = 200
	}

	return &Scheduler{
		storage: storage,
		client:  client,
		cfg:     cfg,
		logger:  logger,
		updates: make(chan string, buffer*2),
	}
}

func (s *Scheduler) Start(parent context.Context) {
	if s.client == nil || s.storage == nil {
		s.logger.Warn("search scheduler not started: missing dependencies")
		return
	}

	s.startOnce.Do(func() {
		if parent == nil {
			parent = context.Background()
		}

		s.ctx, s.cancel = context.WithCancel(parent)
		go s.loop()
	})
}

func (s *Scheduler) Stop() {
	s.stopOnce.Do(func() {
		if s.cancel != nil {
			s.cancel()
		}
	})
}

func (s *Scheduler) NotifyTaskChanged(taskID string) {
	if taskID == "" || s.client == nil {
		return
	}

	select {
	case s.updates <- taskID:
	default:
		s.logger.Warn("search scheduler queue is full, dropping task id", zap.String("task_id", taskID))
	}
}

func (s *Scheduler) loop() {
	interval := s.cfg.SchedulerInterval
	if interval <= 0 {
		interval = 30 * time.Second
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	pending := make(map[string]int)

	for {
		select {
		case <-s.ctx.Done():
			return
		case taskID := <-s.updates:
			if taskID != "" {
				pending[taskID] = pending[taskID]
			}
		case <-ticker.C:
			pending = s.process(pending)
		}
	}
}

func (s *Scheduler) process(pending map[string]int) map[string]int {
	ctx := s.ctx
	if ctx == nil {
		return pending
	}

	maxRetries := s.cfg.SchedulerMaxRetries
	if maxRetries <= 0 {
		maxRetries = 3
	}

	nextPending := make(map[string]int, len(pending))
	for id, count := range pending {
		nextPending[id] = count
	}

	processed := make(map[string]struct{})

	if len(pending) > 0 {
		ids := make([]string, 0, len(pending))
		for id := range pending {
			ids = append(ids, id)
		}

		tasks, err := s.storage.GetTasksByIDs(ctx, ids)
		if err != nil {
			s.logger.Error("failed to fetch pending tasks for indexing", zap.Error(err))
			return nextPending
		}

		taskMap := make(map[string]*domain.Task, len(tasks))
		for _, task := range tasks {
			if task != nil && task.ID != "" {
				taskMap[task.ID] = task
			}
		}

		for _, id := range ids {
			task := taskMap[id]
			if task == nil {
				continue
			}
			if s.indexTask(ctx, task) {
				delete(nextPending, id)
				processed[id] = struct{}{}
				continue
			}
			retries := pending[id] + 1
			if retries <= maxRetries {
				nextPending[id] = retries
			} else {
				delete(nextPending, id)
				s.logger.Error("max retries exceeded when indexing task", zap.String("task_id", id))
			}
		}
	}

	cursor, err := s.storage.LoadSearchCursor(ctx)
	if err != nil {
		s.logger.Error("failed to load search cursor", zap.Error(err))
		return nextPending
	}

	batchSize := s.cfg.SchedulerBatchSize
	if batchSize <= 0 {
		batchSize = 200
	}

	fetchCursor := cursor
	maxIndexedTime := cursor

	for {
		tasks, err := s.storage.GetTasksUpdatedAfter(ctx, fetchCursor, batchSize)
		if err != nil {
			s.logger.Error("failed to fetch tasks for indexing", zap.Error(err))
			break
		}
		if len(tasks) == 0 {
			break
		}

		for _, task := range tasks {
			if task == nil || task.ID == "" {
				continue
			}
			if _, alreadyProcessed := processed[task.ID]; alreadyProcessed {
				continue
			}
			if _, pendingAgain := nextPending[task.ID]; pendingAgain {
				continue
			}

			if s.indexTask(ctx, task) {
				processed[task.ID] = struct{}{}
				if task.UpdatedAt.After(maxIndexedTime) {
					maxIndexedTime = task.UpdatedAt
				}
			} else {
				nextPending[task.ID] = 1
			}
		}

		last := tasks[len(tasks)-1]
		if last == nil || !last.UpdatedAt.After(fetchCursor) {
			break
		}
		fetchCursor = last.UpdatedAt

		if len(tasks) < batchSize {
			break
		}
	}

	if maxIndexedTime.After(cursor) {
		if err := s.storage.SaveSearchCursor(ctx, maxIndexedTime); err != nil {
			s.logger.Error("failed to save search cursor", zap.Error(err))
		}
	}

	return nextPending
}

func (s *Scheduler) indexTask(ctx context.Context, task *domain.Task) bool {
	payload := searchintegration.IndexTask{
		TaskID:   task.ID,
		TaskName: strings.TrimSpace(task.Name),
		TaskDesc: strings.TrimSpace(task.Description),
		TaskType: taskTypeFromMeta(task.Meta),
		GeoData:  geoFromMeta(task.Meta),
	}

	if err := s.client.IndexTask(ctx, payload); err != nil {
		s.logger.Error("failed to index task in search", zap.Error(err), zap.String("task_id", task.ID))
		return false
	}

	return true
}

func taskTypeFromMeta(meta json.RawMessage) string {
	data := metaToMap(meta)
	if len(data) == 0 {
		return ""
	}
	if value := strings.TrimSpace(data["task_type"]); value != "" {
		return value
	}
	if value := strings.TrimSpace(data["type"]); value != "" {
		return value
	}
	return ""
}

func geoFromMeta(meta json.RawMessage) string {
	data := metaToMap(meta)
	if len(data) == 0 {
		return ""
	}
	if value := strings.TrimSpace(data["geo_data"]); value != "" {
		return value
	}
	if value := strings.TrimSpace(data["geo"]); value != "" {
		return value
	}
	lat, latOK := data["lat"]
	lon, lonOK := data["lon"]
	if latOK && lonOK {
		return strings.TrimSpace(lat) + "," + strings.TrimSpace(lon)
	}
	return ""
}

func metaToMap(meta json.RawMessage) map[string]string {
	if len(meta) == 0 || string(meta) == "null" {
		return nil
	}

	var stringMap map[string]string
	if err := json.Unmarshal(meta, &stringMap); err == nil {
		return stringMap
	}

	var generic map[string]any
	if err := json.Unmarshal(meta, &generic); err != nil {
		return nil
	}

	result := make(map[string]string, len(generic))
	for key, value := range generic {
		switch v := value.(type) {
		case string:
			result[key] = v
		case float64:
			result[key] = strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			result[key] = strconv.FormatBool(v)
		}
	}

	return result
}

