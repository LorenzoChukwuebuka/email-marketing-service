CREATE TABLE IF NOT EXISTS contacts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    company_id UUID NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    from_origin VARCHAR(255) NOT NULL, -- "from" is a reserved keyword, consider renaming if possible
    is_subscribed BOOLEAN DEFAULT FALSE,
    user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    -- Corrected FOREIGN KEY definition
    FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS contact_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    company_id UUID NOT NULL,
    group_name VARCHAR(255) NOT NULL,
    user_id UUID NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    -- Corrected FOREIGN KEY definition
    FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE
);

-- Corrected user_contact_groups table definition
CREATE TABLE IF NOT EXISTS user_contact_groups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id UUID NOT NULL,
    -- Corrected data types to UUID
    contact_group_id UUID NOT NULL,
    contact_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    -- FOREIGN KEY constraints were already correct in terms of column references
    CONSTRAINT fk_ucg_group FOREIGN KEY (contact_group_id) REFERENCES contact_groups (id) ON DELETE CASCADE,
    CONSTRAINT fk_ucg_contact FOREIGN KEY (contact_id) REFERENCES contacts (id) ON DELETE CASCADE
);

-- Corrected contact_groups table definition
CREATE INDEX idx_contact_groups_id ON contact_groups (id);

CREATE INDEX idx_user_contact_groups_id ON user_contact_groups (id);

-- These indexes were already correct
CREATE INDEX idx_ucg_contact_group_id ON user_contact_groups (contact_group_id);

CREATE INDEX idx_ucg_contact_id ON user_contact_groups (contact_id);