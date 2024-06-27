package email

import (
	"net/mail"
	"net/smtp"
)

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

type Sender struct {
	emailAddress string
	auth         smtp.Auth
	smtpAddress  string
}

func NewSender(emailAddress, emailPassword, smtpAddress string) *Sender {
	auth := smtp.PlainAuth("", emailAddress, emailPassword, smtpAddress)
	return &Sender{
		emailAddress: emailAddress,
		auth:         auth,
		smtpAddress:  smtpAddress,
	}
}

func (s *Sender) Send(email string, subject string, body string) error {
	err := smtp.SendMail(s.smtpAddress, s.auth, s.emailAddress, []string{email}, []byte(subject+"\r\n"+body))

	return err
}
