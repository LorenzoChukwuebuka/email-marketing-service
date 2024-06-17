//go:build wireinject
// +build wireinject

package routes

import (
	"email-marketing-service/api/controllers"
	adminController "email-marketing-service/api/controllers/admin"
	"email-marketing-service/api/repository"
	adminrepository "email-marketing-service/api/repository/admin"
	"email-marketing-service/api/services"
	adminservice "email-marketing-service/api/services/admin"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func InitializeUserController(db *gorm.DB) (*controllers.UserController, error) {
	wire.Build(
		controllers.NewUserController,
		services.NewUserService,
		services.NewOTPService,
		repository.NewOTPRepository,
		repository.NewUserRepository,
		repository.NewPlanRepository,
		repository.NewSubscriptionRepository,
		repository.NewBillingRepository,
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
		repository.NewDailyMailCalcRepository,
		repository.NewUserRepository,
		repository.NewMailStatusRepository,
	)

	return nil, nil
}

func InitializeSubscriptionController(db *gorm.DB) (*controllers.SubscriptionController, error) {
	wire.Build(
		controllers.NewSubscriptionController,
		services.NewSubscriptionService,
		repository.NewSubscriptionRepository,
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

func InitializeSupportTicketController(db *gorm.DB) (*controllers.SupportTicketController, error) {
	wire.Build(
		controllers.NewTicketController,
		services.NewSupportTicketService,
		repository.NewSupportRepository,
		repository.NewUserRepository,
	)

	return nil, nil
}
