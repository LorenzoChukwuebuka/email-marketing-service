
-- Down migration for subscription table alterations
ALTER TABLE subscriptions
    DROP COLUMN IF EXISTS next_billing_date,
    DROP COLUMN IF EXISTS auto_renew,
    DROP COLUMN IF EXISTS cancellation_reason,
    DROP COLUMN IF EXISTS last_payment_date;

-- Down migration for payment table alterations
ALTER TABLE payments
    DROP COLUMN IF EXISTS transaction_reference,
    DROP COLUMN IF EXISTS payment_date,
    DROP COLUMN IF EXISTS billing_period_start,
    DROP COLUMN IF EXISTS billing_period_end,
    DROP COLUMN IF EXISTS refunded_amount,
    DROP COLUMN IF EXISTS refund_date,
    DROP COLUMN IF EXISTS integrity_hash;

-- Drop triggers
DROP TRIGGER IF EXISTS set_subscriptions_timestamp ON subscriptions;
DROP TRIGGER IF EXISTS set_payments_timestamp ON payments;