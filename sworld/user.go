package sworld

import (
	"errors"
)

var (
	// ErrWrongItem is when the item is not valid for an action
	ErrWrongItem = errors.New("The item is not valid for that action")
	// ErrSameSlot is when both slots are the same
	ErrSameSlot = errors.New("Can't use the same slot twice")
)

// User represents a user
type User struct {
	ID         string
	Username   string
	Characters []*Character
	Bags       []Bag
	Gold       int
}

// PickupItem stores an item into the user inventory
func (u *User) PickupItem(item Item) (ItemLocation, error) {
	location, err := u.findEmptyBagSlot(item)
	if err != nil {
		return location, err
	}
	err = u.Bags[location.BagID].StoreItem(item, location.Slot)

	return location, err
}

// TakeCharacterItem moves an item from a character to the user
func (u *User) TakeCharacterItem(characterID string, location ItemLocation) error {
	character, err := u.FindCharacter(characterID)
	if err != nil {
		return err
	}

	if location.BagID >= len(character.Bags) {
		return ErrInvalidBag
	}

	bag := character.Bags[location.BagID]

	// TODO: Lock the bag
	item, err := bag.GetItem(location.Slot)
	if err != nil {
		return err
	}

	// TODO: Lock the inventory
	newLocation, err := u.findEmptyBagSlot(item)
	if err != nil {
		return err
	}

	item, err = bag.DropItem(location.Slot)
	if err != nil {
		return err
	}
	return u.Bags[newLocation.BagID].StoreItem(item, newLocation.Slot)
}

// DropItem drops an item that is at a given location
func (u *User) DropItem(bagID, slot int) error {
	if bagID > len(u.Bags) {
		return ErrInvalidBag
	}
	bag := u.Bags[bagID]

	_, err := bag.DropItem(slot)
	return err
}

// GetItem returns the item at a given location
func (u User) GetItem(location ItemLocation) (Item, error) {
	if location.BagID > len(u.Bags) {
		return nil, ErrInvalidBag
	}
	bag := u.Bags[location.BagID]
	item, err := bag.GetItem(location.Slot)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (u User) findEmptyBagSlot(item Item) (ItemLocation, error) {
	for id, bag := range u.Bags {
		slot, err := bag.FindEmptySlot(item)
		if err == nil {
			return ItemLocation{BagID: id, Slot: slot}, nil
		}
	}
	return ItemLocation{}, ErrInventoryFull
}

// FindCharacter searchs for a character from a user
func (u User) FindCharacter(id string) (*Character, error) {
	for _, character := range u.Characters {
		if character.ID == id {
			return character, nil
		}
	}
	return nil, ErrCharacterNotFound
}

// MergeStones merges two stones
func (u *User) MergeStones(source ItemLocation, target ItemLocation) (ItemLocation, error) {
	if source.SameAs(target) {
		return ItemLocation{}, ErrSameSlot
	}
	// TODO: lock inventory
	item, err := u.GetItem(source)
	if err != nil {
		return ItemLocation{}, err
	}
	stone1, ok := item.(*PortalStone)
	if !ok {
		return ItemLocation{}, ErrWrongItem
	}

	item, err = u.GetItem(target)
	if err != nil {
		return ItemLocation{}, err
	}
	stone2, ok := item.(*PortalStone)
	if !ok {
		return ItemLocation{}, ErrWrongItem
	}

	result, err := stone1.Merge(*stone2)
	if err != nil {
		return ItemLocation{}, err
	}

	err = u.DropItem(source.BagID, source.Slot)
	if err != nil {
		return ItemLocation{}, err
	}
	err = u.DropItem(target.BagID, target.Slot)
	if err != nil {
		return ItemLocation{}, err
	}
	newLocation, err := u.PickupItem(&result)

	return newLocation, nil
}
