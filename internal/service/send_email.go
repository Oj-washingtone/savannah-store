package service

import (
	"os"

	"github.com/resend/resend-go/v2"
)

type EmailParams struct {
	To      string
	From    string
	Text    string
	Subject string
}

func SendEmail(to string, subject string, body string) error {

	apiKey := os.Getenv("RESEND_KEY")

	client := resend.NewClient(apiKey)

	params := &resend.SendEmailRequest{
		To:      []string{to},
		From:    "orders@linxs.co.ke",
		Text:    body,
		Subject: subject,
	}

	_, err := client.Emails.Send(params)

	if err != nil {
		return err
	}

	return nil
}
