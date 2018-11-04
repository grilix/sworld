package server

import (
	"context"
	"errors"

	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/grilix/sworld/sworldservice"
)

// UserDetails represents a user
type UserDetails struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// AuthenticateRequest holds the credentials for authentication
type AuthenticateRequest struct {
	Credentials sworldservice.Credentials
}

// AuthenticateResponse holds the result of the authentication endpoint
type AuthenticateResponse struct {
	User  UserDetails `json:"user,omitempty"`
	Error string      `json:"error,omitempty"`
	Token string      `json:"token,omitempty"`
}

// ViewUserInventoryRequest represents a request for viewing the user inventory
type ViewUserInventoryRequest struct {
}

// ViewUserInventoryResponse represents a response with the user inventory
type ViewUserInventoryResponse struct {
	Bags []*BagDetails `json:"bags"`
}

type ctxSessionKeyType string

var (
	// ErrCantGenerateJWT means something went wrong when generating
	// the JWT token for the response
	ErrCantGenerateJWT = errors.New("JWT Token can't be generated")
)

var (
	jwtSecretKey                   = []byte("SigningString") // TODO: set key
	ctxUserKey   ctxSessionKeyType = "user"
)

func jwtKey(token *stdjwt.Token) (interface{}, error) {
	return jwtSecretKey, nil
}

// EmptyResponse represents a response that has no information to show
// FIXME: This actually sounds weird, can't we just use nil instead?
type EmptyResponse struct{}

// TODO: we need to improve this shit
func userFromContext(s sworldservice.Service) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			claims, ok := ctx.Value(jwt.JWTClaimsContextKey).(*stdjwt.StandardClaims)
			if !ok {
				return next(ctx, request)
			}

			user := s.FindUser(claims.Id)
			if user == nil {
				return EmptyResponse{}, ErrWrongToken
			}

			ctx = context.WithValue(ctx, ctxUserKey, user)
			return next(ctx, request)
		}
	}
}

func authenticatedEndpoint(
	s sworldservice.Service,
	endpointFactory func(s sworldservice.Service) endpoint.Endpoint,
) endpoint.Endpoint {
	return jwt.NewParser(
		jwtKey, stdjwt.SigningMethodHS256, jwt.StandardClaimsFactory,
	)(userFromContext(s)(endpointFactory(s)))
}
