CREATE TABLE IF NOT EXISTS campaign_errors (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    campaign_id UUID NOT NULL REFERENCES campaigns (id) ON DELETE CASCADE,
    error_type VARCHAR(100),
    error_message TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);