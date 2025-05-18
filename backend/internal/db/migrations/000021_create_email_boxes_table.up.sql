CREATE TABLE IF NOT EXISTS email_boxes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_name TEXT,
    "from" TEXT,
    "to" TEXT,
    content BYTEA,
    mailbox TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);