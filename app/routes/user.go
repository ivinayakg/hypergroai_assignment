package routes

import (
	"github.com/gorilla/mux"
	"github.com/ivinayakg/hypergroai_assignment/controllers"
	"github.com/ivinayakg/hypergroai_assignment/middleware"
)

func UserRoutes(r *mux.Router) {
	r.HandleFunc("/sign_in_with_google", controllers.SignInWithGoogle).Methods("GET")
	r.HandleFunc("/signin/callback", controllers.CallbackSignInWithGoogle).Methods("GET")

	protectedR := r.NewRoute().Subrouter()
	protectedR.Use(middleware.Authentication)
	protectedR.HandleFunc("/self", controllers.SelfUser).Methods("GET")
}
