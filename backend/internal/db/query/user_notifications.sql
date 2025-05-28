-- name: CreateUserNotification :one
INSERT INTO
    user_notifications (
        user_id,
        title,
        additional_field
    )
VALUES ($1, $2, $3) RETURNING *;

-- name: GetUserNotifications :many
SELECT *
FROM user_notifications
WHERE
    user_id = $1
ORDER BY created_at DESC;

-- name: MarkNotificationAsRead :exec
UPDATE user_notifications
SET
    read_status = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    user_id = $1
    AND id = $1;

-- name: MarkAllUserNotificationsAsRead :exec
UPDATE user_notifications
SET
    read_status = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    user_id = $1
    AND read_status = FALSE;