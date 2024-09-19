package cronjobs

import (
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/services"
	"gorm.io/gorm"
)

type SendScheduledCampaignJobs struct {
	campaignSVC *services.CampaignService
}

func NewSendScheduledCampaignJobs(db *gorm.DB) *SendScheduledCampaignJobs {
	campaignRepository := repository.NewCampaignRepository(db)
	contactRepository := repository.NewContactRepository(db)
	templateRepository := repository.NewTemplateRepository(db)
	mailUsageRepository := repository.NewMailUsageRepository(db)
	subscriptionRepository := repository.NewSubscriptionRepository(db)
	userRepository := repository.NewUserRepository(db)
	domainRepository := repository.NewDomainRepository(db)
	userNotification := repository.NewUserNotificationRepository(db)
	campaignService := services.NewCampaignService(campaignRepository, contactRepository, templateRepository, mailUsageRepository, subscriptionRepository, userRepository, domainRepository,userNotification)
	return &SendScheduledCampaignJobs{
		campaignSVC: campaignService,
	}
}

func (j *SendScheduledCampaignJobs) Run() {
	j.campaignSVC.SendScheduledCampaigns()
}

func (j *SendScheduledCampaignJobs) Schedule() string {
	return "*/5 * * * *" // Every 5 minutes
}
