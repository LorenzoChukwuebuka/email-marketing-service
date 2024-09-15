package cronjobs

import (
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/services"
	"gorm.io/gorm"
)

type UpdateExpiredSubscriptionJob struct {
	subscriptionService *services.SubscriptionService
}

func NewUpdateExpiredSubscriptionJob(db *gorm.DB) *UpdateExpiredSubscriptionJob {
	subscriptionRepo := repository.NewSubscriptionRepository(db)
	planRepo := repository.NewPlanRepository(db)
	dailyRepo := repository.NewMailUsageRepository(db)
	subscriptionService := services.NewSubscriptionService(subscriptionRepo, dailyRepo, planRepo)

	return &UpdateExpiredSubscriptionJob{
		subscriptionService: subscriptionService,
	}
}

func (j *UpdateExpiredSubscriptionJob) Run() {
	j.subscriptionService.UpdateExpiredSubscription()
}

func (j *UpdateExpiredSubscriptionJob) Schedule() string {
	return "0 0 * * *"
}
