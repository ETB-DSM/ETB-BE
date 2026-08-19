package infra

import (
	"gopkg.in/gomail.v2"

	"github.com/Heiji57/ETB-BE/config"
)

type Mailer struct {
	dialer *gomail.Dialer
	from   string
}

func NewSMTP(cfg *config.Config) *Mailer {
	d := gomail.NewDialer(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.User, cfg.SMTP.Password)
	return &Mailer{dialer: d, from: cfg.SMTP.From}
}

func (m *Mailer) Send(to, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)
	return m.dialer.DialAndSend(msg)
}
