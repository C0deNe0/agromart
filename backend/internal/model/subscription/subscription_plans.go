package subscription

import "github.com/C0deNe0/agromart/internal/model"

type BillingCycle string

const (
	BillingCycleMonthly BillingCycle = "MONTHLY"
	BillingCycleYearly  BillingCycle = "YEARLY"
)

type SubscriptionPlan struct {
	model.Base
	Name                 string       `json:"name" db:"name"`
	Price                string       `json:"price" db:"price"`
	BillingCycle         BillingCycle `json:"billingCycle" db:"billing_cycle"`
	MaxProducts          *int         `json:"maxProducts,omitempty" db:"max_products"`
	MaxProductImages     *int         `json:"maxProductImages,omitempty" db:"max_product_images"`
	MaxVariantPerProduct *int         `json:"maxVariantPerProduct,omitempty" db:"max_variant_per_product"`
	IsActive             bool         `json:"isActive" db:"is_active"`
}
