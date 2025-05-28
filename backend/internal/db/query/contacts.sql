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
    c.user_id = @user_id
    AND c.company_id = @company_id
    AND c.deleted_at IS NULL
    AND (
        c.first_name ILIKE '%' || @search_term || '%'
        OR c.last_name ILIKE '%' || @search_term || '%'
        OR c.email ILIKE '%' || @search_term || '%'
    )
ORDER BY c.first_name, c.last_name
LIMIT @row_limit
OFFSET
    @row_offset;

-- name: CreateContactGroup :one
-- Creates a new contact group
INSERT INTO
    contact_groups (
        company_id,
        group_name,
        user_id,
        description
    )
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetContactGroup :one
-- Gets a contact group by ID
SELECT * FROM contact_groups WHERE id = $1 AND deleted_at IS NULL;

-- name: GetContactGroupByName :one
-- Gets a contact group by name within a company
SELECT *
FROM contact_groups
WHERE
    company_id = $1
    AND group_name = $2
    AND deleted_at IS NULL;

-- name: ListContactGroups :many
-- Lists all contact groups for a company with pagination
SELECT *
FROM contact_groups
WHERE
    company_id = $1
    AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- name: ListContactGroupsForUser :many
-- Lists all contact groups for a specific user with pagination
SELECT *
FROM contact_groups
WHERE
    company_id = $1
    AND user_id = $2
    AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $3
OFFSET
    $4;

-- name: CountContactGroups :one
-- Counts total number of contact groups for a company
SELECT COUNT(*)
FROM contact_groups
WHERE
    company_id = $1
    AND deleted_at IS NULL;

-- name: UpdateContactGroup :one
-- Updates a contact group's details
UPDATE contact_groups
SET
    group_name = $2,
    description = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
    AND user_id = $4
    AND deleted_at IS NULL RETURNING *;

-- name: SoftDeleteContactGroup :exec
-- Soft deletes a contact group
UPDATE contact_groups
SET
    deleted_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
    AND user_id = $2
    AND deleted_at IS NULL;

-- name: RestoreContactGroup :exec
-- Restores a soft-deleted contact group
UPDATE contact_groups
SET
    deleted_at = NULL,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1;

-- name: HardDeleteContactGroup :exec
-- Hard deletes a contact group (use with caution)
DELETE FROM contact_groups WHERE id = $1;

-- name: SearchContactGroups :many
-- Searches contact groups by name or description
SELECT *
FROM contact_groups
WHERE
    company_id = $1
    AND deleted_at IS NULL
    AND (
        group_name ILIKE '%' || $2 || '%'
        OR description ILIKE '%' || $2 || '%'
    )
ORDER BY created_at DESC
LIMIT $3
OFFSET
    $4;

-- name: IsContactGroupNameUnique :one
-- Checks if a contact group name is unique within a company
SELECT NOT EXISTS (
        SELECT 1
        FROM contact_groups
        WHERE
            company_id = @companyID
            AND group_name = @groupname
            AND user_id = @userID
            AND deleted_at IS NULL
    ) AS is_unique;

-- name: AddContactToGroup :one
-- Adds a contact to a group and returns the created entry
INSERT INTO
    user_contact_groups (
        user_id,
        contact_group_id,
        contact_id
    )
VALUES ($1, $2, $3) RETURNING *;

-- name: RemoveContactFromGroup :exec
-- Soft deletes a contact from a group by setting the deleted_at timestamp
UPDATE user_contact_groups
SET
    deleted_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE
    user_id = $1
    AND contact_group_id = $2
    AND contact_id = $3
    AND deleted_at IS NULL;

-- name: IsContactInGroup :one
-- Checks if a contact is already in a specific group
SELECT EXISTS (
        SELECT 1
        FROM user_contact_groups
        WHERE
            user_id = $1
            AND contact_group_id = $2
            AND contact_id = $3
            AND deleted_at IS NULL
    ) AS is_in_group;

-- name: GetGroupsWithContacts :many
-- Fetches all contact groups with their associated contacts for a specific user and company
-- with pagination support using limit and offset
SELECT
    cg.id AS group_id,
    cg.group_name,
    cg.description,
    cg.created_at AS group_created_at,
    c.id AS contact_id,
    c.first_name as contact_first_name,
    c.last_name as contact_last_name,
    c.email as contact_email,
    c.from_origin as contact_from_origin,
    c.is_subscribed as contact_is_subscribed,
    c.created_at AS contact_created_at
FROM
    contact_groups cg
    LEFT JOIN user_contact_groups ucg ON cg.id = ucg.contact_group_id
    AND ucg.deleted_at IS NULL
    LEFT JOIN contacts c ON ucg.contact_id = c.id
    AND c.deleted_at IS NULL
WHERE
    cg.company_id = @company_id
    AND cg.user_id = @user_id
    AND (
        @searchterm = ''
        OR LOWER(cg.group_name) LIKE LOWER('%' || @searchterm || '%')
        OR LOWER(c.first_name) LIKE LOWER('%' || @searchterm || '%')
        OR LOWER(c.last_name) LIKE LOWER('%' || @searchterm || '%')
        OR LOWER(c.email) LIKE LOWER('%' || @searchterm || '%')
    )
    AND cg.deleted_at IS NULL
ORDER BY cg.group_name, c.last_name, c.first_name
LIMIT @rowlimit
OFFSET
    @rowoffset;


-- name: GetSingleGroupWithContacts :many
-- Fetches a specific contact group with all its associated contacts for a specific user and company
SELECT
    cg.id AS group_id,
    cg.group_name,
    cg.description,
    cg.created_at AS group_created_at,
    c.id AS contact_id,
    c.first_name as contact_first_name,
    c.last_name as contact_last_name,
    c.email as contact_email,
    c.from_origin as contact_from_origin,
    c.is_subscribed as contact_is_subscribed,
    c.created_at AS contact_created_at
FROM
    contact_groups cg
    LEFT JOIN user_contact_groups ucg ON cg.id = ucg.contact_group_id
    AND ucg.deleted_at IS NULL
    LEFT JOIN contacts c ON ucg.contact_id = c.id
    AND c.deleted_at IS NULL
WHERE
    cg.id = @group_id
    AND cg.company_id = @company_id
    AND cg.user_id = @user_id
    AND cg.deleted_at IS NULL
ORDER BY c.last_name, c.first_name;

-- name: GetContactUnsubscribedCount :one
-- Get count of unsubscribed contacts for a specific user
SELECT
    COUNT(*) AS count
FROM
    contacts
WHERE
    user_id = @user_id
    AND is_subscribed = false
    AND deleted_at IS NULL;

-- name: GetContactTotalCount :one
-- Get total count of contacts for a specific user
SELECT
    COUNT(*) AS count
FROM
    contacts
WHERE
    user_id = @user_id
    AND deleted_at IS NULL;

-- name: GetNewContactsCount :one
-- Get count of new contacts (less than 10 days old) for a specific user
SELECT
    COUNT(*) AS count
FROM
    contacts
WHERE
    user_id = @user_id
    AND created_at >= @ten_days_ago
    AND deleted_at IS NULL;

-- name: GetContactStats :one
-- Get all contact statistics in a single query
-- name: GetContactStats :one
-- Get all contact statistics in a single query
SELECT
    (SELECT COUNT(*) FROM contacts c1 WHERE c1.user_id = @user_id AND c1.is_subscribed = false AND c1.deleted_at IS NULL) AS unsubscribed_count,
    (SELECT COUNT(*) FROM contacts c2 WHERE c2.user_id = @user_id AND c2.deleted_at IS NULL) AS total_count,
    (SELECT COUNT(*) FROM contacts c3 WHERE c3.user_id = @user_id AND c3.created_at >= @ten_days_ago AND c3.deleted_at IS NULL) AS new_contacts_count,
    (SELECT COUNT(DISTINCT c4.id) 
     FROM contacts c4
     JOIN email_campaign_results ecr ON c4.email = ecr.recipient_email
     WHERE c4.user_id = @user_id 
     AND c4.deleted_at IS NULL
     AND ecr.deleted_at IS NULL
     AND (ecr.opened_at IS NOT NULL OR ecr.clicked_at IS NOT NULL OR ecr.conversion_at IS NOT NULL)
    ) AS engaged_count;

