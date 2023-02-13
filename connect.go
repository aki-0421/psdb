package psdb

import (
	"context"
	"database/sql/driver"
	"errors"

	"github.com/aki-0421/psdb/pkg/api"
)

type Connect struct {
	api *api.Client
}

func (c *Connect) Close() error {
	return nil
}

func (c *Connect) Prepare(query string) (driver.Stmt, error) {
	return nil, errors.New("prepare method not implemented")
}

func (c *Connect) Begin() (driver.Tx, error) {
	return nil, errors.New("begin method not implemented")
}

func (c *Connect) Rollback() (driver.Stmt, error) {
	return nil, errors.New("rollback method not implemented")
}

func (c *Connect) Query(query string, args []driver.Value) (driver.Rows, error) {
	ctx := context.Background()

	return c.QueryContext(ctx, query, args)
}

func (c *Connect) QueryContext(ctx context.Context, query string, args []driver.Value) (driver.Rows, error) {
	// Create session
	if c.api.Session == nil {
		err := c.api.CreateSession(ctx)
		if err != nil {
			return nil, err
		}
	}

	res, err := c.api.Execute(ctx, query)
	if err != nil {
		return nil, err
	}

	return &Rows{
		Fields: res.Result.Fields,
		Rows:   res.Result.Rows,
	}, nil
}

func NewConnect(api *api.Client) *Connect {
	c := Connect{
		api: api,
	}

	return &c
}
