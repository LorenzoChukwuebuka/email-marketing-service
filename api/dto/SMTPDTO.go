package dto

import "email-marketing-service/api/model"

type EmailRequest struct {
	Sender      Sender      `json:"sender"`
	To          []Recipient `json:"to"`
	Subject     string      `json:"subject"`
	HtmlContent *string     `json:"htmlContent"`
	Text        *string     `json:"text"`
}

type Sender struct {
	Email string  `json:"email"`
	Name  *string `json:"name"`
}

type Recipient struct {
	Email string `json:"email"`
}

type SentEmails struct {
	Sender         uint
	Recipient      uint
	MessageContent string
	Status         model.EmailStatus
}
