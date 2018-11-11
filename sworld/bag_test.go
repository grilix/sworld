package sworld

import (
	"testing"
)

func TestStoringItems(t *testing.T) {
	bag := NewStandardBag(10)
	stone := &PortalStone{}
	err := bag.StoreItem(stone, 0)
	if err != nil {
		t.Fatal(err)
	}

	items := bag.Items()
	if len(items) != 10 {
		t.Error("Bag is expected to have 10 items, got ", len(items))
	}

	stoneItem, ok := items[0].(*PortalStone)
	if !ok {
		t.Error("First item was expected to be a stone")
	}

	if stoneItem.Level != stone.Level {
		t.Error("First item was expected to be the stored stone")
	}
}

func TestDropItem(t *testing.T) {
	zone := &Zone{}
	stone1 := PortalStone{Zone: zone}
	stone2 := PortalStone{Zone: zone}

	bag := NewStandardBag(3)
	err := bag.StoreItem(&stone1, 0)
	if err != nil {
		t.Fatal(err)
	}
	err = bag.StoreItem(&stone2, 1)
	if err != nil {
		t.Fatal(err)
	}
	item, err := bag.DropItem(0)
	if err != nil {
		t.Fatal(err)
	}
	if item == nil {
		t.Error("Expected drop to return the item")
	}
	item, err = bag.DropItem(1)
	if err != nil {
		t.Fatal(err)
	}
	if item == nil {
		t.Error("Expected drop to return the item")
	}
	stored, err := bag.GetItem(0)
	if err == nil {
		t.Error("Expected drop to remove the item, got:", stored)
	}
}
