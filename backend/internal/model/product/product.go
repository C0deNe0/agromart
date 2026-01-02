package product

import (
	"github.com/C0deNe0/agromart/internal/model"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	model.Base
	CompanyID   uuid.UUID       `json:"companyId" db:"company_id"`
	CategoryID  *uuid.UUID      `json:"categoryId,omitempty" db:"category_id"`
	Name        string          `json:"name" db:"name"`
	Description *string         `json:"description,omitempty" db:"description"`
	Price       decimal.Decimal `json:"price" db:"price"`
	Unit        string          `json:"unit" db:"unit"`
	Origin      *string         `json:"origin,omitempty" db:"origin"`
	IsActive    bool            `json:"isActive" db:"is_active"`
}

// for the product images
type ProductImage struct {
	model.Base
	ProductID uuid.UUID `json:"productId" db:"product_id"`
	ImageURL  string    `json:"imageUrl" db:"image_url"`
	IsPrimary bool      `json:"isPrimary" db:"is_primary"`
}

type ProductVariant struct {
	model.Base
	ProductID     uuid.UUID       `json:"productId" db:"product_id"`
	Label         string          `json:"label" db:"label"`
	QuantityValue decimal.Decimal `json:"quantityValue" db:"quantity_value"`
	QuantityUnit  string          `json:"quantityUnit" db:"quantity_unit"`
	Price         decimal.Decimal `json:"price" db:"price"`
	IsActive      bool            `json:"isActive" db:"is_active"`
}
