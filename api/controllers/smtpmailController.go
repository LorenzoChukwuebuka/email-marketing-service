package controllers

import "net/http"

type SMTPMailController struct {
}

func NewSMTPMailController() *SMTPMailController {
	return &SMTPMailController{}
}

func (c *SMTPMailController) SendSMTPMail(w http.ResponseWriter, r *http.Request) {
	// Get the value of the "api-key" header
	apiKey := r.Header.Get("api-key")
	if apiKey == "" {
		// The header is not present or has an empty value
		errorRes := map[string]interface{}{
			"status":         http.StatusUnauthorized,
			"error response": "API key not provided",
		}
		response.ErrorResponse(w, errorRes)
		return
	}

	response.SuccessResponse(w, 200, apiKey)

}
