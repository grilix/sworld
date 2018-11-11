package sworld

import (
	"testing"
)

func TestSkill(t *testing.T) {
	char := &Character{Level: 1}
	enemy := &Enemy{Health: 1000}
	skill := NewHitSkill(char)
	skill.Use(enemy)

	if !(enemy.Health < 1000) {
		t.Error("Expected enemy to have received damage, but Health is", enemy.Health)
	}
}
