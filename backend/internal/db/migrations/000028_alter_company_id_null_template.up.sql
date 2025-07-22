-- Make company_id nullable
ALTER TABLE templates 
ALTER COLUMN company_id DROP NOT NULL;

-- Add a constraint to ensure data integrity
-- Gallery templates must have NULL company_id
-- Non-gallery templates must have a valid company_id
ALTER TABLE templates 
ADD CONSTRAINT check_company_id_for_gallery 
CHECK (
    (is_gallery_template = TRUE AND company_id IS NULL) OR
    (is_gallery_template = FALSE AND company_id IS NOT NULL)
);
 