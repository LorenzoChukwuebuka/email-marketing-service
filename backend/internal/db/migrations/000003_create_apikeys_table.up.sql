
CREATE TABLE IF NOT EXISTS api_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    company_id UUID NOT NULL,
    name TEXT NOT NULL,
    api_key TEXT NOT NULL UNIQUE, -- Added UNIQUE constraint directly here for api_key
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE 
);


 CREATE INDEX IF NOT EXISTS idx_api_keys_user_id ON api_keys(user_id);
 CREATE INDEX IF NOT EXISTS idx_api_keys_company_id ON api_keys(company_id);

CREATE TABLE IF NOT EXISTS smtp_master_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    company_id UUID NOT NULL,
    smtp_login TEXT NOT NULL, 
    key_name TEXT NOT NULL,
    password TEXT NOT NULL, 
    status TEXT NOT NULL DEFAULT 'active', 
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE 
);

 CREATE INDEX IF NOT EXISTS idx_smtp_master_keys_user_id ON smtp_master_keys(user_id);
 CREATE INDEX IF NOT EXISTS idx_smtp_master_keys_company_id ON smtp_master_keys(company_id);

CREATE TABLE IF NOT EXISTS smtp_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL,
    user_id UUID NOT NULL,
    key_name TEXT NOT NULL,
    password TEXT NOT NULL,  
    status TEXT NOT NULL DEFAULT 'active',  
    smtp_login TEXT NOT NULL,  
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE  
);

-- Indexes for smtp_keys
CREATE INDEX IF NOT EXISTS idx_smtp_keys_key_name ON smtp_keys(key_name);
 CREATE INDEX IF NOT EXISTS idx_smtp_keys_user_id ON smtp_keys(user_id);
 CREATE INDEX IF NOT EXISTS idx_smtp_keys_company_id ON smtp_keys(company_id);
 CREATE INDEX IF NOT EXISTS idx_smtp_keys_smtp_login ON smtp_keys(smtp_login);

