-- Migration: add_invitations_table.down.sql
DROP TRIGGER IF EXISTS update_invitations_updated_at ON invitations;
DROP FUNCTION IF EXISTS update_updated_at_column();
DROP TABLE IF EXISTS invitations;