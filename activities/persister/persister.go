package persister

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

	"github.com/mcmhav/strava-fetcher/activities/distances"
	"google.golang.org/api/sheets/v4"
)

var spreadsheetID = os.Getenv("SPREADSHEET_ID")
var readRange = os.Getenv("READ_RANGE")
var developerKey = os.Getenv("DEVELOPER_KEY")
var aud = os.Getenv("GCP_AUD")

func PersistDistancesToSpreadsheet(ctx context.Context) (*string, error) {
	// ctxRai := context.WithValue(ctx, oauth2.HTTPClient, &http.Client{
	// 	Transport: &transport.APIKey{Key: developerKey},
	// })

	// client, err := google.DefaultClient(ctxRai, "https://www.googleapis.com/auth/spreadsheets")
	// sheets.New()
	client := urlfetch.Client(ctx)

	srv, err := sheets.NewService(
		ctx,
		option.WithHTTPClient(client),
		option.WithAPIKey(developerKey),
	)

	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
		return nil, err
	}
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	log.Printf("got here, wiii? %v", srv)

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		fmt.Println("Name, Major:")
		for _, row := range resp.Values {
			// Print columns A and E, which correspond to indices 0 and 4.
			fmt.Printf("%s, %s\n", row[0], row[4])
		}
	}

	var rai = "asdifhas"

	return &rai, nil
}

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
