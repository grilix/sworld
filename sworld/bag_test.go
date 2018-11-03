package sworld

import (
	"testing"
)

func TestStoringItems(t *testing.T) {
	bag := NewStandardBag(10)
	stone := &PortalStone{}
	slot, err := bag.StoreItem(stone)
	if err != nil {
		t.Fatal(err)
	}

	if slot != 0 {
		t.Error("Expected slot to be 0, got ", slot)
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
