package company

import (
	"time"

	"github.com/C0deNe0/agromart/internal/model"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CreateCompanyRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=255"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=2048"`
	LogoURL     *string `json:"logoUrl,omitempty" validate:"omitempty,url"`

	BusinessEmail *string `json:"businessEmail,omitempty" validate:"omitempty,email"`
	BusinessPhone *string `json:"businessPhone,omitempty" validate:"omitempty,min=10,max=15"`

	City      *string `json:"city,omitempty" validate:"omitempty,max=100"`
	State     *string `json:"state,omitempty" validate:"omitempty,max=100"`
	Pincode   *string `json:"pincode,omitempty" validate:"omitempty,max=6"`
	GSTNumber *string `json:"gstNumber,omitempty" validate:"omitempty"`
	PANNumber *string `json:"panNumber,omitempty" validate:"omitempty"`

	ProductVisibility *ProductVisibility `json:"productVisibility,omitempty" validate:"omitempty,oneof=PUBLIC FOLLOWERS_ONLY PRIVATE"`
}

func (c *CreateCompanyRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

type UpdateCompanyRequest struct {
	ID          uuid.UUID `json:"id" validate:"required,uuid"`
	Name        *string   `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Description *string   `json:"description,omitempty" validate:"omitempty,max=2048"`
	LogoURL     *string   `json:"logoUrl,omitempty" validate:"omitempty,url"`

	BusinessEmail *string `json:"businessEmail,omitempty" validate:"omitempty,email"`
	BusinessPhone *string `json:"businessPhone,omitempty" validate:"omitempty,min=10,max=15"`

	City    *string `json:"city,omitempty" validate:"omitempty,max=100"`
	State   *string `json:"state,omitempty" validate:"omitempty,max=100"`
	Pincode *string `json:"pincode,omitempty" validate:"omitempty,max=6"`

	GSTNumber *string `json:"gstNumber,omitempty" validate:"omitempty"`
	PANNumber *string `json:"panNumber,omitempty" validate:"omitempty"`

	ProductVisibility *ProductVisibility `json:"productVisibility,omitempty" validate:"omitempty,oneof=PUBLIC FOLLOWERS_ONLY PRIVATE"`

	IsActive *bool `json:"isActive,omitempty" validate:"omitempty,oneof=true false"`
}

func (u *UpdateCompanyRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

type ListCompanyQuery struct {
	Page       int        `query:"page" validate:"min=1"`
	Limit      int        `query:"limit" validate:"min=1,max=100"`
	Search     *string    `query:"search" validate:"omitempty,max=255"`
	IsApproved *bool      `query:"isApproved"`
	IsActive   *bool      `query:"isActive"`
	OwnerID    *uuid.UUID `query:"ownerId" validate:"omitempty,uuid"`
}

func (q *ListCompanyQuery) Validate() error {
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

	return nil
}

type GetCompanyByIDRequest struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (r *GetCompanyByIDRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

type DeleteCompanyRequest struct {
	ID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (r *DeleteCompanyRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

//FOLLOW / UNFOLLOW company

type FollowCompanyRequest struct {
	CompanyID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (r *FollowCompanyRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

type UnfollowCompanyRequest struct {
	CompanyID uuid.UUID `param:"id" validate:"required,uuid"`
}

func (r *UnfollowCompanyRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

// LIST FOLLOWERS QUERY
type ListFollowersQuery struct {
	CompanyID uuid.UUID `param:"id" validate:"required,uuid"`
	Page      int       `query:"page" validate:"min=1"`
	Limit     int       `query:"limit" validate:"min=1,max=100"`
}

func (q *ListFollowersQuery) Validate() error {
	validate := validator.New()
	if err := validate.Struct(q); err != nil {
		return err
	}

	if q.Page == 0 {
		q.Page = 1
	}

	if q.Limit == 0 {
		q.Limit = 20
	}

	return nil
}

// LIST FOLLOWING QUERY
type ListFollowedCompaniesQuery struct {
	// CompanyID uuid.UUID `param:"id" validate:"required,uuid"`
	Page  int `query:"page" validate:"min=1"`
	Limit int `query:"limit" validate:"min=1,max=100"`
}

func (q *ListFollowedCompaniesQuery) Validate() error {
	validate := validator.New()
	if err := validate.Struct(q); err != nil {
		return err
	}

	if q.Page == 0 {
		q.Page = 1
	}

	if q.Limit == 0 {
		q.Limit = 20
	}

	return nil
}

//RESPONSE

type CompanyResponse struct {
	ID          uuid.UUID `json:"id"`
	OwnerID     uuid.UUID `json:"ownerId"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	LogoURL     *string   `json:"logoUrl,omitempty"`

	BusinessEmail *string `json:"businessEmail,omitempty"`
	BusinessPhone *string `json:"businessPhone,omitempty"`

	City    *string `json:"city,omitempty"`
	State   *string `json:"state,omitempty"`
	Pincode *string `json:"pincode,omitempty"`

	IsApproved bool       `json:"isApproved"`
	ApprovedBy *uuid.UUID `json:"approvedBy,omitempty"`
	ApprovedAt *time.Time `json:"approvedAt,omitempty"`
	IsActive   bool       `json:"isActive"`

	FollowerCount     int               `json:"followerCount"`
	ProductVisibility ProductVisibility `json:"productVisibility"`
	IsFollowing       *bool             `json:"isFollowing,omitempty"`
	CreatedAt         time.Time         `json:"createdAt"`
	UpdatedAt         time.Time         `json:"updatedAt"`
}

type CompanyFollowerResponse struct {
	ID         uuid.UUID `json:"id"`
	CompanyID  uuid.UUID `json:"companyId"`
	UserID     uuid.UUID `json:"userId"`
	UserName   string    `json:"userName"`
	UserEmail  string    `json:"userEmail"`
	FollowedAt time.Time `json:"followedAt"`
}

type FollowStatusResponse struct {
	CompanyID   uuid.UUID  `json:"companyId"`
	IsFollowing bool       `json:"isFollowing"`
	FollowedAt  *time.Time `json:"followedAt,omitempty"`
}

//MAPPERS

func ToCompanyResponse(c *Company, isFollowing *bool) *CompanyResponse {
	return &CompanyResponse{
		ID:                c.ID,
		OwnerID:           c.OwnerID,
		Name:              c.Name,
		Description:       c.Description,
		LogoURL:           c.LogoURL,
		BusinessEmail:     c.BusinessEmail,
		BusinessPhone:     c.BusinessPhone,
		City:              c.City,
		State:             c.State,
		Pincode:           c.Pincode,
		IsApproved:        c.IsApproved(),
		ApprovedBy:        c.ReviewedByID,
		ApprovedAt:        c.ReviewedAt,
		IsActive:          c.IsActive,
		FollowerCount:     c.FollowerCount,
		ProductVisibility: c.ProductVisibility,
		IsFollowing:       isFollowing,
		CreatedAt:         c.CreatedAt,
		UpdatedAt:         c.UpdatedAt,
	}
}

func MapCompanyPage(page *model.PaginatedResponse[Company], followStatusMap map[uuid.UUID]bool) *model.PaginatedResponse[CompanyResponse] {
	responses := make([]CompanyResponse, 0, len(page.Data))

	for _, c := range page.Data {
		isFollowing, exists := followStatusMap[c.ID]
		var isFollowingPtr *bool
		if exists {
			isFollowingPtr = &isFollowing
		}
		responses = append(responses, *ToCompanyResponse(&c, isFollowingPtr))
	}

	return &model.PaginatedResponse[CompanyResponse]{
		Data:       responses,
		Page:       page.Page,
		Limit:      page.Limit,
		Total:      page.Total,
		TotalPages: page.TotalPages,
	}
}
