package common

import (
	"time"

	"github.com/google/uuid"
)

type DomainEvent interface {
	Id() string
	Name() string
	Timestamp() time.Time
	EventData() map[string]any
}

func NewBaseDomainEvent(name string) *BaseDomainEvent {
	return &BaseDomainEvent{
		id:        uuid.NewString(),
		name:      name,
		timestamp: time.Now().UTC(),
	}
}

type BaseDomainEvent struct {
	id        string
	name      string
	timestamp time.Time
}

func (d BaseDomainEvent) Id() string {
	return d.id
}

func (d BaseDomainEvent) Name() string {
	return d.name
}

func (d BaseDomainEvent) Timestamp() time.Time {
	return d.timestamp
}

type AggregateRoot interface {
	DomainEvents() []DomainEvent
	ClearDomainEvents()
}
