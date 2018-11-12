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
	// ErrWrongItem is when the item is not valid for an action
	ErrWrongItem = errors.New("The item is not valid for that action")
	// ErrAlreadyHasCharacters is when the user already has an alive character
	ErrAlreadyHasCharacters = errors.New("You already have characters")
	// ErrCharacterIsDead is when the character is dead
	ErrCharacterIsDead = errors.New("Character is dead")
)

// FIXME: I'd say we can get rid of these two
type sPortal struct {
	p *sworld.Portal
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

	SpawnCharacter(user *sworld.User) (*sworld.Character, error)
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
		defaultZone:           createDefaultZone(),
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
	if character.Health <= 0 {
		return ErrCharacterIsDead
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
	userID := user.ID

	for _, portal := range s.portals {
		if portal.p.User.ID == userID {
			portals = append(portals, portal.p)
		}
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
	if character.Health <= 0 {
		return ErrCharacterIsDead
	}

	exploration, err := character.EnterPortal(sportal.p)
	if err != nil {
		return err
	}

	s.handleCharacterAttack(exploration)
	s.handleCharacterMove(exploration)

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
	portal, err := s.openPortal(user, *stone)

	err = user.DropItem(bagID, slot)
	// TODO: unlock inventory
	if err != nil {
		// TODO: if an error occurs here, we should close the portal
		return nil, err
	}

	return portal, nil
}

func (s *swService) OpenDefaultPortal(user *sworld.User) (*sworld.Portal, error) {
	stone := s.defaultStone(user)
	portal, err := s.openPortal(user, stone)

	return portal, err
}

func (s *swService) SpawnCharacter(user *sworld.User) (*sworld.Character, error) {
	if user.HasAliveCharacters() {
		return nil, ErrAlreadyHasCharacters
	}

	character := sworld.NewCharacter()
	character.User = user
	user.Characters = append(user.Characters, character)

	s.characters[character.ID] = character

	return character, nil
}
