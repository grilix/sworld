package sworld

import (
	"testing"
)

func TestClosestEnemy(t *testing.T) {
	portal := &Portal{
		PortalStone: PortalStone{Level: 1},
		enemies:     make([]*Enemy, 0, 2),
	}
	portal.randomEnemyEvent(2)
	portal.randomEnemyEvent(3)

	exploration := &Explorer{
		Portal:   portal,
		position: 1,
	}
	enemy := exploration.ClosestEnemy()
	if enemy == nil {
		t.Fatal("Expected ClosestEnemy to return an enemy")
	}

	if enemy.position != 2 {
		t.Fatal("Expected closest enemy position to be 2, got", enemy.position)
	}

	exploration.position = 3
	enemy = exploration.ClosestEnemy()
	if enemy == nil {
		t.Fatal("Expected ClosestEnemy to return an enemy")
	}

	if enemy.position != 3 {
		t.Fatal("Expected closest enemy position to be 3, got", enemy.position)
	}
}
