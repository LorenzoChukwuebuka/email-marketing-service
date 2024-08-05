package routes 


import "github.com/gorilla/mux"

type Route interface {
    InitRoutes(r *mux.Router)
}