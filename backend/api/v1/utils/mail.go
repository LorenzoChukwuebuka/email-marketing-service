package utils

import (
	"gopkg.in/gomail.v2"
)

// SMTPConfig holds the SMTP server configuration
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

// DefaultSMTPConfig returns the default SMTP configuration
func DefaultSMTPConfig() SMTPConfig {
	config := LoadEnv()
	return SMTPConfig{
		Host:     "sandbox.smtp.mailtrap.io",
		Port:     2525,
		Username: config.MailUsername,
		Password: config.MailPassword,
	}
}

// SendMail sends an email using the provided SMTP configuration
func SendMail(subject, email, message, sender string, smtpConfig *SMTPConfig) error {
	// Use default config if not provided
	if smtpConfig == nil {
		defaultConfig := DefaultSMTPConfig()
		smtpConfig = &defaultConfig
	}

	// Create a new email message
	msg := gomail.NewMessage()
	msg.SetHeader("From", sender)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", message)

	// Initialize the SMTP sender
	d := gomail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.Username, smtpConfig.Password)

	// Send the email
	if err := d.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
