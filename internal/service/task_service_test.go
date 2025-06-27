package service

import (
	"context"
	"testing"
	"time"

	"sword-challenge/internal/models"
	"sword-challenge/pkg/messaging"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTaskRepository struct {
	mock.Mock
}

func (m *MockTaskRepository) Create(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) GetByID(ctx context.Context, id int64) (*models.Task, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepository) GetByTechnicianID(ctx context.Context, technicianID int64) ([]*models.Task, error) {
	args := m.Called(ctx, technicianID)
	return args.Get(0).([]*models.Task), args.Error(1)
}

func (m *MockTaskRepository) GetLastInsertTask(ctx context.Context) (*models.Task, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Task), args.Error(1)
}

func (m *MockTaskRepository) GetAll(ctx context.Context) ([]*models.Task, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(ctx context.Context, task *models.Task) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTaskService_CreateTask(t *testing.T) {
	tests := []struct {
		name           string
		task           *models.Task
		userID         int64
		setupMocks     func(*MockTaskRepository, *MockUserRepository, *messaging.MockBroker)
		expectedError  error
		verifyMessages func(*testing.T, *messaging.MockBroker)
	}{
		{
			name: "successful task creation by technician",
			task: &models.Task{
				Title:       "Test task",
				Summary:     "Test task summary",
				PerformedAt: time.Now(),
			},
			userID: 1,
			setupMocks: func(tr *MockTaskRepository, ur *MockUserRepository, mb *messaging.MockBroker) {
				ur.On("GetByID", mock.Anything, int64(1)).Return(&models.User{
					ID:   1,
					Name: "Test Tech",
					Role: models.RoleTechnician,
				}, nil)
				tr.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
			verifyMessages: func(t *testing.T, mb *messaging.MockBroker) {
				messages := mb.GetMessages()
				assert.Len(t, messages, 1)
				assert.Equal(t, int64(0), messages[0].TaskID)
				assert.Equal(t, int64(1), messages[0].TechnicianID)
				assert.Equal(t, "Test task", messages[0].Title)
			},
		},
		{
			name: "unauthorized - manager cannot create task",
			task: &models.Task{
				Title:       "Test task",
				Summary:     "Test task summary",
				PerformedAt: time.Now(),
			},
			userID: 1,
			setupMocks: func(tr *MockTaskRepository, ur *MockUserRepository, mb *messaging.MockBroker) {
				ur.On("GetByID", mock.Anything, int64(1)).Return(&models.User{
					ID:   1,
					Name: "Test Manager",
					Role: models.RoleManager,
				}, nil)
			},
			expectedError: ErrUnauthorized,
			verifyMessages: func(t *testing.T, mb *messaging.MockBroker) {
				messages := mb.GetMessages()
				assert.Len(t, messages, 0)
			},
		},
		{
			name: "user not found",
			task: &models.Task{
				Title:       "Test task",
				Summary:     "Test task summary",
				PerformedAt: time.Now(),
			},
			userID: 1,
			setupMocks: func(tr *MockTaskRepository, ur *MockUserRepository, mb *messaging.MockBroker) {
				ur.On("GetByID", mock.Anything, int64(1)).Return(nil, nil)
			},
			expectedError: ErrNotFound,
			verifyMessages: func(t *testing.T, mb *messaging.MockBroker) {
				messages := mb.GetMessages()
				assert.Len(t, messages, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTaskRepo := new(MockTaskRepository)
			mockUserRepo := new(MockUserRepository)
			mockBroker := messaging.NewMockBroker()

			tt.setupMocks(mockTaskRepo, mockUserRepo, mockBroker)

			service := NewTaskService(mockTaskRepo, mockUserRepo, mockBroker)
			_, err := service.CreateTask(context.Background(), tt.task, tt.userID)

			assert.Equal(t, tt.expectedError, err)
			mockTaskRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
			tt.verifyMessages(t, mockBroker)
		})
	}
}

func TestTaskService_GetTask(t *testing.T) {
	tests := []struct {
		name          string
		taskID        int64
		userID        int64
		setupMocks    func(*MockTaskRepository, *MockUserRepository, *messaging.MockBroker)
		expectedTask  *models.Task
		expectedError error
	}{
		{
			name:   "successful task retrieval by manager",
			taskID: 1,
			userID: 1,
			setupMocks: func(tr *MockTaskRepository, ur *MockUserRepository, mb *messaging.MockBroker) {
				ur.On("GetByID", mock.Anything, int64(1)).Return(&models.User{
					ID:   1,
					Name: "Test Manager",
					Role: models.RoleManager,
				}, nil)
				tr.On("GetByID", mock.Anything, int64(1)).Return(&models.Task{
					ID:           1,
					TechnicianID: 2,
					Title:        "Test task",
					Summary:      "Test task summary",
					PerformedAt:  time.Now(),
				}, nil)
			},
			expectedTask: &models.Task{
				ID:           1,
				TechnicianID: 2,
				Title:        "Test task",
				Summary:      "Test task summary",
				PerformedAt:  time.Now(),
			},
			expectedError: nil,
		},
		{
			name:   "successful task retrieval by technician - own task",
			taskID: 1,
			userID: 1,
			setupMocks: func(tr *MockTaskRepository, ur *MockUserRepository, mb *messaging.MockBroker) {
				ur.On("GetByID", mock.Anything, int64(1)).Return(&models.User{
					ID:   1,
					Name: "Test Tech",
					Role: models.RoleTechnician,
				}, nil)
				tr.On("GetByID", mock.Anything, int64(1)).Return(&models.Task{
					ID:           1,
					TechnicianID: 1,
					Title:        "Test task",
					Summary:      "Test task summary",
					PerformedAt:  time.Now(),
				}, nil)
			},
			expectedTask: &models.Task{
				ID:           1,
				TechnicianID: 1,
				Title:        "Test task",
				Summary:      "Test task summary",
				PerformedAt:  time.Now(),
			},
			expectedError: nil,
		},
		{
			name:   "unauthorized - technician accessing other's task",
			taskID: 1,
			userID: 1,
			setupMocks: func(tr *MockTaskRepository, ur *MockUserRepository, mb *messaging.MockBroker) {
				ur.On("GetByID", mock.Anything, int64(1)).Return(&models.User{
					ID:   1,
					Name: "Test Tech",
					Role: models.RoleTechnician,
				}, nil)
				tr.On("GetByID", mock.Anything, int64(1)).Return(&models.Task{
					ID:           1,
					TechnicianID: 2,
					Title:        "Test task",
					Summary:      "Test task summary",
					PerformedAt:  time.Now(),
				}, nil)
			},
			expectedTask:  nil,
			expectedError: ErrUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTaskRepo := new(MockTaskRepository)
			mockUserRepo := new(MockUserRepository)
			mockBroker := messaging.NewMockBroker()

			tt.setupMocks(mockTaskRepo, mockUserRepo, mockBroker)

			service := NewTaskService(mockTaskRepo, mockUserRepo, mockBroker)
			task, err := service.GetTask(context.Background(), tt.taskID, tt.userID)

			assert.Equal(t, tt.expectedError, err)
			if tt.expectedTask != nil {
				assert.Equal(t, tt.expectedTask.ID, task.ID)
				assert.Equal(t, tt.expectedTask.TechnicianID, task.TechnicianID)
				assert.Equal(t, tt.expectedTask.Title, task.Title)
				assert.Equal(t, tt.expectedTask.Summary, task.Summary)
			} else {
				assert.Nil(t, task)
			}
			mockTaskRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)
		})
	}
}
