package smtpfactory

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/utils"
	"fmt"
)

type SMTPProcessor struct {
}

var (
	config = utils.LoadEnv()
)

func (s *SMTPProcessor) HandleSendMail(emailRequest *dto.EmailRequest) error {
	switch to := emailRequest.To.(type) {
	case dto.Recipient:
		s.sendMailToRecipient(to, emailRequest)
	case []dto.Recipient:
		for _, recipient := range to {
			s.sendMailToRecipient(recipient, emailRequest)
		}
	default:
		return fmt.Errorf("invalid recipient type")
	}

	return nil
}

func (s *SMTPProcessor) sendMailToRecipient(recipient dto.Recipient, emailRequest *dto.EmailRequest) {
	// Determine the mail content (HTML or Text)
	var mailContent string
	if emailRequest.HtmlContent != nil {
		mailContent = *emailRequest.HtmlContent
	} else if emailRequest.Text != nil {
		mailContent = *emailRequest.Text
	} else {
		fmt.Println("No content to send")
		return
	}

	email := recipient.Email
	subject := emailRequest.Subject
	sender := emailRequest.Sender.Email

	smtpConfig := utils.SMTPConfig{
		Host:     config.SMTP_SERVER,
		Port:     1025,
		Username: emailRequest.AuthUser.Username,
		Password: emailRequest.AuthUser.Password,
	}

	err := utils.AsyncSendMail(subject, email, mailContent, sender, &smtpConfig, &wg)
	if err != nil {
		fmt.Printf("Error sending mail to %s: %v\n", email, err)
	} else {

		fmt.Printf("Mail sent to %s\n", email)
	}
}
