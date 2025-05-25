-- name: Create :exec
INSERT INTO notifications (task_id, message)
VALUES (?, ?);

-- name: GetAll :many
SELECT * FROM notifications;

-- name: GetByID :one
SELECT * FROM notifications WHERE id = ?;

-- name: GetByTaskID :many
SELECT * FROM notifications WHERE task_id = ?;

-- name: GetUnread :many
SELECT * FROM notifications WHERE is_read = 0;

-- name: MarkAsRead :exec
UPDATE notifications SET is_read = 1 WHERE id = ?;

-- name: Delete :exec
DELETE FROM notifications WHERE id = ?;

-- name: DeleteByTaskID :exec
DELETE FROM notifications WHERE task_id = ?;