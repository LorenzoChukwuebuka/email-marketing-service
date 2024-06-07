package observers

import "email-marketing-service/api/dto"

// Event represents an event with relevant data.
type Event struct {
	Type         string
	Message      string
	EmailRequest *dto.EmailRequest
}

// Observer interface for event listeners.
type Observer interface {
	Notify(event Event)
}
