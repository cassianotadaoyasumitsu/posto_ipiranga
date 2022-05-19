package ufo

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type Endpoint struct {
	CreateUfo endpoint.Endpoint
	OfID      endpoint.Endpoint
	List      endpoint.Endpoint
}

// Create a new Ufo.
func NewCreateUfoEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*createUfoRequest)
		if !ok {
			return nil, fmt.Errorf("could not assert request to *createUfoRequest")
		}
		u, err := svc.Create(ctx, req.Model, req.Licence, req.Plate, req.UfoTank, req.Fuel)
		if err != nil {
			return nil, fmt.Errorf("could not create ufo: %v", err)
		}
		response = &createUfoResponse{
			Ufo: u,
			Err: nil,
		}
		return response, nil
	}
}

// Get a Ufo by ID.
func NewOfIDEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(*ofIDRequest)
		if !ok {
			return nil, fmt.Errorf("could not assert request to *ofIDRequest")
		}
		u, err := svc.OfID(ctx, req.UfoID)
		if err != nil {
			return nil, fmt.Errorf("could not get ufo: %v", err)
		}
		response = &ofIDResponse{
			Ufo: u,
			Err: nil,
		}
		return response, nil
	}
}

// Get all Ufos.
func NewListEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		u, err := svc.List(ctx)
		if err != nil {
			return nil, fmt.Errorf("could not get ufos: %v", err)
		}
		response = &listResponse{
			Ufos: u,
			Err:  nil,
		}
		return response, nil
	}
}

type createUfoRequest struct {
	Model   string
	Licence string
	Plate   string
	UfoTank int
	Fuel    string
}

type createUfoResponse struct {
	Ufo *Ufo
	Err error
}

type ofIDRequest struct {
	UfoID uuid.UUID
}

type ofIDResponse struct {
	Ufo *Ufo
	Err error
}

type listResponse struct {
	Ufos []*Ufo
	Err  error
}
