package product

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

//CREATE PRODUCT

type CreateProductRequest struct {
	Name         string     `json:"name" validate:"required,min=3,max=255"`
	Description  *string    `json:"description,omitempty" validate:"omitempty,max=2000"`
	CategoryID   *uuid.UUID `json:"categoryId,omitempty" validate:"omitempty,uuid"`
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
	ID           uuid.UUID  `json:"id"`
	CompanyID    uuid.UUID  `json:"companyId"`
	Name         string     `json:"name"`
	Description  *string    `json:"description,omitempty"`
	CategoryID   *uuid.UUID `json:"categoryId,omitempty"`
	Unit         *string    `json:"unit,omitempty"`
	Origin       *string    `json:"origin,omitempty"`
	PriceDisplay *string    `json:"priceDisplay,omitempty"`
	IsActive     bool       `json:"isActive"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}
