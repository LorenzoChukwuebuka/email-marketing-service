-- name: ArchiveCampaign :one
UPDATE campaigns
SET 
    is_archived = true,
    updated_at = CURRENT_TIMESTAMP
WHERE 
    id = $1
    AND deleted_at IS NULL
RETURNING *;

-- name: ListCampaignsByUserID :many
SELECT 
    c.*,
    u.fullname AS user_fullname,
    u.email AS user_email,
    comp.companyname AS company_name
FROM 
    campaigns c
LEFT JOIN 
    users u ON c.user_id = u.id
LEFT JOIN 
    companies comp ON c.company_id = comp.id
WHERE 
    c.user_id = $1
    AND c.deleted_at IS NULL
ORDER BY 
    c.created_at DESC
LIMIT $2
OFFSET $3;

-- name: CreateCampaign :one
INSERT INTO campaigns (
    company_id,
    name,
    subject,
    preview_text,
    user_id,
    sender_from_name,
    template_id,
    recipient_info,
    status,
    track_type,
    sender,
    scheduled_at,
    has_custom_logo
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
RETURNING *;

-- name: GetCampaignByID :one
SELECT 
    c.*,
    u.fullname AS user_fullname,
    u.email AS user_email,
    u.picture AS user_picture
FROM 
    campaigns c
LEFT JOIN 
    users u ON c.user_id = u.id
WHERE 
    c.id = $1
    AND c.deleted_at IS NULL;

-- name: ListCampaigns :many
SELECT 
    c.*,
    u.fullname AS user_fullname,
    u.email AS user_email,
    comp.companyname AS company_name
FROM 
    campaigns c
LEFT JOIN 
    users u ON c.user_id = u.id
LEFT JOIN 
    companies comp ON c.company_id = comp.id
WHERE 
    c.deleted_at IS NULL
ORDER BY 
    c.created_at DESC
LIMIT $1
OFFSET $2;

-- name: CheckCampaignNameExists :one
SELECT 
    EXISTS (
        SELECT 1 
        FROM campaigns 
        WHERE name = $1 
        AND company_id = $2
        AND deleted_at IS NULL
    ) AS campaign_exists;

-- name: ListCampaignsByCompanyID :many
SELECT 
    c.*,
    u.fullname AS user_fullname,
    u.email AS user_email,
    comp.companyname AS company_name
FROM 
    campaigns c
LEFT JOIN 
    users u ON c.user_id = u.id
LEFT JOIN 
    companies comp ON c.company_id = comp.id
WHERE 
    c.company_id = $1
    AND c.deleted_at IS NULL
ORDER BY 
    c.created_at DESC
LIMIT $2
OFFSET $3;

-- name: GetCampaignsByTemplateType :many
SELECT 
    c.*,
    t.type AS template_type,
    u.fullname AS user_fullname,
    comp.companyname AS company_name
FROM 
    campaigns c
JOIN 
    templates t ON c.template_id = t.id
LEFT JOIN 
    users u ON c.user_id = u.id
LEFT JOIN 
    companies comp ON c.company_id = comp.id
WHERE 
    t.type = $1
    AND c.deleted_at IS NULL
ORDER BY 
    c.created_at DESC
LIMIT $2
OFFSET $3;

-- name: UpdateCampaign :one
UPDATE campaigns
SET 
    name = COALESCE($1, name),
    subject = COALESCE($2, subject),
    preview_text = COALESCE($3, preview_text),
    sender_from_name = COALESCE($4, sender_from_name),
    template_id = COALESCE($5, template_id),
    recipient_info = COALESCE($6, recipient_info),
    status = COALESCE($7, status),
    track_type = COALESCE($8, track_type),
    sender = COALESCE($9, sender),
    scheduled_at = COALESCE($10, scheduled_at),
    has_custom_logo = COALESCE($11, has_custom_logo),
    is_published = COALESCE($12, is_published),
    is_archived = COALESCE($13, is_archived),
    updated_at = CURRENT_TIMESTAMP
WHERE 
    id = $14
    AND deleted_at IS NULL
RETURNING *;

-- name: MarkCampaignAsSent :one
UPDATE campaigns
SET 
    status = 'sent',
    sent_at = CURRENT_TIMESTAMP,
    sent_template_id = template_id,
    updated_at = CURRENT_TIMESTAMP
WHERE 
    id = $1
    AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteCampaign :exec
UPDATE campaigns
SET 
    deleted_at = CURRENT_TIMESTAMP
WHERE 
    id = $1
    AND deleted_at IS NULL;

-- name: HardDeleteCampaign :exec
DELETE FROM campaigns
WHERE id = $1;