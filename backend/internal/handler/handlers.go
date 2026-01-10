package handler

import (
	"github.com/C0deNe0/agromart/internal/service"
)

type Handlers struct {
	Auth    *AuthHandler
	User    *UserHandler
	Company *CompanyHandler
	Product *ProductHandler
	Health  *HealthHandler
	Upload  *UploadHandler
	Admin   *AdminHandler
}

func NewHandlers(s *service.Services) Handlers {
	return Handlers{
		Health:  NewHealthHandler(),
		User:    NewUserHandler(s.User),
		Company: NewCompanyHandler(s.Company),
		Product: NewProductHandler(s.Product),
		Auth:    NewAuthHandler(s.Auth),
		Upload:  NewUploadHandler(s.Upload),
		Admin:   NewAdminHandler(s.Company),
	}
}
