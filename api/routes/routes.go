package routes

import (
	"email-marketing-service/api/controllers"

	"github.com/gorilla/mux"
)

var RegisterRoutes = func(router *mux.Router) {
	router.HandleFunc("/greet", controllers.Welcome).Methods("GET")
	router.HandleFunc("/user-signup", controllers.RegisterUser).Methods("POST")
	router.HandleFunc("/verify-user", controllers.VerifyUser).Methods("POST")
}
