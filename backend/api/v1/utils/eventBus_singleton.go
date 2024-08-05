package utils

import (
	"email-marketing-service/api/v1/observers"
	"sync"
)

var (
	instance *EventBus
	once     sync.Once
)

func GetEventBus() *EventBus {
	once.Do(func() {
		instance = &EventBus{
			observers: make(map[string][]observers.Observer),
		}
	})
	return instance
}
