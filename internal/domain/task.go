package domain

import (
	"encoding/json"
	"time"
)

type Task struct {
	ID               string           `json:"id" db:"id"`
	CustomerID       string           `json:"customer_id" db:"customer_id"`
	Name             string           `json:"name" db:"name"`
	Description      string           `json:"description" db:"description"`
	VerificationType VerificationType `json:"verification_type" db:"verification_type"`
	Cost             int              `json:"cost" db:"cost"`
	MembersCount     int              `json:"members_count" db:"members_count"`
	Meta             json.RawMessage  `json:"meta" db:"meta"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserTask struct {
	UserID string `json:"user_id" db:"user_id"`
	TaskID string `json:"task_id" db:"task_id"`
	Status Status `json:"status" db:"status"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Status string

const (
	StatusInProgress Status = "pending"
	StatusCompleted  Status = "completed"
	StatusApproved   Status = "approved"
	StatusRejected   Status = "rejected"
	StatusCancelled  Status = "cancelled"
)

func (s Status) String() string {
	return string(s)
}
