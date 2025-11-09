package delivery

import (
	taskpb "DobrikaDev/task-service/internal/generated/proto/task"
	"DobrikaDev/task-service/internal/service/task"
	"DobrikaDev/task-service/utils/config"
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	taskService *task.TaskService
	taskpb.UnimplementedTaskServiceServer

	cfg    *config.Config
	logger *zap.Logger
}

func NewServer(ctx context.Context, taskService *task.TaskService, cfg *config.Config, logger *zap.Logger) *Server {
	server := &Server{taskService: taskService, cfg: cfg, logger: logger}
	return server
}

func (s *Server) Register(grpcServer *grpc.Server) {
	taskpb.RegisterTaskServiceServer(grpcServer, s)
}
