package routes

import (
	"github.com/gorilla/mux"
	"github.com/ivinayakg/hypergroai_assignment/controllers"
	"github.com/ivinayakg/hypergroai_assignment/middleware"
)

func MigrationFileRoutes(r *mux.Router) {
	protectedR := r.NewRoute().Subrouter()
	protectedR.Use(middleware.Authentication)
	protectedR.HandleFunc("/upload", controllers.UploadHandler).Methods("POST")
	protectedR.HandleFunc("/run", controllers.RunMigrations).Methods("POST")
	protectedR.HandleFunc("", controllers.GetMigrations).Methods("GET")
}
