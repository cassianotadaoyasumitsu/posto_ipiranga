package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	sqltrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/database/sql"

	// required for migrations
	ddsqlx "gopkg.in/DataDog/dd-trace-go.v1/contrib/jmoiron/sqlx"
	// Enable PostgreSQL support
	_ "github.com/lib/pq"
)

const driver = "postgres"

var ErrInvalidDSN = errors.New("invalid database URL, valid format: postgres://user:password@host:port/db?sslmode=mode")

type Conn struct {
	DB        *sqlx.DB
	MaskedDSN string
}

func NewConnection(dsn string, maxConn int, connectionServiceName string) (*Conn, error) {
	// Make sure to only accept urls in a standard format
	if !strings.HasPrefix(dsn, "postgres://") && !strings.HasPrefix(dsn, "postgresql://") {
		return nil, ErrInvalidDSN
	}
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, ErrInvalidDSN
	}

	mDSN := fmt.Sprintf("%s://%s:%s@%s%s", u.Scheme, u.User.Username(), "********", u.Host, u.RequestURI())

	sqltrace.Register(driver, &pq.Driver{}, sqltrace.WithAnalytics(true), sqltrace.WithAnalyticsRate(0.2), sqltrace.WithServiceName(fmt.Sprintf("%s.db", connectionServiceName)))
	db, err := ddsqlx.Connect(driver, u.String())
	if err != nil {
		return nil, err
	}
	if maxConn > 0 {
		db.SetMaxOpenConns(maxConn)
	}
	c := &Conn{
		DB:        db,
		MaskedDSN: mDSN,
	}
	return c, nil
}

func NewConnectionFromDB(db *sql.DB) *Conn {
	return &Conn{
		DB:        sqlx.NewDb(db, driver),
		MaskedDSN: "",
	}
}

func (c *Conn) Close() error {
	return c.DB.Close()
}
