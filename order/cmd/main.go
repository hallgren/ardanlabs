package main

import (
	"ardanlabs/order"
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
	ser.RegisterTypes(&order.Order{},
		func() interface{} { return &order.Created{} },
		func() interface{} { return &order.ItemAdded{} },
		func() interface{} { return &order.ItemRemoved{} },
		func() interface{} { return &order.PaidWithCreditCard{} },
		func() interface{} { return &order.Deleted{} },
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

	o := order.Create()
	o.AddItem(123)
	o.AddItem(456)
	o.Delete()
	//spew.Dump(o)
	repo.Save(o)

	orderCopy := order.Order{}
	repo.Get(o.ID(), &orderCopy)
	spew.Dump(orderCopy)

	/*
		order := order.Order{}
		repo.Get("jDQPGoPqDbrCzL72ScD9", &order)
		order.Delete()
		spew.Dump(order)
	*/
}
