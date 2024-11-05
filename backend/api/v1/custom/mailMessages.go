package custom

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/utils"
)

// EmailConfig holds configuration for email service
type EmailConfig struct {
	Sender  string
	AppName string
	BaseURL string
}

// TemplateData represents a generic structure for email template replacements
type TemplateData map[string]string

// EmailService manages email-related operations
type EmailService struct {
	config EmailConfig
	wg     *sync.WaitGroup
}

// NewEmailService creates a new instance of EmailService
func NewEmailService() *EmailService {
	config := utils.LoadEnv()
	return &EmailService{
		config: EmailConfig{
			Sender:  config.SENDER,
			AppName: config.APPName,
			BaseURL: "http://localhost:5054",
		},
		wg: &sync.WaitGroup{},
	}
}

// loadTemplate reads and returns the content of an email template
func (s *EmailService) loadTemplate(templateName string) (string, error) {
	templatePath := filepath.Join("api", "v1", "templates", templateName)
	mailTemplate, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template %s: %v", templateName, err)
	}
	return string(mailTemplate), nil
}

// formatEmailTemplate replaces placeholders in the email template
func formatEmailTemplate(template string, data TemplateData) string {
	formattedMail := template
	for placeholder, value := range data {
		formattedMail = strings.Replace(formattedMail, placeholder, value, -1)
	}
	return formattedMail
}

// sendEmail sends an email asynchronously
func (s *EmailService) sendEmail(subject, recipient, body string) error {
	return utils.AsyncSendMail(subject, recipient, body, s.config.Sender, nil, s.wg)
}

// SignUpMail sends a sign-up verification email
func (s *EmailService) SignUpMail(email, username, userID, otp string) error {
	template, err := s.loadTemplate("verifyuser.templ")
	if err != nil {
		return err
	}

	templateData := TemplateData{
		"{{Link}}":     fmt.Sprintf("%s/auth/account-verification", s.config.BaseURL),
		"{{Username}}": username,
		"{{Token}}":    otp,
		"{{AppName}}":  s.config.AppName,
		"{{UserId}}":   userID,
		"{{Email}}":    email,
	}

	formattedMail := formatEmailTemplate(template, templateData)
	return s.sendEmail("Email Verification", email, formattedMail)
}

// ResetPasswordMail sends a password reset email
func (s *EmailService) ResetPasswordMail(email, username, otp string) error {
	template, err := s.loadTemplate("resetpassword.templ")
	if err != nil {
		return err
	}

	templateData := TemplateData{
		"{{Link}}":     fmt.Sprintf("%s/auth/reset-password", s.config.BaseURL),
		"{{Username}}": username,
		"{{Token}}":    otp,
		"{{Email}}":    email,
		"{{AppName}}":  s.config.AppName,
	}

	formattedMail := formatEmailTemplate(template, templateData)
	return s.sendEmail("Password Reset", email, formattedMail)
}

// VerifySenderMail sends a sender verification email
func (s *EmailService) VerifySenderMail(username, userEmail, domainEmail, otp, userID string) error {
	template, err := s.loadTemplate("verifysender.templ")
	if err != nil {
		return err
	}

	templateData := TemplateData{
		"{{Username}}":         username,
		"{{UserEmail}}":        userEmail,
		"{{DomainEmail}}":      domainEmail,
		"{{VerificationLink}}": fmt.Sprintf("%s/verifysender", s.config.BaseURL),
		"{{Token}}":            otp,
		"{{UserId}}":           userID,
	}

	formattedMail := formatEmailTemplate(template, templateData)
	return s.sendEmail("Verify a new Sender [Crabmailer]", domainEmail, formattedMail)
}

// DeviceVerificationMail sends a device verification email
func (s *EmailService) DeviceVerificationMail(username, email string, session *model.UserSession, code string) error {
	template, err := s.loadTemplate("planexpiry.templ")
	if err != nil {
		return err
	}

	templateData := TemplateData{
		".Username":   username,
		".Token":      code,
		".Device":     *session.Device,
		".Browser":    *session.Browser,
		".IP Address": *session.IPAddress,
		".AppName":    s.config.AppName,
	}

	formattedMail := formatEmailTemplate(template, templateData)
	return s.sendEmail("Email Verification", email, formattedMail)
}

// SubscriptionExpiryMail sends a subscription expiry notification
func (s *EmailService) SubscriptionExpiryMail(username, email, planName string) error {
	template, err := s.loadTemplate("planexpiry.templ")
	if err != nil {
		return err
	}

	templateData := TemplateData{
		"{{Username}}": username,
		"{{PlanName}}": planName,
		"{{AppName}}":  s.config.AppName,
	}

	formattedMail := formatEmailTemplate(template, templateData)
	return s.sendEmail("Subscription Expiry Notification", email, formattedMail)
}

// SubscriptionExpiryReminder sends a subscription expiry reminder
func (s *EmailService) SubscriptionExpiryReminder(username, email, planName string) error {
	mailTemplate := `
	<html>
		<body style="font-family: Arial, sans-serif;">
			<h2>Hi .Username ,</h2>
			<p>Please note that your .PlanName will expire in 5 days</p>
			<p>Regards,<br>  .Appname </p>
		</body>
	</html>
	`

	templateData := TemplateData{
		".Username": username,
		".PlanName": planName,
		".AppName":  s.config.AppName,
	}

	formattedMail := formatEmailTemplate(mailTemplate, templateData)
	return s.sendEmail("Service expiry reminder", email, formattedMail)
}
