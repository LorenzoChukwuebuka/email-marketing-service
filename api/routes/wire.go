//go:build wireinject
// +build wireinject

package routes

import (
	"email-marketing-service/api/controllers"
	"email-marketing-service/api/repository"
	"email-marketing-service/api/services"
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


func InitializeAPIKeyController(db *gorm.DB)(*controllers.ApiKeyController,error){
	wire.Build(
		controllers.NewAPIKeyController,
		services.NewAPIKeyService,
		repository.NewAPIkeyRepository,
	)
	return nil, nil
}


func InitializeUserssionController(db *gorm.DB)(*controllers.UserSessionController,error){
	wire.Build(
		controllers.NewUserSessionController,
		services.NewUserSessionService,
		repository.NewUserRepository,
		repository.NewUserSessionRepository,
	)
	return nil, nil
}


func InitializeTransactionController(db *gorm.DB)(*controllers.TransactionController,error){
	wire.Build(
		controllers.NewTransactionController,
		services.NewBillingService,
		repository.NewBillingRepository,
		services.NewSubscriptionService, 
		repository.NewSubscriptionRepository,

	)
	return nil, nil
}

func InitializeSMTPController(db *gorm.DB)(*controllers.SMTPMailController,error){
	wire.Build(
		controllers.NewSMTPMailController,
		services.NewAPIKeyService,
		repository.NewAPIkeyRepository,
		repository.NewSubscriptionRepository,
		services.NewSMTPMailService,
		repository.NewDailyMailCalcRepository,
	)

	return nil, nil
}

func InitializeSubscriptionController(db *gorm.DB)(*controllers.SubscriptionController,error){
	wire.Build(
	controllers.NewSubscriptionController,
	services.NewSubscriptionService,
	repository.NewSubscriptionRepository,
	)
	return nil, nil
}