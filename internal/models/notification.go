package models

import (
	"errors"
	"time"
)

// Notification represents a notification in the system
// @Description Notification information
type Notification struct {
	// @Description The unique identifier of the notification
	ID int64 `json:"id" example:"1"`
	// @Description The ID of the task this notification is about
	TaskID int64 `json:"task_id" example:"1"`
	// @Description The notification message
	Message string `json:"message" example:"The tech John Doe performed the task on 2024-03-20 14:30:00"`
	// @Description Whether the notification has been read
	IsRead bool `json:"is_read" example:"false"`
	// @Description When the notification was created
	CreatedAt time.Time `json:"created_at" example:"2024-03-20T14:30:00Z"`
}

var (
	ErrNilTask       = errors.New("task cannot be nil")
	ErrNilTechnician = errors.New("technician cannot be nil")
)

func NewTaskNotification(task *Task, technician *User) (*Notification, error) {
	if task == nil {
		return nil, ErrNilTask
	}
	if technician == nil {
		return nil, ErrNilTechnician
	}

	message, err := formatNotificationMessage(technician, task)
	if err != nil {
		return nil, err
	}

	return &Notification{
		TaskID:    task.ID,
		Message:   message,
		CreatedAt: time.Now(),
	}, nil
}

func formatNotificationMessage(tech *User, task *Task) (string, error) {
	if tech == nil {
		return "", ErrNilTechnician
	}
	if task == nil {
		return "", ErrNilTask
	}
	return "The tech " + tech.Name + " performed the task on " + task.PerformedAt.Format("2006-01-02 15:04:05"), nil
}
