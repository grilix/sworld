package sworld

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

//"math/rand"

// Enemy represents an enemy
type Enemy struct {
	ID        string
	Level     int
	MaxHealth int
	Health    int
	Skills    []Skill

	position int
	portal   *Portal
	D        chan bool
}

// NewEnemy creates a new enemy
func NewEnemy(portal *Portal, position int) *Enemy {
	health := (rand.Intn(10) * portal.PortalStone.Level) + (portal.PortalStone.Level * 10)

	enemy := &Enemy{
		ID:        RandomID(16),
		Level:     portal.PortalStone.Level,
		MaxHealth: health,
		Health:    health,
		Skills:    make([]Skill, 0),
		D:         make(chan bool),
		portal:    portal,
		position:  position,
	}

	return enemy
}

// AddSkill adds a skill to the enemy
func (e *Enemy) AddSkill(skill Skill) {
	e.Skills = append(e.Skills, skill)
}

// ReceiveDamage handles damage dealt to this enemy
func (e *Enemy) ReceiveDamage(source Skill, amount int) int {
	// TODO: damage reduction should be applied here
	e.Health -= amount

	if e.Health <= 0 {
		close(e.D)
	}
	fmt.Printf("Enemy received %d damage, health is now %d\n", amount, e.Health)

	return e.Health
}

// Damage returns the base damage dealt by the enemy
func (e Enemy) Damage() int {
	return 10 * e.Level
}

// AvailableSkill returns a skill that can be used right away
func (e *Enemy) AvailableSkill() Skill {
	if e.Health <= 0 {
		return nil
	}

	// TODO: Select skill
	for _, skill := range e.Skills {
		if skill.WaitTime() == 0 {
			return skill
		}
		fmt.Printf(" WAIT: %s\n", skill.WaitTime().String())
	}

	return nil
}

// ClosestExplorer returns the closest explorer from an enemy
func (e *Enemy) ClosestExplorer() (*Explorer, int) {
	var closest *Explorer
	var closestDistance int
	var distance int

	for _, explorer := range e.portal.explorers {
		character := explorer.Character
		// TODO: remove dead characters
		if character.Health <= 0 {
			continue
		}
		if explorer.position < e.position {
			distance = e.position - explorer.position
		} else {
			distance = explorer.position - e.position
		}
		if (closest == nil) || distance < closestDistance {
			closest = explorer
			closestDistance = distance
		}
	}

	return closest, closestDistance
}

func (e *Enemy) handleMove() {
	go func() {
		// TODO: define speed somehow
		moveTimer := time.NewTicker(time.Second * 1)
		defer func() {
			// Cleanup sCharacter
			moveTimer.Stop()
			// FIXME: Handle this somewhere else
			//exploration.Character.ReturnToTown(exploration.Portal)
		}()

		for {
			select {
			case <-moveTimer.C:
				explorer, distance := e.ClosestExplorer()

				if explorer != nil {
					if distance >= 1 {
						if explorer.position > e.position {
							e.position++
							log.Printf(" Enemy: Advancing, now at %d\n", e.position)
						} else {
							e.position--
							log.Printf(" Enemy: Going back, now at %d\n", e.position)
						}
					}
				}
			case _, _ = <-e.D:
				return
			case _, _ = <-e.portal.C:
				return
			}
		}
	}()
}

func (e *Enemy) handleAttack() {
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
				explorer, distance := e.ClosestExplorer()

				if explorer != nil {
					if distance < 1 {
						skill := e.AvailableSkill()

						if skill != nil {
							log.Printf(" Enemy: Attacking %v\n", explorer.Character)
							skill.Use(explorer)
						} else {
							log.Printf(" Enemy: No skills to attack!\n")
						}
					}
				}
			case _, _ = <-e.D:
				return
			case _, _ = <-e.portal.C:
				return
			}
		}
	}()
}
