package product

import (
	"errors"
	"time"

	"github.com/C0deNe0/agromart/internal/model"
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
	Price       decimal.Decimal `json:"price" validate:"required,gt=0"`
}

func (r *CreateProductRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// UPDATE PRODUCT

type UpdateProductRequest struct {
	ID          uuid.UUID        `param:"id" validate:"required,uuid"`
	CategoryID  *uuid.UUID       `json:"categoryId,omitempty" validate:"omitempty,uuid"`
	Name        *string          `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Description *string          `json:"description,omitempty" validate:"omitempty,max=2000"`
	Unit        *string          `json:"unit,omitempty" validate:"omitempty,max=50"`
	Origin      *string          `json:"origin,omitempty" validate:"omitempty,max=100"`
	Price       *decimal.Decimal `json:"price,omitempty" validate:"omitempty,gt=0"`
	IsActive    *bool            `json:"isActive,omitempty"`
	// NOTE: No CompanyID - cannot change product's company
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
	Page       int        `query:"page" validate:"min=1"`
	Limit      int        `query:"limit" validate:"min=1,max=100"`
	CompanyID  *uuid.UUID `query:"companyId" validate:"omitempty,uuid"`
	CategoryID *uuid.UUID `query:"categoryId" validate:"omitempty,uuid"`
	Search     *string    `query:"search" validate:"omitempty,min=1"`
	IsActive   *bool      `query:"isActive" validate:"omitempty"`
}

func (q *ListProductsQuery) Validate() error {
	validate := validator.New()

	if err := validate.Struct(q); err != nil {
		return err
	}

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

// RESPONSE PRODUCT

type ProductResponse struct {
	ID          uuid.UUID       `json:"id"`
	CompanyID   uuid.UUID       `json:"companyId"`
	CategoryID  *uuid.UUID      `json:"categoryId,omitempty"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	Unit        string          `json:"unit"`
	Origin      *string         `json:"origin,omitempty"`
	Price       decimal.Decimal `json:"price"`
	IsActive    bool            `json:"isActive"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
}

// Response with categroy details
type ProductWithCategoryResponse struct {
	ProductResponse
	CategoryName *string `json:"categoryName,omitempty"`
}

// this is the mapper which maps the product to product response
func ToProductResponse(p *Product) *ProductResponse {
	return &ProductResponse{
		ID:          p.ID,
		CompanyID:   p.CompanyID,
		CategoryID:  p.CategoryID,
		Name:        p.Name,
		Description: p.Description,
		Unit:        p.Unit,
		Origin:      p.Origin,
		Price:       p.Price,
		IsActive:    p.IsActive,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func MapProductPage(
	page *model.PaginatedResponse[Product],
) *model.PaginatedResponse[ProductResponse] {

	responses := make([]ProductResponse, 0, len(page.Data))
	for _, p := range page.Data {
		responses = append(responses, *ToProductResponse(&p))
	}

	return &model.PaginatedResponse[ProductResponse]{
		Data:       responses,
		Page:       page.Page,
		Limit:      page.Limit,
		Total:      page.Total,
		TotalPages: page.TotalPages,
	}
}

//PRODUCT IMAGES

type ProductImageUploadInput struct {
	ImageURL  string `json:"imageUrl" validate:"required"`
	IsPrimary bool   `json:"isPrimary" validate:"required"`
}

type DeleteProductImageRequest struct {
	ProductID uuid.UUID `param:"productId" validate:"required,uuid"`
	ImageID   uuid.UUID `param:"imageId" validate:"required,uuid"`
}
