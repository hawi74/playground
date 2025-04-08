package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	fmt.Println("Hello, world!")

	id1 := "123"
	id2 := "444"
	eid1 := "666"
	eid2 := "777"

	streams := map[string][]Event{
		"sl-123": {
			ShoppingListCreated{ID: id1, Name: "first list"},
			ShoppingListNameChanged{Name: "Original list"},
			EntryAdded{ID: eid1, Name: "my item"},
			EntryChecked{ID: eid1},
			EntryUnchecked{ID: eid1},
		},
		"sl-444": {
			ShoppingListCreated{ID: id2, Name: "second list"},
			EntryAdded{ID: eid2, Name: "some item"},
			EntryChecked{ID: eid2},
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
	uncommittedChanges []Event
}

func (s *ShoppingList) Print() {
	bytes, _ := json.MarshalIndent(s, "", "\t")
	fmt.Printf("%s\n", bytes)
}

func NewShoppingList(events ...Event) *ShoppingList {
	sl := &ShoppingList{
		Entries: make(map[string]Entry),
	}
	sl.rehydrate(events)
	return sl
}

func (s *ShoppingList) rehydrate(events []Event) error {
	for _, ev := range events {
		if err := s.apply(ev); err != nil {
			return err
		}
	}
	return nil
}

func (s *ShoppingList) ApplyChange(ev Event) error {
	s.uncommittedChanges = append(s.uncommittedChanges, ev)
	return s.apply(ev)
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
