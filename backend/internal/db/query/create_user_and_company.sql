-- name: CreateCompany :one
INSERT INTO companies (companyname) VALUES ($1) RETURNING *;

-- name: GetCompanyByID :one
SELECT * FROM companies WHERE id = $1 AND deleted_at IS NULL;

-- name: ListCompanies :many
SELECT *
FROM companies
WHERE
    deleted_at IS NULL
ORDER BY created_at DESC;

-- name: SoftDeleteCompany :exec
UPDATE companies SET deleted_at = now() WHERE id = $1;

-- name: CreateUser :one
INSERT INTO
    users (
        fullname,
        company_id,
        email,
        phonenumber,
        password,
        google_id,
        picture,
        verified,
        blocked,
        verified_at,
        status,
        scheduled_for_deletion,
        scheduled_deletion_at,
        last_login_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11,
        $12,
        $13,
        $14
    ) RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 AND deleted_at IS NULL;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL;

-- name: ListUsersByCompany :many
SELECT *
FROM users
WHERE
    company_id = $1
    AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateUserLoginTime :exec
UPDATE users
SET
    last_login_at = now(),
    updated_at = now()
WHERE
    id = $1;

-- name: VerifyUser :exec
UPDATE users
SET
    verified = TRUE,
    verified_at = now(),
    updated_at = now()
WHERE
    id = $1;

-- name: BlockUser :exec
UPDATE users SET blocked = TRUE, updated_at = now() WHERE id = $1;

-- name: SoftDeleteUser :exec
UPDATE users SET deleted_at = now() WHERE id = $1;