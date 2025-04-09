package main

import (
	"context"
	"encoding/json"
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")

	id1 := "123"
	id2 := "444"
	eid1 := "666"
	eid2 := "777"

	streams := map[string][]EventEnvelope{
		"sl-123": {
			EventEnvelope{StreamID: "sl-123", Event: ShoppingListCreated{ID: id1, Name: "first list"}},
			EventEnvelope{StreamID: "sl-123", Event: ShoppingListNameChanged{Name: "Original list"}},
			EventEnvelope{StreamID: "sl-123", Event: EntryAdded{ID: eid1, Name: "my item"}},
			EventEnvelope{StreamID: "sl-123", Event: EntryChecked{ID: eid1}},
			EventEnvelope{StreamID: "sl-123", Event: EntryUnchecked{ID: eid1}},
		},
		"sl-444": {
			EventEnvelope{StreamID: "sl-444", Event: ShoppingListCreated{ID: id2, Name: "second list"}},
			EventEnvelope{StreamID: "sl-444", Event: EntryAdded{ID: eid2, Name: "some item"}},
			EventEnvelope{StreamID: "sl-444", Event: EntryChecked{ID: eid2}},
		},
	}

	for streamID, events := range streams {
		sl := NewShoppingList(events...)
		fmt.Printf("shoppingList from stream %s", streamID)
		sl.Print()

		err := sl.ChangeName("Original list")
		if err != nil {
			fmt.Println("canont update shopping list name", sl.ID, err.Error())
		}
		fmt.Printf("new changes for %s: %#v\n", sl.ID, sl.uncommittedChanges)
	}
}

type Entry struct {
	ID      string
	Name    string
	Checked bool
}

type ShoppingList struct {
	ID                 string
	Name               string
	Entries            map[string]Entry
	uncommittedChanges []EventEnvelope
}

type Repository interface {
	Save(ctx context.Context, sl ShoppingList) error
	Get(ctx context.Context, id string) (ShoppingList, error)
}

type EventBus interface {
	Send(context.Context, Event) error
}

func (s *ShoppingList) Print() {
	bytes, _ := json.MarshalIndent(s, "", "\t")
	fmt.Printf("%s\n", bytes)
}

func NewShoppingList(events ...EventEnvelope) *ShoppingList {
	sl := &ShoppingList{
		Entries: make(map[string]Entry),
	}
	sl.rehydrate(events)
	return sl
}

func (s *ShoppingList) rehydrate(events []EventEnvelope) error {
	for _, ev := range events {
		if err := s.apply(ev.Event); err != nil {
			return err
		}
	}
	return nil
}

func (s *ShoppingList) ApplyChange(ev EventEnvelope) error {
	s.uncommittedChanges = append(s.uncommittedChanges, ev)
	return s.apply(ev.Event)
}

func (s *ShoppingList) apply(ev Event) error {
	return ev.Accept(s)
}

func (s *ShoppingList) VisitShoppingListCreated(ev ShoppingListCreated) error {
	s.Name = ev.Name
	s.ID = ev.ID
	return nil
}

func (s *ShoppingList) VisitShoppingListNameChanged(ev ShoppingListNameChanged) error {
	s.Name = ev.Name
	return nil
}

func (s *ShoppingList) VisitEntryAdded(ev EntryAdded) error {
	s.Entries[ev.ID] = Entry{ID: ev.ID, Name: ev.Name, Checked: false}
	return nil
}

func (s *ShoppingList) VisitEntryChecked(ev EntryChecked) error {
	e, ok := s.Entries[ev.ID]
	if !ok {
		return fmt.Errorf("entry not found")
	}
	e.Checked = true
	s.Entries[ev.ID] = e
	return nil
}

func (s *ShoppingList) VisitEntryUnchecked(ev EntryUnchecked) error {
	e, ok := s.Entries[ev.ID]
	if !ok {
		return fmt.Errorf("entry not found")
	}
	e.Checked = false
	s.Entries[ev.ID] = e
	return nil
}

func (s *ShoppingList) ChangeName(newName string) error {
	if newName == s.Name {	
		return fmt.Errorf("new name is the same as current name")
	}
	ev := ShoppingListNameChanged{
		Name: newName,
	}
	s.ApplyChange(ev)
	return nil
}
