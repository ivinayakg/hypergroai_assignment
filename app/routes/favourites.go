package routes

import (
	"github.com/gorilla/mux"
	"github.com/ivinayakg/hypergroai_assignment/controllers"
	"github.com/ivinayakg/hypergroai_assignment/middleware"
)

func FavouritesRoutes(r *mux.Router) {
	protectedR := r.NewRoute().Subrouter()
	protectedR.Use(middleware.Authentication)
	protectedR.HandleFunc("/codes", controllers.GetUserFavouriteCodes).Methods("GET")
	protectedR.HandleFunc("", controllers.AddFavourite).Methods("POST")
	protectedR.HandleFunc("/{code}", controllers.RemoveFavourite).Methods("DELETE")
	protectedR.HandleFunc("", controllers.GetUserFavourite).Methods("GET")
}
