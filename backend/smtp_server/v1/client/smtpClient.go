package smtpserver

import (
	"log"
	"net/smtp"
)

const (
	smtpServer = "localhost:1025"
	from       = "sender@example.com"
	to         = "lobi@enugudisco.com"
	username   = "username"
	password   = "password"
)

func main() {
	// Set up authentication information.
	auth := smtp.PlainAuth("", username, password, "localhost")

	// Compose the message
	msg := []byte("To: " + to + "\r\n" +
		"Subject: Test Email\r\n" +
		"\r\n" +
		"This is a test email body.\r\n")

	// Send the email
	err := smtp.SendMail(smtpServer, auth, from, []string{to}, msg)
	if err != nil {
		log.Fatalf("Error sending email: %v", err)
	}

	log.Println("Email sent successfully!")
}
