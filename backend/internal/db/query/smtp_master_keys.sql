-- name: CreateSMTPMasterKey :one
INSERT INTO
    smtp_master_keys (
        user_id,
        company_id,
        smtp_login,
        key_name,
        password,
        status
    )
VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetSMTPMasterKeyAndPass :one
SELECT *
FROM smtp_master_keys
WHERE
    key_name = $1
    AND password = $2
    AND deleted_at IS NULL
LIMIT 1;

-- name: CheckSMTPMasterKeyExists :one
SELECT EXISTS (
        SELECT 1
        FROM smtp_master_keys
        WHERE
            smtp_login = $1
            AND password = $2
            AND deleted_at IS NULL
    ) AS exists;

-- name: UpdateSMTPKeyMasterPasswordAndLogin :exec
UPDATE smtp_master_keys
SET password = $1,
smtp_login = $2
WHERE
    user_id = $3;

-- name: GetMasterSMTPKey :one 
SELECT * FROm smtp_master_keys WHERE user_id = $1;