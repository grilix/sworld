package server

import (
	"context"
	"errors"
	"fmt"

	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/endpoint"
	"github.com/grilix/sworld/sworld"
	"github.com/grilix/sworld/sworldservice"
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
	AuthenticateEndpoint endpoint.Endpoint

	ViewCharacterEndpoint          endpoint.Endpoint
	ListCharactersEndpoint         endpoint.Endpoint
	ViewCharacterInventoryEndpoint endpoint.Endpoint

	OpenPortalEndpoint    endpoint.Endpoint
	ExplorePortalEndpoint endpoint.Endpoint
	ViewPortalEndpoint    endpoint.Endpoint
	ListPortalsEndpoint   endpoint.Endpoint
}

// MakeServerEndpoints creates an endpoints list for a server
func MakeServerEndpoints(s sworldservice.Service) Endpoints {
	return Endpoints{
		AuthenticateEndpoint: MakeAuthenticateEndpoint(s),

		ViewCharacterEndpoint:          authenticatedEndpoint(s, MakeViewCharacterEndpoint),
		ListCharactersEndpoint:         authenticatedEndpoint(s, MakeListCharactersEndpoint),
		ViewCharacterInventoryEndpoint: authenticatedEndpoint(s, MakeViewCharacterInventoryEndpoint),

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
func MakeAuthenticateEndpoint(s sworldservice.Service) endpoint.Endpoint {
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

// MakeViewCharacterEndpoint creates the ViewCharacter endopint
func MakeViewCharacterEndpoint(s sworldservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		user, ok := ctx.Value(ctxUserKey).(*sworld.User)
		if !ok {
			return ViewCharacterResponse{}, ErrNoAccount
		}

		character := user.Character

		return ViewCharacterResponse{
			Character: &CharacterDetails{
				ID:        character.ID,
				Exploring: character.Exploring,
			},
		}, nil
	}
}

// MakeListCharactersEndpoint creates the ListCharacters endpoint
func MakeListCharactersEndpoint(s sworldservice.Service) endpoint.Endpoint {
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
				Exploring: character.Exploring,
			})
		}

		return ListCharactersResponse{
			Characters: charactersList,
		}, nil
	}
}

// MakeViewCharacterInventoryEndpoint creates the ViewCharacterInventory endpoint
func MakeViewCharacterInventoryEndpoint(s sworldservice.Service) endpoint.Endpoint {
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
func MakeListPortalsEndpoint(s sworldservice.Service) endpoint.Endpoint {
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
			portalsList = append(portalsList, &PortalDetails{
				ID:       portal.ID,
				Duration: portal.PortalStone.Duration.String(),
				Level:    portal.PortalStone.Level,
				Zone: ZoneDetails{
					ID:   portal.PortalStone.Zone.ID,
					Name: portal.PortalStone.Zone.Name,
				},
			})
		}

		return ListPortalsResponse{
			Portals: portalsList,
		}, nil
	}
}

// MakeViewPortalEndpoint creates the ViewPortal endpoint
func MakeViewPortalEndpoint(s sworldservice.Service) endpoint.Endpoint {
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
			Portal: &PortalDetails{
				ID:       portal.ID,
				Duration: portal.PortalStone.Duration.String(),
				Level:    portal.PortalStone.Level,
				Zone: ZoneDetails{
					ID:   portal.PortalStone.Zone.ID,
					Name: portal.PortalStone.Zone.Name,
				},
			},
		}, err
	}
}

// MakeExplorePortalEndpoint creates the endpoint for exploring a portal
func MakeExplorePortalEndpoint(s sworldservice.Service) endpoint.Endpoint {
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
func MakeOpenPortalEndpoint(s sworldservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user, ok := ctx.Value(ctxUserKey).(*sworld.User)
		if !ok {
			return OpenPortalResponse{}, ErrNoAccount
		}

		portalReq, ok := request.(OpenPortalRequest)
		if !ok {
			return OpenPortalResponse{}, WrongRequestError{Endpoint: "OpenPortal"}
		}

		portal, err := s.OpenPortal(user, portalReq.StoneID)
		if err != nil {
			return OpenPortalResponse{Error: err.Error()}, err
		}

		return OpenPortalResponse{
			Portal: &PortalDetails{
				ID:       portal.ID,
				Duration: portal.PortalStone.Duration.String(),
				Level:    portal.PortalStone.Level,
				Zone: ZoneDetails{
					ID:   portal.PortalStone.Zone.ID,
					Name: portal.PortalStone.Zone.Name,
				},
			},
		}, nil
	}
}
