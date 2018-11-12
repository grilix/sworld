package server

// TakeCharacterItemResponse represents a response after dropping an item from the character inventory
type TakeCharacterItemResponse struct {
	// TODO: what to respond here?
}

// DropCharacterItemResponse represents a response after dropping an item from the character inventory
type DropCharacterItemResponse struct {
	// TODO: what to respond here?
}

// ViewCharacterInventoryResponse represents a response with the character inventory
type ViewCharacterInventoryResponse struct {
	Bags []*BagDetails `json:"bags"`
}

// ListCharactersResponse represents a response with the characters list
type ListCharactersResponse struct {
	Characters []*CharacterDetails `json:"characters"`
}

// ViewCharacterResponse represents a response with the character details
// TODO: This could also include the inventory
type ViewCharacterResponse struct {
	Character *CharacterDetails `json:"character,omitempty"`
	Error     string            `json:"error,omitempty"`
}

// AuthenticateResponse holds the result of the authentication endpoint
type AuthenticateResponse struct {
	User  UserDetails `json:"user,omitempty"`
	Error string      `json:"error,omitempty"`
	Token string      `json:"token,omitempty"`
}

// ViewUserInventoryResponse represents a response with the user inventory
type ViewUserInventoryResponse struct {
	Bags []*BagDetails `json:"bags"`
}

// SpawnCharacterResponse represents a response for spawning a character
type SpawnCharacterResponse struct {
	Character *CharacterDetails `json:"character"`
}

// EmptyResponse represents a response that has no information to show
// FIXME: This actually sounds weird, can't we just use nil instead?
type EmptyResponse struct{}

// MergeStonesResponse represents the response of merging stones
type MergeStonesResponse struct {
	ResultLocation ItemLocation `json:"location"`
}

// OpenPortalResponse holds the result of creating a portal
type OpenPortalResponse struct {
	Portal *PortalDetails `json:"portal,omitempty"`
	Error  string         `json:"error,omitempty"`
}

// ViewPortalResponse represents a response with the portal details
type ViewPortalResponse struct {
	Portal *PortalDetails `json:"portal"`
}

// ListPortalsResponse represents a response with the portals list
type ListPortalsResponse struct {
	Portals []*PortalDetails `json:"portals"`
}

// ZoneInformationResponse holds information about a zone
type ZoneInformationResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ListZonesResponse holds the zones list
type ListZonesResponse struct {
	Zones []ZoneInformationResponse `json:"zones"`
	Error string                    `json:"error,omitempty"`
}

// ExplorePortalResponse represents the result of an explore portal request
// TODO: We might want to include extra information here, like, for example:
//   - Portal information?
//   - Time remaining
type ExplorePortalResponse struct {
	//Error string `json:"error,omitempty"`
}
