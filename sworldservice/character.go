package sworldservice

import (
	"github.com/grilix/sworld/sworld"
)

func (s *swService) createCharacter(user *sUser) *sworld.Character {
	health := 10000

	character := &sworld.Character{
		Level:     1,
		Health:    health,
		MaxHealth: health,
		Gold:      0,
		Exploring: false,
		ID:        sworld.RandomID(16),
	}
	user.u.Character = character
	return character
}
