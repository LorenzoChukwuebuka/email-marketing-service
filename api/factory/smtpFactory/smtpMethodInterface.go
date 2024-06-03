package smtpfactory

import "email-marketing-service/api/dto"

type SmtpMethodInterface interface {
	SendMail(userId int) error
	HandleSendMail(d dto.EmailRequest) error
}
