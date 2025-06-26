package service

import (
	"context"
	"errors"
	"sword-challenge/internal/models"
	"sword-challenge/internal/repository"
	"sword-challenge/pkg/messaging"
)

var (
	ErrUnauthorized = errors.New("unauthorized access")
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidInput = errors.New("invalid input")
)

type TaskService struct {
	taskRepo      repository.TaskRepository
	userRepo      repository.UserRepository
	messageBroker messaging.MessageBroker
}

func NewTaskService(
	taskRepo repository.TaskRepository,
	userRepo repository.UserRepository,
	messageBroker messaging.MessageBroker,
) *TaskService {
	return &TaskService{
		taskRepo:      taskRepo,
		userRepo:      userRepo,
		messageBroker: messageBroker,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, task *models.Task, userID int64) error {
	user, err := s.userRepo.GetByID(ctx, userID) // don't trust in user input
	if err != nil {
		return err
	}
	if user == nil {
		return ErrNotFound
	}
	if !user.IsTechnician() {
		return ErrUnauthorized
	}

	// Sanitize input
	task.Sanitize()

	// Validate input
	if err := task.Validate(); err != nil {
		return ErrInvalidInput
	}

	task.TechnicianID = userID
	if err := s.taskRepo.Create(ctx, &models.Task{
		TechnicianID: userID,
		Title:        task.Title,
		Summary:      task.Summary,
		PerformedAt:  task.PerformedAt,
	}); err != nil {
		return err
	}

	task, err = s.taskRepo.GetLastInsertTask(ctx)
	if err != nil {
		return err
	}

	// Publish task created event
	return s.messageBroker.PublishTaskCreated(ctx, task.ID, userID, task.Title)
}

func (s *TaskService) GetTask(ctx context.Context, taskID int64, userID int64) (*models.Task, error) {
	user, err := s.userRepo.GetByID(ctx, userID) // don't trust in user input
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNotFound
	}

	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, ErrNotFound
	}

	// Technicians can only see their own tasks
	if user.IsTechnician() && task.TechnicianID != userID {
		return nil, ErrUnauthorized
	}

	return &models.Task{
		ID:           task.ID,
		TechnicianID: task.TechnicianID,
		Title:        task.Title,
		Summary:      task.Summary,
		PerformedAt:  task.PerformedAt,
	}, nil
}

func (s *TaskService) GetTasks(ctx context.Context, userID int64) ([]*models.Task, error) {
	user, err := s.userRepo.GetByID(ctx, userID) // don't trust in user input
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrNotFound
	}

	if user.IsTechnician() {
		return s.taskRepo.GetByTechnicianID(ctx, userID)
	}
	return s.taskRepo.GetAll(ctx)
}

func (s *TaskService) UpdateTask(ctx context.Context, task *models.Task, userID int64) error {
	user, err := s.userRepo.GetByID(ctx, userID) // don't trust in user input
	if err != nil {
		return err
	}
	if user == nil {
		return ErrNotFound
	}

	existingTask, err := s.taskRepo.GetByID(ctx, task.ID)
	if err != nil {
		return err
	}
	if existingTask == nil {
		return ErrNotFound
	}

	// Only technicians can update their own tasks
	if user.IsTechnician() {
		if existingTask.TechnicianID != userID {
			return ErrUnauthorized
		}
	}

	// Sanitize input
	task.Sanitize()

	// Validate input
	if err := task.Validate(); err != nil {
		return err
	}

	return s.taskRepo.Update(ctx, task)
}

func (s *TaskService) DeleteTask(ctx context.Context, taskID int64, userID int64) error {
	user, err := s.userRepo.GetByID(ctx, userID) // don't trust in user input
	if err != nil {
		return err
	}
	if user == nil {
		return ErrNotFound
	}

	// Only managers can delete tasks
	if !user.IsManager() {
		return ErrUnauthorized
	}

	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}
	if task == nil {
		return ErrNotFound
	}

	return s.taskRepo.Delete(ctx, taskID)
}
