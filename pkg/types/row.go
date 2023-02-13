package types

import (
	"encoding/base64"
	"strconv"
)

type Row struct {
	Lengths []string `json:"lengths"`
	Values  string   `json:"values"`
}

func (r *Row) Decode() ([][]byte, error) {
	c := make([][]byte, len(r.Lengths))

	d, err := base64.StdEncoding.DecodeString(r.Values)
	if err != nil {
		return c, err
	}

	v := string(d)
	o := 0
	for i := range c {
		w, _ := strconv.Atoi(r.Lengths[i])
		if w < 0 {
			continue
		}

		c[i] = []byte(v[o : o+w])
		o += w
	}

	return c, nil
}
