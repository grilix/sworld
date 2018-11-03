package sworldservice

import (
	"github.com/grilix/sworld/sworld"
)

func (s *swService) createCharacter(user *sUser) *sworld.Character {
	health := 10000

	character := &sworld.Character{
		ID:        sworld.RandomID(16),
		Level:     1,
		Health:    health,
		MaxHealth: health,
		Gold:      0,
		Exploring: false,
		User:      user.u,
		Bags: []sworld.Bag{
			sworld.NewStandardBag(10),
		},
	}
	user.u.Character = character

	s.characters[character.ID] = character

	return character
}
