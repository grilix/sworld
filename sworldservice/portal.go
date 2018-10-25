package sworldservice

import (
	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/encryptio/alias"
	"github.com/grilix/sworld/sworld"
)

var (
	// ErrCharacterBusy means the character is busy and cannot perform the requested task
	ErrCharacterBusy = errors.New("The character is busy")
	// ErrInvalidCharacterOwner means the user does not own that character
	ErrInvalidCharacterOwner = errors.New("Not the owner of that character")
	// ErrCantEnterPortal means the user is not allowed to enter that portal
	ErrCantEnterPortal = errors.New("The portal is not accessible")
	// ErrPortalNotFound means the portal does not exist
	ErrPortalNotFound = errors.New("The portal was not found")
	// ErrPortalIsClosed means the portal is closed
	ErrPortalIsClosed = errors.New("The portal is already closed")
)

func defaultZone(user *sworld.User) *sworld.Zone {
	return &sworld.Zone{
		ID:   sworld.RandomID(16),
		Name: "Forest", // TODO
		DropRate: sworld.DropRate{ // TODO
			Gold:     8,
			Enemy:    20,
			Item:     4,
			Stone:    9,
			Material: 5,
			Nothing:  60,
		},
	}
}

func (s *swService) defaultStone(user *sworld.User) sworld.PortalStone {
	return sworld.PortalStone{
		Level:    1,
		Zone:     defaultZone(user),
		Duration: s.defaultPortalDuration,
	}
}

func (s *swService) exploreEvent(portal *sPortal, character *sCharacter, eventID int) {
	// TODO
	switch eventID {
	case 0: // Gold
		character.gold++
		log.Printf("   -> Gold spawn at %s", portal.p.ID)
	case 1: // Enemy
		character.enemies++
		log.Printf("   -> Enemy spawn at %s", portal.p.ID)
	case 2: // Item
		character.items++
		log.Printf("   -> Item spawn at %s", portal.p.ID)
	case 3: // Stone
		character.stones++
		log.Printf("   -> Stone spawn at %s", portal.p.ID)
	case 4: // Material
		character.materials++
		log.Printf("   -> Material spawn at %s", portal.p.ID)
	}
}

func (s *swService) characterReturn(character *sCharacter, portal *sPortal) {
	log.Printf(" -> Character %s left portal %s", character.c.ID, portal.p.ID)
	log.Printf("   -> Gold: %d", character.gold)
	log.Printf("   -> Enemies: %d", character.enemies)
	log.Printf("   -> Items: %d", character.items)
	log.Printf("   -> Stones: %d", character.stones)
	log.Printf("   -> Materials: %d", character.materials)

	u := character.u.u
	u.Gold += character.gold
	u.Items += character.items
	u.Stones += character.stones
	u.Materials += character.materials

	log.Printf(" -> Current stats for user:")
	log.Printf("   -> Gold: %d", u.Gold)
	log.Printf("   -> Items: %d", u.Items)
	log.Printf("   -> Stones: %d", u.Stones)
	log.Printf("   -> Materials: %d", u.Materials)
}

func (s *swService) portalExploreHandler(portal *sPortal, character *sCharacter) {
	if character.tExplore != nil {
		panic(" portalExploreHandler doesn't expect event timer to be set.")
	}

	if portal.run == nil {
		panic("Trying to explore an uninitialized portal")
	}

	source := rand.NewSource(time.Now().UnixNano())
	seed := rand.New(source)

	rate := portal.p.PortalStone.Zone.DropRate
	events, err := alias.New([]float64{
		rate.Gold,
		rate.Enemy,
		rate.Item,
		rate.Stone,
		rate.Material,
		rate.Nothing,
	})
	if err != nil {
		panic("Can't initialize events")
	}

	defer func() {
		// Cleanup sCharacter
		character.tExplore.Stop()
		character.c.Exploring = false
		s.characterReturn(character, portal)
	}()

	character.tExplore = time.NewTicker(time.Second * 1) // TODO: character speed

	log.Printf(" -> Character %s is exploring portal %s", character.c.ID, portal.p.ID)
	for {
		select {
		case <-character.tExplore.C:
			s.exploreEvent(portal, character, int(events.Gen(seed)))
		case _, ok := <-portal.run:
			if !ok {
				return
			}
			panic("Received data from portal.run channel")
		}
	}
}

func (s *swService) portalHandler(portal *sPortal) {
	if portal.run != nil {
		panic("Trying to handle a portal that is already running")
	}

	defer func() {
		portal.p.IsOpen = false
		close(portal.run)
		portal.tClose.Stop()
		delete(s.portals, portal.p.ID)
		log.Printf(" -> Portal closed: %s\n", portal.p.ID)
	}()

	portal.run = make(chan bool)

	log.Printf(" -> Portal open: %s", portal.p.ID)
	for {
		select {
		case _, _ = <-portal.run:
			return
		case _, _ = <-portal.tClose.C:
			return
		}
	}
}
