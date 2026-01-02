package company

import (
	"time"

	"github.com/C0deNe0/agromart/internal/model"
	"github.com/google/uuid"
)

type ProductVisibility string

const (
	ProductVisibilityPublic        ProductVisibility = "PUBLIC"
	ProductVisibilityFollowersOnly ProductVisibility = "FOLLOWERS_ONLY"
	ProductVisibilityPrivate       ProductVisibility = "PRIVATE"
)

type Company struct {
	model.Base
	OwnerID     uuid.UUID `json:"ownerId" db:"owner_id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description,omitempty" db:"description"`
	LogoURL     *string   `json:"logoUrl,omitempty" db:"logo_url"`

	BusinessEmail *string `json:"businessEmail,omitempty" db:"business_email"`
	BusinessPhone *string `json:"businessPhone,omitempty" db:"business_phone"`

	City    *string `json:"city,omitempty" db:"city"`
	State   *string `json:"state,omitempty" db:"state"`
	Pincode *string `json:"pincode,omitempty" db:"pincode"`

	IsApproved bool       `json:"isApproved" db:"is_approved"`
	ApprovedBy *uuid.UUID `json:"approvedBy,omitempty" db:"approved_by"`
	ApprovedAt *time.Time `json:"approvedAt,omitempty" db:"approved_at"`
	IsActive   bool       `json:"isActive" db:"is_active"`

	FollowerCount     int               `json:"followerCount" db:"follower_count"`
	ProductVisibility ProductVisibility `json:"productVisibility" db:"product_visibility"`
}

type CompanyFollower struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CompanyID  uuid.UUID `json:"companyId" db:"company_id"`
	UserID     uuid.UUID `json:"userId" db:"user_id"`
	FollowedAt time.Time `json:"followedAt" db:"followed_at"`
}
