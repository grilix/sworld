package server

import "github.com/grilix/sworld/sworldservice"

// MergeStonesRequest represents a request for merging stones
type MergeStonesRequest struct {
	SourceLocation ItemLocation `json:"source"`
	TargetLocation ItemLocation `json:"target"`
}

// OpenPortalRequest holds the parameters for the new portal
type OpenPortalRequest struct {
	StoneLocation *ItemLocation `json:"stone_location"`
}

// ViewPortalRequest represents a request for viewing a portal
type ViewPortalRequest struct {
	ID string `json:"id"`
}

// ListPortalsRequest represents a request for listing the portals
type ListPortalsRequest struct {
}

// ListZonesRequest holds the parameters for listing zones
type ListZonesRequest struct{}

// ExplorePortalRequest represents a request to explore a portal
type ExplorePortalRequest struct {
	PortalID    string `json:"portal_id"`
	CharacterID string `json:"character_id"`
}

// SpawnCharacterRequest represents a request for spawning a character
type SpawnCharacterRequest struct {
}

// AuthenticateRequest holds the credentials for authentication
type AuthenticateRequest struct {
	// TODO: use something else
	Credentials sworldservice.Credentials
}

// ViewUserInventoryRequest represents a request for viewing the user inventory
type ViewUserInventoryRequest struct {
}

// TakeCharacterItemRequest represents a request for dropping an item from the character inventory
type TakeCharacterItemRequest struct {
	CharacterID  string       `json:"id"`
	ItemLocation ItemLocation `json:"location"`
}

// DropCharacterItemRequest represents a request for dropping an item from the character inventory
type DropCharacterItemRequest struct {
	CharacterID  string       `json:"id"`
	ItemLocation ItemLocation `json:"location"`
}

// ViewCharacterInventoryRequest represents a request for viewing the character inventory
type ViewCharacterInventoryRequest struct {
	CharacterID string `json:"character_id"`
}

// ListCharactersRequest represents a request for listing the characters
type ListCharactersRequest struct {
}

// ViewCharacterRequest represents a request for viewing a character
type ViewCharacterRequest struct {
	ID string `json:"id"`
}
