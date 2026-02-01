package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/C0deNe0/agromart/internal/lib/aws"
	"github.com/C0deNe0/agromart/internal/model"
	"github.com/C0deNe0/agromart/internal/model/company"
	"github.com/C0deNe0/agromart/internal/model/product"
	"github.com/C0deNe0/agromart/internal/repository"
	"github.com/C0deNe0/agromart/internal/repository/productRepo"
	"github.com/google/uuid"
)

type ProductService struct {
	productRepo        *productRepo.ProductRepository
	productImageRepo   *productRepo.ProductImageRepository
	productVariantRepo *productRepo.ProductVariantRepository
	companyRepo        *repository.CompanyRepository
	categoryRepo       *repository.CategoryRepository
	S3Service          *aws.S3Service
}

func NewProductService(
	productRepo *productRepo.ProductRepository,
	productImageRepo *productRepo.ProductImageRepository,
	productVariantRepo *productRepo.ProductVariantRepository,
	companyRepo *repository.CompanyRepository,

	s3 *aws.S3Service,
) *ProductService {
	return &ProductService{
		productRepo:        productRepo,
		productImageRepo:   productImageRepo,
		productVariantRepo: productVariantRepo,
		companyRepo:        companyRepo,
		S3Service:          s3,
	}
}

func (s *ProductService) Create(ctx context.Context, userID uuid.UUID, req *product.CreateProductRequest) (*product.ProductResponse, error) {

	approvedCompany, err := s.companyRepo.GetApprovedCompanyByOwner(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check company approval: %w", err)
	}

	if approvedCompany == nil {
		return nil, errors.New("you must have an approved company before creating products. Please create a company and wait for admin approval")
	}

	if req.CompanyID != approvedCompany.ID {
		return nil, errors.New("you can only create products for your own approved company")
	}

	if req.CategoryID != nil {
		_, err := s.categoryRepo.GetByID(ctx, *req.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid category: %w", err)
		}
	}

	p := &product.Product{
		CompanyID:      req.CompanyID,
		CategoryID:     req.CategoryID,
		Name:           req.Name,
		Description:    req.Description,
		Unit:           req.Unit,
		Origin:         req.Origin,
		BasePrice:      req.BasePrice,
		ApprovalStatus: company.ApprovalStatusPending,
		IsActive:       true,
	}

	created, err := s.productRepo.Create(ctx, p)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	//creating the variant
	variants := make([]product.ProductVariant, 0, len(req.Variants))
	for _, variantInput := range req.Variants {
		variant := &product.ProductVariant{
			ProductID:         created.ID,
			Label:             variantInput.Label,
			QuantityValue:     variantInput.QuantityValue,
			QuantityUnit:      variantInput.QuantityUnit,
			Price:             variantInput.Price,
			StockQuantity:     variantInput.StockQuantity,
			LowStockThreshold: variantInput.LowStockThreshold,
			IsAvailable:       true,
		}

		createdVariant, err := s.productVariantRepo.Create(ctx, variant)
		if err != nil {
			return nil, fmt.Errorf("failed to create variant: %w", err)
		}
		variants = append(variants, *createdVariant)
	}

	return product.ToProductResponse(created, []product.ProductImage{}, variants), nil
}

func (s *ProductService) List(ctx context.Context, userID *uuid.UUID, filter productRepo.ProductFilter) (*model.PaginatedResponse[product.ProductResponse], error) {
	if userID == nil && filter.ApprovalStatus == nil {
		approvalStatus := company.ApprovalStatusApproved
		filter.ApprovalStatus = &approvalStatus
	}

	//get the products
	products, err := s.productRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	if len(products.Data) == 0 {
		return &model.PaginatedResponse[product.ProductResponse]{
			Data:       []product.ProductResponse{},
			Page:       products.Page,
			Limit:      products.Limit,
			Total:      0,
			TotalPages: 0,
		}, nil
	}

	productIDs := make([]uuid.UUID, len(products.Data))
	for i, p := range products.Data {
		productIDs[i] = p.ID
	}
	//get the images
	images, err := s.productImageRepo.ListByProductIDs(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to list images: %w", err)
	}

	//get the variants
	variants, err := s.productVariantRepo.ListByProductIDs(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to list variants: %w", err)
	}

	return product.MapProductPage(products, images, variants), nil
}

func (s *ProductService) Update(ctx context.Context, userID uuid.UUID, productID uuid.UUID, updates *product.UpdateProductRequest) (*product.Product, error) {
	existing, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Check ownership
	comp, err := s.companyRepo.GetByID(ctx, existing.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get company: %w", err)
	}

	if comp.OwnerID != userID {
		return nil, errors.New("not authorized to update this product")
	}

	// ✅ CRITICAL: Can only update PENDING or REJECTED products
	if !existing.CanBeModified() {
		return nil, fmt.Errorf("cannot modify product with status: %s. Only PENDING or REJECTED products can be modified", existing.ApprovalStatus)
	}

	// Apply updates
	if updates.CategoryID != nil {
		if *updates.CategoryID != uuid.Nil {
			_, err := s.categoryRepo.GetByID(ctx, *updates.CategoryID)
			if err != nil {
				return nil, fmt.Errorf("invalid category: %w", err)
			}
		}
		existing.CategoryID = updates.CategoryID
	}

	if updates.Name != nil {
		existing.Name = *updates.Name
	}
	if updates.Description != nil {
		existing.Description = updates.Description
	}
	if updates.Unit != nil {
		existing.Unit = *updates.Unit
	}
	if updates.Origin != nil {
		existing.Origin = updates.Origin
	}
	if updates.BasePrice != nil {
		existing.BasePrice = *updates.BasePrice
	}

	return s.productRepo.Update(ctx, existing)
}

func (s *ProductService) GetByID(ctx context.Context, id uuid.UUID, userID *uuid.UUID) (*product.ProductResponse, error) {
	p, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}
	if p.ApprovalStatus != company.ApprovalStatusApproved {
		if userID == nil {
			return nil, fmt.Errorf("product is not approved")
		}
		comp, err := s.companyRepo.GetByID(ctx, p.CompanyID)
		if err != nil {
			return nil, fmt.Errorf("company not found")
		}
		if comp.OwnerID != *userID {
			return nil, fmt.Errorf("not authorized to get product")
		}
	}
	//getting the images
	images, err := s.productImageRepo.ListByProductID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get images: %w", err)
	}

	//getting the variants
	variants, err := s.productVariantRepo.ListByProductID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get the variants: %w", err)
	}
	return product.ToProductResponse(p, images, variants), nil
}

func (s *ProductService) Delete(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error {
	existing, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	comp, err := s.companyRepo.GetByID(ctx, existing.CompanyID)
	if err != nil {
		return fmt.Errorf("failed to get company: %w", err)
	}

	if comp.OwnerID != userID {
		return errors.New("not authorized to delete this product")
	}

	// Can only delete PENDING or REJECTED products
	if !existing.CanBeModified() {
		return fmt.Errorf("cannot delete product with status: %s", existing.ApprovalStatus)
	}

	return s.productRepo.Delete(ctx, productID)
}

func (s *ProductService) Resubmit(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error {
	existing, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	comp, err := s.companyRepo.GetByID(ctx, existing.CompanyID)
	if err != nil {
		return fmt.Errorf("failed to get company: %w", err)
	}

	if comp.OwnerID != userID {
		return errors.New("not authorized to resubmit this product")
	}

	if !existing.IsRejected() {
		return errors.New("only rejected products can be resubmitted")
	}

	return s.productRepo.Resubmit(ctx, productID, userID)
}

// ADMINT APPROVE PRODUCTS
func (s *ProductService) Approve(ctx context.Context, productID, adminID uuid.UUID, notes *string) error {
	existing, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	if !existing.IsPending() {
		return fmt.Errorf("only pending products can be approved. Current status: %s", existing.ApprovalStatus)
	}

	return s.productRepo.Approve(ctx, productID, adminID, notes)
}
func (s *ProductService) Reject(ctx context.Context, productID, adminID uuid.UUID, reason string, notes *string) error {
	existing, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	if !existing.IsPending() {
		return fmt.Errorf("only pending products can be rejected. Current status: %s", existing.ApprovalStatus)
	}

	return s.productRepo.Reject(ctx, productID, adminID, reason, notes)
}

func (s *ProductService) GetApprovalHistory(ctx context.Context, productID uuid.UUID) ([]product.ProductApprovalHistory, error) {
	return s.productRepo.GetApprovalHistory(ctx, productID)
}

func (s *ProductService) CountPendingApprovals(ctx context.Context) (int, error) {
	return s.productRepo.CountPendingApprovals(ctx)
}

//PRODUCT IMAGES

func (s *ProductService) GenerateImageUploadURL(ctx context.Context, userID uuid.UUID, req *product.GenerateImageUploadURLRequest) (*product.ImageUploadURLResponse, error) {
	p, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	comp, err := s.companyRepo.GetByID(ctx, p.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("company not found: %w", err)
	}
	if comp.OwnerID != userID {
		return nil, errors.New("not authorized to upload images for this product")
	}

	imageID := uuid.New()
	s3Key := fmt.Sprintf("products/$s/images/%s", req.ProductID.String(), imageID.String())

	//presin url
	uploadURL := s.S3Service.GetPublicURL(s3Key)

	img := &product.ProductImage{
		ID:        imageID,
		ProductID: req.ProductID,
		ImageURL:  uploadURL,
		S3Key:     s3Key,
		IsPrimary: *req.IsPrimary,
	}
	_, err = s.productImageRepo.Create(ctx, img)
	if err != nil {
		return nil, fmt.Errorf("failed to created image record: %w", err)
	}

	return &product.ImageUploadURLResponse{
		UploadURL: uploadURL,
		ImageID:   imageID,
		S3Key:     s3Key,
		ExpiresIn: 900,
	}, nil
}

func (s *ProductService) DeleteImage(ctx context.Context, userID uuid.UUID, productID, imageID uuid.UUID) error {
	// Check product ownership
	p, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	comp, err := s.companyRepo.GetByID(ctx, p.CompanyID)
	if err != nil {
		return fmt.Errorf("failed to get company: %w", err)
	}

	if comp.OwnerID != userID {
		return errors.New("not authorized to delete images for this product")
	}

	// Get image
	img, err := s.productImageRepo.GetByID(ctx, imageID)
	if err != nil {
		return fmt.Errorf("image not found: %w", err)
	}

	if img.ProductID != productID {
		return errors.New("image does not belong to this product")
	}

	// Delete from S3
	err = s.S3Service.DeleteObject(ctx, img.S3Key)
	if err != nil {
		return fmt.Errorf("failed to delete from S3: %w", err)
	}

	// Delete from database
	err = s.productImageRepo.Delete(ctx, imageID)
	if err != nil {
		return fmt.Errorf("failed to delete image record: %w", err)
	}

	return nil
}

func (s *ProductService) SetPrimaryImage(ctx context.Context, userID uuid.UUID, productID, imageID uuid.UUID) error {
	// Check product ownership
	p, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	comp, err := s.companyRepo.GetByID(ctx, p.CompanyID)
	if err != nil {
		return fmt.Errorf("failed to get company: %w", err)
	}

	if comp.OwnerID != userID {
		return errors.New("not authorized to manage images for this product")
	}

	return s.productImageRepo.SetPrimary(ctx, productID, imageID)
}

//VARIANT MANAGEMENT

func (s *ProductService) CreateVariant(ctx context.Context, userID uuid.UUID, req *product.CreateVariantRequest) (*product.ProductVariant, error) {
	// Check product ownership
	p, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	comp, err := s.companyRepo.GetByID(ctx, p.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get company: %w", err)
	}

	if comp.OwnerID != userID {
		return nil, errors.New("not authorized to create variants for this product")
	}

	// Can only add variants to PENDING or REJECTED products
	if !p.CanBeModified() {
		return nil, fmt.Errorf("cannot add variants to product with status: %s", p.ApprovalStatus)
	}

	variant := &product.ProductVariant{
		ProductID: req.ProductID,

		Label:             req.Label,
		QuantityValue:     req.QuantityValue,
		QuantityUnit:      req.QuantityUnit,
		Price:             req.Price,
		StockQuantity:     req.StockQuantity,
		LowStockThreshold: req.LowStockThreshold,
		IsAvailable:       true,
	}

	return s.productVariantRepo.Create(ctx, variant)
}

func (s *ProductService) UpdateVariant(ctx context.Context, userID uuid.UUID, productID, variantID uuid.UUID, updates *product.UpdateVariantRequest) (*product.ProductVariant, error) {
	// Check product ownership
	p, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	comp, err := s.companyRepo.GetByID(ctx, p.CompanyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get company: %w", err)
	}

	if comp.OwnerID != userID {
		return nil, errors.New("not authorized to update variants for this product")
	}

	// Get existing variant
	existing, err := s.productVariantRepo.GetByID(ctx, variantID)
	if err != nil {
		return nil, fmt.Errorf("variant not found: %w", err)
	}

	if existing.ProductID != productID {
		return nil, errors.New("variant does not belong to this product")
	}

	// Apply updates

	if updates.Label != nil {
		existing.Label = *updates.Label
	}
	if updates.QuantityValue != nil {
		existing.QuantityValue = *updates.QuantityValue
	}
	if updates.QuantityUnit != nil {
		existing.QuantityUnit = *updates.QuantityUnit
	}
	if updates.Price != nil {
		existing.Price = *updates.Price
	}
	if updates.StockQuantity != nil {
		existing.StockQuantity = updates.StockQuantity
	}
	if updates.LowStockThreshold != nil {
		existing.LowStockThreshold = updates.LowStockThreshold
	}
	if updates.IsAvailable != nil {
		existing.IsAvailable = *updates.IsAvailable
	}

	return s.productVariantRepo.Update(ctx, existing)
}

func (s *ProductService) DeleteVariant(ctx context.Context, userID uuid.UUID, productID, variantID uuid.UUID) error {
	// Check product ownership
	p, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	comp, err := s.companyRepo.GetByID(ctx, p.CompanyID)
	if err != nil {
		return fmt.Errorf("failed to get company: %w", err)
	}

	if comp.OwnerID != userID {
		return errors.New("not authorized to delete variants for this product")
	}

	// Get existing variant
	existing, err := s.productVariantRepo.GetByID(ctx, variantID)
	if err != nil {
		return fmt.Errorf("variant not found: %w", err)
	}

	if existing.ProductID != productID {
		return errors.New("variant does not belong to this product")
	}

	// Check if this is the last variant
	variants, err := s.productVariantRepo.ListByProductID(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to check variants: %w", err)
	}

	if len(variants) <= 1 {
		return errors.New("cannot delete the last variant. Products must have at least one variant")
	}

	return s.productVariantRepo.Delete(ctx, variantID)
}
