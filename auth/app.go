package stravaauth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"cloud.google.com/go/firestore"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

var clientID = os.Getenv("STRAVA_CLIENT_ID")
var clientSecret = os.Getenv("STRAVA_CLIENT_SECRET")
var homeURL = "https://mcmhav.com"

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func getAuthorizeURL() (string, error) {
	var authorizeURL = "https://www.strava.com/oauth/authorize"
	var url, err = url.Parse(authorizeURL)
	if err != nil {
		return "", err
	}
	q := url.Query()
	q.Add("client_id", clientID)
	q.Add("response_type", "code")
	q.Add("redirect_uri", "https://strava-auth-dot-cake-mcmhav.appspot.com/authorization_successful")
	q.Add("scope", "activity:read")
	q.Add("state", "rairai")
	q.Add("approval_prompt", "force")

	url.RawQuery = q.Encode()

	return url.String(), nil
}

func handleAuthorize(w http.ResponseWriter, r *http.Request) {

	var authorizeURL, err = getAuthorizeURL()
	if err != nil {
		fmt.Println(err)
		return
	}
	http.Redirect(w, r, authorizeURL, http.StatusTemporaryRedirect)
}

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

func handleAuthorizeSuccessful(w http.ResponseWriter, r *http.Request) {
	var code = r.URL.Query().Get("code")

	values := map[string]string{
		"code":          code,
		"client_id":     clientID,
		"client_secret": clientSecret,
		"grant_type":    "authorization_code",
	}
	jsonValue, _ := json.Marshal(values)

	var authorizeSuccessURL = "https://www.strava.com/oauth/token"

	ctx := appengine.NewContext(r)

	client := urlfetch.Client(ctx)

	resp, err := client.Post(authorizeSuccessURL, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Fatalln(err)
		fmt.Fprintln(w, err)
		return
	}
	var user User
	json.NewDecoder(resp.Body).Decode(&user)

	addUserDataToUser(ctx, user)

	// fmt.Fprintf(w, "Hello %s", user.Athlete.FirstName)
	http.Redirect(w, r, homeURL, http.StatusTemporaryRedirect)
}

func addUserDataToUser(ctx context.Context, user User) {
	projectID := appengine.AppID(ctx)

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	defer client.Close()

	_, err = client.Collection("users").Doc(strconv.FormatInt(user.Athlete.ID, 10)).Set(ctx, user)
	if err != nil {
		log.Fatalf("Failed adding user: %v", err)
	}
}

func addHandlers() {
	http.HandleFunc("/authorize", handleAuthorize)
	http.HandleFunc("/", handler)
	http.HandleFunc("/authorization_successful", handleAuthorizeSuccessful)
}

func main() {
	addHandlers()
}

func init() {
	addHandlers()
}
