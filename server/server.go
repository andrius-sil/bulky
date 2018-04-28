package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const Port = "3000"
const AppUrl = "http://localhost:" + Port

const ClientID = "24292"
const BaseUrl = "https://www.strava.com/api/v3"

var clientSecret = os.Getenv("STRAVA_CLIENT_SECRET")

type StravaAthleteModel struct {
	Firstname string
	Lastname  string
}

type StravaTokenModel struct {
	Access_token string
	Athlete      StravaAthleteModel
}

type StravaActivityModel struct {
	Id               int
	Name             string
	Start_date_local string
	Distance         float64
	Start_latlng     [2]float64
	End_latlng       [2]float64
	Commute          bool
	Private          bool
}

type AuthResponse struct {
	Code string
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var authResponse AuthResponse
	err := json.NewDecoder(r.Body).Decode(&authResponse)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	values := url.Values{
		"client_id":     {ClientID},
		"client_secret": {clientSecret},
		"code":          {authResponse.Code},
	}
	resp, err := http.PostForm("https://www.strava.com/oauth/token", values)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	fmt.Fprint(w, string(body))
}

// TODO: use pagination and fetch till empty list is returned
func fetchActivities(accessToken string, afterEpoch, beforeEpoch string) ([]StravaActivityModel, error) {
	fetchUrl := fmt.Sprintf("%s/athlete/activities?before=%s&after=%s", BaseUrl, beforeEpoch, afterEpoch)

	client := &http.Client{}

	req, err := http.NewRequest("GET", fetchUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// TODO: use per page value
	data := make([]StravaActivityModel, 50)
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GetActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("Authorization")

	afterEpoch := r.URL.Query().Get("after")
	beforeEpoch := r.URL.Query().Get("before")

	activities, err := fetchActivities(accessToken, afterEpoch, beforeEpoch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(activities)
}

func run() {
	if clientSecret == "" {
		fmt.Println("'STRAVA_CLIENT_SECRET' env variable is empty or unset")
		os.Exit(1)
	}

	r := mux.NewRouter()

	api := r.PathPrefix("/api/").Subrouter()
	api.HandleFunc("/login", loginHandler).Methods("POST")
	api.HandleFunc("/activities", GetActivitiesHandler).Methods("GET")

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	fmt.Println(AppUrl)
	log.Fatal(http.ListenAndServe(":"+Port, loggedRouter))
}
