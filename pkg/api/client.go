package api

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/aki-0421/psdb/pkg/config"
	"github.com/aki-0421/psdb/pkg/types"
)

type Client struct {
	http    *http.Client
	config  *config.Config
	Session *types.Session
}

func (c *Client) SendRequest(ctx context.Context, method, path string, body []byte) ([]byte, error) {
	// Create endpoint url
	endpoint := url.URL{
		Scheme: "https",
		Host:   c.config.Host,
		User: url.UserPassword(
			c.config.Username,
			c.config.Password,
		),
		Path: path,
	}

	// Init request
	req, err := http.NewRequestWithContext(ctx, method, endpoint.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// Basic Auth
	req.SetBasicAuth(c.config.Username, c.config.Password)

	// Add headers
	req.Header.Set("Host", c.config.Host)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "database-go")

	// Send request
	res, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Request failed
	if res.StatusCode != 200 {
		return nil, errors.New("failed to request")
	}

	return io.ReadAll(res.Body)
}

func New(config *config.Config) (*Client, error) {
	hc := &http.Client{
		Transport: http.DefaultTransport,
	}

	c := Client{
		http:   hc,
		config: config,
	}

	return &c, nil
}
