package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"sword-challenge/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockNotificationRepository is a mock implementation of repository.NotificationRepository
type MockNotificationRepository struct {
	mock.Mock
}

func (m *MockNotificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	args := m.Called(ctx, notification)
	return args.Error(0)
}

func (m *MockNotificationRepository) GetUnread(ctx context.Context) ([]*models.Notification, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Notification), args.Error(1)
}

func (m *MockNotificationRepository) MarkAsRead(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockNotificationRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestNotificationService_GetUnreadNotifications(t *testing.T) {
	tests := []struct {
		name           string
		userID         int64
		mockUser       *models.User
		mockUserErr    error
		mockNotifs     []*models.Notification
		mockNotifsErr  error
		expectedNotifs []*models.Notification
		expectedErr    error
	}{
		{
			name:   "success - manager gets notifications",
			userID: 1,
			mockUser: &models.User{
				ID:   1,
				Role: models.RoleManager,
			},
			mockNotifs: []*models.Notification{
				{ID: 1, TaskID: 1, Message: "Test notification", IsRead: false, CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			expectedNotifs: []*models.Notification{
				{ID: 1, TaskID: 1, Message: "Test notification", IsRead: false, CreatedAt: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
		},
		{
			name:        "error - user not found",
			userID:      1,
			mockUser:    nil,
			mockUserErr: nil,
			expectedErr: ErrNotFound,
		},
		{
			name:   "error - user not manager",
			userID: 1,
			mockUser: &models.User{
				ID:   1,
				Role: models.RoleTechnician,
			},
			expectedErr: ErrUnauthorized,
		},
		{
			name:   "error - repository error",
			userID: 1,
			mockUser: &models.User{
				ID:   1,
				Role: models.RoleManager,
			},
			mockNotifsErr: errors.New("repository error"),
			expectedErr:   errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockUserRepo := new(MockUserRepository)
			mockNotifRepo := new(MockNotificationRepository)

			// Setup expectations
			mockUserRepo.On("GetByID", mock.Anything, tt.userID).Return(tt.mockUser, tt.mockUserErr)
			if tt.mockUser != nil && tt.mockUser.Role == models.RoleManager {
				mockNotifRepo.On("GetUnread", mock.Anything).Return(tt.mockNotifs, tt.mockNotifsErr)
			}

			// Create service
			service := NewNotificationService(mockNotifRepo, mockUserRepo)

			// Execute
			notifs, err := service.GetUnreadNotifications(context.Background(), tt.userID)

			// Assert
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
				assert.Nil(t, notifs)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedNotifs, notifs)
			}

			// Verify all expectations were met
			mockUserRepo.AssertExpectations(t)
			mockNotifRepo.AssertExpectations(t)
		})
	}
}

func TestNotificationService_MarkAsRead(t *testing.T) {
	tests := []struct {
		name           string
		notificationID int64
		userID         int64
		mockUser       *models.User
		mockUserErr    error
		mockNotifErr   error
		expectedErr    error
	}{
		{
			name:           "success - manager marks notification as read",
			notificationID: 1,
			userID:         1,
			mockUser: &models.User{
				ID:   1,
				Role: models.RoleManager,
			},
		},
		{
			name:           "error - user not found",
			notificationID: 1,
			userID:         1,
			mockUser:       nil,
			mockUserErr:    nil,
			expectedErr:    ErrNotFound,
		},
		{
			name:           "error - user not manager",
			notificationID: 1,
			userID:         1,
			mockUser: &models.User{
				ID:   1,
				Role: models.RoleTechnician,
			},
			expectedErr: ErrUnauthorized,
		},
		{
			name:           "error - repository error",
			notificationID: 1,
			userID:         1,
			mockUser: &models.User{
				ID:   1,
				Role: models.RoleManager,
			},
			mockNotifErr: errors.New("repository error"),
			expectedErr:  errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockUserRepo := new(MockUserRepository)
			mockNotifRepo := new(MockNotificationRepository)

			// Setup expectations
			mockUserRepo.On("GetByID", mock.Anything, tt.userID).Return(tt.mockUser, tt.mockUserErr)
			if tt.mockUser != nil && tt.mockUser.Role == models.RoleManager {
				mockNotifRepo.On("MarkAsRead", mock.Anything, tt.notificationID).Return(tt.mockNotifErr)
			}

			// Create service
			service := NewNotificationService(mockNotifRepo, mockUserRepo)

			// Execute
			err := service.MarkAsRead(context.Background(), tt.notificationID, tt.userID)

			// Assert
			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedErr, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify all expectations were met
			mockUserRepo.AssertExpectations(t)
			mockNotifRepo.AssertExpectations(t)
		})
	}
}
