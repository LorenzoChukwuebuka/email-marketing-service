package cronjobs

import (
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/services"

	"gorm.io/gorm"
)

type DailyFreeMailCalcJob struct {
	SMTPMailSVC *services.SMTPMailService
}

func NewDailMailCalculationJob(db *gorm.DB) *DailyFreeMailCalcJob {

	apiKeyRepository := repository.NewAPIkeyRepository(db)
	apiKeyService := services.NewAPIKeyService(apiKeyRepository)
	subscriptionRepository := repository.NewSubscriptionRepository(db)
	mailUsageRepository := repository.NewMailUsageRepository(db)
	userRepository := repository.NewUserRepository(db)
	mailStatusRepository := repository.NewMailStatusRepository(db)
	planRepository := repository.NewPlanRepository(db)
	smtpkeyRepo := repository.NewSMTPkeyRepository(db)
	smtpMailService := services.NewSMTPMailService(apiKeyService, subscriptionRepository, mailUsageRepository, userRepository, mailStatusRepository, planRepository, smtpkeyRepo)

	return &DailyFreeMailCalcJob{
		SMTPMailSVC: smtpMailService,
	}
}

func (j *DailyFreeMailCalcJob) Run() {
	j.SMTPMailSVC.CreateRecordForDailyMailCalculation()
}

func (j *DailyFreeMailCalcJob) Schedule() string {
	return "0 0 0 * * *"
}
