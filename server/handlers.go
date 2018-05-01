package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type ActivitiesUpdateModel struct {
	Private map[int]bool
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
		"client_secret": {config.StravaCientSecret},
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

func GetActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("Authorization")

	afterEpoch := r.URL.Query().Get("after")
	beforeEpoch := r.URL.Query().Get("before")

	activities, err := fetchActivities(accessToken, afterEpoch, beforeEpoch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Append urls to fetched activities.
	for i, activity := range activities {
		activities[i].Url = fmt.Sprintf("https://www.strava.com/activities/%d", activity.Id)
	}

	json.NewEncoder(w).Encode(activities)
}

func ActivitiesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.Header.Get("Authorization")

	var data ActivitiesUpdateModel
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for activityId, privateValue := range data.Private {
		err := updateActivity(accessToken, activityId, privateValue)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
