package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"time"

	"github.com/resend/resend-go/v3"
)

type ResendMailer struct {
	fromEmail string
	apiKey    string
	client    *resend.Client
}

func NewResendClient(apiKey, fromEmail string) *ResendMailer {
	client := resend.NewClient(apiKey)
	return &ResendMailer{
		fromEmail: fromEmail,
		apiKey:    apiKey,
		client:    client,
	}
}

func (m *ResendMailer) Send(templateFile, username, email string, data any, isSandbox bool) (string, error) {
	// TODO: use isSandbox to not emails in development
	// template parsing & building
	tmpl, err := template.ParseFS(FS, "templates/"+templateFile)
	if err != nil {
		return "", err
	}

	subject := new(bytes.Buffer)
	if err = tmpl.ExecuteTemplate(subject, "subject", data); err != nil {
		return "", err
	}

	body := new(bytes.Buffer)
	if err = tmpl.ExecuteTemplate(body, "body", data); err != nil {
		return "", err
	}

	// Send email
	params := &resend.SendEmailRequest{
		From:    fmt.Sprintf("%s <%s>", FromName, m.fromEmail),
		To:      []string{email},
		Html:    body.String(),
		Subject: subject.String(),
	}
	return m.sendWithRetry(params)
}

// helper for exponential backoff retries
func (m *ResendMailer) sendWithRetry(params *resend.SendEmailRequest) (string, error) {
	var retryErr error
	for i := range MaxRetries {
		sent, retryErr := m.client.Emails.Send(params)
		if retryErr == nil {
			log.Printf("Email sent successfully: %v", sent.Id)
			return sent.Id, nil
		}
		// cant use logger??
		log.Printf("send attempt %d failed: %v", i+1, retryErr)
		time.Sleep(time.Second * time.Duration(i+1))
	}
	return "", fmt.Errorf("failed to send email after %d attempts, error: %v", MaxRetries, retryErr)
}
