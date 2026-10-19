package utilities

import (
	"os"

	"github.com/resend/resend-go/v4"
)

func SendEmail(text string) error {
	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))
	params := &resend.SendEmailRequest{
		To:   []string{os.Getenv("RECIPIENT_EMAIL")},
		From: os.Getenv("MAIL_FROM"),
		Html: text,
		Text: text,
	}

	_, err := client.Emails.Send(params)

	return err
}
