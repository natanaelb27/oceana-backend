package utils

import (
	"os"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"
)

func SendMail(to []string, subject, message string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", os.Getenv("MAIL_AUTH_EMAIL"))
	mailer.SetHeader("To", strings.Join(to, ","))
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", message)

	portStr := os.Getenv("MAIL_SMTP_PORT")
	port, err := strconv.Atoi(portStr)

	if err != nil {
		return err
	}

	dialer := gomail.NewDialer(
		os.Getenv("MAIL_SMTP_HOST"),
		port,
		os.Getenv("MAIL_AUTH_EMAIL"),
		os.Getenv("MAIL_AUTH_PASSWORD"),
	)

	err = dialer.DialAndSend(mailer)

	if err != nil {

		return err
	}

	return nil

}
