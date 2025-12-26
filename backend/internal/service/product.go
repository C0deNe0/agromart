package service

import (
	"context"
	"fmt"

	"github.com/C0deNe0/agromart/internal/model"
	"github.com/C0deNe0/agromart/internal/model/product"
	"github.com/C0deNe0/agromart/internal/repository"
	"github.com/google/uuid"
)

type ProductCreateInput struct {
	CompanyID    uuid.UUID
	Name         string
	Description  *string
	CategoryID   *uuid.UUID
	Unit         *string
	Origin       *string
	PriceDisplay *string
}
type ProductUpdateInput struct {
	ID           uuid.UUID
	Name         *string
	Description  *string
	CategoryID   *uuid.UUID
	Unit         *string
	Origin       *string
	PriceDisplay *string
	IsActive     *bool
}

type ProductService struct {
	productRepo *repository.ProductRepository
	companyRepo *repository.CompanyRepository
}

func NewProductService(productRepo *repository.ProductRepository, companyRepo *repository.CompanyRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		companyRepo: companyRepo,
	}
}

func (s *ProductService) Create(ctx context.Context, userID uuid.UUID, input ProductCreateInput) (*product.Product, error) {

	//check if comp exists
	company, err := s.companyRepo.GetByID(ctx, input.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get company by id: %w", err)
	}

	if company.OwnerID != userID {
		return nil, fmt.Errorf("not autorized to create product for this company")
	}

	if !company.IsApproved {
		return nil, fmt.Errorf("company is not approved")
	}

	p := &product.Product{
		CompanyID:    input.CompanyID,
		Name:         input.Name,
		Description:  input.Description,
		CategoryID:   input.CategoryID,
		Unit:         input.Unit,
		Origin:       input.Origin,
		PriceDisplay: input.PriceDisplay,
		IsActive:     true,
	}

	return s.productRepo.Create(ctx, p)
}

func (s *ProductService) List(ctx context.Context, filter repository.ProductFilter) (*model.PaginatedResponse[product.Product], error) {
	return s.productRepo.List(ctx, filter)
}

func (s *ProductService) Update(ctx context.Context, userID uuid.UUID, p ProductUpdateInput) (*product.Product, error) {
	existing, err := s.productRepo.GetByID(ctx, p.ID)
	if err != nil {
		return nil, fmt.Errorf("couldn't find product: %w", err)
	}

	company, err := s.companyRepo.GetByID(ctx, existing.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("company doesn't exist")
	}

	if company.OwnerID != userID {
		return nil, fmt.Errorf("not authorized to update product")
	}

	if p.Name != nil {
		existing.Name = *p.Name
	}
	if p.Description != nil {
		existing.Description = p.Description
	}
	if p.CategoryID != nil {
		existing.CategoryID = p.CategoryID
	}
	if p.Unit != nil {
		existing.Unit = p.Unit
	}
	if p.Origin != nil {
		existing.Origin = p.Origin
	}
	if p.PriceDisplay != nil {
		existing.PriceDisplay = p.PriceDisplay
	}
	if p.IsActive != nil {
		existing.IsActive = *p.IsActive
	}

	return s.productRepo.Update(ctx, existing)
}

func (s *ProductService) GetByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	p, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}
	return p, nil
}

func (s *ProductService) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	existing, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("couldn't find product: %w", err)
	}

	company, err := s.companyRepo.GetByID(ctx, existing.CompanyID)
	if err != nil {
		return fmt.Errorf("company doesn't exist")
	}

	if company.OwnerID != userID {
		return fmt.Errorf("not authorized to delete product")
	}

	return s.productRepo.Delete(ctx, id)
}
