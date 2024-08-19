package custom

import (
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/utils"
	"strings"
)

type Mail struct {
}

var (
	config = utils.LoadEnv()
)

func (m *Mail) SignUpMail(email string, username string, userId string, otp string) error {
	mailTemplate := `
	<html>
	<body style="font-family: Arial, sans-serif;">
		<h2>Hi .Username ,</h2>
		<p>Thank you for registering with our service. Please use the link to verify your email address and complete your account setup:</p>
	   <p><a href="http://localhost:5054/auth/account-verification?email=.Email&token=.Token&userId=.UserId" style="padding: 10px 20px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 5px;">Verify Account </a></p>
		<p>If you did not attempt to register with our service, please ignore this email.</p>
		<br>
		<p>Regards,<br> .AppName </p>
	</body>
</html>
`
	//replace placeholders

	replacements := map[string]string{
		".Username": username,
		".Token":    otp,
		".AppName":  config.APPName,
		".UserId":   userId,
		".Email":    email,
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

func (m *Mail) ResetPasswordMail(email string, username string, otp string) error {

	mailTemplate :=
		`<html>
    <body style="font-family: Arial, sans-serif;">
        <h2>Hello .Username,</h2>
       <p>  
         We’ve received a request to change your password. If it was you, please confirm the password change. </p>
        <p>Click the link below to reset your password:</p>
        <p><a href="http://localhost:5174/auth/reset-password?email=.Email&token=.Token" style="padding: 10px 20px; background-color: #4CAF50; color: white; text-decoration: none; border-radius: 5px;">Reset Password</a></p>
        <p>If you did not attempt to reset your password, please ignore this email.</p>
		<br>
        <p>Regards,<br>.AppName</p>
    </body>
</html>
`
	replacements := map[string]string{
		".Username": username,
		".Token":    otp,
		".Email":    email,
		".AppName":  config.APPName,
	}

	formattedMail := mailTemplate

	for placeholder, value := range replacements {
		formattedMail = strings.Replace(formattedMail, placeholder, value, -1)
	}

	err := utils.SendMail("Password Reset", email, formattedMail)

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
