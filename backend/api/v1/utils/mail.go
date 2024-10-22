package utils

import (
	"bytes"
	adminmodel "email-marketing-service/api/v1/model/admin"
	adminrepository "email-marketing-service/api/v1/repository/admin"
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"strings"
	"sync"
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
		Host:     config.SMTP_SERVER,
		Port:     1025,
		Username: config.MailUsername,
		Password: config.MailPassword,
	}
}

// AsyncSendMail sends an email asynchronously using goroutines
func AsyncSendMail(subject, email, message, sender string, smtpConfig *SMTPConfig, wg *sync.WaitGroup) error {
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := sendMail(subject, email, message, sender, smtpConfig)
		if err != nil {
			log.Printf("Error sending email: %v", err)
		}
	}()

	return nil
}

func sendMail(subject, email, message, sender string, smtpConfig *SMTPConfig) error {
	// Use default config if not provided
	defaultConfig := DefaultSMTPConfig()
	if smtpConfig == nil {
		smtpConfig = &defaultConfig
	}

	// Create a new email message
	msg := gomail.NewMessage()
	msg.SetHeader("From", sender)
	msg.SetHeader("To", email)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", message)

	// Convert the email to bytes for signing
	var emailBuffer bytes.Buffer
	_, err := msg.WriteTo(&emailBuffer)
	if err != nil {
		return fmt.Errorf("failed to write email to buffer: %w", err)
	}
	emailBytes := emailBuffer.Bytes()

	// Sign the email if using default configurations
	if *smtpConfig == defaultConfig {
		smtpSettings, err := ReadSMTPSettingsFromFile(extractDomain(sender))
		if err != nil {
			return fmt.Errorf("failed to read SMTP settings: %w", err)
		}

		signedEmail, err := SignEmail(&emailBytes, smtpSettings.Domain, smtpSettings.DKIMSelector, smtpSettings.DKIMPrivateKey)
		if err != nil {
			log.Printf("Failed to sign email: %v. Proceeding with unsigned email.", err)
		}
		emailBytes = signedEmail
	}

	// Initialize the SMTP sender
	d := gomail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.Username, smtpConfig.Password)

	// Send the email
	if err := d.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}

 

func ReadSMTPSettingsFromDB(domain string, repo adminrepository.SystemRepository) (*adminmodel.SystemsSMTPSetting, error) {
	// Fetch the SMTP settings for the given domain
	settings, err := repo.GetSMTPSettings(domain)
	if err != nil {
		return nil, fmt.Errorf("failed to read SMTP settings from the database: %w", err)
	}

	return settings, nil
}

// extractDomain extracts the domain from an email address
func extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}
