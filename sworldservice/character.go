package sworldservice

import (
	"github.com/grilix/sworld/sworld"
)

func (s *swService) createCharacter(user *sUser) *sworld.Character {
	character := sworld.NewCharacter()
	character.User = user.u
	user.u.Characters = append(user.u.Characters, character)

	s.characters[character.ID] = character

	return character
}

func (s *swService) respawnCharacter(user *sworld.User, character *sworld.Character) *sworld.Character {
	delete(s.characters, character.ID)

	character = sworld.NewCharacter()
	character.User = user
	user.Characters = append(user.Characters, character)

	s.characters[character.ID] = character

	return character
}
