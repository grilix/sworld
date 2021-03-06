package sworld

import (
	"math/rand"
	"time"

	"github.com/encryptio/alias"
)

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

	startedAt  time.Time
	enemies    []*Enemy
	explorers  []*Explorer
	eventsRate *alias.Alias // TODO: rename eventDrops
	drops      *alias.Alias
	cleared    int
	seed       *rand.Rand
}

// PortalEvent is generated by the portal and sent to a character
// Eeach character is given its own event
type PortalEvent struct {
	Item  Item
	Enemy *Enemy
	Gold  int
}

// TimeLeft is the amount of time until the portal closes
func (p *Portal) TimeLeft() time.Duration {
	return p.PortalStone.Duration - time.Since(p.startedAt)
}

// RandomEnemyEvent returns a random enemy
func (p *Portal) RandomEnemyEvent(position int) *PortalEvent {
	if p.PortalStone.Level < 1 {
		return nil
	}

	enemy := NewEnemy(p, position)

	// FIXME: I don't really like this cross-dependency
	enemy.AddSkill(NewHitSkill(enemy))
	p.enemies = append(p.enemies, enemy)

	enemy.handleAttack()
	enemy.handleMove()

	return &PortalEvent{Enemy: enemy}
}

// DeadEnemies returns the dead enemies on this portal
func (p *Portal) DeadEnemies() []*Enemy {
	enemies := make([]*Enemy, 0, len(p.enemies))
	for _, enemy := range p.enemies {
		if enemy.Health <= 0 {
			enemies = append(enemies, enemy)
		}
	}
	return enemies
}

// OpenPortal opens a portal and sets a timer for closing it
func OpenPortal(user *User, stone PortalStone, closeFn func(*Portal)) (*Portal, error) {
	source := rand.NewSource(time.Now().UnixNano())

	p := &Portal{
		ID:          RandomID(16),
		PortalStone: stone,
		IsOpen:      true,
		User:        user,
		seed:        rand.New(source),
		enemies:     make([]*Enemy, 0, 10),
		explorers:   make([]*Explorer, 0, 1),
		startedAt:   time.Now(),
	}

	stone.Zone.InitializePortal(p)

	p.C = make(chan bool)
	go func() {
		defer func() {
			p.IsOpen = false
			close(p.C)
		}()

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
