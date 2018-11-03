package sworld

import (
	"errors"
)

var (
	// ErrInventoryFull is when there's no more space on a bag
	// TODO: Maybe find a better name, not sure
	ErrInventoryFull = errors.New("There is no space left for that item")
)

// Item is the interface for an item
type Item interface {
	//
}

// Bag is the interface for a bag
// A bag is used to hold items either by a character or a user
type Bag interface {
	StoreItem(item Item) (int, error)
	Items() []Item
}

// StandardBag is the main bag type
type StandardBag struct {
	items []Item
}

// NewStandardBag creates a standard bag of a given capacity
func NewStandardBag(capacity int) *StandardBag {
	return &StandardBag{
		items: make([]Item, capacity),
	}
}

// Items returns the items from a bag
func (b StandardBag) Items() []Item {
	return b.items
}

func (b StandardBag) findEmptySlot(item Item) (int, error) {
	for slot := range b.items {
		if b.items[slot] == nil {
			return slot, nil
		}
	}
	return 0, ErrInventoryFull
}

// StoreItem puts an item on an empty slot, if available
func (b *StandardBag) StoreItem(item Item) (int, error) {
	slot, err := b.findEmptySlot(item)
	if err != nil {
		return 0, err
	}
	b.items[slot] = item

	return slot, nil
}
