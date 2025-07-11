package mailer

import (
	"email-marketing-service/internal/common"
	"email-marketing-service/internal/config"
	"email-marketing-service/internal/helper"
	"fmt"
	"github.com/google/uuid"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// EmailConfig holds configuration for email service
type EmailConfig struct {
	Sender  string
	AppName string
	BaseURL string
}

// EmailData represents the data structure for email sending
type EmailData struct {
	To          string
	Subject     string
	Body        string
	Data        map[string]any
	Attachments []string
}

// EmailService manages email-related operations
type EmailService struct {
	config      EmailConfig
	wg          *sync.WaitGroup
	templateDir string
}

// NewEmailService creates a new instance of EmailService
func NewEmailService() *EmailService {
	cfg := config.LoadEnv()
	return &EmailService{
		config: EmailConfig{
			Sender:  cfg.SENDER,
			AppName: cfg.APPName,
			BaseURL: "http://localhost:5054",
		},
		wg:          &sync.WaitGroup{},
		templateDir: filepath.Join("internal", "mailer", "templates"),
	}
}

// LoadTemplate reads and returns the content of an email template
func (s *EmailService) LoadTemplate(templateName string) (string, error) {
	templatePath := filepath.Join(s.templateDir, templateName)
	mailTemplate, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %v", templateName, err)
	}
	return string(mailTemplate), nil
}

// FormatTemplate replaces placeholders in the HTML template with provided data
func (s *EmailService) FormatTemplate(templateContent string, data map[string]any) (string, error) {
	// Add common values to all templates
	if data == nil {
		data = make(map[string]any)
	}

	// Add default values
	if _, exists := data["AppName"]; !exists {
		data["AppName"] = s.config.AppName
	}
	if _, exists := data["BaseURL"]; !exists {
		data["BaseURL"] = s.config.BaseURL
	}

	// Parse the template with proper delimiters
	tmpl, err := template.New("emailTemplate").Delims("{{", "}}").Parse(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}

	var result strings.Builder
	if err := tmpl.Execute(&result, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}

	return result.String(), nil
}

// SendEmail sends an email using the provided template and data
func (s *EmailService) SendEmail(templateName string, emailData EmailData) error {
	// Load the template
	templateContent, err := s.LoadTemplate(templateName)
	if err != nil {
		return err
	}

	// Format the template with the provided data
	formattedBody, err := s.FormatTemplate(templateContent, emailData.Data)
	if err != nil {
		return err
	}

	// If Body is empty, use the formatted template as body
	body := emailData.Body
	if body == "" {
		body = formattedBody
	}

	// Send the email asynchronously
	return helper.AsyncSendMail(
		emailData.Subject,
		emailData.To,
		body,
		s.config.Sender,
		nil,
		s.wg,
	)
}

// SignUpMail sends a sign-up verification email
func (s *EmailService) SignUpMail(email string, username string, userID uuid.UUID, otp string) error {

	emailData := EmailData{
		To:      email,
		Subject: "Email Verification",
		Data: map[string]any{
			"Link":     fmt.Sprintf("%s/auth/account-verification", s.config.BaseURL),
			"Username": username,
			"Token":    otp,
			"UserId":   userID,
			"Email":    email,
		},
	}
	fmt.Printf("%v", emailData)
	return s.SendEmail(common.VerifyUserTemplate, emailData)
}

// verify user login
func (s *EmailService) VerifyUserLogin(email, username, otp, userID string) error {
	emailData := EmailData{
		To:      email,
		Subject: "Email Verification",
		Data: map[string]any{
			"Username": username,
			"Token":    otp,
			"Email":    email,
		},
	}
	return s.SendEmail(common.VerifyUserLoginTemplate, emailData)
}

// ResetPasswordMail sends a password reset email
func (s *EmailService) ResetPasswordMail(email, username, otp string) error {
	emailData := EmailData{
		To:      email,
		Subject: "Password Reset",
		Data: map[string]any{
			"Link":     fmt.Sprintf("%s/auth/reset-password", s.config.BaseURL),
			"Username": username,
			"Token":    otp,
			"Email":    email,
		},
	}
	return s.SendEmail(common.ResetPasswordTemplate, emailData)
}

// VerifySenderMail sends a sender verification email
func (s *EmailService) VerifySenderMail(username, userEmail, domainEmail, otp, userID string) error {
	emailData := EmailData{
		To:      domainEmail,
		Subject: "Verify a new Sender [Crabmailer]",
		Data: map[string]any{
			"Username":         username,
			"UserEmail":        userEmail,
			"DomainEmail":      domainEmail,
			"VerificationLink": fmt.Sprintf("%s/verifysender", s.config.BaseURL),
			"Token":            otp,
			"UserId":           userID,
		},
	}

	return s.SendEmail(common.VerifySenderTemplate, emailData)
}

// SubscriptionExpiryMail sends a subscription expiry notification
func (s *EmailService) SubscriptionExpiryMail(username, email, planName string) error {
	emailData := EmailData{
		To:      email,
		Subject: "Subscription Expiry Notification",
		Data: map[string]any{
			"Username": username,
			"PlanName": planName,
		},
	}

	return s.SendEmail(common.PlanExpiryTemplate, emailData)
}

// SubscriptionExpiryReminder sends a subscription expiry reminder
func (s *EmailService) SubscriptionExpiryReminder(username, email, planName string) error {
	emailData := EmailData{
		To:      email,
		Subject: "Service expiry reminder",
		Data: map[string]any{
			"Username": username,
			"PlanName": planName,
		},
	}

	return s.SendEmail(common.PlanExpiryReminderTemplate, emailData)
}
