package psdb

import (
	"database/sql"
	"database/sql/driver"

	"github.com/aki-0421/psdb/pkg/api"
	"github.com/aki-0421/psdb/pkg/config"
)

type Driver struct{}

func (d Driver) Open(dsn string) (driver.Conn, error) {
	cnf, err := config.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}

	a, err := api.New(cnf)
	if err != nil {
		return nil, err
	}

	c := NewConnect(a)

	return c, nil
}

func init() {
	sql.Register("planetscale", &Driver{})
}
