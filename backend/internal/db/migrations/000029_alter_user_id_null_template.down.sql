-- DOWN MIGRATION: alter_user_id_null_template.down.sql
-- Remove the user_id constraint first
ALTER TABLE templates 
DROP CONSTRAINT IF EXISTS check_user_id_for_gallery;

-- Before making user_id NOT NULL again, handle existing NULL values
-- Option A: Delete gallery templates (if acceptable)
-- DELETE FROM templates WHERE is_gallery_template = TRUE AND user_id IS NULL;

-- Option B: Assign a default user_id (create a system user first if needed)
-- You'll need to replace 'your-system-user-uuid-here' with an actual admin user UUID
-- UPDATE templates 
-- SET user_id = 'your-system-user-uuid-here' 
-- WHERE user_id IS NULL;

-- Make user_id NOT NULL again
ALTER TABLE templates 
ALTER COLUMN user_id SET NOT NULL;