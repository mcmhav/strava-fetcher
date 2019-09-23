package distances

import (
	"context"
	"log"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/mcmhav/strava-fetcher/activities/activities"
	"github.com/mcmhav/strava-fetcher/activities/users"
)

func getTotalDistanceFromStravaActivities(activities []activities.Activity) float64 {
	var sumDistance float64
	for _, a := range activities {
		sumDistance += a.Distance
	}

	return sumDistance / 1000
}

func getDistanceForUser(ctx context.Context, user *users.User) (float64, error) {
	user, err := users.CheckIfTokenIsExpired(ctx, user)
	if err != nil {
		return -1, err
	}
	log.Printf(user.AccessToken)
	stravaActivities, err := activities.GetActivitiesForUser(ctx, user)
	if err != nil {
		return -1, err
	}

	sumDistance := getTotalDistanceFromStravaActivities(*stravaActivities)

	return sumDistance, nil
}

type UserDistance struct {
	Distance  float64 `json:"distance"`
	UserID    string  `json:"userId"`
	FirstName string  `json:"firstname"`
	LastName  string  `json:"lastname"`
}

func GetDistancesForUsers(ctx context.Context, userIter []*firestore.DocumentSnapshot) (*[]UserDistance, error) {

	var userDistances []UserDistance
	for _, userDoc := range userIter {
		var user users.User
		userDoc.DataTo(&user)

		sumDistance, err := getDistanceForUser(ctx, &user)
		if err != nil {
			log.Fatalf("Got error when handling user: %v", err)
		}
		log.Printf("Distance for user: %v, %v", user.Athlete.FirstName, sumDistance)

		var userDistance UserDistance

		userDistance.Distance = sumDistance
		userDistance.UserID = strconv.FormatInt(user.Athlete.ID, 10)
		userDistance.FirstName = user.Athlete.FirstName
		userDistance.LastName = user.Athlete.LastName

		userDistances = append(userDistances, userDistance)
	}

	return &userDistances, nil
}
