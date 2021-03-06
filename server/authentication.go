package server

import (
	"context"
	"errors"

	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/grilix/sworld/sworldservice"
)

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
