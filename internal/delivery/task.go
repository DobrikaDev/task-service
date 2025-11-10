package delivery

import (
	"DobrikaDev/task-service/internal/domain"
	taskpb "DobrikaDev/task-service/internal/generated/proto/task"
	"DobrikaDev/task-service/internal/service/task"
	"context"
	"encoding/json"
	"errors"

	"github.com/dr3dnought/gospadi"
	"go.uber.org/zap"
)

func (s *Server) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.CreateTaskResponse, error) {
	if req.GetTask() == nil {
		return &taskpb.CreateTaskResponse{
			Error: validationError("task is required"),
		}, nil
	}

	payload := req.GetTask()

	if payload.Id == "" {
		return &taskpb.CreateTaskResponse{
			Error: validationError("id is required"),
		}, nil
	}

	if payload.CustomerId == "" {
		return &taskpb.CreateTaskResponse{
			Error: validationError("customer id is required"),
		}, nil
	}

	if payload.Name == "" {
		return &taskpb.CreateTaskResponse{
			Error: validationError("name is required"),
		}, nil
	}

	metaJSON, err := convertProtoMetaToJSON(payload.Meta)
	if err != nil {
		return &taskpb.CreateTaskResponse{
			Error: validationError("meta is invalid"),
		}, nil
	}

	task := &domain.Task{
		ID:               payload.Id,
		CustomerID:       payload.CustomerId,
		Name:             payload.Name,
		Description:      payload.Description,
		VerificationType: convertVerificationTypeToDomain(payload.VerificationType),
		Cost:             int(payload.Cost),
		MembersCount:     int(payload.MembersCount),
		Meta:             metaJSON,
	}

	task, err = s.taskService.CreateTask(ctx, task)
	if err != nil {
		return &taskpb.CreateTaskResponse{
			Error: convertErrorToProto(err),
		}, nil
	}

	s.logger.Info("task created", zap.Any("task", task))

	return &taskpb.CreateTaskResponse{
		Task: convertTaskToProto(task),
	}, nil
}

func (s *Server) GetTasks(ctx context.Context, req *taskpb.GetTasksRequest) (*taskpb.GetTasksResponse, error) {
	tasks, count, err := s.taskService.GetTasks(ctx, req.GetCustomerId(), int(req.GetLimit()), int(req.GetOffset()))
	if err != nil {
		return &taskpb.GetTasksResponse{
			Error: convertErrorToProto(err),
		}, nil
	}

	s.logger.Info("tasks fetched", zap.Int("count", count))

	return &taskpb.GetTasksResponse{
		Tasks: gospadi.Map(tasks, convertTaskToProto),
		Total: int32(count),
	}, nil
}

func (s *Server) GetTaskByID(ctx context.Context, req *taskpb.GetTaskByIDRequest) (*taskpb.GetTaskByIDResponse, error) {
	if req.GetId() == "" {
		return &taskpb.GetTaskByIDResponse{
			Error: validationError("id is required"),
		}, nil
	}

	task, err := s.taskService.GetTaskByID(ctx, req.GetId())
	if err != nil {
		return &taskpb.GetTaskByIDResponse{
			Error: convertErrorToProto(err),
		}, nil
	}

	s.logger.Info("task fetched", zap.Any("task", task))

	return &taskpb.GetTaskByIDResponse{
		Task: convertTaskToProto(task),
	}, nil
}

func (s *Server) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskRequest) (*taskpb.UpdateTaskResponse, error) {
	if req.GetTask() == nil {
		return &taskpb.UpdateTaskResponse{
			Error: validationError("task is required"),
		}, nil
	}

	payload := req.GetTask()

	if payload.Id == "" {
		return &taskpb.UpdateTaskResponse{
			Error: validationError("id is required"),
		}, nil
	}

	metaJSON, err := convertProtoMetaToJSON(payload.Meta)
	if err != nil {
		return &taskpb.UpdateTaskResponse{
			Error: validationError("meta is invalid"),
		}, nil
	}

	task := &domain.Task{
		ID:               payload.Id,
		CustomerID:       payload.CustomerId,
		Name:             payload.Name,
		Description:      payload.Description,
		VerificationType: convertVerificationTypeToDomain(payload.VerificationType),
		Cost:             int(payload.Cost),
		MembersCount:     int(payload.MembersCount),
		Meta:             metaJSON,
	}

	task, err = s.taskService.UpdateTask(ctx, task)
	if err != nil {
		return &taskpb.UpdateTaskResponse{
			Error: convertErrorToProto(err),
		}, nil
	}

	s.logger.Info("task updated", zap.Any("task", task))

	return &taskpb.UpdateTaskResponse{
		Task: convertTaskToProto(task),
	}, nil
}

func (s *Server) DeleteTask(ctx context.Context, req *taskpb.DeleteTaskRequest) (*taskpb.DeleteTaskResponse, error) {
	if req.GetMaxId() == "" {
		return &taskpb.DeleteTaskResponse{
			Error: validationError("max id is required"),
		}, nil
	}

	err := s.taskService.DeleteTask(ctx, req.GetMaxId())
	if err != nil {
		return &taskpb.DeleteTaskResponse{
			Error: convertErrorToProto(err),
		}, nil
	}

	s.logger.Info("task deleted", zap.String("task_id", req.GetMaxId()))

	return &taskpb.DeleteTaskResponse{
		MaxId: req.GetMaxId(),
	}, nil
}

func convertTaskToProto(task *domain.Task) *taskpb.Task {
	meta := convertTaskMetaToProto(task.Meta)

	return &taskpb.Task{
		Id:               task.ID,
		CustomerId:       task.CustomerID,
		Name:             task.Name,
		Description:      task.Description,
		VerificationType: convertVerificationTypeToProto(task.VerificationType),
		Cost:             int32(task.Cost),
		MembersCount:     int32(task.MembersCount),
		Meta:             meta,
		CreatedAt:        int32(task.CreatedAt.Unix()),
		UpdatedAt:        int32(task.UpdatedAt.Unix()),
	}
}

func convertVerificationTypeToDomain(verificationType taskpb.VerificationType) domain.VerificationType {
	switch verificationType {
	case taskpb.VerificationType_VERIFICATION_TYPE_KYC:
		return domain.VerificationTypeKYC
	case taskpb.VerificationType_VERIFICATION_TYPE_NONE:
		return domain.VerificationTypeNone
	case taskpb.VerificationType_VERIFICATION_TYPE_OTHER:
		return domain.VerificationTypeOther
	default:
		return domain.VerificationType("")
	}
}

func convertVerificationTypeToProto(verificationType domain.VerificationType) taskpb.VerificationType {
	switch verificationType {
	case domain.VerificationTypeKYC:
		return taskpb.VerificationType_VERIFICATION_TYPE_KYC
	case domain.VerificationTypeNone:
		return taskpb.VerificationType_VERIFICATION_TYPE_NONE
	case domain.VerificationTypeOther:
		return taskpb.VerificationType_VERIFICATION_TYPE_OTHER
	default:
		return taskpb.VerificationType_VERIFICATION_TYPE_UNSPECIFIED
	}
}

func convertTaskMetaToProto(raw json.RawMessage) []*taskpb.Meta {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}

	var metaMap map[string]string
	if err := json.Unmarshal(raw, &metaMap); err != nil {
		return nil
	}

	result := make([]*taskpb.Meta, 0, len(metaMap))
	for key, value := range metaMap {
		result = append(result, &taskpb.Meta{
			Key:   key,
			Value: value,
		})
	}

	return result
}

func convertProtoMetaToJSON(meta []*taskpb.Meta) (json.RawMessage, error) {
	if len(meta) == 0 {
		return nil, nil
	}

	metaMap := make(map[string]string, len(meta))
	for _, item := range meta {
		if item == nil {
			continue
		}
		key := item.GetKey()
		if key == "" {
			return nil, errors.New("meta key is empty")
		}
		metaMap[key] = item.GetValue()
	}

	if len(metaMap) == 0 {
		return nil, nil
	}

	bytes, err := json.Marshal(metaMap)
	if err != nil {
		return nil, err
	}

	return json.RawMessage(bytes), nil
}

func validationError(msg string) *taskpb.Error {
	return &taskpb.Error{
		Code:    taskpb.ErrorCode_ERROR_CODE_VALIDATION,
		Message: msg,
	}
}

func convertErrorToProto(err error) *taskpb.Error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, task.ErrTaskNotFound):
		return &taskpb.Error{
			Code:    taskpb.ErrorCode_ERROR_CODE_NOT_FOUND,
			Message: err.Error(),
		}
	case errors.Is(err, task.ErrTaskAlreadyExists):
		return &taskpb.Error{
			Code:    taskpb.ErrorCode_ERROR_CODE_ALREADY_EXISTS,
			Message: err.Error(),
		}
	case errors.Is(err, task.ErrTaskInvalid):
		return &taskpb.Error{
			Code:    taskpb.ErrorCode_ERROR_CODE_VALIDATION,
			Message: err.Error(),
		}
	case errors.Is(err, task.ErrTaskInternal):
		return &taskpb.Error{
			Code:    taskpb.ErrorCode_ERROR_CODE_INTERNAL,
			Message: err.Error(),
		}
	case errors.Is(err, task.ErrUserTaskNotFound):
		return &taskpb.Error{
			Code:    taskpb.ErrorCode_ERROR_CODE_NOT_FOUND,
			Message: err.Error(),
		}
	case errors.Is(err, task.ErrUserTaskAlreadyExists):
		return &taskpb.Error{
			Code:    taskpb.ErrorCode_ERROR_CODE_ALREADY_EXISTS,
			Message: err.Error(),
		}
	case errors.Is(err, task.ErrUserTaskInvalid):
		return &taskpb.Error{
			Code:    taskpb.ErrorCode_ERROR_CODE_VALIDATION,
			Message: err.Error(),
		}
	case errors.Is(err, task.ErrUserTaskInternal):
		return &taskpb.Error{
			Code:    taskpb.ErrorCode_ERROR_CODE_INTERNAL,
			Message: err.Error(),
		}
	default:
		return &taskpb.Error{
			Code:    taskpb.ErrorCode_ERROR_CODE_UNSPECIFIED,
			Message: err.Error(),
		}
	}
}
