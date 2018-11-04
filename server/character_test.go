package server

import (
	"github.com/grilix/sworld/sworld"
	"testing"
)

func TestInventoryDetails(t *testing.T) {
	bags := []sworld.Bag{
		sworld.NewStandardBag(5),
		sworld.NewStandardBag(5),
	}
	zone := &sworld.Zone{}
	bags[0].StoreItem(&sworld.PortalStone{Level: 3, Zone: zone}, 0)

	details := inventoryDetails(bags)
	if len(details) != 2 {
		t.Error("Expected two bags, got ", len(details))
	}

	bagOne := details[0]
	bagTwo := details[1]

	if len(bagOne.Items) != 5 {
		t.Error("Expected first bag to have 5 items, got ", len(bagOne.Items))
	}
	if len(bagTwo.Items) != 5 {
		t.Error("Expected second bag to have 5 items, got ", len(bagTwo.Items))
	}

	if bagOne.Items[0].Stone == nil {
		t.Error("Expected slot 0 on first bag to be a stone")
	}

	if bagTwo.Items[0].Stone != nil {
		t.Error("Expected slot 0 on second bag to not be a stone, got ", bagTwo.Items[0].Stone)
	}
}
