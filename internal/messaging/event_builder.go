package messaging

import (
	"encoding/json"
	"time"
)

type EventBuilder struct {
	event Event
}

func NewEventBuilder() *EventBuilder {
	return &EventBuilder{
		event: Event{
			Timestamp: time.Now(),
		},
	}
}

func (eb *EventBuilder) WithType(eventType string) *EventBuilder {
	eb.event.Type = eventType
	return eb
}

func (eb *EventBuilder) WithPayload(payload interface{}) *EventBuilder {
	payloadBytes, _ := json.Marshal(payload)
	eb.event.Payload = payloadBytes
	return eb
}

func (eb *EventBuilder) WithSource(source string) *EventBuilder {
	eb.event.Source = source
	return eb
}

func (eb *EventBuilder) Build() Event {
	return eb.event
}
