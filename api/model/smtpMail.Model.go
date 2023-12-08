package model

type EmailRequest struct {
	Sender      Sender     `json:"sender"`
	To          []Recipient `json:"to"`
	Subject     string     `json:"subject"`
	HtmlContent string     `json:"htmlContent"`
}

type Sender struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Recipient struct {
	Email string `json:"email"`
}