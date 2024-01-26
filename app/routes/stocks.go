package routes

import (
	"github.com/gorilla/mux"
	"github.com/ivinayakg/hypergroai_assignment/controllers"
	"github.com/ivinayakg/hypergroai_assignment/middleware"
)

func StockRoutes(r *mux.Router) {
	r.HandleFunc("/unverified", controllers.GetStocksUnverified).Methods("GET")

	protectedR := r.NewRoute().Subrouter()
	protectedR.Use(middleware.Authentication)
	protectedR.HandleFunc("/top", controllers.GetTopStocks).Methods("GET")
	protectedR.HandleFunc("", controllers.GetStocks).Methods("GET")
	protectedR.HandleFunc("/{id}", controllers.GetStockDetail).Methods("GET")
}
