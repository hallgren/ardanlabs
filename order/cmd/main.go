package main

import (
	"ardanlabs"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/hallgren/eventsourcing"
	sqles "github.com/hallgren/eventsourcing/eventstore/sql"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "test.db")
	if err != nil {
		panic(err)
	}
	ser := eventsourcing.NewSerializer(json.Marshal, json.Unmarshal)
	ser.RegisterTypes(&ardanlabs.Order{},
		func() interface{} { return &ardanlabs.Created{} },
		func() interface{} { return &ardanlabs.ItemAdded{} },
		func() interface{} { return &ardanlabs.ItemRemoved{} },
		func() interface{} { return &ardanlabs.Paid{} },
		func() interface{} { return &ardanlabs.Deleted{} },
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

	o := ardanlabs.Create()
	o.AddItem(123)
	o.AddItem(456)
	o.Delete()
	//spew.Dump(o)
	repo.Save(o)

	orderCopy := ardanlabs.Order{}
	repo.Get(o.ID(), &orderCopy)
	spew.Dump(orderCopy)

	/*
		order := ardanlabs.Order{}
		repo.Get("jDQPGoPqDbrCzL72ScD9", &order)
		order.Delete()
		spew.Dump(order)
	*/
}
