package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/grilix/sworld/server"
	"github.com/grilix/sworld/sworld"
	svc "github.com/grilix/sworld/sworldservice"
)

// Client holds the client data
type Client struct {
	e     server.Endpoints
	token string
	user  *sworld.User
}

var (
	// ErrWrongResponse means the server returned something we didn't expect
	ErrWrongResponse = errors.New("Got wrong response from service")
)

func openPortal(ctx context.Context, client *Client) (server.OpenPortalResponse, error) {
	req := server.OpenPortalRequest{
		StoneID: "",
	}

	res, err := client.e.OpenPortalEndpoint(ctx, req)
	if err != nil {
		return server.OpenPortalResponse{}, err
	}
	portalRes, ok := res.(server.OpenPortalResponse)
	if !ok {
		return server.OpenPortalResponse{}, ErrWrongResponse
	}
	fmt.Printf(" Portal open: %s\n", portalRes.Portal.ID)

	return portalRes, nil
}

func characterDetails(ctx context.Context, client *Client) (server.CharacterDetailsResponse, error) {
	req := server.CharacterDetailsRequest{}
	res, err := client.e.CharacterDetailsEndpoint(ctx, req)
	if err != nil {
		return server.CharacterDetailsResponse{}, err
	}

	charRes, ok := res.(server.CharacterDetailsResponse)
	if !ok {
		return server.CharacterDetailsResponse{}, ErrWrongResponse
	}

	return charRes, nil
}

func explorePortal(ctx context.Context, client *Client, characterID, portalID string) (server.ExplorePortalResponse, error) {
	req2 := server.ExplorePortalRequest{
		PortalID:    portalID,
		CharacterID: characterID,
	}
	res, err := client.e.ExplorePortalEndpoint(ctx, req2)
	if err != nil {
		return server.ExplorePortalResponse{}, err
	}

	exploreRes, ok := res.(server.ExplorePortalResponse)
	if !ok {
		return server.ExplorePortalResponse{}, ErrWrongResponse
	}

	return exploreRes, nil
}

func authenticate(ctx context.Context, client *Client, c svc.Credentials) (server.AuthenticateResponse, error) {
	req := server.AuthenticateRequest{Credentials: c}

	res, err := client.e.AuthenticateEndpoint(ctx, req)
	if err != nil {
		return server.AuthenticateResponse{}, err
	}
	authResp, ok := res.(server.AuthenticateResponse)
	if !ok {
		return server.AuthenticateResponse{}, ErrWrongResponse
	}
	client.token = authResp.Token
	client.user = &sworld.User{
		ID: authResp.User.Username,
	}

	return authResp, nil
}

func main() {
	e, err := server.MakeHTTPClientEndpoints("localhost:8089")
	if err != nil {
		panic(err)
	}
	client := &Client{e: e}

	ctx := context.TODO()

	_, err = authenticate(ctx, client, svc.Credentials{
		Username: "grilix",
	})
	if err != nil {
		panic(err)
	}
	ctx = context.WithValue(ctx, jwt.JWTTokenContextKey, client.token)

	character, err := characterDetails(ctx, client)
	if err != nil {
		panic(err)
	}

	if character.Character.Exploring {
		fmt.Println("Character is already exploring.")
	} else {
		portal, err := openPortal(ctx, client)
		if err != nil {
			panic(err)
		}

		_, err = explorePortal(ctx, client, character.Character.ID, portal.Portal.ID)
		if err != nil {
			panic(err)
		}
	}
}
