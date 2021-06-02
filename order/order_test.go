package ardanlabs_test

import "testing"
import "ardanlabs"

//import "github.com/davecgh/go-spew/spew"

func TestCreate(t *testing.T) {
	o := ardanlabs.Create()
	if len(o.Events()) != 1 {
		t.Fatalf("expected one event got %d", len(o.Events()))
	}
	if o.Events()[0].Reason != "Created" {
		t.Fatalf("expected event Created got %s", o.Events()[0].Reason)
	}
	if o.State != ardanlabs.Empty {
		t.Fatalf("expected order state to be Empty was %d", o.State)
	}
	//spew.Dump(o.Events())
}

func TestItemAddRemove(t *testing.T) {
	o := ardanlabs.Create()
	err := o.AddItem(1)
	if err != nil {
		t.Fatal(err)
	}
	if o.Events()[1].Reason != "ItemAdded" {
		t.Fatalf("expected event ItemAdded got %s", o.Events()[1].Reason)
	}
	if o.State != ardanlabs.Ongoing {
		t.Fatalf("expected order state ongoing was %d", o.State)
	}
	// remove none existing item
	err = o.RemoveItem(999)
	if err == nil {
		t.Fatal("expected error when removing none existing item")
	}
	err = o.RemoveItem(1)
	if err != nil {
		t.Fatal(err)
	}
	if o.Events()[2].Reason != "ItemRemoved" {
		t.Fatalf("expected event ItemRemoved got %s", o.Events()[2].Reason)
	}
	if o.State != ardanlabs.Empty {
		t.Fatalf("expected order state empty was %d", o.State)
	}
	//spew.Dump(o.Events())
}

func TestPayEmptyOrder(t *testing.T) {
	o := ardanlabs.Create()
	err := o.Pay()
	if err == nil {
		t.Fatal("expected err when paying an empty order")
	}
}

func TestDeleteEmptyOrder(t *testing.T) {
	o := ardanlabs.Create()
	err := o.Delete()
	if err != nil {
		t.Fatal("expected no err when deleting empty order")
	}
	if o.Events()[1].Reason != "Deleted" {
		t.Fatalf("expected event Deleted got %s", o.Events()[1].Reason)
	}
}
