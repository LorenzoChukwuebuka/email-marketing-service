CREATE TABLE IF NOT EXISTS companies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    companyname TEXT DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ
);


CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    fullname TEXT NOT NULL,
    company_id UUID NOT NULL,
    email TEXT NOT NULL UNIQUE,
    phonenumber VARCHAR(255) DEFAULT NULL,
    password TEXT DEFAULT NULL,
    google_id TEXT DEFAULT NULL,
    picture TEXT DEFAULT NULL,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    blocked BOOLEAN NOT NULL DEFAULT FALSE,
    verified_at TIMESTAMP DEFAULT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    scheduled_for_deletion BOOLEAN NOT NULL DEFAULT FALSE,
    scheduled_deletion_at TIMESTAMP,
    last_login_at TIMESTAMP,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (company_id) REFERENCES companies (id) ON DELETE CASCADE
);

 


