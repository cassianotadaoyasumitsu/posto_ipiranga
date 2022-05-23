package ufo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"git.wealth-park.com/cassiano/posto_ipiranga/internal/log"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error implements the error interface.
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewError(code string, msg string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
	}
}

func EncodeHTTPErrorResponse(ctx context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	e := NewError("internal_error", err.Error())
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(e)
}

func RegisterRoutes(r *mux.Router, eps *Endpoints, l log.Logger) {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(EncodeHTTPErrorResponse),
	}
	r.Methods(http.MethodPost).Path("/ufos").Handler(
		httptransport.NewServer(
			eps.CreateUfo,
			decodeCreateUfoRequest,
			encodeCreateUfoResponse,
			options...,
		))
	l.Info("Added route", "method", http.MethodPost, "path", "/ufos")
	r.Methods(http.MethodGet).Path("/ufos/{id}").Handler(
		httptransport.NewServer(
			eps.OfID,
			decodeOfIDRequest,
			encodeOfIDResponse,
			options...,
		))
	l.Info("Added route", "method", http.MethodGet, "path", "/ufos/{id}")
	r.Methods(http.MethodGet).Path("/ufos").Handler(
		httptransport.NewServer(
			eps.List,
			httptransport.NopRequestDecoder,
			encodeListResponse,
			options...,
		))
	l.Info("Added route", "method", http.MethodPatch, "path", "/ufos")
}

func encodeCreateUfoResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	resp, ok := response.(*createUfoResponse)
	if !ok {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return json.NewEncoder(w).Encode(&resp.Ufo)
}

func decodeCreateUfoRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Model   string `json:"model"`
		Licence string `json:"licence"`
		Plate   string `json:"plate"`
		UfoTank int    `json:"tank"`
		Fuel    string `json:"fuel"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, fmt.Errorf("invalid request: %v", err)
	}
	return &createUfoRequest{
		Model:   body.Model,
		Licence: body.Licence,
		Plate:   body.Plate,
		UfoTank: body.UfoTank,
		Fuel:    body.Fuel,
	}, nil
}

func encodeOfIDResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	resp, ok := response.(*ofIDResponse)
	if !ok {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return json.NewEncoder(w).Encode(&resp.Ufo)
}

func decodeOfIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, fmt.Errorf("missing id")
	}
	return &ofIDRequest{UfoID: id}, nil
}

func encodeListResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	resp, ok := response.(*listResponse)
	if !ok {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return json.NewEncoder(w).Encode(&resp.Ufos)
}
