package server

// ZoneDetails represents the details of a zone
type ZoneDetails struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// PortalDetails holds the details for a portal
type PortalDetails struct {
	ID       string      `json:"id"`
	Duration string      `json:"duration"`
	Level    int         `json:"level"`
	Zone     ZoneDetails `json:"zone"`
}

// OpenPortalRequest holds the parameters for the new portal
type OpenPortalRequest struct {
	StoneLocation *ItemLocation `json:"stone_location"`
}

// OpenPortalResponse holds the result of creating a portal
type OpenPortalResponse struct {
	Portal *PortalDetails `json:"portal,omitempty"`
	Error  string         `json:"error,omitempty"`
}

// ViewPortalRequest represents a request for viewing a portal
type ViewPortalRequest struct {
	ID string `json:"id"`
}

// ViewPortalResponse represents a response with the portal details
type ViewPortalResponse struct {
	Portal *PortalDetails `json:"portal"`
}

// ListPortalsRequest represents a request for listing the portals
type ListPortalsRequest struct {
}

// ListPortalsResponse represents a response with the portals list
type ListPortalsResponse struct {
	Portals []*PortalDetails `json:"portals"`
}

// ListZonesRequest holds the parameters for listing zones
type ListZonesRequest struct{}

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

// ExplorePortalRequest represents a request to explore a portal
type ExplorePortalRequest struct {
	PortalID    string `json:"portal_id"`
	CharacterID string `json:"character_id"`
}

// ExplorePortalResponse represents the result of an explore portal request
// TODO: We might want to include extra information here, like, for example:
//   - Portal information?
//   - Time remaining
type ExplorePortalResponse struct {
	//Error string `json:"error,omitempty"`
}
