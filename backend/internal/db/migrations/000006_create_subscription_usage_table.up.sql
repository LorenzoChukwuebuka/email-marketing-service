CREATE TABLE If not EXISTS subscription_usages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id UUID NOT NULL,
    company_id UUID NOT NULL,
    subscription_id UUID NOT NULL,
    feature_name VARCHAR(255) NOT NULL, -- e.g., 'positions'
    used_count INTEGER DEFAULT 0,
    max_allowed INTEGER DEFAULT 0,
    reset_at TIMESTAMP NOT NULL, -- When the counter resets (next billing cycle)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,

-- Foreign key constraints (adjust table names as needed)
CONSTRAINT fk_subscription_usages_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_subscription_usages_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    CONSTRAINT fk_subscription_usages_subscription FOREIGN KEY (subscription_id) REFERENCES subscriptions(id) ON DELETE CASCADE
);