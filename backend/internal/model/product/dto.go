package product

import (
	"errors"
	"time"

	"github.com/C0deNe0/agromart/internal/model"
	"github.com/C0deNe0/agromart/internal/model/company"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

//CREATE PRODUCT

type CreateProductRequest struct {
	CompanyID   uuid.UUID       `json:"companyId" validate:"required,uuid"`
	CategoryID  *uuid.UUID      `json:"categoryId,omitempty" validate:"omitempty,uuid"`
	Name        string          `json:"name" validate:"required,min=3,max=255"`
	Description *string         `json:"description,omitempty" validate:"omitempty,max=2000"`
	Unit        string          `json:"unit" validate:"required,max=50"`
	Origin      *string         `json:"origin,omitempty" validate:"omitempty,max=100"`
	BasePrice   decimal.Decimal `json:"basePrice" validate:"required,gt=0"`

	Variants []CreateVariantInput `json:"variants" validate:"required,min=1,dive"`
}

func (r *CreateProductRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

type CreateVariantInput struct {
	// SKU           *string         `json:"sku,omitempty" validate:"omitempty,max=50"`
	Label         string          `json:"label" validate:"required,min=2,max=100"`
	QuantityValue decimal.Decimal `json:"quantityValue" validate:"required,gt=0"`
	QuantityUnit  string          `json:"quantityUnit" validate:"required,max=20"`
	Price         decimal.Decimal `json:"price" validate:"required,gt=0"`

	StockQuantity     *int `json:"stockQuantity,omitempty" validate:"omitempty,gte=0"`
	LowStockThreshold *int `json:"lowStockThreshold,omitempty" validate:"omitempty,gte=0"`
}

// UPDATE PRODUCT

type UpdateProductRequest struct {
	ID          uuid.UUID        `param:"id" validate:"required,uuid"`
	CategoryID  *uuid.UUID       `json:"categoryId,omitempty" validate:"omitempty,uuid"`
	Name        *string          `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Description *string          `json:"description,omitempty" validate:"omitempty,max=2000"`
	Unit        *string          `json:"unit,omitempty" validate:"omitempty,max=50"`
	Origin      *string          `json:"origin,omitempty" validate:"omitempty,max=100"`
	BasePrice   *decimal.Decimal `json:"basePrice,omitempty" validate:"omitempty,gt=0"`
}

func (r *UpdateProductRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// LIST PRODUCT

// type ListProductsQuery struct {
// 	Page       *int       `query:"page" validate:"required,min=1"`
// 	Limit      *int       `query:"limit" validate:"required,min=1"`
// 	CompanyID  *uuid.UUID `query:"companyId" validate:"omitempty,uuid"`
// 	CategoryID *uuid.UUID `query:"categoryId" validate:"omitempty,uuid"`
// 	Search     *string    `query:"search" validate:"omitempty,min=1"`
// 	IsActive   *bool      `query:"isActive" validate:"omitempty"`
// }

type ListProductsQuery struct {
	Page           int                     `query:"page" validate:"min=1"`
	Limit          int                     `query:"limit" validate:"min=1,max=100"`
	CompanyID      *uuid.UUID              `query:"companyId" validate:"omitempty,uuid"`
	CategoryID     *uuid.UUID              `query:"categoryId" validate:"omitempty,uuid"`
	Search         *string                 `query:"search" validate:"omitempty,min=1"`
	ApprovalStatus *company.ApprovalStatus `query:"approvalStatus" validate:"omitempty"`
	IsActive       *bool                   `query:"isActive" validate:"omitempty"`
}

func (q *ListProductsQuery) Validate() error {

	if q.Page == 0 {
		q.Page = 1
	}

	if q.Limit == 0 {
		q.Limit = 10
	}

	if q.Page < 1 {
		return errors.New("page must be greater than 0")
	}

	if q.Limit < 1 {
		return errors.New("limit must be greater than 0")
	}

	validate := validator.New()

	if err := validate.Struct(q); err != nil {
		return err
	}

	return nil
}

type GetProductByIDRequest struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (r *GetProductByIDRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

type DeleteProductRequest struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (r *DeleteProductRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

//PRODUCT IMAGES

type GenerateImageUploadURLRequest struct {
	ProductID   uuid.UUID `param:"productId" validate:"required,uuid"`
	FileName    string    `json:"fileName" validate:"required"`
	ContentType string    `json:"contentType" validate:"required,oneof=image/jpeg image/png image/webp"`
	IsPrimary   *bool     `json:"isPrimary,omitempty" validate:"omitempty"`
}

func (r *GenerateImageUploadURLRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

type ImageUploadURLResponse struct {
	UploadURL string    `json:"uploadUrl"`
	ImageID   uuid.UUID `json:"imageId"`
	S3Key     string    `json:"s3Key"`
	ExpiresIn int       `json:"expiresIn"` // seconds
}

type DeleteProductImageRequest struct {
	ProductID uuid.UUID `param:"id" validate:"required,uuid"`
	ImageID   uuid.UUID `param:"imageId" validate:"required,uuid"`
}

func (r *DeleteProductImageRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

type CreateVariantRequest struct {
	ProductID uuid.UUID `param:"productId" validate:"required,uuid"`
	CreateVariantInput
}

func (r *CreateVariantRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

type UpdateVariantRequest struct {
	ProductID uuid.UUID `param:"productId" validate:"required,uuid"`
	VariantID uuid.UUID `param:"variantId" validate:"required,uuid"`

	// SKU           *string          `json:"sku,omitempty" validate:"omitempty,max=50"`
	Label         *string          `json:"label,omitempty" validate:"omitempty,min=2,max=100"`
	QuantityValue *decimal.Decimal `json:"quantityValue,omitempty" validate:"omitempty,gt=0"`
	QuantityUnit  *string          `json:"quantityUnit,omitempty" validate:"omitempty,max=20"`
	Price         *decimal.Decimal `json:"price,omitempty" validate:"omitempty,gt=0"`

	StockQuantity     *int  `json:"stockQuantity,omitempty" validate:"omitempty,gte=0"`
	LowStockThreshold *int  `json:"lowStockThreshold,omitempty" validate:"omitempty,gte=0"`
	IsAvailable       *bool `json:"isAvailable,omitempty"`
}

func (r *UpdateVariantRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

type DeleteVariantRequest struct {
	ProductID uuid.UUID `param:"id" validate:"required,uuid"`
	VariantID uuid.UUID `param:"variantId" validate:"required,uuid"`
}

func (r *DeleteVariantRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// admin actions
type ApproveProductRequest struct {
	ProductID uuid.UUID `param:"id" validate:"required,uuid"`
	Notes     *string   `json:"notes,omitempty" validate:"omitempty,max=500"`
}

func (r *ApproveProductRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

type RejectProductRequest struct {
	ProductID uuid.UUID `param:"id" validate:"required,uuid"`
	Reason    string    `json:"reason" validate:"required,min=10,max=500"`
	Notes     *string   `json:"notes,omitempty" validate:"omitempty,max=500"`
}

func (r *RejectProductRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

type ResubmitProductRequest struct {
	ProductID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (r *ResubmitProductRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

type CountPendingApprovalsRequest struct{}

func (r *CountPendingApprovalsRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// Responses
type ProductResponse struct {
	ID         uuid.UUID  `json:"id"`
	CompanyID  uuid.UUID  `json:"companyId"`
	CategoryID *uuid.UUID `json:"categoryId,omitempty"`

	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Unit        string          `json:"unit"`
	Origin      *string         `json:"origin,omitempty"`
	BasePrice   decimal.Decimal `json:"price"`

	ApprovalStatus  company.ApprovalStatus `json:"approvalStatus"`
	SubmittedAt     time.Time              `json:"submittedAt"`
	ReviewedByID    *uuid.UUID             `json:"reviewedById,omitempty"`
	ReviewedAt      *time.Time             `json:"reviewedAt,omitempty"`
	RejectionReason *string                `json:"rejectionReason,omitempty"`

	IsActive bool `json:"isActive"`

	Images   []ProductImageResponse   `json:"images"`
	Variants []ProductVariantResponse `json:"variants"`

	CanBeModified bool `json:"canBeModified"`
	IsVisible     bool `json:"isVisible"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ProductImageResponse struct {
	ID           uuid.UUID `json:"id"`
	ImageURL     string    `json:"imageUrl"`
	IsPrimary    bool      `json:"isPrimary"`
	DisplayOrder int       `json:"displayOrder"`
}

type ProductVariantResponse struct {
	ID        uuid.UUID `json:"id"`
	ProductID uuid.UUID `json:"productId"`

	SKU   *string `json:"sku,omitempty"`
	Label string  `json:"label"`

	QuantityValue decimal.Decimal `json:"quantityValue"`
	QuantityUnit  string          `json:"quantityUnit"`
	Price         decimal.Decimal `json:"price"`

	StockQuantity     *int `json:"stockQuantity,omitempty"`
	LowStockThreshold *int `json:"lowStockThreshold,omitempty"`
	IsLowStock        bool `json:"isLowStock"`

	IsAvailable bool `json:"isAvailable"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// this is the mapper which maps the product to product response
func ToProductResponse(
	p *Product,
	images []ProductImage,
	variants []ProductVariant,
) *ProductResponse {

	imageResponses := make([]ProductImageResponse, len(images))
	for i, img := range images {
		imageResponses[i] = ProductImageResponse{
			ID:        img.ID,
			ImageURL:  img.ImageURL,
			IsPrimary: img.IsPrimary,
			// DisplayOrder: img.DisplayOrder,
		}
	}

	variantResponses := make([]ProductVariantResponse, len(variants))
	for i, v := range variants {
		variantResponses[i] = ProductVariantResponse{
			ID:        v.ID,
			ProductID: v.ProductID,
			// SKU:               v.SKU,
			Label:             v.Label,
			QuantityValue:     v.QuantityValue,
			QuantityUnit:      v.QuantityUnit,
			Price:             v.Price,
			StockQuantity:     v.StockQuantity,
			LowStockThreshold: v.LowStockThreshold,
			IsLowStock:        v.IsLowStock(),
			IsAvailable:       v.IsAvailable,
			CreatedAt:         v.CreatedAt,
			UpdatedAt:         v.UpdatedAt,
		}
	}

	return &ProductResponse{
		ID:              p.ID,
		CompanyID:       p.CompanyID,
		CategoryID:      p.CategoryID,
		Name:            p.Name,
		Description:     p.Description,
		Unit:            p.Unit,
		Origin:          p.Origin,
		BasePrice:       p.BasePrice,
		ApprovalStatus:  p.ApprovalStatus,
		SubmittedAt:     p.SubmittedAt,
		ReviewedByID:    p.ReviewedByID,
		ReviewedAt:      p.ReviewedAt,
		RejectionReason: p.RejectionReason,
		IsActive:        p.IsActive,
		Images:          imageResponses,
		Variants:        variantResponses,
		CanBeModified:   p.CanBeModified(),
		IsVisible:       p.IsVisible(),
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

func MapProductPage(
	page *model.PaginatedResponse[Product],
	imagesMap map[uuid.UUID][]ProductImage,
	variantsMap map[uuid.UUID][]ProductVariant,
) *model.PaginatedResponse[ProductResponse] {
	responses := make([]ProductResponse, 0, len(page.Data))

	for _, p := range page.Data {
		images := imagesMap[p.ID]
		if images == nil {
			images = []ProductImage{}
		}

		variants := variantsMap[p.ID]
		if variants == nil {
			variants = []ProductVariant{}
		}

		responses = append(responses, *ToProductResponse(&p, images, variants))
	}

	return &model.PaginatedResponse[ProductResponse]{
		Data:       responses,
		Page:       page.Page,
		Limit:      page.Limit,
		Total:      page.Total,
		TotalPages: page.TotalPages,
	}
}
