package utils

import (
	"gopkg.in/gomail.v2"
)

func SendMail(subject string, email string, message string) error {

	config := LoadEnv()
	// Mailtrap SMTP server settings
	smtpHost := "sandbox.smtp.mailtrap.io"
	smtpPort := 2525
	smtpUsername := config.MailUsername
	smtpPassword := config.MailPassword

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
