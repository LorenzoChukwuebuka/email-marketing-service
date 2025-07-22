-- UP MIGRATION: alter_user_id_null_template.up.sql
-- Make user_id nullable for gallery templates
ALTER TABLE templates 
ALTER COLUMN user_id DROP NOT NULL;

-- Add constraint to ensure data integrity for user_id
-- Gallery templates must have NULL user_id
-- Non-gallery templates must have a valid user_id
ALTER TABLE templates 
ADD CONSTRAINT check_user_id_for_gallery 
CHECK (
    (is_gallery_template = TRUE AND user_id IS NULL) OR
    (is_gallery_template = FALSE AND user_id IS NOT NULL)
);
