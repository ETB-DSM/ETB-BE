package infra

import (
	"github.com/resend/resend-go/v3"

	"github.com/Heiji57/ETB-BE/config"
)

type Mailer struct {
	client *resend.Client
	from   string
}

func NewMailer(cfg *config.Config) *Mailer {
	return &Mailer{
		client: resend.NewClient(cfg.Resend.APIKey),
		from:   cfg.Resend.From,
	}
}

func (m *Mailer) Send(to, subject, body string) error {
	params := &resend.SendEmailRequest{
		From:    m.from,
		To:      []string{to},
		Subject: subject,
		Html:    body,
	}
	_, err := m.client.Emails.Send(params)
	return err
}
