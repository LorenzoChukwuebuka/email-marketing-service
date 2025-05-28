package enums 

type EmailStatus string

const (
	Sending EmailStatus = "sending"
	Failed  EmailStatus = "failed"
	Success EmailStatus = "success"
)
