package mailsender

import (
	"errors"
	"fmt"
	"net/mail"
	"net/smtp"
)

var (
	ErrorInvalidEmail = errors.New("invalid email")
)

type Mailsender struct {
	UserName     string
	SmtpPassword string
	SmtpHost     string
	SmtpPort     string
}

func NewMailSender(username, smtpPassword, smtpHost, smtpPort string) *Mailsender {
	return &Mailsender{
		UserName:     username,
		SmtpPassword: smtpPassword,
		SmtpHost:     smtpHost,
		SmtpPort:     smtpPort,
	}
}

func (m *Mailsender) Send(to, body, subject string) error {
	toAddress, err := mail.ParseAddress(to)
	if err != nil {
		return ErrorInvalidEmail
	}

	auth := smtp.PlainAuth("", m.UserName, m.SmtpPassword, m.SmtpHost)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body +
		"\r\n")

	err = smtp.SendMail(m.SmtpHost+":"+m.SmtpPort, auth, m.UserName, []string{toAddress.String()}, msg)
	if err != nil {
		return fmt.Errorf("Can't send email  %v", err)
	}
	return nil
}
