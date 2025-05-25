package models

import (
	"errors"
	"html"
	"strings"
	"time"
)

var (
	ErrSummaryTooLong   = errors.New("summary exceeds maximum length of 2500 characters")
	ErrTitleTooLong     = errors.New("title exceeds maximum length of 255 characters")
	ErrInvalidDateRange = errors.New("performed_at date must be between 1900-01-01 and 2100-12-31")
	ErrEmptyTitle       = errors.New("title cannot be empty")
	ErrEmptySummary     = errors.New("summary cannot be empty")
)

// Task represents a task in the system
// @Description Task information
type Task struct {
	// @Description The unique identifier of the task
	ID int64 `json:"id" example:"1"`
	// @Description The ID of the technician who performed the task
	TechnicianID int64 `json:"technician_id" example:"1"`
	// @Description The title of the task
	Title string `json:"title" example:"Fix air conditioning"`
	// @Description The detailed summary of the task
	Summary string `json:"summary" example:"Replaced filters and recharged coolant"`
	// @Description When the task was performed
	PerformedAt time.Time `json:"performed_at" example:"2024-03-20T14:30:00Z"`
	// @Description When the task was created
	CreatedAt time.Time `json:"created_at" example:"2024-03-20T14:30:00Z"`
	// @Description When the task was last updated
	UpdatedAt time.Time `json:"updated_at" example:"2024-03-20T14:30:00Z"`
}

func (t *Task) Validate() error {
	// Trim whitespace
	t.Title = strings.TrimSpace(t.Title)
	t.Summary = strings.TrimSpace(t.Summary)

	// Check for empty fields
	if t.Title == "" {
		return ErrEmptyTitle
	}
	if t.Summary == "" {
		return ErrEmptySummary
	}

	// Check length constraints
	if len(t.Title) > 255 {
		return ErrTitleTooLong
	}
	if len(t.Summary) > 2500 {
		return ErrSummaryTooLong
	}

	// Validate date range
	minDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	maxDate := time.Date(2100, 12, 31, 23, 59, 59, 0, time.UTC)
	if t.PerformedAt.Before(minDate) || t.PerformedAt.After(maxDate) {
		return ErrInvalidDateRange
	}

	return nil
}

func (t *Task) Sanitize() {
	// Sanitize HTML/script content
	t.Title = html.EscapeString(t.Title)
	t.Summary = html.EscapeString(t.Summary)
}

// CreateTaskRequest represents the request body for creating a task
// @Description Request body for creating a new task
type CreateTaskRequest struct {
	// @Description The title of the task (max 255 characters)
	Title string `json:"title" binding:"required,max=255" example:"Fix air conditioning"`
	// @Description The detailed summary of the task (max 2500 characters)
	Summary string `json:"summary" binding:"required,max=2500" example:"Replaced filters and recharged coolant"`
	// @Description When the task was performed (ISO 8601 format)
	PerformedAt string `json:"performed_at" binding:"required" example:"2024-03-20T14:30:00Z"`
}

// UpdateTaskRequest represents the request body for updating a task
// @Description Request body for updating an existing task
type UpdateTaskRequest struct {
	// @Description The title of the task (max 255 characters)
	Title string `json:"title" binding:"required,max=255" example:"Fix air conditioning"`
	// @Description The detailed summary of the task (max 2500 characters)
	Summary string `json:"summary" binding:"required,max=2500" example:"Replaced filters and recharged coolant"`
	// @Description When the task was performed (ISO 8601 format)
	PerformedAt string `json:"performed_at" binding:"required" example:"2024-03-20T14:30:00Z"`
}
