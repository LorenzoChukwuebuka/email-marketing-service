package cronjobs

import (
	db "email-marketing-service/internal/db/sqlc"

	"context"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

type JobFactory struct {
	DB  *db.Store
	ctx context.Context
}

func NewJobFactory(db *db.Store, ctx context.Context) *JobFactory {
	return &JobFactory{DB: db, ctx: ctx}
}

func (f *JobFactory) CreateJobs() []Job {
	return []Job{
		NewUpdateExpiredSubscriptionJob(*f.DB, f.ctx),
	}
}

func SetupCronJobs(db *db.Store) *cron.Cron {
	c := cron.New(cron.WithSeconds(), cron.WithChain(
		cron.SkipIfStillRunning(cron.DefaultLogger),
		cron.Recover(cron.DefaultLogger),
	))
	ctx := context.Background()
	factory := NewJobFactory(db, ctx)
	jobs := factory.CreateJobs()

	for _, job := range jobs {
		_, err := c.AddFunc(job.Schedule(), wrapJobRun(job))
		if err != nil {
			log.Printf("Error adding job to cron: %v", err)
		}
	}

	return c
}

func wrapJobRun(job Job) func() {
	return func() {
		start := time.Now()
		log.Printf("Starting job: %T", job)
		job.Run()
		log.Printf("Finished job: %T, duration: %v", job, time.Since(start))
	}
}
