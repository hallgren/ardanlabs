package main

import (
	"fmt"
	"github.com/hallgren/eventsourcing"
)

// Payed event
type Payed struct {
	From   string
	To     string
	Amount float64
}

// Payment entity
type Payment struct {
	eventsourcing.AggregateRoot
	Amount float64
}

// Transition build the Payment aggregate
func (p *Payment) Transition(event eventsourcing.Event) {
	switch e := event.Data.(type) {
	case *Payed:
		p.Amount = e.Amount
	}

}

// Pay construtor
func Pay(to, from string, amount float64) *Payment {
	p := &Payment{}
	p.SetID("1")
	p.TrackChange(p, &Payed{From: from, To: to, Amount: amount})
	return p
}

func main() {
	fmt.Println("vim-go")
	p := Pay("bill", "morgan", 100)
	p2 := Pay("bill", "morgan", 100)
	fmt.Println(p.Events())
	fmt.Println(p.Amount)
	fmt.Println(p2.Events())
	fmt.Println(p2.Amount)
}
