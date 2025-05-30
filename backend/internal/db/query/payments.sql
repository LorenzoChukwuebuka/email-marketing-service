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

-- name: UpdatePaymentHash :exec
UPDATE payments
SET
    integrity_hash = $1,
    updated_at = now()
WHERE
    id = $2;

-- name: GetLastPaymentByCompanyID :one
SELECT *
FROM payments
WHERE
    company_id = $1
    AND deleted_at IS NULL
ORDER BY payment_date DESC, created_at DESC
LIMIT 1;

-- name: CheckPaymentIntentExists :one
SELECT EXISTS (
        SELECT 1
        FROM payments
        WHERE
            payment_id = $1
            AND deleted_at IS NULL
    ) AS exists;

-- name: GetPaymentsByCompanyAndUser :many
SELECT
    p.id,
    p.company_id,
    p.user_id,
    p.subscription_id,
    p.payment_id,
    p.amount,
    p.currency,
    p.payment_method,
    p.status,
    p.notes,
    p.created_at,
    p.updated_at,
    p.deleted_at,
    -- Company details
    c.companyname as companyname,
    c.created_at as companycreatedat,
    c.updated_at as companyupdatedat,
    -- User details
    u.fullname as userfullname,
    u.email as useremail,
    u.phonenumber as userphonenumber,
    u.picture as userpicture,
    u.verified as userverified,
    u.blocked as userblocked,
    u.status as userstatus,
    u.last_login_at as userlastloginat,
    u.created_at as usercreatedat,
    -- Subscription details
    s.plan_id as subscriptionplanid,
    s.amount as subscriptionamount,
    s.billing_cycle as subscriptionbillingcycle,
    s.trial_starts_at as subscriptiontrialstartsat,
    s.trial_ends_at as subscriptiontrialendsat,
    s.starts_at as subscriptionstartsat,
    s.ends_at as subscriptionendsat,
    s.status as subscriptionstatus,
    s.created_at as subscriptioncreatedat
FROM
    payments p
    INNER JOIN companies c ON p.company_id = c.id
    INNER JOIN users u ON p.user_id = u.id
    INNER JOIN subscriptions s ON p.subscription_id = s.id
WHERE
    p.company_id = $1
    AND p.user_id = $2
    AND p.deleted_at IS NULL
    AND c.deleted_at IS NULL
    AND u.deleted_at IS NULL
    AND s.deleted_at IS NULL
ORDER BY p.created_at DESC
LIMIT $3
OFFSET
    $4;

-- name: GetPaymentsByCompanyAndUserSimple :many
SELECT
    p.*,
    c.companyname,
    u.fullname,
    u.email,
    s.billing_cycle,
    s.status as subscription_status
FROM
    payments p
    LEFT JOIN companies c ON p.company_id = c.id
    LEFT JOIN users u ON p.user_id = u.id
    LEFT JOIN subscriptions s ON p.subscription_id = s.id
WHERE
    p.company_id = $1
    -- AND p.user_id = $2
    AND p.deleted_at IS NULL
ORDER BY p.created_at DESC;

-- name: GetPaymentCounts :one
SELECT COUNT(*) 
FROM payments
WHERE company_id = $1 
AND deleted_at IS NULL;