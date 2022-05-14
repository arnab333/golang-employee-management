package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/arnab333/golang-employee-management/helpers"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GoogleServiceAccountConfig struct {
	Type                string `json:"type"`
	ProjectID           string `json:"project_id"`
	PrivateKeyID        string `json:"private_key_id"`
	PrivateKey          string `json:"private_key"`
	ClientEmail         string `json:"client_email"`
	ClientID            string `json:"client_id"`
	AuthURI             string `json:"auth_uri"`
	TokenURI            string `json:"token_uri"`
	AuthProviderCertURL string `json:"auth_provider_x509_cert_url"`
	ClientCertURL       string `json:"client_x509_cert_url"`
}

func calendarInit() []Holiday {
	ctx := context.Background()
	srv, err := calendar.NewService(ctx, option.WithCredentialsFile("google-credentials.json"))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	t := time.Now().Format(time.RFC3339)

	events, err := srv.Events.List(os.Getenv(helpers.EnvKeys.GOOGLE_CALENDAR_ID)).ShowDeleted(false).
		SingleEvents(true).TimeMin(t).MaxResults(250).OrderBy("startTime").Do()

	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}

	fmt.Println("Upcoming events:")
	var holidays []Holiday
	if len(events.Items) == 0 {
		fmt.Println("No upcoming events found.")
	} else {
		for _, item := range events.Items {
			holidays = append(holidays, Holiday{
				Summary: item.Summary,
				Date:    item.Start.Date,
			})
			// date := item.Start.DateTime
			// if date == "" {
			// 	date = item.Start.Date
			// }
			// fmt.Printf("%v (%v)\n", item.Summary, date)
		}

		return holidays
	}

	return nil
}
