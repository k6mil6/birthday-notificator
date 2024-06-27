package email

import (
	"gopkg.in/gomail.v2"
	"net/mail"
)

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

type Sender struct {
	emailAddress  string
	emailPassword string
	smtpAddress   string
	smtpPort      int
}

func NewSender(emailAddress, emailPassword, smtpAddress string, smtpPort int) *Sender {
	return &Sender{
		emailAddress:  emailAddress,
		emailPassword: emailPassword,
		smtpAddress:   smtpAddress,
		smtpPort:      smtpPort,
	}
}

func (s *Sender) Send(email, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.emailAddress)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(s.smtpAddress, s.smtpPort, s.emailAddress, s.emailPassword)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
