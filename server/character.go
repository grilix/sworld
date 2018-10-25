package server

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/grilix/sworld/sworld"
	"github.com/grilix/sworld/sworldservice"
)

var (
	ErrCharacterNotFound = errors.New("Character not found")
)

type CharacterDetails struct {
	ID        string `json:"id"`
	Exploring bool   `json:"exploring"`
}

type CharacterDetailsRequest struct {
}

type CharacterDetailsResponse struct {
	Character *CharacterDetails `json:"character,omitempty"`
	Error     string            `json:"error,omitempty"`
}

func MakeCharacterDetailsEndpoint(s sworldservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		user, ok := ctx.Value(ctxUserKey).(*sworld.User)
		if !ok {
			return CharacterDetailsResponse{}, ErrNoAccount
		}

		character := user.Character

		return CharacterDetailsResponse{
			Character: &CharacterDetails{
				ID:        character.ID,
				Exploring: character.Exploring,
			},
		}, nil
	}
}
