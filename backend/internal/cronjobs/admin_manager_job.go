package cronjobs

import (
	"context"
	"database/sql"
	db "email-marketing-service/internal/db/sqlc"
	"log"
)

type AdminJobManager struct {
	store db.Store
	ctx   context.Context
}

func NewAdminJobManager(store db.Store, ctx context.Context) *AdminJobManager {
	return &AdminJobManager{store: store, ctx: ctx}
}

func (a *AdminJobManager) EnableJob(jobName string) error {
	err := a.store.EnableJob(a.ctx, jobName)
	if err != nil {
		return err
	}
	log.Printf("Job %s has been enabled", jobName)
	return nil
}

func (a *AdminJobManager) DisableJob(jobName string) error {
	err := a.store.DisableJob(a.ctx, jobName)
	if err != nil {
		return err
	}
	log.Printf("Job %s has been disabled", jobName)
	return nil
}

func (a *AdminJobManager) UpdateJobSchedule(jobName, newSchedule, description string, timeoutSeconds, maxRetries *int32) error {
	params := db.UpdateJobScheduleParams{
		JobName:      jobName,
		CronSchedule: newSchedule,
	}

	if description != "" {
		params.Description = sql.NullString{String: description, Valid: true}
	}
	if timeoutSeconds != nil {
		params.TimeoutSeconds = sql.NullInt32{Int32: *timeoutSeconds, Valid: true}
	}
	if maxRetries != nil {
		params.MaxRetries = sql.NullInt32{Int32: *maxRetries, Valid: true}
	}

	err := a.store.UpdateJobSchedule(a.ctx, params)
	if err != nil {
		return err
	}
	log.Printf("Job %s schedule updated to: %s", jobName, newSchedule)
	return nil
}

func (a *AdminJobManager) GetJobHistory(jobName string, limit int32) ([]db.GetJobExecutionHistoryRow, error) {
	params := db.GetJobExecutionHistoryParams{
		JobName: jobName,
		Limit:   limit,
	}
	result, _ := a.store.GetJobExecutionHistory(a.ctx, params)

	return result, nil
}

func (a *AdminJobManager) GetJobStatus(jobName string) (*db.JobSchedule, error) {
	result, _ := a.store.GetJobScheduleByName(a.ctx, jobName)

	return &result, nil
}
