package api

import (
	"context"
	"net/http"

	"github.com/aki-0421/psdb/pkg/types"
	"github.com/pquerna/ffjson/ffjson"
)

//go:generate ffjson $GOFILE
type CreateSessionRequest struct {
}

//go:generate ffjson $GOFILE
type CreateSessionResponse struct {
	Branch  string         `json:"branch"`
	User    *types.User    `json:"user"`
	Session *types.Session `json:"session"`
}

func (c *Client) CreateSession(ctx context.Context) error {
	// Create query parameters
	req := &CreateSessionRequest{}

	// Build request
	method := http.MethodPost
	path := "/psdb.v1alpha1.Database/CreateSession"
	body, err := ffjson.Marshal(req)
	if err != nil {
		return err
	}

	// Send request
	buf, err := c.SendRequest(ctx, method, path, body)
	if err != nil {
		return err
	}

	// Unmarshal json response
	var res CreateSessionResponse
	err = ffjson.Unmarshal(buf, &res)
	if err != nil {
		return err
	}

	c.Session = res.Session

	return nil
}

//go:generate ffjson $GOFILE
type ExecuteRequest struct {
	Query   string         `json:"query"`
	Session *types.Session `json:"session"`
}

//go:generate ffjson $GOFILE
type ExecuteResponse struct {
	Session *types.Session       `json:"session"`
	Result  *types.ExecuteResult `json:"result"`
	Error   *types.Error         `json:"error"`
	Timing  float64              `json:"timing"`
}

func (c *Client) Execute(ctx context.Context, query string) (*ExecuteResponse, error) {
	// Create query parameters
	req := &ExecuteRequest{
		Query:   query,
		Session: c.Session,
	}

	// Build request
	method := http.MethodPost
	path := "/psdb.v1alpha1.Database/Execute"
	body, err := ffjson.Marshal(req)
	if err != nil {
		return nil, err
	}

	// Send request
	buf, err := c.SendRequest(ctx, method, path, body)
	if err != nil {
		return nil, err
	}

	// Unmarshal json response
	var res ExecuteResponse
	err = ffjson.Unmarshal(buf, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
