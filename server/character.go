package server

import (
	"errors"
	// "github.com/grilix/sworld/sworld"
)

var (
	// ErrCharacterNotFound is when the character wasn't found
	// FIXME: We also have this on sworldservice, maybe we can reuse one of
	// them and remove the other
	ErrCharacterNotFound = errors.New("Character not found")
)

// ItemLocation represents an inventory location, used for referring to an item
type ItemLocation struct {
	BagID int `json:"bag_id"`
	Slot  int `json:"slot"`
}
