-- name: CreatePayment :one
INSERT INTO
    payments (
        company_id,
        user_id,
        subscription_id,
        payment_id,
        amount,
        currency,
        payment_method,
        status,
        notes,
        transaction_reference,
        payment_date,
        billing_period_start,
        billing_period_end
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

-- name: GetPaymentByID :one
SELECT * FROM payments WHERE id = $1 AND deleted_at IS NULL;

-- name: GetPaymentByPaymentID :one
SELECT * FROM payments WHERE payment_id = $1 AND deleted_at IS NULL;

-- name: ListPaymentsByCompanyID :many
SELECT *
FROM payments
WHERE
    company_id = $1
    AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListPaymentsBySubscriptionID :many
SELECT *
FROM payments
WHERE
    subscription_id = $1
    AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: ListPaymentsByUserID :many
SELECT *
FROM payments
WHERE
    user_id = $1
    AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdatePaymentStatus :one
UPDATE payments
SET
    status = $1,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $2
    AND deleted_at IS NULL RETURNING *;

-- name: RecordRefund :one
UPDATE payments
SET
    refunded_amount = $1,
    refund_date = CURRENT_TIMESTAMP,
    status = 'refunded',
    notes = COALESCE(notes || E '\n', '') || 'Refund: ' || $2,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $3
    AND deleted_at IS NULL RETURNING *;

-- name: GetCompanyPaymentsWithSubscriptionInfo :many
SELECT p.*, s.plan_id, s.status as subscription_status, s.billing_cycle
FROM payments p
    JOIN subscriptions s ON p.subscription_id = s.id
WHERE
    p.company_id = $1
    AND p.deleted_at IS NULL
ORDER BY p.created_at DESC
LIMIT $2
OFFSET
    $3;

-- name: GetPaymentSuccessRate :one
SELECT 
    COUNT(*) FILTER (WHERE status = 'successful') AS successful_count,
    COUNT(*) AS total_count,
    CASE 
        WHEN COUNT(*) = 0 THEN 0
        ELSE COUNT(*) FILTER (WHERE status = 'successful')::DECIMAL / COUNT(*)::DECIMAL * 100
    END AS success_rate
FROM payments
WHERE created_at >= $1 AND deleted_at IS NULL;

-- name: GetTotalPaymentAmount :one
SELECT
    COALESCE(SUM(amount), 0) as total_amount,
    COALESCE(
        SUM(amount) FILTER (
            WHERE
                status = 'successful'
        ),
        0
    ) as successful_amount
FROM payments
WHERE
    company_id = $1
    AND deleted_at IS NULL;

/* -- name: GetRecentPaymentActivity :many
SELECT p.*, c.name as companyname
FROM payments p
    JOIN companies c ON p.company_id = c.id
WHERE
    p.deleted_at IS NULL
ORDER BY p.created_at DESC
LIMIT $1; */

-- name: GetPaymentsByDateRange :many
SELECT *
FROM payments
WHERE
    created_at BETWEEN $1 AND $2
    AND deleted_at IS NULL
ORDER BY created_at DESC;