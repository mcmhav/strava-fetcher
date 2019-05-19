package activities

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"

	"google.golang.org/api/sheets/v4"
)

var spreadsheetID = "1K-ZVHbi-S7PK3l95fhEcAhCGZ3vy9CaPaR2fdff-UQs"
var readRange = "distances!A2:E"
var developerKey = "AIzaSyCws6S6MtVP8R84ZXeOwdF-9BeC2DS3sGk"
var aud = "688928621123-2vd96hff89aur35c7n200i6dvr6vkrji.apps.googleusercontent.com"

func persistDistancesToSpreadsheet(ctx context.Context) (*string, error) {
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

func persistDistancesToFirestore(ctx context.Context, userDistances *[]UserDistance) (*string, error) {
	projectID := appengine.AppID(ctx)
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	for _, userDistance := range *userDistances {
		log.Printf("userid: %v", userDistance.UserID)
		log.Printf("distance: %v", userDistance.Distance)
		wr, err := client.Collection("distances").Doc(userDistance.UserID).Set(ctx, map[string]interface{}{
			"distance": userDistance.Distance,
		})
		if err != nil {
			return nil, err
		}
		log.Printf("writeresult %s", wr)
	}

	response := "ok"

	return &response, nil
}
