package users

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

var clientID = os.Getenv("STRAVA_CLIENT_ID")
var clientSecret = os.Getenv("STRAVA_CLIENT_SECRET")
var grantType = "refresh_token"

type Athlete struct {
	ID            int64  `json:"id"`
	Username      string `json:"username"`
	ResourceState int    `json:"resource_state"`
	FirstName     string `json:"firstname"`
	LastName      string `json:"lastname"`
	City          string `json:"city"`
	State         string `json:"state"`
	Country       string `json:"country"`
	Sex           string `json:"sex"`
	Premium       bool   `json:"premium"`
	Summit        bool   `json:"summit"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	BadgeTypeID   int    `json:"badge_type_id"`
	ProfileMedium string `json:"profile_medium"`
	Profile       string `json:"profile"`
	Friend        string `json:"friend"`
	Follower      string `json:"follower"`
}

type User struct {
	TokenType    string `json:"token_type"`
	ExpiresAt    int64  `json:"expires_at"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	Athlete      Athlete
}

func GetUsers(ctx context.Context) ([]*firestore.DocumentSnapshot, error) {
	projectID := appengine.AppID(ctx)
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	userIter := client.Collection("users").Documents(ctx)
	users, err := userIter.GetAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func RefreshTokenForUser(ctx context.Context, user *User) (*User, error) {
	log.Println("refreshing token")
	client := urlfetch.Client(ctx)

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
		return nil, err
	}

	resp, err := client.Post("https://www.strava.com/oauth/token", "application/json", bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		return nil, err
	}

	var refreshTokenResponse RefreshTokenResponse

	json.NewDecoder(resp.Body).Decode(&refreshTokenResponse)

	user.AccessToken = refreshTokenResponse.AccessToken
	user.ExpiresAt = refreshTokenResponse.ExpiresAt

	return user, nil
}

func CheckIfTokenIsExpired(ctx context.Context, user *User) (*User, error) {
	now := time.Now().UnixNano() / 1000000
	if now > user.ExpiresAt {
		user, err := RefreshTokenForUser(ctx, user)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	return user, nil
}
