package sworld

import "log"

// Explorer is the character session on a portal
type Explorer struct {
	Character *Character
	Portal    *Portal
	position  int
}

// ClosestEnemy returns the closest enemy from an explorer
func (e *Explorer) ClosestEnemy() *Enemy {
	var closest *Enemy
	var closestDistance int
	var distance int

	for _, enemy := range e.Portal.enemies {
		// TODO: remove dead enemies
		if enemy.Health <= 0 {
			continue
		}
		if enemy.position < e.position {
			distance = e.position - enemy.position
		} else {
			distance = enemy.position - e.position
		}
		if (closest == nil) || distance < closestDistance {
			closest = enemy
			closestDistance = distance
		}
	}

	return closest
}

// ReceiveDamage handles the damage received by an explorer
func (e *Explorer) ReceiveDamage(source Skill, amount int) int {
	e.Character.Health -= amount
	log.Printf("Character: Received %d damage, health is now: %d\n", amount, e.Character.Health)

	if e.Character.Health <= 0 {
		e.Character.Die()
	}

	return e.Character.Health
}

// Position returns the current position of an explorer
func (e Explorer) Position() int {
	return e.position
}

// Advance moves the explorer forward
func (e *Explorer) Advance() *PortalEvent {
	e.position++
	p := e.Portal

	//cleared := len(p.portalMap)
	cleared := p.cleared

	var event *PortalEvent

	if cleared < e.position {
		event = p.RandomEvent(e.position)
	}

	if event != nil {
		e.Character.EncounterEvent(event)
	}

	return event
}
