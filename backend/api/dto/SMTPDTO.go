package dto

import (
	"email-marketing-service/api/model"
	"encoding/json"
	"errors"
)

type EmailRequest struct {
	Sender      Sender      `json:"sender"`
	To          interface{} `json:"to"`
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

// Custom unmarshaler for EmailRequest to handle different types for the To field.
func (e *EmailRequest) UnmarshalJSON(data []byte) error {
	type Alias EmailRequest
	aux := &struct {
		To json.RawMessage `json:"to"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var recipient Recipient
	if err := json.Unmarshal(aux.To, &recipient); err == nil {
		e.To = recipient
		return nil
	}

	var recipients []Recipient
	if err := json.Unmarshal(aux.To, &recipients); err == nil {
		e.To = recipients
		return nil
	}

	return errors.New("invalid type for To field")
}
