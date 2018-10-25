package sworldservice

import (
	"context"
	"errors"
	"time"

	"github.com/grilix/sworld/sworld"
)

var (
	// ErrNotImplemented means a service call is not implemented
	ErrNotImplemented = errors.New("This service call is not implemented")
	// ErrPortalStoneNotFound means the stone does not exist
	ErrPortalStoneNotFound = errors.New("The stone does not exist")
)

type sCharacter struct {
	c        *sworld.Character
	u        *sUser
	tExplore *time.Ticker

	// TODO: inventory
	gold,
	enemies,
	items,
	stones,
	materials int
}

type sPortal struct {
	p      sworld.Portal
	c      *sCharacter
	run    chan bool
	tClose *time.Timer
}

type sUser struct {
	u *sworld.User
}

// Service is the sworld's super service
type Service interface {
	Authenticate(ctx context.Context, c Credentials) (*sworld.User, error)
	FindUser(id string) *sworld.User
	OpenPortal(user *sworld.User, stoneID string) (*sworld.Portal, error)
	ExplorePortal(user *sworld.User, portalID, characterID string) error
}

type swService struct {
	users                 map[string]*sUser
	portals               map[string]*sPortal
	world                 *sworld.World
	defaultPortalDuration time.Duration
}

// NewService creates the service
func NewService(world *sworld.World) Service {
	return &swService{
		world:                 world,
		users:                 make(map[string]*sUser),
		portals:               make(map[string]*sPortal),
		defaultPortalDuration: time.Second * 5, // TODO: set default duration
	}
}

func (s *swService) FindUser(id string) *sworld.User {
	user := s.users[id]
	if user != nil {
		return user.u
	}

	return nil
}

func (s *swService) Authenticate(ctx context.Context, c Credentials) (*sworld.User, error) {
	user, err := s.userByUsername(c.Username)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *swService) ExplorePortal(user *sworld.User, portalID, characterID string) error {
	sportal := s.portals[portalID]
	if sportal == nil {
		return ErrPortalNotFound
	}
	if !sportal.p.IsOpen {
		return ErrPortalIsClosed
	}
	if sportal.p.User.ID != user.ID {
		return ErrCantEnterPortal
	}
	if user.Character.ID != characterID {
		return ErrInvalidCharacterOwner
	}

	if user.Character.Exploring {
		return ErrCharacterBusy
	}

	user.Character.Exploring = true
	sportal.p.Character = user.Character
	sportal.c = &sCharacter{
		c: user.Character,
		u: s.users[user.ID],
	}

	go s.portalExploreHandler(sportal, sportal.c)

	return nil
}

func (s *swService) OpenPortal(user *sworld.User, stoneID string) (*sworld.Portal, error) {
	if stoneID != "" {
		return nil, ErrPortalStoneNotFound
	}

	// TODO: fetch user stones
	stone := s.defaultStone(user)

	sportal := &sPortal{
		p: sworld.Portal{
			ID:          sworld.RandomID(16),
			PortalStone: stone,
			IsOpen:      true,
			User:        user,
		},
		tClose: time.NewTimer(stone.Duration),
	}

	s.portals[sportal.p.ID] = sportal

	go s.portalHandler(sportal)

	return &sportal.p, nil
}
