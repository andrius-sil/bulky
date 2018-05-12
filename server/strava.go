package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type StravaSummaryActivityModel struct {
	Id           int
	Name         string
	Start_date   string
	Distance     float64
	Start_latlng [2]float64
	End_latlng   [2]float64
	Commute      bool
	Private      bool
	Url          string
}

type StravaUpdatableActivityModel struct {
	Private bool `json:"private"`
	Commute bool `json:"commute"`
}

var httpClient *http.Client

const BaseUrl = "https://www.strava.com/api/v3"

func init() {
	httpClient = &http.Client{}
}

func fetchActivities(accessToken string, afterEpoch, beforeEpoch string) ([]StravaSummaryActivityModel, error) {
	fetchMethod := "GET"
	fetchUrl := fmt.Sprintf("%s/athlete/activities?before=%s&after=%s", BaseUrl, beforeEpoch, afterEpoch)
	fmt.Printf("%s sent - %s\n", fetchMethod, fetchUrl)

	req, err := http.NewRequest(fetchMethod, fetchUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", accessToken)

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		return nil, fmt.Errorf("Strava API returned %d (%s)\n", resp.StatusCode, bodyString)
	}

	data := make([]StravaSummaryActivityModel, 0)
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// TODO: check status code
func updateActivity(accessToken string, id int, privateValue, commuteValue bool) error {
	updateMethod := "PUT"
	updateUrl := fmt.Sprintf("%s/activities/%d", BaseUrl, id)
	fmt.Printf("%s sent - %s\n", updateMethod, updateUrl)

	payload := StravaUpdatableActivityModel{privateValue, commuteValue}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(updateMethod, updateUrl, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", accessToken)
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
