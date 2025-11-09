package delivery

import (
	"DobrikaDev/task-service/internal/domain"
	"DobrikaDev/task-service/internal/generated/proto/task"
	"context"
)

func (s *Server) UserJoinTask(ctx context.Context, req *task.UserJoinTaskRequest) (*task.UserJoinTaskResponse, error) {
	_, err := s.taskService.UserJoinTask(ctx, req.UserId, req.TaskId)
	if err != nil {
		return &task.UserJoinTaskResponse{
			Error: convertErrorToProto(err),
		}, nil
	}
	return &task.UserJoinTaskResponse{
		Error: convertErrorToProto(err),
	}, nil
}

func (s *Server) UserLeaveTask(ctx context.Context, req *task.UserLeaveTaskRequest) (*task.UserLeaveTaskResponse, error) {
	_, err := s.taskService.UpdateUserTaskStatus(ctx, req.UserId, req.TaskId, domain.StatusCancelled)
	if err != nil {
		return &task.UserLeaveTaskResponse{
			Error: convertErrorToProto(err),
		}, nil
	}
	return &task.UserLeaveTaskResponse{
		Error: convertErrorToProto(err),
	}, nil
}

func (s *Server) UserConfirmTask(ctx context.Context, req *task.UserConfirmTaskRequest) (*task.UserConfirmTaskResponse, error) {
	_, err := s.taskService.UpdateUserTaskStatus(ctx, req.UserId, req.TaskId, domain.StatusCompleted)
	if err != nil {
		return &task.UserConfirmTaskResponse{
			Error: convertErrorToProto(err),
		}, nil
	}
	return &task.UserConfirmTaskResponse{
		Error: convertErrorToProto(err),
	}, nil
}

func (s *Server) ApproveTask(ctx context.Context, req *task.ApproveTaskRequest) (*task.ApproveTaskResponse, error) {
	_, err := s.taskService.UpdateUserTaskStatus(ctx, req.UserId, req.TaskId, domain.StatusApproved)
	if err != nil {
		return &task.ApproveTaskResponse{
			Error: convertErrorToProto(err),
		}, nil
	}
	return &task.ApproveTaskResponse{
		Error: convertErrorToProto(err),
	}, nil
}

func (s *Server) RejectTask(ctx context.Context, req *task.RejectTaskRequest) (*task.RejectTaskResponse, error) {
	_, err := s.taskService.UpdateUserTaskStatus(ctx, req.UserId, req.TaskId, domain.StatusRejected)
	if err != nil {
		return &task.RejectTaskResponse{
			Error: convertErrorToProto(err),
		}, nil
	}
	return &task.RejectTaskResponse{
		Error: convertErrorToProto(err),
	}, nil
}
