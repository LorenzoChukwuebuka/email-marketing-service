//go:build wireinject
// +build wireinject

package routes

import (
	"email-marketing-service/api/v1/controllers"
	adminController "email-marketing-service/api/v1/controllers/admin"
	"email-marketing-service/api/v1/repository"
	adminrepository "email-marketing-service/api/v1/repository/admin"
	"email-marketing-service/api/v1/services"
	adminservice "email-marketing-service/api/v1/services/admin"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitializeUserController(db *gorm.DB) (*controllers.UserController, error) {
	wire.Build(
		controllers.NewUserController,
		services.NewUserService,
		services.NewOTPService,
		services.NewSenderServices,
		repository.NewOTPRepository,
		repository.NewUserRepository,
		repository.NewPlanRepository,
		repository.NewSubscriptionRepository,
		repository.NewBillingRepository,
		repository.NewMailUsageRepository,
		repository.NewSMTPkeyRepository,
		repository.NewSenderRepository,
		repository.NewDomainRepository,
		repository.NewUserNotificationRepository,
		adminrepository.NewAdminNoficationRepository,
	)

	return nil, nil
}

func InitializePlanController(db *gorm.DB) (*controllers.PlanController, error) {
	wire.Build(
		controllers.NewPlanController,
		services.NewPlanService,
		repository.NewPlanRepository,
	)
	return nil, nil
}

func InitializeAPIKeyController(db *gorm.DB) (*controllers.ApiKeyController, error) {
	wire.Build(
		controllers.NewAPIKeyController,
		services.NewAPIKeyService,
		repository.NewAPIkeyRepository,
	)
	return nil, nil
}

func InitializeUserssionController(db *gorm.DB) (*controllers.UserSessionController, error) {
	wire.Build(
		controllers.NewUserSessionController,
		services.NewUserSessionService,
		repository.NewUserRepository,
		repository.NewUserSessionRepository,
	)
	return nil, nil
}

func InitializeTransactionController(db *gorm.DB) (*controllers.TransactionController, error) {
	wire.Build(
		controllers.NewTransactionController,
		services.NewBillingService,
		repository.NewBillingRepository,
		services.NewSubscriptionService,
		repository.NewSubscriptionRepository,
		repository.NewUserRepository,
		repository.NewPlanRepository,
		repository.NewMailUsageRepository,
	)
	return nil, nil
}

func InitializeSMTPController(db *gorm.DB) (*controllers.SMTPMailController, error) {
	wire.Build(
		controllers.NewSMTPMailController,
		services.NewAPIKeyService,
		repository.NewAPIkeyRepository,
		repository.NewSubscriptionRepository,
		services.NewSMTPMailService,
		repository.NewMailUsageRepository,
		repository.NewUserRepository,
		repository.NewMailStatusRepository,
		repository.NewPlanRepository,
	)

	return nil, nil
}

func InitializeSMTPKeyController(db *gorm.DB) (*controllers.SMTPKeyController, error) {
	wire.Build(
		controllers.NewSMTPKeyController,
		services.NewSMTPKeyService,
		repository.NewSMTPkeyRepository,
	)

	return nil, nil
}

func InitializeSubscriptionController(db *gorm.DB) (*controllers.SubscriptionController, error) {
	wire.Build(
		controllers.NewSubscriptionController,
		services.NewSubscriptionService,
		repository.NewSubscriptionRepository,
		repository.NewMailUsageRepository,
		repository.NewPlanRepository,
	)
	return nil, nil
}

func InitialiazePlanController(db *gorm.DB) (*controllers.PlanController, error) {
	wire.Build(
		controllers.NewPlanController,
		services.NewPlanService,
		repository.NewPlanRepository,
	)

	return nil, nil
}

func InitializeAdminController(db *gorm.DB) (*adminController.AdminController, error) {
	wire.Build(
		adminController.NewAdminController,
		adminservice.NewAdminService,
		adminrepository.NewAdminRepository,
	)
	return nil, nil
}

func InitializeContactController(db *gorm.DB) (*controllers.ContactController, error) {
	wire.Build(
		controllers.NewContactController,
		services.NewContactService,
		repository.NewContactRepository,
		repository.NewUserRepository,
		repository.NewSubscriptionRepository,
	)

	return nil, nil
}

func InitializeTemplateController(db *gorm.DB) (*controllers.TemplateController, error) {
	wire.Build(
		controllers.NewTemplateController,
		services.NewTemplateService,
		repository.NewTemplateRepository,
		repository.NewMailUsageRepository,
		repository.NewSubscriptionRepository,
		repository.NewUserRepository,
	)

	return nil, nil
}

func InitalizeCampaignController(db *gorm.DB) (*controllers.CampaignController, error) {
	wire.Build(
		controllers.NewCampaignController,
		services.NewCampaignService,
		repository.NewCampaignRepository,
		repository.NewContactRepository,
		repository.NewTemplateRepository,
		repository.NewMailUsageRepository,
		repository.NewSubscriptionRepository,
		repository.NewUserRepository,
		repository.NewDomainRepository,
		repository.NewUserNotificationRepository,
	)
	return nil, nil
}

func InitializeDomainController(db *gorm.DB) (*controllers.DomainController, error) {
	wire.Build(
		controllers.NewDomainController,
		services.NewDomainService,
		repository.NewDomainRepository,
	)
	return nil, nil
}

func InitializeSenderController(db *gorm.DB) (*controllers.SenderController, error) {
	wire.Build(
		controllers.NewSenderController,
		services.NewSenderServices,
		repository.NewSenderRepository,
		repository.NewDomainRepository,
	)

	return nil, nil
}

func InitializeSupportTicketController(db *gorm.DB) (*controllers.SupportTicketController, error) {
	wire.Build(
		controllers.NewSupportTicketController,
		services.NewSupportTicketService,
		repository.NewSupportRepository,
		repository.NewUserRepository,
		repository.NewUserNotificationRepository,
		adminrepository.NewAdminNoficationRepository,
		adminrepository.NewAdminRepository,
	)

	return nil, nil
}

func InitializeAdminUsersController(db *gorm.DB) (*adminController.AdminUsersController, error) {
	wire.Build(
		adminController.NewAdminUsersController,
		adminservice.NewAdminUsersService,
		adminrepository.NewAdminUsersRepository,
		// repository.NewUserNotificationRepository,
		// adminrepository.NewAdminNoficationRepository,
	)

	return nil, nil
}

func InitialiazeAdminSupportController(db *gorm.DB) (*adminController.AdminSupportTicketController, error) {
	wire.Build(
		adminController.NewAdminSupportTicketController,
		adminservice.NewAdminSupportService,
		adminrepository.NewAdminSupportRepository,
		repository.NewUserNotificationRepository,
		adminrepository.NewAdminNoficationRepository,
		repository.NewSupportRepository,
	)

	return nil, nil
}

func InitializeAdminCampaignController(db *gorm.DB) (*adminController.AdminCampaignController, error) {
	wire.Build(
		adminController.NewAdminCampaginController,
		adminservice.NewAdminCampaignService,
		adminrepository.NewAdminCampaignRepository,
	)

	return nil, nil
}
