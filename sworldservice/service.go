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
	// ErrWrongItem is when the item is not valid for an action
	ErrWrongItem = errors.New("The item is not valid for that action")
)

type sCharacter struct {
	c                *sworld.Character
	u                *sUser
	exploringPortal  *sworld.Portal
	exploredDIstance int
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
	MergeStones(user *sworld.User, source sworld.ItemLocation, target sworld.ItemLocation) (sworld.ItemLocation, error)

	ListCharacters(user *sworld.User) ([]*sworld.Character, error)
	ViewCharacterInventory(characterID string) ([]sworld.Bag, error)
	DropCharacterItem(user *sworld.User, characterID string, bagID, slot int) error
	TakeCharacterItem(user *sworld.User, characterID string, bagID, slot int) error

	OpenDefaultPortal(user *sworld.User) (*sworld.Portal, error)
	OpenPortalWithStone(user *sworld.User, bagID, slot int) (*sworld.Portal, error)
	ExplorePortal(user *sworld.User, portalID, characterID string) error
	ViewPortal(portalID string) (*sworld.Portal, error)
	ListPortals(user *sworld.User) ([]*sworld.Portal, error)
}

type swService struct {
	users                 map[string]*sUser
	portals               map[string]*sPortal
	defaultPortalDuration time.Duration
	characters            map[string]*sworld.Character
	defaultZone           *sworld.Zone
}

// NewService creates the service
func NewService() Service {
	return &swService{
		users:      make(map[string]*sUser),
		portals:    make(map[string]*sPortal),
		characters: make(map[string]*sworld.Character),
		// TODO: this should be on settings
		defaultPortalDuration: time.Second * 10,
		defaultZone: &sworld.Zone{
			ID:   sworld.RandomID(16),
			Name: "Forest",
			DropRate: sworld.DropRate{
				Gold:    8,
				Enemy:   20,
				Item:    20,
				Nothing: 60,
			},
		},
	}
}

func (s *swService) TakeCharacterItem(user *sworld.User, characterID string, bagID, slot int) error {
	// TODO: Fail if the character is exploring
	return user.TakeCharacterItem(characterID, sworld.ItemLocation{
		BagID: bagID,
		Slot:  slot,
	})
}

func (s *swService) DropCharacterItem(user *sworld.User, characterID string, bagID, slot int) error {
	// TODO: fail if the character is exploring
	character, err := user.FindCharacter(characterID)
	if err != nil {
		return err
	}
	_, err = character.DropItem(bagID, slot)

	return err
}

func (s *swService) MergeStones(user *sworld.User, source sworld.ItemLocation, target sworld.ItemLocation) (sworld.ItemLocation, error) {
	// FIXME: call user.MergeStones directly?
	return user.MergeStones(source, target)
}

func (s *swService) ViewUserInventory(user *sworld.User) ([]sworld.Bag, error) {
	// FIXME: call user.Bags directly?
	return user.Bags, nil
}

func (s *swService) ViewCharacterInventory(characterID string) ([]sworld.Bag, error) {
	// TODO: Fail if the character is explorig
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
	if sportal.p.User.ID != user.ID {
		return ErrCantEnterPortal
	}
	character, err := user.FindCharacter(characterID)
	if err != nil {
		return err
	}

	exploration, err := character.EnterPortal(sportal.p)
	if err != nil {
		return err
	}

	s.handleCharacterAttack(exploration)
	s.handleCharacterMove(exploration)

	// TODO: What's this?
	sportal.c = &sCharacter{
		c: character,
		u: s.users[user.ID],
	}

	return nil
}

func (s *swService) OpenPortalWithStone(user *sworld.User, bagID, slot int) (*sworld.Portal, error) {
	// TODO: lock inventory
	item, err := user.GetItem(sworld.ItemLocation{BagID: bagID, Slot: slot})
	if err != nil {
		return nil, err
	}
	stone, ok := item.(*sworld.PortalStone)
	if !ok {
		return nil, ErrWrongItem
	}

	portal, err := sworld.OpenPortal(user, *stone, func(portal *sworld.Portal) {
		log.Printf("Portal closed: %s\n", portal.ID)
		delete(s.portals, portal.ID)
	})
	if err != nil {
		return portal, err
	}
	err = user.DropItem(bagID, slot)
	// TODO: unlock inventory
	if err != nil {
		// TODO: if an error occurs here, we should close the portal
		return nil, err
	}

	log.Printf("Portal open: %s\n", portal.ID)

	sportal := &sPortal{
		p: portal,
	}

	s.portals[sportal.p.ID] = sportal

	return sportal.p, nil
}

func (s *swService) OpenDefaultPortal(user *sworld.User) (*sworld.Portal, error) {
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
