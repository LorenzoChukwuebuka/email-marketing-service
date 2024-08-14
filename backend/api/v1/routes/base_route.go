package routes

import (
	"email-marketing-service/api/v1/middleware"
	"github.com/gorilla/mux"
)

type BaseRoute struct {
	Router  *mux.Router
	useAuth bool
}

func (br *BaseRoute) ApplyAuthMiddleware() {
	if br.useAuth {
		br.Router.Use(middleware.AuthMiddleware)
	}
}
