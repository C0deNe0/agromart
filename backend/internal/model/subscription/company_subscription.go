package subscription

import (
	"time"

	"github.com/C0deNe0/agromart/internal/model"
	"github.com/google/uuid"
)

type SubscriptionStatus string

const (
	SubActive    SubscriptionStatus = "ACTIVE"
	SubInactive  SubscriptionStatus = "INACTIVE"
	SubExpired   SubscriptionStatus = "EXPIRED"
	SubCancelled SubscriptionStatus = "CANCELLED"
)

type CompanySubscription struct {
	model.Base

	CompanyID uuid.UUID          `json:"companyId" db:"company_id"`
	PlanID    uuid.UUID          `json:"planId" db:"plan_id"`
	Status    SubscriptionStatus `json:"status" db:"status"`
	StartDate time.Time          `json:"startDate" db:"start_date"`
	EndDate   *time.Time         `json:"endDate" db:"end_date"`
}

// FREE plan:
//   max_product_images = 3

// User uploads image:
//   → fetch company subscription
//   → check limit
//   → allow / reject
