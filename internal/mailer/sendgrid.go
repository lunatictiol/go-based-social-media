package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	fromEmail string
	apiKey    string
	client    *sendgrid.Client
}

func NewMailer(apikey, fromEmail string) *SendGridMailer {
	client := sendgrid.NewSendClient(apikey)
	return &SendGridMailer{
		fromEmail: fromEmail,
		apiKey:    apikey,
		client:    client,
	}
}
func (m *SendGridMailer) Send(templateFile, username, email string, data any, issandbox bool) error {
	from := mail.NewEmail(fromName, m.fromEmail)
	to := mail.NewEmail(username, email)
	emailTemplate, err := template.ParseFS(fs, "templates"+templateFile)
	if err != nil {
		return nil
	}

	subject := new(bytes.Buffer)

	err = emailTemplate.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}

	body := new(bytes.Buffer)

	err = emailTemplate.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}

	message := mail.NewSingleEmail(from, subject.String(), to, "", body.String())

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &issandbox,
		},
	})

	response, err := m.client.Send(message)
	for i := 0; i < maxTries; i++ {
		if err != nil {
			log.Printf("failed to send email to %v attempt %d of %d", email, i+1, maxTries)
			log.Panicf("error %v", err)

			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		log.Printf("email sent status: %v", response.StatusCode)

		return nil
	}

	return fmt.Errorf("failed to send the email")

}
