-- name: CreateUserNotification :one
INSERT INTO user_notifications (
    user_id,
    title,
    additional_field
)
VALUES ($1, $2, $3) RETURNING *;

-- name: MarkNotificationAsRead :exec
UPDATE user_notifications
SET
    read_status = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    user_id = $1
    AND id = $2;

-- name: MarkAllUserNotificationsAsRead :exec
UPDATE user_notifications
SET
    read_status = TRUE,
    updated_at = CURRENT_TIMESTAMP
WHERE
    user_id = $1
    AND read_status = FALSE;

-- name: GetUserNotifications :many
SELECT *
FROM user_notifications
WHERE
    user_id = $1
ORDER BY created_at DESC
LIMIT 50;

-- name: GetUnreadNotifications :many
SELECT *
FROM user_notifications
WHERE
    user_id = $1
    AND read_status = false
ORDER BY created_at DESC;

-- name: GetUserNotificationsSinceID :many
WITH reference_notification AS (
  SELECT created_at 
  FROM user_notifications
  WHERE user_notifications.id = $2
)
SELECT n.* 
FROM user_notifications n, reference_notification rn
WHERE n.user_id = $1 
  AND n.created_at > rn.created_at
ORDER BY n.created_at DESC;

-- name: GetNotificationCount :one
SELECT COUNT(*)
FROM user_notifications
WHERE
    user_id = $1
    AND read_status = false;
