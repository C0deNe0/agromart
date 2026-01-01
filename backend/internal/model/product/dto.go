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
	Name         string     `json:"name" validate:"required,min=3,max=255"`
	Description  *string    `json:"description,omitempty" validate:"omitempty,max=2000"`
	CategoryID   *uuid.UUID `json:"categoryId,omitempty" validate:"omitempty,uuid"`
	CompanyID    *uuid.UUID `json:"companyId,omitempty" validate:"omitempty,uuid"`
	Unit         *string    `json:"unit,omitempty" validate:"omitempty,max=50"`
	Origin       *string    `json:"origin,omitempty" validate:"omitempty,max=100"`
	PriceDisplay *string    `json:"priceDisplay,omitempty" validate:"omitempty,max=255"`
}

func (r *CreateProductRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// UPDATE PRODUCT

type UpdateProductRequest struct {
	ID           uuid.UUID  `param:"id" validate:"required,uuid"`
	Name         *string    `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Description  *string    `json:"description,omitempty" validate:"omitempty,max=2000"`
	CategoryID   *uuid.UUID `json:"categoryId,omitempty" validate:"omitempty,uuid"`
	CompanyID    *uuid.UUID `json:"companyId,omitempty" validate:"omitempty,uuid"`
	Unit         *string    `json:"unit,omitempty" validate:"omitempty,max=50"`
	Origin       *string    `json:"origin,omitempty" validate:"omitempty,max=100"`
	PriceDisplay *string    `json:"priceDisplay,omitempty" validate:"omitempty,max=255"`

	IsActive *bool `json:"isActive,omitempty" validate:"omitempty"`
}

func (r *UpdateProductRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// LIST PRODUCT

type ListProductsQuery struct {
	Page       *int       `query:"page" validate:"required,min=1"`
	Limit      *int       `query:"limit" validate:"required,min=1"`
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

	if q.Page == nil {
		p := 1
		q.Page = &p
	}

	if q.Limit == nil {
		l := 10
		q.Limit = &l
	}

	if *q.Page < 1 {
		return errors.New("page must be greater than 0")
	}

	if *q.Limit < 1 {
		return errors.New("limit must be greater than 0")
	}

	return nil
}

// RESPONSE PRODUCT

type ProductResponse struct {
	ID          uuid.UUID       `json:"id"`
	CompanyID   uuid.UUID       `json:"companyId"`
	Name        string          `json:"name"`
	Description *string         `json:"description,omitempty"`
	CategoryID  *uuid.UUID      `json:"categoryId,omitempty"`
	Unit        string          `json:"unit"`
	Origin      *string         `json:"origin,omitempty"`
	Price       decimal.Decimal `json:"price"`
	IsActive    bool            `json:"isActive"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
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

type ProductCreateInput struct {
	CompanyID   uuid.UUID
	Name        string
	Description *string
	CategoryID  *uuid.UUID
	Unit        *string
	Origin      *string
	Price       decimal.Decimal
}

func (p *ProductCreateInput) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type ProductUpdateInput struct {
	ID          uuid.UUID
	Name        *string
	Description *string
	CategoryID  *uuid.UUID
	Unit        *string
	Origin      *string
	Price       decimal.Decimal
	IsActive    *bool
}

func (p *ProductUpdateInput) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

// this is the mapper which maps the product to product response
func ToProductResponse(p *Product) *ProductResponse {
	return &ProductResponse{
		ID:          p.ID,
		CompanyID:   p.CompanyID,
		Name:        p.Name,
		Description: p.Description,
		Unit:        p.Unit,
		Origin:      p.Origin,
		Price:       p.Price,
		IsActive:    p.IsActive,
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

type ProductImageUploadInput struct {
	ImageURL  string `json:"imageUrl" validate:"required"`
	IsPrimary bool   `json:"isPrimary" validate:"required"`
}

type DeleteProductImageRequest struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}
