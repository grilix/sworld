package sworld

import (
	"errors"
	"log"
	"time"
)

var (
	// ErrCharacterNotFound is when the character cannot be found
	ErrCharacterNotFound = errors.New("That character was not found")
)

// Delete item:
//     a[i] = a[len(a)-1]
//     a[len(a)-1] = nil
//     a = a[:len(a)-1]
//
// Killing enemy xp:
// xp := int64(
//     math.Round(
//         (float64(enemy.Level) / float64(c.Level)) * float64(enemy.Level),
//     ),
// )

// Next lvl xp:
// return int64(math.Round((4 * math.Pow(float64(c.Level), 3)) / 5))

// Character represents a user character
type Character struct {
	ID              string
	Health          int
	MaxHealth       int
	Gold            int
	Exploring       bool
	Level           int
	ExploreInterval time.Duration
	User            *User
	Bags            []Bag

	// TODO: this is so we can debug things
	enemies int
}

func (c Character) findEmptyBagSlot(item Item) (int, int, error) {
	for id, bag := range c.Bags {
		slot, err := bag.FindEmptySlot(item)
		if err == nil {
			return id, slot, nil
		}
	}
	return 0, 0, ErrInventoryFull
}

// TODO: This could be exported
func (c *Character) pickupItem(item Item) (int, int, error) {
	bagID, slot, err := c.findEmptyBagSlot(item)
	if err != nil {
		return bagID, slot, err
	}
	c.Bags[bagID].StoreItem(item, slot)

	return bagID, slot, nil
}

// DropItem discards an item that is at a given location
func (c *Character) DropItem(bagID, slot int) (Item, error) {
	if bagID >= len(c.Bags) {
		return nil, ErrInvalidBag
	}
	bag := c.Bags[bagID]
	return bag.DropItem(slot)
}

// ReturnToTown makes the character leave the "exploring" state
func (c *Character) ReturnToTown(portal *Portal) {
	c.Exploring = false

	// FIXME: This is just for debugging
	for bagID, bag := range c.Bags {
		items := bag.Items()

		for slot := range items {
			if items[slot] != nil {
				log.Printf(" -> Bag %d slot %d: %T%v", bagID, slot, items[slot], items[slot])
			}
		}
	}

	log.Printf(" -> Character %s left portal %s", c.ID, portal.ID)
	log.Printf("   -> Gold: %d", c.Gold)
	log.Printf("   -> Enemies: %d", c.enemies)

	u := c.User
	u.Gold += c.Gold

	log.Printf(" -> Current stats for user:")
	log.Printf("   -> Gold: %d", u.Gold)
}

// EncounterEvent lets a character handle an event
func (c *Character) EncounterEvent(event PortalEvent) error {
	if event.Item != nil {
		c.pickupItem(event.Item)
	}
	if event.Gold > 0 {
		c.Gold += event.Gold
	}

	return nil
}
