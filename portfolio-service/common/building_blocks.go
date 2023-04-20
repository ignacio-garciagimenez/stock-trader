package common

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Handler[K any, V any] interface {
	Handle(context.Context, K) (V, error)
}

type DomainEvent interface {
	Id() string
	Name() string
	Timestamp() time.Time
	EventData() map[string]any
}

type DomainEventEntity struct {
	Id        string            `gorm:"column:id"`
	Timestamp time.Time         `gorm:"column:timestamp"`
	Name      string            `gorm:"column:name"`
	EventData datatypes.JSONMap `gorm:"column:event_data"`
}

func (DomainEventEntity) TableName() string {
	return "event_journal"
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
