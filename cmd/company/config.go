package main

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"xm/internal/auth"
)

type (
	Config struct {
		Server  Server
		DB      DB
		Company Company
		Auth    Auth
	}

	Server struct {
		Port int
	}

	DB struct {
		Endpoint string
		Port     int
		Name     string
		User     string
		Password string
	}

	Company struct {
		Subroute string
	}

	Auth struct {
		Subroute string
		Auth     auth.Config
	}
)

func Parse(path string) (Config, error) {
	if path == "" {
		return Config{}, errors.New("no path")
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return Config{}, errors.Wrap(err, "read file")
	}

	expandedData := os.ExpandEnv(string(bytes))

	var config Config
	err = yaml.UnmarshalStrict([]byte(expandedData), &config)

	return config, errors.Wrap(err, "unmarshal")
}

//nolint:nosprintfhostport // called only by me
func (db DB) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		db.User, db.Password, db.Endpoint, db.Port, db.Name)
}

func (c *Company) Validate() {
	if c.Subroute != "" {
		return
	}

	c.Subroute = "/"
}
