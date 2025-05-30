-- name: CreateSubscription :one
INSERT INTO
    subscriptions (
        company_id,
        plan_id,
        amount,
        billing_cycle,
        trial_starts_at,
        trial_ends_at,
        starts_at,
        ends_at,
        status,
        next_billing_date,
        auto_renew
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
        $11
    ) RETURNING *;

-- name: GetSubscriptionByID :one
SELECT * FROM subscriptions WHERE id = $1 AND deleted_at IS NULL;

-- name: GetActiveSubscriptionByCompanyID :one
SELECT *
FROM subscriptions
WHERE
    company_id = $1
    AND status = 'active'
    AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT 1;

-- name: ListSubscriptionsByCompanyID :many
SELECT *
FROM subscriptions
WHERE
    company_id = $1
    AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: GetCurrentRunningSubscription :one
SELECT
    s.id AS subscription_id,
    s.company_id,
    s.amount AS subscription_amount,
    s.billing_cycle,
    s.trial_starts_at,
    s.trial_ends_at,
    s.starts_at,
    s.ends_at,
    s.status AS subscription_status,
    s.created_at AS subscription_created_at,
    s.updated_at AS subscription_updated_at,
    p.id AS plan_id,
    p.name AS plan_name,
    p.description AS plan_description,
    p.price AS plan_price,
    p.billing_cycle AS plan_billing_cycle,
    p.status AS plan_status
FROM subscriptions s
    JOIN plans p ON s.plan_id = p.id
WHERE
    s.deleted_at IS NULL
    AND s.company_id = $1
ORDER BY s.created_at DESC
LIMIT 1;

-- name: UpdateSubscriptionStatus :one
UPDATE subscriptions
SET
    status = $2,
    updated_at = NOW()
WHERE
    id = $1
    AND deleted_at IS NULL RETURNING *;

-- name: GetAllActiveSubscriptions :many
SELECT *
FROM subscriptions
WHERE
    status = 'active'
    AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: GetActiveSubscriptionsEndingInDays :many
-- Find all active subscriptions that end within a specific number of days
SELECT *
FROM subscriptions
WHERE
    status = 'active'
    AND deleted_at IS NULL
    AND ends_at IS NOT NULL
    AND ends_at > NOW()
    AND ends_at <= NOW() + INTERVAL '%d days'
ORDER BY ends_at ASC;

-- name: GetActiveSubscriptionsEndingIn5Days :many
-- Find all active subscriptions that end within 5 days
SELECT *
FROM subscriptions
WHERE
    status = 'active'
    AND deleted_at IS NULL
    AND ends_at IS NOT NULL
    AND ends_at > NOW()
    AND ends_at <= NOW() + INTERVAL '5 days'
ORDER BY ends_at ASC;

-- name: GetActiveSubscriptionsNotExpired :many
-- Find all active subscriptions where ends_at is greater than current time
SELECT *
FROM subscriptions
WHERE
    status = 'active'
    AND deleted_at IS NULL
    AND (ends_at IS NULL OR ends_at > NOW())
ORDER BY created_at DESC;

-- name: GetExpiredActiveSubscriptions :many
-- Find subscriptions that are marked as 'active' but have actually expired
SELECT *
FROM subscriptions
WHERE
    status = 'active'
    AND deleted_at IS NULL
    AND ends_at IS NOT NULL
    AND ends_at <= NOW()
ORDER BY ends_at DESC;