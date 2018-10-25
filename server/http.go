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

// TODO: check which ones are being used
var (
	ErrBadRouting = errors.New(
		"inconsistent mapping between route and handler (programmer error)",
	)
	ErrNotSignedIn = errors.New("Not signed in")
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
	r.Methods("GET").Path("/api/v1/character").Handler(CharacterDetailsHTTPServer(e, options))
	r.Methods("POST").Path("/api/v1/portals").Handler(OpenPortalHTTPServer(e, options))
	r.Methods("POST").Path("/api/v1/portals/{id}/explore").Handler(ExplorePortalHTTPServer(e, options))

	return r
}

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
		AuthenticateEndpoint:     AuthenticateHTTPClient(tgt, options),
		CharacterDetailsEndpoint: CharacterDetailsHTTPClient(tgt, options),
		OpenPortalEndpoint:       OpenPortalHTTPClient(tgt, options),
		ExplorePortalEndpoint:    ExplorePortalHTTPClient(tgt, options),
	}, nil
}

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

func OpenPortalHTTPServer(endpoints Endpoints, options []httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(endpoints.OpenPortalEndpoint,
		func(_ context.Context, r *http.Request) (request interface{}, err error) {
			var req OpenPortalRequest
			if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
				return nil, e
			}

			return req, nil
		},
		encodeResponse,
		options...,
	)
}

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

func CharacterDetailsHTTPServer(endpoints Endpoints, options []httptransport.ServerOption) *httptransport.Server {
	return httptransport.NewServer(endpoints.CharacterDetailsEndpoint,
		func(_ context.Context, r *http.Request) (request interface{}, err error) {
			return CharacterDetailsRequest{}, nil
		},
		encodeResponse,
		options...,
	)
}

func CharacterDetailsHTTPClient(tgt *url.URL, options []httptransport.ClientOption) endpoint.Endpoint {
	return httptransport.NewClient("GET", tgt,
		func(ctx context.Context, req *http.Request, request interface{}) error {
			req.URL.Path = "/api/v1/character"
			return encodeRequest(ctx, req, nil)
		},
		func(_ context.Context, resp *http.Response) (interface{}, error) {
			var response CharacterDetailsResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			return response, err
		},
		options...,
	).Endpoint()
}

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
