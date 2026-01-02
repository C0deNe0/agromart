package handler

import (
	"github.com/C0deNe0/agromart/internal/service"
)

type Handlers struct {
	Health  *HealthHandler
	User    *UserHandler
	Company *CompanyHandler
	Product *ProductHandler
	Auth    *AuthHandler
	Upload  *UploadHandler
}

func NewHandlers(s *service.Services) Handlers {
	return Handlers{
		Health:  NewHealthHandler(),
		User:    NewUserHandler(s.User),
		Company: NewCompanyHandler(s.Company),
		Product: NewProductHandler(s.Product),
		Auth:    NewAuthHandler(s.Auth),
		Upload:  NewUploadHandler(s.Upload),
	}
}
