package product

import (
	"time"

	"github.com/C0deNe0/agromart/internal/model"
	"github.com/C0deNe0/agromart/internal/model/company"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	model.Base
	CompanyID  uuid.UUID  `json:"companyId" db:"company_id"`
	CategoryID *uuid.UUID `json:"categoryId,omitempty" db:"category_id"`

	Name        string  `json:"name" db:"name"`
	Description *string `json:"description,omitempty" db:"description"`
	Unit        string  `json:"unit" db:"unit"`
	Origin      *string `json:"origin,omitempty" db:"origin"`

	BasePrice decimal.Decimal `json:"price" db:"price"`

	ApprovalStatus company.ApprovalStatus `json:"approvalStatus" db:"approval_status"`
	SubmittedAt    time.Time              `json:"submittedAt" db:"submitted_at"`

	ReviewedByID    *uuid.UUID `json:"reviewedById,omitempty" db:"reviewed_by_id"`
	ReviewedAt      *time.Time `json:"reviewedAt,omitempty" db:"reviewed_at"`
	RejectionReason *string    `json:"rejectionReason,omitempty" db:"rejection_reason"`

	IsActive bool `json:"isActive" db:"is_active"`
	// Variants []ProductVariant `json:"variants" db:"variants"`
}

func (p *Product) IsPending() bool {
	return p.ApprovalStatus == company.ApprovalStatusPending
}

func (p *Product) IsApproved() bool {
	return p.ApprovalStatus == company.ApprovalStatusApproved
}

func (p *Product) IsRejected() bool {
	return p.ApprovalStatus == company.ApprovalStatusRejected
}

func (p *Product) CanBeModified() bool {
	return p.IsPending() || p.IsRejected()
}

func (p *Product) IsVisible() bool {
	return p.IsApproved() && p.IsActive
}

// for the product images
type ProductImage struct {
	ID        uuid.UUID `json:"id" db:"id"`
	ProductID uuid.UUID `json:"productId" db:"product_id"`
	ImageURL  string    `json:"imageUrl" db:"image_url"`
	S3Key     string    `json:"s3Key" db:"s3_key"`
	IsPrimary bool      `json:"isPrimary" db:"is_primary"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type ProductVariant struct {
	model.Base
	ProductID uuid.UUID `json:"productId" db:"product_id"`

	Label string `json:"label" db:"label"`

	QuantityValue decimal.Decimal `json:"quantityValue" db:"quantity_value"`
	QuantityUnit  string          `json:"quantityUnit" db:"quantity_unit"`

	Price decimal.Decimal `json:"price" db:"price"`

	StockQuantity     *int `json:"stockQuantity,omitempty" db:"stock_quantity"`
	LowStockThreshold *int `json:"lowStockThreshold,omitempty" db:"low_stock_threshold"`

	IsAvailable bool `json:"isAvailable" db:"is_available"`
}

// Helper methods
func (v *ProductVariant) IsLowStock() bool {
	if v.StockQuantity == nil || v.LowStockThreshold == nil {
		return false
	}
	return *v.StockQuantity <= *v.LowStockThreshold
}

type ProductApprovalHistory struct {
	ID            uuid.UUID              `json:"id" db:"id"`
	ProductID     uuid.UUID              `json:"productId" db:"product_id"`
	Action        company.ApprovalAction `json:"action" db:"action"`
	PerformedByID uuid.UUID              `json:"performedById" db:"performed_by_id"`
	Reason        *string                `json:"reason,omitempty" db:"reason"`
	Notes         *string                `json:"notes,omitempty" db:"notes"`
	CreatedAt     time.Time              `json:"createdAt" db:"created_at"`
}
