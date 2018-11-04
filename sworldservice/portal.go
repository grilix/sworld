package sworldservice

import (
	"errors"
	"math/rand"
	"time"

	"github.com/grilix/sworld/sworld"
)

var (
	// ErrCharacterBusy means the character is busy and cannot perform the requested task
	ErrCharacterBusy = errors.New("The character is busy")
	// ErrCantEnterPortal means the user is not allowed to enter that portal
	ErrCantEnterPortal = errors.New("The portal is not accessible")
	// ErrPortalNotFound means the portal does not exist
	ErrPortalNotFound = errors.New("The portal was not found")
	// ErrPortalIsClosed means the portal is closed
	ErrPortalIsClosed = errors.New("The portal is already closed")
)

// TODO: this stuff should be handled by some sort of settings
func defaultZone(user *sworld.User) *sworld.Zone {
	return &sworld.Zone{
		ID:   sworld.RandomID(16),
		Name: "Forest",
		DropRate: sworld.DropRate{
			Gold:    8,
			Enemy:   20,
			Item:    20,
			Nothing: 60,
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

func (s *swService) handleExplore(portal *sworld.Portal, character *sworld.Character) {
	if character.Exploring {
		panic("explorePortal received a character that is already exploring")
	}

	if portal.C == nil {
		panic("explorePortal received a portal that is not initialized")
	}

	character.Exploring = true

	source := rand.NewSource(time.Now().UnixNano())
	seed := rand.New(source)

	go func() {
		// TODO: Not sure what to do here, we could probably use a character speed or something
		// I didn't want the portal to have a predefined interval for events, but idk
		exploreTimer := time.NewTicker(time.Second * 1)
		defer func() {
			// Cleanup sCharacter
			exploreTimer.Stop()
			character.ReturnToTown(portal)
		}()

		for {
			select {
			case <-exploreTimer.C:
				event := portal.RandomEvent(seed)
				err := character.EncounterEvent(event)
				if err != nil {
					// TODO:
				}

			case _, _ = <-portal.C:
				// Portal closed
				return
			}
		}
	}()
}
