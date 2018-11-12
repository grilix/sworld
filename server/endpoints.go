package server

import (
	"context"
	"errors"
	"fmt"

	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	"github.com/grilix/sworld/sworld"
	svc "github.com/grilix/sworld/sworldservice"
)

var (
	// ErrNoAccount means there's no user information in the request
	ErrNoAccount = errors.New("No account provided")

	// ErrWrongToken means the token is not valid
	ErrWrongToken = errors.New("Err wrong token")
)

// WrongRequestError is the error returned when an endpoint receives
// a request that can't be cast into it's expected type
// FIXME: Hmmmm
type WrongRequestError struct {
	Endpoint string
}

// Endpoints hold the endpoints
type Endpoints struct {
	AuthenticateEndpoint      endpoint.Endpoint
	ViewUserInventoryEndpoint endpoint.Endpoint
	MergeStonesEndpoint       endpoint.Endpoint

	SpawnCharacterEndpoint         endpoint.Endpoint
	ViewCharacterEndpoint          endpoint.Endpoint
	ListCharactersEndpoint         endpoint.Endpoint
	ViewCharacterInventoryEndpoint endpoint.Endpoint
	DropCharacterItemEndpoint      endpoint.Endpoint
	TakeCharacterItemEndpoint      endpoint.Endpoint

	OpenPortalEndpoint    endpoint.Endpoint
	ExplorePortalEndpoint endpoint.Endpoint
	ViewPortalEndpoint    endpoint.Endpoint
	ListPortalsEndpoint   endpoint.Endpoint
}

// MakeServerEndpoints creates an endpoints list for a server
func MakeServerEndpoints(s svc.Service) Endpoints {
	return Endpoints{
		AuthenticateEndpoint:      MakeAuthenticateEndpoint(s),
		ViewUserInventoryEndpoint: authenticatedEndpoint(s, MakeViewUserInventoryEndpoint),
		MergeStonesEndpoint:       authenticatedEndpoint(s, MakeMergeStonesEndpoint),

		SpawnCharacterEndpoint:         authenticatedEndpoint(s, MakeSpawnCharacterEndpoint),
		ViewCharacterEndpoint:          authenticatedEndpoint(s, MakeViewCharacterEndpoint),
		ListCharactersEndpoint:         authenticatedEndpoint(s, MakeListCharactersEndpoint),
		ViewCharacterInventoryEndpoint: authenticatedEndpoint(s, MakeViewCharacterInventoryEndpoint),
		DropCharacterItemEndpoint:      authenticatedEndpoint(s, MakeDropCharacterItemEndpoint),
		TakeCharacterItemEndpoint:      authenticatedEndpoint(s, MakeTakeCharacterItemEndpoint),

		OpenPortalEndpoint:    authenticatedEndpoint(s, MakeOpenPortalEndpoint),
		ExplorePortalEndpoint: authenticatedEndpoint(s, MakeExplorePortalEndpoint),
		ViewPortalEndpoint:    authenticatedEndpoint(s, MakeViewPortalEndpoint),
		ListPortalsEndpoint:   authenticatedEndpoint(s, MakeListPortalsEndpoint),
	}
}

func (e WrongRequestError) Error() string {
	return fmt.Sprintf("Endpoint %s received a wrong request type", e.Endpoint)
}

// MakeAuthenticateEndpoint creates the Authenticate endpoint
func MakeAuthenticateEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AuthenticateRequest)
		user, err := s.Authenticate(ctx, req.Credentials)
		if err != nil {
			return AuthenticateResponse{Error: err.Error()}, err
		}

		token := stdjwt.NewWithClaims(
			stdjwt.SigningMethodHS256, stdjwt.StandardClaims{
				Id: user.ID,
			},
		)

		tokenString, err := token.SignedString(jwtSecretKey)
		if err != nil {
			return AuthenticateResponse{Error: err.Error()}, ErrCantGenerateJWT
		}

		return AuthenticateResponse{
			User: UserDetails{
				ID:       user.ID,
				Username: user.Username,
			},
			Token: tokenString,
		}, nil
	}
}

// MakeSpawnCharacterEndpoint creates the SpawnCharacter endopint
func MakeSpawnCharacterEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		user, ok := ctx.Value(ctxUserKey).(*sworld.User)
		if !ok {
			return SpawnCharacterResponse{}, ErrNoAccount
		}

		character, err := s.SpawnCharacter(user)
		if err != nil {
			return SpawnCharacterResponse{}, err
		}

		return SpawnCharacterResponse{
			Character: &CharacterDetails{
				ID:        character.ID,
				Health:    character.Health,
				MaxHealth: character.MaxHealth,
				Exploring: character.Exploring,
			},
		}, nil
	}
}

// MakeViewCharacterEndpoint creates the ViewCharacter endopint
func MakeViewCharacterEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		user, ok := ctx.Value(ctxUserKey).(*sworld.User)
		if !ok {
			return ViewCharacterResponse{}, ErrNoAccount
		}

		charReq, ok := request.(ViewCharacterRequest)
		if !ok {
			return ViewCharacterResponse{}, WrongRequestError{Endpoint: "ViewCharacterEndpoint"}
		}

		character, err := user.FindCharacter(charReq.ID)
		if err != nil {
			return ViewCharacterResponse{}, err
		}

		return ViewCharacterResponse{
			Character: &CharacterDetails{
				ID:        character.ID,
				Health:    character.Health,
				MaxHealth: character.MaxHealth,
				Exploring: character.Exploring,
			},
		}, nil
	}
}

// MakeListCharactersEndpoint creates the ListCharacters endpoint
func MakeListCharactersEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user, ok := ctx.Value(ctxUserKey).(*sworld.User)
		if !ok {
			return ListCharactersResponse{}, ErrNoAccount
		}

		characters, err := s.ListCharacters(user)
		if err != nil {
			return ListCharactersResponse{}, err
		}

		charactersList := make([]*CharacterDetails, 0, len(characters))
		for _, character := range characters {
			charactersList = append(charactersList, &CharacterDetails{
				ID:        character.ID,
				Health:    character.Health,
				MaxHealth: character.MaxHealth,
				Exploring: character.Exploring,
			})
		}

		return ListCharactersResponse{
			Characters: charactersList,
		}, nil
	}
}

// MakeViewUserInventoryEndpoint creates the ViewUserInventory endpoint
func MakeViewUserInventoryEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user, ok := ctx.Value(ctxUserKey).(*sworld.User)
		if !ok {
			return ViewUserInventoryResponse{}, ErrNoAccount
		}

		inventory, err := s.ViewUserInventory(user)
		if err != nil {
			return ViewUserInventoryResponse{}, err
		}

		userBags := inventoryDetails(inventory)

		return ViewUserInventoryResponse{
			Bags: userBags,
		}, nil
	}
}

// MakeMergeStonesEndpoint creates the endpoint for dropping character items
func MakeMergeStonesEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user, ok := ctx.Value(ctxUserKey).(*sworld.User)
		if !ok {
			return MergeStonesResponse{}, ErrNoAccount
		}

		mergeReq, ok := request.(MergeStonesRequest)
		if !ok {
			return MergeStonesResponse{}, WrongRequestError{Endpoint: "MergeStones"}
		}

		source := sworld.ItemLocation{
			BagID: mergeReq.SourceLocation.BagID,
			Slot:  mergeReq.SourceLocation.Slot,
		}
		target := sworld.ItemLocation{
			BagID: mergeReq.TargetLocation.BagID,
			Slot:  mergeReq.TargetLocation.Slot,
		}

		location, err := s.MergeStones(user, source, target)

		return MergeStonesResponse{
			ResultLocation: ItemLocation{
				BagID: location.BagID,
				Slot:  location.Slot,
			},
		}, err
	}
}

// MakeViewCharacterInventoryEndpoint creates the ViewCharacterInventory endpoint
func MakeViewCharacterInventoryEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// TODO: Do we care about the user here?
		//
		// user, ok := ctx.Value(ctxUserKey).(*sworld.User)
		// if !ok {
		//     return ViewCharacterInventoryResponse{}, ErrNoAccount
		// }

		listReq, ok := request.(ViewCharacterInventoryRequest)
		if !ok {
			return ViewCharacterInventoryResponse{}, WrongRequestError{
				Endpoint: "ViewCharacterInventoryEndpoint",
			}
		}

		inventory, err := s.ViewCharacterInventory(listReq.CharacterID)
		if err != nil {
			return ViewCharacterInventoryResponse{}, err
		}

		characterBags := inventoryDetails(inventory)

		return ViewCharacterInventoryResponse{
			Bags: characterBags,
		}, nil
	}
}

// MakeListPortalsEndpoint creates the ListPortals endpoint
func MakeListPortalsEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user, ok := ctx.Value(ctxUserKey).(*sworld.User)
		if !ok {
			return ListPortalsResponse{}, ErrNoAccount
		}

		portals, err := s.ListPortals(user)
		if err != nil {
			return ListPortalsResponse{}, err
		}

		portalsList := make([]*PortalDetails, 0, len(portals))
		for _, portal := range portals {
			portalsList = append(portalsList, portalDetails(portal, true))
		}

		return ListPortalsResponse{
			Portals: portalsList,
		}, nil
	}
}

// MakeViewPortalEndpoint creates the ViewPortal endpoint
func MakeViewPortalEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// TODO: Do we care about the user here?
		//
		// user, ok := ctx.Value(ctxUserKey).(*sworld.User)
		// if !ok {
		//     return ViewPortalResponse{}, ErrNoAccount
		// }

		exploreReq, ok := request.(ViewPortalRequest)
		if !ok {
			return ViewPortalResponse{}, WrongRequestError{Endpoint: "ViewPortal"}
		}

		portal, err := s.ViewPortal(exploreReq.ID)

		return ViewPortalResponse{
			Portal: portalDetails(portal, false),
		}, err
	}
}

// MakeTakeCharacterItemEndpoint creates the endpoint for dropping character items
func MakeTakeCharacterItemEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user, ok := ctx.Value(ctxUserKey).(*sworld.User)
		if !ok {
			return TakeCharacterItemResponse{}, ErrNoAccount
		}

		dropReq, ok := request.(TakeCharacterItemRequest)
		if !ok {
			return TakeCharacterItemResponse{}, WrongRequestError{Endpoint: "TakeCharacterItem"}
		}

		err := s.TakeCharacterItem(user, dropReq.CharacterID, dropReq.ItemLocation.BagID, dropReq.ItemLocation.Slot)

		return TakeCharacterItemResponse{}, err
	}
}

// MakeDropCharacterItemEndpoint creates the endpoint for dropping character items
func MakeDropCharacterItemEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user, ok := ctx.Value(ctxUserKey).(*sworld.User)
		if !ok {
			return DropCharacterItemResponse{}, ErrNoAccount
		}

		dropReq, ok := request.(DropCharacterItemRequest)
		if !ok {
			return DropCharacterItemResponse{}, WrongRequestError{Endpoint: "DropCharacterItem"}
		}

		err := s.DropCharacterItem(user, dropReq.CharacterID, dropReq.ItemLocation.BagID, dropReq.ItemLocation.Slot)

		return DropCharacterItemResponse{}, err
	}
}

// MakeExplorePortalEndpoint creates the endpoint for exploring a portal
func MakeExplorePortalEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user, ok := ctx.Value(ctxUserKey).(*sworld.User)
		if !ok {
			return ExplorePortalResponse{}, ErrNoAccount
		}

		exploreReq, ok := request.(ExplorePortalRequest)
		if !ok {
			return ExplorePortalResponse{}, WrongRequestError{Endpoint: "ExplorePortal"}
		}

		err := s.ExplorePortal(user, exploreReq.PortalID, exploreReq.CharacterID)

		return ExplorePortalResponse{}, err
	}
}

// MakeOpenPortalEndpoint makes the endpoint for creating a portal
func MakeOpenPortalEndpoint(s svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user, ok := ctx.Value(ctxUserKey).(*sworld.User)
		if !ok {
			return OpenPortalResponse{}, ErrNoAccount
		}

		portalReq, ok := request.(OpenPortalRequest)
		if !ok {
			return OpenPortalResponse{}, WrongRequestError{Endpoint: "OpenPortal"}
		}
		var err error
		var portal *sworld.Portal

		location := portalReq.StoneLocation

		if location != nil {
			portal, err = s.OpenPortalWithStone(user, location.BagID, location.Slot)
		} else {
			portal, err = s.OpenDefaultPortal(user)
		}
		if err != nil {
			return OpenPortalResponse{Error: err.Error()}, err
		}

		return OpenPortalResponse{
			Portal: portalDetails(portal, false),
		}, nil
	}
}
