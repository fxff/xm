package storage

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"

	"xm/internal/model"
)

type storage struct {
	conn *pgx.Conn
}

const (
	table = "companies.companies"
)

var (
	insertColumns = []string{"name", "description", "employees", "registered", "type"}
	selectColumns = append([]string{"id"}, insertColumns...)
)

func New(conn *pgx.Conn) *storage {
	return &storage{conn: conn}
}

func (s *storage) Insert(ctx context.Context, company *model.Company) (uuid.UUID, error) {
	b := sq.Insert(table).
		Columns(insertColumns...).
		Values(
			company.Name,
			company.Description,
			company.Employees,
			company.Registered,
			company.Type,
		).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	query, args, err := b.ToSql()
	if err != nil {
		return uuid.Nil, errors.Wrap(err, "build query")
	}

	var uuid uuid.UUID
	err = pgxscan.Get(ctx, s.conn, &uuid, query, args...)
	return uuid, errors.Wrap(err, "exec")
}

func (s *storage) Update(ctx context.Context, companyID uuid.UUID, company *model.Company) error {
	b := sq.Update(table).
		Set("name", company.Name).
		Set("description", company.Description).
		Set("employees", company.Employees).
		Set("registered", company.Registered).
		Set("type", company.Type).
		Where(sq.Eq{"id": companyID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := b.ToSql()
	if err != nil {
		return errors.Wrap(err, "build query")
	}

	tag, err := s.conn.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "exec")
	}
	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}

func (s *storage) Get(ctx context.Context, companyID uuid.UUID) (model.Company, error) {
	b := sq.Select(selectColumns...).
		From(table).
		Where(sq.Eq{"id": companyID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := b.ToSql()
	if err != nil {
		return model.Company{}, errors.Wrap(err, "build query")
	}

	var company model.Company
	err = pgxscan.Get(ctx, s.conn, &company, query, args...)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.Company{}, model.ErrNotFound
	}

	return company, errors.Wrap(err, "select")
}

func (s *storage) Delete(ctx context.Context, companyID uuid.UUID) error {
	b := sq.Delete(table).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": companyID})
	query, args, err := b.ToSql()
	if err != nil {
		return errors.Wrap(err, "build query")
	}

	tag, err := s.conn.Exec(ctx, query, args...)
	if err != nil {
		return errors.Wrap(err, "exec")
	}

	if tag.RowsAffected() == 0 {
		return model.ErrNotFound
	}

	return nil
}
