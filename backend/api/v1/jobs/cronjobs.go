package cronjobs

import (
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"log"
	"time"
)

type JobFactory struct {
	DB *gorm.DB
}

func NewJobFactory(db *gorm.DB) *JobFactory {
	return &JobFactory{DB: db}
}

func (f *JobFactory) CreateJobs() []Job {
	return []Job{
		NewUpdateExpiredSubscriptionJob(f.DB),
		NewDailMailCalculationJob(f.DB),
		NewSendScheduledCampaignJobs(f.DB),
		// Add other jobs here
	}
}

func SetupCronJobs(db *gorm.DB) *cron.Cron {
	c := cron.New(cron.WithSeconds(), cron.WithChain(
		cron.SkipIfStillRunning(cron.DefaultLogger),
		cron.Recover(cron.DefaultLogger),
	))

	factory := NewJobFactory(db)
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
