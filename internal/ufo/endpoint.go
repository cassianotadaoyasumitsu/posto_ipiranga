package ufo

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
)

type Endpoints struct {
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

		if err != nil {
			return nil, fmt.Errorf("could not create ufo: %v", err)
		}

		u, err := svc.Create(ctx, req.Model, req.License, req.Plate, req.UfoTank, req.Fuel)
		if err != nil {
			busErr, ok := err.(*BusinessError)
			if !ok {
				return nil, fmt.Errorf("(Business Error) could not create ufo: %v", err)
			}
			return &createUfoResponse{
				Ufo: u,
				Err: busErr,
			}, nil
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
		ufoID, err := uuid.Parse(req.UfoID)
		if err != nil {
			return nil, fmt.Errorf("could not parse ufoID: %v", err)
		}
		u, err := svc.OfID(ctx, ufoID)
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
	License string
	Plate   string
	UfoTank int
	Fuel    string
}

type createUfoResponse struct {
	Ufo *Ufo
	Err error
}

type ofIDRequest struct {
	UfoID string
}

type ofIDResponse struct {
	Ufo *Ufo
	Err error
}

type listResponse struct {
	Ufos []*Ufo
	Err  error
}
