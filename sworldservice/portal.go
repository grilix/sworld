package sworldservice

import (
	"errors"
	"log"
	"time"

	"github.com/grilix/sworld/sworld"
)

var (
	// ErrCantEnterPortal means the user is not allowed to enter that portal
	ErrCantEnterPortal = errors.New("The portal is not accessible")
	// ErrPortalNotFound means the portal does not exist
	ErrPortalNotFound = errors.New("The portal was not found")
)

func (s *swService) defaultStone(user *sworld.User) sworld.PortalStone {
	return sworld.PortalStone{
		Level:    1,
		Zone:     s.defaultZone,
		Duration: s.defaultPortalDuration,
	}
}

func (s *swService) handleCharacterMove(exploration *sworld.Explorer) {
	// TODO: Not sure what was this for
	if exploration.Portal.C == nil {
		panic("explorePortal received a portal that is not initialized")
	}

	go func() {
		// TODO: define speed somehow
		moveTimer := time.NewTicker(time.Second * 1)
		defer func() {
			// Cleanup sCharacter
			moveTimer.Stop()

			// FIXME: Handle this somewhere else
			if exploration.Character.Health > 0 {
				exploration.Character.ReturnToTown(exploration.Portal)
			} else {
				//
			}
		}()

		for {
			select {
			case <-moveTimer.C:
				enemy := exploration.ClosestEnemy()

				if enemy == nil {
					_ = exploration.Advance()
					log.Printf(" Character: Advancing, now at %d\n", exploration.Position())
				}
			case _, _ = <-exploration.Character.D:
				return
			case _, _ = <-exploration.Portal.C:
				return
			}
		}
	}()
}

func (s *swService) handleCharacterAttack(exploration *sworld.Explorer) {
	// TODO: Not sure what was this for
	if exploration.Portal.C == nil {
		panic("explorePortal received a portal that is not initialized")
	}

	go func() {
		// TODO: define speed somehow
		attackTimer := time.NewTicker(time.Millisecond * 500)
		defer func() {
			// Cleanup sCharacter
			attackTimer.Stop()
		}()

		for {
			select {
			case <-attackTimer.C:
				character := exploration.Character
				enemy := exploration.ClosestEnemy()

				if enemy != nil {
					skill := character.AvailableSkill()
					if skill != nil {
						log.Printf(" Character: Attacking %v\n", enemy)
						skill.Use(enemy)
					} else {
						log.Printf(" Character: No skills to attack!\n")
					}
				}
			case _, _ = <-exploration.Character.D:
				return
			case _, _ = <-exploration.Portal.C:
				return
			}
		}
	}()
}
