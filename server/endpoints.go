package server

import (
	"errors"
	"fmt"

	"github.com/go-kit/kit/endpoint"
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
type WrongRequestError struct {
	Endpoint string
}

// Endpoints hold the endpoints
type Endpoints struct {
	AuthenticateEndpoint     endpoint.Endpoint
	CharacterDetailsEndpoint endpoint.Endpoint
	OpenPortalEndpoint       endpoint.Endpoint
	ExplorePortalEndpoint    endpoint.Endpoint
}

// MakeServerEndpoints creates an endpoints list for a server
func MakeServerEndpoints(s sworldservice.Service) Endpoints {
	return Endpoints{
		AuthenticateEndpoint:     MakeAuthenticateEndpoint(s),
		CharacterDetailsEndpoint: authenticatedEndpoint(s, MakeCharacterDetailsEndpoint),
		OpenPortalEndpoint:       authenticatedEndpoint(s, MakeOpenPortalEndpoint),
		ExplorePortalEndpoint:    authenticatedEndpoint(s, MakeExplorePortalEndpoint),
	}
}

func (e WrongRequestError) Error() string {
	return fmt.Sprintf("Endpoint %s received a wrong request type", e.Endpoint)
}
