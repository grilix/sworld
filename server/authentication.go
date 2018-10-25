package server

import (
	"context"
	"errors"

	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	//"github.com/grilix/sworld/sworld"
	"github.com/grilix/sworld/sworldservice"
)

// AuthenticateRequest holds the credentials for authentication
type AuthenticateRequest struct {
	Credentials sworldservice.Credentials
}

// UserDetails represents a user
type UserDetails struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// AuthenticateResponse holds the result of the authentication endpoint
type AuthenticateResponse struct {
	User  UserDetails `json:"user,omitempty"`
	Error string      `json:"error,omitempty"`
	Token string      `json:"token,omitempty"`
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
				//return next(ctx, request)
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
