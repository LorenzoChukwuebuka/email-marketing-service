-- name: CreateJobSchedule :one
INSERT INTO job_schedules (
    job_name,
    job_type,
    cron_schedule,
    enabled,
    description,
    timeout_seconds,
    max_retries,
    failure_count,
    total_runs,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, 0, 0, NOW(), NOW()
) RETURNING id;

-- name: GetEnabledJobSchedules :many
SELECT 
    id,
    job_name,
    job_type,
    cron_schedule,
    enabled,
    description,
    timeout_seconds,
    max_retries,
    last_run_at,
    last_success_at,
    last_failure_at,
    failure_count,
    total_runs,
    created_at,
    updated_at
FROM job_schedules 
WHERE enabled = true
ORDER BY job_name;

-- name: GetJobScheduleByName :one
SELECT 
    id,
    job_name,
    job_type,
    cron_schedule,
    enabled,
    description,
    timeout_seconds,
    max_retries,
    last_run_at,
    last_success_at,
    last_failure_at,
    failure_count,
    total_runs,
    created_at,
    updated_at
FROM job_schedules 
WHERE job_name = $1;

-- name: UpdateJobLastRun :exec
UPDATE job_schedules 
SET 
    last_run_at = $2,
    total_runs = total_runs + 1,
    updated_at = NOW()
WHERE job_name = $1;

-- name: UpdateJobSuccess :exec
UPDATE job_schedules 
SET 
    last_success_at = $2,
    failure_count = 0,
    updated_at = NOW()
WHERE job_name = $1;

-- name: UpdateJobFailure :exec
UPDATE job_schedules 
SET 
    last_failure_at = $2,
    failure_count = failure_count + 1,
    updated_at = NOW()
WHERE job_name = $1;

-- name: CreateJobExecutionLog :one
INSERT INTO job_execution_logs (
    job_schedule_id,
    job_name,
    started_at,
    status
) VALUES (
    $1, $2, $3, $4
) RETURNING id;

-- name: UpdateJobExecutionLog :exec
UPDATE job_execution_logs 
SET 
    finished_at = $2,
    status = $3,
    duration_ms = $4,
    error_message = $5,
    output_data = $6
WHERE id = $1;

-- name: GetJobExecutionHistory :many
SELECT 
    id,
    job_name,
    started_at,
    finished_at,
    status,
    duration_ms,
    error_message,
    output_data,
    created_at
FROM job_execution_logs 
WHERE job_name = $1 
ORDER BY started_at DESC 
LIMIT $2;

-- name: EnableJob :exec
UPDATE job_schedules 
SET 
    enabled = true,
    updated_at = NOW()
WHERE job_name = $1;

-- name: DisableJob :exec
UPDATE job_schedules 
SET 
    enabled = false,
    updated_at = NOW()
WHERE job_name = $1;

-- name: UpdateJobSchedule :exec
UPDATE job_schedules 
SET 
    cron_schedule = $2,
    description = COALESCE($3, description),
    timeout_seconds = COALESCE($4, timeout_seconds),
    max_retries = COALESCE($5, max_retries),
    updated_at = NOW()
WHERE job_name = $1;