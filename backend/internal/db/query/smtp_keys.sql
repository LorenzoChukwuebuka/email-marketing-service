-- name: CreateSMTPKey :one
INSERT INTO
    smtp_keys (
        company_id,
        user_id,
        key_name,
        password,
        status,
        smtp_login
    )
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetSMTPKeyUserAndPass :one
SELECT *
FROM smtp_keys
WHERE
    key_name = $1
    AND password = $2
    AND deleted_at IS NULL
LIMIT 1;

-- name: CheckSMTPKeyExists :one
SELECT EXISTS (
        SELECT 1
        FROM smtp_keys
        WHERE
            key_name = $1
            AND password = $2
            AND deleted_at IS NULL
    ) AS exists;

-- name: GetUserSMTPKey :many
SELECT * FROM smtp_keys WHERE user_id = $1;

-- name: UpdateSMTPKeyLogin :exec
UPDATE smtp_keys SET smtp_login = $1 WHERE user_id = $2;

-- name: GetUserSmtpKeys :many
SELECT * FROM smtp_keys WHERE user_id = $1;

-- name: GetSMTPKeyByID :one
SELECT * FROM smtp_keys WHERE user_id = $1 AND id = $2 LIMIT 1;

-- name: UpdateSMTPKeyStatus :exec
UPDATE smtp_keys
SET
    status = $1,
    updated_at = now()
WHERE
    id = $2
    AND user_id = $3;

-- name: SoftDeleteSMTPKey :exec
UPDATE smtp_keys
SET
    updated_at = now(),
    deleted_at = now()
WHERE
    id = $1;