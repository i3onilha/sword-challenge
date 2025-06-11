package mysql

import (
	"context"
	"database/sql"
	"sword-challenge/internal/models"
	"sword-challenge/internal/repository"
	"sword-challenge/internal/repository/mysql/tasks"
)

type taskRepository struct {
	query tasks.Queries
}

func NewTaskRepository(db *sql.DB) repository.TaskRepository {
	return &taskRepository{query: *tasks.New(db)}
}

func (r *taskRepository) Create(ctx context.Context, task *models.Task) error {
	return r.query.Create(ctx, tasks.CreateParams{
		TechnicianID: task.TechnicianID,
		Title:        task.Title,
		Summary:      task.Summary,
		PerformedAt:  task.PerformedAt,
	})
}

func (r *taskRepository) GetByID(ctx context.Context, id int64) (*models.Task, error) {
	task, err := r.query.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &models.Task{
		ID:           task.ID,
		TechnicianID: task.TechnicianID,
		Title:        task.Title,
		Summary:      task.Summary,
		PerformedAt:  task.PerformedAt,
		CreatedAt:    task.CreatedAt.Time,
		UpdatedAt:    task.UpdatedAt.Time,
	}, nil
}

func (r *taskRepository) GetByTechnicianID(ctx context.Context, technicianID int64) ([]*models.Task, error) {
	tallTasks, err := r.query.GetByTechnicianID(ctx, technicianID)
	if err != nil {
		return nil, err
	}
	tasks := make([]*models.Task, 0, len(tallTasks))
	for _, task := range tallTasks {
		tasks = append(tasks, &models.Task{
			ID:           task.ID,
			TechnicianID: task.TechnicianID,
			Title:        task.Title,
			Summary:      task.Summary,
		})
	}
	return tasks, nil
}

func (r *taskRepository) GetAll(ctx context.Context) ([]*models.Task, error) {
	tallTasks, err := r.query.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	tasks := make([]*models.Task, 0, len(tallTasks))
	for _, task := range tallTasks {
		tasks = append(tasks, &models.Task{
			ID:           task.ID,
			TechnicianID: task.TechnicianID,
			Title:        task.Title,
			Summary:      task.Summary,
		})
	}
	return tasks, nil
}

func (r *taskRepository) Update(ctx context.Context, task *models.Task) error {
	return r.query.Update(ctx, tasks.UpdateParams{
		ID:          task.ID,
		Title:       task.Title,
		Summary:     task.Summary,
		PerformedAt: task.PerformedAt,
	})
}

func (r *taskRepository) Delete(ctx context.Context, id int64) error {
	return r.query.Delete(ctx, id)
}
