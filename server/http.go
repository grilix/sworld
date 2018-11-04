package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	klog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/grilix/sworld/sworldservice"
)

var (
	// ErrBadRouting is when a handler receives the wrong information from a route
	ErrBadRouting = errors.New("inconsistent mapping between route and handler")
)

type errorer interface {
	error() error
}

// MakeHTTPServer creates an http server for the endpoints
func MakeHTTPServer(s sworldservice.Service, logger klog.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerBefore(jwt.HTTPToContext()),
	}

	r.Methods("POST").Path("/api/v1/auth").Handler(AuthenticateHTTPServer(e, options))
	r.Methods("GET").Path("/api/v1/inventory").Handler(ViewUserInventoryHTTPServer(e, options))

	r.Methods("GET").Path("/api/v1/characters").Handler(ListCharactersHTTPServer(e, options))
	r.Methods("GET").Path("/api/v1/characters/{id}").Handler(ViewCharacterHTTPServer(e, options))
	r.Methods("GET").Path("/api/v1/characters/{id}/inventory").Handler(ViewCharacterInventoryHTTPServer(e, options))
	r.Methods("POST").Path("/api/v1/characters/{id}/drop").Handler(DropCharacterItemHTTPServer(e, options))
	r.Methods("POST").Path("/api/v1/characters/{id}/take").Handler(TakeCharacterItemHTTPServer(e, options))

	r.Methods("POST").Path("/api/v1/portals").Handler(OpenPortalHTTPServer(e, options))
	r.Methods("GET").Path("/api/v1/portals").Handler(ListPortalsHTTPServer(e, options))
	r.Methods("GET").Path("/api/v1/portal/{id}").Handler(ViewPortalHTTPServer(e, options))
	r.Methods("POST").Path("/api/v1/portals/{id}/explore").Handler(ExplorePortalHTTPServer(e, options))

	return r
}

// DumpRequest dumps the request to stdout for debugging
func DumpRequest() httptransport.RequestFunc {
	return func(ctx context.Context, r *http.Request) context.Context {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			return ctx
		}
		fmt.Println(string(dump))
		return ctx
	}
}

// DumpResponse dumps the response to stdout for debugging
func DumpResponse() httptransport.ClientResponseFunc {
	return func(ctx context.Context, r *http.Response) context.Context {
		dump, err := httputil.DumpResponse(r, true)
		if err != nil {
			return ctx
		}
		fmt.Println(string(dump))
		return ctx
	}
}

// MakeHTTPClientEndpoints initializes the endpoints to be used by a client
func MakeHTTPClientEndpoints(instance string) (Endpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	tgt, err := url.Parse(instance)
	if err != nil {
		return Endpoints{}, err
	}
	tgt.Path = ""

	options := []httptransport.ClientOption{
		httptransport.ClientBefore(jwt.ContextToHTTP()),
		httptransport.ClientBefore(DumpRequest()),
		httptransport.ClientAfter(DumpResponse()),
	}

	return Endpoints{
		AuthenticateEndpoint:      AuthenticateHTTPClient(tgt, options),
		ViewUserInventoryEndpoint: ViewUserInventoryHTTPClient(tgt, options),

		ViewCharacterEndpoint:          ViewCharacterHTTPClient(tgt, options),
		ListCharactersEndpoint:         ListCharactersHTTPClient(tgt, options),
		ViewCharacterInventoryEndpoint: ViewCharacterInventoryHTTPClient(tgt, options),
		DropCharacterItemEndpoint:      DropCharacterItemHTTPClient(tgt, options),
		TakeCharacterItemEndpoint:      TakeCharacterItemHTTPClient(tgt, options),

		OpenPortalEndpoint:    OpenPortalHTTPClient(tgt, options),
		ExplorePortalEndpoint: ExplorePortalHTTPClient(tgt, options),
		ListPortalsEndpoint:   ListPortalsHTTPClient(tgt, options),
		ViewPortalEndpoint:    ViewPortalHTTPClient(tgt, options),
	}, nil
}

// AuthenticateHTTPServer serves the AuthenticateEndpoint
func AuthenticateHTTPServer(endpoints Endpoints, options []httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(endpoints.AuthenticateEndpoint,
		func(_ context.Context, r *http.Request) (request interface{}, err error) {
			var req AuthenticateRequest
			if e := json.NewDecoder(r.Body).Decode(&req.Credentials); e != nil {
				return nil, e
			}
			return req, nil
		},
		func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
			if e, ok := response.(errorer); ok && e.error() != nil {
				encodeError(ctx, e.error(), w)
				return nil
			}
			authResponse, ok := response.(AuthenticateResponse)
			if !ok {
				encodeError(ctx, errors.New("Auth system error: Wrong response"), w)
				return nil
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", authResponse.Token))
			return json.NewEncoder(w).Encode(response)
		},
		options...,
	)
}

// AuthenticateHTTPClient calls the AuthenticateEndpoint
func AuthenticateHTTPClient(tgt *url.URL, options []httptransport.ClientOption) endpoint.Endpoint {
	return httptransport.NewClient("POST", tgt,
		func(ctx context.Context, req *http.Request, request interface{}) error {
			req.URL.Path = "/api/v1/auth"
			authRequest, ok := request.(AuthenticateRequest)
			if !ok {
				panic("Wrong request type")
			}
			return encodeRequest(ctx, req, authRequest.Credentials)
		},
		func(_ context.Context, resp *http.Response) (interface{}, error) {
			var response AuthenticateResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			return response, err
		},
		options...,
	).Endpoint()
}

// OpenPortalHTTPServer serves the OpenPortalEndpoint
func OpenPortalHTTPServer(endpoints Endpoints, options []httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(endpoints.OpenPortalEndpoint,
		func(_ context.Context, r *http.Request) (request interface{}, err error) {
			var req OpenPortalRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				return nil, err
			}

			return req, nil
		},
		encodeResponse,
		options...,
	)
}

// OpenPortalHTTPClient calls the OpenPortalEndpoint
func OpenPortalHTTPClient(tgt *url.URL, options []httptransport.ClientOption) endpoint.Endpoint {
	return httptransport.NewClient("POST", tgt,
		func(ctx context.Context, req *http.Request, request interface{}) error {
			req.URL.Path = "/api/v1/portals"
			portalRequest, ok := request.(OpenPortalRequest)
			if !ok {
				panic("Wrong request type")
			}
			return encodeRequest(ctx, req, portalRequest)
		},
		func(_ context.Context, resp *http.Response) (interface{}, error) {
			var response OpenPortalResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			return response, err
		},
		options...,
	).Endpoint()
}

// ListPortalsHTTPServer serves the ListPortalsEndpoint
func ListPortalsHTTPServer(endpoints Endpoints, options []httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(endpoints.ListPortalsEndpoint,
		func(_ context.Context, r *http.Request) (request interface{}, err error) {
			return ListPortalsRequest{}, nil
		},
		encodeResponse,
		options...,
	)
}

// ListPortalsHTTPClient calls the ListPortalsEndpoint
func ListPortalsHTTPClient(tgt *url.URL, options []httptransport.ClientOption) endpoint.Endpoint {
	return httptransport.NewClient("GET", tgt,
		func(ctx context.Context, req *http.Request, request interface{}) error {
			req.URL.Path = "/api/v1/portals"
			listRequest, ok := request.(ListPortalsRequest)
			if !ok {
				panic("Wrong request type")
			}
			return encodeRequest(ctx, req, listRequest)
		},
		func(_ context.Context, resp *http.Response) (interface{}, error) {
			var response ListPortalsResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			return response, err
		},
		options...,
	).Endpoint()
}

// ViewPortalHTTPServer serves the ViewPortalEndpoint
func ViewPortalHTTPServer(endpoints Endpoints, options []httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(endpoints.ViewPortalEndpoint,
		func(_ context.Context, r *http.Request) (request interface{}, err error) {
			vars := mux.Vars(r)
			id, ok := vars["id"]
			if !ok {
				return nil, ErrBadRouting
			}

			return ViewPortalRequest{
				ID: id,
			}, nil
		},
		encodeResponse,
		options...,
	)
}

// ViewPortalHTTPClient calls the ViewPortalEndpoint
func ViewPortalHTTPClient(tgt *url.URL, options []httptransport.ClientOption) endpoint.Endpoint {
	return httptransport.NewClient("GET", tgt,
		func(ctx context.Context, req *http.Request, request interface{}) error {
			portalRequest, ok := request.(ViewPortalRequest)
			if !ok {
				panic("Wrong request type")
			}
			req.URL.Path = fmt.Sprintf("/api/v1/portals/%s", portalRequest.ID)
			return encodeRequest(ctx, req, nil)
		},
		func(_ context.Context, resp *http.Response) (interface{}, error) {
			var response ViewPortalResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			return response, err
		},
		options...,
	).Endpoint()
}

// ViewUserInventoryHTTPServer serves the ViewUserInventoryEndpoint
func ViewUserInventoryHTTPServer(endpoints Endpoints, options []httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(endpoints.ViewUserInventoryEndpoint,
		func(_ context.Context, r *http.Request) (request interface{}, err error) {
			return ViewUserInventoryRequest{}, nil
		},
		encodeResponse,
		options...,
	)
}

// ViewUserInventoryHTTPClient calls the ViewUserInventoryEndpoint
func ViewUserInventoryHTTPClient(tgt *url.URL, options []httptransport.ClientOption) endpoint.Endpoint {
	return httptransport.NewClient("GET", tgt,
		func(ctx context.Context, req *http.Request, request interface{}) error {
			req.URL.Path = "/api/v1/inventory"
			return encodeRequest(ctx, req, nil)
		},
		func(_ context.Context, resp *http.Response) (interface{}, error) {
			var response ViewUserInventoryResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			return response, err
		},
		options...,
	).Endpoint()
}

// ViewCharacterInventoryHTTPServer serves the ViewCharacterInventoryEndpoint
func ViewCharacterInventoryHTTPServer(endpoints Endpoints, options []httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(endpoints.ViewCharacterInventoryEndpoint,
		func(_ context.Context, r *http.Request) (request interface{}, err error) {
			vars := mux.Vars(r)
			id, ok := vars["id"]
			if !ok {
				return nil, ErrBadRouting
			}

			return ViewCharacterInventoryRequest{
				CharacterID: id,
			}, nil
		},
		encodeResponse,
		options...,
	)
}

// ViewCharacterInventoryHTTPClient calls the ViewCharacterInventoryEndpoint
func ViewCharacterInventoryHTTPClient(tgt *url.URL, options []httptransport.ClientOption) endpoint.Endpoint {
	return httptransport.NewClient("GET", tgt,
		func(ctx context.Context, req *http.Request, request interface{}) error {
			listReq, ok := request.(ViewCharacterInventoryRequest)
			if !ok {
				panic("Wrong request type")
			}
			req.URL.Path = fmt.Sprintf("/api/v1/characters/%s/inventory", listReq.CharacterID)
			return encodeRequest(ctx, req, nil)
		},
		func(_ context.Context, resp *http.Response) (interface{}, error) {
			var response ViewCharacterInventoryResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			return response, err
		},
		options...,
	).Endpoint()
}

// ListCharactersHTTPServer serves the ListCharactersEndpoint
func ListCharactersHTTPServer(endpoints Endpoints, options []httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(endpoints.ListCharactersEndpoint,
		func(_ context.Context, r *http.Request) (request interface{}, err error) {
			return ListCharactersRequest{}, nil
		},
		encodeResponse,
		options...,
	)
}

// ListCharactersHTTPClient serves the ListCharactersEndpoint
func ListCharactersHTTPClient(tgt *url.URL, options []httptransport.ClientOption) endpoint.Endpoint {
	return httptransport.NewClient("GET", tgt,
		func(ctx context.Context, req *http.Request, request interface{}) error {
			req.URL.Path = "/api/v1/characters"
			return encodeRequest(ctx, req, nil)
		},
		func(_ context.Context, resp *http.Response) (interface{}, error) {
			var response ListCharactersResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			return response, err
		},
		options...,
	).Endpoint()
}

// ViewCharacterHTTPServer serves the ViewCharacterEndpoint
func ViewCharacterHTTPServer(endpoints Endpoints, options []httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(endpoints.ViewCharacterEndpoint,
		func(_ context.Context, r *http.Request) (request interface{}, err error) {
			vars := mux.Vars(r)
			id, ok := vars["id"]
			if !ok {
				return nil, ErrBadRouting
			}

			return ViewCharacterRequest{
				ID: id,
			}, nil
		},
		encodeResponse,
		options...,
	)
}

// ViewCharacterHTTPClient calls the ViewCharacterEndpoint
func ViewCharacterHTTPClient(tgt *url.URL, options []httptransport.ClientOption) endpoint.Endpoint {
	return httptransport.NewClient("GET", tgt,
		func(ctx context.Context, req *http.Request, request interface{}) error {
			characterReq, ok := request.(ViewPortalRequest)
			if !ok {
				panic("Wrong request type")
			}
			req.URL.Path = fmt.Sprintf("/api/v1/characters/%s", characterReq.ID)
			return encodeRequest(ctx, req, nil)
		},
		func(_ context.Context, resp *http.Response) (interface{}, error) {
			var response ViewCharacterResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			return response, err
		},
		options...,
	).Endpoint()
}

// TakeCharacterItemHTTPServer serves the TakeCharacterItemEndpoint
func TakeCharacterItemHTTPServer(endpoints Endpoints, options []httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(endpoints.TakeCharacterItemEndpoint,
		func(_ context.Context, r *http.Request) (request interface{}, err error) {
			vars := mux.Vars(r)
			id, ok := vars["id"]
			if !ok {
				return nil, ErrBadRouting
			}

			var req TakeCharacterItemRequest
			if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
				return nil, e
			}
			req.CharacterID = id

			return req, nil
		},
		encodeResponse,
		options...,
	)
}

// TakeCharacterItemHTTPClient calls the ViewCharacterEndpoint
func TakeCharacterItemHTTPClient(tgt *url.URL, options []httptransport.ClientOption) endpoint.Endpoint {
	return httptransport.NewClient("GET", tgt,
		func(ctx context.Context, req *http.Request, request interface{}) error {
			dropReq, ok := request.(TakeCharacterItemRequest)
			if !ok {
				panic("Wrong request type")
			}
			req.URL.Path = fmt.Sprintf("/api/v1/characters/%s/take", dropReq.CharacterID)
			return encodeRequest(ctx, req, req)
		},
		func(_ context.Context, resp *http.Response) (interface{}, error) {
			var response TakeCharacterItemResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			return response, err
		},
		options...,
	).Endpoint()
}

// DropCharacterItemHTTPServer serves the DropCharacterItemEndpoint
func DropCharacterItemHTTPServer(endpoints Endpoints, options []httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(endpoints.DropCharacterItemEndpoint,
		func(_ context.Context, r *http.Request) (request interface{}, err error) {
			vars := mux.Vars(r)
			id, ok := vars["id"]
			if !ok {
				return nil, ErrBadRouting
			}
			var req DropCharacterItemRequest
			if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
				return nil, e
			}
			req.CharacterID = id

			return req, nil
		},
		encodeResponse,
		options...,
	)
}

// DropCharacterItemHTTPClient calls the ViewCharacterEndpoint
func DropCharacterItemHTTPClient(tgt *url.URL, options []httptransport.ClientOption) endpoint.Endpoint {
	return httptransport.NewClient("GET", tgt,
		func(ctx context.Context, req *http.Request, request interface{}) error {
			dropReq, ok := request.(DropCharacterItemRequest)
			if !ok {
				panic("Wrong request type")
			}
			req.URL.Path = fmt.Sprintf("/api/v1/characters/%s/drop", dropReq.CharacterID)
			return encodeRequest(ctx, req, dropReq)
		},
		func(_ context.Context, resp *http.Response) (interface{}, error) {
			var response DropCharacterItemResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			return response, err
		},
		options...,
	).Endpoint()
}

// ExplorePortalHTTPServer serves the ExplorePortalEndpoint
func ExplorePortalHTTPServer(endpoints Endpoints, options []httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(endpoints.ExplorePortalEndpoint,
		func(_ context.Context, r *http.Request) (request interface{}, err error) {
			vars := mux.Vars(r)
			id, ok := vars["id"]
			if !ok {
				return nil, ErrBadRouting
			}

			var req ExplorePortalRequest
			if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
				return nil, e
			}
			req.PortalID = id

			return req, nil
		},
		encodeResponse,
		options...,
	)
}

// ExplorePortalHTTPClient calls the ExplorePortalEndpoint
func ExplorePortalHTTPClient(tgt *url.URL, options []httptransport.ClientOption) endpoint.Endpoint {
	return httptransport.NewClient("POST", tgt,
		func(ctx context.Context, req *http.Request, request interface{}) error {
			exploreReq, ok := request.(ExplorePortalRequest)
			if !ok {
				panic("Wrong request type")
			}

			req.URL.Path = fmt.Sprintf("/api/v1/portals/%s/explore", exploreReq.PortalID)
			return encodeRequest(ctx, req, request)
		},
		func(_ context.Context, resp *http.Response) (interface{}, error) {
			var response ExplorePortalResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			return response, err
		},
		options...,
	).Endpoint()
}

// encodeRequest likewise JSON-encodes the request to the HTTP request body.
// Don't use it directly as a transport/http.Client EncodeRequestFunc:
// profilesvc endpoints require mutating the HTTP method and request path.
func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrCharacterNotFound:
		return http.StatusNotFound
	case ErrNoAccount:
		return http.StatusUnauthorized
	default:
		switch err.(type) {
		case WrongRequestError:
			return http.StatusBadRequest
		default:
			return http.StatusInternalServerError
		}
	}
}
