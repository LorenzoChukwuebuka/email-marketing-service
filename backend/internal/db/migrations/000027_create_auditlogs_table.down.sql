-- Drop indexes first (optional, but avoids errors in some versions of Postgres)
DROP INDEX IF EXISTS idx_audit_logs_user_id;
DROP INDEX IF EXISTS idx_audit_logs_resource;
DROP INDEX IF EXISTS idx_audit_logs_occurred_at;

-- Drop the table
DROP TABLE IF EXISTS audit_logs;

-- Drop the enum type
DROP TYPE IF EXISTS audit_action;
