package smtpfactory

import "email-marketing-service/internal/domain"

type SmtpMethodInterface interface {
	HandleSendMail(d *domain.EmailRequest) error
}
