package company

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"xm/internal/model"
)

type (
	storage interface {
		Insert(ctx context.Context, company *model.Company) (uuid.UUID, error)
		Update(ctx context.Context, companyID uuid.UUID, company *model.Company) error
		Get(ctx context.Context, companyID uuid.UUID) (model.Company, error)
		Delete(ctx context.Context, companyID uuid.UUID) error
	}

	service struct {
		logger  *zap.Logger
		storage storage
	}
)

func NewService(
	logger *zap.Logger,
	storage storage,
) *service {
	return &service{
		logger:  logger,
		storage: storage,
	}
}

func (s *service) New(ctx context.Context, company *model.Company) (uuid.UUID, error) {
	uuid, err := s.storage.Insert(ctx, company)
	return uuid, errors.Wrap(err, "insert")
}

func (s *service) Delete(ctx context.Context, companyID uuid.UUID) error {
	err := s.storage.Delete(ctx, companyID)
	return errors.Wrap(err, "delete")
}

func (s *service) Update(ctx context.Context, companyID uuid.UUID, company *model.Company) error {
	err := s.storage.Update(ctx, companyID, company)
	return errors.Wrap(err, "update")
}

func (s *service) Get(ctx context.Context, companyID uuid.UUID) (model.Company, error) {
	company, err := s.storage.Get(ctx, companyID)
	if err != nil {
		return model.Company{}, errors.Wrap(err, "get")
	}

	return company, nil
}
