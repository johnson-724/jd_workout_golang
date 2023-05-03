package email

import (
	"os"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Email struct {
	FromName  string
	FromEmail string
	ToName    string
	ToEmail   string
	Subject   string
	Content   string
}

func Send(email Email) *error {
	from := mail.NewEmail(email.FromName, email.FromEmail)
	to := mail.NewEmail(email.ToName, email.ToEmail)

	println(email.FromEmail)

	message := mail.NewSingleEmail(from, email.Subject, to, "", email.Content)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	response, err := client.Send(message)

	if err != nil || response.StatusCode != 202 {
		return &err
	}

	return nil
}
