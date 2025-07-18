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
    fullname = COALESCE(
        sqlc.narg ('fullname'),
        fullname
    ),
    email = COALESCE(sqlc.narg ('email'), email),
    phonenumber = COALESCE(
        sqlc.narg ('phonenumber'),
        phonenumber
    ),
    updated_at = now()
WHERE
    id = $1;

-- name: UpdateCompanyName :exec
UPDATE companies
SET
    companyname = COALESCE($2, companyname),
    updated_at = now()
WHERE
    id = $1;

-- name: GetAllUsers :many
SELECT 
    u.id,
    u.fullname,
    u.email,
    u.phonenumber,
    u.picture,
    u.verified,
    u.blocked,
    u.verified_at,
    u.status,
    u.scheduled_for_deletion,
    u.scheduled_deletion_at,
    u.last_login_at,
    u.created_at,
    u.updated_at,
    c.id as company_id,
    c.companyname
FROM users u
LEFT JOIN companies c ON u.company_id = c.id
WHERE u.deleted_at IS NULL
    AND c.deleted_at IS NULL
    AND ($1::text = '' OR u.fullname ILIKE '%' || $1 || '%' OR u.email ILIKE '%' || $1 || '%' OR c.companyname ILIKE '%' || $1 || '%')
ORDER BY u.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountAllUsers :one
SELECT COUNT(*) FROM users u WHERE u.deleted_at IS NULL;

-- name: UnblockUser :one
UPDATE users
SET
    blocked = false,
    updated_at = now()
WHERE
    id = $1
    AND deleted_at IS NULL RETURNING id,
    fullname,
    email,
    phonenumber,
    picture,
    verified,
    blocked,
    verified_at,
    status,
    scheduled_for_deletion,
    scheduled_deletion_at,
    last_login_at,
    created_at,
    updated_at,
    company_id;

-- name: GetVerifiedUsers :many
SELECT 
    u.id,
    u.fullname,
    u.email,
    u.phonenumber,
    u.picture,
    u.verified,
    u.blocked,
    u.verified_at,
    u.status,
    u.scheduled_for_deletion,
    u.scheduled_deletion_at,
    u.last_login_at,
    u.created_at,
    u.updated_at,
    c.id as company_id,
    c.companyname
FROM users u
LEFT JOIN companies c ON u.company_id = c.id
WHERE u.deleted_at IS NULL
    AND c.deleted_at IS NULL
    AND u.verified = true
    AND ($1::text = '' OR u.fullname ILIKE '%' || $1 || '%' OR u.email ILIKE '%' || $1 || '%' OR c.companyname ILIKE '%' || $1 || '%')
ORDER BY u.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountVerifiedUsers :one
SELECT COUNT(*)
FROM users u
WHERE
    u.deleted_at IS NULL
    AND u.verified = true;

-- name: GetUnVerifiedUsers :many
SELECT 
    u.id,
    u.fullname,
    u.email,
    u.phonenumber,
    u.picture,
    u.verified,
    u.blocked,
    u.verified_at,
    u.status,
    u.scheduled_for_deletion,
    u.scheduled_deletion_at,
    u.last_login_at,
    u.created_at,
    u.updated_at,
    c.id as company_id,
    c.companyname
FROM users u
LEFT JOIN companies c ON u.company_id = c.id
WHERE u.deleted_at IS NULL
    AND c.deleted_at IS NULL
    AND u.verified = false
    AND ($1::text = '' OR u.fullname ILIKE '%' || $1 || '%' OR u.email ILIKE '%' || $1 || '%' OR c.companyname ILIKE '%' || $1 || '%')
ORDER BY u.created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountUnVerifiedUsers :one
SELECT COUNT(*)
FROM users u
WHERE
    u.deleted_at IS NULL
    AND u.verified = false;

-- name: GetSingleUser :one
SELECT
    u.id,
    u.fullname,
    u.email,
    u.phonenumber,
    u.picture,
    u.verified,
    u.blocked,
    u.verified_at,
    u.status,
    u.scheduled_for_deletion,
    u.scheduled_deletion_at,
    u.last_login_at,
    u.created_at,
    u.updated_at,
    c.id as company_id,
    c.companyname
FROM users u
    LEFT JOIN companies c ON u.company_id = c.id
WHERE
    u.id = $1
    AND u.deleted_at IS NULL
    AND c.deleted_at IS NULL;

--- Counts for user stats  ---

-- name: CountUserContacts :one
SELECT COUNT(*)
FROM contacts
WHERE
    user_id = $1
    AND deleted_at IS NULL;

-- name: CountUserCampaigns :one
SELECT COUNT(*)
FROM campaigns
WHERE
    user_id = $1
    AND deleted_at IS NULL;

-- name: CountUserTemplates :one
SELECT COUNT(*)
FROM templates
WHERE
    user_id = $1
    AND deleted_at IS NULL;

-- name: CountUserCampaignsSent :one
SELECT COUNT(*)
FROM campaigns
WHERE
    user_id = $1
    AND status = 'SENT'
    AND deleted_at IS NULL;

-- name: CountUserGroups :one
SELECT COUNT(*)
FROM contact_groups
WHERE
    user_id = $1
    AND deleted_at IS NULL;

-- name: GetAllVerifiedUserEmails :many
SELECT id, fullname, email FROM users WHERE verified = true;