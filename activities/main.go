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
	ctx := appengine.NewContext(r)

	userIter, err := users.GetUsers(ctx)
	log.Printf("rairai")
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
	log.Printf(*response)
	fmt.Fprintf(w, "Activites handled, Distance:")
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func handleActivites(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	ctx := appengine.NewContext(r)

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
	// http.HandleFunc("/Activities", handleActivities)
}

func main() {
	log.Printf("rairai")
	addHandlers()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
