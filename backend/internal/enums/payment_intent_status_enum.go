package enums

type PaymentIntentStatus string

const (
	PaymentIntentProcessing PaymentIntentStatus = "processing"
	PaymentIntentSuccessful PaymentIntentStatus = "successful"
	PaymentIntentFailed     PaymentIntentStatus = "failed"
)
