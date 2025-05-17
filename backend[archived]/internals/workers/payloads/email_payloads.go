package payloads

type EmailPayload struct {
	Subject string `json:"subject"`
	Email   string `json:"email"`
	Message string `json:"message"`
	Sender  string `json:"sender"`
}
