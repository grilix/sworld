package sworldservice

import (
	"context"
	"errors"
	"log"
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
	c *sworld.Character
	u *sUser
}

type sPortal struct {
	p      *sworld.Portal
	c      *sCharacter
	tClose *time.Timer
}

type sUser struct {
	u *sworld.User
}

// Service is the sworld's super service
type Service interface {
	Authenticate(ctx context.Context, c Credentials) (*sworld.User, error)
	FindUser(id string) *sworld.User
	ViewUserInventory(user *sworld.User) ([]sworld.Bag, error)

	ListCharacters(user *sworld.User) ([]*sworld.Character, error)
	ViewCharacterInventory(characterID string) ([]sworld.Bag, error)
	DropCharacterItem(user *sworld.User, characterID string, bagID, slot int) error
	TakeCharacterItem(user *sworld.User, characterID string, bagID, slot int) error

	OpenPortal(user *sworld.User, stoneID string) (*sworld.Portal, error)
	ExplorePortal(user *sworld.User, portalID, characterID string) error
	ViewPortal(portalID string) (*sworld.Portal, error)
	ListPortals(user *sworld.User) ([]*sworld.Portal, error)
}

type swService struct {
	users                 map[string]*sUser
	portals               map[string]*sPortal
	defaultPortalDuration time.Duration
	characters            map[string]*sworld.Character
}

// NewService creates the service
func NewService() Service {
	return &swService{
		users:      make(map[string]*sUser),
		portals:    make(map[string]*sPortal),
		characters: make(map[string]*sworld.Character),
		// TODO: this should be on settings
		defaultPortalDuration: time.Second * 10,
	}
}

func (s *swService) TakeCharacterItem(user *sworld.User, characterID string, bagID, slot int) error {
	// FIXME: Maybe just remove this service method and use user directly?
	return user.TakeCharacterItem(characterID, bagID, slot)
}

func (s *swService) DropCharacterItem(user *sworld.User, characterID string, bagID, slot int) error {
	character, err := user.FindCharacter(characterID)
	if err != nil {
		return err
	}
	if bagID >= len(character.Bags) {
		return sworld.ErrInvalidBag
	}
	bag := character.Bags[bagID]
	_, err = bag.DropItem(slot)

	return err
}

func (s *swService) ViewUserInventory(user *sworld.User) ([]sworld.Bag, error) {
	// FIXME: call user.Bags directly?
	return user.Bags, nil
}

func (s *swService) ViewCharacterInventory(characterID string) ([]sworld.Bag, error) {
	character := s.characters[characterID]
	if character == nil {
		return nil, sworld.ErrCharacterNotFound
	}

	return character.Bags, nil
}

func (s *swService) ListCharacters(user *sworld.User) ([]*sworld.Character, error) {
	// FIXME: Not sure about this
	return user.Characters, nil
}

func (s *swService) ViewPortal(portalID string) (*sworld.Portal, error) {
	portal := s.portals[portalID]
	if portal == nil {
		return nil, ErrPortalNotFound
	}

	return portal.p, nil
}

func (s *swService) ListPortals(user *sworld.User) ([]*sworld.Portal, error) {
	portals := make([]*sworld.Portal, 0, len(s.portals))

	for _, portal := range s.portals {
		portals = append(portals, portal.p)
	}

	return portals, nil
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
	character, err := user.FindCharacter(characterID)
	if err != nil {
		return err
	}

	if character.Exploring {
		return ErrCharacterBusy
	}

	s.handleExplore(sportal.p, character)

	// TODO: What's this?
	sportal.c = &sCharacter{
		c: character,
		u: s.users[user.ID],
	}

	return nil
}

func (s *swService) OpenPortal(user *sworld.User, stoneID string) (*sworld.Portal, error) {
	if stoneID != "" {
		return nil, ErrPortalStoneNotFound
	}

	// TODO: if a stoneID is given, use that stone from the user inventory
	// Also, we could receive the coordinates of the item (bag, slot), instead of the ID
	// or we can just search everywhere until we find that stone
	stone := s.defaultStone(user)

	portal, err := sworld.OpenPortal(user, stone, func(portal *sworld.Portal) {
		log.Printf("Portal closed: %s\n", portal.ID)
		delete(s.portals, portal.ID)
	})
	if err != nil {
		return portal, err
	}
	log.Printf("Portal open: %s\n", portal.ID)

	sportal := &sPortal{
		p: portal,
	}

	s.portals[sportal.p.ID] = sportal

	return sportal.p, nil
}
