package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Config struct {
	StravaCientSecret string
}

var config Config

const (
	Port   = "3000"
	AppUrl = "http://localhost:" + Port

	ClientID        = "24292"
	ClientSecretEnv = "STRAVA_CLIENT_SECRET"
)

func main() {
	config.StravaCientSecret = os.Getenv(ClientSecretEnv)
	if config.StravaCientSecret == "" {
		fmt.Printf("'%s' env variable is empty or unset.\n", ClientSecretEnv)
		os.Exit(1)
	}

	r := mux.NewRouter()

	api := r.PathPrefix("/api/").Subrouter()
	api.HandleFunc("/login", loginHandler).Methods("POST")
	api.HandleFunc("/activities", GetActivitiesHandler).Methods("GET")
	api.HandleFunc("/activities_update", ActivitiesUpdateHandler).Methods("PUT")

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	fmt.Printf("Running on %s\n", AppUrl)
	log.Fatal(http.ListenAndServe(":"+Port, loggedRouter))
}
