CREATE TABLE IF NOT EXISTS systems_smtp_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    txt_record TEXT,
    dmarc_record TEXT,
    dkim_selector TEXT,
    dkim_public_key TEXT,
    dkim_private_key TEXT,
    spf_record TEXT,
    verified BOOLEAN DEFAULT FALSE,
    mx_record TEXT,
    domain TEXT
);

 