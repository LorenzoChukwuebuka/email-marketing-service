-- Optional: Enum for consistent action names
CREATE TYPE audit_action AS ENUM ('CREATE', 'UPDATE', 'DELETE',  'LOGIN',
    'LOGOUT','LOGIN_FAILED');

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id UUID NOT NULL,
    action audit_action NOT NULL, -- Enum ensures only valid actions
    resource VARCHAR(50) NOT NULL, -- E.g. "User", "Invoice"
    resource_id UUID, -- Specific record ID
    method VARCHAR(10), -- HTTP method if applicable
    endpoint VARCHAR(255), -- API endpoint or route
    ip_address INET, -- Client IP
    occurred_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    success BOOLEAN DEFAULT TRUE,
    request_body JSONB,
    changes JSONB
);

-- Indexes for faster querying
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs (user_id);

CREATE INDEX IF NOT EXISTS idx_audit_logs_resource ON audit_logs (resource, resource_id);

CREATE INDEX IF NOT EXISTS idx_audit_logs_occurred_at ON audit_logs (occurred_at);