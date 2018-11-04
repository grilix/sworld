package sworld

import (
	"errors"
)

var (
	// ErrInventoryFull is when there's no more space on a bag
	// TODO: Maybe find a better name, not sure
	ErrInventoryFull = errors.New("There is no space left for that item")
	// ErrInvalidBagSlot is when the slot is < 0 or > capacity
	ErrInvalidBagSlot = errors.New("That's an invalid slot number")
	// ErrEmptyBagSlot is when the slot is empty
	ErrEmptyBagSlot = errors.New("That slot is empty")
	// ErrNotEmptyBagSlot means the bag slot already has an item
	ErrNotEmptyBagSlot = errors.New("That slot is not empty")
	// ErrInvalidBag is when the bag does not exist on character or user
	ErrInvalidBag = errors.New("That bag does not exist")
)

// Item is the interface for an item
type Item interface {
	//
}

// Bag is the interface for a bag
// A bag is used to hold items either by a character or a user
type Bag interface {
	StoreItem(item Item, slot int) error
	GetItem(slot int) (Item, error)
	FindEmptySlot(item Item) (int, error)
	DropItem(slot int) (Item, error)
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

// GetItem returns the item that is at a given slot
func (b StandardBag) GetItem(slot int) (Item, error) {
	if slot >= len(b.items) {
		return nil, ErrInvalidBagSlot
	}
	item := b.items[slot]
	if item == nil {
		return nil, ErrEmptyBagSlot
	}
	return item, nil
}

// DropItem removes an item from the bag
func (b *StandardBag) DropItem(slot int) (Item, error) {
	if slot < 0 || slot >= len(b.items) {
		return nil, ErrInvalidBagSlot
	}
	if b.items[slot] == nil {
		return nil, ErrEmptyBagSlot
	}
	item := b.items[slot]
	b.items[slot] = nil
	return item, nil
}

// Items returns the items from a bag
func (b StandardBag) Items() []Item {
	return b.items
}

// FindEmptySlot returns a slot available for a given item
func (b StandardBag) FindEmptySlot(item Item) (int, error) {
	for slot := range b.items {
		if b.items[slot] == nil {
			return slot, nil
		}
	}
	return 0, ErrInventoryFull
}

// StoreItem puts an item on an empty slot, if available
func (b *StandardBag) StoreItem(item Item, slot int) error {
	if b.items[slot] != nil {
		return ErrNotEmptyBagSlot
	}

	b.items[slot] = item

	return nil
}
