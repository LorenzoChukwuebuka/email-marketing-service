package cronjobs

import (
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
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
		// Add other jobs here
	}
}

func SetupCronJobs(db *gorm.DB) *cron.Cron {
	c := cron.New()
	factory := NewJobFactory(db)
	jobs := factory.CreateJobs()

	for _, job := range jobs {
		c.AddFunc(job.Schedule(), job.Run)
	}

	return c
}
