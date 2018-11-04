package sworld

// User represents a user
type User struct {
	ID         string
	Username   string
	Characters []*Character
	Bags       []Bag
	Gold       int
}

// TakeCharacterItem moves an item from a character to the user
func (u *User) TakeCharacterItem(characterID string, bagID, slot int) error {
	character, err := u.FindCharacter(characterID)
	if err != nil {
		return err
	}

	if bagID >= len(character.Bags) {
		return ErrInvalidBag
	}

	bag := character.Bags[bagID]

	// TODO: Lock the bag
	item, err := bag.GetItem(slot)
	if err != nil {
		return err
	}

	// TODO: Lock the inventory
	newBagID, newSlot, err := u.findEmptyBagSlot(item)
	if err != nil {
		return err
	}

	item, err = bag.DropItem(slot)
	if err != nil {
		return err
	}
	u.Bags[newBagID].StoreItem(item, newSlot)

	return nil
}

// DropItem drops an item that is at a given location
func (u User) DropItem(bagID, slot int) error {
	if bagID > len(u.Bags) {
		return ErrInvalidBag
	}
	bag := u.Bags[bagID]

	_, err := bag.DropItem(slot)
	return err
}

// GetItem returns the item at a given location
func (u User) GetItem(bagID, slot int) (Item, error) {
	if bagID > len(u.Bags) {
		return nil, ErrInvalidBag
	}
	bag := u.Bags[bagID]
	item, err := bag.GetItem(slot)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (u User) findEmptyBagSlot(item Item) (int, int, error) {
	for id, bag := range u.Bags {
		slot, err := bag.FindEmptySlot(item)
		if err == nil {
			return id, slot, nil
		}
	}
	return 0, 0, ErrInventoryFull
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
