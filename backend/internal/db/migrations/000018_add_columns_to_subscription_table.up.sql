-- First, create the trigger function
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Subscription table alterations
ALTER TABLE IF EXISTS subscriptions
ADD COLUMN IF NOT EXISTS next_billing_date TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS auto_renew BOOLEAN,
ADD COLUMN IF NOT EXISTS cancellation_reason TEXT,
ADD COLUMN IF NOT EXISTS last_payment_date TIMESTAMPTZ;

-- Payment table alterations
ALTER TABLE IF EXISTS payments
ADD COLUMN IF NOT EXISTS transaction_reference VARCHAR(255),
ADD COLUMN IF NOT EXISTS payment_date TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS billing_period_start TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS billing_period_end TIMESTAMPTZ,
ADD COLUMN IF NOT EXISTS refunded_amount DECIMAL(10, 2) DEFAULT 0.00,
ADD COLUMN IF NOT EXISTS integrity_hash TEXT,
ADD COLUMN IF NOT EXISTS refund_date TIMESTAMPTZ;

-- Drop existing triggers if they exist before recreating (to avoid duplicate trigger errors)
DROP TRIGGER IF EXISTS set_subscriptions_timestamp ON subscriptions;

DROP TRIGGER IF EXISTS set_payments_timestamp ON payments;

-- Add trigger for auto-updating 'updated_at' on your existing tables
CREATE TRIGGER set_subscriptions_timestamp
BEFORE UPDATE ON subscriptions
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_payments_timestamp
BEFORE UPDATE ON payments
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();