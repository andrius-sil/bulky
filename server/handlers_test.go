package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var testToken string

const TestTokenEnv = "TEST_TOKEN"

func init() {
	testToken = os.Getenv(TestTokenEnv)
	if testToken == "" {
		fmt.Printf("'%s' env variable is empty or unset.\n", TestTokenEnv)
		os.Exit(1)
	}
}

func TestGetActivitiesHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/activities", nil)
	if err != nil {
		t.Fatal(err)
	}

	q := req.URL.Query()
	q.Add("before", "1525215599")
	q.Add("after", "1524524400")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+testToken)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetActivitiesHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v expected %v",
			status, http.StatusOK)
	}
}
