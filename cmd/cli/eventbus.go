package main

import "context"

type InMemoryEventBus struct {
	listeners []Visitor
}

func NewInMemoryEventBus() *InMemoryEventBus {
	return &InMemoryEventBus{
		listeners: make([]Visitor, 0),
	}
}

func (eb *InMemoryEventBus) RegisterListener(l Visitor) error {
	eb.listeners = append(eb.listeners, l)
	return nil
}

func (eb *InMemoryEventBus) Send(ctx context.Context, ev Event) error {
	for _, l := range eb.listeners {
		if err := ev.Accept(l); err != nil {
			return err
		}
	}
	return nil
}
