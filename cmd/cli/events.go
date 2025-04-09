package main

type Visitor interface {
	VisitShoppingListCreated(ShoppingListCreated) error
	VisitShoppingListNameChanged(ShoppingListNameChanged) error
	VisitEntryAdded(EntryAdded) error
	VisitEntryChecked(EntryChecked) error
	VisitEntryUnchecked(EntryUnchecked) error
}

type Event interface {
	Accept(v Visitor) error
}

type ShoppingListCreated struct {
	ID   string
	Name string
}

func (ev ShoppingListCreated) Accept(v Visitor) error {
	return v.VisitShoppingListCreated(ev)
}

type ShoppingListNameChanged struct {
	Name string
}

func (ev ShoppingListNameChanged) Accept(v Visitor) error {
	return v.VisitShoppingListNameChanged(ev)
}

type EntryAdded struct {
	ID   string
	Name string
}

func (ev EntryAdded) Accept(v Visitor) error {
	return v.VisitEntryAdded(ev)
}

type EntryChecked struct {
	ID string
}

func (ev EntryChecked) Accept(v Visitor) error {
	return v.VisitEntryChecked(ev)
}

type EntryUnchecked struct {
	ID string
}

func (ev EntryUnchecked) Accept(v Visitor) error {
	return v.VisitEntryUnchecked(ev)
}

type EventEnvelope struct {
	StreamID string
	Event Event
}

