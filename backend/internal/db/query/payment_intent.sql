-- name: CreatePaymentIntent :one
INSERT INTO
    payment_intents (
        company_id,
        user_id,
        subscription_id,
        payment_intent_id,
        amount,
        currency,
        payment_method_types,
        status,
        client_secret,
        description,
        metadata,
        automatic_payment_methods,
        receipt_email,
        setup_future_usage,
        confirmation_method,
        capture_method,
        payment_method_id,
        last_payment_error,
        next_action,
        canceled_at,
        succeeded_at,
        expires_at
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
        $18,
        $19,
        $20,
        $21,
        $22
    ) RETURNING *;

-- name: GetPaymentIntent :one
SELECT * FROM payment_intents WHERE id = $1 AND deleted_at IS NULL;

-- name: GetPaymentIntentByPaymentIntentID :one
SELECT *
FROM payment_intents
WHERE
    payment_intent_id = $1
    AND deleted_at IS NULL;

-- name: GetPaymentIntentsByCompanyID :many
SELECT *
FROM payment_intents
WHERE
    company_id = $1
    AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: GetPaymentIntentsByUserID :many
SELECT *
FROM payment_intents
WHERE
    user_id = $1
    AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: GetPaymentIntentsBySubscriptionID :many
SELECT *
FROM payment_intents
WHERE
    subscription_id = $1
    AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: GetPaymentIntentsByStatus :many
SELECT *
FROM payment_intents
WHERE
    company_id = $1
    AND status = $2
    AND deleted_at IS NULL
ORDER BY created_at DESC;

-- name: UpdatePaymentIntent :one
UPDATE payment_intents
SET
    subscription_id = COALESCE(
        sqlc.narg ('subscription_id'),
        subscription_id
    ),
    amount = COALESCE(sqlc.narg ('amount'), amount),
    currency = COALESCE(
        sqlc.narg ('currency'),
        currency
    ),
    payment_method_types = COALESCE(
        sqlc.narg ('payment_method_types'),
        payment_method_types
    ),
    payment_intent_id = COALESCE(
        sqlc.narg ('payment_intent_id'),
        payment_intent_id
    ),
    status = COALESCE(sqlc.narg ('status'), status),
    client_secret = COALESCE(
        sqlc.narg ('client_secret'),
        client_secret
    ),
    description = COALESCE(
        sqlc.narg ('description'),
        description
    ),
    metadata = COALESCE(
        sqlc.narg ('metadata'),
        metadata
    ),
    automatic_payment_methods = COALESCE(
        sqlc.narg ('automatic_payment_methods'),
        automatic_payment_methods
    ),
    receipt_email = COALESCE(
        sqlc.narg ('receipt_email'),
        receipt_email
    ),
    setup_future_usage = COALESCE(
        sqlc.narg ('setup_future_usage'),
        setup_future_usage
    ),
    confirmation_method = COALESCE(
        sqlc.narg ('confirmation_method'),
        confirmation_method
    ),
    capture_method = COALESCE(
        sqlc.narg ('capture_method'),
        capture_method
    ),
    payment_method_id = COALESCE(
        sqlc.narg ('payment_method_id'),
        payment_method_id
    ),
    last_payment_error = COALESCE(
        sqlc.narg ('last_payment_error'),
        last_payment_error
    ),
    next_action = COALESCE(
        sqlc.narg ('next_action'),
        next_action
    ),
    canceled_at = COALESCE(
        sqlc.narg ('canceled_at'),
        canceled_at
    ),
    succeeded_at = COALESCE(
        sqlc.narg ('succeeded_at'),
        succeeded_at
    ),
    expires_at = COALESCE(
        sqlc.narg ('expires_at'),
        expires_at
    ),
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = sqlc.arg ('id')
    AND deleted_at IS NULL RETURNING *;

-- name: UpdatePaymentIntentError :one
UPDATE payment_intents
SET
    last_payment_error = $1,
    status = $2,
    updated_at = CURRENT_TIMESTAMP
WHERE
    payment_intent_id = $3
    AND deleted_at IS NULL RETURNING *;

