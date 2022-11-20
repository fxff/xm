package model

import (
	"errors"

	"github.com/google/uuid"
)

type CompanyType string

const (
	Corporation        CompanyType = "corporation"
	NonProfit          CompanyType = "non-profit"
	Cooperative        CompanyType = "cooperative"
	SoleProprietorship CompanyType = "sole-proprietorship"
)

type Company struct {
	ID          uuid.UUID   `db:"id"`
	Name        string      `db:"name"`
	Description string      `db:"description"`
	Employees   int         `db:"employees"`
	Registered  bool        `db:"registered"`
	Type        CompanyType `db:"type"`
}

var (
	ErrNotFound = errors.New("not found")

	ErrNoName      = errors.New("no name")
	ErrNoEmployees = errors.New("no employees")
	ErrUnknownType = errors.New("unknown type")
)

func (c *Company) Validate() error {
	if c.Name == "" {
		return ErrNoName
	}
	if c.Employees == 0 {
		return ErrNoEmployees
	}

	switch c.Type {
	case Corporation, NonProfit, Cooperative, SoleProprietorship:
	default:
		return ErrUnknownType
	}

	return nil
}
