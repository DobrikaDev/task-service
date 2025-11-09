package sql

import "errors"

var (
	ErrTaskNotFound      = errors.New("task not found")
	ErrTaskAlreadyExists = errors.New("task already exists")
	ErrTaskInvalid       = errors.New("task invalid")
	ErrTaskInternal      = errors.New("task internal error")

	ErrUserTaskNotFound      = errors.New("user task not found")
	ErrUserTaskAlreadyExists = errors.New("user task already exists")
	ErrUserTaskInvalid       = errors.New("user task invalid")
	ErrUserTaskInternal      = errors.New("user task internal error")

	ErrFeedbackNotFound      = errors.New("feedback not found")
	ErrFeedbackInternal      = errors.New("feedback internal error")
	ErrFeedbackInvalid       = errors.New("feedback invalid")
	ErrFeedbackAlreadyExists = errors.New("feedback already exists")
)
