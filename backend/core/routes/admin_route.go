package routes

import (
	"email-marketing-service/core/handler/admin/auth/controller"
	authService "email-marketing-service/core/handler/admin/auth/services"
	campaigncontroller "email-marketing-service/core/handler/admin/campaigns/controller"
	campaignservice "email-marketing-service/core/handler/admin/campaigns/services"
	emailTemplateController "email-marketing-service/core/handler/admin/email-templates/controller"
	emailTemplateService "email-marketing-service/core/handler/admin/email-templates/services"
	planController "email-marketing-service/core/handler/admin/plans/controller"
	planservice "email-marketing-service/core/handler/admin/plans/service"
	supportController "email-marketing-service/core/handler/admin/support/controller"
	supportService "email-marketing-service/core/handler/admin/support/service"
	sysController "email-marketing-service/core/handler/admin/systems/controllers"
	sysService "email-marketing-service/core/handler/admin/systems/services"
	userController "email-marketing-service/core/handler/admin/users/controller"
	userService "email-marketing-service/core/handler/admin/users/services"
	"email-marketing-service/core/middleware"
	db "email-marketing-service/internal/db/sqlc"
	"github.com/gorilla/mux"
)
type AdminRoute struct {
	store db.Store
}

func NewAdminRoute(store db.Store) *AdminRoute {
	return &AdminRoute{
		store: store,
	}
}

func (a *AdminRoute) InitRoutes(r *mux.Router) {
	authRouter := r.PathPrefix("/auth").Subrouter()
	authSvc := authService.NewAdminAuthService(a.store)
	authHandler := controller.NewAdminAuthController(authSvc)

	{
		authRouter.HandleFunc("/create", authHandler.CreateAdmin).Methods("POST", "OPTIONS")
		authRouter.HandleFunc("/login", authHandler.AdminLogin).Methods("POST", "OPTIONS")
		authRouter.HandleFunc("/refresh-token", authHandler.RefreshTokenHandler).Methods("POST", "OPTIONS")
	}

	systemRouter := r.PathPrefix("/system-settings").Subrouter()
	sysSvc := sysService.NewAdminSystemsService(a.store)
	systemRouter.Use(middleware.AdminJWTMiddleware)
	sysHandler := sysController.NewAdminSystemsController(sysSvc)

	{
		systemRouter.HandleFunc("/create", sysHandler.CreateRecords).Methods("POST", "OPTIONS")
		systemRouter.HandleFunc("/fetch/{domain}", sysHandler.GetDNSRecords).Methods("GET", "OPTIONS")
		systemRouter.HandleFunc("/delete/{domain}", sysHandler.DeleteDNSRecords).Methods("DELETE", "OPTIONS")
		systemRouter.HandleFunc("/logs/app", sysHandler.ReadAppLogs).Methods("GET", "OPTIONS")
		systemRouter.HandleFunc("/logs/request", sysHandler.ReadRequestLogs).Methods("GET", "OPTIONS")
	}

	planRoute := r.PathPrefix("/plans").Subrouter()
	planService := planservice.NewPlanService(a.store)
	planRoute.Use(middleware.AdminJWTMiddleware)
	planController := planController.NewPlanController(planService)

	{
		planRoute.HandleFunc("/create", planController.CreatePlan).Methods("POST", "OPTIONS")
		planRoute.HandleFunc("/get", planController.GetAllPlans).Methods("GET", "OPTIONS")
		planRoute.HandleFunc("/get/{planId}", planController.GetPlanByID).Methods("GET", "OPTIONS")
		planRoute.HandleFunc("/update/{planId}", planController.UpdatePlan).Methods("PUT", "OPTIONS")
		planRoute.HandleFunc("/delete/{planId}", planController.DeletePlan).Methods("DELETE", "OPTIONS")
	}

	supportRouter := r.PathPrefix("/support").Subrouter()
	supportRouter.Use(middleware.AdminJWTMiddleware)
	supportService := supportService.NewAdminSupportService(a.store)
	supportController := supportController.NewAdminSupportController(supportService)

	{
		supportRouter.HandleFunc("/reply/{ticketId}", supportController.ReplyTicket).Methods("PUT", "OPTIONS")
		supportRouter.HandleFunc("/get/all", supportController.GetAllTickets).Methods("GET", "OPTIONS")
		supportRouter.HandleFunc("/get/pending", supportController.GetPendingTickets).Methods("GET", "OPTIONS")
		supportRouter.HandleFunc("/get/closed", supportController.GetClosedTickets).Methods("GET", "OPTIONS")
	}

	userRoute := r.PathPrefix("/users").Subrouter()
	userRoute.Use(middleware.AdminJWTMiddleware)
	userSvc := userService.NewAdminUsersServices(a.store)
	userCtrl := userController.NewAdminUsersController(userSvc)

	{
		userRoute.HandleFunc("/get", userCtrl.GetAllUsers).Methods("GET", "OPTIONS")
		userRoute.HandleFunc("/get/verified", userCtrl.GetVerifiedUsers).Methods("GET", "OPTIONS")
		userRoute.HandleFunc("/get/unverified", userCtrl.GetUnVerfiedUsers).Methods("GET", "OPTIONS")
		userRoute.HandleFunc("/block/{userId}", userCtrl.BlockUser).Methods("PUT", "OPTIONS")
		userRoute.HandleFunc("/unblock/{userId}", userCtrl.UnblockUser).Methods("PUT", "OPTIONS")
		userRoute.HandleFunc("/verify/{userId}", userCtrl.VerifyUser).Methods("PUT", "OPTIONS")
		userRoute.HandleFunc("/delete/{userId}", userCtrl.DeleteUser).Methods("DELETE", "OPTIONS")
		userRoute.HandleFunc("/get/{userId}", userCtrl.GetUserByID).Methods("GET", "OPTIONS")
		userRoute.HandleFunc("/stats/{userId}", userCtrl.GetUserStats).Methods("GET", "OPTIONS")
	}

	campaignRoute := r.PathPrefix("/campaigns").Subrouter()
	campaignRoute.Use(middleware.AdminJWTMiddleware)
	cmpService := campaignservice.NewAdminCampaignService(a.store)
	cmpController := campaigncontroller.NewAdminCampaignController(cmpService)

	{
		campaignRoute.HandleFunc("/get/{userId}/{companyId}", cmpController.GetAllUserCampaigns).Methods("GET", "OPTIONS")
		campaignRoute.HandleFunc("/get/single/{userId}/{companyId}/{campaignId}", cmpController.GetSingleCampaign).Methods("GET", "OPTIONS")
		campaignRoute.HandleFunc("/get-campaign-recipients/{campaignId}/{companyId}", cmpController.GetAllRecipientsForACampaign).Methods("GET", "OPTIONS")
	}

	emailTemplateRoute := r.PathPrefix("/gallery-templates").Subrouter()
	emailTemplateRoute.Use(middleware.AdminJWTMiddleware)
	emailTempService := emailTemplateService.NewAdminTemplatesService(a.store)
	emailTempController := emailTemplateController.NewAdminTemplateController(emailTempService)

	{
		emailTemplateRoute.HandleFunc("/create", emailTempController.CreateGalleryTemplate).Methods("POST", "OPTIONS")
		emailTemplateRoute.HandleFunc("/get/{type}", emailTempController.GetTemplatesByType).Methods("GET", "OPTIONS")
		emailTemplateRoute.HandleFunc("/get/{templateId}", emailTempController.GetTemplateById).Methods("GET", "OPTIONS")

	}

	templateRoute := r.PathPrefix("/templates").Subrouter()
	templateRoute.Use(middleware.AdminJWTMiddleware)
	templateService := emailTemplateService.NewAdminUserTemplatesService(a.store)
	adminUserTemplateController := emailTemplateController.NewAdminUserTemplateController(templateService)

	{
		templateRoute.HandleFunc("/get/{userId}/{type}", adminUserTemplateController.GetUserTemplates).Methods("GET", "OPTIONS")
		templateRoute.HandleFunc("/get/single/{userId}/{templateId}", adminUserTemplateController.GetSingleTemplate).Methods("GET", "OPTIONS")
	}

}
