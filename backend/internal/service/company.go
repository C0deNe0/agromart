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

func (s *CompanyService) Create(ctx context.Context, userID uuid.UUID, input company.CreateCompanyInput) (*company.Company, error) {
	c := &company.Company{
		OwnerID:       userID,
		Name:          input.Name,
		Description:   input.Description,
		LogoURL:       &input.LogoURL,
		BusinessEmail: &input.BusinessEmail,
		BusinessPhone: &input.BusinessPhone,
		City:          input.City,
		State:         input.State,
		Pincode:       input.Pincode,
		IsApproved:    false,
		IsActive:      true,
	}

	return s.companyRepo.Create(ctx, c)
}
func (s *CompanyService) ListByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]company.Company, error) {
	return s.companyRepo.ListByOwnerID(ctx, ownerID)
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
