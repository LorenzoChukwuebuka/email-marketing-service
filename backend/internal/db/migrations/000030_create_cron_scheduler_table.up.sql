CREATE TABLE job_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_name VARCHAR(100) NOT NULL UNIQUE,
    job_type VARCHAR(50) NOT NULL,
    cron_schedule VARCHAR(50) NOT NULL,
    enabled BOOLEAN DEFAULT true,
    description TEXT,
    timeout_seconds INTEGER DEFAULT 300,
    max_retries INTEGER DEFAULT 3,
    last_run_at TIMESTAMP,
    last_success_at TIMESTAMP,
    last_failure_at TIMESTAMP,
    failure_count INTEGER DEFAULT 0,
    total_runs INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Table to log job executions
CREATE TABLE job_execution_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_schedule_id UUID REFERENCES job_schedules(id) ON DELETE CASCADE,
    job_name VARCHAR(100) NOT NULL,
    started_at TIMESTAMP NOT NULL,
    finished_at TIMESTAMP,
    status VARCHAR(20) NOT NULL CHECK (status IN ('running', 'completed', 'failed', 'timeout')),
    duration_ms INTEGER,
    error_message TEXT,
    output_data JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_job_schedules_enabled ON job_schedules(enabled);
CREATE INDEX idx_job_schedules_next_run ON job_schedules(last_run_at) WHERE enabled = true;
CREATE INDEX idx_job_execution_logs_job_schedule ON job_execution_logs(job_schedule_id);
CREATE INDEX idx_job_execution_logs_started_at ON job_execution_logs(started_at);

-- Insert initial job configurations
INSERT INTO job_schedules (job_name, job_type, cron_schedule, description, enabled) VALUES 
('auto_close_support_tickets', 'AutoCloseSupportTicket', '0 0 0 * * *', 'Automatically close stale support tickets that have no replies for 48+ hours', true),
('update_expired_subscriptions', 'UpdateExpiredSubscription', '0 0 2 * * *', 'Update expired subscription statuses', true);

-- Function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Trigger to automatically update updated_at
CREATE TRIGGER update_job_schedules_updated_at 
    BEFORE UPDATE ON job_schedules 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();