package utils

import (
	"gopkg.in/gomail.v2"
	"os"
)

func SendMail(subject string, email string, message string) error {

	LoadEnv()
	// Mailtrap SMTP server settings
	smtpHost := "sandbox.smtp.mailtrap.io"
	smtpPort := 2525
	smtpUsername := os.Getenv("MAIL_USERNAME")
	smtpPassword := os.Getenv("MAIL_PASSWORD")

	// Create a new email message
	msg := gomail.NewMessage()
	msg.SetHeader("From", "sender@example.com")
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", message)

	// Initialize the SMTP sender
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)

	// Send the email
	if err := d.DialAndSend(msg); err != nil {
		return err
	}

	return nil

}
