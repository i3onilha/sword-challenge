package repository

import (
	"context"
	"sword-challenge/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int64) error
}

type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id int64) (*models.Task, error)
	GetByTechnicianID(ctx context.Context, technicianID int64) ([]*models.Task, error)
	GetAll(ctx context.Context) ([]*models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, id int64) error
}

type NotificationRepository interface {
	Create(ctx context.Context, notification *models.Notification) error
	GetUnread(ctx context.Context) ([]*models.Notification, error)
	MarkAsRead(ctx context.Context, id int64) error
	Delete(ctx context.Context, id int64) error
}
