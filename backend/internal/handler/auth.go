package handler

import (
	"net/http"

	"github.com/C0deNe0/agromart/internal/lib/utils"
	"github.com/C0deNe0/agromart/internal/model/auth"
	"github.com/C0deNe0/agromart/internal/model/user"
	"github.com/C0deNe0/agromart/internal/service"
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

func (h *AuthHandler) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing refresh token")
		}

		const prefix = "Bearer "
		if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header")
		}

		rawRefreshToken := authHeader[len(prefix):]
		if err := h.authService.Logout(c.Request().Context(), rawRefreshToken); err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
		return c.JSON(http.StatusOK, echo.Map{
			"message": "logged out successfully",
		})
	}
}

func (h *AuthHandler) Refresh() echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing refresh token")
		}

		const prefix = "Bearer "
		if len(authHeader) <= len(prefix) || authHeader[:len(prefix)] != prefix {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header")
		}

		rawRefreshToken := authHeader[len(prefix):]

		resp, err := h.authService.Refresh(
			c.Request().Context(),
			rawRefreshToken,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}

		return c.JSON(http.StatusOK, resp)
	}
}

// package handler...

func (h *AuthHandler) LoginWithGoogleIDToken() echo.HandlerFunc {
	return Handle(
		&auth.GoogleIDTokenRequest{},
		func(c echo.Context, req *auth.GoogleIDTokenRequest) (*user.AuthResponse, error) {

			googleUserClaims, err := utils.VerifyGoogleIDToken(c.Request().Context(), req.IDToken)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, "Invalid Google ID Token"+err.Error())
			}

			// 2. Call the service layer to handle login/registration using the verified claims
			resp, err := h.authService.LoginWithGoogle(
				c.Request().Context(),
				googleUserClaims.Sub,
				googleUserClaims.Email,
				googleUserClaims.Name,
				&googleUserClaims.Picture,
			)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
			}

			return resp, nil
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

// GOOGLE AUTH
// HTTP → Handler → AuthService
//      → GoogleOAuth.Exchange
//      → FetchGoogleUser
//      → AuthService.LoginWithGoogle
//      → TokenManager.Generate
//      → Response

// func (h *AuthHandler) GoogleLogin() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		state := uuid.New().String()

// 		//store state in HTTP only cookie for 5 min

// 		c.SetCookie(&http.Cookie{
// 			Name:     "oauth_state",
// 			Value:    state,
// 			Path:     "/",
// 			HttpOnly: true,
// 			Secure:   true,
// 			MaxAge:   300,
// 			// Expires:  time.Now().Add(5 * time.Minute),
// 			SameSite: http.SameSiteLaxMode,
// 		})

// 		return c.Redirect(
// 			http.StatusTemporaryRedirect,
// 			h.googleOAuth.AuthURL(state),
// 		)
// 	}
// }

// func (h *AuthHandler) GoogleCallback() echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		state := c.QueryParam("state")
// 		if state == "" {
// 			return echo.NewHTTPError(http.StatusBadRequest, "state is required")
// 		}
// 		code := c.QueryParam("code")
// 		if code == "" {
// 			return echo.NewHTTPError(http.StatusBadRequest, "code is required")
// 		}

// 		//validate state
// 		cookie, err := c.Cookie("oauth_state")
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusBadRequest, "state is required")
// 		}
// 		if cookie.Value != state {
// 			return echo.NewHTTPError(http.StatusBadRequest, "state does not match")
// 		}

// 		//clearing after use
// 		c.SetCookie(&http.Cookie{
// 			Name:     "oauth_state",
// 			Value:    "",
// 			Path:     "/",
// 			HttpOnly: true,
// 			Secure:   true,
// 			MaxAge:   -1,
// 			SameSite: http.SameSiteLaxMode,
// 		})

// 		token, err := h.googleOAuth.Exchange(c.Request().Context(), code)
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
// 		}

// 		googleUser, err := utils.FetchGoogleUser(token.AccessToken)
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
// 		}

// 		resp, err := h.authService.LoginWithGoogle(c.Request().Context(),
// 			googleUser.Sub,
// 			googleUser.Email,
// 			googleUser.Name,
// 			&googleUser.Picture,
// 		)
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
// 		}

// 		return c.JSON(http.StatusOK, resp)

// 	}
// }
