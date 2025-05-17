

CREATE TABLE IF NOT EXISTS plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) DEFAULT 0.00,
    billing_cycle VARCHAR(50) DEFAULT 'monthly', -- Allowed: 'monthly', 'yearly' (or 'annually')
    status VARCHAR(100) DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), -- Note: To auto-update this on row changes, a trigger is needed in PostgreSQL
    deleted_at TIMESTAMPTZ
);

-- Plan Features
CREATE TABLE IF NOT EXISTS plan_features (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    plan_id UUID NOT NULL,
    name VARCHAR(255),
    description TEXT,
    value TEXT, -- Could be numeric, boolean, or text depending on the feature
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), -- Note: To auto-update this on row changes, a trigger is needed
    deleted_at TIMESTAMPTZ,
    CONSTRAINT fk_plan_features_plan_id FOREIGN KEY (plan_id) REFERENCES plans (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_plan_features_plan_id ON plan_features (plan_id);

-- Mailing Limits
CREATE TABLE IF NOT EXISTS mailing_limits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    plan_id UUID NOT NULL,
    daily_limit INTEGER DEFAULT 0 CHECK (daily_limit >= 0),
    monthly_limit INTEGER DEFAULT 0 CHECK (monthly_limit >= 0),
    max_recipients_per_mail INTEGER DEFAULT 0 CHECK (max_recipients_per_mail >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), -- Note: To auto-update this on row changes, a trigger is needed
    deleted_at TIMESTAMPTZ,
    CONSTRAINT uq_mailing_limits_plan_id UNIQUE (plan_id), -- Ensures one mailing_limit entry per plan
    CONSTRAINT fk_mailing_limits_plan_id FOREIGN KEY (plan_id) REFERENCES plans (id) ON DELETE CASCADE
);

