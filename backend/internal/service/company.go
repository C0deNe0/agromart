package service

import (
	"context"

	"github.com/C0deNe0/agromart/internal/model/company"
	"github.com/C0deNe0/agromart/internal/repository"
	"github.com/google/uuid"
)

type CompanyService struct {
	companyRepo repository.CompanyRepositoryImp
}

func NewCompanyService(companyRepo repository.CompanyRepositoryImp) *CompanyService {
	return &CompanyService{
		companyRepo: companyRepo,
	}
}

func (s *CompanyService) ListPending(ctx context.Context) ([]company.Company, error) {
	return s.companyRepo.ListPending(ctx)
}

func (s *CompanyService) Approve(ctx context.Context, adminID uuid.UUID, companyID uuid.UUID) error {
	return s.companyRepo.Approve(ctx, adminID, companyID)
}

func (s *CompanyService) Reject(ctx context.Context, adminID uuid.UUID, companyID uuid.UUID) error {
	return s.companyRepo.Reject(ctx, adminID, companyID)
}

