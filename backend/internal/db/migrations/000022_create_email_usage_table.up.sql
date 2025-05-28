-- Email Usage Tracking
CREATE TABLE IF NOT EXISTS email_usage (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    subscription_id UUID NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
    usage_period_start DATE NOT NULL, -- Start of the tracking period (daily/monthly)
    usage_period_end DATE NOT NULL,   -- End of the tracking period
    period_type VARCHAR(20) NOT NULL DEFAULT 'daily', -- 'daily' or 'monthly'
    emails_sent INTEGER DEFAULT 0,
    emails_limit INTEGER NOT NULL,    -- Limit for this period based on plan
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Ensure one record per company per period
    UNIQUE(company_id, usage_period_start, period_type)
);
