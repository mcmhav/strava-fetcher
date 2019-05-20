package activities

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"google.golang.org/appengine"
)

func handleFetchActivites(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	userIter, err := getUsers(ctx)
	if err != nil {
		fmt.Fprintf(w, "Errorrr %v", err)
		return
	}
	userDistances, err := getDistancesForUsers(ctx, userIter)
	if err != nil {
		fmt.Fprintf(w, "Errorrr %v", err)
		return
	}

	response, err := persistDistancesToFirestore(ctx, userDistances)
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

	var userDistances []UserDistance

	for _, userDistanceDoc := range userDistancesDocs {
		var userDinstance UserDistance

		userDistanceDoc.DataTo(&userDinstance)

		userDinstance.UserID = userDistanceDoc.Ref.ID

		userDistances = append(userDistances, userDinstance)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userDistances)
}

func addHandlers() {
	http.HandleFunc("/fetchFromStrava", handleFetchActivites)
	http.HandleFunc("/activites", handleActivites)
}

func main() {
	log.Printf("rairai")
}

func init() {
	addHandlers()
}
