package psdb

import (
	"database/sql/driver"
	"io"

	"github.com/aki-0421/psdb/pkg/types"
)

type Rows struct {
	pos    int
	Fields []types.Field
	Rows   []types.Row
}

var _ driver.Rows = (*Rows)(nil)

func (r *Rows) Columns() []string {
	var cols []string
	for _, f := range r.Fields {
		cols = append(cols, f.Name)
	}
	return cols
}

func (r *Rows) Close() error {
	return nil
}

func (r *Rows) Next(dest []driver.Value) error {
	if r.pos+1 > len(r.Rows) {
		return io.EOF
	}

	row := r.Rows[r.pos]
	dr, err := row.Decode()
	if err != nil {
		return err
	}

	for i := 0; i != len(dr); i++ {
		dest[i] = dr[i]
	}

	r.pos++
	return nil
}
