package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/C0deNe0/agromart/internal/model"
	"github.com/C0deNe0/agromart/internal/model/product"
	"github.com/C0deNe0/agromart/internal/repository"
	"github.com/google/uuid"
)

type ProductService struct {
	productRepo  *repository.ProductRepository
	companyRepo  *repository.CompanyRepository
	categoryRepo *repository.CategoryRepository
}

func NewProductService(productRepo *repository.ProductRepository, companyRepo *repository.CompanyRepository, categoryRepo *repository.CategoryRepository) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		companyRepo:  companyRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *ProductService) Create(ctx context.Context, userID uuid.UUID, p *product.Product) (*product.Product, error) {
	
	approvedCompany, err := s.companyRepo.GetApprovedCompanyByOwner(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check company approval: %w", err)
	}

	if approvedCompany == nil {
		return nil, errors.New("you must have an approved company before creating products. Please create a company and wait for admin approval")
	}

	if p.CompanyID != approvedCompany.ID {
		return nil, errors.New("you can only create products for your own approved company")
	}

	company, err := s.companyRepo.GetByID(ctx, p.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get company: %w", err)
	}

	if company.OwnerID != userID {
		return nil, errors.New("not authorized to create product for this company")
	}

	if !company.IsApproved() {
		return nil, errors.New("company must be approved before creating products")
	}

	if !company.IsActive {
		return nil, errors.New("company is not active")
	}

	return s.productRepo.Create(ctx, p)
}

func (s *ProductService) List(ctx context.Context, filter repository.ProductFilter) (*model.PaginatedResponse[product.Product], error) {
	return s.productRepo.List(ctx, filter)
}

func (s *ProductService) ListWithCategory(ctx context.Context, filter repository.ProductFilter) (*model.PaginatedResponse[product.ProductWithCategoryResponse], error) {
	return s.productRepo.ListWithCategory(ctx, filter)
}

func (s *ProductService) Update(ctx context.Context, userID uuid.UUID, productID uuid.UUID, p *product.UpdateProductRequest) (*product.Product, error) {
	existing, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	company, err := s.companyRepo.GetByID(ctx, existing.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("company doesn't exist")
	}

	if company.OwnerID != userID {
		return nil, fmt.Errorf("not authorized to update product")
	}

	if p.CategoryID != nil {
		_, err := s.categoryRepo.GetByID(ctx, *p.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid category:%w", err)
		}
		existing.CategoryID = p.CategoryID
	}

	if p.Name != nil {
		existing.Name = *p.Name
	}
	if p.Description != nil {
		existing.Description = p.Description
	}

	if p.Unit != nil {
		existing.Unit = *p.Unit
	}
	if p.Origin != nil {
		existing.Origin = p.Origin
	}
	if p.Price != nil {
		existing.Price = *p.Price
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
		return fmt.Errorf("product not found: %w", err)
	}

	company, err := s.companyRepo.GetByID(ctx, existing.CompanyID)
	if err != nil {
		return fmt.Errorf("company not found: %w", err)
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
	if !company.IsApproved() {
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
