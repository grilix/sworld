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
