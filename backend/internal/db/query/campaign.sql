-- name: ArchiveCampaign :one
UPDATE campaigns
SET
    is_archived = true,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
    AND deleted_at IS NULL RETURNING *;

-- name: ListCampaignsByUserID :many
SELECT
    c.*,
    u.fullname AS user_fullname,
    u.email AS user_email,
    comp.companyname AS company_name
FROM
    campaigns c
    LEFT JOIN users u ON c.user_id = u.id
    LEFT JOIN companies comp ON c.company_id = comp.id
WHERE
    c.user_id = $1
    AND c.deleted_at IS NULL
ORDER BY c.created_at DESC
LIMIT $2
OFFSET
    $3;

-- name: CreateCampaign :one
INSERT INTO
    campaigns (
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
        $13
    ) RETURNING *;

-- name: GetCampaignByID :one
SELECT
    -- Campaign information (all columns)
    c.*,
    -- User information (with user_ prefix)
    u.id AS user_id,
    u.fullname AS user_fullname,
    u.email AS user_email,
    u.phonenumber AS user_phonenumber,
    u.picture AS user_picture,
    u.verified AS user_verified,
    u.blocked AS user_blocked,
    u.verified_at AS user_verified_at,
    u.status AS user_status,
    u.last_login_at AS user_last_login_at,
    u.created_at AS user_created_at,
    u.updated_at AS user_updated_at,
    -- Company information (with company_ prefix)
    comp.id AS company_id_ref,
    comp.companyname AS company_name,
    comp.created_at AS company_created_at,
    comp.updated_at AS company_updated_at,
    -- Template information (all columns with template_ prefix)
    t.id AS template_id_ref,
    t.user_id AS template_user_id,
    t.company_id AS template_company_id,
    t.template_name AS template_name,
    t.sender_name AS template_sender_name,
    t.from_email AS template_from_email,
    t.subject AS template_subject,
    t.type AS template_type,
    t.email_html AS template_email_html,
    t.email_design AS template_email_design,
    t.is_editable AS template_is_editable,
    t.is_published AS template_is_published,
    t.is_public_template AS template_is_public_template,
    t.is_gallery_template AS template_is_gallery_template,
    t.tags AS template_tags,
    t.description AS template_description,
    t.image_url AS template_image_url,
    t.is_active AS template_is_active,
    t.editor_type AS template_editor_type,
    t.created_at AS template_created_at,
    t.updated_at AS template_updated_at,
    t.deleted_at AS template_deleted_at
FROM
    campaigns c
    INNER JOIN users u ON c.user_id = u.id
    AND u.deleted_at IS NULL
    AND u.blocked = FALSE
    INNER JOIN companies comp ON c.company_id = comp.id
    AND comp.deleted_at IS NULL
    LEFT JOIN templates t ON c.template_id = t.id
    AND t.deleted_at IS NULL
    AND t.is_active = TRUE
WHERE
    c.company_id = $1
    AND c.id = $3
    AND c.deleted_at IS NULL
    AND c.user_id = $2;

-- name: CheckCampaignNameExists :one
SELECT EXISTS (
        SELECT 1
        FROM campaigns
        WHERE
            name = $1
            AND company_id = $2
            AND deleted_at IS NULL
    ) AS campaign_exists;

-- name: ListCampaignsByCompanyID :many
SELECT 
    -- Campaign information (all columns)
    c.*,
    -- User information (with user_ prefix)
    u.id AS user_id,
    u.fullname AS user_fullname,
    u.email AS user_email,
    u.phonenumber AS user_phonenumber,
    u.picture AS user_picture,
    u.verified AS user_verified,
    u.blocked AS user_blocked,
    u.verified_at AS user_verified_at,
    u.status AS user_status,
    u.last_login_at AS user_last_login_at,
    u.created_at AS user_created_at,
    u.updated_at AS user_updated_at,
    -- Company information (with company_ prefix)
    comp.id AS company_id_ref,
    comp.companyname AS company_name,
    comp.created_at AS company_created_at,
    comp.updated_at AS company_updated_at,
    -- Template information (all columns with template_ prefix)
    t.id AS template_id_ref,
    t.user_id AS template_user_id,
    t.company_id AS template_company_id,
    t.template_name AS template_name,
    t.sender_name AS template_sender_name,
    t.from_email AS template_from_email,
    t.subject AS template_subject,
    t.type AS template_type,
    t.email_html AS template_email_html,
    t.email_design AS template_email_design,
    t.is_editable AS template_is_editable,
    t.is_published AS template_is_published,
    t.is_public_template AS template_is_public_template,
    t.is_gallery_template AS template_is_gallery_template,
    t.tags AS template_tags,
    t.description AS template_description,
    t.image_url AS template_image_url,
    t.is_active AS template_is_active,
    t.editor_type AS template_editor_type,
    t.created_at AS template_created_at,
    t.updated_at AS template_updated_at,
    t.deleted_at AS template_deleted_at
FROM 
    campaigns c
INNER JOIN 
    users u ON c.user_id = u.id 
    AND u.deleted_at IS NULL 
    AND u.blocked = FALSE
INNER JOIN 
    companies comp ON c.company_id = comp.id 
    AND comp.deleted_at IS NULL
LEFT JOIN 
    templates t ON c.template_id = t.id 
    AND t.deleted_at IS NULL 
    AND t.is_active = TRUE
WHERE 
    c.company_id = $1
    AND c.deleted_at IS NULL
    AND  c.user_id = $2  
     AND (
        $5::TEXT IS NULL OR $5 = '' OR (
            -- Search in campaign fields
            LOWER(c.campaign_name) LIKE LOWER('%' || $5 || '%') OR
            LOWER(c.subject) LIKE LOWER('%' || $5 || '%') OR
            LOWER(c.description) LIKE LOWER('%' || $5 || '%') OR
            LOWER(c.from_name) LIKE LOWER('%' || $5 || '%') OR
            LOWER(c.from_email) LIKE LOWER('%' || $5 || '%') OR
            -- Search in user fields
            LOWER(u.fullname) LIKE LOWER('%' || $5 || '%') OR
            LOWER(u.email) LIKE LOWER('%' || $5 || '%') OR
            -- Search in template fields
            LOWER(t.template_name) LIKE LOWER('%' || $5 || '%') OR
            LOWER(t.sender_name) LIKE LOWER('%' || $5 || '%') OR
            LOWER(t.subject) LIKE LOWER('%' || $5 || '%') OR
            LOWER(t.description) LIKE LOWER('%' || $5 || '%')
        )
    )
ORDER BY 
    c.created_at DESC
LIMIT $3
OFFSET $4;

-- name: GetCampaignsByTemplateType :many
SELECT
    c.*,
    t.type AS template_type,
    u.fullname AS user_fullname,
    comp.companyname AS company_name
FROM
    campaigns c
    JOIN templates t ON c.template_id = t.id
    LEFT JOIN users u ON c.user_id = u.id
    LEFT JOIN companies comp ON c.company_id = comp.id
WHERE
    t.type = $1
    AND c.deleted_at IS NULL
ORDER BY c.created_at DESC
LIMIT $2
OFFSET
    $3;

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
    AND user_id = $15
    AND deleted_at IS NULL RETURNING *;

-- name: MarkCampaignAsSent :one
UPDATE campaigns
SET
    status = 'sent',
    sent_at = CURRENT_TIMESTAMP,
    sent_template_id = template_id,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
    AND deleted_at IS NULL RETURNING *;

-- name: SoftDeleteCampaign :exec
UPDATE campaigns
SET
    deleted_at = CURRENT_TIMESTAMP
WHERE
    id = $1
    AND deleted_at IS NULL;

-- name: HardDeleteCampaign :exec
DELETE FROM campaigns WHERE id = $1;

-- name: GetCampaignCounts :one
SELECT COUNT(*)
FROM campaigns
WHERE
    user_id = $1
    AND company_id = $2
    AND deleted_at IS NULL;