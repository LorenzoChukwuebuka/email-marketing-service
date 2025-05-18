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