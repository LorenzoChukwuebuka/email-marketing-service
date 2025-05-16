-- Campaigns
CREATE TABLE IF NOT EXISTS campaigns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    company_id UUID NOT NULL REFERENCES companies (id) ON DELETE CASCADE, -- Added foreign key constraint
    name VARCHAR(255) NOT NULL,
    subject TEXT,
    preview_text TEXT,
    user_id UUID NOT NULL,
    sender_from_name TEXT,
    template_id UUID,
    sent_template_id UUID,
    recipient_info TEXT,
    is_published BOOLEAN DEFAULT FALSE,
    status VARCHAR(20) DEFAULT 'draft',
    track_type VARCHAR(50),
    is_archived BOOLEAN DEFAULT FALSE,
    sent_at TIMESTAMP,
    sender TEXT,
    scheduled_at TIMESTAMP,
    has_custom_logo BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Campaign Groups (linking campaigns to contact groups)
CREATE TABLE IF NOT EXISTS campaign_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    campaign_id UUID NOT NULL REFERENCES campaigns (id) ON DELETE CASCADE, -- Added/Ensured foreign key constraint
    contact_group_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Email Campaign Results
CREATE TABLE IF NOT EXISTS email_campaign_results (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    company_id UUID NOT NULL REFERENCES companies (id) ON DELETE CASCADE, -- Added foreign key constraint
    campaign_id UUID NOT NULL REFERENCES campaigns (id) ON DELETE CASCADE, -- Added foreign key constraint
    recipient_email VARCHAR(255) NOT NULL,
    recipient_name VARCHAR(255),
    version VARCHAR(10),
    sent_at TIMESTAMP,
    opened_at TIMESTAMP,
    open_count INTEGER DEFAULT 0,
    clicked_at TIMESTAMP,
    click_count INTEGER DEFAULT 0,
    conversion_at TIMESTAMP,
    bounce_status VARCHAR(20),
    unsubscribed_at TIMESTAMP,
    complaint_status BOOLEAN DEFAULT FALSE,
    device_type VARCHAR(50),
    location VARCHAR(100),
    retry_count INTEGER DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Indexes
-- Safe re-creation of indexes
DROP INDEX IF EXISTS idx_campaigns_status;

CREATE INDEX idx_campaigns_status ON campaigns (status);

DROP INDEX IF EXISTS idx_campaigns_scheduled_at;

CREATE INDEX idx_campaigns_scheduled_at ON campaigns (scheduled_at);

DROP INDEX IF EXISTS idx_email_results_campaign_id;

CREATE INDEX idx_email_results_campaign_id ON email_campaign_results (campaign_id);

DROP INDEX IF EXISTS idx_email_results_recipient_email;

CREATE INDEX idx_email_results_recipient_email ON email_campaign_results (recipient_email);