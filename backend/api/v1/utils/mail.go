package utils

import (
	"bytes"
	adminmodel "email-marketing-service/api/v1/model/admin"
	"encoding/json"
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"os"
	"path/filepath"
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

// ReadSMTPSettingsFromFile reads SMTP settings from a JSON file
func ReadSMTPSettingsFromFile(domain string) (*adminmodel.SystemsSMTPSetting, error) {
	var dir string
	if os.Getenv("SERVER_MODE") == "production" {
		dir = "/app/backend/smtp_settings" // Absolute path for production
	} else {
		dir = "./smtp_settings" // Relative path for development
	}

	// Create the file path
	filePath := filepath.Join(dir, fmt.Sprintf("%s_smtp_settings.json", domain))

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read SMTP settings file: %w", err)
	}

	var settingsMap map[string]interface{}
	if err := json.Unmarshal(data, &settingsMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal SMTP settings: %w", err)
	}

	// Convert map back to struct
	smtpSetting := &adminmodel.SystemsSMTPSetting{
		Domain:         settingsMap["Domain"].(string),
		TXTRecord:      settingsMap["TXTRecord"].(string),
		DMARCRecord:    settingsMap["DMARCRecord"].(string),
		DKIMSelector:   settingsMap["DKIMSelector"].(string),
		DKIMPublicKey:  settingsMap["DKIMPublicKey"].(string),
		DKIMPrivateKey: settingsMap["DKIMPrivateKey"].(string),
		SPFRecord:      settingsMap["SPFRecord"].(string),
		MXRecord:       settingsMap["MXRecord"].(string),
	}

	return smtpSetting, nil
}

// extractDomain extracts the domain from an email address
func extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}
