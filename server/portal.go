package server

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/grilix/sworld/sworld"
	"github.com/grilix/sworld/sworldservice"
)

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
	StoneID string `json:"stone_id"`
}

// OpenPortalResponse holds the result of creating a portal
type OpenPortalResponse struct {
	Portal *PortalDetails `json:"portal,omitempty"`
	Error  string         `json:"error,omitempty"`
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
type ExplorePortalResponse struct {
	//Error string `json:"error,omitempty"`
}

// MakeExplorePortalEndpoint creates the endpoint for exploring a portal
func MakeExplorePortalEndpoint(s sworldservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user, ok := ctx.Value(ctxUserKey).(*sworld.User)
		if !ok {
			return ExplorePortalResponse{}, ErrNoAccount
		}

		exploreReq, ok := request.(ExplorePortalRequest)
		if !ok {
			return ExplorePortalResponse{}, WrongRequestError{Endpoint: "ExplorePortal"}
		}

		err := s.ExplorePortal(user, exploreReq.PortalID, exploreReq.CharacterID)

		// TODO: response?
		return ExplorePortalResponse{}, err
	}
}

// MakeOpenPortalEndpoint makes the endpoint for creating a portal
func MakeOpenPortalEndpoint(s sworldservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		user, ok := ctx.Value(ctxUserKey).(*sworld.User)
		if !ok {
			return OpenPortalResponse{}, ErrNoAccount
		}

		portalReq, ok := request.(OpenPortalRequest)
		if !ok {
			return OpenPortalResponse{}, WrongRequestError{Endpoint: "OpenPortal"}
		}

		portal, err := s.OpenPortal(user, portalReq.StoneID)
		if err != nil {
			return OpenPortalResponse{Error: err.Error()}, err
		}

		return OpenPortalResponse{
			Portal: &PortalDetails{
				ID:       portal.ID,
				Duration: portal.PortalStone.Duration.String(),
				Level:    portal.PortalStone.Level,
				Zone: ZoneDetails{
					ID:   portal.PortalStone.Zone.ID,
					Name: portal.PortalStone.Zone.Name,
				},
			},
		}, nil
	}
}
