package services

import (
	"os"

	"github.com/arnab333/golang-employee-management/helpers"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailDetails struct {
	Name    string
	Address string
}

func SendEmail(sendTo *EmailDetails) (*rest.Response, error) {
	from := mail.NewEmail("Arnab", "arnab3111@gmail.com")
	subject := "Sending with SendGrid is Fun"
	to := mail.NewEmail(sendTo.Name, sendTo.Address)
	htmlContent := "<strong>and easy to do anywhere, even with Go</strong>"
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
	client := sendgrid.NewSendClient(os.Getenv(helpers.EnvKeys.SENDGRID_API_KEY))
	response, err := client.Send(message)
	if err != nil {
		return nil, err
	}

	return response, nil
}
