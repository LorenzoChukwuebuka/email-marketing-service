package custom

import (
	"email-marketing-service/api/utils"
	"strings"
)

type Mail struct {
}

func (m *Mail) SignUpMail(email string, username string, otp string) error {
	mailTemplate := `
	<html>
	<body style="font-family: Arial, sans-serif;">
		<h2>Hi .Username ,</h2>
		<p>Thank you for registering with our service. Please use the following One-Time Password (OTP) to verify your email address and complete your account setup:</p>
		<h3>OTP: .Token </h3>
		<p>Please note that this OTP can only be used once and is valid for a limited time.</p>
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
		".AppName":  "Appname",
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
        <h2>Hi .Username ,</h2>
        <p>Please use the following One-Time Password (OTP) to reset your password:</p>
        <h3>OTP:  .Token </h3>
        <p>Please note that this OTP can only be used once and is valid for a limited time.</p>
        <p>If you did not attempt to reset your password, please ignore this email.</p>
        <br>
        <p>Regards,<br>  .Appname </p>
    </body>
</html>
`
	replacements := map[string]string{
		".Username": username,
		".Token":    otp,
		".AppName":  "Appname",
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
		".AppName":  "Appname",
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
		".AppName":  "Appname",
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
