package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/C0deNe0/agromart/internal/model"
	"github.com/C0deNe0/agromart/internal/model/company"
	"github.com/C0deNe0/agromart/internal/repository"
	"github.com/google/uuid"
)

type CompanyService struct {
	companyRepo         *repository.CompanyRepository
	companyFollowerRepo *repository.CompanyFollowerRepository
}

func NewCompanyService(companyRepo *repository.CompanyRepository, companyFollowerRepo *repository.CompanyFollowerRepository) *CompanyService {
	return &CompanyService{
		companyRepo:         companyRepo,
		companyFollowerRepo: companyFollowerRepo,
	}
}

func (s *CompanyService) Create(ctx context.Context, userID uuid.UUID, c company.Company) (*company.Company, error) {

	existing, err := s.companyRepo.GetByOwnerAndName(ctx, userID, c.Name)
	if err != nil {
		return nil, fmt.Errorf("fialed to check existing company:%w", err)
	}

	if existing != nil {
		return nil, errors.New("company already present by your name")
	}

	c.OwnerID = userID
	c.IsActive = true
	c.ApprovalStatus = company.ApprovalStatusPending

	if c.ProductVisibility == "" {
		c.ProductVisibility = company.ProductVisibilityPublic
	}

	created, err := s.companyRepo.Create(ctx, &c)
	if err != nil {
		return nil, fmt.Errorf("failed to create company: %w", err)
	}
	return created, nil

}
func (s *CompanyService) GetByID(ctx context.Context, id uuid.UUID, userID *uuid.UUID) (*company.CompanyResponse, error) {
	c, err := s.companyRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}

	var isFollowing *bool
	if userID != nil {
		following, err := s.companyFollowerRepo.IsFollowing(ctx, id, *userID)
		if err != nil {
			return nil, fmt.Errorf("failed to check follow status: %w", err)
		}
		isFollowing = &following
	}
	// hideSensitve := true
	// if userID != nil && c.OwnerID == *userID {
	// 	hideSensitve = false
	// }

	return company.ToCompanyResponse(c, isFollowing), nil
}
func (s *CompanyService) List(ctx context.Context, userID *uuid.UUID, filter repository.CompanyFilter) (*model.PaginatedResponse[company.CompanyResponse], error) {
	result, err := s.companyRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list companies:%w", err)
	}

	followStatusMap := make(map[uuid.UUID]bool)
	if userID != nil {
		companyIDs := make([]uuid.UUID, len(result.Data))
		for i, company := range result.Data {
			companyIDs[i] = company.ID
		}

		followStatusMap, err = s.companyFollowerRepo.GetFollowStatusBatch(ctx, companyIDs, *userID)
		if err != nil {
			return nil, fmt.Errorf("failed to get follow statuses: %w", err)
		}
	}

	return company.MapCompanyPage(result, followStatusMap), nil

}

func (s *CompanyService) Update(ctx context.Context, userID uuid.UUID, companyID uuid.UUID, updates *company.UpdateCompanyRequest) (*company.Company, error) {

	existing, err := s.companyRepo.GetByID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}

	if existing.OwnerID != userID {
		return nil, errors.New("unauthorized to update the company")
	}

	if !existing.CanBeModified() {
		return nil, fmt.Errorf("cannot modify company with status: %s. Only PENDING and REJECTED companies can be modified", existing.ApprovalStatus)
	}

	if updates.Name != nil {
		duplicate, err := s.companyRepo.GetByOwnerAndName(ctx, userID, *updates.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to check duplicate name: %w", err)
		}

		if duplicate != nil && duplicate.ID != companyID {
			return nil, errors.New("company name already exists")
		}
		existing.Name = *updates.Name
	}

	if updates.Description != nil {
		existing.Description = updates.Description
	}
	if updates.LogoURL != nil {
		existing.LogoURL = updates.LogoURL
	}
	if updates.BusinessEmail != nil {
		existing.BusinessEmail = updates.BusinessEmail
	}
	if updates.BusinessPhone != nil {
		existing.BusinessPhone = updates.BusinessPhone
	}
	if updates.City != nil {
		existing.City = updates.City
	}
	if updates.State != nil {
		existing.State = updates.State
	}
	if updates.Pincode != nil {
		existing.Pincode = updates.Pincode
	}
	if updates.GSTNumber != nil {
		existing.GSTNumber = updates.GSTNumber
	}
	if updates.PANNumber != nil {
		existing.PANNumber = updates.PANNumber
	}
	if updates.ProductVisibility != nil {
		existing.ProductVisibility = *updates.ProductVisibility
	}
	if updates.IsActive != nil {
		existing.IsActive = *updates.IsActive
	}

	return s.companyRepo.Update(ctx, existing)
}

func (s *CompanyService) Delete(ctx context.Context, userID uuid.UUID, companyID uuid.UUID) error {
	existing, err := s.companyRepo.GetByID(ctx, companyID)
	if err != nil {
		return err
	}
	if existing.OwnerID != userID {
		return errors.New("unauthorized to delete the company")
	}
	if !existing.CanBeModified() {
		return fmt.Errorf("cannot delete company with status: %s", existing.ApprovalStatus)
	}
	return s.companyRepo.Delete(ctx, companyID)
}

func (s *CompanyService) Resubmit(ctx context.Context, userID uuid.UUID, companyID uuid.UUID) error {
	existing, err := s.companyRepo.GetByID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("company not found: %w", err)
	}

	if existing.OwnerID != userID {
		return errors.New("not authorized to resubmit this company")
	}

	if !existing.IsRejected() {
		return errors.New("only rejected companies can be resubmitted")
	}

	return s.companyRepo.Resubmit(ctx, companyID, userID)
}

func (s *CompanyService) Approve(ctx context.Context, companyID uuid.UUID, adminID uuid.UUID, notes *string) error {
	existing, err := s.companyRepo.GetByID(ctx, companyID)
	if err != nil {
		return err
	}

	if !existing.IsPending() {
		return fmt.Errorf("only pending companies can be approved. Current status: %s", existing.ApprovalStatus)
	}

	return s.companyRepo.Approve(ctx, companyID, adminID, notes)
}

func (s *CompanyService) Reject(ctx context.Context, companyID, adminID uuid.UUID, reason string, notes *string) error {
	existing, err := s.companyRepo.GetByID(ctx, companyID)
	if err != nil {
		return fmt.Errorf("company not found :%w", err)

	}

	if !existing.IsPending() {
		return fmt.Errorf("only pending companies can be rejected. Current status: %s", existing.ApprovalStatus)
	}

	return s.companyRepo.Reject(ctx, companyID, adminID, reason, notes)
}

func (s *CompanyService) GetApprovalHistory(ctx context.Context, companyID uuid.UUID) ([]company.CompanyApprovalHistory, error) {
	return s.companyRepo.GetApprovalHistory(ctx, companyID)
}

func (s *CompanyService) CountPendingApprovals(ctx context.Context) (int, error) {
	return s.companyRepo.CountPendingApprovals(ctx)
}

//follow methids

func (s *CompanyService) Follow(ctx context.Context, companID, userID uuid.UUID) error {
	comp, err := s.companyRepo.GetByID(ctx, companID)
	if err != nil {
		return fmt.Errorf("company not found: %w", err)
	}

	if !comp.CanBeFollowed() {
		return fmt.Errorf("cannot follow compnay with status:%s", comp.ApprovalStatus)
	}
	if !comp.IsActive {
		return errors.New("connot follow inactive company")
	}

	if comp.OwnerID == userID {
		return errors.New("connot follow own company")
	}

	_, err = s.companyFollowerRepo.Follow(ctx, companID, userID)
	if err != nil {
		return fmt.Errorf("failed to follow company: %w", err)
	}

	return nil
}

func (s *CompanyService) Unfollow(ctx context.Context, companID, userID uuid.UUID) error {
	return s.companyFollowerRepo.Unfollow(ctx, companID, userID)
}

func (s *CompanyService) GetFollowStatus(ctx context.Context, companyID uuid.UUID, userID uuid.UUID) (*company.FollowStatusResponse, error) {
	follower, err := s.companyFollowerRepo.GetFollowStatus(ctx, companyID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get follow status : %w", err)

	}

	resp := &company.FollowStatusResponse{
		CompanyID:   companyID,
		IsFollowing: follower != nil,
	}

	if follower != nil {
		resp.FollowedAt = &follower.FollowedAt
	}
	return resp, nil
}

func (s *CompanyService) ListFollowers(ctx context.Context, companyID uuid.UUID, page, limit int) (*model.PaginatedResponse[company.CompanyFollowerResponse], error) {
	return s.companyFollowerRepo.ListFollowers(ctx, companyID, page, limit)
}

func (s *CompanyService) ListFollowedCompanies(ctx context.Context, userID uuid.UUID, page, limit int) (*model.PaginatedResponse[company.CompanyResponse], error) {
	res, err := s.companyFollowerRepo.ListFollowedCompanies(ctx, userID, page, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list followed companies: %w", err)
	}

	followStatusMap := make(map[uuid.UUID]bool)
	for _, c := range res.Data {
		followStatusMap[c.ID] = true
	}

	return company.MapCompanyPage(res, followStatusMap), nil
}

func (s *CompanyService) CanViewProducts(ctx context.Context, companyID uuid.UUID, userID *uuid.UUID) (bool, error) {
	return s.companyFollowerRepo.CanViewProducts(ctx, companyID, userID)
}

func (s *CompanyService) UserHasApprovedCompany(ctx context.Context, userID uuid.UUID) (bool, error) {
	approvedCompany, err := s.companyRepo.GetApprovedCompanyByOwner(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to check approved company: %w", err)
	}

	return approvedCompany != nil, nil
}
func (s *CompanyService) GetUserApprovedCompany(ctx context.Context, userID uuid.UUID) (*company.Company, error) {
	approvedCompany, err := s.companyRepo.GetApprovedCompanyByOwner(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get approved company: %w", err)
	}

	if approvedCompany == nil {
		return nil, errors.New("you don't have an approved company. Please create and get approval first")
	}

	return approvedCompany, nil
}
