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
	from := mail.NewEmail("Employee Management", os.Getenv(helpers.EnvKeys.SENDGRID_FROM_EMAIL))
	subject := "Account Created"
	to := mail.NewEmail(sendTo.Name, sendTo.Address)
	htmlContent := "An account with your email address has been created!"
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)
	client := sendgrid.NewSendClient(os.Getenv(helpers.EnvKeys.SENDGRID_API_KEY))
	response, err := client.Send(message)
	if err != nil {
		return nil, err
	}

	return response, nil
}
