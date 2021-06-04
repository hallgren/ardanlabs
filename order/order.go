package order

import (
	"fmt"
	"time"

	"github.com/hallgren/eventsourcing"
)

// State current state of the order
type State int

const (
	// Empty no items in order
	Empty State = iota
	// Ongoing is when items in order but not paid for
	Ongoing
	// PaidFor items are paid for
	PaidFor
	// Canceled - the order has been deleted
	Canceled
)

// Order an order
type Order struct {
	eventsourcing.AggregateRoot
	State     State
	Items     map[int]struct{}
	UpdatedAt time.Time
	CreatedAt time.Time
}

// Created is the initial event on the order
type Created struct{}

// ItemAdded event
type ItemAdded struct {
	ItemID int
}

// ItemRemoved event
type ItemRemoved struct {
	ItemID int
}

// PaidWithCreditCard when order is paid for with credit card
type PaidWithCreditCard struct{}

// Deleted when the order has been deleted
type Deleted struct{}

// Transition transform the order
func (o *Order) Transition(event eventsourcing.Event) {
	switch e := event.Data.(type) {
	case *Created:
		o.State = Empty
		o.Items = make(map[int]struct{})
		o.UpdatedAt = event.Timestamp
		o.CreatedAt = event.Timestamp
	case *ItemAdded:
		o.State = Ongoing
		o.Items[e.ItemID] = struct{}{}
		o.UpdatedAt = event.Timestamp
	case *ItemRemoved:
		o.UpdatedAt = event.Timestamp
		delete(o.Items, e.ItemID)
		if len(o.Items) == 0 {
			o.State = Empty
		}
	case *PaidWithCreditCard:
		o.State = PaidFor
	case *Deleted:
		o.State = Canceled
	}
}

// Create is the constructor
func Create() *Order {
	order := Order{}
	order.TrackChange(&order, &Created{})
	return &order
}

// AddItem adds item if in correct state
func (o *Order) AddItem(itemID int) error {
	if o.State != Ongoing && o.State != Empty {
		return fmt.Errorf("order in wrong state")
	}
	_, ok := o.Items[itemID]
	if ok {
		return fmt.Errorf("item already present")
	}
	o.TrackChange(o, &ItemAdded{ItemID: itemID})
	return nil
}

// RemoveItem remove item or error if not present
func (o *Order) RemoveItem(itemID int) error {
	if o.State != Ongoing {
		return fmt.Errorf("order in wrong state")
	}
	_, ok := o.Items[itemID]
	if !ok {
		return fmt.Errorf("item not present")
	}
	o.TrackChange(o, &ItemRemoved{ItemID: itemID})
	return nil
}

// PayWithCreditCard pays for the order
func (o *Order) PayWithCreditCard() error {
	if o.State != Ongoing {
		return fmt.Errorf("order in wrong state")
	}
	o.TrackChange(o, &PaidWithCreditCard{})
	return nil
}

// Delete deletes the order
func (o *Order) Delete() error {
	if o.State != Empty && o.State != Ongoing {
		return fmt.Errorf("order in wrong state")
	}
	o.TrackChange(o, &Deleted{})
	return nil
}
