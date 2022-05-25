package inmem

import (
	"context"
	"fmt"

	"git.wealth-park.com/cassiano/posto_ipiranga/internal/ufo"
	"github.com/google/uuid"
)

type ufoRepository struct {
	store map[string]*ufo.Ufo
}

func NewUfoRepository() ufo.Repository {
	return &ufoRepository{
		store: make(map[string]*ufo.Ufo),
	}
}

func (r *ufoRepository) OfID(ctx context.Context, ufoID uuid.UUID) (*ufo.Ufo, error) {
	u, ok := r.store[ufoID.String()]
	if !ok {
		return nil, fmt.Errorf("ufo not found")
	}
	return u, nil
}

func (r ufoRepository) CreateUfo(ctx context.Context, u *ufo.Ufo) error {
	r.store[u.ID.String()] = u
	return nil
}

func (r ufoRepository) List(ctx context.Context) ([]*ufo.Ufo, error) {
	var ufos []*ufo.Ufo
	for _, u := range r.store {
		ufos = append(ufos, u)
	}
	return ufos, nil
}
