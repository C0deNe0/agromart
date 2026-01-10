package company

import (
	"time"

	"github.com/C0deNe0/agromart/internal/model"
	"github.com/google/uuid"
)

type ApprovalStatus string

const (
	ApprovalStatusPending  ApprovalStatus = "PENDING"
	ApprovalStatusApproved ApprovalStatus = "APPROVED"
	ApprovalStatusRejected ApprovalStatus = "REJECTED"
)

type ProductVisibility string

const (
	ProductVisibilityPublic        ProductVisibility = "PUBLIC"
	ProductVisibilityFollowersOnly ProductVisibility = "FOLLOWERS_ONLY"
	ProductVisibilityPrivate       ProductVisibility = "PRIVATE"
)

type Company struct {
	model.Base
	OwnerID uuid.UUID `json:"ownerId" db:"owner_id"`

	Name        string  `json:"name" db:"name"`
	Description *string `json:"description,omitempty" db:"description"`
	LogoURL     *string `json:"logoUrl,omitempty" db:"logo_url"`

	BusinessEmail *string `json:"businessEmail,omitempty" db:"business_email"`
	BusinessPhone *string `json:"businessPhone,omitempty" db:"business_phone"`

	City    *string `json:"city,omitempty" db:"city"`
	State   *string `json:"state,omitempty" db:"state"`
	Pincode *string `json:"pincode,omitempty" db:"pincode"`

	GSTNumber *string `json:"gstNumber,omitempty" db:"gst_number"`
	PANNumber *string `json:"panNumber,omitempty" db:"pan_number"`

	// Approval Workflow
	ApprovalStatus ApprovalStatus `json:"approvalStatus" db:"approval_status"`
	SubmittedAt    time.Time      `json:"submittedAt" db:"submitted_at"`

	ReviewedByID    *uuid.UUID `json:"reviewedById,omitempty" db:"reviewed_by_id"`
	ReviewedAt      *time.Time `json:"reviewedAt,omitempty" db:"reviewed_at"`
	RejectionReason *string    `json:"rejectionReason,omitempty" db:"rejection_reason"`

	IsActive bool `json:"isActive" db:"is_active"`

	FollowerCount     int               `json:"followerCount" db:"follower_count"`
	ProductVisibility ProductVisibility `json:"productVisibility" db:"product_visibility"`
}

func (c *Company) IsPending() bool {
	return c.ApprovalStatus == ApprovalStatusPending
}

func (c *Company) IsApproved() bool {
	return c.ApprovalStatus == ApprovalStatusApproved
}

func (c *Company) IsRejected() bool {
	return c.ApprovalStatus == ApprovalStatusRejected
}
func (c *Company) CanCreateProducts() bool {
	return c.IsApproved() && c.IsActive
}

func (c *Company) CanBeFollowed() bool {
	return c.IsApproved() && c.IsActive
}

func (c *Company) CanBeModified() bool {
	// Can only modify PENDING or REJECTED companies
	return c.IsPending() || c.IsRejected()
}

// COMPANY FOLLOWER MODEL

type CompanyFollower struct {
	ID         uuid.UUID `json:"id" db:"id"`
	CompanyID  uuid.UUID `json:"companyId" db:"company_id"`
	UserID     uuid.UUID `json:"userId" db:"user_id"`
	FollowedAt time.Time `json:"followedAt" db:"followed_at"`
}

type ApprovalAction string

const (
	ActionSubmitted   ApprovalAction = "SUBMITTED"
	ActionApproved    ApprovalAction = "APPROVED"
	ActionRejected    ApprovalAction = "REJECTED"
	ActionResubmitted ApprovalAction = "RESUBMITTED"
)

type CompanyApprovalHistory struct {
	ID            uuid.UUID      `json:"id" db:"id"`
	CompanyID     uuid.UUID      `json:"companyId" db:"company_id"`
	Action        ApprovalAction `json:"action" db:"action"`
	PerformedByID uuid.UUID      `json:"performedById" db:"performed_by_id"`
	Reason        *string        `json:"reason,omitempty" db:"reason"`
	Notes         *string        `json:"notes,omitempty" db:"notes"`
	CreatedAt     time.Time      `json:"createdAt" db:"created_at"`
}
 

