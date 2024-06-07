package observers

 

// Event represents an event with relevant data.
type Event struct {
	Type         string
	Message      string
	Payload interface{}
}

// Observer interface for event listeners.
type Observer interface {
	Notify(event Event)
}
