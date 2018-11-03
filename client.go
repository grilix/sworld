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

func listPortals(ctx context.Context, client *Client) (server.ListPortalsResponse, error) {
	req := server.ListPortalsRequest{}

	res, err := client.e.ListPortalsEndpoint(ctx, req)
	if err != nil {
		return server.ListPortalsResponse{}, err
	}
	portalRes, ok := res.(server.ListPortalsResponse)
	if !ok {
		return server.ListPortalsResponse{}, ErrWrongResponse
	}

	return portalRes, nil
}

func listCharacters(ctx context.Context, client *Client) (server.ListCharactersResponse, error) {
	req := server.ListCharactersRequest{}
	res, err := client.e.ListCharactersEndpoint(ctx, req)
	if err != nil {
		return server.ListCharactersResponse{}, err
	}

	charRes, ok := res.(server.ListCharactersResponse)
	if !ok {
		return server.ListCharactersResponse{}, ErrWrongResponse
	}

	return charRes, nil
}

func characterInventory(ctx context.Context, client *Client, characterID string) (server.ViewCharacterInventoryResponse, error) {
	req := server.ViewCharacterInventoryRequest{
		CharacterID: characterID,
	}
	res, err := client.e.ViewCharacterInventoryEndpoint(ctx, req)
	if err != nil {
		return server.ViewCharacterInventoryResponse{}, err
	}

	charRes, ok := res.(server.ViewCharacterInventoryResponse)
	if !ok {
		return server.ViewCharacterInventoryResponse{}, ErrWrongResponse
	}

	return charRes, nil
}

func characterDetails(ctx context.Context, client *Client, characterID string) (server.ViewCharacterResponse, error) {
	req := server.ViewCharacterRequest{
		ID: characterID,
	}
	res, err := client.e.ViewCharacterEndpoint(ctx, req)
	if err != nil {
		return server.ViewCharacterResponse{}, err
	}

	charRes, ok := res.(server.ViewCharacterResponse)
	if !ok {
		return server.ViewCharacterResponse{}, ErrWrongResponse
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

	//reader := bufio.NewReader(os.Stdin)
	fmt.Print(" Username: ")

	var username string
	_, err = fmt.Scanln(&username)
	if err != nil {
		panic("Can't read from stdin")
	}

	_, err = authenticate(ctx, client, svc.Credentials{
		Username: username,
	})
	if err != nil {
		panic(err)
	}
	ctx = context.WithValue(ctx, jwt.JWTTokenContextKey, client.token)

	var character sworld.Character
	var portal sworld.Portal

	// characters
	charRes, err := listCharacters(ctx, client)
	if err != nil {
		fmt.Println(err.Error())
	}
	character.ID = charRes.Characters[0].ID

	// open portal
	portalRes, err := openPortal(ctx, client)
	if err != nil {
		panic(err)
	}
	portal.ID = portalRes.Portal.ID

	// explore
	_, err = explorePortal(ctx, client, character.ID, portal.ID)
	if err != nil {
		panic(err)
	}

	for {
		var command string
		fmt.Print(" -> ")
		_, err := fmt.Scanln(&command)
		if err != nil {
			panic("Can't read from stdin")
		}
		switch command {
		case "characters":
			charRes, err := listCharacters(ctx, client)
			if err != nil {
				fmt.Println(err.Error())
			}
			character.ID = charRes.Characters[0].ID
		case "open-portal":
			portalRes, err := openPortal(ctx, client)
			if err != nil {
				panic(err)
			}
			portal.ID = portalRes.Portal.ID
		case "inventory":
			_, err = characterInventory(ctx, client, character.ID)
			if err != nil {
				panic(err)
			}
		case "explore":
			_, err = explorePortal(ctx, client, character.ID, portal.ID)
			if err != nil {
				panic(err)
			}
		case "portals":
			portalsRes, err := listPortals(ctx, client)
			if err != nil {
				panic(err)
			}

			for _, portal := range portalsRes.Portals {
				fmt.Printf(" -> %s\n", portal.ID)
			}
		case "q", "quit":
			return
		default:
			fmt.Printf(" Invalid command: %v\n", command)
		}
	}
}
