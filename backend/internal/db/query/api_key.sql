
-- name: CreateAPIKey :one
INSERT INTO
    api_keys (
        user_id,
        company_id,
        name,
        api_key
    )
VALUES ($1, $2, $3, $4) RETURNING id,
    user_id,
    company_id,
    name,
    api_key,
    created_at,
    updated_at;

-- name: GetAPIKeysByUserID :many
SELECT
    id,
    user_id,
    company_id,
    name,
    api_key,
    created_at,
    updated_at
FROM api_keys
WHERE
    user_id = $1
    AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: GetAPIKeysByCompanyID :many
SELECT
    id,
    user_id,
    company_id,
    name,
    api_key,
    created_at,
    updated_at
FROM api_keys
WHERE
    company_id = $1
    AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: DeleteAPIKey :exec
UPDATE api_keys
SET
    deleted_at = now(),
    updated_at = now()
WHERE
    id = $1
    AND deleted_at IS NULL;

-- name: FindUserWithAPIKey :one
SELECT
    u.id AS user_id,
    u.email,
    u.fullname,
    ak.company_id,
    ak.id AS api_key_id,
    ak.name AS api_key_name
FROM api_keys ak
    JOIN users u ON ak.user_id = u.id
WHERE
    ak.api_key = $1
    AND ak.deleted_at IS NULL
    AND u.deleted_at IS NULL
LIMIT 1;