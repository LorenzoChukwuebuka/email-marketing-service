package routes

import (
	"email-marketing-service/core/handler/admin/auth/controller"
	authService "email-marketing-service/core/handler/admin/auth/services"
	planController "email-marketing-service/core/handler/admin/plans/controller"
	planservice "email-marketing-service/core/handler/admin/plans/service"
	sysController "email-marketing-service/core/handler/admin/systems/controllers"
	sysService "email-marketing-service/core/handler/admin/systems/services"
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
	sysHandler := sysController.NewAdminSystemsController(sysSvc)
	{
		systemRouter.HandleFunc("/create", sysHandler.CreateRecords).Methods("POST", "OPTIONS")
		systemRouter.HandleFunc("/fetch/{domain}", sysHandler.GetDNSRecords).Methods("GET", "OPTIONS")
		systemRouter.HandleFunc("/delete/{domain}", sysHandler.DeleteDNSRecords).Methods("DELETE", "OPTIONS")
	}

	planRoute := r.PathPrefix("/plans").Subrouter()
	planService := planservice.NewPlanService(a.store)
	planController := planController.NewPlanController(planService)

	{
		planRoute.HandleFunc("/create", planController.CreatePlan).Methods("POST", "OPTIONS")
		planRoute.HandleFunc("/get", planController.GetAllPlans).Methods("GET", "OPTIONS")
		planRoute.HandleFunc("/get/{planId}", planController.GetPlanByID).Methods("GET", "OPTIONS")
		planRoute.HandleFunc("/update/{planId}", planController.UpdatePlan).Methods("PUT", "OPTIONS")
		planRoute.HandleFunc("/delete/{planId}", planController.DeletePlan).Methods("DELETE", "OPTIONS")
	}

}
