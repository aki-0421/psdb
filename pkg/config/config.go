package config

import (
	"errors"
)

type Config struct {
	Host     string
	Username string
	Password string
	Database string
}

func New(host, user, pass, db string) (*Config, error) {
	if host == "" {
		return nil, errors.New("cannot be left host blank")
	}
	if user == "" {
		return nil, errors.New("cannot be left user blank")
	}
	if pass == "" {
		return nil, errors.New("cannot be left pass blank")
	}
	if db == "" {
		return nil, errors.New("cannot be left db blank")
	}

	c := Config{
		Host:     host,
		Username: user,
		Password: pass,
		Database: db,
	}

	return &c, nil
}
