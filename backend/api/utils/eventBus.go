package utils

import (
	"email-marketing-service/api/observers"
	"fmt"
	"sync"
)

// EventBus manages observers and notifies them of events.
type EventBus struct {
	observers map[string][]observers.Observer
	mu        sync.Mutex
}

// NewEventBus creates a new EventBus.
func NewEventBus() *EventBus {
	return &EventBus{
		observers: make(map[string][]observers.Observer),
	}
}

// Register adds an observer for a specific event type.
func (eb *EventBus) Register(eventType string, observer observers.Observer) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.observers[eventType] = append(eb.observers[eventType], observer)
}

// Unregister removes an observer for a specific event type.
func (eb *EventBus) Unregister(eventType string, observer observers.Observer) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	observers := eb.observers[eventType]
	for i, obs := range observers {
		if obs == observer {
			eb.observers[eventType] = append(observers[:i], observers[i+1:]...)
			break
		}
	}
}

// Notify all registered observers of an event and unregister them after.
func (eb *EventBus) Notify(event observers.Event) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	observers := eb.observers[event.Type]

	fmt.Printf("EventBus: Notifying observers for event type: %s\n", event.Type) // Debug print

	for _, obs := range observers {
		fmt.Printf("EventBus: Notifying observer\n") // Debug print
		obs.Notify(event)
	}

	// Unregister observers after notification
	fmt.Printf("EventBus: Unregistering observers for event type: %s\n", event.Type) // Debug print
	delete(eb.observers, event.Type)
}
