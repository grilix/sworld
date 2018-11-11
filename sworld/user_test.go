package sworld

import (
	"testing"
	"time"
)

func TestUserMergeStones(t *testing.T) {
	userZone1 := &Zone{ID: "Zone 1", Name: "Test zone 1"}

	stone1 := PortalStone{
		Level:    1,
		Duration: 14 * time.Second,
		Zone:     userZone1,
	}
	stone2 := PortalStone{
		Level:    1,
		Duration: 12 * time.Second,
		Zone:     userZone1,
	}

	user := User{
		Bags: []Bag{NewStandardBag(10)},
	}

	source, err := user.PickupItem(&stone1)
	target, err := user.PickupItem(&stone2)

	_, err = user.MergeStones(source, source)
	if err == nil {
		t.Error("Expected merge to fail for same slot")
	}

	result, err := user.MergeStones(source, target)
	if err != nil {
		t.Fatal(err)
	}
	item, err := user.GetItem(result)
	stone, ok := item.(*PortalStone)
	if !ok {
		t.Error("Expected item to be a stone, got:", item)
	}

	if stone.Level != 2 {
		t.Error("Expected level to be 2, got:", stone.Level)
	}

	item, err = user.GetItem(target)
	if err == nil {
		t.Error("Expected item at 0,0 to have been dropped, got:", item)
	}
}
