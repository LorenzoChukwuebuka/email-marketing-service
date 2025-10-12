-- Make company_id nullable
ALTER TABLE templates 
ALTER COLUMN company_id DROP NOT NULL;

-- Safely add constraint only if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM pg_constraint
        WHERE conname = 'check_company_id_for_gallery'
          AND conrelid = 'templates'::regclass
    ) THEN
        ALTER TABLE templates
        ADD CONSTRAINT check_company_id_for_gallery
        CHECK (
            (is_gallery_template = TRUE AND company_id IS NULL) OR
            (is_gallery_template = FALSE AND company_id IS NOT NULL)
        );
    END IF;
END
$$;
