-- Corrected templates table definition
CREATE TABLE IF NOT EXISTS templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id UUID NOT NULL,
    company_id UUID NOT NULL, -- company_id must be NOT NULL as it's a foreign key
    template_name VARCHAR(255) NOT NULL,
    sender_name VARCHAR(255),
    from_email VARCHAR(255),
    subject TEXT,
    type VARCHAR(20) NOT NULL CHECK (
        type IN ('transactional', 'marketing')
    ),
    email_html TEXT,
    email_design JSONB,
    is_editable BOOLEAN DEFAULT FALSE,
    is_published BOOLEAN DEFAULT FALSE,
    is_public_template BOOLEAN DEFAULT FALSE,
    is_gallery_template BOOLEAN DEFAULT FALSE,
    tags VARCHAR(255), -- Made tags nullable, as NOT NULL might be too restrictive if not always required
    description TEXT,
    image_url TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    editor_type VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    -- Corrected FOREIGN KEY definition for company_id
    FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE
);

-- Index on the primary key is often automatically created, but explicitly defining doesn't hurt
CREATE INDEX IF NOT EXISTS idx_templates_id ON templates (id);

CREATE INDEX IF NOT EXISTS idx_templates_user_id ON templates (user_id);

CREATE INDEX IF NOT EXISTS idx_templates_type ON templates(type);
