package handlers

import (
	"net/http"

	"github.com/C0deNe0/agromart/internal/model/user"
	"github.com/C0deNe0/agromart/internal/service"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register() echo.HandlerFunc {
	return Handle(
		&user.RegisterRequest{},
		func(c echo.Context, req *user.RegisterRequest) (*user.AuthResponse, error) {
			return h.authService.RegisterWithEmail(c.Request().Context(), req.Email, req.Password, req.Name)
		},
		http.StatusCreated,
	)
}

func (h *AuthHandler) Login() echo.HandlerFunc {
	return Handle(
		&user.LoginRequest{},
		func(c echo.Context, req *user.LoginRequest) (*user.AuthResponse, error) {
			return h.authService.LoginWithEmail(c.Request().Context(),
				req.Email,
				req.Password,
			)
		},
		http.StatusOK,
	)

}


//REGISTER----------

// HTTP → Handler → AuthService
//      → UserRepo.Create
//      → UserAuthMethodRepo.Create
//      → TokenManager.Generate
//      → 


// LOGIN------------

// HTTP → Handler → AuthService
//      → UserAuthMethodRepo.GetByEmail
//      → VerifyPassword
//      → TokenManager.Generate
//      → Response