package service

import (
	"context"
	"errors"
	"time"

	"github.com/C0deNe0/agromart/internal/model/company"
	"github.com/C0deNe0/agromart/internal/repository"
	"github.com/google/uuid"
)

type CompanyService struct {
	companyRepo *repository.CompanyRepository
}

func NewCompanyService(companyRepo *repository.CompanyRepository) *CompanyService {
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

func (s *CompanyService) DeleteCompany(ctx context.Context, adminID uuid.UUID, companyID uuid.UUID) error {

	c, err :=s.companyRepo.GetByID(ctx, companyID)
	if err != nil {
		return err
	}

	if !c.IsActive {
		return errors.New("company already deleted")
	}
	

	return s.companyRepo.SoftDelete(ctx, companyID)
}
func (s *CompanyService) ApproveCompany(ctx context.Context, adminID uuid.UUID, companyID uuid.UUID) error {
	c, err := s.companyRepo.GetByID(ctx, companyID)
	if err != nil {
		return err
	}

	now := time.Now()
	if c.IsApproved {
		return errors.New("company already approved")
	}
	c.IsApproved = true
	c.ApprovedBy = &adminID
	c.ApprovedAt = &now

	return s.companyRepo.Approve(ctx, adminID, companyID)
}

func (s *CompanyService) RejectCompany(ctx context.Context, companyID uuid.UUID) error {
	c, err := s.companyRepo.GetByID(ctx, companyID)
	if err != nil {
		return err
	}
	if c.IsApproved {
		return errors.New("company already approved")
	}
	return s.companyRepo.Reject(ctx, companyID)
}
