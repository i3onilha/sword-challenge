-- name: Create :exec
INSERT INTO tasks (technician_id, title, summary, performed_at)
VALUES (?, ?, ?, ?);

-- name: GetLastInsertID :one
SELECT * FROM users WHERE id = LAST_INSERT_ID();

-- name: GetByID :one
SELECT * FROM tasks WHERE id = ?;

-- name: GetAll :many
SELECT * FROM tasks;

-- name: GetByTechnicianID :many
SELECT * FROM tasks WHERE technician_id = ?;

-- name: Update :exec
UPDATE tasks SET title = ?, summary = ?, performed_at = ? WHERE id = ?;

-- name: Delete :exec
DELETE FROM tasks WHERE id = ?;