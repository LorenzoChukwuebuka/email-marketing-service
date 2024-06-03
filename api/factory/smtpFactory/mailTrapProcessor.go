package smtpfactory

import "email-marketing-service/api/dto"

type MailTrapProcessor struct {
}

func (s *MailTrapProcessor) SendMail(userId int) error {
	print("hello world")
	return nil
}

func (s *MailTrapProcessor) HandleSendMail(d dto.EmailRequest) error {
	print ("welcome mail sent successfully")
	return nil
}
