package task

import "errors"

var ErrTaskNotFound = errors.New("task not found")
var ErrTaskAlreadyExists = errors.New("task already exists")
var ErrTaskInternal = errors.New("task internal error")
var ErrTaskInvalid = errors.New("task invalid")

var ErrUserTaskAlreadyExists = errors.New("user task already exists")
var ErrUserTaskInternal = errors.New("user task internal error")
var ErrUserTaskInvalid = errors.New("user task invalid")
var ErrUserTaskNotFound = errors.New("user task not found")