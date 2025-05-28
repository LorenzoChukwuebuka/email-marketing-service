-- name: CreateTemplate :one
INSERT INTO
    templates (
        user_id,
        company_id,
        template_name,
        sender_name,
        from_email,
        subject,
        type,
        email_html,
        email_design,
        is_editable,
        is_published,
        is_public_template,
        is_gallery_template,
        tags,
        description,
        image_url,
        is_active,
        editor_type
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


-- name: ListTemplates :many
SELECT
    t.*,
    u.fullname AS user_fullname,
    u.email AS user_email,
    c.companyname AS company_name
FROM
    templates t
    LEFT JOIN users u ON t.user_id = u.id
    LEFT JOIN companies c ON t.company_id = c.id
WHERE
    t.deleted_at IS NULL
ORDER BY t.created_at DESC
LIMIT $1
OFFSET
    $2;

-- name: ListTemplatesByCompanyID :many
SELECT
    t.*,
    u.fullname AS user_fullname,
    u.email AS user_email,
    u.picture AS user_picture,
    c.companyname AS company_name
FROM
    templates t
    LEFT JOIN users u ON t.user_id = u.id
    LEFT JOIN companies c ON t.company_id = c.id
WHERE
    t.company_id = $1
    AND t.deleted_at IS NULL
ORDER BY t.created_at DESC
LIMIT $2
OFFSET
    $3;

-- name: ListTemplatesByUserID :many
SELECT
    t.*,
    u.fullname AS user_fullname,
    u.email AS user_email,
    u.picture AS user_picture,
    c.companyname AS company_name
FROM
    templates t
    LEFT JOIN users u ON t.user_id = u.id
    LEFT JOIN companies c ON t.company_id = c.id
WHERE
    t.user_id = $1
    AND t.deleted_at IS NULL
ORDER BY t.created_at DESC
LIMIT $2
OFFSET
    $3;


-- name: UpdateTemplate :one
UPDATE templates
SET
    template_name = COALESCE($1, template_name),
    sender_name = COALESCE($2, sender_name),
    from_email = COALESCE($3, from_email),
    subject = COALESCE($4, subject),
    type = COALESCE($5, type),
    email_html = COALESCE($6, email_html),
    email_design = COALESCE($7, email_design),
    is_editable = COALESCE($8, is_editable),
    is_published = COALESCE($9, is_published),
    is_public_template = COALESCE($10, is_public_template),
    is_gallery_template = COALESCE($11, is_gallery_template),
    tags = COALESCE($12, tags),
    description = COALESCE($13, description),
    image_url = COALESCE($14, image_url),
    is_active = COALESCE($15, is_active),
    editor_type = COALESCE($16, editor_type),
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $17
    AND user_id = $18
    AND deleted_at IS NULL RETURNING *;

-- name: SoftDeleteTemplate :exec
UPDATE templates
SET
    deleted_at = CURRENT_TIMESTAMP
WHERE
    id = $1
    AND user_id = $2
    AND deleted_at IS NULL;

-- name: HardDeleteTemplate :exec
DELETE FROM templates WHERE id = $1;

-- name: CheckTemplateNameExists :one
SELECT EXISTS (
        SELECT 1
        FROM templates
        WHERE
            template_name = $1
            AND user_id = $2
            AND deleted_at IS NULL
    ) AS template_exists;

-- name: GetTemplateByID :one
SELECT
    t.*,
    u.fullname AS user_fullname,
    u.email AS user_email,
    u.picture AS user_picture,
    c.companyname AS company_name
FROM
    templates t
    LEFT JOIN users u ON t.user_id = u.id
    LEFT JOIN companies c ON t.company_id = c.id
WHERE
    t.id = $1
    AND t.user_id = $2
    AND t.type = $3
    AND t.deleted_at IS NULL
ORDER BY t.created_at DESC
LIMIT 1;

-- name: ListTemplatesByType :many
SELECT 
    t.*,
    u.fullname AS user_fullname,
    u.email AS user_email,
    u.picture AS user_picture,
    c.companyname AS company_name
FROM 
    templates t
LEFT JOIN 
    users u ON t.user_id = u.id
LEFT JOIN 
    companies c ON t.company_id = c.id
WHERE 
    t.type = $1
    AND t.user_id = $2
    AND t.deleted_at IS NULL
    AND ($5 = '' OR t.template_name ILIKE '%' || $5 || '%')
ORDER BY 
    t.created_at DESC
LIMIT $3
OFFSET $4;

-- name: CountTemplatesByUserID :one
SELECT COUNT(*) 
FROM templates 
WHERE user_id = $1 
AND deleted_at IS NULL;