package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/mcmhav/strava-fetcher/activities/distances"
	"github.com/mcmhav/strava-fetcher/activities/persister"
	"github.com/mcmhav/strava-fetcher/activities/users"
	"google.golang.org/appengine"
)

func handleFetchActivities(w http.ResponseWriter, r *http.Request) {
	// ctx := appengine.NewContext(r)
	ctx := r.Context()

	userIter, err := users.GetUsers(ctx)
	if err != nil {
		fmt.Fprintf(w, "Errorrr %v", err)
		return
	}
	userDistances, err := distances.GetDistancesForUsers(ctx, userIter)
	if err != nil {
		fmt.Fprintf(w, "Errorrr %v", err)
		return
	}
	response, err := persister.PersistDistancesToFirestore(ctx, userDistances)
	if err != nil {
		fmt.Fprintf(w, "Errorrr %v", err)
		return
	}
	log.Println(*response)
	fmt.Fprintf(w, "Activites handled, Distance:")
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func handleActivities(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	ctx := r.Context()

	projectID := appengine.AppID(ctx)
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		fmt.Fprintf(w, "Errorrr %v", err)
		return
	}
	defer client.Close()

	userDistancesDocumentsIter := client.Collection("distances").Documents(ctx)
	userDistancesDocs, err := userDistancesDocumentsIter.GetAll()

	var userDistances []distances.UserDistance

	for _, userDistanceDoc := range userDistancesDocs {
		var userDinstance distances.UserDistance

		userDistanceDoc.DataTo(&userDinstance)

		userDinstance.UserID = userDistanceDoc.Ref.ID

		userDistances = append(userDistances, userDinstance)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userDistances)
}

func addHandlers() {
	http.HandleFunc("/fetchFromStrava", handleFetchActivities)
	http.HandleFunc("/activites", handleActivities)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	addHandlers()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("Defaulting to port", port)
	}

	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
