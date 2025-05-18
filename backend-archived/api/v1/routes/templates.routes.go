package routes

import (
	"email-marketing-service/api/v1/middleware"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type TemplateRoute struct {
	db *gorm.DB
}

func NewTemplateRoute(db *gorm.DB) *TemplateRoute {
	return &TemplateRoute{db: db}
}

func (ur *TemplateRoute) InitRoutes(router *mux.Router) {
	templateController, _ := InitializeTemplateController(ur.db)
	router.HandleFunc("/create-marketing-template", middleware.JWTMiddleware(templateController.CreateAndUpdateTemplate)).Methods("POST", "OPTIONS")
	router.HandleFunc("/get-all-marketing-templates", middleware.JWTMiddleware(templateController.GetAllMarketingTemplates)).Methods("GET", "OPTIONS")
	router.HandleFunc("/get-all-transactional-templates", middleware.JWTMiddleware(templateController.GetAllTransactionalTemplates)).Methods("GET", "OPTIONS")
	router.HandleFunc("/get-transactional-template/{templateId}", middleware.JWTMiddleware(templateController.GetTransactionalTemplate)).Methods("GET", "OPTIONS")
	router.HandleFunc("/get-marketing-template/{templateId}", middleware.JWTMiddleware(templateController.GetMarketingTemplate)).Methods("GET", "OPTIONS")
	router.HandleFunc("/update-template/{templateId}", middleware.JWTMiddleware(templateController.UpdateTemplate)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/send-test-mails", middleware.JWTMiddleware(templateController.SendTestMail)).Methods("POST", "OPTIONS")
	router.HandleFunc("/delete-template/{templateId}", middleware.JWTMiddleware(templateController.DeleteTemplate)).Methods("DELETE", "OPTIONS")
}
