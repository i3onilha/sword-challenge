package service

import (
	"context"
	"sword-challenge/internal/models"
	"sword-challenge/internal/repository"
)

type NotificationService struct {
	notificationRepo repository.NotificationRepository
	userRepo         repository.UserRepository
}

func NewNotificationService(
	notificationRepo repository.NotificationRepository,
	userRepo repository.UserRepository,
) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
	}
}

func (s *NotificationService) GetUnreadNotifications(ctx context.Context, userID int64) ([]*models.Notification, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNotFound
	}

	// Only managers can view notifications
	if !user.IsManager() {
		return nil, ErrUnauthorized
	}

	return s.notificationRepo.GetUnread(ctx)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID int64, userID int64) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrNotFound
	}

	// Only managers can mark notifications as read
	if !user.IsManager() {
		return ErrUnauthorized
	}

	return s.notificationRepo.MarkAsRead(ctx, notificationID)
}
