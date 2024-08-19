// routes/route_interface.go
package routes

import "github.com/gorilla/mux"

type RouteInterface interface {
	InitRoutes(r *mux.Router)
}