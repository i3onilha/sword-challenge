package mysql

import (
	"context"
	"database/sql"
	"sword-challenge/internal/models"
	"sword-challenge/internal/repository"
	"sword-challenge/internal/repository/mysql/notifications"
)

type notificationRepository struct {
	query notifications.Queries
}

func NewNotificationRepository(db *sql.DB) repository.NotificationRepository {
	return &notificationRepository{query: *notifications.New(db)}
}

func (r *notificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	return r.query.Create(ctx, notifications.CreateParams{
		TaskID:  notification.TaskID,
		Message: notification.Message,
	})
}

func (r *notificationRepository) GetUnread(ctx context.Context) ([]*models.Notification, error) {
	allNotifications, err := r.query.GetUnread(ctx)
	if err != nil {
		return nil, err
	}
	notifications := make([]*models.Notification, 0, len(allNotifications))
	for _, notification := range allNotifications {
		notifications = append(notifications, &models.Notification{
			ID:      notification.ID,
			TaskID:  notification.TaskID,
			Message: notification.Message,
		})
	}
	return notifications, nil
}

func (r *notificationRepository) MarkAsRead(ctx context.Context, id int64) error {
	return r.query.MarkAsRead(ctx, id)
}

func (r *notificationRepository) Delete(ctx context.Context, id int64) error {
	return r.query.Delete(ctx, id)
}
