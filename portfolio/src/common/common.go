package common

import (
	"errors"
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

type AggregateRoot[K comparable] interface {
	Id() K
	DomainEvents() []DomainEvent
	ClearDomainEvents()
}

type Repository[K comparable, V AggregateRoot[K]] interface {
	FindById(K) (V, error)
	Save(V) error
}

type InMemoryBaseRepository[K comparable, V AggregateRoot[K]] struct {
	entities     map[K]V
	domainEvents []DomainEvent
}

func (r *InMemoryBaseRepository[K, V]) FindById(key K) (V, error) {
	var entity V
	if entity, found := r.entities[key]; found {
		return entity, nil
	}
	return entity, errors.New("entity not found")

}

func (r *InMemoryBaseRepository[K, V]) Save(entity V) error {
	r.entities[entity.Id()] = entity

	//For Transactional Outbox implementation
	r.domainEvents = append(r.domainEvents, entity.DomainEvents()...)

	return nil
}

func NewInMemoryBaseRepository[K comparable, V AggregateRoot[K]]() *InMemoryBaseRepository[K, V] {
	return &InMemoryBaseRepository[K, V]{
		entities:     map[K]V{},
		domainEvents: []DomainEvent{},
	}
}
