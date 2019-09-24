package activities

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/mcmhav/strava-fetcher/activities/users"
)

var after = "1546297200"

func getActivitiesURL() (*url.URL, error) {
	var activitiesURL = "https://www.strava.com/api/v3/athlete/activities"
	var url, err = url.Parse(activitiesURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Got error when handling user: %v", err)
		return url, err
	}
	q := url.Query()
	q.Add("after", after)
	// q.Add("before", "30")
	// q.Add("page", "30")
	// q.Add("per_page", "30")

	url.RawQuery = q.Encode()

	return url, nil
}

type Activity struct {
	Distance float64 `json:"distance"`
}

func GetActivitiesForUser(ctx context.Context, user *users.User) (*[]Activity, error) {
	activitiesURL, err := getActivitiesURL()
	var activites []Activity
	if err != nil {
		fmt.Fprintf(os.Stderr, "Got error when handling user: %v", err)
		return nil, err
	}
	client := http.Client{}

	req, err := http.NewRequest("GET", activitiesURL.String(), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Got error when handling user: %v", err)
		return nil, err
	}

	log.Printf("Using accessToken: %v", user.AccessToken)

	req.Header.Add("authorization", "Bearer "+user.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Got error when handling user: %v", err)
		return nil, err
	}

	json.NewDecoder(resp.Body).Decode(&activites)

	return &activites, nil
}
