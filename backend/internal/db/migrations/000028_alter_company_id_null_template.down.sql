-- Remove the constraint first
ALTER TABLE templates 
DROP CONSTRAINT IF EXISTS check_company_id_for_gallery;

-- Before making company_id NOT NULL again, we need to handle existing NULL values
-- Option A: Delete gallery templates (if acceptable)
-- DELETE FROM templates WHERE is_gallery_template = TRUE AND company_id IS NULL;

-- Option B: Assign a default company_id (create a system company first if needed)
-- UPDATE templates 
-- SET company_id = 'your-system-company-uuid-here' 
-- WHERE company_id IS NULL;

-- Make company_id NOT NULL again
ALTER TABLE templates 
ALTER COLUMN company_id SET NOT NULL;