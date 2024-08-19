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

	err = utils.SendMail("Email Verification", email, formattedMail)

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

	err = utils.SendMail("Password Reset", email, formattedMail)

	if err != nil {
		return err
	}
	return nil
}

func (m *Mail) DeviceVerificationMail(username string, email string, d *model.UserSession, code string) error {
	mailTemplate := `
	<html>
	<body style="font-family: Arial, sans-serif;">
		<h2>Hi .Username ,</h2>
		<p>Device: .Device </p>
		<p>Browser: .Browser </p>
		<p>IP Address: .IP </p>
		<p>Use the code below to verify your device:</p>
		<h3>Code: .Token </h3>
		<p>Please note that this code can only be used once and is valid for a limited time.</p>
		
		<br>
		<p>Regards,<br> .AppName </p>
	</body>
</html>
`
	//replace placeholders

	replacements := map[string]string{
		".Username":   username,
		".Token":      code,
		".Device":     *d.Device,
		".Browser":    *d.Browser,
		".IP Address": *d.IPAddress,
		".AppName":    config.APPName,
	}

	formattedMail := mailTemplate

	for placeholder, value := range replacements {
		formattedMail = strings.Replace(formattedMail, placeholder, value, -1)
	}

	err := utils.SendMail("Email Verification", email, formattedMail)

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
			<p>Please note that your .PlanName has expired</p>
			
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

	err := utils.SendMail("Subscription Expiry Notification", email, formattedMail)

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

	err := utils.SendMail("Service expiry reminder", email, formattedMail)

	if err != nil {
		return err
	}
	return nil
}
