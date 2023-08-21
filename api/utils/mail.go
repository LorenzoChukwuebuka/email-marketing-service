package utils

import (
	"gopkg.in/gomail.v2"
)

func SendMail(subject string, email string, message string) {
	// Mailtrap SMTP server settings
	smtpHost := "sandbox.smtp.mailtrap.io"
	smtpPort := 587
	smtpUsername := "your_mailtrap_username"
	smtpPassword := "your_mailtrap_password"

	// Create a new email message
	msg := gomail.NewMessage()
	msg.SetHeader("From", "sender@example.com")
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", message)

	// Initialize the SMTP sender
	d := gomail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)

	// Send the email
	if err := d.DialAndSend(msg); err != nil {
		panic(err)
	}

	println("Email sent successfully!")

}
