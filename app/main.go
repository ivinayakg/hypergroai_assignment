package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ivinayakg/hypergroai_assignment/helpers"
	"github.com/ivinayakg/hypergroai_assignment/middleware"
	"github.com/ivinayakg/hypergroai_assignment/routes"
	"github.com/ivinayakg/hypergroai_assignment/utils"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func setupRoutes(router *mux.Router) {
	routes.UserRoutes(router.PathPrefix("/user").Subrouter())
	routes.MigrationFileRoutes(router.PathPrefix("/migration").Subrouter())
	routes.StockRoutes(router.PathPrefix("/stock").Subrouter())
	routes.FavouritesRoutes(router.PathPrefix("/favourite").Subrouter())
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	}

	PORT := os.Getenv("PORT")
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "./client_service.json")
	helpers.SaveAsJSON(os.Getenv("SERVICE_ACCOUNT_KEY"))

	allowed_origins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), " ")
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   allowed_origins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	r := mux.NewRouter()
	helpers.CreateDBInstance()
	helpers.RedisSetup()
	r.Use(middleware.LogMW)
	helpers.CreateGoogleCloudStorageClient()

	// create migrations
	go utils.RunMigrations()

	setupRoutes(r)
	routerProtected := corsHandler.Handler(r)

	fmt.Println("Starting the server on port " + PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", PORT), routerProtected))
}
