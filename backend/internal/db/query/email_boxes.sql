-- name: CreateEmailBox :one
INSERT INTO
    email_boxes (
        user_name,
        "from",
        "to",
        content,
        mailbox
    )
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetEmailBoxByID :one
SELECT *
FROM email_boxes
WHERE
    id = $1
    AND deleted_at IS NULL
LIMIT 1;

-- name: ListEmailBoxesByUser :many
SELECT *
FROM email_boxes
WHERE
    user_name = $1
    AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListEmailBoxesByMailbox :many
SELECT *
FROM email_boxes
WHERE
    user_name = $1
    AND mailbox = $2
    AND deleted_at IS NULL
ORDER BY created_at DESC;