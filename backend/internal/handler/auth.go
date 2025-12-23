package handler

import (
	"net/http"

	"github.com/C0deNe0/agromart/internal/lib/utils"
	"github.com/C0deNe0/agromart/internal/model/user"
	"github.com/C0deNe0/agromart/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// OAuth is not business logic, and not HTTP logic either.
// It is an external integration dependency.

// So it should be:
// ------Created once
// ------Injected
// ------Not recreated per request

type AuthHandler struct {
	authService *service.AuthService
	googleOAuth *utils.GoogleOAuth
}

func NewAuthHandler(authService *service.AuthService, googleOAuth *utils.GoogleOAuth) *AuthHandler {
	return &AuthHandler{authService: authService, googleOAuth: googleOAuth}
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

func (h *AuthHandler) GoogleLogin() echo.HandlerFunc {
	return func(c echo.Context) error {
		state := uuid.New().String()
		return c.Redirect(
			http.StatusTemporaryRedirect,
			h.googleOAuth.AuthURL(state),
		)
	}
}

func (h *AuthHandler) GoogleCallback() echo.HandlerFunc {
	return func(c echo.Context) error {
		code := c.QueryParam("code")
		if code == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "code is required")
		}
		token, err := h.googleOAuth.Exchange(c.Request().Context(), code)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		googleUser, err := utils.FetchGoogleUser(token.AccessToken)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		resp, err := h.authService.LoginWithGoogle(c.Request().Context(),
			googleUser.Sub,
			googleUser.Email,
			googleUser.Name,
			&googleUser.Picture,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, resp)

	}
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

// GOOGLE AUTH
// HTTP → Handler → AuthService
//      → GoogleOAuth.Exchange
//      → FetchGoogleUser
//      → AuthService.LoginWithGoogle
//      → TokenManager.Generate
//      → Response
