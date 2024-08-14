// routes/base.go
package routes

import (
	//"email-marketing-service/api/v1/middleware"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Route struct {
	router  *mux.Router
	db      *gorm.DB
	useAuth bool
}

func NewRoute(db *gorm.DB, useAuth bool) Route {
	return Route{
		db:      db,
		useAuth: useAuth,
	}
}

func (r *Route) InitBaseRouter(router *mux.Router) {
	r.router = router
	if r.useAuth {
		//r.router.Use(middleware.JWTMiddleware)
		print("hello")
	}
}

// InitRoutes is left to be implemented by the embedding structs
func (r *Route) InitRoutes() {}
