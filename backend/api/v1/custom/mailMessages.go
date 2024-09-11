package custom

import (
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/utils"
	"os"
	"path/filepath"
	"strings"
)

type Mail struct {
}

var (
	config = utils.LoadEnv()
	sender = "noreply@crabmailer.app"
)

func (m *Mail) SignUpMail(email string, username string, userId string, otp string) error {
	// Read the template file
	templatePath := filepath.Join("api", "v1", "templates", "verifyuser.templ")
	mailTemplate, err := os.ReadFile(templatePath)
	if err != nil {
		return err
	}

	replacements := map[string]string{
		"{{Link}}":     "http://localhost:5054/auth/account-verification",
		"{{Username}}": username,
		"{{Token}}":    otp,
		"{{AppName}}":  config.APPName,
		"{{UserId}}":   userId,
		"{{Email}}":    email,
	}

	formattedMail := string(mailTemplate)

	for placeholder, value := range replacements {
		formattedMail = strings.Replace(formattedMail, placeholder, value, -1)
	}

	err = utils.SendMail("Email Verification", email, formattedMail, sender,nil)

	if err != nil {
		return err
	}

	return nil
}

func (m *Mail) ResetPasswordMail(email string, username string, otp string) error {

	// Read the template file
	templatePath := filepath.Join("api", "v1", "templates", "resetpassword.templ")
	mailTemplate, err := os.ReadFile(templatePath)
	if err != nil {
		return err
	}

	replacements := map[string]string{
		"{{Link}}":     "http://localhost:5054/auth/reset-password",
		"{{Username}}": username,
		"{{Token}}":    otp,
		"{{Email}}":    email,
		"{{AppName}}":  config.APPName,
	}

	formattedMail := string(mailTemplate)

	for placeholder, value := range replacements {
		formattedMail = strings.Replace(formattedMail, placeholder, value, -1)
	}

	err = utils.SendMail("Password Reset", email, formattedMail, sender,nil)

	if err != nil {
		return err
	}
	return nil
}

func (m *Mail) VerifySenderMail() error {
	return nil
}

func (m *Mail) DeviceVerificationMail(username string, email string, d *model.UserSession, code string) error {

	templatePath := filepath.Join("api", "v1", "templates", "planexpiry.templ")
	mailTemplate, err := os.ReadFile(templatePath)
	if err != nil {
		return err
	}
	//replace placeholders

	replacements := map[string]string{
		".Username":   username,
		".Token":      code,
		".Device":     *d.Device,
		".Browser":    *d.Browser,
		".IP Address": *d.IPAddress,
		".AppName":    config.APPName,
	}

	formattedMail := string(mailTemplate)

	for placeholder, value := range replacements {
		formattedMail = strings.Replace(formattedMail, placeholder, value, -1)
	}

	err = utils.SendMail("Email Verification", email, formattedMail, sender,nil)

	if err != nil {
		return err
	}

	return nil
}

func (m *Mail) SubscriptionExpiryMail(username string, email string, planName string) error {
	mailTemplate :=
		`<html>
		<body style="font-family: Arial, sans-serif;">
			<h2>Hi .Username ,</h2>
			<p>Please note that your .PlanName plan has expired</p>
			
			<p>Regards,<br>  .Appname </p>
		</body>
		</html>
       `
	replacements := map[string]string{
		"{{Username}}": username,
		"{{PlanName}}": planName,
		"{{AppName}}":  config.APPName,
	}

	formattedMail := mailTemplate

	for placeholder, value := range replacements {
		formattedMail = strings.Replace(formattedMail, placeholder, value, -1)
	}

	err := utils.SendMail("Subscription Expiry Notification", email, formattedMail, sender,nil)

	if err != nil {
		return err
	}
	return nil
}

func (m *Mail) SubscriptionExpiryReminder(username string, email string, planName string) error {
	mailTemplate := `
	<html>
		<body style="font-family: Arial, sans-serif;">
			<h2>Hi .Username ,</h2>
			<p>Please note that your .PlanName will expire in 5 days</p>
			<p>Regards,<br>  .Appname </p>
		</body>
		</html>
	`

	replacements := map[string]string{
		".Username": username,
		".PlanName": planName,
		".AppName":  config.APPName,
	}

	formattedMail := mailTemplate

	for placeholder, value := range replacements {
		formattedMail = strings.Replace(formattedMail, placeholder, value, -1)
	}

	err := utils.SendMail("Service expiry reminder", email, formattedMail, sender,nil)

	if err != nil {
		return err
	}
	return nil
}
