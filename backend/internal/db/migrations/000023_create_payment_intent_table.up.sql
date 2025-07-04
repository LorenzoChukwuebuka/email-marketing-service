-- Payment Intents
CREATE TABLE IF NOT EXISTS payment_intents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    subscription_id UUID REFERENCES subscriptions(id),
    payment_intent_id VARCHAR(255)  NOT NULL, -- External payment gateway intent ID
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(10) DEFAULT 'NGN',
    payment_method_types TEXT[], -- Array of allowed payment methods ['card', 'bank_transfer', etc.]
    status VARCHAR(50) DEFAULT 'requires_payment_method', -- requires_payment_method/requires_confirmation/requires_action/processing/succeeded/canceled
    client_secret VARCHAR(255), -- For frontend payment confirmation
    description TEXT,
    metadata JSONB, -- Flexible field for additional data
    automatic_payment_methods BOOLEAN DEFAULT true,
    receipt_email VARCHAR(255),
    setup_future_usage VARCHAR(50), -- 'on_session', 'off_session', null
    confirmation_method VARCHAR(50) DEFAULT 'automatic', -- automatic/manual
    capture_method VARCHAR(50) DEFAULT 'automatic', -- automatic/manual
    payment_method_id VARCHAR(255), -- Reference to attached payment method
    last_payment_error JSONB, -- Store error details if payment fails
    next_action JSONB, -- Instructions for additional authentication
    canceled_at TIMESTAMP,
    succeeded_at TIMESTAMP,
    expires_at TIMESTAMP, -- When the intent expires
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Indexes for better performance
CREATE INDEX IF NOT EXISTS idx_payment_intents_company_id ON payment_intents(company_id);
CREATE INDEX IF NOT EXISTS idx_payment_intents_user_id ON payment_intents(user_id);
CREATE INDEX IF NOT EXISTS idx_payment_intents_subscription_id ON payment_intents(subscription_id);
CREATE INDEX IF NOT EXISTS idx_payment_intents_status ON payment_intents(status);
CREATE INDEX IF NOT EXISTS idx_payment_intents_created_at ON payment_intents(created_at);
CREATE INDEX IF NOT EXISTS idx_payment_intents_payment_intent_id ON payment_intents(payment_intent_id);