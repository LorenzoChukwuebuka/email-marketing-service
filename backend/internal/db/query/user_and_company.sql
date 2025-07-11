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
        verified,
        verified_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8
    ) RETURNING *;

-- name: GetUserByEmail :one
SELECT
    u.*,
    c.id AS company_id,
    c.companyname,
    c.created_at AS company_created_at,
    c.updated_at AS company_updated_at,
    c.deleted_at AS company_deleted_at
FROM users u
    JOIN companies c ON u.company_id = c.id
WHERE
    u.email = $1
    AND u.deleted_at IS NULL
    AND c.deleted_at IS NULL;

-- name: GetUserByID :one
SELECT
    u.*,
    c.id AS company_id,
    c.companyname,
    c.created_at AS company_created_at,
    c.updated_at AS company_updated_at,
    c.deleted_at AS company_deleted_at
FROM users u
    JOIN companies c ON u.company_id = c.id
WHERE
    u.id = $1
    AND u.deleted_at IS NULL
    AND c.deleted_at IS NULL;

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

-- name: ResetUserPassword :exec
UPDATE users SET password = $1, updated_at = now() WHERE id = $2;

-- name: MarkUserForDeletion :one
UPDATE users
SET
    scheduled_for_deletion = TRUE,
    scheduled_deletion_at = $2,
    status = $3,
    updated_at = now()
WHERE
    id = $1 RETURNING *;

-- name: CancelUserDeletion :one
UPDATE users
SET
    scheduled_for_deletion = FALSE,
    scheduled_deletion_at = NULL,
    updated_at = now()
WHERE
    id = $1 RETURNING *;

-- name: DeleteScheduledUsers :many
UPDATE users
SET
    deleted_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE
    scheduled_for_deletion = TRUE
    AND deleted_at IS NULL
    AND scheduled_deletion_at IS NOT NULL
    AND scheduled_deletion_at < CURRENT_TIMESTAMP - INTERVAL '30 days' RETURNING *;

-- name: UpdateUserRecords :exec
UPDATE users
SET
    fullname = COALESCE(sqlc.narg('fullname'), fullname),
    email = COALESCE(sqlc.narg('email'), email),
    phonenumber = COALESCE(sqlc.narg('phonenumber'), phonenumber),
    updated_at = now()
WHERE
    id = $1;

-- name: UpdateCompanyName :exec
UPDATE companies
SET 
    companyname = COALESCE($2, companyname),
    updated_at = now()
WHERE id = $1;