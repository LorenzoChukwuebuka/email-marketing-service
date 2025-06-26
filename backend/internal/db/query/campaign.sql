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
            LOWER(c.name) LIKE LOWER('%' || $5 || '%') OR
            LOWER(c.subject) LIKE LOWER('%' || $5 || '%') OR
            LOWER(c.sender) LIKE LOWER('%' || $5 || '%') OR
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
    AND user_id = $2
    AND company_id = $3
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

-- name: CreateCampaignGroups :exec
INSERT INTO
    campaign_groups (campaign_id, contact_group_id)
VALUES ($1, $2);

-- name: CampaignGroupExists :one
SELECT EXISTS (
        SELECT 1
        FROM campaign_groups
        WHERE
            campaign_id = $1
            AND deleted_at IS NULL
    ) AS campaign_group_exists;

-- name: UpdateCampaignGroup :exec
UPDATE campaign_groups
SET
    contact_group_id = $1
WHERE
    campaign_id = $2;

-- name: UpdateCampaignStatus :exec
UPDATE campaigns
SET
    status = $1,
    sent_at = COALESCE($2, sent_at)
WHERE
    id = $3
    AND user_id = $4;
    
-- name: ListScheduledCampaignsByCompanyID :many
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
    AND c.scheduled_at IS NULL
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

/* -- name: GetCampaignWithGroups :many
SELECT
c.id as campaign_id,
c.name as campaign_name,
c.sender,
c.template,
c.sent_at,
c.created_at,
c.scheduled_at,
cg.id as group_id,
cg.name as group_name,
cg.description as group_description
FROM
campaigns c
LEFT JOIN campaign_groups cg ON c.id = cg.campaign_id
WHERE
c.id = $1
AND c.user_id = $2; */

-- name: GetCampaignContactEmails :many
SELECT DISTINCT
    c.email
FROM
    contacts c
    JOIN user_contact_groups ucg ON c.id = ucg.contact_id
    JOIN contact_groups cg ON ucg.contact_group_id = cg.id
    JOIN campaign_groups camp_g ON cg.id = camp_g.contact_group_id
WHERE
    camp_g.campaign_id = $1
    AND c.is_subscribed = true
    AND c.deleted_at IS NULL
    AND ucg.deleted_at IS NULL
    AND cg.deleted_at IS NULL
    AND camp_g.deleted_at IS NULL;

-- name: GetCampaignContactGroups :many
SELECT cg.id, cg.group_name, cg.description, cg.created_at
FROM
    contact_groups cg
    JOIN campaign_groups camp_g ON cg.id = camp_g.contact_group_id
WHERE
    camp_g.campaign_id = $1
    AND cg.deleted_at IS NULL
    AND camp_g.deleted_at IS NULL;

-- name: CreateEmailCampaignResult :one
INSERT INTO
    email_campaign_results (
        company_id,
        campaign_id,
        recipient_email,
        recipient_name,
        version,
        sent_at,
        opened_at,
        open_count,
        clicked_at,
        click_count,
        conversion_at,
        bounce_status,
        unsubscribed_at,
        complaint_status,
        device_type,
        location,
        retry_count,
        notes
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
        $14,
        $15,
        $16,
        $17,
        $18
    ) RETURNING *;

-- name: UpdateEmailCampaignResult :one
UPDATE email_campaign_results
SET
    recipient_name = COALESCE($2, recipient_name),
    version = COALESCE($3, version),
    sent_at = COALESCE($4, sent_at),
    opened_at = COALESCE($5, opened_at),
    open_count = COALESCE($6, open_count),
    clicked_at = COALESCE($7, clicked_at),
    click_count = COALESCE($8, click_count),
    conversion_at = COALESCE($9, conversion_at),
    bounce_status = COALESCE($10, bounce_status),
    unsubscribed_at = COALESCE($11, unsubscribed_at),
    complaint_status = COALESCE($12, complaint_status),
    device_type = COALESCE($13, device_type),
    location = COALESCE($14, location),
    retry_count = COALESCE($15, retry_count),
    notes = COALESCE($16, notes),
    updated_at = CURRENT_TIMESTAMP
WHERE
    campaign_id = $1
    AND recipient_email = $17
    AND deleted_at IS NULL RETURNING *;

-- name: GetEmailCampaignResult :one
SELECT *
FROM email_campaign_results
WHERE
    campaign_id = $1
    AND recipient_email = $2
    AND deleted_at IS NULL;

-- name: GetEmailCampaignResultsByCampaign :many
SELECT *
FROM email_campaign_results
WHERE
    campaign_id = $1
    AND company_id = $2
    AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: GetEmailCampaignResultsByRecipient :many
SELECT *
FROM email_campaign_results
WHERE
    recipient_email = $1
    AND company_id = $2
    AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdateEmailOpened :one
UPDATE email_campaign_results
SET
    opened_at = COALESCE($3, opened_at),
    open_count = open_count + 1,
    device_type = COALESCE($4, device_type),
    location = COALESCE($5, location),
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
    AND company_id = $2
    AND deleted_at IS NULL RETURNING *;

-- name: UpdateEmailClicked :one
UPDATE email_campaign_results
SET
    clicked_at = COALESCE($3, clicked_at),
    click_count = click_count + 1,
    device_type = COALESCE($4, device_type),
    location = COALESCE($5, location),
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
    AND company_id = $2
    AND deleted_at IS NULL RETURNING *;

-- name: UpdateEmailBounced :one
UPDATE email_campaign_results
SET
    bounce_status = $3,
    retry_count = retry_count + 1,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
    AND company_id = $2
    AND deleted_at IS NULL RETURNING *;

-- name: UpdateEmailUnsubscribed :one
UPDATE email_campaign_results
SET
    unsubscribed_at = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
    AND company_id = $2
    AND deleted_at IS NULL RETURNING *;

-- name: UpdateEmailComplaint :one
UPDATE email_campaign_results
SET
    complaint_status = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
    AND company_id = $2
    AND deleted_at IS NULL RETURNING *;

-- name: UpdateEmailConversion :one
UPDATE email_campaign_results
SET
    conversion_at = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
    AND company_id = $2
    AND deleted_at IS NULL RETURNING *;

-- name: SoftDeleteEmailCampaignResult :exec
UPDATE email_campaign_results
SET
    deleted_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1
    AND company_id = $2;

-- name: GetEmailCampaignStats :one
SELECT
    COUNT(*) as total_sent,
    COUNT(opened_at) as total_opened,
    COUNT(clicked_at) as total_clicked,
    COUNT(conversion_at) as total_conversions,
    COUNT(
        CASE
            WHEN bounce_status IS NOT NULL THEN 1
        END
    ) as total_bounced,
    COUNT(unsubscribed_at) as total_unsubscribed,
    COUNT(
        CASE
            WHEN complaint_status = true THEN 1
        END
    ) as total_complaints
FROM email_campaign_results
WHERE
    campaign_id = $1
    AND company_id = $2
    AND deleted_at IS NULL;

-- name: GetUserCampaignStats :one
SELECT
    COUNT(*) as total_emails_sent,
    COALESCE(SUM(open_count), 0) as total_opens,
    COUNT(CASE WHEN open_count > 0 THEN 1 END) as unique_opens,
    COALESCE(SUM(click_count), 0) as total_clicks,
    COUNT(CASE WHEN click_count > 0 THEN 1 END) as unique_clicks,
    COUNT(CASE WHEN bounce_status = 'soft' THEN 1 END) as soft_bounces,
    COUNT(CASE WHEN bounce_status = 'hard' THEN 1 END) as hard_bounces,
    COUNT(CASE WHEN bounce_status IN ('soft', 'hard') THEN 1 END) as total_bounces,
    COUNT(*) - COUNT(CASE WHEN bounce_status IN ('soft', 'hard') THEN 1 END) as total_deliveries
FROM email_campaign_results ecr
WHERE ecr.campaign_id IN (
    SELECT c.id 
    FROM campaigns c 
    WHERE c.user_id = $1 
    AND c.deleted_at IS NULL
)
AND ecr.deleted_at IS NULL;

-- name: GetAllCampaignsByUser :many
SELECT 
    c.id as campaign_id,
    c.name,
    c.sent_at
FROM campaigns c
WHERE c.user_id = $1 
AND c.deleted_at IS NULL;

-- name: GetCampaignStats :one
SELECT
    COUNT(*) as total_emails_sent,
    COALESCE(SUM(open_count), 0) as total_opens,
    COUNT(CASE WHEN open_count > 0 THEN 1 END) as unique_opens,
    COALESCE(SUM(click_count), 0) as total_clicks,
    COUNT(CASE WHEN click_count > 0 THEN 1 END) as unique_clicks,
    COUNT(CASE WHEN bounce_status = 'soft' THEN 1 END) as soft_bounces,
    COUNT(CASE WHEN bounce_status = 'hard' THEN 1 END) as hard_bounces,
    COUNT(CASE WHEN bounce_status IN ('soft', 'hard') THEN 1 END) as total_bounces,
    COUNT(*) - COUNT(CASE WHEN bounce_status IN ('soft', 'hard') THEN 1 END) as total_deliveries,
    COUNT(unsubscribed_at) as unsubscribed,
    COUNT(CASE WHEN complaint_status = true THEN 1 END) as complaints
FROM email_campaign_results
WHERE campaign_id = $1
AND deleted_at IS NULL;

-- name: CreateCampaignError :exec
INSERT INTO campaign_errors (
    campaign_id,
    error_type,
    error_message
) VALUES (
    $1, $2, $3
);

-- name: GetScheduledCampaignsDue :many
-- Get campaigns that are scheduled and due to be sent
SELECT 
    id,
    company_id,
    name,
    subject,
    preview_text,
    user_id,
    sender_from_name,
    template_id,
    sent_template_id,
    recipient_info,
    is_published,
    status,
    track_type,
    is_archived,
    sent_at,
    sender,
    scheduled_at,
    has_custom_logo,
    created_at,
    updated_at,
    deleted_at
FROM campaigns 
WHERE 
    scheduled_at IS NOT NULL 
    AND scheduled_at <= $1 
    AND (sent_at IS NULL OR sent_at = '1970-01-01 00:00:00')
    AND (is_archived IS NULL OR is_archived = false)
    AND deleted_at IS NULL
    AND status IN ('draft', 'scheduled')
ORDER BY scheduled_at ASC;

-- name: ClearCampaignSchedule :exec
-- Clear the scheduled_at field after processing to prevent reprocessing
UPDATE campaigns 
SET scheduled_at = NULL, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

/* -- name: GetCampaignContactGroups :many
-- Get contact groups associated with a campaign
SELECT 
    cg.id,
    cg.campaign_id,
    cg.contact_group_id,
    cg.created_at,
    cg.updated_at
FROM campaign_groups cg
WHERE cg.campaign_id = $1 
AND cg.deleted_at IS NULL;
  */
