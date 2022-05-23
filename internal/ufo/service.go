package ufo

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

const (
	smallTank  = 100
	mediumTank = 200
	largeTank  = 300
)

var (
	tanks = []int{smallTank, mediumTank, largeTank}
)

type Ufo struct {
	ID      uuid.UUID `json:"id"`
	Model   string    `json:"model"`
	Licence string    `json:"licence"`
	Plate   string    `json:"plate"`
	Tank    int       `json:"tank"`
	Fuel    string    `json:"fuel"`
}

type BusinessError struct {
	Err error
}

func (e *BusinessError) Error() string {
	return e.Err.Error()
}

type Repository interface {
	OfID(ctx context.Context, ufoID uuid.UUID) (*Ufo, error)
	Persist(ctx context.Context, u *Ufo) error
	List(ctx context.Context) ([]*Ufo, error)
}

type Service interface {
	// CreateUfo creates a new Ufo.
	Create(ctx context.Context, model string, licence string, plate string, ufoTank int, fuel string) (*Ufo, error)
	// GetUfo returns the Ufo with the given ID.
	OfID(ctx context.Context, ufoID uuid.UUID) (*Ufo, error)
	// ListUfos returns all Ufos.
	List(ctx context.Context) ([]*Ufo, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func validateTankSize(ufoTank int) error {
	for _, tank := range tanks {
		if tank == ufoTank {
			return nil
		}
	}
	return fmt.Errorf("tank size %d is not supported", ufoTank)
}

func (s *service) Create(ctx context.Context, model string, licence string, plate string, ufoTank int, fuel string) (*Ufo, error) {
	err := validateTankSize(ufoTank)
	if err != nil {
		return nil, err
	}
	if model == "" {
		return nil, fmt.Errorf("ufo model is required")
	}
	if plate == "" || licence == "" {
		return nil, fmt.Errorf("ufo plate or licence is required")
	}
	if fuel == "" {
		return nil, fmt.Errorf("ufo fuel is required")
	}
	u := &Ufo{
		ID:      uuid.New(),
		Model:   model,
		Licence: licence,
		Plate:   plate,
		Tank:    ufoTank,
		Fuel:    fuel,
	}
	err = s.repo.Persist(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("failed to persist ufo: %w", err)
	}
	return u, nil
}

// OfID returns the Ufo with the given ID.
func (s *service) OfID(ctx context.Context, ufoID uuid.UUID) (*Ufo, error) {
	u, err := s.repo.OfID(ctx, ufoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ufo: %w", err)
	}
	return u, nil
}

// List returns all Ufos.
func (s *service) List(ctx context.Context) ([]*Ufo, error) {
	ufos, err := s.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list ufos: %w", err)
	}
	return ufos, nil
}
