-- name: CreatePlan :one
INSERT INTO plans (
    id,
    name,
    description,
    price,
    billing_cycle,
    status,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: PlanExists :one
SELECT EXISTS(SELECT 1 FROM plans WHERE name = $1);

-- name: GetPlanByID :one
SELECT * FROM plans
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetPlanByName :one
SELECT * FROM plans
WHERE name = $1 AND deleted_at IS NULL;

-- name: ListActivePlans :many
SELECT * FROM plans
WHERE status = 'active' AND deleted_at IS NULL
ORDER BY price ASC;

-- name: UpdatePlan :one
UPDATE plans
SET 
    name = COALESCE($1, name),
    description = COALESCE($2, description),
    price = COALESCE($3, price),
    billing_cycle = COALESCE($4, billing_cycle),
    status = COALESCE($5, status),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $6 AND deleted_at IS NULL
RETURNING *;

-- name: ArchivePlan :one
UPDATE plans
SET 
    status = 'archived',
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeletePlan :exec
UPDATE plans
SET 
    deleted_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: CreatePlanFeature :one
INSERT INTO plan_features (
    id,
    plan_id,
    name,
    description,
    value,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetPlanFeaturesByPlanID :many
SELECT * FROM plan_features
WHERE plan_id = $1 AND deleted_at IS NULL;

-- name: UpdatePlanFeature :one
UPDATE plan_features
SET 
    name = COALESCE($1, name),
    description = COALESCE($2, description),
    value = COALESCE($3, value),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $4 AND deleted_at IS NULL
RETURNING *;

-- name: DeletePlanFeature :exec
UPDATE plan_features
SET 
    deleted_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: CreateMailingLimit :one
INSERT INTO mailing_limits (
    id,
    plan_id,
    daily_limit,
    monthly_limit,
    max_recipients_per_mail,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetMailingLimitByPlanID :one
SELECT * FROM mailing_limits
WHERE plan_id = $1 AND deleted_at IS NULL;

-- name: UpdateMailingLimit :one
UPDATE mailing_limits
SET 
    daily_limit = COALESCE($1, daily_limit),
    monthly_limit = COALESCE($2, monthly_limit),
    max_recipients_per_mail = COALESCE($3, max_recipients_per_mail),
    updated_at = CURRENT_TIMESTAMP
WHERE plan_id = $4 AND deleted_at IS NULL
RETURNING *;

-- name: GetPlanWithDetails :one
SELECT 
    p.*,
    (
        SELECT json_agg(pf.*)
        FROM plan_features pf
        WHERE pf.plan_id = p.id AND pf.deleted_at IS NULL
    ) as features,
    ml.daily_limit,
    ml.monthly_limit,
    ml.max_recipients_per_mail
FROM plans p
LEFT JOIN mailing_limits ml ON p.id = ml.plan_id AND ml.deleted_at IS NULL
WHERE p.id = $1 AND p.deleted_at IS NULL;

-- name: ListPlansWithDetails :many
SELECT 
    p.*,
    (
        SELECT json_agg(pf.*)
        FROM plan_features pf
        WHERE pf.plan_id = p.id AND pf.deleted_at IS NULL
    ) as features,
    ml.daily_limit,
    ml.monthly_limit,
    ml.max_recipients_per_mail,
    (
        SELECT COUNT(*)
        FROM subscriptions s
        WHERE s.plan_id = p.id AND s.status = 'active' AND s.deleted_at IS NULL
    ) as active_subscriptions_count
FROM plans p
LEFT JOIN mailing_limits ml ON p.id = ml.plan_id AND ml.deleted_at IS NULL
WHERE p.deleted_at IS NULL
ORDER BY p.price ASC;