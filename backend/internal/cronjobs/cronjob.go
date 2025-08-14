package cronjobs

import (
	"context"
	"database/sql"
	db "email-marketing-service/internal/db/sqlc"
	"github.com/robfig/cron/v3"
	"log"
)

type JobFactory struct {
	DB         db.Store
	ctx        context.Context
	jobManager *JobManager
}

func NewJobFactory(dbStore db.Store, ctx context.Context) *JobFactory {
	return &JobFactory{
		DB:         dbStore,
		ctx:        ctx,
		jobManager: NewJobManager(dbStore, ctx),
	}

}

func (f *JobFactory) CreateJobs() []Job {
	jobs := []Job{
		NewUpdateExpiredSubscriptionJob(f.DB, f.ctx),
		NewAutoCloseSupportTicket(f.DB, f.ctx),
		NewDeleteScheduledUsers(f.DB, f.ctx),
	}

	// Set store for all jobs (in case they need it)
	for _, job := range jobs {
		job.SetStore(f.DB)
	}

	return jobs
}

// SetupCronJobs creates and configures the cron scheduler with database integration
func SetupCronJobs(dbStore db.Store) *cron.Cron {
	c := cron.New(cron.WithSeconds(), cron.WithChain(
		cron.SkipIfStillRunning(cron.DefaultLogger),
		cron.Recover(cron.DefaultLogger),
	))

	ctx := context.Background()
	factory := NewJobFactory(dbStore, ctx)
	if err := factory.InitializeJobSchedules(); err != nil {
		log.Printf("Error initializing job schedules: %v", err)
	}

	jobs := factory.CreateJobs()

	for _, job := range jobs {
		// Use the job manager's execute method instead of direct job.Run()
		_, err := c.AddFunc(job.Schedule(), wrapJobWithManager(job, factory.jobManager))
		if err != nil {
			log.Printf("Error adding job %s to cron: %v", job.Name(), err)
		} else {
			log.Printf("Successfully scheduled job: %s (%s) with schedule: %s",
				job.Name(), job.Type(), job.Schedule())
		}
	}

	return c
}

// wrapJobWithManager wraps job execution with the JobManager for database tracking
func wrapJobWithManager(job Job, manager *JobManager) func() {
	return func() {
		manager.ExecuteJob(job)
	}
}

// SyncJobSchedulesFromDB syncs job schedules from the database
// This allows runtime updates to cron schedules without restart
func (f *JobFactory) SyncJobSchedulesFromDB() error {
	schedules, err := f.DB.GetEnabledJobSchedules(f.ctx)
	if err != nil {
		return err
	}

	log.Printf("Found %d enabled job schedules in database", len(schedules))

	for _, schedule := range schedules {
		log.Printf("Job: %s, Schedule: %s, Last Run: %v",
			schedule.JobName,
			schedule.CronSchedule,
			schedule.LastRunAt, 
		)
	}

	return nil
}

// InitializeJobSchedules ensures all jobs have entries in the job_schedules table
func (jf *JobFactory) InitializeJobSchedules() error {
	jobs := jf.CreateJobs()

	for _, job := range jobs {
		// Check if job schedule exists
		_, err := jf.DB.GetJobScheduleByName(jf.ctx, job.Name())
		if err == sql.ErrNoRows {
			// Create the job schedule
			err = jf.jobManager.ensureJobScheduleExists(job)
			if err != nil {
				log.Printf("Failed to initialize schedule for job %s: %v", job.Name(), err)
				continue
			}
			log.Printf("Initialized schedule for job: %s", job.Name())
		} else if err != nil {
			log.Printf("Error checking job schedule for %s: %v", job.Name(), err)
			continue
		}
	}

	return nil
}
