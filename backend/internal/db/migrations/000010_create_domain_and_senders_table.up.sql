
CREATE TABLE IF NOT EXISTS domains (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id UUID NOT NULL,
    company_id UUID NOT NULL, -- company_id must be NOT NULL as it's a foreign key
    domain VARCHAR(255) NOT NULL,
    txt_record VARCHAR(255),
    dmarc_record VARCHAR(255),
    dkim_selector VARCHAR(255),
    dkim_public_key TEXT,
    dkim_private_key TEXT,
    spf_record VARCHAR(255),
    verified BOOLEAN DEFAULT FALSE,
    mx_record TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    -- Corrected FOREIGN KEY definition
    FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE
);

-- Index on the primary key is often automatically created, but explicitly defining doesn't hurt
CREATE INDEX idx_domains_id ON domains (id);

CREATE INDEX idx_domains_user_id ON domains (user_id);

CREATE INDEX idx_domains_domain ON domains (domain);


-- Corrected senders table definition
CREATE TABLE IF NOT EXISTS senders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    user_id UUID NOT NULL,
    -- Corrected company_id definition and moved FOREIGN KEY
    company_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    verified BOOLEAN DEFAULT FALSE,
    is_signed BOOLEAN DEFAULT FALSE,
    -- domain_id definition is correct as an inline foreign key
    domain_id UUID NOT NULL REFERENCES domains (id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    -- Added FOREIGN KEY constraint for company_id
    FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE
);

-- Index on the primary key is often automatically created, but explicitly defining doesn't hurt
CREATE INDEX idx_senders_id ON senders (id);

CREATE INDEX idx_senders_user_id ON senders (user_id);

CREATE INDEX idx_senders_domain_id ON senders (domain_id);
