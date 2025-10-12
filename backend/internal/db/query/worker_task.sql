-- name: CreateTask :one
INSERT INTO tasks (
    task_type,
    payload,
    status,
    max_retries,
    scheduled_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: ClaimNextTask :one
UPDATE tasks
SET 
    status = 'processing',
    started_at = NOW(),
    updated_at = NOW()
WHERE id = (
    SELECT id
    FROM tasks
    WHERE status IN ('pending', 'failed')
        AND scheduled_at <= NOW()
        AND retry_count < max_retries
    ORDER BY scheduled_at ASC
    LIMIT 1
    FOR UPDATE SKIP LOCKED
)
RETURNING *;

-- name: MarkTaskCompleted :exec
UPDATE tasks
SET 
    status = 'completed',
    completed_at = NOW(),
    updated_at = NOW()
WHERE id = $1;

-- name: MarkTaskFailed :exec
UPDATE tasks
SET 
    status = 'failed',
    retry_count = retry_count + 1,
    error_message = $2,
    updated_at = NOW(),
    scheduled_at = NOW() + ($3 || ' seconds')::interval
WHERE id = $1;

-- name: GetTaskByID :one
SELECT * FROM tasks WHERE id = $1;

-- name: GetPendingTasksCount :one
SELECT COUNT(*) FROM tasks WHERE status = 'pending';

-- name: GetProcessingTasksCount :one
SELECT COUNT(*) FROM tasks WHERE status = 'processing';

-- name: GetFailedTasksCount :one
SELECT COUNT(*) FROM tasks WHERE status = 'failed' AND retry_count >= max_retries;

-- name: GetTasksByStatus :many
SELECT * FROM tasks 
WHERE status = $1 
ORDER BY created_at DESC 
LIMIT $2 OFFSET $3;

-- name: CleanupOldTasks :exec
DELETE FROM tasks 
WHERE status = 'completed' 
    AND completed_at < NOW() - ($1 || ' days')::interval;

-- name: ResetStaleTasks :exec
UPDATE tasks
SET 
    status = 'pending',
    started_at = NULL,
    updated_at = NOW()
WHERE status = 'processing'
    AND started_at < NOW() - ($1 || ' minutes')::interval;