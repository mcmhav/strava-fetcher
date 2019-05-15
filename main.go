package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

var clientID = os.Getenv("STRAVA_CLIENT_ID")
var clientSecret = os.Getenv("STRAVA_CLIENT_SECRET")
var grantType = "refresh_token"

type User struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	ExpiresAt    int64  `json:"expires_at"`
}

type Users struct {
	Users []User `json:"users"`
}

func readUsers() (Users, error) {
	jsonFile, err := os.Open("data/users.json")
	var users Users
	if err != nil {
		fmt.Println("Json read error")
		return users, err
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &users)

	return users, nil
}

func refreshTokenForUser(user User) User {
	fmt.Println("refreshing token")
	client := &http.Client{}

	type RefreshTokenResponse struct {
		TokenType    string `json:"token_type"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresAt    int64  `json:"expires_at"`
	}
	message := map[string]interface{}{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"grant_type":    grantType,
		"refresh_token": user.RefreshToken,
	}
	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		fmt.Println("dfsf")
	}

	resp, err := client.Post("https://www.strava.com/oauth/token", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		fmt.Println(err)
	}

	var refreshTokenResponse RefreshTokenResponse

	json.NewDecoder(resp.Body).Decode(&refreshTokenResponse)

	user.AccessToken = refreshTokenResponse.AccessToken
	user.ExpiresAt = refreshTokenResponse.ExpiresAt

	return user
}

func checkIfTokenIsExpired(user User) User {
	now := time.Now().UnixNano() / 1000000
	if now > user.ExpiresAt {
		user = refreshTokenForUser(user)
	}

	return user
}

type Activity struct {
	Distance float64 `json:"distance"`
}

func getActivitiesURL() (*url.URL, error) {
	var activitiesURL = "https://www.strava.com/api/v3/athlete/activities"
	var url, err = url.Parse(activitiesURL)
	if err != nil {
		return url, err
	}
	q := url.Query()
	// q.Set("after", "30")
	// q.Set("before", "30")
	// q.Set("page", "30")
	q.Add("per_page", "30")

	url.RawQuery = q.Encode()

	return url, nil
}

func getActivitiesForUser(user User) ([]Activity, error) {
	activitiesURL, err := getActivitiesURL()
	var activites []Activity
	if err != nil {
		return activites, err
	}
	fmt.Println(activitiesURL.String())
	req, err := http.NewRequest("GET", activitiesURL.String(), nil)
	if err != nil {
		return activites, err
	}

	req.Header.Add("authorization", "Bearer "+user.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return activites, err
	}

	json.NewDecoder(resp.Body).Decode(&activites)

	return activites, nil
}

func getTotalDistanceFromStravaActivities(activities []Activity) float64 {
	var sumDistance float64
	for i, a := range activities {
		fmt.Println(i, a.Distance)
		sumDistance += a.Distance
	}

	return sumDistance / 1000
}

func handleUser(user User) {
	user = checkIfTokenIsExpired(user)
	activities, err := getActivitiesForUser(user)
	if err != nil {
		fmt.Println(err)
	}

	sumDistance := getTotalDistanceFromStravaActivities(activities)

	fmt.Println(sumDistance)

}

func main() {
	var users, err = readUsers()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(len(users.Users))

	for i := 0; i < len(users.Users); i++ {
		fmt.Println("refresh token:", users.Users[i].RefreshToken)
		fmt.Println("access token:", users.Users[i].AccessToken)
		fmt.Println("expires in:", users.Users[i].ExpiresAt)

		handleUser(users.Users[i])
	}
}
