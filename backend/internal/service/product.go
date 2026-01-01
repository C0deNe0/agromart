package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/C0deNe0/agromart/internal/model"
	"github.com/C0deNe0/agromart/internal/model/product"
	"github.com/C0deNe0/agromart/internal/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

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

func (s *ProductService) Create(ctx context.Context, userID uuid.UUID, input product.ProductCreateInput) (*product.Product, error) {

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
		CompanyID:   input.CompanyID,
		Name:        input.Name,
		Description: input.Description,
		// CategoryID:  input.CategoryID,
		Unit:     *input.Unit,
		Origin:   input.Origin,
		Price:    input.Price,
		IsActive: true,
	}

	return s.productRepo.Create(ctx, p)
}

func (s *ProductService) List(ctx context.Context, filter repository.ProductFilter) (*model.PaginatedResponse[product.Product], error) {
	return s.productRepo.List(ctx, filter)
}

func (s *ProductService) Update(ctx context.Context, userID uuid.UUID, p product.ProductUpdateInput) (*product.Product, error) {
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
	// if p.CategoryID != nil {
	// 	existing.CategoryID = p.CategoryID
	// }
	if p.Unit != nil {
		existing.Unit = *p.Unit
	}
	if p.Origin != nil {
		existing.Origin = p.Origin
	}
	if p.Price != decimal.Zero {
		existing.Price = p.Price
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

func (s *ProductService) AuthorizeProductMutation(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error {

	// 1️⃣ Product must exist
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return errors.New("product not found")
	}

	// 2️⃣ Product must be active
	if !product.IsActive {
		return errors.New("product is inactive")
	}

	// 3️⃣ Company must exist
	company, err := s.companyRepo.GetByID(ctx, product.CompanyID)
	if err != nil {
		return errors.New("company not found")
	}

	// 4️⃣ Company must be active
	if !company.IsActive {
		return errors.New("company is inactive")
	}

	// 5️⃣ Company must be approved
	if !company.IsApproved {
		return errors.New("company is not approved")
	}

	// 6️⃣ Ownership check (USER only for now)
	if company.OwnerID != userID {
		return errors.New("you do not own this product")
	}

	return nil
}

// func (s *ProductService) GenerateImageUploadURL(
// 	ctx context.Context,
// 	productID uuid.UUID,
// 	userID uuid.UUID,
// 	input product.ProductImageUploadInput,
// ) (string, error) {

// }
