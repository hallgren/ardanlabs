package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/hallgren/eventsourcing"
	sqles "github.com/hallgren/eventsourcing/eventstore/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Person struct {
	eventsourcing.AggregateRoot
	Name string
	Age  int
}

// Transition the person state dependent on the events
func (person *Person) Transition(event eventsourcing.Event) {
	switch e := event.Data.(type) {
	case *Born:
		person.Age = 0
		person.Name = e.Name
	case *AgedOneYear:
		person.Age += 1
	}
}

// Initial event
type Born struct {
	Name string
}

// Event that happens once a year
type AgedOneYear struct{}

// CreatePerson constructor for Person
func CreatePerson(name string) (*Person, error) {
	if name == "" {
		return nil, errors.New("name can't be blank")
	}
	person := Person{}
	person.TrackChange(&person, &Born{Name: name})
	return &person, nil
}

// GrowOlder command
func (person *Person) GrowOlder() {
	person.TrackChange(person, &AgedOneYear{})
}

func main() {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}
	ser := eventsourcing.NewSerializer(json.Marshal, json.Unmarshal)
	ser.RegisterTypes(&Person{},
		func() interface{} { return &Born{} },
		func() interface{} { return &AgedOneYear{} },
	)
	sqlEventStore := sqles.Open(db, *ser)
	if err != nil {
		panic(err)
	}
	err = sqlEventStore.Migrate()
	if err != nil {
		//panic(err)
		fmt.Println(err)
	}
	// Create a repo to handle event sourced
	repo := eventsourcing.NewRepository(sqlEventStore, nil)

	p, err := CreatePerson("Morgan")
	if err != nil {
		panic(err)
	}
	p.GrowOlder()
	p.GrowOlder()
	p.GrowOlder()
	p.GrowOlder()
	p.GrowOlder()

	fmt.Println("person age", p.Age)
	spew.Dump(p.Events())

	repo.Save(p)
}
