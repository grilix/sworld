package sworld

import (
	"testing"
	"time"
)

func TestEncounterEvent(t *testing.T) {
	bags := []Bag{
		NewStandardBag(2),
		NewStandardBag(2),
	}
	char := Character{Bags: bags}
	itemEvent := &PortalEvent{
		Item: &PortalStone{Level: 2},
	}
	char.EncounterEvent(itemEvent)
}

func TestAvailableSkill(t *testing.T) {
	char := &Character{
		Health: 100,
		Skills: make([]Skill, 2),
	}

	s1 := NewHitSkill(char)
	s2 := NewHitSkill(char)

	char.Skills[0] = s1
	char.Skills[1] = s2

	available := char.AvailableSkill()

	if available == nil {
		t.Fatal("Expected AvailableSkill to return a skill")
	}

	s1.cooldown = time.Minute
	s2.cooldown = time.Minute
	s1.lastUse = time.Now()
	s2.lastUse = time.Now()

	available = char.AvailableSkill()

	if available != nil {
		t.Fatal("Expected AvailableSkill to not return a skill, got", available)
	}
}

func TestPickupItem(t *testing.T) {
	bags := []Bag{
		NewStandardBag(2),
		NewStandardBag(2),
	}
	char := Character{Bags: bags}
	bag, slot, err := char.pickupItem(&PortalStone{Level: 3})
	if err != nil {
		t.Fatal("Can't pick up item, ", err)
	}
	if bag != 0 {
		t.Error("Expected item to be on bag 0, got ", bag)
	}
	if slot != 0 {
		t.Error("Expected item to be on slot 0, got ", slot)
	}

	bagOne, ok := char.Bags[0].(*StandardBag)
	if !ok {
		t.Fatal("Can't cast bag")
	}
	bagTwo, ok := char.Bags[1].(*StandardBag)
	if !ok {
		t.Fatal("Can't cast bag")
	}

	if bagOne.items[0] == nil {
		t.Fatal("Item not picked up")
	}
	if bagTwo.items[0] != nil {
		t.Fatal("Item picked up twice")
	}

	// Picks up 3 more items
	bag, slot, err = char.pickupItem(&PortalStone{Level: 3})
	if err != nil {
		t.Fatal("Can't pick up item, ", err)
	}
	if bag != 0 {
		t.Error("Expected item to be on bag 0, got ", bag)
	}
	if slot != 1 {
		t.Error("Expected item to be on slot 1, got ", slot)
	}

	bag, slot, err = char.pickupItem(&PortalStone{Level: 3})
	if err != nil {
		t.Fatal("Can't pick up item, ", err)
	}
	if bag != 1 {
		t.Error("Expected item to be on bag 1, got ", bag)
	}
	if slot != 0 {
		t.Error("Expected item to be on slot 0, got ", slot)
	}

	bag, slot, err = char.pickupItem(&PortalStone{Level: 3})
	if err != nil {
		t.Fatal("Can't pick up item, ", err)
	}
	if bag != 1 {
		t.Error("Expected item to be on bag 1, got ", bag)
	}
	if slot != 1 {
		t.Error("Expected item to be on slot 1, got ", slot)
	}

}
