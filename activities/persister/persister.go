package persister

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/appengine"

	"github.com/mcmhav/strava-fetcher/activities/distances"
)

var spreadsheetID = os.Getenv("SPREADSHEET_ID")
var readRange = os.Getenv("READ_RANGE")
var developerKey = os.Getenv("DEVELOPER_KEY")
var aud = os.Getenv("GCP_AUD")

func PersistDistancesToFirestore(ctx context.Context, userDistances *[]distances.UserDistance) (*string, error) {
	projectID := appengine.AppID(ctx)
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	for _, userDistance := range *userDistances {
		log.Printf("userid: %v", userDistance.UserID)
		log.Printf("distance: %v", userDistance.Distance)
		wr, err := client.Collection("distances").Doc(userDistance.UserID).Set(ctx, userDistance)
		if err != nil {
			return nil, err
		}
		log.Printf("writeresult %s", wr)
	}

	response := "ok"

	return &response, nil
}
