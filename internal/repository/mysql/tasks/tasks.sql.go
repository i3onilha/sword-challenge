// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: tasks.sql

package tasks

import (
	"context"
	"time"
)

const create = `-- name: Create :exec
INSERT INTO tasks (technician_id, title, summary, performed_at)
VALUES (?, ?, ?, ?)
`

type CreateParams struct {
	TechnicianID int64
	Title        string
	Summary      string
	PerformedAt  time.Time
}

func (q *Queries) Create(ctx context.Context, arg CreateParams) error {
	_, err := q.db.ExecContext(ctx, create,
		arg.TechnicianID,
		arg.Title,
		arg.Summary,
		arg.PerformedAt,
	)
	return err
}

const delete = `-- name: Delete :exec
DELETE FROM tasks WHERE id = ?
`

func (q *Queries) Delete(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, delete, id)
	return err
}

const getAll = `-- name: GetAll :many
SELECT id, technician_id, title, summary, performed_at, created_at, updated_at FROM tasks
`

func (q *Queries) GetAll(ctx context.Context) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, getAll)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.TechnicianID,
			&i.Title,
			&i.Summary,
			&i.PerformedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getByID = `-- name: GetByID :one
SELECT id, technician_id, title, summary, performed_at, created_at, updated_at FROM tasks WHERE id = ?
`

func (q *Queries) GetByID(ctx context.Context, id int64) (Task, error) {
	row := q.db.QueryRowContext(ctx, getByID, id)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.TechnicianID,
		&i.Title,
		&i.Summary,
		&i.PerformedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getByTechnicianID = `-- name: GetByTechnicianID :many
SELECT id, technician_id, title, summary, performed_at, created_at, updated_at FROM tasks WHERE technician_id = ?
`

func (q *Queries) GetByTechnicianID(ctx context.Context, technicianID int64) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, getByTechnicianID, technicianID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.TechnicianID,
			&i.Title,
			&i.Summary,
			&i.PerformedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLastInsertTask = `-- name: GetLastInsertTask :one
SELECT id, technician_id, title, summary, performed_at, created_at, updated_at FROM tasks WHERE id = LAST_INSERT_ID()
`

func (q *Queries) GetLastInsertTask(ctx context.Context) (Task, error) {
	row := q.db.QueryRowContext(ctx, getLastInsertTask)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.TechnicianID,
		&i.Title,
		&i.Summary,
		&i.PerformedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getLastInsertUser = `-- name: GetLastInsertUser :one
SELECT id, name, email, password_hash, role, created_at, updated_at FROM users WHERE id = LAST_INSERT_ID()
`

func (q *Queries) GetLastInsertUser(ctx context.Context) (User, error) {
	row := q.db.QueryRowContext(ctx, getLastInsertUser)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.PasswordHash,
		&i.Role,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const update = `-- name: Update :exec
UPDATE tasks SET title = ?, summary = ?, performed_at = ? WHERE id = ?
`

type UpdateParams struct {
	Title       string
	Summary     string
	PerformedAt time.Time
	ID          int64
}

func (q *Queries) Update(ctx context.Context, arg UpdateParams) error {
	_, err := q.db.ExecContext(ctx, update,
		arg.Title,
		arg.Summary,
		arg.PerformedAt,
		arg.ID,
	)
	return err
}
