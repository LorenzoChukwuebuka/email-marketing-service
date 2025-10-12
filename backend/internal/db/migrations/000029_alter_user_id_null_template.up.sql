-- Make user_id nullable for gallery templates
ALTER TABLE templates 
ALTER COLUMN user_id DROP NOT NULL;

-- Safely add constraint only if it doesn't exist
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 
        FROM pg_constraint
        WHERE conname = 'check_user_id_for_gallery'
          AND conrelid = 'templates'::regclass
    ) THEN
        ALTER TABLE templates
        ADD CONSTRAINT check_user_id_for_gallery
        CHECK (
            (is_gallery_template = TRUE AND user_id IS NULL) OR
            (is_gallery_template = FALSE AND user_id IS NOT NULL)
        );
    END IF;
END
$$;
