package auth

import "github.com/go-playground/validator/v10"

type GoogleIDTokenRequest struct {
	// The ID token obtained from the native Google Sign-In SDK on the mobile device
	IDToken string `json:"id_token" validate:"required"`
}

func (g *GoogleIDTokenRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(g)
}
