package smtpfactory

import "email-marketing-service/api/v1/dto"

type SmtpMethodInterface interface {
	HandleSendMail(d *dto.EmailRequest) error
}
