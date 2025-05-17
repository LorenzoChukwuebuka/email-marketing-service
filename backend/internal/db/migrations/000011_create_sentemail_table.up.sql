
-- Corrected sent_emails table definition
CREATE TABLE IF NOT EXISTS sent_emails (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    company_id UUID NOT NULL,
    -- Corrected FOREIGN KEY to reference 'senders' table based on context
    -- If it should reference a 'users' table, change 'senders' back to 'users'
    sender_id UUID NOT NULL REFERENCES senders (id), -- Renamed column to sender_id for clarity
    recipient VARCHAR(255) NOT NULL,
    message_content TEXT NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (
        status IN (
            'sending',
            'failed',
            'success'
        )
    ),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NULL, -- updated_at is typically nullable as it's not set on creation
    -- Corrected FOREIGN KEY definition for company_id
    FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE
);

-- Index on the primary key is often automatically created, but explicitly defining doesn't hurt
CREATE INDEX IF NOT EXISTS idx_sent_emails_id ON sent_emails (id);

-- Corrected index name to match the column name sender_id
CREATE INDEX IF NOT EXISTS idx_sent_emails_sender_id ON sent_emails (sender_id);
