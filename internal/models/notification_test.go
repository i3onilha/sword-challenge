package models

import (
	"errors"
	"testing"
	"time"
)

func TestNewTaskNotification(t *testing.T) {
	tests := []struct {
		name       string
		task       *Task
		technician *User
		wantErr    error
	}{
		{
			name: "valid notification",
			task: &Task{
				ID:          1,
				PerformedAt: time.Now(),
			},
			technician: &User{
				ID:   1,
				Name: "John Doe",
			},
			wantErr: nil,
		},
		{
			name:       "nil task",
			task:       nil,
			technician: &User{Name: "John Doe"},
			wantErr:    ErrNilTask,
		},
		{
			name: "nil technician",
			task: &Task{
				ID:          1,
				PerformedAt: time.Now(),
			},
			technician: nil,
			wantErr:    ErrNilTechnician,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTaskNotification(tt.task, tt.technician)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewTaskNotification() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil {
				if got.TaskID != tt.task.ID {
					t.Errorf("NewTaskNotification().TaskID = %v, want %v", got.TaskID, tt.task.ID)
				}
				if got.IsRead {
					t.Error("NewTaskNotification().IsRead = true, want false")
				}
				if got.CreatedAt.IsZero() {
					t.Error("NewTaskNotification().CreatedAt is zero")
				}
			} else if got != nil {
				t.Error("NewTaskNotification() returned non-nil notification with error")
			}
		})
	}
}

func TestFormatNotificationMessage(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		tech    *User
		task    *Task
		want    string
		wantErr error
	}{
		{
			name: "valid message",
			tech: &User{
				Name: "John Doe",
			},
			task: &Task{
				PerformedAt: now,
			},
			want:    "The tech John Doe performed the task on " + now.Format("2006-01-02 15:04:05"),
			wantErr: nil,
		},
		{
			name:    "nil technician",
			tech:    nil,
			task:    &Task{PerformedAt: now},
			wantErr: ErrNilTechnician,
		},
		{
			name:    "nil task",
			tech:    &User{Name: "John Doe"},
			task:    nil,
			wantErr: ErrNilTask,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatNotificationMessage(tt.tech, tt.task)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("formatNotificationMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil && got != tt.want {
				t.Errorf("formatNotificationMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
