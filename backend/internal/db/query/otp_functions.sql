-- name: CreateOTP :one
INSERT INTO
    otps (
        user_id,
        token,
        created_at,
        updated_at,
        expires_at
    )
VALUES (
        $1,
        $2,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP,
        $3
    ) RETURNING *;

-- name: GetOTPByToken :one
SELECT * FROM otps WHERE token = $1 AND deleted_at IS NULL LIMIT 1;

-- name: DeleteOTPById :exec
UPDATE otps SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: HardDeleteOTPById :exec
DELETE FROM otps WHERE id = $1;