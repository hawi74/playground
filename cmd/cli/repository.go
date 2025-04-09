package main

import (
	"context"
	"fmt"
)

type InMemoryRepository struct {
	streams map[string][]EventEnvelope
	eventBus EventBus
}

func NewInMemoryRepository(eb EventBus) *InMemoryRepository {
	return &InMemoryRepository{
		streams: make(map[string][]EventEnvelope),
		eventBus: eb,
	}
}

func (r *InMemoryRepository) Save(ctx context.Context, sl ShoppingList) error {
	r.streams[sl.ID] = append(r.streams[sl.ID], sl.uncommittedChanges...)
	for _, ev := range sl.uncommittedChanges {
		r.eventBus.Send(ctx, ev)
	}
	return nil
}

func (r *InMemoryRepository) Get(ctx context.Context, id string) (ShoppingList, error) {
	events, ok := r.streams[id]
	if !ok {
		return ShoppingList{}, fmt.Errorf("shopping list with if %s not found", id)
	}
	return *NewShoppingList(events...), nil
}
