package task

import "errors"

var ErrTaskNotFound = errors.New("task not found")
var ErrTaskAlreadyExists = errors.New("task already exists")
var ErrTaskInternal = errors.New("task internal error")
var ErrTaskInvalid = errors.New("task invalid")

var ErrFeedbackNotFound = errors.New("feedback not found")
var ErrFeedbackInternal = errors.New("feedback internal error")
var ErrFeedbackInvalid = errors.New("feedback invalid")
var ErrFeedbackAlreadyExists = errors.New("feedback already exists")
