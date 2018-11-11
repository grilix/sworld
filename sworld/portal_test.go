package sworld

import (
	"math/rand"
	"testing"
	"time"
)

func TestRandomEnemyEvent(t *testing.T) {
	source := rand.NewSource(time.Now().UnixNano())
	seed := rand.New(source)

	portal := Portal{
		seed:        seed,
		PortalStone: PortalStone{Level: 0},
	}
	event := portal.randomEnemyEvent(1)
	if event != nil {
		t.Error("Expected level 0 portal to not spawn enemies")
	}

	portal = Portal{
		PortalStone: PortalStone{Level: 1},
		seed:        seed,
	}
	event = portal.randomEnemyEvent(1)
	if event == nil {
		t.Fatal("Expected level 1 portal to spawn enemies")
	}

	if event.Enemy == nil {
		t.Fatal("Expected event to be Enemy")
	}
	if event.Enemy.position != 1 {
		t.Error("Expected enemy to be at position 1, got", event.Enemy.position)
	}
}

func TestRandomItemEvent(t *testing.T) {
	source := rand.NewSource(time.Now().UnixNano())
	seed := rand.New(source)

	portal := Portal{
		seed: seed,
	}
	event := portal.randomItemEvent()
	if event.Item == nil {
		t.Fatal("Expected event to be Item")
	}
	_, ok := event.Item.(*PortalStone)
	if !ok {
		t.Fatal("Expected item to be a stone")
	}
}
