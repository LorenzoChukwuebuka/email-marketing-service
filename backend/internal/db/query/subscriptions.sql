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