package smtpfactory

import "email-marketing-service/api/dto"

type SmtpMethodInterface interface {
	HandleSendMail(d *dto.EmailRequest) error
}
