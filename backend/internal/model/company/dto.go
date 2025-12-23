package company

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CreateCompanyRequest struct {
	Name          string  `json:"name"`
	Description   *string `json:"description,omitempty"`
	LogoURL       string  `json:"logoUrl,omitempty"`
	BusinessEmail string  `json:"businessEmail,omitempty"`
	BusinessPhone string  `json:"businessPhone,omitempty"`
	City          *string `json:"city,omitempty"`
	State         *string `json:"state,omitempty"`
	Pincode       *string `json:"pincode,omitempty"`
}

type CreateCompanyInput struct {
	Name          string  `json:"name"`
	Description   *string `json:"description,omitempty"`
	LogoURL       string  `json:"logoUrl,omitempty"`
	BusinessEmail string  `json:"businessEmail,omitempty"`
	BusinessPhone string  `json:"businessPhone,omitempty"`
	City          *string `json:"city,omitempty"`
	State         *string `json:"state,omitempty"`
	Pincode       *string `json:"pincode,omitempty"`
}

func (c *CreateCompanyInput) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

type CompanyResponse struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	IsApproved bool      `json:"isApproved"`
	IsActive   bool      `json:"isActive"`
}

func ToCompanyResponse(company Company) *CompanyResponse {
	return &CompanyResponse{
		ID:         company.ID,
		Name:       company.Name,
		IsApproved: company.IsApproved,
		IsActive:   company.IsActive,
	}
}
