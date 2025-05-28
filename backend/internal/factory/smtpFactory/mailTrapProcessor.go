package smtpfactory

import (
	"email-marketing-service/internal/domain"
	"email-marketing-service/internal/helper"
	"fmt"
	"sync"
)

var wg sync.WaitGroup

type MailTrapProcessor struct {
}

func (s *MailTrapProcessor) HandleSendMail(emailRequest *domain.EmailRequest) error {
	switch to := emailRequest.To.(type) {
	case domain.Recipient:
		s.sendMailToRecipient(to, emailRequest)
	case []domain.Recipient:
		for _, recipient := range to {
			s.sendMailToRecipient(recipient, emailRequest)
		}
	default:
		return fmt.Errorf("invalid recipient type")
	}

	return nil
}

func (s *MailTrapProcessor) sendMailToRecipient(recipient domain.Recipient, emailRequest *domain.EmailRequest) {
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
	sender := emailRequest.Sender.Email

	err := helper.AsyncSendMail(subject, email, mailContent, sender, nil, &wg)
	if err != nil {
		fmt.Printf("Error sending mail to %s: %v\n", email, err)
	} else {

		fmt.Printf("Mail sent to %s\n", email)
	}
}
