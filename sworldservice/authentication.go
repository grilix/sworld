package sworldservice

import (
	"github.com/grilix/sworld/sworld"
)

// Credentials represent the user credentials for signing in
type Credentials struct {
	Username string
	Password string
}

func (s *swService) createUser(username string) (*sworld.User, error) {
	user := &sUser{
		u: &sworld.User{
			ID:       sworld.RandomID(16),
			Username: username,
		},
	}
	s.createCharacter(user)
	s.users[user.u.ID] = user
	return user.u, nil
}

func (s *swService) userByUsername(username string) (*sworld.User, error) {
	for _, user := range s.users {
		if user.u.Username == username {
			return user.u, nil
		}
	}

	return s.createUser(username)
}
