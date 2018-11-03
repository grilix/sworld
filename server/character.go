package server

import (
	"errors"

	"github.com/grilix/sworld/sworld"
)

var (
	// ErrCharacterNotFound is when the character wasn't found
	// FIXME: We also have this on sworldservice, maybe we can reuse one of
	// them and remove the other
	ErrCharacterNotFound = errors.New("Character not found")
)

// CharacterDetails represents a character in a response
type CharacterDetails struct {
	ID        string `json:"id"`
	Exploring bool   `json:"exploring"`
}

// StoneDetails holds the information about a stone item in a response
type StoneDetails struct {
	Level    int         `json:"level"`
	Zone     ZoneDetails `json:"zone"`
	Duration string      `json:"duration"`
}

// BagSlotDetails represents a bag slot on a response
type BagSlotDetails struct {
	// TODO:
	Slot  int           `json:"slot"`
	Stone *StoneDetails `json:"stone,omitempty"`
}

// BagDetails represents a bag in a response
type BagDetails struct {
	ID    int               `json:"id"`
	Items []*BagSlotDetails `json:"items"`
}

// ViewCharacterInventoryRequest represents a request for viewing the character inventory
type ViewCharacterInventoryRequest struct {
	CharacterID string `json:"character_id"`
}

// ViewCharacterInventoryResponse represents a response with the character inventory
type ViewCharacterInventoryResponse struct {
	Bags []*BagDetails `json:"bags"`
}

// ListCharactersRequest represents a request for listing the characters
type ListCharactersRequest struct {
}

// ListCharactersResponse represents a response with the characters list
type ListCharactersResponse struct {
	Characters []*CharacterDetails `json:"characters"`
}

// ViewCharacterRequest represents a request for viewing a character
type ViewCharacterRequest struct {
	ID string `json:"id"`
}

// ViewCharacterResponse represents a response with the character details
// TODO: This could also include the inventory
type ViewCharacterResponse struct {
	Character *CharacterDetails `json:"character,omitempty"`
	Error     string            `json:"error,omitempty"`
}

// TODO: Not sure what could be the best approach here, we'll just try every item type for now
func bagSlotDetails(slot int, item sworld.Item) *BagSlotDetails {
	details := &BagSlotDetails{
		Slot: slot,
	}

	stoneItem, ok := item.(*sworld.PortalStone)
	if ok {
		details.Stone = &StoneDetails{
			Level:    stoneItem.Level,
			Duration: stoneItem.Duration.String(),
		}
		return details
	}
	return details
}

func inventoryDetails(inventory []sworld.Bag) []*BagDetails {
	characterBags := make([]*BagDetails, 0, len(inventory))
	for id, bag := range inventory {
		items := bag.Items()

		bagItems := make([]*BagSlotDetails, 0, len(items))
		for slot, item := range items {
			bagItems = append(bagItems, bagSlotDetails(slot, item))
		}

		characterBags = append(characterBags, &BagDetails{
			ID:    id,
			Items: bagItems,
		})
	}
	return characterBags
}
