CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    actor_id UUID, -- Admin or user who performed the action
    actor_type VARCHAR(50) NOT NULL, -- "admin" or "user"
    action VARCHAR(255) NOT NULL, -- e.g. "BLOCK_USER"
    target_id UUID, -- Optional: user affected
    target_type VARCHAR(50), -- Optional: e.g. "user", "campaign"
    metadata JSONB, -- Optional: store extra data
    ip_address VARCHAR(64),
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
