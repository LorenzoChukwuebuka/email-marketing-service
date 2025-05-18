-- name: CreateAdmin :one
INSERT INTO
    admins (
        firstname,
        middlename,
        lastname,
        email,
        password,
        type,
        created_at,
        updated_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        NOW(),
        NOW()
    ) RETURNING *;

-- name: UpsertAdmin :one
INSERT INTO
    admins (
        firstname,
        middlename,
        lastname,
        email,
        password,
        type,
        created_at,
        updated_at
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        NOW(),
        NOW()
    ) ON CONFLICT (email) DO
UPDATE
SET
    firstname = EXCLUDED.first_name,
    middlename = EXCLUDED.middle_name,
    lastname = EXCLUDED.last_name,
    password = EXCLUDED.password,
    type = EXCLUDED.type,
    updatedat = NOW() RETURNING *;

-- name: GetAdminByEmail :one
SELECT * FROM admins WHERE email = $1 AND deleted_at IS NULL;

-- name: GetAdminByID :one
SELECT * FROM admins WHERE id = $1 AND deleted_at IS NULL;

-- name: ListAdmins :many
SELECT *
FROM admins
WHERE
    deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateAdmin :one
UPDATE admins
SET
    firstname = COALESCE($1, first_name),
    middlename = COALESCE($2, middle_name),
    lastname = COALESCE($3, last_name),
    password = COALESCE($4, password),
    type = COALESCE($5, type),
    updated_at = NOW()
WHERE
    id = $6
    AND deleted_at IS NULL RETURNING *;

-- name: SoftDeleteAdmin :exec
UPDATE admins
SET
    deleted_at = NOW()
WHERE
    id = $1
    AND deleted_at IS NULL;

 