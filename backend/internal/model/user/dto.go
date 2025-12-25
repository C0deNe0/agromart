package user

import (
	"github.com/go-playground/validator/v10"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Name     string `json:"name" validate:"required,min=3,max=100"`
}

func (r *RegisterRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (l *LoginRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(l)
}
