package sworld

import (
	"log"
	"math/rand"
	"time"

	"github.com/encryptio/alias"
)

// DropRate represents the rates for drops
type DropRate struct {
	Gold    float64
	Enemy   float64
	Item    float64
	Nothing float64
}

// Zone defines the type of enemies that will be found
// It's the base for creating a portal
type Zone struct {
	ID       string
	Name     string
	DropRate DropRate
}

// PortalStone is used to open a portal
type PortalStone struct {
	Level        int
	Zone         *Zone
	Duration     time.Duration
	DropInterval time.Duration
}

// Portal is an instance of a Zone, where players can teleport to
type Portal struct {
	ID          string
	PortalStone PortalStone
	User        *User
	Character   *Character

	// TODO: Would this be the same as checking for C != nil ?
	IsOpen bool

	// C is the channel that communicates the portal closing event
	C chan bool

	eventsRate *alias.Alias
}

// PortalEvent is generated by the portal and sent to a character
// Eeach character is given its own event
type PortalEvent struct {
	Item Item
	Gold int
}

// RandomItemEvent creates an event with a random item
func (p Portal) RandomItemEvent(seed *rand.Rand) PortalEvent {
	// TODO: Random item
	stone := &PortalStone{}
	log.Printf("   -> Item spawn %T%v", stone, stone)
	return PortalEvent{Item: stone}
}

// RandomEvent creates a random event
func (p Portal) RandomEvent(seed *rand.Rand) PortalEvent {
	eventID := int(p.eventsRate.Gen(seed))

	// TODO
	switch eventID {
	case 0: // Gold
		return PortalEvent{Gold: 1}
		//c.gold++
		log.Printf("   -> Gold spawn at %s", p.ID)
	case 1: // Enemy
		//c.enemies++
		log.Printf("   -> Enemy spawn at %s", p.ID)
	case 2: // Item
		return p.RandomItemEvent(seed)
	}

	return PortalEvent{}
}

// GetRates returns an alias.Alias instance for drop rates
func (r DropRate) GetRates() (*alias.Alias, error) {
	return alias.New([]float64{
		r.Gold,
		r.Enemy,
		r.Item,
		r.Nothing,
	})
}

// OpenPortal opens a portal and sets a timer for closing it
func OpenPortal(user *User, stone PortalStone, closeFn func(*Portal)) (*Portal, error) {
	p := &Portal{
		ID:          RandomID(16),
		PortalStone: stone,
		IsOpen:      true,
		User:        user,
	}

	rate := stone.Zone.DropRate
	events, err := rate.GetRates()
	if err != nil {
		// FIXME: return err instead?
		panic("Can't initialize events")
	}
	p.eventsRate = events

	p.C = make(chan bool)
	go func() {
		defer close(p.C)

		tClose := time.NewTimer(stone.Duration)

		defer tClose.Stop()
		defer closeFn(p)
		for {
			select {
			case _, _ = <-p.C:
				return
			case _, _ = <-tClose.C:
				return
			}
		}
	}()

	return p, nil
}
