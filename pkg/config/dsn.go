package config

import (
	"errors"
	"net/url"
)

func (c *Config) FormatDSN() string {
	u := url.URL{
		Scheme: "https",
		Host:   c.Host,
		User: url.UserPassword(
			c.Username,
			c.Password,
		),
		Path: c.Database,
	}

	return u.String()
}

func ParseDSN(name string) (*Config, error) {
	u, err := url.Parse(name)
	if err != nil {
		return nil, err
	}

	pass, ok := u.User.Password()
	if !ok {
		return nil, errors.New("password not found")
	}

	return New(u.Host, u.User.Username(), pass, u.Path)
}
