package cronjobs

import (
	"context"
	"database/sql"
	db "email-marketing-service/internal/db/sqlc"
	"fmt"
	"time"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

// JobManager handles job execution with database tracking
type JobManager struct {
	store db.Store
	ctx   context.Context
}

func NewJobManager(store db.Store, ctx context.Context) *JobManager {
	return &JobManager{
		store: store,
		ctx:   ctx,
	}
}

// ExecuteJob runs a job with full database tracking and error handling
func (jm *JobManager) ExecuteJob(job Job) {
	jobName := job.Name()
	startTime := time.Now()

	// Check if job is enabled
	schedule, err := jm.store.GetJobScheduleByName(jm.ctx, jobName)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Printf("Job %s not found in schedule table, skipping\n", jobName)
			// Optionally: create the job schedule entry automatically
			jm.ensureJobScheduleExists(job)
			return
		}
		fmt.Printf("Error checking job schedule for %s: %v\n", jobName, err)
		return
	}

	if !schedule.Enabled.Bool {
		fmt.Printf("Job %s is disabled, skipping\n", jobName)
		return
	}

	// Update last run time
	err = jm.store.UpdateJobLastRun(jm.ctx, db.UpdateJobLastRunParams{
		JobName:   jobName,
		LastRunAt: sql.NullTime{Time: startTime, Valid: true},
	})
	if err != nil {
		fmt.Printf("Warning: Failed to update last run time for job %s: %v\n", jobName, err)
	}

	// Create execution log entry with better error handling
	var logID uuid.UUID
	logID, err = jm.store.CreateJobExecutionLog(jm.ctx, db.CreateJobExecutionLogParams{
		JobScheduleID: uuid.NullUUID{UUID: schedule.ID, Valid: true},
		JobName:       jobName,
		StartedAt:     startTime,
		Status:        "running",
	})

	if err != nil {
		fmt.Printf("ERROR: Failed to create execution log for job %s: %v\n", jobName, err)
		// Continue execution but we'll have to handle the missing log entry
		logID = uuid.Nil
	} else {
		fmt.Printf("Created execution log with ID: %s for job: %s\n", logID.String(), jobName)
	}

	fmt.Printf("Starting job: %s (Type: %s)\n", jobName, job.Type())

	// Execute the job with timeout handling
	jobChan := make(chan error, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				jobChan <- fmt.Errorf("job panicked: %v", r)
			}
		}()
		jobChan <- job.Run()
	}()

	var jobErr error
	timeout := time.Duration(schedule.TimeoutSeconds.Int32) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Minute // Default timeout
	}

	select {
	case jobErr = <-jobChan:
		// Job completed normally
	case <-time.After(timeout):
		jobErr = fmt.Errorf("job timed out after %v", timeout)
	}

	finishTime := time.Now()
	duration := finishTime.Sub(startTime)

	// Determine status and handle result
	var status string
	if jobErr != nil {
		status = "failed"
		if jobErr.Error() == fmt.Sprintf("job timed out after %v", timeout) {
			status = "timeout"
		}

		fmt.Printf("Job %s failed: %v (duration: %v)\n", jobName, jobErr, duration)

		// Update failure count
		err = jm.store.UpdateJobFailure(jm.ctx, db.UpdateJobFailureParams{
			JobName:       jobName,
			LastFailureAt: sql.NullTime{Time: finishTime, Valid: true},
		})
		if err != nil {
			fmt.Printf("Warning: Failed to update job failure for %s: %v\n", jobName, err)
		}
	} else {
		status = "completed"
		fmt.Printf("Job %s completed successfully (duration: %v)\n", jobName, duration)

		// Update success time
		err = jm.store.UpdateJobSuccess(jm.ctx, db.UpdateJobSuccessParams{
			JobName:       jobName,
			LastSuccessAt: sql.NullTime{Time: finishTime, Valid: true},
		})
		if err != nil {
			fmt.Printf("Warning: Failed to update job success for %s: %v\n", jobName, err)
		}
	}

	// Update execution log - ALWAYS try to update, even if creation failed
	if logID != uuid.Nil {
		var errorMsg sql.NullString
		if jobErr != nil {
			errorMsg = sql.NullString{String: jobErr.Error(), Valid: true}
		}

		err = jm.store.UpdateJobExecutionLog(jm.ctx, db.UpdateJobExecutionLogParams{
			ID:           logID,
			FinishedAt:   sql.NullTime{Time: finishTime, Valid: true},
			Status:       status,
			DurationMs:   sql.NullInt32{Int32: int32(duration.Milliseconds()), Valid: true},
			ErrorMessage: errorMsg,
			OutputData:   pqtype.NullRawMessage{RawMessage: nil, Valid: false},
		})
		if err != nil {
			fmt.Printf("ERROR: Failed to update execution log for job %s: %v\n", jobName, err)
		} else {
			fmt.Printf("Updated execution log for job %s with status: %s\n", jobName, status)
		}
	} else {
		// Log creation failed, try to create a minimal log entry
		fmt.Printf("Attempting to create fallback execution log for job %s\n", jobName)
		jm.createFallbackExecutionLog(schedule.ID, jobName, startTime, finishTime, status, duration, jobErr)
	}
}

// ensureJobScheduleExists creates a job schedule entry if it doesn't exist
func (jm *JobManager) ensureJobScheduleExists(job Job) error {
	_, err := jm.store.CreateJobSchedule(jm.ctx, db.CreateJobScheduleParams{
		JobName:        job.Name(),
		CronSchedule:   job.Schedule(),
		Description:    sql.NullString{String: job.Description(), Valid: true},
		Enabled:        sql.NullBool{Bool: true, Valid: true},
		TimeoutSeconds: sql.NullInt32{Int32: 1800, Valid: true}, // 30 minutes default
		MaxRetries:     sql.NullInt32{Int32: 3, Valid: true},
	})

	if err != nil {
		fmt.Printf("Failed to create job schedule for %s: %v\n", job.Name(), err)
		return err
	}

	fmt.Printf("Created job schedule entry for %s\n", job.Name())
	return nil
}

// createFallbackExecutionLog attempts to create a complete execution log when the initial creation failed
func (jm *JobManager) createFallbackExecutionLog(scheduleID uuid.UUID, jobName string, startTime, finishTime time.Time, status string, duration time.Duration, jobErr error) {
	var errorMsg sql.NullString
	if jobErr != nil {
		errorMsg = sql.NullString{String: jobErr.Error(), Valid: true}
	}

	logID, err := jm.store.CreateJobExecutionLog(jm.ctx, db.CreateJobExecutionLogParams{
		JobScheduleID: uuid.NullUUID{UUID: scheduleID, Valid: true},
		JobName:       jobName,
		StartedAt:     startTime,
		Status:        "running", // Start with running status
	})

	if err != nil {
		fmt.Printf("ERROR: Fallback execution log creation also failed for job %s: %v\n", jobName, err)
		return
	}

	// Immediately update with completion data
	err = jm.store.UpdateJobExecutionLog(jm.ctx, db.UpdateJobExecutionLogParams{
		ID:           logID,
		FinishedAt:   sql.NullTime{Time: finishTime, Valid: true},
		Status:       status,
		DurationMs:   sql.NullInt32{Int32: int32(duration.Milliseconds()), Valid: true},
		ErrorMessage: errorMsg,
		OutputData:   pqtype.NullRawMessage{RawMessage: nil, Valid: false},
	})

	if err != nil {
		fmt.Printf("ERROR: Failed to update fallback execution log for job %s: %v\n", jobName, err)
	} else {
		fmt.Printf("Successfully created fallback execution log for job %s\n", jobName)
	}
}

// IsJobEnabled checks if a job is enabled without running it
func (jm *JobManager) IsJobEnabled(jobName string) bool {
	schedule, err := jm.store.GetJobScheduleByName(jm.ctx, jobName)
	if err != nil {
		return false
	}
	return schedule.Enabled.Bool
}
