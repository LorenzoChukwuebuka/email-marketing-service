package smtpfactory

import (
	"email-marketing-service/api/v1/dto"
	//"email-marketing-service/api/observers"
	"email-marketing-service/api/v1/utils"
	"fmt"
)

type MailTrapProcessor struct {
}

func (s *MailTrapProcessor) HandleSendMail(emailRequest *dto.EmailRequest) error {
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

func (s *MailTrapProcessor) sendMailToRecipient(recipient dto.Recipient, emailRequest *dto.EmailRequest) {
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

	//eventBus := utils.GetEventBus()

	email := recipient.Email
	subject := emailRequest.Subject

	err := utils.SendMail(subject, email, mailContent)
	if err != nil {
		fmt.Printf("Error sending mail to %s: %v\n", email, err)
	} else {
		//	eventBus.Notify(observers.Event{Type: "send_success", Message: "email sent successfully", Payload: emailRequest})
		fmt.Printf("Mail sent to %s\n", email)
	}
}
