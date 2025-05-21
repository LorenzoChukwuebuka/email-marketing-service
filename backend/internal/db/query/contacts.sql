-- name: CreateContact :exec
INSERT INTO
    contacts (
        company_id,
        first_name,
        last_name,
        email,
        from_origin,
        is_subscribed,
        user_id,
        created_at
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
    );

-- name: CheckIfContactEmailExists :one
SELECT EXISTS (
        SELECT 1
        FROM contacts
        WHERE
            email = $1
            AND user_id = $2
            AND deleted_at IS NULL
    ) AS exists;

-- name: DeleteContact :exec
UPDATE contacts
SET
    deleted_at = now()
WHERE
    id = $1
    AND user_id = $2
    AND deleted_at IS NULL;

-- name: UpdateContact :exec
UPDATE contacts
SET
    first_name = COALESCE(
        sqlc.narg (first_name),
        first_name
    ),
    last_name = COALESCE(
        sqlc.narg (last_name),
        last_name
    ),
    email = COALESCE(sqlc.narg (email), email),
    from_origin = COALESCE(
        sqlc.narg (from_origin),
        from_origin
    ),
    is_subscribed = sqlc.arg (is_subscribed),
    updated_at = now()
WHERE
    id = sqlc.arg (id)
    AND user_id = sqlc.arg (user_id)
    AND deleted_at IS NULL;

-- name: GetContactsCount :one
SELECT COUNT(*) AS total_count
FROM contacts c
WHERE
    c.user_id = $1
    AND c.company_id = $2
    AND c.deleted_at IS NULL;

-- name: GetAllContacts :many
SELECT
    c.id AS contact_id,
    c.company_id,
    c.first_name,
    c.last_name,
    c.email,
    c.from_origin,
    c.is_subscribed,
    c.user_id,
    c.created_at AS contact_created_at,
    c.updated_at AS contact_updated_at,
    ucg.id AS user_contact_group_id,
    ucg.user_id AS ucg_user_id,
    ucg.contact_group_id,
    ucg.contact_id,
    ucg.created_at AS ucg_created_at,
    ucg.updated_at AS ucg_updated_at,
    ucg.deleted_at AS ucg_deleted_at,
    cg.id AS group_id,
    cg.group_name,
    cg.description,
    cg.user_id AS group_creator_id,
    cg.created_at AS group_created_at,
    cg.updated_at AS group_updated_at
FROM
    contacts c
    LEFT JOIN user_contact_groups ucg ON c.id = ucg.contact_id
    AND ucg.deleted_at IS NULL
    LEFT JOIN contact_groups cg ON ucg.contact_group_id = cg.id
WHERE
    c.user_id = $1
    AND c.company_id = $2
    AND c.deleted_at IS NULL
ORDER BY c.first_name, c.last_name
LIMIT $3
OFFSET
    $4;