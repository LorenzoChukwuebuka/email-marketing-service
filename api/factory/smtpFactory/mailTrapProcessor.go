package smtpfactory

import (
	"email-marketing-service/api/dto"
	"email-marketing-service/api/utils"
	"fmt"
)

type MailTrapProcessor struct {
}

func (s *MailTrapProcessor) SendMail(userId int) error {
	return nil
}

func (s *MailTrapProcessor) HandleSendMail(emailRequest *dto.EmailRequest) error {
	for _, recipient := range emailRequest.To {
		// Determine the mail content (HTML or Text)
		var mailContent string
		if emailRequest.HtmlContent != nil {
			mailContent = *emailRequest.HtmlContent
		} else if emailRequest.Text != nil {
			mailContent = *emailRequest.Text
		} else {
			continue
		}

		email := recipient.Email

		subject := emailRequest.Subject

		// Send the email
		err := utils.SendMail(subject, email, mailContent)
		if err != nil {
			fmt.Printf("Error sending mail to %s: %v\n", recipient.Email, err)
		} else {
			fmt.Printf("Mail sent to %s\n", recipient.Email)
		}

	}

	return nil
}
