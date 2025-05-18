-- name: CreateContact :exec
INSERT INTO contacts (
    id, company_id, first_name, last_name, email, from_origin, is_subscribed, user_id, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
);

-- name: CheckIfEmailExists :one
SELECT EXISTS(
    SELECT 1 FROM contacts
    WHERE email = $1 AND user_id = $2 AND deleted_at IS NULL
) AS exists;

/* -- name: BulkCreateContacts :copyfrom
INSERT INTO contacts (
    id, company_id, first_name, last_name, email, from_origin, is_subscribed, user_id, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
); */

-- name: GetContactByID :one
SELECT 
    c.id,
    c.company_id,
    c.first_name,
    c.last_name,
    c.email,
    c.from_origin,
    c.is_subscribed,
    c.user_id,
    c.created_at,
    c.updated_at,
    c.deleted_at
FROM contacts c
WHERE c.id = $1 AND c.user_id = $2 AND c.deleted_at IS NULL;

-- name: ListContacts :many
SELECT 
    c.id,
    c.company_id,
    c.first_name,
    c.last_name,
    c.email,
    c.from_origin,
    c.is_subscribed,
    c.user_id,
    c.created_at,
    c.updated_at,
    c.deleted_at
FROM contacts c
WHERE c.user_id = $1 
AND c.deleted_at IS NULL
AND ($2::text = '' OR 
    c.first_name ILIKE concat('%', $2, '%') OR 
    c.last_name ILIKE concat('%', $2, '%') OR 
    c.email ILIKE concat('%', $2, '%')
)
ORDER BY c.created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountContacts :one
SELECT 
    COUNT(*) as total_count
FROM contacts c
WHERE c.user_id = $1 
AND c.deleted_at IS NULL
AND ($2::text = '' OR 
    c.first_name ILIKE concat('%', $2, '%') OR 
    c.last_name ILIKE concat('%', $2, '%') OR 
    c.email ILIKE concat('%', $2, '%')
);

-- name: DeleteContact :exec
UPDATE contacts 
SET deleted_at = now() 
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;

-- name: UpdateContact :exec
UPDATE contacts
SET 
    first_name = CASE WHEN $3::text != '' THEN $3 ELSE first_name END,
    last_name = CASE WHEN $4::text != '' THEN $4 ELSE last_name END,
    email = CASE WHEN $5::text != '' THEN $5 ELSE email END,
    from_origin = CASE WHEN $6::text != '' THEN $6 ELSE from_origin END,
    is_subscribed = $7,
    updated_at = now()
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;

-- name: UpdateSubscriptionStatus :exec
UPDATE contacts
SET is_subscribed = false, updated_at = now()
WHERE email = $1 AND deleted_at IS NULL;

-- name: CreateContactGroup :exec
INSERT INTO contact_groups (
    id, company_id, group_name, user_id, description, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: AddContactToGroup :exec
INSERT INTO user_contact_groups (
    id, user_id, contact_group_id, contact_id, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: ListContactGroups :many
SELECT 
    cg.id,
    cg.company_id,
    cg.group_name,
    cg.user_id,
    cg.description,
    cg.created_at,
    cg.updated_at,
    cg.deleted_at
FROM contact_groups cg
WHERE cg.user_id = $1 
AND cg.deleted_at IS NULL
AND ($2::text = '' OR cg.group_name ILIKE concat('%', $2, '%'))
ORDER BY cg.created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountContactGroups :one
SELECT 
    COUNT(*) as total_count
FROM contact_groups cg
WHERE cg.user_id = $1 
AND cg.deleted_at IS NULL
AND ($2::text = '' OR cg.group_name ILIKE concat('%', $2, '%'));

-- name: GetContactGroupByID :one
SELECT 
    cg.id,
    cg.company_id,
    cg.group_name,
    cg.user_id,
    cg.description,
    cg.created_at,
    cg.updated_at,
    cg.deleted_at
FROM contact_groups cg
WHERE cg.id = $1 AND cg.user_id = $2 AND cg.deleted_at IS NULL;

-- name: GetContactsInGroup :many
SELECT 
    c.id,
    c.company_id,
    c.first_name,
    c.last_name,
    c.email,
    c.from_origin,
    c.is_subscribed,
    c.user_id,
    c.created_at,
    c.updated_at,
    c.deleted_at
FROM contacts c
JOIN user_contact_groups ucg ON ucg.contact_id = c.id
WHERE ucg.contact_group_id = $1 
AND ucg.user_id = $2
AND c.deleted_at IS NULL
AND ucg.deleted_at IS NULL;

-- name: DeleteContactGroup :exec
-- This is a transaction that will be handled in Go code
UPDATE contact_groups
SET deleted_at = now()
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;

-- name: DeleteContactGroupRelations :exec
UPDATE user_contact_groups
SET deleted_at = now()
WHERE contact_group_id = $1;

-- name: RemoveContactFromGroup :exec
UPDATE user_contact_groups
SET deleted_at = now()
WHERE contact_group_id = $1 AND user_id = $2 AND contact_id = $3 AND deleted_at IS NULL;

-- name: UpdateContactGroup :exec
UPDATE contact_groups
SET 
    group_name = $3,
    description = $4,
    updated_at = now()
WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL;

-- name: CheckIfGroupNameExists :one
SELECT EXISTS(
    SELECT 1 FROM contact_groups
    WHERE group_name = $1 AND user_id = $2 AND deleted_at IS NULL
) AS exists;

-- name: GetContactCount :one
SELECT 
    COUNT(*) as total_count,
    SUM(CASE WHEN created_at >= NOW() - INTERVAL '30 days' THEN 1 ELSE 0 END) as recent_count
FROM contacts
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: GetContactSubscriptionStatsDashboard :one
SELECT
    COUNT(*) as total_count,
    SUM(CASE WHEN is_subscribed = false THEN 1 ELSE 0 END) as unsubscribed_count,
    SUM(CASE WHEN created_at >= NOW() - INTERVAL '10 days' THEN 1 ELSE 0 END) as new_contacts_count
FROM contacts
WHERE user_id = $1 AND deleted_at IS NULL;

-- name: GetEngagedContactsCount :one
SELECT 
    COUNT(DISTINCT c.id) as engaged_count
FROM contacts c
JOIN email_campaign_results ecr ON c.email = ecr.recipient_email
WHERE c.user_id = $1 
AND c.deleted_at IS NULL
AND (ecr.opened_at IS NOT NULL OR ecr.clicked_at IS NOT NULL OR ecr.conversion_at IS NOT NULL);