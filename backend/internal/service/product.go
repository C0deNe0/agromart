package service

import (
	"context"
	"fmt"

	"github.com/C0deNe0/agromart/internal/model"
	"github.com/C0deNe0/agromart/internal/model/product"
	"github.com/C0deNe0/agromart/internal/repository"
	"github.com/google/uuid"
)

type ProductService struct {
	productRepo repository.ProductRepositoryImp
	companyRepo repository.CompanyRepositoryImp
}

func NewProductService(productRepo repository.ProductRepositoryImp, companyRepo repository.CompanyRepositoryImp) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		companyRepo: companyRepo,
	}
}

func (s *ProductService) Create(ctx context.Context, userID uuid.UUID, input product.CreateProductRequest) (*product.Product, error) {

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
		CompanyID:    *input.CompanyID,
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

func (s *ProductService) Update(ctx context.Context, userID uuid.UUID, p *product.Product) (*product.Product, error) {
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

	return s.productRepo.Update(ctx, p)
}
 


