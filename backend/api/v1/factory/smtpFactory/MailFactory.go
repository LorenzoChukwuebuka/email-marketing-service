package smtpfactory

import "fmt"

func MailFactory(mailhost string) (SmtpMethodInterface, error) {
	var sI SmtpMethodInterface
	switch mailhost {
	case "mailtrap":
		sI = &MailTrapProcessor{}
		return sI, nil
	case "smtp":
		return nil, nil
	default:
		return nil, fmt.Errorf("invalid mail processor: %s", mailhost)
	}

}
