package sworld

import (
	"errors"
	"fmt"
	"log"
)

var (
	// ErrCharacterNotFound is when the character cannot be found
	ErrCharacterNotFound = errors.New("That character was not found")
	// ErrPortalIsClosed means the portal is closed
	ErrPortalIsClosed = errors.New("The portal is already closed")
	// ErrCharacterBusy means the character is busy and cannot perform the requested task
	ErrCharacterBusy = errors.New("The character is busy")
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
	ID        string
	Level     int
	Health    int
	MaxHealth int
	Gold      int
	Exploring bool
	User      *User
	Bags      []Bag
	Skills    []Skill
	D         chan bool

	// TODO: this is so we can debug things
	enemies int
}

// NewCharacter creates a character
func NewCharacter() *Character {
	health := 100

	character := &Character{
		ID:        RandomID(16),
		Level:     1,
		Health:    health,
		MaxHealth: health,
		Gold:      0,
		Exploring: false,
		Skills:    make([]Skill, 1),
		D:         make(chan bool),
		Bags: []Bag{
			NewStandardBag(10),
		},
	}

	// FIXME: I don't really like this cross-dependency
	character.Skills[0] = NewHitSkill(character)

	return character
}

// Die kills the character
func (c *Character) Die() {
	c.Health = 0
	for _, bag := range c.Bags {
		bag.Empty()
	}

	if c.User == nil {
		log.Println(" --->  Character doe not have a user!")
	}

	log.Printf("Character: Died.\n")
	close(c.D)
}

// Damage returns the base damage dealt by the character
func (c Character) Damage() int {
	return c.Level * 20
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

// AvailableSkill will return a skill that can be used right away
func (c *Character) AvailableSkill() Skill {
	if c.Health <= 0 {
		return nil
	}

	// TODO: Select skill
	for _, skill := range c.Skills {
		if skill.WaitTime() == 0 {
			return skill
		}
		fmt.Printf(" WAIT: %s\n", skill.WaitTime().String())
	}

	return nil
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
// TODO: remove?
func (c *Character) EncounterEvent(event *PortalEvent) error {
	if event.Item != nil {
		c.pickupItem(event.Item)
	}

	if event.Enemy != nil {
		//c.fightEnemy(event.Enemy)
		log.Printf("   -> Enemy spawn %T%v\n", event.Enemy, event.Enemy)
	}

	if event.Gold > 0 {
		c.Gold += event.Gold
		log.Printf("   -> Gold spawn\n")
	}

	return nil
}

// EnterPortal makes a character enter a portal
func (c *Character) EnterPortal(portal *Portal) (*Explorer, error) {
	if !portal.IsOpen {
		return nil, ErrPortalIsClosed
	}

	if c.Exploring {
		return nil, ErrCharacterBusy
	}

	c.Exploring = true

	exploration := &Explorer{
		Portal:    portal,
		Character: c,
	}
	portal.explorers = append(portal.explorers, exploration)

	return exploration, nil
}
