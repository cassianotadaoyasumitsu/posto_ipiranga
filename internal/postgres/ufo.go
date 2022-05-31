package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"git.wealth-park.com/cassiano/posto_ipiranga/internal/ufo"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ufoRepository struct {
	db *sqlx.DB
}

func NewUfoRepository(db *sqlx.DB) ufo.Repository {
	return &ufoRepository{db: db}
}

type rawUfo struct {
	ID        uuid.UUID  `json:"id"`
	Model     string     `json:"model"`
	License   string     `json:"license"`
	Plate     string     `json:"plate"`
	Tank      int        `json:"tank"`
	Fuel      string     `json:"fuel"`
	CreatedAt *time.Time `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func (u *rawUfo) ToUfo() *ufo.Ufo {
	return &ufo.Ufo{
		ID:      u.ID,
		Model:   u.Model,
		License: u.License,
		Plate:   u.Plate,
		Tank:    u.Tank,
		Fuel:    u.Fuel,
	}
}

func (u *ufoRepository) OfID(ctx context.Context, ufoID uuid.UUID) (*ufo.Ufo, error) {
	var raw rawUfo
	err := u.db.GetContext(ctx, &raw, `
		SELECT id,
			   model,
			   license,
			   plate,
			   tank,
			   fuel,
			   created_at,
			   updated_at
		FROM ufos
		WHERE id = $1;
	`, ufoID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return raw.ToUfo(), nil
}

func (u *ufoRepository) CreateUfo(ctx context.Context, ufo *ufo.Ufo) error {
	_, err := u.db.ExecContext(ctx, `
		INSERT INTO ufos (id, model, license, plate, tank, fuel, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		ON CONFLICT(id) DO UPDATE SET plate=$4, tank=$5, fuel=$6, updated_at=$8;
	`, ufo.ID.String(), ufo.Model, ufo.License, ufo.Plate, ufo.Tank, ufo.Fuel, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("failed to persist ufo: %w", err)
	}
	return nil
}

func (u *ufoRepository) List(ctx context.Context) ([]*ufo.Ufo, error) {
	var raws []rawUfo
	err := u.db.SelectContext(ctx, &raws, `
		SELECT id,
			   model,
			   license,
			   plate,
			   tank,
			   fuel,
			   created_at,
			   updated_at
		FROM ufos
	`)
	if err != nil {
		return nil, err
	}
	ufos := make([]*ufo.Ufo, len(raws))
	for i, raw := range raws {
		ufos[i] = raw.ToUfo()
	}
	return ufos, nil
}
